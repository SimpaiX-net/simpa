package bodyparser

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// BodyParser object
type BodyParser struct {
	req *http.Request // request context
}

/*
the methods a body parser has to implement
*/
type BodyParserI interface {
	// a parser to support multiple content type parsing
	Parse(dest interface{}, ct Binding) error
	// gets the corresponding request context
	GetRequest() *http.Request
	// function that sets request object into the bodyparser struct field 'req'
	New(r *http.Request)
}

/*
Creates a new bodyparser object
*/
func (b *BodyParser) New(r *http.Request) {
	b.req = r
}

/*
Parse request body to given binding, parses it into dest, so dest should be a pointer struct.
Currently only parsing JSON body's is supported

Soon we'll add support for multiple bindings
*/
func (b *BodyParser) Parse(dest interface{}, ct Binding) error {
	switch ct {
	case JSON:
		if err := b.parse_json(dest); err != nil {
			return err
		}
	default:
		return errors.New("given binding is currently not supported")
	}

	return nil
}

/*
Gets the corresponding request context
*/
func (b *BodyParser) GetRequest() *http.Request {
	return b.req
}

func (b *BodyParser) parse_json(dest interface{}) error {
	req := b.GetRequest()
	if req.Header.Get("content-type") != "application/json" {
		return errors.New("'Content-Type' is not JSON")
	}

	bd, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bd, dest)
}
