package flow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/validator"
)

func (c *Controller) getBookFlow(bookID string, encrypt bool, step *isucandar.BenchmarkStep) flow {
	if bookID == "" {
		step.AddError(fmt.Errorf("GET /api/books/:id: %w", failure.NewError(model.ErrCritical, fmt.Errorf("memberID is empty"))))
	}

	findable := false
	if _, err := c.br.GetBookByID(bookID); err == nil {
		findable = true
	}

	if encrypt {
		var err error
		bookID, err = c.cr.Encrypt(bookID)
		if err != nil {
			step.AddError(fmt.Errorf("GET /api/books/:id: %w", failure.NewError(model.ErrCritical, err)))
			return nil
		}
	}

	return func(ctx context.Context) {
		res, err := c.ba.GetBook(ctx, bookID, encrypt)
		if err != nil {
			step.AddError(fmt.Errorf("GET /api/books/%s: %w", bookID, err))
			return
		}

		if findable {
			err = validator.Validate(res,
				validator.WithStatusCode(http.StatusOK),
				validator.WithContentType("application/json"),
				validator.WithJsonValidation(
					func(body model.Book) error {
						v, err := c.br.GetBookByID(body.ID)
						if err != nil {
							return failure.NewError(model.ErrInvalidBody, err)
						}
						return validator.JsonEquals(v.Book)(body)
					}),
			)
			if err != nil {
				step.AddError(fmt.Errorf("GET /api/books/%s: %w", bookID, err))
				return
			}
		} else {
			err = validator.Validate(res,
				validator.WithStatusCode(http.StatusNotFound),
			)
			if err != nil {
				step.AddError(fmt.Errorf("GET /api/books/%s: %w", bookID, err))
				return
			}
		}
	}
}
