package main

import "github.com/codegangsta/cli"

const (
	FLAG_URL      = "url"
	FLAG_USERNAME = "username"
	FLAG_PASSWORD = "password"
	FLAG_VERBOSE  = "verbose"
	FLAG_TRACE    = "trace"
	FLAG_DRY_RUN  = "dry-run"
	FLAG_LONG     = "long"

	FLAG_AGENT_ID          = "agent-id"
	FLAG_AGENT_NAME        = "agent-name"
	FLAG_AGENT_POOL_NAME   = "agent-pool-name"
	FLAG_AGENT_PROPERTY    = "agent-property"
	FLAG_ALL_AGENTS        = "all-agents"
	FLAG_DISCONNECTED_ONLY = "disconnected-only"

	FLAG_BUILD_ID            = "build-id"
	FLAG_BUILD_TYPE_ID       = "build-type-id"
	FLAG_BUILD_ARTIFACT_NAME = "artifact-name"

	FLAG_LAST = "last"
)

var (
	FlagUrl      = cli.StringFlag{FLAG_URL, "", "url of the TeamCity server"}
	FlagUsername = cli.StringFlag{FLAG_USERNAME, "", "TeamCity username"}
	FlagPassword = cli.StringFlag{FLAG_PASSWORD, "", "TeamCity password"}
	FlagVerbose  = cli.BoolFlag{FLAG_VERBOSE, "additional content"}
	FlagTrace    = cli.BoolFlag{FLAG_TRACE, "developer level details"}
	FlagDryRun   = cli.BoolFlag{FLAG_DRY_RUN, "dry-run, don't execute anything"}
	FlagLong     = cli.BoolFlag{FLAG_LONG, "show additional details"}

	FlagAgentId          = cli.StringSliceFlag{FLAG_AGENT_ID, &cli.StringSlice{}, "filter agent name (regexp)"}
	FlagAgentName        = cli.StringSliceFlag{FLAG_AGENT_NAME, &cli.StringSlice{}, "filter by agent name (regexp)"}
	FlagAgentProperty    = cli.StringSliceFlag{FLAG_AGENT_PROPERTY, &cli.StringSlice{}, "filter by agent property ex. teamcity.agent.name=[12]"}
	FlagAgentPoolName    = cli.StringFlag{FLAG_AGENT_POOL_NAME, "", "specify an agent pool by name (regexp)"}
	FlagAllAgents        = cli.BoolFlag{FLAG_ALL_AGENTS, "include all agents"}
	FlagDisconnectedOnly = cli.BoolFlag{FLAG_DISCONNECTED_ONLY, "include disconnected agents only"}

	FlagBuildId      = cli.StringFlag{FLAG_BUILD_ID, "", "the build to retrieve details for"}
	FlagBuildTypeId  = cli.StringFlag{FLAG_BUILD_TYPE_ID, "", "the build type id"}
	FlagArtifactName = cli.StringFlag{FLAG_BUILD_ARTIFACT_NAME, "", "the filename of the artifact to download"}

	FlagLast = cli.IntFlag{FLAG_LAST, 1, "how many builds to retrieve"}
)

type Options struct {
	Verbose bool
	Trace   bool
	DryRun  bool
}

func options(c *cli.Context) Options {
	return Options{
		Verbose: c.Bool(FLAG_VERBOSE) || c.Bool(FLAG_TRACE),
		Trace:   c.Bool(FLAG_TRACE),
		DryRun:  c.Bool(FLAG_DRY_RUN),
	}
}
