package engine

import "fmt"

// default error handling function
var defaultErrHandler Handler = func(c *Ctx) error {
	if c.Error != nil {
		c.String(500, fmt.Sprintf("[Error]: %s", c.Error.Error()))
	}

	return nil
}
