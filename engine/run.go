package engine

import (
	"net/http"
)

func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, e.router)
}
