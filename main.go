package main

import (
	"crypto/aes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SimpaiX-net/simpa/engine"
	"github.com/SimpaiX-net/simpa/engine/crypt"
	"github.com/SimpaiX-net/simpa/engine/parsers/bodyparser"
)

func hello(c *engine.Ctx) error {
	name := c.Req.URL.Query().Get("name")
	if name == "" {
		return c.String(403, c.Error.Error())
	}
	return nil
}

func main() {
	app := engine.New()
	{
		app.MaxBodySize = 1000000 // 1MB

		{
			/*
				example
			*/
			randKey := make([]byte, 32)
			rand.Read(randKey)

			aes, err := aes.NewCipher(randKey)
			if err != nil {
				log.Fatal(err)
			}

			hmac := hmac.New(sha512.New, []byte("secret123"))
			app.SecureCookie = crypt.New_AES_GCM(aes, hmac)
		}
	}

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
