package stack

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tagMiddleware(tag string) Middleware {
	return func(c *C, h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(tag))
			h.ServeHTTP(w, r)
		})
	}
}

var testApp = HandlerFunc(func(c C, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("app\n"))
})

var testContextApp = HandlerFunc(func(c C, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(c.Env["foo"].(string) + "\n"))
})

func TestNew(t *testing.T) {
	h1 := func(c *C, h http.Handler) http.Handler {
		return nil
	}

	h2 := func(c *C, h http.Handler) http.Handler {
		return http.StripPrefix("potato", nil)
	}

	middlewares := []Middleware{h1, h2}
	stack := New(middlewares...)

	assert.Equal(t, stack.middlewares[0], middlewares[0])
	assert.Equal(t, stack.middlewares[1], middlewares[1])
}

func TestThen(t *testing.T) {
	t1 := tagMiddleware("t1\n")
	t2 := tagMiddleware("t2\n")
	t3 := tagMiddleware("t3\n")

	stacked := New(t1, t2, t3).Then(testApp)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	stacked.ServeHTTP(w, r)

	assert.Equal(t, w.Body.String(), "t1\nt2\nt3\napp\n")
}

func TestThenWithContext(t *testing.T) {
	h1 := func(c *C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			c.Env["foo"] = "bar"
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}

	stacked := New(h1).Then(testContextApp)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	stacked.ServeHTTP(w, r)

	assert.Equal(t, w.Body.String(), "bar\n")
}

func TestThenWorksWithNoMiddleware(t *testing.T) {
	stacked := New().Then(testApp)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	stacked.ServeHTTP(w, r)

	assert.Equal(t, w.Body.String(), "app\n")
}

func TestThenTreatsNilAsDefaultServeMux(t *testing.T) {
	stacked := New().Then(nil)
	assert.Equal(t, stacked, http.DefaultServeMux)
}

func TestThenFuncTreatsNilAsDefaultServeMux(t *testing.T) {
	stacked := New().ThenFunc(nil)
	assert.Equal(t, stacked, http.DefaultServeMux)
}

func TestAppendAddsHandlersCorrectly(t *testing.T) {
	stack := New(tagMiddleware("t1\n"), tagMiddleware("t2\n"))
	newStack := stack.Append(tagMiddleware("t3\n"), tagMiddleware("t4\n"))

	assert.Equal(t, len(stack.middlewares), 2)
	assert.Equal(t, len(newStack.middlewares), 4)

	stacked := newStack.Then(testApp)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	stacked.ServeHTTP(w, r)

	assert.Equal(t, w.Body.String(), "t1\nt2\nt3\nt4\napp\n")
}

func TestAppendRespectsImmutability(t *testing.T) {
	stack := New(tagMiddleware(""))
	newStack := stack.Append(tagMiddleware(""))

	assert.NotEqual(t, &stack.middlewares[0], &newStack.middlewares[0])
}
