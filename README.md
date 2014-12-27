# stack [![GoDoc](https://godoc.org/github.com/nmerouze/stack?status.png)](https://godoc.org/github.com/nmerouze/stack) [![Build Status](https://travis-ci.org/nmerouze/stack.svg?branch=master)](https://travis-ci.org/nmerouze/stack)


stack is a framework to build [JSON-APIs](http://jsonapi.org) faster. It is based on the series of articles ["Build Your Own Web Framework in Go"](http://nicolasmerouze.com/build-web-framework-golang/) I began to write a few weeks ago.

The public API of the package is stable and you can use it right now to make your application. I don't indend to break any existing feature but will add new features to make the framework more useful for production. [Look at the documentation](http://godoc.org/github.com/nmerouze/stack/jsonapi).

I am also developing new ideas that could become new features in the [experimental branch](https://github.com/nmerouze/stack/tree/experimental). There are some exciting stuff in it, so check it out!

# Getting started

``` go
package main

import (
  "net/http"
  "github.com/nmerouze/stack/jsonapi"
)

type Tea struct {
  Name string `json:"name"`
}

type TeaCollection struct {
  Data []Tea `json:"data"`  
}

type TeaResource struct {
  Data Tea `json:"data"`  
}

func teaHandler(w http.ResponseWriter, r *http.Request) {
  res := getTeas() // Returns a *TeaCollection
  jsonapi.Write(w, res)  
}

func teaHandler(w http.ResponseWriter, r *http.Request) {
  res := getTea(mux.Params(r).ByName("id")) // Returns a *TeaResource
  jsonapi.Write(w, res)  
}

func createTeaHandler(w http.ResponseWriter, r *http.Request) {
  res := createTea(jsonapi.Body(r).(*TeaResource))
  jsonapi.Write(w, res)
}

func main() {
  m := jsonapi.New()
  m.Get("/teas").ThenFunc(teasHandler)
  m.Get("/teas/:id").ThenFunc(teaHandler)
  m.Post("/teas").Use(jsonapi.ContentTypeHandler, jsonapi.BodyHandler(TeaResource{})).ThenFunc(createTeaHandler)
}
```