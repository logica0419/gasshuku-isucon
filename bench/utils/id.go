package utils

import "github.com/oklog/ulid/v2"

func GenerateID() string {
	return ulid.Make().String()
}
