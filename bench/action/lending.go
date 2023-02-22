package action

import (
	"context"
	"net/http"

	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

type LendingController interface {
	PostLendings(ctx context.Context, memberID string, bookIDs []string) (*http.Response, error)
	GetLendings(ctx context.Context, overDue bool) (*http.Response, error)
	ReturnLendings(ctx context.Context, memberID string, bookIDs []string) (*http.Response, error)
}

var _ LendingController = &Controller{}

type PostLendingsRequest struct {
	BookIDs  []string `json:"book_ids"`
	MemberID string   `json:"member_id"`
}

// POST /api/lendings
func (c *Controller) PostLendings(ctx context.Context, memberID string, bookIDs []string) (*http.Response, error) {
	body := PostLendingsRequest{
		BookIDs:  bookIDs,
		MemberID: memberID,
	}

	reader, err := utils.EncodeJson(body)
	if err != nil {
		return nil, failure.NewError(model.ErrCritical, err)
	}

	agent := c.libAgent()

	req, err := agent.POST("/api/lendings", reader)
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

// GET /api/lendings
func (c *Controller) GetLendings(ctx context.Context, overDue bool) (*http.Response, error) {
	agent := c.libAgent()

	url := "/api/lendings"
	if overDue {
		url += "?over_due=true"
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

type ReturnLendingsRequest struct {
	BookIDs  []string `json:"book_ids"`
	MemberID string   `json:"member_id"`
}

// POST /api/lendings/return
func (c *Controller) ReturnLendings(ctx context.Context, memberID string, bookIDs []string) (*http.Response, error) {
	body := PostLendingsRequest{
		BookIDs:  bookIDs,
		MemberID: memberID,
	}

	reader, err := utils.EncodeJson(body)
	if err != nil {
		return nil, failure.NewError(model.ErrCritical, err)
	}

	agent := c.libAgent()

	req, err := agent.POST("/api/lendings/return", reader)
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
