package engine

import (
	"net/http"
)

func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, e.handler())
}

func (e *Engine) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		route, err := e.GetRoute(path)
		if err != nil {
			w.WriteHeader(404)
			return
		}

		c := &Ctx{
			Req:   *r,
			Res:   w,
			Error: nil,
		}
		for _, v := range route.handlers {
			if err := v(c); err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))

				c.Error = err
				return
			}

			if c.Error != nil {
				return
			}
		}
	}
}
