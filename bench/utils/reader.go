package utils

import (
	"bytes"
	"io"

	"github.com/bytedance/sonic"
)

func StructToReader(s any) (io.Reader, error) {
	b, err := sonic.Marshal(s)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
