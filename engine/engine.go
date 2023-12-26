package engine

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/SimpaiX-net/simpa/engine/binding"
	"github.com/SimpaiX-net/simpa/engine/parsers/bodyparser"
	"github.com/gorilla/securecookie"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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
		validator    binding.ValidatorImpl                                       // validator engine
		template     *template.Template                                          // template
		MaxBodySize  int64                                                       // max request body size
		SecureCookie *securecookie.SecureCookie                                  // secure cookie impl
		bodyparser   bodyparser.BodyParserI                                      // body parser
	}
)

/*
Creates new engine with default config
*/
func New() *Engine {
	return &Engine{
		panicHandler: func(w http.ResponseWriter, r *http.Request, i interface{}) {
			w.WriteHeader(500)
			fmt.Println("recovered")
		},
		errHandler:  defaultErrHandler,
		router:      httprouter.New(),
		validator:   binding.DefaultValidator,
		MaxBodySize: 1042 * 4,
		bodyparser:  nil,
	}
}

func (e *Engine) SetBodyParser(parser bodyparser.BodyParserI) {
	e.bodyparser = parser
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
func (e *Engine) SetValidatorEngine(validator binding.ValidatorImpl) {
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

	var h http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		route, err := e.GetRoute(name, method)
		if err != nil {
			w.WriteHeader(404)
			return
		}

		p := e.bodyparser
		if p == nil {
			p = bodyparser.DefaultBodyParser
		}

		p.New(r) // initialize request context with bodyparser

		v, ok := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
		if !ok {
			v = []httprouter.Param{}
		}
		c := &Ctx{
			Req:        r,
			Res:        w,
			Error:      nil,
			Params:     v,
			BodyParser: p,
			engine:     e,
		}

		for _, v := range route.handlers {
			if err := v(c); err != nil {
				c.Error = err
				e.errHandler(c)

				return
			}
		}
	}

	e.router.Handler(method, name, http.MaxBytesHandler(h2c.NewHandler(h, &http2.Server{}), e.MaxBodySize))
}
