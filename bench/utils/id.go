package utils

import "github.com/oklog/ulid/v2"

// ULIDを生成する
func GenerateID() string {
	return ulid.Make().String()
}
