package logger

import (
	"io"
	"log"
	"os"
)

var (
	// 競技者に見せてもいい内容を書くロガー
	Contestant *log.Logger
	// 運営だけが見れる内容を書くロガー
	Admin *log.Logger
)

func init() {
	contestantWriter := io.MultiWriter(os.Stdout, os.Stderr)
	Contestant = log.New(contestantWriter, "", log.Lmicroseconds)
	Admin = log.New(os.Stderr, "[Admin] ", log.Lmicroseconds)
}
