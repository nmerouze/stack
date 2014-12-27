package jsonapi

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/context"
	"github.com/justinas/alice"
	"github.com/nmerouze/stack/mux"
)

// New() returns a new mux with handlers useful to create JSON APIs.
func New() *mux.Mux {
	m := mux.New()
	m.Chain = alice.New(RecoverHandler, LoggingHandler, AcceptHandler)
	return m
}

// Err contains more information than standard go errors to create useful error messages for API consumers.
type Err struct {
	Id     string `json:"id"`
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// Err implements the error interface.
func (err Err) Error() string {
	return err.Detail
}

var (
	ErrBadRequest           = &Err{"bad_request", 400, "Bad request", "Request body is not well-formed. It must be JSON."}
	ErrUnauthorized         = &Err{"unauthorized", 401, "Unauthorized", "Access token is invalid."}
	ErrNotFound             = &Err{"not_found", 404, "Not found", "Route not found."}
	ErrNotAcceptable        = &Err{"not_acceptable", 406, "Not acceptable", "Accept HTTP header must be \"application/vnd.api+json\"."}
	ErrUnsupportedMediaType = &Err{"unsupported_media_type", 415, "Unsupported Media Type", "Content-Type header must be \"application/vnd.api+json\"."}
	ErrInternalServer       = &Err{"internal_server_error", 500, "Internal Server Error", "Something went wrong."}
)

// Encodes the response into JSON and sets the adequate Content-Type HTTP header.
func Write(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(v)
}

// Encodes the error into JSON, sets the adequate Content-Type HTTP header and set the status code related to the error.
func Error(w http.ResponseWriter, err *Err) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(map[string][]*Err{"errors": []*Err{err}})
}

// Returns an error if the Content-Type HTTP header is not "application/vnd.api+json".
func ContentTypeHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/vnd.api+json" {
			Error(w, ErrUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// Returns an error if the Accept HTTP header is not "application/vnd.api+json".
func AcceptHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/vnd.api+json" {
			Error(w, ErrNotAcceptable)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// Decodes a request body into the struct passed to the middleware.
// If the request body is not JSON, it will return a 400 Bad Request error.
// Stores the decoded body into a context object.
func BodyHandler(v interface{}) func(http.Handler) http.Handler {
	t := reflect.TypeOf(v)

	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			val := reflect.New(t).Interface()
			err := json.NewDecoder(r.Body).Decode(val)

			if err != nil {
				Error(w, ErrBadRequest)
				return
			}

			if next != nil {
				context.Set(r, "body", val)
				next.ServeHTTP(w, r)
			}
		}

		return http.HandlerFunc(fn)
	}

	return m
}

// Body(r *http.Request) is a function to get the decoded body from the request context
func Body(r *http.Request) interface{} {
	return context.Get(r, "body")
}

// If the code panics, it logs the error and returns a 500 Internal Server Error error.
func RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				Error(w, ErrInternalServer)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// Logs every request.
func LoggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}
