package engine

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type H map[string]any

/*
request and response context
*/
type Ctx struct {
	Error  error               // represents an error
	Req    http.Request        // http request
	Res    http.ResponseWriter // http response
	Params httprouter.Params   // http params
	temp   *template.Template
}

// Sends string with custom status code
func (c *Ctx) String(status int, data string) error {
	c.Res.WriteHeader(status)
	if _, err := c.Res.Write([]byte(data)); err != nil {
		return err
	}

	return nil
}

// Sends JSON with 'application/json' content type.
// 'data' is a pointer to the struct, and it is a JSON unmarshalled object
// this function marshalls the JSON and sends it to the client
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

/*
Renders given HTML template file with the std go templating engine
*/
func (c *Ctx) RenderHTML(name string, data H) error {
	return c.temp.ExecuteTemplate(c.Res, name, data)
}
