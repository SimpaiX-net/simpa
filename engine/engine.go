package engine

import (
	"errors"
)

type (
	Handler func(c *Ctx) error
	Route   struct {
		name     string
		handlers []Handler
	}
)

type (
	Engine struct {
		routes     []*Route
		errHandler Handler
	}
)

func New() *Engine {
	return &Engine{
		errHandler: defaultErrHandler,
	}
}

func (e *Engine) SetErrorHandler(h Handler) {
	e.errHandler = h
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
