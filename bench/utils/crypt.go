package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/logica0419/helpisu"
)

// AES + CTRモード + base64による暗号 / 復号化ツール
type Crypt struct {
	block cipher.Block
	cache *helpisu.Cache[string, string]
}

// 暗号化ツールを生成
func NewCrypt(key string) (*Crypt, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	return &Crypt{
		block: block,
		cache: helpisu.NewCache[string, string](),
	}, nil
}

// 暗号化
func (c *Crypt) Encrypt(plainText string) (string, error) {
	if v, ok := c.cache.Get(plainText); ok {
		return v, nil
	}

	cipherText := make([]byte, aes.BlockSize+len([]byte(plainText)))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	encryptStream := cipher.NewCTR(c.block, iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], []byte(plainText))
	encryptText := base64.URLEncoding.EncodeToString(cipherText)

	c.cache.Set(plainText, encryptText)
	return encryptText, nil
}

// 復号化
func (c *Crypt) Decrypt(cipherText string) (string, error) {
	if v, ok := c.cache.Get(cipherText); ok {
		return v, nil
	}

	cipherByte, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	decryptedText := make([]byte, len([]byte(cipherByte[aes.BlockSize:])))
	decryptStream := cipher.NewCTR(c.block, []byte(cipherByte[:aes.BlockSize]))
	decryptStream.XORKeyStream(decryptedText, []byte(cipherByte[aes.BlockSize:]))

	c.cache.Set(cipherText, string(decryptedText))
	return string(decryptedText), nil
}
