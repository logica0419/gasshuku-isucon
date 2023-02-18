package validator

import (
	"net/http"

	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

func WithBody[V comparable](v V) ValidateOpt {
	return func(res http.Response) error {
		var body V
		if err := utils.ReaderToStruct(res.Body, &body); err != nil {
			return nil
		}

		if body != v {
			return nil
		}
		return nil
	}
}
