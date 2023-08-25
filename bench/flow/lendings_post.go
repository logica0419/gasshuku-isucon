package flow

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/grader"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/validator"
)

func (c *Controller) postLendingFlow(memberID string, num int, step *isucandar.BenchmarkStep) flow {
	return func(ctx context.Context) {
		_, err := c.mr.GetMemberByID(memberID)
		if err != nil {
			step.AddError(fmt.Errorf("GET /api/member/%s: %w", memberID, failure.NewError(model.ErrCritical, err)))
			return
		}

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			c.getMemberFlow(memberID, true, step)(ctx)
			wg.Done()
		}()

		books, err := c.br.GetNotLendingBooks(num)
		if err != nil {
			return
		}

		bookIDs := []string{}
		for _, book := range books {
			book := book
			wg.Add(1)
			go func() {
				c.getBookFlow(book.ID, true, step)(ctx)
				wg.Done()
			}()
			bookIDs = append(bookIDs, book.ID)
		}

		wg.Wait()

		res, err := c.la.PostLendings(ctx, memberID, bookIDs)
		if err != nil {
			step.AddError(fmt.Errorf("POST /api/lendings: %w", err))
			return
		}

		lendings := []*model.LendingWithNames{}

		err = validator.Validate(res,
			validator.WithStatusCode(http.StatusCreated),
			validator.WithContentType("application/json"),
			validator.WithSliceJsonValidation(
				validator.SliceJsonCheckEach(
					func(body model.LendingWithNames) error {
						lendings = append(lendings, &body)

						for _, book := range books {
							if body.BookID == book.ID && body.MemberID == memberID && body.BookTitle == book.Title {
								return nil
							}
						}
						return failure.NewError(model.ErrInvalidBody, fmt.Errorf("invalid lending: %s", body.ID))
					},
				),
			),
		)
		if err != nil {
			step.AddError(fmt.Errorf("POST /api/lendings: %w", err))
			return
		}

		c.lr.AddLendings(lendings)

		step.AddScore(grader.ScorePostLendings)
	}
}
