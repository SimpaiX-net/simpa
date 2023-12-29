package main

import (
	"crypto/aes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"log"

	"github.com/SimpaiX-net/simpa/engine"
	"github.com/SimpaiX-net/simpa/engine/crypt"
)

func main() {
	app := engine.New()
	{
		app.MaxBodySize = 1000000 // 1MB

		{
			// do not use this example in a production app.
			// you should not randomize the AES key.
			// rather import and load it from your environment and make sure it is 256 bit / 32 bytes
			randKey := make([]byte, 32)
			rand.Read(randKey)

			aes, err := aes.NewCipher(randKey)
			if err != nil {
				log.Fatal(err)
			}

			hmac := hmac.New(sha512.New, []byte("secret123"))
			app.SecureCookie = crypt.New_AES_CTR(aes, hmac)
		}
	}

	app.Get("/", func(c *engine.Ctx) error {
		return c.String(200, "hello")
	})

	app.Run(":2000")
}
