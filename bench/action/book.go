package action

import (
	"context"
	"net/http"
	"strconv"

	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

type BookController interface {
	PostBooks(ctx context.Context, body []PostBooksRequest) (*http.Response, error)
	GetBooks(ctx context.Context, query GetBooksQuery) (*http.Response, error)
	GetBook(ctx context.Context, id string, encrypted bool) (*http.Response, error)
	GetBookQRCode(ctx context.Context, id string) (*http.Response, error)
}

var _ BookController = &Controller{}

type PostBooksRequest struct {
	Title  string      `json:"title"`
	Author string      `json:"author"`
	Genre  model.Genre `json:"genre"`
}

// POST /api/books
func (c *Controller) PostBooks(ctx context.Context, body []PostBooksRequest) (*http.Response, error) {
	reader, err := utils.EncodeJson(body)
	if err != nil {
		return nil, failure.NewError(model.ErrCritical, err)
	}

	agent := c.libAgent()

	req, err := agent.POST("/api/books", reader)
	if err != nil {
		return nil, failure.NewError(model.ErrCritical, err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := agent.Do(ctx, req)
	if err != nil {
		return nil, processErr(ctx, err)
	}

	return res, nil
}

type GetBooksQuery struct {
	Title      string
	Author     string
	Genre      model.Genre
	Page       int
	LastBookID string
}

type GetBooksResponse struct {
	Books []model.BookWithLending `json:"books"`
	Total int                     `json:"total"`
}

// GET /api/books
func (c *Controller) GetBooks(ctx context.Context, query GetBooksQuery) (*http.Response, error) {
	agent := c.searchAgent()

	url := "/api/books?"
	if query.Title != "" {
		url += "title=" + query.Title + "&"
	}
	if query.Author != "" {
		url += "author=" + query.Author + "&"
	}
	if query.Genre >= 0 {
		url += "genre=" + query.Genre.String() + "&"
	}
	if query.Page > 1 {
		url += "page=" + strconv.Itoa(query.Page) + "&"
	}
	if query.LastBookID != "" {
		url += "last_book_id=" + query.LastBookID + "&"
	}
	url = url[:len(url)-1] // 最後の一文字(?か&)を削除する

	req, err := agent.GET(url)
	if err != nil {
		return nil, failure.NewError(model.ErrCritical, err)
	}

	res, err := agent.Do(ctx, req)
	if err != nil {
		return nil, processErr(ctx, err)
	}

	return res, nil
}

// GET /api/books/:id
func (c *Controller) GetBook(ctx context.Context, id string, encrypted bool) (*http.Response, error) {
	agent := c.libAgent()

	url := "/api/books/" + id
	if encrypted {
		url += "?encrypted=true"
	}

	req, err := agent.GET(url)
	if err != nil {
		return nil, failure.NewError(model.ErrCritical, err)
	}

	res, err := agent.Do(ctx, req)
	if err != nil {
		return nil, processErr(ctx, err)
	}

	return res, nil
}

// GET /api/members/:id/qrcode
func (c *Controller) GetBookQRCode(ctx context.Context, id string) (*http.Response, error) {
	agent := c.libAgent()

	req, err := agent.GET("/api/books/" + id + "/qrcode")
	if err != nil {
		return nil, failure.NewError(model.ErrCritical, err)
	}

	res, err := agent.Do(ctx, req)
	if err != nil {
		return nil, processErr(ctx, err)
	}

	return res, nil
}
