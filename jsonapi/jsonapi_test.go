package jsonapi_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/context"
	"github.com/nmerouze/stack/jsonapi"
)

func TestWrite(t *testing.T) {
	w := httptest.NewRecorder()
	r := map[string]interface{}{"id": "123456", "name": "sencha"}
	jsonapi.Write(w, r)

	if w.Body.String() != "{\"id\":\"123456\",\"name\":\"sencha\"}\n" {
		t.Fatalf("Response body expected: %#v, got: %#v", "{\"id\":\"123456\",\"name\":\"sencha\"}\n", w.Body.String())
	}

	if w.Header().Get("Content-Type") != "application/vnd.api+json" {
		t.Fatalf(`Content-type "application/vnd.api+json" expected, got "%#v"`, w.Header().Get("Content-Type"))
	}
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	err := &jsonapi.Err{Id: "bad_request", Status: 400, Title: "Bad Request", Detail: "Bad Request"}
	jsonapi.Error(w, err)

	expect := "{\"errors\":[{\"id\":\"bad_request\",\"status\":400,\"title\":\"Bad Request\",\"detail\":\"Bad Request\"}]}\n"

	if w.Code != 400 {
		t.Fatalf("Response status expected: %#v, got: %#v", 400, w.Code)
	}

	if w.Body.String() != expect {
		t.Fatalf("Response body expected: %#v, got: %#v", expect, w.Body.String())
	}
}

func TestContentTypeHandler(t *testing.T) {
	h := jsonapi.ContentTypeHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sencha\n"))
	}))

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, r)
	if w.Code != 415 {
		t.Fatalf("Response status expected: %#v, got: %#v", 415, w.Code)
	}

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/", nil)
	r.Header.Set("Content-Type", "application/vnd.api+json")
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("Response status expected: %#v, got: %#v", 200, w.Code)
	}
}

func TestAcceptHandler(t *testing.T) {
	h := jsonapi.AcceptHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sencha\n"))
	}))

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, r)
	if w.Code != 406 {
		t.Fatalf("Response status expected: %#v, got: %#v", 406, w.Code)
	}

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/", nil)
	r.Header.Set("Accept", "application/vnd.api+json")
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("Response status expected: %#v, got: %#v", 200, w.Code)
	}
}

type testBody struct {
	Data struct {
		Title string `json:"title"`
	} `json:"data"`
}

func TestBodyHandler(t *testing.T) {
	h := jsonapi.BodyHandler(testBody{})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := context.Get(r, "body").(*testBody)
		fmt.Fprintf(w, "%s\n", c.Data.Title)
	}))

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/", bytes.NewBufferString(``))
	h.ServeHTTP(w, r)
	if w.Code != 400 {
		t.Fatalf("Response status expected: %#v, got: %#v", 400, w.Code)
	}

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("POST", "/", bytes.NewBufferString(`{"data":{"title":"foobar"}}`))
	h.ServeHTTP(w, r)
	if w.Body.String() != "foobar\n" {
		t.Fatalf("Response body expected: %#v, got: %#v", "foobar\n", w.Body.String())
	}
}
