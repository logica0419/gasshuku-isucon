package validator

import (
	"fmt"
	"net/http"

	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

// JSON配列ボディの検証をするための関数
type SliceJsonValidationOpt[V comparable] func(body []V) error

// JSON配列ボディのデコードと検証を行う
func WithSliceJsonValidation[V comparable](opt ...SliceJsonValidationOpt[V]) ValidateOpt {
	return func(res *http.Response) error {
		var body []V
		if err := utils.ReaderToStruct(res.Body, &body); err != nil {
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

// JSON配列ボディの長さが期待する値と等しいか検証
func SliceJsonLengthEquals[V comparable](length int) SliceJsonValidationOpt[V] {
	return func(body []V) error {
		if len(body) != length {
			return failure.NewError(model.ErrInvalidBody,
				fmt.Errorf("expected: %v, actual: %v", length, len(body)))
		}
		return nil
	}
}

// JSON配列ボディの長さが期待する値の範囲内か検証
func SliceJsonLengthRange[V comparable](min, max int) SliceJsonValidationOpt[V] {
	return func(body []V) error {
		if len(body) < min || len(body) > max {
			return failure.NewError(model.ErrInvalidBody,
				fmt.Errorf("expected: %v, actual: %v", min, len(body)))
		}
		return nil
	}
}

// 全てのJSON配列ボディの要素に対して検証を行う
func SliceJsonCheckEach[V comparable](f JsonValidateOpt[V]) SliceJsonValidationOpt[V] {
	return func(body []V) error {
		for _, v := range body {
			if err := f(v); err != nil {
				return err
			}
		}
		return nil
	}
}
