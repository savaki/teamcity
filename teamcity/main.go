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

func Get80Client(c*cli.Context) *v80.TeamCity {
	codebase := c.String("url")
	if codebase == "" {
		log.Fatalln("ERROR: no TeamCity url specified with --url")
	}

	username := os.Getenv("TEAMCITY_USERNAME")
	if value := c.String("username"); value != "" {
		username = value
	}
	if username == "" {
		log.Fatalln("ERROR: TeamCity username not set")
	}

	password := os.Getenv("TEAMCITY_PASSWORD")
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
			Flags:globalFlags,
			Action: func(c *cli.Context) {
				client := Get80Client(c)
				serverInfo, err := client.ServerInfo()
				Print(serverInfo, err)
			},
		},
		{
			Name:  "project",
			Usage: "project related commands",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the projects on this server",
					Flags:globalFlags,
					Action: func(c *cli.Context) {
						client := Get80Client(c)
						projects, err := client.Projects()
						Print(projects, err)
					},
				},
			},
		},
	}
	app.Run(os.Args)
}
