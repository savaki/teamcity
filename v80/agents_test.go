package v80

import (
	"encoding/xml"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestAgents(t *testing.T) {
	Convey("Given an xml list of agents", t, func() {
		content := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><agents count="2" href="/app/rest/agents"><agent id="1" name="agent-1" typeId="1" href="/app/rest/agents/id:1"/><agent id="2" name="agent-2" typeId="2" href="/app/rest/agents/id:2"/></agents>`

		Convey("When I decode them", func() {
			agents := Agents{}
			err := xml.NewDecoder(strings.NewReader(content)).Decode(&agents)
			So(err, ShouldBeNil)

			Convey("Then I expect 2 agents", func() {
				So(len(agents.Agents), ShouldEqual, 2)
			})

			Convey("And I expect the first agent's attributes to be set", func() {
				agent := agents.Agents[0]
				So(agent.Id, ShouldEqual, "1")
				So(agent.Name, ShouldEqual, "agent-1")
				So(agent.TypeId, ShouldEqual, "1")
				So(agent.Href, ShouldEqual, "/app/rest/agents/id:1")
			})

			Convey("And I expect the second agent's attributes to be set", func() {
				agent := agents.Agents[1]
				So(agent.Id, ShouldEqual, "2")
				So(agent.Name, ShouldEqual, "agent-2")
				So(agent.TypeId, ShouldEqual, "2")
				So(agent.Href, ShouldEqual, "/app/rest/agents/id:2")
			})
		})
	})
}
