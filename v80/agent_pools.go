package v80

import (
	"errors"
	"log"
	"net/url"
	"regexp"
	"strings"
)

type AgentPools struct {
	Href  string       `xml:"href,attr,omitempty" json:"href,omitempty"`
	Pools []*AgentPool `xml:"agentPool,omitempty" json:"agent-pools,omitempty"`
}

type AgentPool struct {
	Id   string `xml:"id,attr,omitempty" json:"id,attr,omitempty"`
	Name string `xml:"name,attr,omitempty" json:"name,attr,omitempty"`
	Href string `xml:"href,attr,omitempty" json:"href,attr,omitempty"`

	Projects *Projects `xml:"projects,omitempty" json:"projects,omitempty"`
	Agents   *Agents   `xml:"agents,omitempty" json:"agents,omitempty"`
}

func (tc *TeamCity) AgentPools() ([]*AgentPool, error) {
	server, err := tc.Server()
	if err != nil {
		return nil, err
	}

	pools := &AgentPools{}
	err = tc.get(server.AgentPools.Href, url.Values{}, pools)
	if err != nil {
		return nil, err
	}

	return pools.Pools, nil
}

type AgentPoolFilter func(*AgentPool) bool

func NewAgentPoolFilter(name string) AgentPoolFilter {
	matcher, err := regexp.Compile(name)
	if err != nil {
		log.Fatalf("unable to filter by name, invalid regexp for name, %s\n", name)
	}

	return func(agent *AgentPool) bool {
		return matcher.Match([]byte(agent.Name))
	}
}

func (tc *TeamCity) FindAgentPools(poolFilter AgentPoolFilter) ([]*AgentPool, error) {
	pools, err := tc.AgentPools()
	if err != nil {
		return nil, err
	}

	filtered := []*AgentPool{}
	for _, pool := range pools {
		if poolFilter(pool) {
			filtered = append(filtered, pool)
		}
	}

	return filtered, nil
}

// assignAgentsToPool - sadly this involves screen scraping
func (tc *TeamCity) assignAgentsToPool(agents []*Agent, pool *AgentPool) error {
	typeId := "pool-agent-types-popup"

	// Request #1 - sets the pool context
	params := url.Values{}
	params.Add("typeId", typeId)
	params.Add("contextId", pool.Id)

	err := tc.get("/popupDialog.html", params, nil)
	if err != nil {
		return err
	}

	// Request #2 - apply
	params = url.Values{}
	params.Add("typeId", typeId)
	params.Add("action", "apply")
	for _, agent := range agents {
		params.Add("checkedItemId", agent.TypeId)
	}
	_, err = tc.httpFn("POST", "/popupDialog.html", url.Values{}, strings.NewReader(params.Encode()), "application/x-www-form-urlencoded; charset=UTF-8")
	return err
}

// AssignAgentsToPool returns the number of agents that were assigned to a pool
func (tc *TeamCity) AssignAgentsToPool(agentFilters AgentFilters, poolFilter AgentPoolFilter) (int, error) {
	pools, err := tc.FindAgentPools(poolFilter)
	if err != nil {
		return 0, err

	} else if len(pools) == 0 {
		return 0, errors.New("pool filter didn't match any pool")

	} else if len(pools) != 1 {
		return 0, errors.New("pool filter matched more than one pool")
	}

	allAgents, err := tc.FindAgents(agentFilters)
	if err != nil {
		return 0, err
	}

	agents := []*Agent{}
	pool := pools[0]
	for _, agent := range allAgents {
		if !poolFilter(agent.Pool) {
			if Verbose {
				log.Printf("assigning agent, %s (id:%s; typeId:%s; ip:%s), to pool, %s\n", agent.Name, agent.Id, agent.TypeId, agent.Ip, pool.Name)
			}
			agents = append(agents, agent)
		}
	}

	if !DryRun {
		err = tc.assignAgentsToPool(agents, pool)
		if err != nil {
			return len(agents), err
		}

	}

	return len(agents), err
}
