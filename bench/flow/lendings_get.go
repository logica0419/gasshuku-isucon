package flow

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/isucon/isucandar"
	"github.com/logica0419/gasshuku-isucon/bench/grader"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/validator"
	"golang.org/x/sync/errgroup"
)

const banPeriod = 3000

func (c *Controller) getLendingsFlow(step *isucandar.BenchmarkStep) flow {
	return func(ctx context.Context) {
		overdue := rand.Intn(4) != 0

		res, err := c.la.GetLendings(ctx, overdue)
		if err != nil {
			step.AddError(fmt.Errorf("GET /api/lendings: %w", err))
			return
		}

		now := time.Now()
		banUserIDs := []string{}

		err = validator.Validate(res,
			validator.WithStatusCode(http.StatusOK),
			validator.WithContentType("application/json"),
			validator.WithSliceJsonValidation(
				validator.SliceJsonCheckOrder(func(v model.LendingWithNames) string { return v.ID }, validator.Asc),
				validator.SliceJsonCheckEach(func(body model.LendingWithNames) error {
					if body.Due.Add(banPeriod * time.Millisecond).Before(now) {
						exist := false
						for _, u := range banUserIDs {
							if u == body.MemberID {
								exist = true
								break
							}
						}
						if !exist {
							banUserIDs = append(banUserIDs, body.MemberID)
						}
					}

					v, err := c.lr.GetLendingByID(body.ID)
					if err != nil {
						return nil
					}
					v.MemberName = body.MemberName
					return validator.JsonEquals(*v)(body)
				}),
			),
		)
		if err != nil {
			step.AddError(fmt.Errorf("GET /api/lendings: %w", err))
			return
		}

		step.AddScore(grader.ScoreGetLendings)

		eg := errgroup.Group{}
		for _, id := range banUserIDs {
			id := id

			eg.Go(func() error {
				res, err := c.ma.BanMember(ctx, id)
				if err != nil {
					step.AddError(fmt.Errorf("DELETE /api/members/%s: %w", id, err))
					return err
				}

				if res.StatusCode == http.StatusNotFound {
					return nil
				}

				err = validator.Validate(res,
					validator.WithStatusCode(http.StatusNoContent),
				)
				if err != nil {
					step.AddError(fmt.Errorf("DELETE /api/members/%s: %w", id, err))
					return err
				}

				c.mr.DeleteMember(id)
				step.AddScore(grader.ScoreBanMember)
				return nil
			})
		}
		err = eg.Wait()
		if err != nil {
			return
		}
	}
}
