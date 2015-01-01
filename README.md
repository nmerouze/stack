# stack [![GoDoc](https://godoc.org/github.com/nmerouze/stack?status.png)](https://godoc.org/github.com/nmerouze/stack/jsonapi) [![Build Status](https://travis-ci.org/nmerouze/stack.svg?branch=master)](https://travis-ci.org/nmerouze/stack)

You are on the experimental branch of stack, a framework to build JSON APIs faster. The master branch currently contains the basic components, this branch will have more higher level components.

# Current ideas

## Entities

To make APIs with stack, you would just need an entity. It's a simple structure:

``` go
type Tea struct {
  Id string `json:"id"`
  Name string `json:"name"`
}
```

Then it would need to implement interfaces: `stack.Read`, `stack.Find`, `stack.Create`, `stack.Update` and/or `stack.Delete`. Here's an example implementing `stack.Read`:

``` go
type Tea struct {
  Id string `json:"id"`
  Name string `json:"name"`
}

func (t *Tea) Read(c *stack.Context) error {
  results, err := c.Db.GetTeas(c.Params["id"])
  if err != nil {
    return err
  }

  c.ResBody = results

  return nil
}
```

Now the entity can be use in the router:

``` go
func main() {
  api := New()
  api.Read("/teas", new(Tea))
  http.ListenAndServe(":8080", api)
}
```

This way the router handles almost every possible error, you just have to focus on the data layer. If things has to be cutomized, there will be a fallback on lower level components. And there will always have middlewares so features like authentication and authorizations can be implemented without disrupting anything.

## Data validation

I think a crucial feature is data validation. URL params and request bodies need to be validated. Here's how I see it:

``` go
var getTeasSchema = NewSchema(
  NewParam("limit").Int().Max(10)
)

var postTeaSchema = NewSchema(
  NewParam("name").String().Required()
)

func main() {
  api := New()
  api.Get("/teas", new(Tea)).Params(getTeasSchema)
  api.Post("/teas", new(Tea)).ReqBody(postTeaSchema)
  api.Get("/teas/:id", new(Tea))
  http.ListenAndServe(":8080", api)
}
```

## Documentation generation

Developers (sadly) don't want to bother with documentation. Maintaining it can be very time consuming so a part of building JSON APIs faster is to automate the documentation as much as it is possible.

To be able to generate the documentation, all the features must be made in a way it's possible to pass through them to extract information. We typically need:

- Method & URL
- Params: type, validations, description, etc.
- Possible responses (2xx, 4xx)
- Code examples

I think the documentation of the [Heroku Platform API](https://devcenter.heroku.com/articles/platform-api-reference) is really great. The overview contains every global feature of the API and could be generated easily by the framework. There is one block with all possible errors, a far more practical approach for documentation generation that having to display possible errors on each endpoint (because we may not know which errors can happen for a certain endpoint). Then the rest of the documentation can focus on the endpoints with: title, description, params, examples.

# Conclusion

This is a rough draft of what I would like to create in the coming weeks. There are far more things I could discuss right now but these 3 features are the most important in my opinion.