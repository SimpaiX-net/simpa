package engine

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type (
	Handler func(c *Ctx) error
	Route   struct {
		method   string
		name     string
		handlers []Handler
	}
)

type (
	Engine struct {
		router       *httprouter.Router
		routes       []*Route
		errHandler   Handler
		panicHandler func(w http.ResponseWriter, r *http.Request, i interface{})
	}
)

func New() *Engine {
	return &Engine{
		panicHandler: func(w http.ResponseWriter, r *http.Request, i interface{}) {
			w.WriteHeader(500)
		},
		errHandler: defaultErrHandler,
		router:     httprouter.New(),
	}
}

func (e *Engine) SetErrorHandler(h Handler) {
	e.errHandler = h
}

// get existing route by name
func (e *Engine) GetRoute(name string, method string) (*Route, error) {
	for _, h := range e.routes {
		if h.name == name && h.method == method {
			return h, nil
		}
	}

	return nil, errors.New("Cannot find route")
}

// register route
func (e *Engine) Get(name string, handler ...Handler) {
	e.RegisterRoute(name, http.MethodGet, handler...)
}

// register route
func (e *Engine) Post(name string, handler ...Handler) {
	e.RegisterRoute(name, http.MethodPost, handler...)
}

func (e *Engine) RegisterRoute(name, method string, handler ...Handler) {
	e.routes = append(e.routes, &Route{
		name:     name,
		method:   method,
		handlers: handler,
	})
	e.router.Handle(http.MethodPost, name, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		route, err := e.GetRoute(name, http.MethodPost)
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
