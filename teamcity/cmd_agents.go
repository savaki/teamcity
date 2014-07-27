package main

import (
	"github.com/codegangsta/cli"
	"github.com/savaki/teamcity/v80"
)

var agentCommand = cli.Command{
	Name: "agent",
	Subcommands: []cli.Command{
		{
			Name: "list",
			Flags: []cli.Flag{
				FlagVerbose,
			},
			Action: agentListAction,
		},
		{
			Name: "find",
			Flags: []cli.Flag{
				FlagAgentName,
				FlagAgentId,
				FlagVerbose,
			},
			Action: agentFindAction,
		},
	},
}

func agentListAction(c *cli.Context) {
	client := Get80Client(c)
	Print(client.Agents())
}

func agentFindAction(c *cli.Context) {
	client := Get80Client(c)

	filters := []v80.AgentFilter{}

	// filter by id
	if values := c.StringSlice(FLAG_AGENT_ID); values != nil {
		for _, nameFilter := range values {
			filter := v80.NewFilter(nameFilter, v80.AgentIdAccessor)
			filters = append(filters, filter)
		}
	}

	// filter by name
	if values := c.StringSlice(FLAG_AGENT_NAME); values != nil {
		for _, nameFilter := range values {
			filter := v80.NewFilter(nameFilter, v80.AgentNameAccessor)
			filters = append(filters, filter)
		}
	}

	Print(client.FindAgents(filters))
}
