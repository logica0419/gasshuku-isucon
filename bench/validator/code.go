package validator

import (
	"fmt"
	"net/http"

	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
)

func WithStatusCode(code int) ValidateOpt {
	return func(res *http.Response) error {
		if res.StatusCode != code {
			return failure.NewError(model.ErrInvalidStatusCode,
				fmt.Errorf("different status code - expected: %d, actual: %d", code, res.StatusCode))
		}
		return nil
	}
}
