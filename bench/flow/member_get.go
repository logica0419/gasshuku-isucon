package flow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/action"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
	"github.com/logica0419/gasshuku-isucon/bench/validator"
)

const memberPageLimit = 100

func (c *FlowController) membersGetFlow(step *isucandar.BenchmarkStep) flow {
	return func(ctx context.Context) {
		page := 1
		lastMemberID := ""
		order := utils.WeightedSelect(
			[]utils.Choice[string]{
				{Val: "", Weight: 2},
				{Val: "name_asc"},
				{Val: "name_desc"},
			},
		)

		idxFunc := func(v model.Member) string { return v.ID }
		ord := validator.Asc
		switch order {
		case "name_asc":
			idxFunc = func(v model.Member) string { return v.Name }
		case "name_desc":
			idxFunc = func(v model.Member) string { return v.Name }
			ord = validator.Desc
		}

		limit := false

		for {
			query := action.GetMembersQuery{
				Page:         page,
				LastMemberID: lastMemberID,
				Order:        order,
			}

			res, err := c.ma.GetMembers(ctx, query)
			if model.IsErrTimeout(err) {
				step.AddError(fmt.Errorf("GET /api/members: %w", failure.NewError(model.ErrTimeout, nil)))
				return
			}
			if err != nil {
				step.AddError(fmt.Errorf("GET /api/members: %w", failure.NewError(model.ErrRequestFailed, err)))
				return
			}

			if limit {
				err = validator.Validate(res,
					validator.WithStatusCode(http.StatusNotFound),
				)
				if err != nil {
					step.AddError(fmt.Errorf("GET /api/members: %w", err))
				}
				return
			}

			err = validator.Validate(res,
				validator.WithStatusCode(http.StatusOK),
				validator.WithSliceJsonValidation(
					validator.SliceJsonLengthRange[model.Member](1, memberPageLimit),
					validator.SliceJsonCheckOrder(idxFunc, ord),
					validator.SliceJsonCheckEach(func(body model.Member) error {
						v, err := c.mr.GetMemberByID(body.ID)
						if err != nil {
							return failure.NewError(model.ErrInvalidBody, err)
						}
						return validator.JsonEquals(v.Member)(body)
					}),
				),
			)
			if err != nil {
				step.AddError(fmt.Errorf("GET /api/members: %w", err))
				return
			}
		}
	}
}
