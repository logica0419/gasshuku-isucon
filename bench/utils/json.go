package utils

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/bytedance/sonic"
)

// []byteをJSONとしてデコード
func DecodeJson(b []byte, s any) error {
	return sonic.Unmarshal(b, s)
}

// JSONとしてio.Readerにエンコード
func EncodeJson(s any) (io.Reader, error) {
	b, err := sonic.Marshal(s)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

// []byteをJSONとしてデコード
//
//	予期しないJSONが来る可能性があるので、標準パッケージでデコードする
func DecodeJsonWithStandard(b []byte, s any) error {
	return json.Unmarshal(b, s)
}
