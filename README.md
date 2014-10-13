# Stack

Stack is a middleware system with simple contexts. It is a fork of [alice](https://github.com/justinas/alice) and uses a context system similar to [Goji](https://goji.io/).

``` go
func authHandler(c *stack.C, next http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
    c.Env["user"] = "ochasuki"
    next.ServeHTTP(c, w, r)
  }

  return stack.HandlerFunc(fn)
}

func appHandler(c stack.C, w http.ResponseWriter, r *http.Request) {
  w.Write([]byte(c.Env["user"].(string)))
}

func main() {
  http.Handle("/", stack.New(authHandler).ThenFunc(appHandler))
  http.ListenAndServe(":8080", nil)
}
```

# Why

- alice doesn't support contexts
- negroni doesn't support contexts, and its 3rd argument is useless ([read this blog post](http://nicolasmerouze.com/middlewares-golang-best-practices-examples/))
- goji is a framework. stack is BYOR (Bring Your Own Router).