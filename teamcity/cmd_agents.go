package main

import (
	"github.com/codegangsta/cli"
	"github.com/savaki/teamcity/v80"
)

const (
	FLAG_AGENT_NAME = "agent-name"
	FLAG_AGENT_ID   = "agent-id"
)

var agentCommand = cli.Command{
	Name: "agent",
	Subcommands: []cli.Command{
		{
			Name:   "list",
			Action: agentListAction,
		},
		{
			Name: "find",
			Flags: []cli.Flag{
				cli.StringSliceFlag{FLAG_AGENT_NAME, &cli.StringSlice{}, "filter by agent name (regexp)"},
				cli.StringSliceFlag{FLAG_AGENT_ID, &cli.StringSlice{}, "filter agent name (regexp)"},
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
