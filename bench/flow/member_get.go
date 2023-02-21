package flow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/grader"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/validator"
)

func (c *FlowController) memberGetFlow(memberID string, encrypt bool, step *isucandar.BenchmarkStep) flow {
	if memberID == "" {
		step.AddError(fmt.Errorf("GET /api/member/:id: %w", failure.NewError(model.ErrCritical, fmt.Errorf("memberID is empty"))))
	}

	findable := false
	if _, err := c.mr.GetMemberByID(memberID); err == nil {
		findable = true
	}

	if encrypt {
		var err error
		memberID, err = c.cr.Encrypt(memberID)
		if err != nil {
			step.AddError(fmt.Errorf("GET /api/member/:id: %w", failure.NewError(model.ErrCritical, err)))
			return nil
		}
	}

	return func(ctx context.Context) {
		res, err := c.ma.GetMember(ctx, memberID, encrypt)
		if err != nil {
			step.AddError(fmt.Errorf("GET /api/member/%s: %w", memberID, err))
			return
		}

		if findable {
			err = validator.Validate(res,
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
				step.AddError(fmt.Errorf("GET /api/members/%s: %w", memberID, err))
				return
			}
		} else {
			err = validator.Validate(res,
				validator.WithStatusCode(http.StatusNotFound),
			)
			if err != nil {
				step.AddError(fmt.Errorf("GET /api/member/%s: %w", memberID, err))
				return
			}
		}

		step.AddScore(grader.ScoreGetMember)
	}
}
