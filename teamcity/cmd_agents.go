package main

import "github.com/codegangsta/cli"

var agentCommand = cli.Command{
	Name: "agent",
	Subcommands: []cli.Command{
		{
			Name:   "list",
			Action: agentsListAction,
		},
	},
}

func agentsListAction(c *cli.Context) {
	client := Get80Client(c)
	Print(client.Agents())
}
