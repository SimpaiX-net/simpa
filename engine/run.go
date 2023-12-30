package engine

import (
	"log"
	"net/http"
)

// run the http server on given addr
func (e *Engine) Run(addr string) {
	if e.SecureCookie == nil {
		log.Fatalf("SecureCookie crypter is not, set. Set it before starting the web server!")
	}

	e.router.PanicHandler = e.panicHandler
	http.ListenAndServe(addr, e.router)
}
