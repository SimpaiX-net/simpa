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
- [x] Support to provide custom body parser

### Todo
- [ ] Change secure cookie implementation with custom one, to fit AES-GCM only
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
	"fmt"
	"net/http"
	"html/template"

	"github.com/SimpaiX-net/simpa/engine"
	"github.com/SimpaiX-net/simpa/engine/parsers/bodyparser"
	"github.com/gorilla/securecookie"
)

func hello(c *engine.Ctx) error {
	name := c.Req.URL.Query().Get("name")
	if name == "" {
		return c.String(403, c.Error.Error())
	}
	return nil
}

func main() {
	app := engine.New()
	{
		app.MaxBodySize = 1000000 // 1MB
		app.SecureCookie = securecookie.New(
			securecookie.GenerateRandomKey(32),
			securecookie.GenerateRandomKey(32),
		)
	}

	temp := template.Must(template.
		New("views").
		Funcs(template.FuncMap{}).
		ParseGlob("views/*"),
	)

	app.SetTemplate(temp)

	app.Get("/", func(c *engine.Ctx) error {
		return c.RenderHTML("index.html", engine.H{
			"title": struct{ Name string }{
				Name: "Welcome Screen",
			},
		})
	})

	app.Get("/set", func(c *engine.Ctx) error {
		if err := c.SetCookie(&http.Cookie{Name: "hello", Value: "123", MaxAge: 3600}); err != nil {
			return err
		}

		return c.String(200, "success")
	})

	app.Get("/get", func(c *engine.Ctx) error {
		cookie := &http.Cookie{}

		if err := c.DecodeCookie("hello", &cookie.Value); err != nil {
			return err
		}

		fmt.Println(cookie.Value)
		return c.String(200, "success")
	})

	app.Get("/hello/:id", hello, func(c *engine.Ctx) error {
		fmt.Println("id:", c.Params.ByName("id"))
		return c.String(200, c.Req.URL.Query().Get("name"))
	})

	app.Post("/json", func(c *engine.Ctx) error {
		dummy := struct {
			Name string `json:"name"`
		}{}

		if err := c.BodyParser.Parse(&dummy, bodyparser.JSON); err != nil {
			return c.String(403, err.Error())
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
