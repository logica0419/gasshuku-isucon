package flow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/validator"
)

func (c *FlowController) memberGetFlow(memberID string, encrypt bool, step *isucandar.BenchmarkStep) flow {
	if encrypt {
		memberID, _ = c.cr.Encrypt(memberID)
	}

	return func(ctx context.Context) {
		res, b, err := c.ma.GetMember(ctx, memberID, encrypt)
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
			validator.WithContentType("application/json"),
			validator.WithJsonValidation(
				func(body model.Member) error {
					v, err := c.mr.GetMemberByID(body.ID)
					if err != nil {
						return failure.NewError(model.ErrInvalidBody, err)
					}
					return validator.JsonEquals(v.Member)(body)
				}),
		)
		if err != nil {
			step.AddError(fmt.Errorf("GET /api/members: %w", err))
			return
		}
	}
}
