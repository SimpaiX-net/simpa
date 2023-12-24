package engine

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Ctx struct {
	Error  error
	Req    http.Request
	Res    http.ResponseWriter
	Params httprouter.Params
}

func (c *Ctx) String(status int, data string) error {
	c.Res.WriteHeader(status)
	if _, err := c.Res.Write([]byte(data)); err != nil {
		return err
	}

	return nil
}

func (c *Ctx) JSON(status int, data interface{}) error {
	c.Res.WriteHeader(status)
	c.Res.Header().Set("content-type", "application/json")

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := c.Res.Write(b); err != nil {
		return err
	}
	return nil
}

// Parses JSON body to dest, it has to be a pointer.
func (c *Ctx) ParseJSON(dest interface{}) error {
	if c.Req.Header.Get("content-type") != "application/json" {
		return errors.New("'Content-Type' is not JSON")
	}

	b, err := io.ReadAll(c.Req.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, dest)
}
