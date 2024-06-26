<img src="https://github.com/SimpaiX-net/.github/assets/48758770/af960480-aa63-4be4-94bf-66d43453bb83" width="200" style="position: absolute; left:0;"><br>

# Simpa: A Web Framework Inspired by ExpressJS

The main goal of Simpa is to redefine elegantness. Exposing an engine that is extensible in every aspect to it's core. Some examples of these can be Simpa's crypter interface.
While erasing many complex duties of the developer it did not breach the ability for anyone to extent or modify.

#### Disclaimer & note
Please do note that Simpa is far from complete. There has been a revision planned on the core. Additionally unitests files would be added with the upcoming revision.


Scroll further for a list of features.

##### Install

> go get simpaix.net/simpa/v2

### Benchmark
- ``30k / rqps`` on session + crypter route with ``~200ms``
- ``600k / rqps`` on JSON processor and marshall route with ``~10-50ms``

TODO: Video showcase and unitest files 

### Features

- [X] HTTP2 & HTTP1.1 support
- [X] JSON body parser
- [X] JSON binding support
- [X] Validator engine
- [X] Using the Fastest HTTP router
- [X] Built upon STD library ``net/http``
- [X] Supports dynamic path for routes
- [X] ExpressJS like MVC
- [X] Templating & rendering
- [X] Limit request body
- [X] Secure cookie implementation
- [X] Support to provide custom body parser
- [X] AES_GCM and AES_CTR default crypters for secure cookie & session
- [X] Support to provide custom crypter
- [X] Session Implementation w/custom crypter or default crypter support
- [X] Storage: mongo driver for storage

### Todo

- [ ] XML binding support
- [ ] XML body parser
- [ ] JSON body parser support for ``map[any]any``
- [ ] JWT ware implementation

> You can give feedbacks using the 'Issues' tab in this repository.

## Example

### File: `main.go`

```go
package main

import (
	"context"
	"crypto/aes"
	"crypto/hmac"
	"crypto/sha512"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"simpaix.net/simpa/v2/engine"
	"simpaix.net/simpa/v2/engine/crypt"
	"simpaix.net/simpa/v2/engine/parsers/bodyparser"
	"simpaix.net/simpa/v2/engine/sessions"
	"simpaix.net/simpa/v2/engine/sessions/drivers"
)

func hello(c *engine.Ctx) error {
	name := c.Req.URL.Query().Get("name")
	if name == "" {
		return c.String(403, c.Error.Error())
	}
	return nil
}

func main() {
	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	sess := c.Database("testdb").Collection("sessions")

	app := engine.New()
	{

		app.MaxBodySize = 1000000 // 1MB

		{
			/*
				example
			*/
			key := []byte("hallowereld1234secret32323232acs")

			aes, err := aes.NewCipher(key)
			if err != nil {
				log.Fatal(err)
			}

			hmac := hmac.New(sha512.New, []byte("secret123"))
			app.SecureCookie = crypt.New_AES_GCM(aes, hmac)
			app.Storage = drivers.NewMongoStore(sess, time.Duration(time.Second*5), app.SecureCookie)
		}

	}

	app.Get("/", func(c *engine.Ctx) error {
		ck := &sessions.Config{
			Name: "SESS_ID1",
		}
		sess, err := c.Session(ck)
		if err != nil {
			return err
		}

		v, ok := sess.Get("world").(float64)
		if !ok {
			sess.Set("world", 1)
		} else {
			sess.Set("world", int(v+1))
		}

		if err := sess.Save(c.Res); err != nil {
			return err
		}
		return c.String(200, "")
	})

	app.Get("/print-world", func(c *engine.Ctx) error {
		ck := &sessions.Config{
			Name: "SESS_ID1",
		}
		sess, err := c.Session(ck)
		if err != nil {
			return err
		}

		v, ok := sess.Get("world").(float64)
		if !ok {
			return errors.New("unexpected data type")
		}

		return c.String(200, fmt.Sprintln(int(v)))
	})

	app.Get("/set", func(c *engine.Ctx) error {
		if err := c.SetCookie(&http.Cookie{Name: "hello", Value: "123", Secure: false, Expires: time.Now().Add(time.Second * 10), Path: "/"}); err != nil {
			return err
		}

		return c.String(200, "success")
	})

	app.Get("/get", func(c *engine.Ctx) error {
		cookie := &http.Cookie{}

		if err := c.DecodeCookie("hello", cookie); err != nil {
			return err
		}

		fmt.Println("cookie value", cookie)
		return c.String(200, "success")
	})

	app.Get("/hello/:id", hello, func(c *engine.Ctx) error {
		fmt.Println("id:", c.Params.ByName("id"))
		return c.String(200, c.Req.URL.Query().Get("name"))
	})

	app.Post("/json", func(c *engine.Ctx) error {
		dummy := struct {
			Name string `json:"name"`
		}{}

		if err := c.BodyParser.Parse(&dummy, bodyparser.JSON); err != nil {
			return c.String(403, err.Error())
		}

		return c.JSON(200, dummy)
	})
	app.Run(":2000")
}


```

## Author

- **z3ntl3**
- **Simpaix**

For any inquiries or contributions, please feel free to reach out. Your feedback is highly appreciated!
