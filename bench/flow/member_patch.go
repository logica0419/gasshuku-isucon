package flow

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/isucon/isucandar"
	"github.com/logica0419/gasshuku-isucon/bench/action"
	"github.com/logica0419/gasshuku-isucon/bench/grader"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/repository"
	"github.com/logica0419/gasshuku-isucon/bench/validator"
)

func (c *Controller) patchMemberFlow(memberID string, step *isucandar.BenchmarkStep) flow {
	return func(ctx context.Context) {
		c.getMemberFlow(memberID, true, step)(ctx)

		req := action.PatchMemberRequest{}
		if rand.Intn(2) == 0 {
			req.Name = model.NewMemberName()
		}
		if rand.Intn(2) == 0 {
			req.Address = model.NewMemberAddress()
		}
		if (req.Name == "" && req.Address == "") || rand.Intn(2) == 0 {
			req.PhoneNumber = model.NewMemberPhoneNumber()
		}

		if _, err := c.mr.GetMemberByID(memberID); err != nil {
			return
		}

		res, err := c.ma.PatchMember(ctx, memberID, req)
		if err != nil {
			step.AddError(fmt.Errorf("PATCH /api/members/%s: %w", memberID, err))
			return
		}

		if res.StatusCode != http.StatusNotFound {
			return
		}

		err = validator.Validate(res,
			validator.WithStatusCode(http.StatusNoContent),
		)
		if err != nil {
			step.AddError(fmt.Errorf("PATCH /api/members/%s: %w", memberID, err))
			return
		}

		_ = c.mr.UpdateMember(memberID, repository.MemberUpdateQuery(req))

		step.AddScore(grader.ScoreUpdateMember)
	}
}
