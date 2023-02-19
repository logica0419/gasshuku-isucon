package validator

import (
	"net/http"
)

// http.Responseを検証するための関数
type ValidateOpt func(*http.Response, []byte) error

// http.Responseを検証
//
//	渡されたoptsの順にvalidateし、どこかでValidateに失敗したらエラーを返す
func Validate(res *http.Response, body []byte, opts ...ValidateOpt) error {
	for _, opt := range opts {
		if err := opt(res, body); err != nil {
			return err
		}
	}
	return nil
}
