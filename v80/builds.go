package v80

import (
	"fmt"
	"net/url"
)

type BuildType struct {
	Id          string `xml:"id,attr" json:"id"`
	Name        string `xml:"name,attr" json:"name"`
	Href        string `xml:"href,attr" json:"href"`
	WebUrl      string `xml:"webUrl,attr" json:"webUrl"`
	ProjectId   string `xml:"projectId,attr" json:"projectId"`
	ProjectName string `xml:"projectName,attr" json:"projectName"`
}

type BuildTypes struct {
	Count      int          `xml:"count,attr" json:"count"`
	BuildTypes []*BuildType `xml:"buildType" json:"buildType"`
}

func (tc *TeamCity) BuildTypes() (*BuildTypes, error) {
	buildTypes := &BuildTypes{}
	err := tc.get("/app/rest/buildTypes", url.Values{}, buildTypes)
	return buildTypes, err
}

type BuildTypeLocator struct {
	Id string
}

func (b BuildTypeLocator) String() string {
	return fmt.Sprintf("id:%s", b.Id)
}

type BuildTypeDetail struct {
}

func (tc *TeamCity) BuildTypeDetail(locator BuildTypeLocator) (*BuildTypeDetail, error) {
	detail := &BuildTypeDetail{}
	path := fmt.Sprintf("/app/rest/buildTypes/%s", locator)
	err := tc.get(path, url.Values{}, detail)
	return detail, err
}

type Build struct {
	Id            string `xml:"id,attr" json:"id"`
	BuildTypeId   string `xml:"buildTypeId,attr" json:"buildTypeId"`
	Number        string `xml:"number,attr" json:"number"`
	Status        string `xml:"status,attr" json:"status"`
	State         string `xml:"state,attr" json:"state"`
	BranchName    string `xml:"branchName,attr" json:"branchName"`
	DefaultBranch string `xml:"defaultBranch,attr" json:"defaultBranch"`
	Href          string `xml:"href,attr" json:"href"`
	WebUrl        string `xml:"webUrl,attr" json:"webUrl"`
}

type Builds struct {
	Count  int      `xml:"count,attr" json:"count"`
	Builds []*Build `xml:"build" json:"build"`
}

func (tc *TeamCity) Builds(locator BuildTypeLocator) (*Builds, error) {
	builds := &Builds{}
	path := fmt.Sprintf("/app/rest/buildTypes/%s/builds/", locator)
	err := tc.get(path, url.Values{}, builds)
	return builds, err
}

type Link struct {
	Href string `xml:"href,attr" json:"href"`
}

type BuildDetail struct {
	Id            string     `xml:"id,attr" json:"id"`
	BuildTypeId   string     `xml:"buildTypeId,attr" json:"buildTypeId"`
	Number        string     `xml:"number,attr" json:"number"`
	Status        string     `xml:"status,attr" json:"status"`
	State         string     `xml:"state,attr" json:"state"`
	BranchName    string     `xml:"branchName,attr" json:"branchName"`
	DefaultBranch string     `xml:"defaultBranch,attr" json:"defaultBranch"`
	Href          string     `xml:"href,attr" json:"href"`
	WebUrl        string     `xml:"webUrl,attr" json:"webUrl"`
	StatusText    string     `xml:"statusText" json:"statusText"`
	BuildType     *BuildType `xml:"buildType" json:"buildType"`
	QueuedDate    string     `xml:"queuedDate" json:"queuedDate"`
	StartDate     string     `xml:"startDate" json:"startDate"`
	FinishDate    string     `xml:"finishDate" json:"finishDate"`
	Artifacts     Link       `xml:"artifacts" json:"artifacts"`
}

func (tc *TeamCity) BuildDetail(buildId string) (*BuildDetail, error) {
	detail := &BuildDetail{}
	path := fmt.Sprintf("/app/rest/builds/id:%s", buildId)
	err := tc.get(path, url.Values{}, detail)
	return detail, err
}

type Artifact struct {
	ModificationTime string `xml:"modificationTime,attr" json:"modificationTime"`
	Name             string `xml:"name,attr" json:"name"`
	Href             string `xml:"href,attr" json:"href"`
	Children         *Link  `xml:"children,omitempty" json:"children,omitempty"`
	Content          *Link  `xml:"content,omitempty" json:"content,omitempty"`
}

type Artifacts struct {
	Artifacts []*Artifact `xml:"file" json:"file"`
}

func (tc *TeamCity) BuildArtifacts(buildId string) (*Artifacts, error) {
	artifacts := &Artifacts{}
	path := fmt.Sprintf("/app/rest/builds/id:%s/artifacts/children", buildId)
	err := tc.get(path, url.Values{}, artifacts)
	return artifacts, err
}
