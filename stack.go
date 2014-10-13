// Package stack provides a convenient way to chain http handlers with contexts.
package stack

import "net/http"

type C struct {
	Env map[string]interface{}
}

type Middleware func(c *C, h http.Handler) http.Handler

type Handler interface {
	ServeHTTP(C, http.ResponseWriter, *http.Request)
}

type HandlerFunc func(C, http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(c C, w http.ResponseWriter, r *http.Request) {
	f(c, w, r)
}

// Stack acts as a list of http.Handler middlewares.
// Stack is effectively immutable:
// once created, it will always hold
// the same set of middlewares in the same order.
type Stack struct {
	middlewares []Middleware
}

// New creates a new chain,
// memorizing the given list of middlewares.
// New serves no other function,
// middlewares are only called upon a call to Then().
func New(middlewares ...Middleware) Stack {
	s := Stack{}
	s.middlewares = append(s.middlewares, middlewares...)

	return s
}

// Then stacks the middleware and returns the final http.Handler.
//     New(m1, m2, m3).Then(h)
// is equivalent to:
//     m1(m2(m3(h)))
// When the request comes in, it will be passed to m1, then m2, then m3
// and finally, the given handler
// (assuming every middleware calls the following one).
//
// A stack can be safely reused by calling Then() several times.
//     stdStack := stack.New(ratelimitHandler, csrfHandler)
//     indexPipe = stdStack.Then(indexHandler)
//     authPipe = stdStack.Then(authHandler)
// Note that middlewares are called on every call to Then()
// and thus several instances of the same middleware will be created
// when a stack is reused in this way.
// For proper middleware, this should cause no problems.
//
// Then() treats nil as http.DefaultServeMux.
func (s Stack) Then(h Handler) http.Handler {
	c := C{}
	c.Env = map[string]interface{}{}

	var final http.Handler

	if h != nil {
		final = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(c, w, r)
		})
	} else {
		final = http.DefaultServeMux
	}

	for i := len(s.middlewares) - 1; i >= 0; i-- {
		final = s.middlewares[i](&c, final)
	}

	return final
}

// ThenFunc works identically to Then, but takes
// a HandlerFunc instead of a Handler.
//
// The following two statements are equivalent:
//     s.Then(stack.HandlerFunc(fn))
//     s.ThenFunc(fn)
//
// ThenFunc provides all the guarantees of Then.
func (s Stack) ThenFunc(fn HandlerFunc) http.Handler {
	if fn == nil {
		return s.Then(nil)
	}
	return s.Then(HandlerFunc(fn))
}

// Append extends a stack, adding the specified middlewares
// as the last ones in the request flow.
//
// Append returns a new stack, leaving the original one untouched.
//
//     stdStack := stack.New(m1, m2)
//     extStack := stdStack.Append(m3, m4)
//     // requests in stdStack go m1 -> m2
//     // requests in extStack go m1 -> m2 -> m3 -> m4
func (s Stack) Append(middlewares ...Middleware) Stack {
	newMiddlewares := make([]Middleware, len(s.middlewares))
	copy(newMiddlewares, s.middlewares)
	newMiddlewares = append(newMiddlewares, middlewares...)

	newStack := New(newMiddlewares...)
	return newStack
}
