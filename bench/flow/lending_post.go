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

func (c *Controller) postLendingsFlow(num int, step *isucandar.BenchmarkStep) flow {
	return func(ctx context.Context) {
		member, err := c.mr.GetNotLendingMember()
		if err != nil {
			return
		}

		books, err := c.br.GetNotLendingBooks(num)
		if err != nil {
			return
		}

		bookIDs := []string{}
		for _, book := range books {
			bookIDs = append(bookIDs, book.ID)
		}

		res, err := c.la.PostLendings(ctx, member.ID, bookIDs)
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
							if body.BookID == book.ID && body.MemberID == member.ID &&
								body.BookTitle == book.Title && body.MemberName == member.Name {
								return nil
							}
						}
						return failure.NewError(model.ErrInvalidBody, nil)
					},
				),
			),
		)
		if err != nil {
			step.AddError(fmt.Errorf("POST /api/lendings: %w", err))
			return
		}

		c.lr.AddLendings(lendings)

		step.AddScore(grader.ScorePostBooks)
	}
}
