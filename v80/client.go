package v80

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"github.com/savaki/teamcity"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type TeamCity struct {
	auth     *teamcity.Auth
	get      func(string, url.Values, interface{}) error
	Download func(string, url.Values) (io.ReadCloser, error)
}

func getFn(client *http.Client, auth *teamcity.Auth, codebase string) func(path string, params url.Values) (io.ReadCloser, error) {
	return func(path string, params url.Values) (io.ReadCloser, error) {
		theUrl := fmt.Sprintf("%s%s", codebase, path)
		queryParams := params.Encode()
		if len(queryParams) > 0 {
			theUrl = theUrl + "?" + queryParams
		}

		request, err := http.NewRequest("GET", theUrl, nil)
		if err != nil {
			return nil, err
		}
		request.SetBasicAuth(auth.Username, auth.Password)

		response, err := client.Do(request)
		if err != nil {
			return nil, err
		}

		return response.Body, nil
	}
}

func getXmlFn(client *http.Client, auth *teamcity.Auth, codebase string) func(path string, params url.Values, result interface{}) error {
	get := getFn(client, auth, codebase)

	return func(path string, params url.Values, result interface{}) error {
		body, err := get(path, params)
		if err != nil {
			return err
		}

		defer body.Close()
		switch value := result.(type) {
		case *string:
			data, err := ioutil.ReadAll(body)
			if err != nil {
				return nil
			}
			*value = string(data)
			return nil

		default:
			return xml.NewDecoder(body).Decode(result)
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
		auth:     auth,
		get:      getXmlFn(client, auth, codebase),
		Download: getFn(client, auth, codebase),
	}
}
