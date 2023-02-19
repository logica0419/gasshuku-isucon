package action

import (
	"context"
	"net/http"

	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

type InitializeActionController interface {
	Initialize(ctx context.Context, key string) (*http.Response, error)
}

type InitializeHandlerRequest struct {
	Key string `json:"key"`
}

type InitializeHandlerResponse struct {
	Language string `json:"language"`
}

func (c *ActionController) Initialize(ctx context.Context, key string) (*http.Response, error) {
	agent := c.initializeAgent

	body, err := utils.StructToReader(InitializeHandlerRequest{
		Key: key,
	})
	if err != nil {
		return nil, failure.NewError(model.ErrCritical, err)
	}

	req, err := agent.POST("/api/initialize", body)
	if err != nil {
		return nil, err
	}

	res, err := agent.Do(ctx, req)
	if err != nil {
		return nil, failure.NewError(model.ErrRequestFailed, err)
	}

	return res, nil
}
