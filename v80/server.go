package v80

import "net/url"

type Server struct {
	Version      string `xml:"version,attr,omitempty" json:"version,attr,omitempty"`
	VersionMajor string `xml:"versionMajor,attr,omitempty" json:"versionMajor,attr,omitempty"`
	VersionMinor string `xml:"versionMinor,attr,omitempty" json:"versionMinor,attr,omitempty"`
	StartTime    string `xml:"startTime,attr,omitempty" json:"startTime,attr,omitempty"`
	CurrentTime  string `xml:"currentTime,attr,omitempty" json:"currentTime,attr,omitempty"`
	BuildNumber  string `xml:"buildNumber,attr,omitempty" json:"buildNumber,attr,omitempty"`
	BuildDate    string `xml:"buildDate,attr,omitempty" json:"buildDate,attr,omitempty"`
	InternalId   string `xml:"internalId,attr,omitempty" json:"internalId,attr,omitempty"`

	Projects   *Projects   `xml:"projects,omitempty" json:"projects,omitempty"`
	Agents     *Agents     `xml:"agents,omitempty" json:"agents,omitempty"`
	AgentPools *AgentPools `xml:"agentPools,omitempty" json:"agent-pools,omitempty"`
}

func (tc *TeamCity) Server() (*Server, error) {
	result := &Server{}
	err := tc.get("/app/rest/server", url.Values{}, result)
	return result, err
}
