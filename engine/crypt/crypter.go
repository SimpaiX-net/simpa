// # Introduction
//
// A crypter package is intended to securely encrypt or decrypt cookie values.
// Intended for securecookies and session cookies (which are also secure).
//
// # This API contains default crypters to provide example implementation to the interface
//
// [CrypterI]: https://pkg.go.dev/github.com/SimpaiX-net/simpa/engine/crypt#CrypterI
package crypt

import (
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"hash"
	"log"
)

// # CrypterI object
// Exposes the special methods a crypter has to introduce to satisfy the crypter interface.
//
// For examples please consider looking at the [simpa/engine/crypt] examples.
// Like:
//   - AES_GCM
//   - AES_CTR
//
// This crypter is used for securecookie's and sessions.
//
// The crypter object should be set using:
// [Engine SecureCookie]: https://pkg.go.dev/github.com/SimpaiX-net/simpa/engine#Engines
type CrypterI interface {
	/*
		Should encrypt given data and return it in base64 format
	*/
	Encrypt(data string) (string, error)
	/*
		Should Decrypt given data and return its plaintext.
		'data' here is the  base64 encoded string, when the data is decoded it represents
		the encryption cipher text in bytes
	*/
	Decrypt(data string) (string, error)
}

/*
Default crypter type;
uses  AES GCM
*/
type AES_GCM struct {
	aes_gcm cipher.AEAD
}

/*
Creates a new AES GCM object
*/
func New_AES_GCM(block func() cipher.Block) *AES_GCM {
	gcm, err := cipher.NewGCM(block())
	if err != nil {
		log.Fatal(err)
	}

	return &AES_GCM{
		aes_gcm: gcm,
	}
}

/*
Encrypts data and returns the base64 encoded string of the encrypted data or error on failure
*/
func (c *AES_GCM) Encrypt(data string) (string, error) {
	nonce := make([]byte, c.aes_gcm.NonceSize())

	_, err := rand.Read(nonce)
	if err != nil {
		return "", err
	}

	var encr []byte
	encr = append(nonce, c.aes_gcm.Seal(nil, nonce, []byte(data), nil)...)
	return base64.StdEncoding.EncodeToString(encr), nil
}

/*
Decrypts data and returns the plaintext string of the encrypted data or error.
The returned errors can be related to authentication or some sort of other failure
*/
func (c *AES_GCM) Decrypt(data string) (string, error) {
	enc, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	decr, err := c.aes_gcm.Open(nil, enc[:c.aes_gcm.NonceSize()], enc[c.aes_gcm.NonceSize():], nil)
	return string(decr), err
}

/*
Default crypter type;
uses  AES CTR HMAC
*/
type AES_CTR struct {
	block cipher.Block
	hmac  hash.Hash
}

/*
Creates a new AES GCM object
*/
func New_AES_CTR(block cipher.Block, hmac hash.Hash) *AES_CTR {
	return &AES_CTR{
		block,
		hmac,
	}
}

/*
Encrypts data and returns the base64 encoded string of the encrypted data or error on failure
*/
func (c *AES_CTR) Encrypt(data string) (string, error) {
	ciphertext := make([]byte, c.block.BlockSize()+c.hmac.Size()+len(data))
	iv := ciphertext[:c.block.BlockSize()]
	mac := ciphertext[c.block.BlockSize() : c.block.BlockSize()+c.hmac.Size()]

	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	stream := cipher.NewCTR(c.block, iv)
	stream.XORKeyStream(ciphertext[c.block.BlockSize()+c.hmac.Size():], []byte(data))
	{
		c.hmac.Reset()
		if _, err := c.hmac.Write(iv); err != nil {
			return "", err
		}
		if _, err := c.hmac.Write(ciphertext[c.block.BlockSize()+c.hmac.Size():]); err != nil {
			return "", err
		}

		copy(mac, c.hmac.Sum(nil))
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

/*
Decrypts data and returns the plaintext string of the encrypted data or error.
The returned errors can be related to authentication or some sort of other failure
*/
func (c *AES_CTR) Decrypt(data string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	iv := decoded[:c.block.BlockSize()]
	mac := decoded[c.block.BlockSize() : c.block.BlockSize()+c.hmac.Size()]
	ciph := decoded[c.block.BlockSize()+c.hmac.Size():]

	{
		c.hmac.Reset()
		if _, err := c.hmac.Write(iv); err != nil {
			return "", err
		}
		if _, err := c.hmac.Write(ciph); err != nil {
			return "", err
		}

		if !hmac.Equal(mac, c.hmac.Sum(nil)) {
			return "", errors.New("Authentication failed")
		}
	}

	decr := make([]byte, len(ciph))

	stream := cipher.NewCTR(c.block, iv)
	stream.XORKeyStream(decr, ciph)

	return string(decr), nil
}
