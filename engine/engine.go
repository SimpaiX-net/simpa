package engine

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
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
		router     *httprouter.Router
		routes     []*Route
		errHandler Handler
	}
)

func New() *Engine {
	return &Engine{
		errHandler: defaultErrHandler,
		router:     httprouter.New(),
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
	e.router.Handle(http.MethodGet, name, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		route, err := e.GetRoute(name)
		if err != nil {
			w.WriteHeader(404)
			return
		}

		c := &Ctx{
			Req:    *r,
			Res:    w,
			Params: p,
			Error:  nil,
		}
		for _, v := range route.handlers {
			if err := v(c); err != nil {
				e.errHandler(c)
				return
			}

			if c.Error != nil {
				e.errHandler(c)
				return
			}
		}
	})
}
