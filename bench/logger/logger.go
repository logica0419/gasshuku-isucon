package logger

import (
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
	Contestant = log.New(os.Stdout, "", log.Ltime)
	Admin = log.New(os.Stderr, "[Admin] ", log.Lmicroseconds)
}
