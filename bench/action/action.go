package action

import (
	"fmt"

	"github.com/isucon/isucandar/agent"
	"github.com/isucon/isucandar/failure"
	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

// Actionパッケージ内でしか使わないものを管理する構造体
// Agentsは初期化以降ReadOnlyなため、ロックを取る機構は用意していない
type ActionController struct {
	libAgents    []utils.Choice[*agent.Agent]
	searchAgents []utils.Choice[*agent.Agent]
}

func NewActionController(libAgentsNum int, searchAgentsNum int, baseURL string) (*ActionController, error) {
	libAgents := make([]utils.Choice[*agent.Agent], libAgentsNum)
	searchAgents := make([]utils.Choice[*agent.Agent], searchAgentsNum)
	var err error

	for i := 0; i < libAgentsNum; i++ {
		libAgents[i].Val, err = agent.NewAgent(agent.WithBaseURL(baseURL), agent.WithDefaultTransport())
		if err != nil {
			return nil, failure.NewError(model.ErrCritical, err)
		}
		libAgents[i].Val.Name = fmt.Sprintf("Isulibrary-LibAgent-%d", i+1)
		libAgents[i].Weight = 1
	}
	for i := 0; i < searchAgentsNum; i++ {
		searchAgents[i].Val, err = agent.NewAgent(agent.WithBaseURL(baseURL), agent.WithDefaultTransport())
		if err != nil {
			return nil, failure.NewError(model.ErrCritical, err)
		}
		searchAgents[i].Val.Name = fmt.Sprintf("Isulibrary-SearchAgent-%d", i+1)
		searchAgents[i].Weight = 1
	}

	return &ActionController{
		libAgents:    libAgents,
		searchAgents: searchAgents,
	}, nil
}

func (c *ActionController) libAgent() *agent.Agent {
	return utils.WeightedSelect(c.libAgents)
}

func (c *ActionController) searchAgent() *agent.Agent {
	return utils.WeightedSelect(c.searchAgents)
}
