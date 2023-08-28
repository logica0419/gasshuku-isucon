package validator

import (
	"io"
	"net/http"
)

// http.Responseを検証するための関数
type ValidateOpt func(*http.Response) error

// http.Responseを検証
//
//	渡されたoptsの順にvalidateし、どこかでValidateに失敗したらエラーを返す
func Validate(res *http.Response, opts ...ValidateOpt) error {
	if res.Body != nil {
		// ボディが残っているとHTTP keep-aliveができないので読み捨てて閉じる
		defer func() {
			io.Copy(io.Discard, res.Body)
			res.Body.Close()
		}()
	}

	for _, opt := range opts {
		if err := opt(res); err != nil {
			return err
		}
	}
	return nil
}
