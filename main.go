package main

import (
	"fmt"
	"html/template"
	"net/http"

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
			c.Error = err
			return c.String(403, c.Error.Error())
		}

		return c.JSON(200, dummy)
	})
	app.Run(":2000")
}
