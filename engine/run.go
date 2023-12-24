package engine

import (
	"net/http"
)

func (e *Engine) Run(addr string) {
	e.router.PanicHandler = e.panicHandler
	http.ListenAndServe(addr, e.router)
}
