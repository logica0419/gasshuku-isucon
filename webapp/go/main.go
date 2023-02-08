package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"

	"github.com/labstack/echo/v4"
)

var block cipher.Block

func encrypt(plainText string) (string, error) {
	cipherText := make([]byte, aes.BlockSize+len([]byte(plainText)))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], []byte(plainText))
	return string(cipherText), nil
}

func decrypt(cipherText string) string {
	decryptedText := make([]byte, len([]byte(cipherText[aes.BlockSize:])))
	decryptStream := cipher.NewCTR(block, []byte(cipherText[:aes.BlockSize]))
	decryptStream.XORKeyStream(decryptedText, []byte(cipherText[aes.BlockSize:]))
	return string(decryptedText)
}

func main() {
	key := []byte("1234567890123456")

	var err error

	block, err = aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello World!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
