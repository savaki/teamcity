package v80

import "net/url"

type Agents struct {
	Href   string  `xml:"href,attr,omitempty" json:"href,omitempty"`
	Count  int     `xml:"count,attr,omitempty" json:"count,omitempty"`
	Agents []Agent `xml:"agent,omitempty" json:"agents,omitempty"`
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

	Properties []Property `xml:"property,omitempty" json:"property,omitempty"`
	Pool       *AgentPool `xml:"pool,omitempty" json:"pool,omitempty"`
}

type Property struct {
	Name  string `xml:"name,attr,omitempty" json:"name,attr,omitempty"`
	Value string `xml:"value,attr,omitempty" json:"value,attr,omitempty"`
}

type AgentPool struct {
	Id   string `xml:"id,attr,omitempty" json:"id,attr,omitempty"`
	Name string `xml:"name,attr,omitempty" json:"name,attr,omitempty"`
	Href string `xml:"href,attr,omitempty" json:"href,attr,omitempty"`

	Projects Projects `xml:"projects,omitempty" json:"projects,omitempty"`
	Agents   Agents   `xml:"agents,omitempty" json:"agents,omitempty"`
}

func (tc *TeamCity) Agents() (*Agents, error) {
	agents := &Agents{}
	err := tc.get("/app/rest/agents", url.Values{}, agents)
	return agents, err
}
