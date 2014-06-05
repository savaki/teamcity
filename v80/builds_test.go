package v80

import (
	"encoding/xml"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestBuildTypes(t *testing.T) {
	var buildTypes *BuildTypes = nil
	var content string = ""
	var err error = nil

	Convey("Given a buildTypes xml response", t, func() {
		content = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><buildTypes count="1"><buildType id="Pencil_RailsBackend_10CompileAndTest" name="1.0 Run Tests" projectName="Pencil :: Rails Backend" projectId="Pencil_RailsBackend" href="/httpAuth/app/rest/buildTypes/id:Pencil_RailsBackend_10CompileAndTest" webUrl="https://ec2-54-183-24-208.us-west-1.compute.amazonaws.com/viewType.html?buildTypeId=Pencil_RailsBackend_10CompileAndTest"/></buildTypes>`

		Convey("When I parse this via XML", func() {
			buildTypes = &BuildTypes{}
			err = xml.NewDecoder(strings.NewReader(content)).Decode(buildTypes)

			Convey("Then I expect no parse errors", func() {
				So(err, ShouldBeNil)
			})

			Convey("And I expect the parameters to be set", func() {
				So(len(buildTypes.BuildTypes), ShouldEqual, 1)
			})
		})
	})

}
