package crypt

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
)

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
uses  AES GCM algorithm
*/
type Crypter struct {
	aes_gcm cipher.AEAD
}

const Delimiter = byte('%')

func New(block func() cipher.Block) *Crypter {
	gcm, err := cipher.NewGCM(block())
	if err != nil {
		log.Fatal(err)
	}

	return &Crypter{
		aes_gcm: gcm,
	}
}

func (c *Crypter) Encrypt(data string) (string, error) {
	nonce := make([]byte, c.aes_gcm.NonceSize())

	_, err := rand.Read(nonce)
	if err != nil {
		return "", err
	}

	encr := c.aes_gcm.Seal(nil, nonce, []byte(data), nil)
	encr = append(encr, Delimiter)
	encr = append(encr, nonce...)

	return base64.StdEncoding.EncodeToString(encr), nil
}

func (c *Crypter) Decrypt(data string) (string, error) {
	enc, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	decr, err := c.aes_gcm.Open(nil, enc[:c.aes_gcm.NonceSize()], enc[c.aes_gcm.NonceSize():], nil)
	return string(decr), err
}

func (c *Crypter) deserialize(encr []byte) ([]byte, []byte, error) {
	nonceSize := c.aes_gcm.NonceSize()
	if len(encr) < nonceSize {
		return nil, nil, errors.New("message too small")
	}

	nonce := encr[:nonceSize]
	message := encr[nonceSize:]
	return message, nonce, nil
}
