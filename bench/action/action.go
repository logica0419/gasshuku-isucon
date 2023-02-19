package action

import (
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
type ActionController struct {
	initializeAgent *agent.Agent

	initializeTimeout time.Duration
	requestTimeout    time.Duration

	libAgents    []utils.Choice[*agent.Agent]
	searchAgents []utils.Choice[*agent.Agent]
}

func NewActionController(c config.Config) (*ActionController, error) {
	initializeAgent, err := agent.NewAgent(agent.WithBaseURL(c.BaseURL), agent.WithDefaultTransport())
	if err != nil {
		return nil, failure.NewError(model.ErrCritical, err)
	}
	initializeAgent.Name = "Isulibrary-InitializeAgent"

	libAgents := make([]utils.Choice[*agent.Agent], libAgentsNum)
	searchAgents := make([]utils.Choice[*agent.Agent], searchAgentsNum)

	for i := 0; i < libAgentsNum; i++ {
		libAgents[i].Val, err = agent.NewAgent(agent.WithBaseURL(c.BaseURL), agent.WithDefaultTransport())
		if err != nil {
			return nil, failure.NewError(model.ErrCritical, err)
		}
		libAgents[i].Val.Name = fmt.Sprintf("Isulibrary-LibAgent-%d", i+1)
		libAgents[i].Weight = 1
	}
	for i := 0; i < searchAgentsNum; i++ {
		searchAgents[i].Val, err = agent.NewAgent(agent.WithBaseURL(c.BaseURL), agent.WithDefaultTransport())
		if err != nil {
			return nil, failure.NewError(model.ErrCritical, err)
		}
		searchAgents[i].Val.Name = fmt.Sprintf("Isulibrary-SearchAgent-%d", i+1)
		searchAgents[i].Weight = 1
	}

	return &ActionController{
		initializeAgent:   initializeAgent,
		initializeTimeout: time.Duration(c.InitializeTimeout) * time.Millisecond,
		requestTimeout:    time.Duration(c.RequestTimeout) * time.Millisecond,
		libAgents:         libAgents,
		searchAgents:      searchAgents,
	}, nil
}

func (c *ActionController) libAgent() *agent.Agent {
	return utils.WeightedSelect(c.libAgents)
}

func (c *ActionController) searchAgent() *agent.Agent {
	return utils.WeightedSelect(c.searchAgents)
}
