package mux

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func New() *Mux {
	return &Mux{Router: httprouter.New()}
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		context.Set(r, "params", ps)
		context.ClearHandler(h).ServeHTTP(w, r)
	}
}

type Mux struct {
	Router *httprouter.Router
	Chain  alice.Chain
}

func (m *Mux) Use(middlewares ...alice.Constructor) {
	m.Chain = m.Chain.Append(middlewares...)
}

func (m *Mux) Get(p string) *route {
	return &route{mux: m, method: "GET", pattern: p, chain: m.Chain}
}

func (m *Mux) Head(p string) *route {
	return &route{mux: m, method: "HEAD", pattern: p, chain: m.Chain}
}

func (m *Mux) Post(p string) *route {
	return &route{mux: m, method: "POST", pattern: p, chain: m.Chain}
}

func (m *Mux) Patch(p string) *route {
	return &route{mux: m, method: "PATCH", pattern: p, chain: m.Chain}
}

func (m *Mux) Put(p string) *route {
	return &route{mux: m, method: "PUT", pattern: p, chain: m.Chain}
}

func (m *Mux) Delete(p string) *route {
	return &route{mux: m, method: "DELETE", pattern: p, chain: m.Chain}
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Router.ServeHTTP(w, r)
}

type route struct {
	mux     *Mux
	method  string
	pattern string
	chain   alice.Chain
}

func (r *route) Use(middlewares ...alice.Constructor) *route {
	r.chain = r.chain.Append(middlewares...)
	return r
}

func (r *route) Then(h http.Handler) {
	r.mux.Router.Handle(r.method, r.pattern, wrapHandler(r.chain.Then(h)))
}

func (r *route) ThenFunc(f http.HandlerFunc) {
	r.mux.Router.Handle(r.method, r.pattern, wrapHandler(r.chain.ThenFunc(f)))
}

// Params(r *http.Request) is a function to get URL params from the request context
func Params(r *http.Request) httprouter.Params {
	return context.Get(r, "params").(httprouter.Params)
}
