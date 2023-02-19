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

func (c *FlowController) InitializeFlow(step *isucandar.BenchmarkStep) worker.WorkerFunc {
	return func(ctx context.Context, _ int) {
		res, b, err := c.ia.Initialize(ctx, c.key)
		if model.IsErrTimeout(err) {
			step.AddError(fmt.Errorf("GET /api/members: %w", failure.NewError(model.ErrTimeout, nil)))
			return
		}
		if err != nil {
			step.AddError(fmt.Errorf("GET /api/members: %w", failure.NewError(model.ErrRequestFailed, err)))
			return
		}

		err = validator.Validate(res, b,
			validator.WithStatusCode(http.StatusOK),
			validator.WithJsonValidation(
				validator.JsonEquals(action.InitializeHandlerResponse{
					Language: "Go",
				}),
			),
		)
		if err != nil {
			step.AddError(fmt.Errorf("POST /api/initialize: %w", err))
			return
		}
	}
}
