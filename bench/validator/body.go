package validator

import (
	"fmt"
	"net/http"

	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

func WithBody[V comparable](v V) ValidateOpt {
	return func(res *http.Response) error {
		var body V
		if err := utils.ReaderToStruct(res.Body, &body); err != nil {
			return failure.NewError(model.ErrUndecodableBody, err)
		}

		if body != v {
			return failure.NewError(model.ErrInvalidBody,
				fmt.Errorf("different body - expected: %v, actual: %v", v, body))
		}
		return nil
	}
}
