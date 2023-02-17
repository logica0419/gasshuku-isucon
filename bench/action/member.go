package action

import (
	"context"

	"github.com/isucon/isucandar/agent"
	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

type GetMembersQueries struct {
	Page         int
	LastMemberID string
	Order        string
}

func GetMembers(ctx context.Context, a agent.Agent) []string {}

type PostMemberRequest struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

func PostMember(ctx context.Context, a agent.Agent, body PostMemberRequest) (model.Member, error) {
	reader, err := utils.StructToReader(body)
	if err != nil {
		return model.Member{}, failure.NewError(model.ErrCritical, err)
	}

	req, err := a.POST("/member", reader)
	if err != nil {
		return model.Member{}, failure.NewError(model.ErrCritical, err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := a.Do(ctx, req)
	if err != nil {
		return model.Member{}, err
	}
}
