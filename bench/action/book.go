package action

import (
	"context"
	"net/http"

	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

type BookController interface {
	PostBooks(ctx context.Context, body []PostBooksRequest) (*http.Response, error)
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
