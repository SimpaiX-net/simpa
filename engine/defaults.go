package engine

import (
	"fmt"

	"github.com/go-errors/errors"
)

// default error handling function
var defaultErrHandler Handler = func(c *Ctx) error {
	if c.Error != nil {
		c.String(400, fmt.Sprintf("[Error]: %s", errors.New(c.Error).ErrorStack()))
	}

	return nil
}
