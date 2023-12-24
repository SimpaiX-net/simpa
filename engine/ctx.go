package engine

import "net/http"

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