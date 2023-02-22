package flow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/isucon/isucandar"
	"github.com/logica0419/gasshuku-isucon/bench/action"
	"github.com/logica0419/gasshuku-isucon/bench/grader"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/validator"
)

func (c *Controller) postMemberFlow(num int, step *isucandar.BenchmarkStep) flow {
	return func(ctx context.Context) {
		members := []*model.MemberWithLending{}

		for i := 0; i < num; i++ {
			name := model.NewMemberName()
			address := model.NewMemberAddress()
			phoneNumber := model.NewMemberPhoneNumber()
			q := action.PostMemberRequest{
				Name:        name,
				Address:     address,
				PhoneNumber: phoneNumber,
			}

			res, err := c.ma.PostMember(ctx, q)
			if err != nil {
				step.AddError(fmt.Errorf("POST /api/members: %w", err))
				return
			}

			var id string

			err = validator.Validate(res,
				validator.WithStatusCode(http.StatusCreated),
				validator.WithContentType("application/json"),
				validator.WithJsonValidation(
					validator.JsonFieldValidate[model.Member]("Name",
						validator.JsonEquals(name),
					),
					validator.JsonFieldValidate[model.Member]("Address",
						validator.JsonEquals(address),
					),
					validator.JsonFieldValidate[model.Member]("PhoneNumber",
						validator.JsonEquals(phoneNumber),
					),
					validator.JsonFieldValidate[model.Member]("Banned",
						validator.JsonEquals(false),
					),
					func(m model.Member) error {
						id = m.ID

						members = append(members, &model.MemberWithLending{
							Member: m,
						})
						return nil
					},
				),
			)
			if err != nil {
				step.AddError(fmt.Errorf("POST /api/members: %w", err))
				return
			}

			res, err = c.ma.GetMemberQRCode(ctx, id)
			if err != nil {
				step.AddError(fmt.Errorf("GET /api/members/%s/qrcode: %w", id, err))
				return
			}

			err = validator.Validate(res,
				validator.WithStatusCode(http.StatusOK),
				validator.WithContentType("image/png"),
				validator.WithQRCodeEqual(id, c.cr.Decrypt),
			)
			if err != nil {
				step.AddError(fmt.Errorf("GET /api/members/%s/qrcode: %w", id, err))
				return
			}
		}
		c.mr.AddMembers(members)

		step.AddScore(grader.ScorePostMember)
	}
}
