package validator

import (
	"fmt"
	"net/http"

	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

// JSONボディの検証をするための関数
type JsonValidateOpt[V any] func(body V) error

// JSONボディのデコードと検証を行う
func WithJsonValidation[V any](opt ...JsonValidateOpt[V]) ValidateOpt {
	return func(res *http.Response) error {
		var body V
		if err := utils.DecodeJsonWithStandard(res.Body, &body); err != nil {
			return failure.NewError(model.ErrUndecodableBody, err)
		}

		for _, o := range opt {
			if err := o(body); err != nil {
				return err
			}
		}
		return nil
	}
}

// JSONボディが期待する値と等しいか検証
func JsonEquals[V comparable](v V) JsonValidateOpt[V] {
	return func(body V) error {
		if body != v {
			return failure.NewError(model.ErrInvalidBody,
				fmt.Errorf("expected: %v, actual: %v", v, body))
		}
		return nil
	}
}
