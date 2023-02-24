package flow

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/action"
	"github.com/logica0419/gasshuku-isucon/bench/grader"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/validator"
)

const bookPageLimit = 50

func (c *Controller) searchBooksFlow(step *isucandar.BenchmarkStep) flow {
	return func(ctx context.Context) {
		book := c.br.GetRandomBook()

		q := action.GetBooksQuery{}
		if rand.Intn(2) == 0 {
			titleSlice := strings.Split(book.Title, " ")
			if len(titleSlice) > 1 {
				q.Title = titleSlice[rand.Intn(len(titleSlice)-1)]
			} else {
				titleSlice = strings.Split(book.Title, "")
				start := rand.Intn(len(titleSlice) - 1)
				end := rand.Intn(len(titleSlice)-start) + start + 1
				q.Title = strings.Join(titleSlice[start:end], "")
			}
		}
		if rand.Intn(2) == 0 {
			authorSlice := strings.Split(book.Author, " ")
			if len(authorSlice) > 1 {
				q.Author = authorSlice[rand.Intn(len(authorSlice)-1)]
			} else {
				authorSlice = strings.Split(book.Author, "")
				start := rand.Intn(len(authorSlice) - 1)
				end := rand.Intn(len(authorSlice)-start) + start + 1
				q.Author = strings.Join(authorSlice[start:end], "")
			}
		}
		if (q.Title == "" && q.Author == "") || rand.Intn(2) == 0 {
			q.Genre = book.Genre
		} else {
			q.Genre = model.Genre(-1)
		}

		page := 1
		lastBookID := ""

		for {
			q.Page = page
			q.LastBookID = lastBookID

			res, err := c.ba.GetBooks(ctx, q)
			if err != nil {
				step.AddError(fmt.Errorf("GET /api/books: %w", err))
				return
			}

			found := false

			err = validator.Validate(res,
				validator.WithStatusCode(http.StatusOK),
				validator.WithContentType("application/json"),
				validator.WithJsonValidation(
					validator.JsonSliceFieldValidate[action.GetBooksResponse]("Books",
						validator.SliceJsonLengthRange[model.BookWithLending](1, bookPageLimit),
						validator.SliceJsonCheckOrder(func(v model.BookWithLending) string { return v.ID }, validator.Asc),
						validator.SliceJsonCheckEach(func(body model.BookWithLending) error {
							if body.ID == book.ID {
								found = true
							}

							if q.Title != "" && !strings.Contains(body.Title, q.Title) {
								return failure.NewError(model.ErrInvalidBody, errors.New("invalid title"))
							}
							if q.Author != "" && !strings.Contains(body.Author, q.Author) {
								return failure.NewError(model.ErrInvalidBody, errors.New("invalid author"))
							}
							if q.Genre >= 0 && body.Genre != q.Genre {
								return failure.NewError(model.ErrInvalidBody, errors.New("invalid genre"))
							}

							v, err := c.br.GetBookByID(body.ID)
							if err != nil {
								return nil
							}
							return validator.JsonEquals(v.Book)(body.Book)
						}),
						func(body []model.BookWithLending) error {
							lastBookID = body[len(body)-1].ID
							return nil
						},
					),
					validator.JsonFieldValidate[action.GetBooksResponse]("Total",
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
				step.AddError(fmt.Errorf("GET /api/books: %w", err))
				return
			}

			if found {
				break
			}

			page++
		}

		step.AddScore(grader.ScoreSearchBooks)
	}
}
