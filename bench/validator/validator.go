package validator

import (
	"net/http"
)

type ValidateOpt func(http.Response) error

// http.Responseを検証
// 渡されたoptsの順にvalidateし、どこかでValidateに失敗したらエラーを返す
func Validate(res http.Response, opts ...ValidateOpt) error {
	for _, opt := range opts {
		if err := opt(res); err != nil {
			return err
		}
	}
	return nil
}
