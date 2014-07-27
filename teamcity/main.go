package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/savaki/teamcity"
	"github.com/savaki/teamcity/v80"
	"io"
	"log"
	"net/url"
	"os"
)

func Print(value interface{}, err error) {
	if err != nil {
		log.Fatalln(err)
	}

	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(data))
}

func Get80Client(c *cli.Context) *v80.TeamCity {
	codebase := os.Getenv(teamcity.TEAMCITY_URL)
	if value := c.String("url"); value != "" {
		codebase = value
	}
	if codebase == "" {
		log.Fatalln("ERROR: TeamCity url not set")
	}

	username := os.Getenv(teamcity.TEAMCITY_USERNAME)
	if value := c.String("username"); value != "" {
		username = value
	}
	if username == "" {
		log.Fatalln("ERROR: TeamCity username not set")
	}

	password := os.Getenv(teamcity.TEAMCITY_PASSWORD)
	if value := c.String("password"); value != "" {
		password = value
	}
	if username == "" {
		log.Fatalln("ERROR: TeamCity password not set")
	}

	auth := teamcity.New(username, password)
	return v80.New(auth, codebase)
}

func main() {
	globalFlags := []cli.Flag{
		cli.StringFlag{"url", "", "url of the TeamCity server"},
		cli.StringFlag{"username, u", "", "TeamCity username"},
		cli.StringFlag{"password, p", "", "TeamCity password"},
	}

	app := cli.NewApp()
	app.Name = "teamcity"
	app.Usage = "a command line interface for TeamCity"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "retrieve info on the server",
			Flags: globalFlags,
			Action: func(c *cli.Context) {
				client := Get80Client(c)
				server, err := client.Server()
				Print(server, err)
			},
		},
		projectCommand,
		agentCommand,
		{
			Name:  "build",
			Usage: "build elated commands",
			Subcommands: []cli.Command{
				{
					Name:  "list-build-types",
					Usage: "list the build types on this server",
					Flags: globalFlags,
					Action: func(c *cli.Context) {
						client := Get80Client(c)
						buildTypes, err := client.BuildTypes()
						Print(buildTypes, err)
					},
				},
				{
					Name:  "history",
					Usage: "list the builds that have been executed for a given project",
					Flags: append(globalFlags, []cli.Flag{
						cli.StringFlag{"build-type-id", "", "the build type id"},
						cli.IntFlag{"last", 1, "how many builds to retrieve"},
					}...),
					Action: func(c *cli.Context) {
						client := Get80Client(c)
						buildTypeId := c.String("build-type-id")
						last := c.Int("last")
						if buildTypeId == "" {
							log.Fatalln("ERROR: required parameter, build-type-id, not specified")
						}
						locator := v80.BuildTypeLocator{Id: buildTypeId}
						builds, err := client.Builds(locator)

						filtered := []*v80.Build{}
						for i := 0; builds.Builds != nil && i < len(builds.Builds) && i < last; i++ {
							filtered = append(filtered, builds.Builds[i])
						}

						Print(v80.Builds{Builds: filtered}, err)
					},
				},
				{
					Name:  "status",
					Usage: "the status of the last build of this type",
					Flags: append(globalFlags, []cli.Flag{
						cli.StringFlag{"build-type-id", "", "the build type id"},
					}...),
					Action: func(c *cli.Context) {
						client := Get80Client(c)
						buildTypeId := c.String("build-type-id")
						if buildTypeId == "" {
							log.Fatalln("ERROR: required parameter, build-type-id, not specified")
						}
						locator := v80.BuildTypeLocator{Id: buildTypeId}
						builds, err := client.Builds(locator)
						if err != nil {
							log.Fatalln(err)
						}

						if builds.Builds == nil || len(builds.Builds) == 0 {
							log.Fatalf("ERROR: no builds yet for this project\n")
						}

						build := builds.Builds[0]
						fmt.Println(build.Status)
						if build.Status == "SUCCESS" {
							os.Exit(0)
						} else {
							os.Exit(1)
						}
					},
				},
				{
					Name:  "details",
					Usage: "list the builds that have been executed for a given project",
					Flags: append(globalFlags, []cli.Flag{
						cli.StringFlag{"build-id", "", "the build to retrieve details for"},
					}...),
					Action: func(c *cli.Context) {
						client := Get80Client(c)
						buildId := c.String("build-id")
						if buildId == "" {
							log.Fatalln("ERROR: required parameter, build-id, not specified")
						}
						build, err := client.BuildDetail(buildId)
						Print(build, err)
					},
				},
				{
					Name:  "list-artifacts",
					Usage: "list the artifacts for this build",
					Flags: append(globalFlags, []cli.Flag{
						cli.StringFlag{"build-id", "", "the build to retrieve details for"},
					}...),
					Action: func(c *cli.Context) {
						client := Get80Client(c)
						buildId := c.String("build-id")
						if buildId == "" {
							log.Fatalln("ERROR: required parameter, build-id, not specified")
						}
						artifacts, err := client.BuildArtifacts(buildId)
						Print(artifacts, err)
					},
				},
				{
					Name:  "download-artifact",
					Usage: "download the artifact with the specified name",
					Flags: append(globalFlags, []cli.Flag{
						cli.StringFlag{"build-id", "", "the build to retrieve details for"},
						cli.StringFlag{"artifact-name", "", "the filename of the artifact to download"},
					}...),
					Action: func(c *cli.Context) {
						client := Get80Client(c)
						buildId := c.String("build-id")
						artifactName := c.String("artifact-name")
						if buildId == "" {
							log.Fatalln("ERROR: required parameter, build-id, not specified")
						}
						if artifactName == "" {
							log.Fatalln("ERROR: required parameter, artifact-name, not specified")
						}
						artifacts, err := client.BuildArtifacts(buildId)
						if err != nil {
							log.Fatalln(err)
						}
						for _, artifact := range artifacts.Artifacts {
							if artifact.Name == artifactName && artifact.Content != nil {
								fmt.Printf("Downloading artfact, %s\n", artifact.Name)
								content, err := client.Download(artifact.Content.Href, url.Values{})
								defer content.Close()
								if err != nil {
									log.Fatalln(err)
								}

								fi, err := os.Create(artifact.Name)
								if err != nil {
									log.Fatalln(err)
								}
								defer fi.Close()

								io.Copy(fi, content)
								return
							}
						}

						log.Fatalf("ERROR: unable to find artifact, %s, in build id, %s\n", artifactName, buildId)
					},
				},
			},
		},
	}
	app.Run(os.Args)
}
