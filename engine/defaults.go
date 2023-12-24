package engine

var defaultErrHandler Handler = func(c *Ctx) error {
	if c.Error == nil {
		c.Res.WriteHeader(500)
		c.Res.Write([]byte(c.Error.Error()))
	}

	return nil
}
