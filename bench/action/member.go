package action

import (
	"context"
	"net/http"
	"strconv"

	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

type GetMembersQuery struct {
	Page         int
	LastMemberID string
	Order        string
}

func (c *ActionController) GetMembers(ctx context.Context, query GetMembersQuery) (*http.Response, error) {
	agent := c.libAgent()

	url := "/api/members?"
	if query.Page != 0 {
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
		return nil, failure.NewError(model.ErrRequestFailed, err)
	}

	return res, nil
}

type PostMemberRequest struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

func (c *ActionController) PostMember(ctx context.Context, body PostMemberRequest) (*http.Response, error) {
	reader, err := utils.StructToReader(body)
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
		return nil, failure.NewError(model.ErrRequestFailed, err)
	}

	return res, nil
}