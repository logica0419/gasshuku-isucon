package action

import (
	"context"
	"io"
	"net/http"

	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

type InitializeActionController interface {
	Initialize(ctx context.Context, key string) (*http.Response, []byte, error)
}

type InitializeHandlerRequest struct {
	Key string `json:"key"`
}

type InitializeHandlerResponse struct {
	Language string `json:"language"`
}

// POST /api/initialize
func (c *ActionController) Initialize(ctx context.Context, key string) (*http.Response, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, c.initializeTimeout)
	defer cancel()

	agent := c.initializeAgent

	body, err := utils.EncodeJson(InitializeHandlerRequest{
		Key: key,
	})
	if err != nil {
		return nil, nil, failure.NewError(model.ErrCritical, err)
	}

	req, err := agent.POST("/api/initialize", body)
	if err != nil {
		return nil, nil, failure.NewError(model.ErrCritical, err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := agent.Do(ctx, req)
	if err != nil {
		return nil, nil, failure.NewError(model.ErrRequestFailed, err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, failure.NewError(model.ErrUndecodableBody, err)
	}

	return res, b, nil
}
