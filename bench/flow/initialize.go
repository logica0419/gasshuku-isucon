package flow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/failure"
	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/action"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/validator"
)

var langList = []string{"Go"}

func (c *FlowController) InitializeFlow(step *isucandar.BenchmarkStep) worker.WorkerFunc {
	return func(ctx context.Context, _ int) {
		res, err := c.ia.Initialize(ctx, c.key)
		if model.IsErrTimeout(err) {
			step.AddError(fmt.Errorf("POST /api/initialize: %w", failure.NewError(model.ErrTimeout, err)))
			return
		}
		if err != nil {
			step.AddError(fmt.Errorf("POST /api/initialize: %w", failure.NewError(model.ErrRequestFailed, err)))
			return
		}

		err = validator.Validate(res,
			validator.WithStatusCode(http.StatusOK),
			validator.WithJsonValidation(
				validator.JsonFieldValidate[action.InitializeHandlerResponse]("Language",
					func(lang string) error {
						for _, implemented := range langList {
							if lang == implemented {
								return nil
							}
						}
						return fmt.Errorf("language not implemented: %s", lang)
					},
				),
			),
		)
		if err != nil {
			step.AddError(fmt.Errorf("POST /api/initialize: %w", err))
			return
		}
	}
}
