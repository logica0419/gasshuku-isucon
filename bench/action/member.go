package action

import (
	"context"
	"net/http"
	"strconv"

	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

type MemberController interface {
	PostMember(ctx context.Context, body PostMemberRequest) (*http.Response, error)
	GetMembers(ctx context.Context, query GetMembersQuery) (*http.Response, error)
	GetMember(ctx context.Context, id string, encrypted bool) (*http.Response, error)
	BanMember(ctx context.Context, id string) (*http.Response, error)
	PatchMember(ctx context.Context, id string, body PatchMemberRequest) (*http.Response, error)
	GetMemberQRCode(ctx context.Context, id string) (*http.Response, error)
}

var _ MemberController = &Controller{}

type PostMemberRequest struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

// POST /api/members
func (c *Controller) PostMember(ctx context.Context, body PostMemberRequest) (*http.Response, error) {
	reader, err := utils.EncodeJson(body)
	if err != nil {
		return nil, failure.NewError(model.ErrCritical, err)
	}

	agent := c.libAgent()

	req, err := agent.POST("/api/members", reader)
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

type GetMembersQuery struct {
	Page         int
	LastMemberID string
	Order        string
}

type GetMembersResponse struct {
	Members []model.Member `json:"members"`
	Total   int            `json:"total"`
}

// GET /api/members
func (c *Controller) GetMembers(ctx context.Context, query GetMembersQuery) (*http.Response, error) {
	agent := c.libAgent()

	url := "/api/members?"
	if query.Page > 1 {
		url += "page=" + strconv.Itoa(query.Page) + "&"
	}
	if query.LastMemberID != "" {
		url += "last_member_id=" + query.LastMemberID + "&"
	}
	if query.Order != "" {
		url += "order=" + query.Order + "&"
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

// GET /api/members/:id
func (c *Controller) GetMember(ctx context.Context, id string, encrypted bool) (*http.Response, error) {
	agent := c.libAgent()

	url := "/api/members/" + id
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

// DELETE /api/members/:id
func (c *Controller) BanMember(ctx context.Context, id string) (*http.Response, error) {
	agent := c.libAgent()

	req, err := agent.DELETE("/api/members/"+id, nil)
	if err != nil {
		return nil, failure.NewError(model.ErrCritical, err)
	}

	res, err := agent.Do(ctx, req)
	if err != nil {
		return nil, processErr(ctx, err)
	}

	return res, nil
}

type PatchMemberRequest struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

// PATCH /api/members/:id
func (c *Controller) PatchMember(ctx context.Context, id string, body PatchMemberRequest) (*http.Response, error) {
	reader, err := utils.EncodeJson(body)
	if err != nil {
		return nil, failure.NewError(model.ErrCritical, err)
	}

	agent := c.libAgent()

	req, err := agent.PATCH("/api/members/"+id, reader)
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
func (c *Controller) GetMemberQRCode(ctx context.Context, id string) (*http.Response, error) {
	agent := c.libAgent()

	req, err := agent.GET("/api/members/" + id + "/qrcode")
	if err != nil {
		return nil, failure.NewError(model.ErrCritical, err)
	}

	res, err := agent.Do(ctx, req)
	if err != nil {
		return nil, processErr(ctx, err)
	}

	return res, nil
}
