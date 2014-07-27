package v80

import "net/url"

type Project struct {
	Id              string `xml:"id,attr" json:"id"`
	ParentProjectId string `xml:"parentProjectId,attr" json:"parentProjectId"`
	Name            string `xml:"name,attr" json:"name"`
	Description     string `xml:"description,attr" json:"description"`
	Href            string `xml:"href,attr" json:"href"`
	WebUrl          string `xml:"webUrl,attr" json:"webUrl"`
}

type Projects struct {
	Href     string     `xml:"href,attr,omitempty" json:"href,attr,omitempty"`
	Projects []*Project `xml:"project,omitempty" json:"projects,omitempty"`
}

func (tc *TeamCity) Projects() (*Projects, error) {
	projects := &Projects{}
	err := tc.get("/app/rest/projects", url.Values{}, projects)
	return projects, err
}
