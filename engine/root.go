package engine

import (
	"errors"
	"net/http"
)

type Ctx struct {
	Error error
	Req   http.Request
	Res   http.ResponseWriter
}

func (c *Ctx) String(status int, data string) error {
	c.Res.WriteHeader(status)
	if _, err := c.Res.Write([]byte(data)); err != nil {
		return err
	}

	return nil
}

type (
	Handler func(c *Ctx) error
	Route   struct {
		name     string
		handlers []Handler
	}
)

type (
	Engine struct {
		routes []*Route
	}
)

func New() *Engine {
	return &Engine{}
}

// get existing route by name
func (e *Engine) GetRoute(name string) (*Route, error) {
	for _, h := range e.routes {
		if h.name == name {
			return h, nil
		}
	}

	return nil, errors.New("Cannot find route")
}

// register route
func (e *Engine) Get(name string, handler ...Handler) {
	e.routes = append(e.routes, &Route{
		name:     name,
		handlers: handler,
	})
}
