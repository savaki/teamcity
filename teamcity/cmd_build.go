package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/savaki/teamcity/v80"
	"io"
	"log"
	"net/url"
	"os"
)

var buildCommand = cli.Command{
	Name:  "build",
	Usage: "build elated commands",
	Subcommands: []cli.Command{
		{
			Name:  "list-build-types",
			Usage: "list the build types on this server",
			Flags: []cli.Flag{
				FlagUrl,
				FlagUsername,
				FlagPassword,
				FlagVerbose,
				FlagTrace,
			},
			Action: func(c *cli.Context) {
				client := Get80Client(c)
				buildTypes, err := client.BuildTypes()
				Print(buildTypes, err)
			},
		},
		{
			Name:  "history",
			Usage: "list the builds that have been executed for a given project",
			Flags: []cli.Flag{
				FlagUrl,
				FlagUsername,
				FlagPassword,
				FlagVerbose,
				FlagTrace,
				FlagBuildTypeId,
				FlagLast,
			},
			Action: func(c *cli.Context) {
				client := Get80Client(c)
				buildTypeId := c.String(FLAG_BUILD_TYPE_ID)
				last := c.Int(FLAG_LAST)
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
			Flags: []cli.Flag{
				FlagUrl,
				FlagUsername,
				FlagPassword,
				FlagVerbose,
				FlagTrace,
				FlagBuildTypeId,
			},
			Action: func(c *cli.Context) {
				client := Get80Client(c)
				buildTypeId := c.String(FLAG_BUILD_TYPE_ID)
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
			Flags: []cli.Flag{
				FlagUrl,
				FlagUsername,
				FlagPassword,
				FlagVerbose,
				FlagTrace,
				FlagBuildId,
			},
			Action: func(c *cli.Context) {
				client := Get80Client(c)
				buildId := c.String(FLAG_BUILD_ID)
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
			Flags: []cli.Flag{
				FlagUrl,
				FlagUsername,
				FlagPassword,
				FlagVerbose,
				FlagTrace,
				FlagBuildId,
			},
			Action: func(c *cli.Context) {
				client := Get80Client(c)
				buildId := c.String(FLAG_BUILD_ID)
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
			Flags: []cli.Flag{
				FlagUrl,
				FlagUsername,
				FlagPassword,
				FlagVerbose,
				FlagTrace,
				FlagBuildId,
				FlagArtifactName,
			},
			Action: buildDownloadArtifact,
		},
	},
}

func buildDownloadArtifact(c *cli.Context) {
	client := Get80Client(c)
	buildId := c.String(FLAG_BUILD_ID)
	artifactName := c.String(FLAG_BUILD_ARTIFACT_NAME)
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
}
