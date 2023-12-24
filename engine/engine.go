package engine

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/SimpaiX-net/simpa/engine/binding"
	"github.com/julienschmidt/httprouter"
)

type (
	// route handling function
	Handler func(c *Ctx) error
	// defines a route
	Route struct {
		method   string    // route methpd
		name     string    // route name
		handlers []Handler // route handler(s)
	}
)

type (
	// Engine describes the engine for the HTTP server
	Engine struct {
		router       *httprouter.Router                                          // http router
		routes       []*Route                                                    // routes context
		errHandler   Handler                                                     // error handler
		panicHandler func(w http.ResponseWriter, r *http.Request, i interface{}) // panic handler
		validator    *binding.ValidatorImpl                                      // validator engine
		template     *template.Template                                          // template
	}
)

/*
Creates new engine with default config
*/
func New() *Engine {
	return &Engine{
		panicHandler: func(w http.ResponseWriter, r *http.Request, i interface{}) {
			w.WriteHeader(500)
		},
		errHandler: defaultErrHandler,
		router:     httprouter.New(),
		validator:  &binding.DefaultValidator,
	}
}

/*
Set templating
*/
func (e *Engine) SetTemplate(temp *template.Template) {
	e.template = temp
}

/*
Gets the underlying templating engine if it is present
*/
func (e *Engine) GetTemplate() *template.Template {
	return e.template
}

/*
Define custom validator engine. Keep in mind that validator should be a struct pointer
See: '/binding/validator.go' for example
*/
func (e *Engine) SetValidatorEngine(validator *binding.ValidatorImpl) {
	e.validator = validator
}

/*
Set custom error handling function
*/
func (e *Engine) SetErrorHandler(h Handler) {
	e.errHandler = h
}

// Get existing route by name and it's method
func (e *Engine) GetRoute(name string, method string) (*Route, error) {
	for _, h := range e.routes {
		if h.name == name && h.method == method {
			return h, nil
		}
	}

	return nil, errors.New("Cannot find route")
}

// Register POST route; shorthand for 'RegisterRoute'
func (e *Engine) Get(name string, handler ...Handler) {
	e.RegisterRoute(name, http.MethodGet, handler...)
}

// Register GET route; shorthand for 'RegisterRoute'
func (e *Engine) Post(name string, handler ...Handler) {
	e.RegisterRoute(name, http.MethodPost, handler...)
}

// Function to register route with all HTTP methods supported
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
			temp:   e.template,
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
