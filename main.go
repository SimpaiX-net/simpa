package main

import (
	"errors"
	"fmt"

	"github.com/SimpaiX-net/simpa/engine"
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

	app.Run(":2000")
}
