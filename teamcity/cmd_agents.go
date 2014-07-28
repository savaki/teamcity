package main

import (
	"github.com/codegangsta/cli"
	"github.com/savaki/teamcity/v80"
	"log"
)

var agentCommand = cli.Command{
	Name: "agent",
	Subcommands: []cli.Command{
		{
			Name: "list",
			Flags: []cli.Flag{
				FlagVerbose,
				FlagTrace,
			},
			Action: agentListAction,
		},
		{
			Name: "find",
			Flags: []cli.Flag{
				FlagAgentName,
				FlagAgentId,
				FlagVerbose,
				FlagTrace,
			},
			Action: agentFindAction,
		},
		{
			Name: "authorize",
			Flags: []cli.Flag{
				FlagAgentName,
				FlagAgentId,
				FlagVerbose,
				FlagTrace,
			},
			Action: agentAuthorizeAction,
		},
		{
			Name: "assign-to-pool",
			Flags: []cli.Flag{
				FlagAgentId,
				FlagAgentName,
				FlagAgentPoolName,
				FlagVerbose,
				FlagTrace,
			},
			Action: agentAssignToPoolAction,
		},
	},
}

func agentFilters(c *cli.Context) v80.AgentFilters {
	filters := []v80.AgentFilter{}

	// filter by id
	if values := c.StringSlice(FLAG_AGENT_ID); values != nil {
		for _, nameFilter := range values {
			filter := v80.NewAgentFilter(nameFilter, v80.AgentIdAccessor)
			filters = append(filters, filter)
		}
	}

	// filter by name
	if values := c.StringSlice(FLAG_AGENT_NAME); values != nil {
		for _, nameFilter := range values {
			filter := v80.NewAgentFilter(nameFilter, v80.AgentNameAccessor)
			filters = append(filters, filter)
		}
	}

	return filters
}

func agentListAction(c *cli.Context) {
	client := Get80Client(c)
	Print(client.Agents())
}

func agentFindAction(c *cli.Context) {
	client := Get80Client(c)
	filters := agentFilters(c)

	Print(client.FindAgents(filters))
}

func agentAuthorizeAction(c *cli.Context) {
	client := Get80Client(c)
	filters := agentFilters(c)

	client.AuthorizeAgents(filters)
}

func agentAssignToPoolAction(c *cli.Context) {
	client := Get80Client(c)
	filters := agentFilters(c)
	opts := options(c)
	v80.Trace = opts.Trace
	v80.Verbose = opts.Verbose

	agentPoolName := c.String(FLAG_AGENT_POOL_NAME)
	if agentPoolName == "" {
		log.Fatalf("no agent pool name specified via flag --%s\n", FLAG_AGENT_POOL_NAME)
	}
	poolFilter := v80.NewAgentPoolFilter(agentPoolName)

	agentsAssigned, err := client.AssignAgentsToPool(filters, poolFilter)
	if err != nil {
		log.Fatalln(err)
	}

	if opts.Verbose {
		log.Printf("%d agent(s) assigned to pool\n", agentsAssigned)
	}
}
