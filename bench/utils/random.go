package utils

import (
	"math/rand"
)

var (
	codes         = make([]rune, 0)
	codesWithSign = make([]rune, 0)
)

func init() {
	initializeCharCode()
	initializeCharCodeWithSign()
}

// 英数字をリストに登録
func initializeCharCode() {
	for i := '0'; i <= '9'; i++ {
		codes = append(codes, i)
	}
	for i := 'a'; i <= 'z'; i++ {
		codes = append(codes, i)
	}
	for i := 'A'; i <= 'Z'; i++ {
		codes = append(codes, i)
	}
}

// ", '. `, \ を覗く英数字と記号をリストに登録
//
//	ref: https://github.com/githayu/apps.hayu.io/blob/master/src/app/random/unicode-blocks.ts
//	(大本: https://ja.wikipedia.org/wiki/%E3%83%96%E3%83%AD%E3%83%83%E3%82%AF_(Unicode))
func initializeCharCodeWithSign() {
	for i := '0'; i <= '9'; i++ {
		codesWithSign = append(codesWithSign, i)
	}
	for i := 'a'; i <= 'z'; i++ {
		codesWithSign = append(codesWithSign, i)
	}
	for i := 'A'; i <= 'Z'; i++ {
		codesWithSign = append(codesWithSign, i)
	}
	codesWithSign = append(codesWithSign, rune(33))
	for i := 35; i <= 38; i++ {
		codesWithSign = append(codesWithSign, rune(i))
	}
	for i := 40; i <= 47; i++ {
		codesWithSign = append(codesWithSign, rune(i))
	}
	for i := 58; i <= 64; i++ {
		codesWithSign = append(codesWithSign, rune(i))
	}
	codesWithSign = append(codesWithSign, rune(91))
	for i := 93; i <= 95; i++ {
		codesWithSign = append(codesWithSign, rune(i))
	}
	for i := 123; i <= 126; i++ {
		codesWithSign = append(codesWithSign, rune(i))
	}
}

// 英数字のランダム文字列を生成
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(codes[rand.Intn(len(codes))])
	}
	return string(b)
}

// 英数字と記号のランダム文字列を生成
func RandStringWithSign(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(codesWithSign[rand.Intn(len(codesWithSign))])
	}
	return string(b)
}
