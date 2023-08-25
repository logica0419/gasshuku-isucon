package flow

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/action"
	"github.com/logica0419/gasshuku-isucon/bench/grader"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
	"github.com/logica0419/gasshuku-isucon/bench/validator"
)

const memberPageLimit = 100

func (c *Controller) getMembersFlow(memberID string, step *isucandar.BenchmarkStep) flow {
	findable := false
	if memberID != "" {
		if _, err := c.mr.GetMemberByID(memberID); err == nil {
			findable = true
		}
	}

	page := 1
	lastMemberID := ""
	order, _ := utils.WeightedSelect(
		[]utils.Choice[string]{
			{Val: "", Weight: 2},
			{Val: "name_asc", Weight: 1},
			{Val: "name_desc", Weight: 1},
		}, false,
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

	return func(ctx context.Context) {
		for {
			query := action.GetMembersQuery{
				Page:         page,
				LastMemberID: lastMemberID,
				Order:        order,
			}

			res, err := c.ma.GetMembers(ctx, query)
			if err != nil {
				step.AddError(fmt.Errorf("GET /api/members: %w", err))
				return
			}

			if res.StatusCode == http.StatusNotFound && page > 1 {
				break
			}

			found := false

			err = validator.Validate(res,
				validator.WithStatusCode(http.StatusOK),
				validator.WithContentType("application/json"),
				validator.WithJsonValidation(
					validator.JsonSliceFieldValidate[action.GetMembersResponse]("Members",
						validator.SliceJsonLengthRange[model.Member](1, memberPageLimit),
						validator.SliceJsonCheckOrder(idxFunc, ord),
						validator.SliceJsonCheckEach(func(body model.Member) error {
							if body.ID == memberID {
								found = true
							}

							v, err := c.mr.GetMemberByID(body.ID)
							if err != nil {
								return nil
							}
							if v.ID != body.ID {
								return failure.NewError(model.ErrInvalidBody, fmt.Errorf("member ID is not equal"))
							}
							if v.CreatedAt != body.CreatedAt {
								return failure.NewError(model.ErrInvalidBody, fmt.Errorf("created at is not equal"))
							}
							return nil
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

			if found {
				if !findable {
					step.AddError(fmt.Errorf("GET /api/members: %w", failure.NewError(model.ErrInvalidBody, errors.New("found invalid member"))))
				}
				break
			}

			page++
		}

		step.AddScore(grader.ScoreGetMembers)
	}
}
