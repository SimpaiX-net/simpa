package engine

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SimpaiX-net/simpa/engine/parsers/bodyparser"
	"github.com/julienschmidt/httprouter"
)

type H map[string]any

/*
request and response context
*/
type Ctx struct {
	Error      error                  // represents an error
	Req        *http.Request          // http request
	Res        http.ResponseWriter    // http response
	Params     httprouter.Params      // http params
	BodyParser bodyparser.BodyParserI // body parser
	engine     *Engine                // underlying app engine
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

/*
Renders given HTML template file with the std go templating engine
*/
func (c *Ctx) RenderHTML(name string, data H) error {
	c.Res.Header().Set("content-type", "text/html")
	return c.engine.template.ExecuteTemplate(c.Res, name, data)
}

/*
Sets secure cookie if `engine.SecureCookie` is set, otherwise it sets a cookie whom's value is not encrypted
*/
func (c *Ctx) SetCookie(cookie *http.Cookie) error {
	if c.engine.SecureCookie != nil {
		enc, err := c.engine.SecureCookie.Encrypt(cookie.String())
		if err != nil {
			fmt.Println("hiero")
			return err
		}
		cookie.Value = enc
	}

	http.SetCookie(c.Res, cookie)
	return nil
}

/*
Decodes encrypted cookie named 'name' to dest
*/
func (c *Ctx) DecodeCookie(name string, dest *http.Cookie) error {
	cookie, err := c.Req.Cookie(name)
	if err != nil {
		return err
	}

	decr, err := c.engine.SecureCookie.Decrypt(cookie.Value)
	if err != nil {
		return err
	}

	fmt.Println(decr)

	cookie.Value = decr
	*dest = *cookie

	return nil
}
