package validator

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
)

// ステータスコードが期待する値と等しいか検証
func WithStatusCode(code int) ValidateOpt {
	return func(res *http.Response) error {
		if res.StatusCode != code {
			return failure.NewError(model.ErrInvalidStatusCode,
				fmt.Errorf("expected: %d, actual: %d", code, res.StatusCode))
		}
		return nil
	}
}

// Content-Typeヘッダーが期待する値と等しいか検証
func WithContentType(contentType string) ValidateOpt {
	return func(res *http.Response) error {
		actual := res.Header.Get("Content-Type")
		if !strings.HasPrefix(actual, contentType) {
			return failure.NewError(model.ErrInvalidContentType,
				fmt.Errorf("expected: %s, actual: %s", contentType, actual))
		}
		return nil
	}
}
