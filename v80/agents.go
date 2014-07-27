package v80

import (
	"log"
	"net/url"
	"regexp"
)

type Agents struct {
	Href   string   `xml:"href,attr,omitempty" json:"href,omitempty"`
	Count  int      `xml:"count,attr,omitempty" json:"count,omitempty"`
	Agents []*Agent `xml:"agent,omitempty" json:"agents,omitempty"`
}

type Agent struct {
	Id     string `xml:"id,attr,omitempty" json:"id,attr,omitempty"`
	Name   string `xml:"name,attr,omitempty" json:"name,attr,omitempty"`
	TypeId string `xml:"typeId,attr,omitempty" json:"typeId,attr,omitempty"`
	Href   string `xml:"href,attr,omitempty" json:"href,attr,omitempty"`

	Connected  bool   `xml:"connected,attr,omitempty" json:"connected,attr,omitempty"`
	Enabled    bool   `xml:"enabled,attr,omitempty" json:"enabled,attr,omitempty"`
	Authorized bool   `xml:"authorized,attr,omitempty" json:"authorized,attr,omitempty"`
	UpToDate   bool   `xml:"updatodate,attr,omitempty" json:"updatodate,attr,omitempty"`
	Ip         string `xml:"ip,attr,omitempty" json:"ip,attr,omitempty"`

	Properties []*Property `xml:"properties>property,omitempty" json:"properties,omitempty"`
	Pool       *AgentPool  `xml:"pool,omitempty" json:"pool,omitempty"`
}

type Property struct {
	Name  string `xml:"name,attr,omitempty" json:"name,attr,omitempty"`
	Value string `xml:"value,attr,omitempty" json:"value,attr,omitempty"`
}

type AgentPool struct {
	Id   string `xml:"id,attr,omitempty" json:"id,attr,omitempty"`
	Name string `xml:"name,attr,omitempty" json:"name,attr,omitempty"`
	Href string `xml:"href,attr,omitempty" json:"href,attr,omitempty"`

	Projects *Projects `xml:"projects,omitempty" json:"projects,omitempty"`
	Agents   *Agents   `xml:"agents,omitempty" json:"agents,omitempty"`
}

type AgentFilters []AgentFilter

type AgentFilter func(*Agent) bool

type AgentAccessor func(*Agent) string

func AgentIdAccessor(agent *Agent) string {
	return agent.Id
}

func AgentNameAccessor(agent *Agent) string {
	return agent.Name
}

func NewFilter(name string, accessor AgentAccessor) AgentFilter {
	matcher, err := regexp.Compile(name)
	if err != nil {
		log.Fatalf("unable to filter by name, invalid regexp for name, %s\n", name)
	}

	return func(agent *Agent) bool {
		value := accessor(agent)
		return matcher.Match([]byte(value))
	}
}

func (tc *TeamCity) UpdateAgent(filters AgentFilters, field string, value string) error {
	agents, err := tc.FindAgents(filters)
	if err != nil {
		return err
	}

	for _, agent := range agents {
		println(agent.Href)
	}

	return nil
}

func (tc *TeamCity) Agents() (*Agents, error) {
	server, err := tc.Server()
	if err != nil {
		return nil, err
	}

	agents := &Agents{}
	err = tc.get(server.Agents.Href, url.Values{}, agents)
	return agents, err
}

func (tc *TeamCity) FindAgents(filters AgentFilters) ([]*Agent, error) {
	agents, err := tc.Agents()
	if err != nil {
		return nil, err
	}

	filteredAgents := []*Agent{}

	for _, agent := range agents.Agents {
		includeAgent := true
		for _, filter := range filters {
			if !filter(agent) {
				includeAgent = false
			}
		}

		if includeAgent {
			filteredAgents = append(filteredAgents, agent)
		}
	}

	return filteredAgents, nil
}
