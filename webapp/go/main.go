package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/labstack/echo/v4"
)

const qrCodeFileName = "../images/qr.png"

var block cipher.Block

func encrypt(plainText string) (string, error) {
	cipherText := make([]byte, aes.BlockSize+len([]byte(plainText)))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], []byte(plainText))
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func decrypt(cipherText string) (string, error) {
	cipherByte, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	decryptedText := make([]byte, len([]byte(cipherByte[aes.BlockSize:])))
	decryptStream := cipher.NewCTR(block, []byte(cipherByte[:aes.BlockSize]))
	decryptStream.XORKeyStream(decryptedText, []byte(cipherByte[aes.BlockSize:]))
	return string(decryptedText), nil
}

func generateQRCode(id string) ([]byte, error) {
	encryptedID, err := encrypt(id)
	if err != nil {
		return nil, err
	}

	_, err = exec.
		Command("qrencode", "-o", qrCodeFileName, "-t", "PNG", "-s", "1", "-v", "5", "--strict-version", "-l", "M", encryptedID).
		Output()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(qrCodeFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
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
