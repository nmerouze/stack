package mux_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nmerouze/stack/mux"
)

func testMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test\n"))
		h.ServeHTTP(w, r)
	})
}

func TestRouter(t *testing.T) {
	m := mux.New()
	m.Get("/teas").ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sencha\n"))
	})

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/teas", nil)
	if err != nil {
		t.Fatal(err)
	}

	m.ServeHTTP(w, r)

	if w.Body.String() != "sencha\n" {
		t.Fatalf("response body expected: %#v, got: %#v", "sencha\n", w.Body.String())
	}
}

func TestRoute(t *testing.T) {
	m := mux.New()
	m.Get("/teas").Use(testMiddleware).Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sencha\n"))
	}))

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/teas", nil)
	if err != nil {
		t.Fatal(err)
	}

	m.ServeHTTP(w, r)

	if w.Body.String() != "test\nsencha\n" {
		t.Fatalf("response body expected: %#v, got: %#v", "test\nsencha\n", w.Body.String())
	}
}

func TestMiddlewares(t *testing.T) {
	m := mux.New()
	m.Use(testMiddleware)

	m.Get("/teas").ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sencha\n"))
	})

	m.Get("/cars").ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("audi\n"))
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teas", nil)
	m.ServeHTTP(w, r)
	if w.Body.String() != "test\nsencha\n" {
		t.Fatalf("response body expected: %#v, got: %#v", "test\nsencha\n", w.Body.String())
	}

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/cars", nil)
	m.ServeHTTP(w, r)
	if w.Body.String() != "test\naudi\n" {
		t.Fatalf("response body expected: %#v, got: %#v", "test\naudi\n", w.Body.String())
	}
}

func TestParams(t *testing.T) {
	m := mux.New()
	m.Get("/teas/:id").ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		p := mux.Params(r)
		fmt.Fprintf(w, "%s\n", p.ByName("id"))
	})

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/teas/hojicha", nil)
	if err != nil {
		t.Fatal(err)
	}

	m.ServeHTTP(w, r)

	if w.Body.String() != "hojicha\n" {
		t.Fatalf("response body expected: %#v, got: %#v", "hojicha\n", w.Body.String())
	}
}

func TestRoutes(t *testing.T) {
	var get, head, post, patch, put, delete bool

	m := mux.New()
	m.Get("/get").ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		get = true
	})

	m.Head("/head").ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		head = true
	})

	m.Post("/post").ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		post = true
	})

	m.Patch("/patch").ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		patch = true
	})

	m.Put("/put").ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		put = true
	})

	m.Delete("/delete").ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		delete = true
	})

	w := httptest.NewRecorder()

	r, _ := http.NewRequest("GET", "/get", nil)
	m.ServeHTTP(w, r)
	if !get {
		t.Fatalf("routing GET failed")
	}

	r, _ = http.NewRequest("HEAD", "/head", nil)
	m.ServeHTTP(w, r)
	if !head {
		t.Fatalf("routing HEAD failed")
	}

	r, _ = http.NewRequest("POST", "/post", nil)
	m.ServeHTTP(w, r)
	if !post {
		t.Fatalf("routing POST failed")
	}

	r, _ = http.NewRequest("PATCH", "/patch", nil)
	m.ServeHTTP(w, r)
	if !patch {
		t.Fatalf("routing PATCH failed")
	}

	r, _ = http.NewRequest("PUT", "/put", nil)
	m.ServeHTTP(w, r)
	if !put {
		t.Fatalf("routing PUT failed")
	}

	r, _ = http.NewRequest("DELETE", "/delete", nil)
	m.ServeHTTP(w, r)
	if !delete {
		t.Fatalf("routing DELETE failed")
	}
}
