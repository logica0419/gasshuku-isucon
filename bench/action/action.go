package action

import (
	"context"
	"fmt"
	"time"

	"github.com/isucon/isucandar/agent"
	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/config"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

const (
	libAgentsNum    = 10 // 図書館職員エージェント数
	searchAgentsNum = 10 // 検索端末エージェント数
)

// Actionパッケージ内でしか使わないものを管理する構造体
//
//	Agentsは初期化以降ReadOnlyなため、ロックを取る機構は用意していない
type Controller struct {
	initializeAgent *agent.Agent

	initializeTimeout time.Duration
	requestTimeout    time.Duration

	libAgents    []utils.Choice[*agent.Agent]
	searchAgents []utils.Choice[*agent.Agent]
}

func NewController(c *config.Config) (*Controller, error) {
	initializeAgent, err := agent.NewAgent(agent.WithBaseURL(c.BaseURL), agent.WithDefaultTransport(),
		agent.WithTimeout(time.Duration(c.InitializeTimeout)*time.Millisecond))
	if err != nil {
		return nil, failure.NewError(model.ErrCritical, err)
	}
	initializeAgent.Name = "Isulibrary-InitializeAgent"

	libAgents := make([]utils.Choice[*agent.Agent], libAgentsNum)
	searchAgents := make([]utils.Choice[*agent.Agent], searchAgentsNum)

	for i := 0; i < libAgentsNum; i++ {
		libAgents[i].Weight = 1
		libAgents[i].Val, err = agent.NewAgent(agent.WithBaseURL(c.BaseURL), agent.WithDefaultTransport(),
			agent.WithTimeout(time.Duration(c.RequestTimeout)*time.Millisecond))
		if err != nil {
			return nil, failure.NewError(model.ErrCritical, err)
		}
		libAgents[i].Val.Name = fmt.Sprintf("Isulibrary-LibAgent-%d", i+1)
	}
	for i := 0; i < searchAgentsNum; i++ {
		searchAgents[i].Weight = 1
		searchAgents[i].Val, err = agent.NewAgent(agent.WithBaseURL(c.BaseURL), agent.WithDefaultTransport(),
			agent.WithTimeout(time.Duration(c.RequestTimeout)*time.Millisecond))
		if err != nil {
			return nil, failure.NewError(model.ErrCritical, err)
		}
		searchAgents[i].Val.Name = fmt.Sprintf("Isulibrary-SearchAgent-%d", i+1)
	}

	return &Controller{
		initializeAgent:   initializeAgent,
		initializeTimeout: time.Duration(c.InitializeTimeout) * time.Millisecond,
		requestTimeout:    time.Duration(c.RequestTimeout) * time.Millisecond,
		libAgents:         libAgents,
		searchAgents:      searchAgents,
	}, nil
}

func (c *Controller) libAgent() *agent.Agent {
	a, _ := utils.WeightedSelect(c.libAgents, false)
	return a
}

func (c *Controller) searchAgent() *agent.Agent {
	a, _ := utils.WeightedSelect(c.searchAgents, false)
	return a
}

func processErr(ctx context.Context, err error) error {
	select {
	case <-ctx.Done():
		return failure.NewError(model.ErrDeadline, err)
	default:
		if model.IsErrTimeout(err) {
			return failure.NewError(model.ErrTimeout, err)
		}
		return failure.NewError(model.ErrRequestFailed, err)
	}
}
