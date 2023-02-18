package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

type Crypt struct {
	block cipher.Block
}

func NewCrypt(key string) (*Crypt, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	return &Crypt{block: block}, nil
}

func (c *Crypt) Encrypt(plainText string) (string, error) {
	cipherText := make([]byte, aes.BlockSize+len([]byte(plainText)))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	encryptStream := cipher.NewCTR(c.block, iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], []byte(plainText))
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (c *Crypt) Decrypt(cipherText string) (string, error) {
	cipherByte, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	decryptedText := make([]byte, len([]byte(cipherByte[aes.BlockSize:])))
	decryptStream := cipher.NewCTR(c.block, []byte(cipherByte[:aes.BlockSize]))
	decryptStream.XORKeyStream(decryptedText, []byte(cipherByte[aes.BlockSize:]))
	return string(decryptedText), nil
}
