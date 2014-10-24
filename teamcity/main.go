package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/savaki/teamcity"
	"github.com/savaki/teamcity/v80"
	"log"
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
	opts := options(c)
	v80.Verbose = opts.Verbose

	codebase := os.Getenv(teamcity.TEAMCITY_URL)
	if value := c.String("url"); value != "" {
		codebase = value
	}
	if codebase == "" {
		log.Fatalln("ERROR: TEAMCITY_URL not set")
	}

	username := os.Getenv(teamcity.TEAMCITY_USERNAME)
	if value := c.String("username"); value != "" {
		username = value
	}
	if username == "" {
		log.Fatalln("ERROR: TEAMCITY_USERNAME not set")
	}

	password := os.Getenv(teamcity.TEAMCITY_PASSWORD)
	if value := c.String("password"); value != "" {
		password = value
	}
	if username == "" {
		log.Fatalln("ERROR: TEAMCITY_PASSWORD not set")
	}

	auth := teamcity.New(username, password)
	return v80.New(auth, codebase)
}

func main() {
	app := cli.NewApp()
	app.Name = "teamcity"
	app.Usage = "a command line interface for TeamCity"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		agentCommand,
		agentPoolsCommand,
		buildCommand,
		projectCommand,
		serverCommand,
	}
	app.Run(os.Args)
}
