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
Creates a new bodyparser object
*/
func New(r *http.Request) *BodyParser {
	return &BodyParser{r}
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

func (b *BodyParser) parse_json(dest interface{}) error {
	if b.req.Header.Get("content-type") != "application/json" {
		return errors.New("'Content-Type' is not JSON")
	}

	bd, err := io.ReadAll(b.req.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bd, dest)
}
