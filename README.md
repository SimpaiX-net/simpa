# Simpa: A Web Framework Inspired by ExpressJS

Simpa is a web framework designed to cater to the specific needs of Simpaix Telegram bot integration, providing a secure HTTP server endpoint for retrieving bot updates through a webhook. While Simpa is currently in active development, it is not yet fully covered and complete.

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
