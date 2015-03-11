# stack [![GoDoc](https://godoc.org/github.com/nmerouze/stack?status.png)](https://godoc.org/github.com/nmerouze/stack/jsonapi) [![Build Status](https://travis-ci.org/nmerouze/stack.svg?branch=master)](https://travis-ci.org/nmerouze/stack)

stack is a framework to build [JSON-APIs](http://jsonapi.org) faster. It is based on the series of articles ["Build Your Own Web Framework in Go"](http://nicolasmerouze.com/build-web-framework-golang/) I began to write a few weeks ago.

The public API of the package is stable and you can use it right now to make your application. I don't indend to break any existing feature but will add new features to make the framework more useful for production. [Look at the documentation](http://godoc.org/github.com/nmerouze/stack/jsonapi).

Checkout the ["Unstable Features" section](#unstable-features) below to read about upcoming/work-in-progress features.

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

func teasHandler(w http.ResponseWriter, r *http.Request) {
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

# Unstable Features

## `schema` package

You can find this package on the [`schema` branch](https://github.com/nmerouze/stack/tree/schema). It aims at validating URL params and JSON request bodies. It also generates a JSON schema from the schema definition which can serve as a documentation, or to auto-build client libraries.

## `data` package

You can find this package on the [`data` branch](https://github.com/nmerouze/stack/tree/data). It aims at providing interfaces to connect the data layer of an application to the router. This way you can just write your models following these interfaces and the router will take care of the rest. No `http.Handler` to write. As you see in the current ["Getting Started" section](#getting-started), you need to manually make a handler to call your data layer then write the response. This code is always the same and will not be necessary once this package becomes stable.