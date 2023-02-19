package validator

import (
	"fmt"
	"reflect"

	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
)

func JsonFieldValidate[V, F any](name string, opts ...JsonValidateOpt[F]) JsonValidateOpt[V] {
	return func(body V) error {
		bodyRv := reflect.ValueOf(body)
		fieldI := bodyRv.FieldByName(name).Interface()

		field, ok := fieldI.(F)
		if !ok {
			return failure.NewError(model.ErrInvalidBody, fmt.Errorf("unable to get field %s", name))
		}

		for _, o := range opts {
			if err := o(field); err != nil {
				return err
			}
		}
		return nil
	}
}

func JsonSliceFieldValidate[V, F any](name string, opts ...SliceJsonValidateOpt[F]) JsonValidateOpt[V] {
	return func(body V) error {
		bodyRv := reflect.ValueOf(body)
		fieldI := bodyRv.FieldByName(name).Interface()

		field, ok := fieldI.([]F)
		if !ok {
			return failure.NewError(model.ErrInvalidBody, fmt.Errorf("unable to get field %s", name))
		}

		for _, o := range opts {
			if err := o(field); err != nil {
				return err
			}
		}
		return nil
	}
}
