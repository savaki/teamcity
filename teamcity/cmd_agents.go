package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/savaki/teamcity/v80"
	"log"
	"strings"
)

var agentCommand = cli.Command{
	Name: "agent",
	Subcommands: []cli.Command{
		{
			Name: "list",
			Flags: []cli.Flag{
				FlagAgentName,
				FlagAgentId,
				FlagAgentProperty,
				FlagLong,
				FlagVerbose,
				FlagTrace,
			},
			Action: agentListAction,
		},
		{
			Name: "authorize",
			Flags: []cli.Flag{
				FlagAgentName,
				FlagAgentId,
				FlagAgentProperty,
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
				FlagAgentProperty,
				FlagDisconnectedOnly,
				FlagVerbose,
				FlagTrace,
			},
			Action: agentDeauthorizeAction,
		},
		{
			Name: "remove",
			Flags: []cli.Flag{
				FlagAgentName,
				FlagAgentId,
				FlagAgentProperty,
				FlagAllAgents,
				FlagVerbose,
				FlagTrace,
				FlagDryRun,
			},
			Action: agentRemoveDeauthorizedAction,
		},
		{
			Name: "assign-to-pool",
			Flags: []cli.Flag{
				FlagAgentPoolName,
				FlagAgentId,
				FlagAgentProperty,
				FlagAgentName,
				FlagAllAgents,
				FlagVerbose,
				FlagTrace,
				FlagDryRun,
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

	if c.Bool(FLAG_DISCONNECTED_ONLY) {
		value := fmt.Sprintf("%#v", false)
		filters = append(filters, v80.NewAgentFilter(value, v80.AgentConnectedAccessor))
	}

	// filter by agent property
	if values := c.StringSlice(FLAG_AGENT_PROPERTY); values != nil {
		for _, value := range values {
			parts := strings.Split(value, "=")
			if len(parts) != 2 {
				log.Fatalf("invalid value for --%s, %s.  expected something like --%s teamcity.agent-name=agent[12]", FLAG_AGENT_PROPERTY, value, FLAG_AGENT_PROPERTY)
			}
			accessor := v80.NewAgentPropertyAccessor(parts[0])
			filter := v80.NewAgentFilter(parts[1], accessor)
			filters = append(filters, filter)
		}
	}

	// filter by id
	if values := c.StringSlice(FLAG_AGENT_ID); values != nil {
		for _, value := range values {
			filter := v80.NewAgentFilter(value, v80.AgentIdAccessor)
			filters = append(filters, filter)
		}
	}

	// filter by name
	if values := c.StringSlice(FLAG_AGENT_NAME); values != nil {
		for _, value := range values {
			filter := v80.NewAgentFilter(value, v80.AgentNameAccessor)
			filters = append(filters, filter)
		}
	}

	return filters
}

func agentListAction(c *cli.Context) {
	client := Get80Client(c)
	filters := agentFilters(c)

	agents, err := client.FindAgents(filters)

	if !c.Bool(FLAG_LONG) {
		for _, agent := range agents {
			agent.Properties = nil
		}
	}

	Print(agents, err)
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
	v80.DryRun = opts.DryRun

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
