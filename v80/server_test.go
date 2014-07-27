package v80

import (
	"encoding/xml"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	Convey("Given an xml list of agents", t, func() {
		content := `
<?xml version="1.0" encoding="UTF-8" standalone="yes"?><server version="8.1.4 (build 30168)" versionMajor="8" versionMinor="1" startTime="20140726T110855-0700" currentTime="20140726T213033-0700" buildNumber="30168" buildDate="20140722T000000-0700" internalId="0dd78d9a-7e0d-4a30-b119-efb9c1e2cdd3"><projects href="/app/rest/projects"/><vcsRoots href="/app/rest/vcs-roots"/><builds href="/app/rest/builds"/><users href="/app/rest/users"/><userGroups href="/app/rest/userGroups"/><agents href="/app/rest/agents"/><buildQueue href="/app/rest/buildQueue"/><agentPools href="/app/rest/agentPools"/></server>`

		Convey("When I decode them", func() {
			server := Server{}
			err := xml.NewDecoder(strings.NewReader(content)).Decode(&server)
			So(err, ShouldBeNil)

			Convey("And I server.Projects.Href to be set", func() {
				So(server.Projects.Href, ShouldEqual, "/app/rest/projects")
			})

			Convey("And I server.Agents.Href to be set", func() {
				So(server.Agents.Href, ShouldEqual, "/app/rest/agents")
			})
		})
	})
}
