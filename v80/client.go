package v80

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/savaki/teamcity"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var Trace bool = false
var Verbose bool = false
var DryRun bool = false

type TeamCity struct {
	httpFn func(method, path string, params url.Values, body io.Reader, contentType string) (io.ReadCloser, error)
}

func httpFn(client *http.Client, auth *teamcity.Auth, codebase string) func(method, path string, params url.Values, body io.Reader, contentType string) (io.ReadCloser, error) {
	return func(method, path string, params url.Values, body io.Reader, contentType string) (io.ReadCloser, error) {
		urlStr := fmt.Sprintf("%s%s", codebase, path)
		queryParams := params.Encode()
		if len(queryParams) > 0 {
			urlStr = urlStr + "?" + queryParams
		}

		if Trace {
			log.Printf("%s %s\n", method, urlStr)
		}

		// calculate the Content-Length
		contentLength := 0
		if body != nil {
			// determine the content length
			data, err := ioutil.ReadAll(body)
			if err != nil {
				return nil, err
			}
			switch readCloser := body.(type) {
			case io.ReadCloser:
				defer readCloser.Close()
			}

			body = bytes.NewReader(data)
			contentLength = len(data)
		}

		req, err := http.NewRequest(method, urlStr, body)
		if err != nil {
			return nil, err
		}

		if contentType != "" {
			req.Header.Add("Content-Type", contentType)
		}

		if contentLength != 0 {
			req.ContentLength = int64(contentLength)
		}

		req.SetBasicAuth(auth.Username, auth.Password)
		response, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		if Trace {
			log.Printf("status code => %d\n", response.StatusCode)
		}

		if response.StatusCode < 200 || response.StatusCode >= 300 {
			if Verbose {
				defer response.Body.Close()

				data, err := ioutil.ReadAll(response.Body)
				if err != nil {
					fmt.Println(string(data))
				}
			}
			return nil, errors.New(fmt.Sprintf("%s %s returned status code, %d", method, urlStr, response.StatusCode))
		}

		return response.Body, nil
	}
}

func (t *TeamCity) put(path string, params url.Values, body io.ReadCloser, contentType string) (io.ReadCloser, error) {
	return t.httpFn("PUT", path, params, body, contentType)
}

func (t *TeamCity) get(path string, params url.Values, target interface{}) error {
	body, err := t.httpFn("GET", path, params, nil, "application/xml")
	if err != nil {
		return err
	}

	if target == nil {
		return nil
	}

	switch value := target.(type) {
	case *string:
		data, err := ioutil.ReadAll(body)
		if err != nil {
			return nil
		}
		*value = string(data)
		return nil

	default:
		return xml.NewDecoder(body).Decode(target)
	}
}

func (t *TeamCity) Download(path string, params url.Values) (io.ReadCloser, error) {
	return t.httpFn("GET", path, params, nil, "application/octet-stream")
}

func New(auth *teamcity.Auth, codebase string) *TeamCity {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Transport: transport,
		Jar:       cookieJar,
	}

	return &TeamCity{
		httpFn: httpFn(client, auth, codebase),
	}
}
