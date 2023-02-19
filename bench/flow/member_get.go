package flow

import (
	"context"
	"errors"
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
			defer res.Body.Close()

			if res.StatusCode == http.StatusNotFound && page > 1 {
				return
			}

			err = validator.Validate(res,
				validator.WithStatusCode(http.StatusOK),
				validator.WithContentType("application/json"),
				validator.WithJsonValidation(
					validator.JsonSliceFieldValidate[action.GetMembersResponse]("Members",
						validator.SliceJsonLengthRange[model.Member](1, memberPageLimit),
						validator.SliceJsonCheckOrder(idxFunc, ord),
						validator.SliceJsonCheckEach(func(body model.Member) error {
							v, err := c.mr.GetMemberByID(body.ID)
							if err != nil {
								return failure.NewError(model.ErrInvalidBody, err)
							}
							return validator.JsonEquals(v.Member)(body)
						}),
						func(body []model.Member) error {
							lastMemberID = body[len(body)-1].ID
							return nil
						},
					),
					validator.JsonFieldValidate[action.GetMembersResponse]("Total",
						func(total int) error {
							if total <= 0 {
								return failure.NewError(model.ErrInvalidBody, errors.New("total is invalid"))
							}
							return nil
						},
					),
				),
			)
			if err != nil {
				step.AddError(fmt.Errorf("GET /api/members: %w", err))
				return
			}

			res.Body.Close()
			page++
		}
	}
}
