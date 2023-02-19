package validator

import (
	"fmt"
	"net/http"

	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

type SliceJsonValidationOpt[V comparable] func(body []V) error

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

func SliceJsonLengthEquals[V comparable](length int) SliceJsonValidationOpt[V] {
	return func(body []V) error {
		if len(body) != length {
			return failure.NewError(model.ErrInvalidBody,
				fmt.Errorf("expected: %v, actual: %v", length, len(body)))
		}
		return nil
	}
}

// 長さの範囲 [min, max]
func SliceJsonLengthRange[V comparable](min, max int) SliceJsonValidationOpt[V] {
	return func(body []V) error {
		if len(body) < min || len(body) > max {
			return failure.NewError(model.ErrInvalidBody,
				fmt.Errorf("expected: %v, actual: %v", min, len(body)))
		}
		return nil
	}
}

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
