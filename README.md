<img src="https://github.com/SimpaiX-net/.github/assets/48758770/af960480-aa63-4be4-94bf-66d43453bb83" width="200" style="position: absolute; left:0;"><br>

# Simpa: A Web Framework Inspired by ExpressJS

Simpa is a web framework designed to cater to the specific needs of Simpaix Telegram bot integration, providing a secure HTTP server endpoint for retrieving bot updates through a webhook. While Simpa is currently in active development, it is not yet fully covered and complete.

### Features
- [x] HTTP2 & HTTP1.1 support
- [x] JSON body parser 
- [x] JSON binding support 
- [x] Validator engine 
- [x] Using the Fastest HTTP router 
- [x] Built upon STD library ``net/http``
- [x] Supports dynamic path for routes
- [x] ExpressJS like MVC
- [x] Templating & rendering
- [x] Limit request body
- [x] Secure cookie implementation
- [x] Support to provide custom body parser
- [x] AES_GCM and AES_CTR default crypters for secure cookie & session
- [x] Support to provide custom crypter
- [x] Session Implementation w/custom crypter or default crypter support
- [x] Storage: mongo driver for storage

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
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SimpaiX-net/simpa/engine"
	"github.com/SimpaiX-net/simpa/engine/crypt"
	"github.com/SimpaiX-net/simpa/engine/parsers/bodyparser"
	"github.com/SimpaiX-net/simpa/engine/sessions"
	"github.com/SimpaiX-net/simpa/engine/sessions/drivers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

		v, ok := sess.Get("world").(string)
		if !ok {
			sess.Set("world", "a-")
		} else {
			sess.Set("world", v+"b-")
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

		return c.String(200, sess.Get("world").(string))
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
