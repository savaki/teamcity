package main

import "github.com/codegangsta/cli"

var projectCommand = cli.Command{
	Name:  "project",
	Usage: "project related commands",
	Subcommands: []cli.Command{
		{
			Name:  "list",
			Usage: "list the projects on this server",
			Action: func(c *cli.Context) {
				client := Get80Client(c)
				projects, err := client.Projects()
				Print(projects, err)
			},
		},
	},
}
