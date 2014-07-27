package main

import (
	"github.com/codegangsta/cli"
	"github.com/savaki/teamcity/v80"
)

var agentPoolsCommand = cli.Command{
	Name: "agent-pool",
	Subcommands: []cli.Command{
		{
			Name: "list",
			Flags: []cli.Flag{
				FlagVerbose,
			},
			Action: agentPoolListAction,
		},
	},
}

func agentPoolListAction(c *cli.Context) {
	client := Get80Client(c)

	opts := options(c)
	v80.Verbose = opts.Verbose

	Print(client.AgentPools())
}
