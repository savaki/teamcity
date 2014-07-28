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
				FlagAllAgents,
				FlagVerbose,
				FlagTrace,
			},
			Action: agentAuthorizeAction,
		},
		{
			Name: "deauthorize",
			Flags: []cli.Flag{
				FlagAgentName,
				FlagAgentId,
				FlagAllAgents,
				FlagVerbose,
				FlagTrace,
			},
			Action: agentDeauthorizeAction,
		},
		{
			Name: "remove-deauthorized",
			Flags: []cli.Flag{
				FlagAgentName,
				FlagAgentId,
				FlagAllAgents,
				FlagVerbose,
				FlagTrace,
				FlagDryRun,
			},
			Action: agentRemoveDeauthorizedAction,
		},
		{
			Name: "remove",
			Flags: []cli.Flag{
				FlagAgentName,
				FlagAgentId,
				FlagAllAgents,
				FlagVerbose,
				FlagTrace,
			},
			Action: agentDeauthorizeAction,
		},
		{
			Name: "assign-to-pool",
			Flags: []cli.Flag{
				FlagAgentPoolName,
				FlagAgentId,
				FlagAgentName,
				FlagAllAgents,
				FlagVerbose,
				FlagTrace,
			},
			Action: agentAssignToPoolAction,
		},
	},
}

func agentFilters(c *cli.Context) v80.AgentFilters {
	filters := []v80.AgentFilter{}

	if c.Bool(FLAG_ALL_AGENTS) {
		return append(filters, v80.NoopAgentFilter(true))
	}

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
	opts := options(c)
	v80.Verbose = opts.Verbose
	v80.Trace = opts.Trace

	if filters.IsEmpty() {
		log.Fatalf("agents must be specified via either --%s or --%s\n", FLAG_AGENT_ID, FLAG_AGENT_NAME)
	}

	agentsAuthorized, err := client.AuthorizeAgents(filters)
	if err != nil {
		log.Fatalln(err)
	}

	if opts.Verbose {
		log.Printf("%d agent(s) authorized\n", agentsAuthorized)
	}
}

func agentDeauthorizeAction(c *cli.Context) {
	client := Get80Client(c)
	filters := agentFilters(c)
	opts := options(c)
	v80.Verbose = opts.Verbose
	v80.Trace = opts.Trace

	if filters.IsEmpty() {
		log.Fatalf("agents must be specified via either --%s or --%s\n", FLAG_AGENT_ID, FLAG_AGENT_NAME)
	}

	agentsAuthorized, err := client.DeauthorizeAgents(filters)
	if err != nil {
		log.Fatalln(err)
	}

	if opts.Verbose {
		log.Printf("%d agent(s) deauthorized\n", agentsAuthorized)
	}
}

func agentAssignToPoolAction(c *cli.Context) {
	client := Get80Client(c)
	filters := agentFilters(c)

	opts := options(c)
	v80.Trace = opts.Trace
	v80.Verbose = opts.Verbose

	if filters.IsEmpty() {
		log.Fatalf("agents must be specified via either --%s or --%s\n", FLAG_AGENT_ID, FLAG_AGENT_NAME)
	}

	// must specify a target pool to put the agents to
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

func agentRemoveDeauthorizedAction(c *cli.Context) {
	client := Get80Client(c)

	opts := options(c)
	v80.Trace = opts.Trace
	v80.Verbose = opts.Verbose

	if opts.DryRun {
		log.Printf("--%s enabled.  agents will not be deleted.\n", FLAG_DRY_RUN)
	}

	_, err := client.RemoveDeauthorizedAgents(opts.DryRun)
	if err != nil {
		log.Fatalln(err)
	}
}
