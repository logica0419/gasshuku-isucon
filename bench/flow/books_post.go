package flow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/action"
	"github.com/logica0419/gasshuku-isucon/bench/grader"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/validator"
)

func (c *Controller) postBooksFlow(num int, step *isucandar.BenchmarkStep) flow {
	return func(ctx context.Context) {
		req := []action.PostBooksRequest{}

		for i := 0; i < num; i++ {
			title := model.NewBookTitle()
			author := model.NewBookAuthor()
			genre := model.NewBookGenre()
			req = append(req, action.PostBooksRequest{
				Title:  title,
				Author: author,
				Genre:  genre,
			})
		}

		res, err := c.ba.PostBooks(ctx, req)
		if err != nil {
			step.AddError(fmt.Errorf("POST /api/books: %w", err))
			return
		}

		var books []*model.BookWithLending

		err = validator.Validate(res,
			validator.WithStatusCode(http.StatusCreated),
			validator.WithContentType("application/json"),
			validator.WithSliceJsonValidation(
				validator.SliceJsonCheckEach(
					func(body model.Book) error {
						books = append(books, &model.BookWithLending{
							Book: body,
						})

						for _, r := range req {
							if r.Title == body.Title && r.Author == body.Author && r.Genre == body.Genre {
								return nil
							}
						}
						return failure.NewError(model.ErrInvalidBody, nil)
					},
				),
			),
		)
		if err != nil {
			step.AddError(fmt.Errorf("POST /api/books: %w", err))
			return
		}

		for _, book := range books {
			res, err = c.ba.GetBookQRCode(ctx, book.ID)
			if err != nil {
				step.AddError(fmt.Errorf("GET /api/books/%s/qrcode: %w", book.ID, err))
				return
			}

			err = validator.Validate(res,
				validator.WithStatusCode(http.StatusOK),
				validator.WithContentType("image/png"),
				validator.WithQRCodeEqual(book.ID, c.cr.Decrypt),
			)
			if err != nil {
				step.AddError(fmt.Errorf("GET /api/books/%s/qrcode: %w", book.ID, err))
				return
			}
		}
		c.br.AddBooks(books)

		step.AddScore(grader.ScorePostBooks)
	}
}
