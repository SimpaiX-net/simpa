<img src="https://github.com/SimpaiX-net/.github/assets/48758770/af960480-aa63-4be4-94bf-66d43453bb83" width="200" style="position: absolute; left:0;"><br>

# Simpa: A Web Framework Inspired by ExpressJS

Simpa is a web framework designed to cater to the specific needs of Simpaix Telegram bot integration, providing a secure HTTP server endpoint for retrieving bot updates through a webhook. While Simpa is currently in active development, it is not yet fully covered and complete.

### Features
- [x] HTTP2 & HTTP1.1 support
- [x] JSON body parser 
- [x] JSON binding support 
- [x] Validator engine 
- [x] Using the Fastest HTTP router 
- [x] Built upon STD library ``net/http``
- [x] Supports dynamic path for routes
- [x] ExpressJS like MVC
- [x] Templating & rendering
- [x] Limit request body
- [x] Secure cookie implementation

### Todo
- [ ] XML binding support
- [ ] XML body parser
- [ ] JSON body parser support for ``map[any]any``
- [ ] Session implementation
- [ ] JWT ware implementation
> You can give feedbacks using the 'Issues' tab in this repository.

## Example

### File: `main.go`

```go
package main

import (
    errors "errors"
    fmt "fmt"
    engine "github.com/SimpaiX-net/simpa/engine"
)

func hello(c *engine.Ctx) error {
    name := c.Req.URL.Query().Get("name")
    if name == "" {
        c.Error = errors.New("no name given") // abort
        return c.String(403, c.Error.Error())
    }
    return nil
}

func main() {
    app := engine.New()
    app.Get("/hello/:id", hello, func(c *engine.Ctx) error {
        fmt.Println("id:", c.Params.ByName("id"))
        return c.String(200, c.Req.URL.Query().Get("name"))
    })

    app.Post("/json", func(c *engine.Ctx) error {
        dummy := struct {
            Name string `json:"name"`
        }{}

        if err := c.ParseJSON(&dummy); err != nil {
            c.Error = err
            return c.String(403, c.Error.Error())
        }

        return c.JSON(200, dummy)
    })
    app.Run(":2000")
}
```

## Author

- **z3ntl3**
- **Simpaix**

For any inquiries or contributions, please feel free to reach out. Your feedback is highly appreciated!
