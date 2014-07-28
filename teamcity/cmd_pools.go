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
				FlagAgentPoolName,
				FlagVerbose,
				FlagTrace,
			},
			Action: agentPoolListAction,
		},
	},
}

func agentPoolListAction(c *cli.Context) {
	client := Get80Client(c)

	opts := options(c)
	v80.Verbose = opts.Verbose

	agentPoolName := c.String(FLAG_AGENT_POOL_NAME)
	poolFilter := v80.NewAgentPoolFilter(agentPoolName)

	Print(client.FindAgentPools(poolFilter))
}
