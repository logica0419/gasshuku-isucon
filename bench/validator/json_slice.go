package validator

import (
	"fmt"
	"net/http"

	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
	"golang.org/x/exp/constraints"
)

// JSON配列ボディの検証をするための関数
type SliceJsonValidateOpt[V any] func(body []V) error

// JSON配列ボディのデコードと検証を行う
func WithSliceJsonValidation[V comparable](opt ...SliceJsonValidateOpt[V]) ValidateOpt {
	return func(res *http.Response) error {
		var body []V
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

// JSON配列ボディの長さが期待する値と等しいか検証
func SliceJsonLengthEquals[V any](length int) SliceJsonValidateOpt[V] {
	return func(body []V) error {
		if len(body) != length {
			return failure.NewError(model.ErrInvalidBody,
				fmt.Errorf("expected: %v, actual: %v", length, len(body)))
		}
		return nil
	}
}

// JSON配列ボディの長さが期待する値の範囲内か検証
func SliceJsonLengthRange[V any](min, max int) SliceJsonValidateOpt[V] {
	return func(body []V) error {
		if len(body) < min || len(body) > max {
			return failure.NewError(model.ErrInvalidBody,
				fmt.Errorf("expected: %v, actual: %v", min, len(body)))
		}
		return nil
	}
}

// 全てのJSON配列ボディの要素に対して検証を行う
func SliceJsonCheckEach[V any](f JsonValidateOpt[V]) SliceJsonValidateOpt[V] {
	return func(body []V) error {
		for _, v := range body {
			if err := f(v); err != nil {
				return err
			}
		}
		return nil
	}
}

type order string

const (
	Asc  order = "asc"  // 昇順
	Desc order = "desc" // 降順
)

// JSON配列ボディの要素が指定した順序でソートされているか検証
func SliceJsonCheckOrder[V any, I constraints.Ordered](idxFunc func(v V) I, ord order) SliceJsonValidateOpt[V] {
	return func(body []V) error {
		if len(body) < 2 {
			return nil
		}

		idxList := make([]I, len(body))
		for i, v := range body {
			idxList[i] = idxFunc(v)
		}

		for i := 0; i < len(idxList)-1; i++ {
			switch ord {
			case Asc:
				if idxList[i] > idxList[i+1] {
					return failure.NewError(model.ErrInvalidBody, fmt.Errorf("expected: %s", ord))
				}

			case Desc:
				if idxList[i] < idxList[i+1] {
					return failure.NewError(model.ErrInvalidBody, fmt.Errorf("expected: %s", ord))
				}
			}
		}

		return nil
	}
}
