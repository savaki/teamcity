package v80

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"github.com/savaki/teamcity"
	"io/ioutil"
	"net/http"
	"net/url"
)

type TeamCity struct {
	auth *teamcity.Auth
	get  func(string, url.Values, interface{}) error
}

func getFn(client *http.Client, auth *teamcity.Auth, codebase string) func(path string, params url.Values, result interface{}) error {
	return func(path string, params url.Values, result interface{}) error {
		theUrl := fmt.Sprintf("%s/httpAuth%s", codebase, path)
		queryParams := params.Encode()
		if len(queryParams) > 0 {
			theUrl = theUrl+"?"+queryParams
		}

		request, err := http.NewRequest("GET", theUrl, nil)
		if err != nil {
			return err
		}
		request.SetBasicAuth(auth.Username, auth.Password)

		response, err := client.Do(request)
		if err != nil {
			return err
		}

		defer response.Body.Close()
		switch value := result.(type) {
		case *string:
			data, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return nil
			}
			*value = string(data)
			return nil

		default:
			return xml.NewDecoder(response.Body).Decode(result)
		}
	}
}

func New(auth *teamcity.Auth, codebase string) *TeamCity {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
	}

	return &TeamCity{
		auth: auth,
		get:  getFn(client, auth, codebase),
	}
}

type ServerInfo struct {
	Version      string `xml:"version,attr" json:"version"`
	VersionMajor string `xml:"versionMajor,attr" json:"versionMajor"`
	VersionMinor string `xml:"versionMinor,attr" json:"versionMinor"`
	StartTime    string `xml:"startTime,attr" json:"startTime"`
	CurrentTime  string `xml:"currentTime,attr" json:"currentTime"`
	BuildNumber  string `xml:"buildNumber,attr" json:"buildNumber"`
	BuildDate    string `xml:"buildDate,attr" json:"buildDate"`
	InternalId   string `xml:"internalId,attr" json:"internalId"`
}

func (tc *TeamCity) ServerInfo() (*ServerInfo, error) {
	result := &ServerInfo{}
	err := tc.get("/app/rest/server", url.Values{}, result)
	return result, err
}
