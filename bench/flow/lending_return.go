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

func (c *Controller) returnLendingsFlow(memberID string, step *isucandar.BenchmarkStep) flow {
	return func(ctx context.Context) {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			c.getMemberFlow(memberID, true, step)(ctx)
			wg.Done()
		}()

		lendings, err := c.lr.GetLendingsByMemberID(memberID)
		if err != nil {
			step.AddError(fmt.Errorf("POST /api/lendings/return: %w", failure.NewError(model.ErrCritical, err)))
			return
		}

		bookIDs := []string{}
		for _, lending := range lendings {
			bookID := lending.BookID
			bookIDs = append(bookIDs, bookID)
			wg.Add(1)
			go func() {
				c.getBookFlow(bookID, true, step)(ctx)
				wg.Done()
			}()
		}

		wg.Wait()

		res, err := c.la.ReturnLendings(ctx, memberID, bookIDs)
		if err != nil {
			step.AddError(fmt.Errorf("POST /api/lendings/return: %w", err))
			return
		}

		err = validator.Validate(res,
			validator.WithStatusCode(http.StatusNoContent),
		)
		if err != nil {
			step.AddError(fmt.Errorf("POST /api/lendings/return: %w", err))
			return
		}

		c.lr.DeleteLendings(memberID)

		step.AddScore(grader.ScoreReturnLendings)
	}
}
