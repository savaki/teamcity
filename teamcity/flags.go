package main

import "github.com/codegangsta/cli"

const (
	FLAG_URL      = "url"
	FLAG_USERNAME = "username"
	FLAG_PASSWORD = "password"
	FLAG_VERBOSE  = "verbose"
	FLAG_TRACE    = "trace"

	FLAG_AGENT_ID        = "agent-id"
	FLAG_AGENT_NAME      = "agent-name"
	FLAG_AGENT_POOL_NAME = "agent-pool-name"
	FLAG_ALL_AGENTS      = "all-agents"

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

	FlagAgentId       = cli.StringSliceFlag{FLAG_AGENT_ID, &cli.StringSlice{}, "filter agent name (regexp)"}
	FlagAgentName     = cli.StringSliceFlag{FLAG_AGENT_NAME, &cli.StringSlice{}, "filter by agent name (regexp)"}
	FlagAgentPoolName = cli.StringFlag{FLAG_AGENT_POOL_NAME, "", "filter by agent pool name"}
	FlagAllAgents     = cli.BoolFlag{FLAG_ALL_AGENTS, "include all agents"}

	FlagBuildId      = cli.StringFlag{FLAG_BUILD_ID, "", "the build to retrieve details for"}
	FlagBuildTypeId  = cli.StringFlag{FLAG_BUILD_TYPE_ID, "", "the build type id"}
	FlagArtifactName = cli.StringFlag{FLAG_BUILD_ARTIFACT_NAME, "", "the filename of the artifact to download"}

	FlagLast = cli.IntFlag{FLAG_LAST, 1, "how many builds to retrieve"}
)

type Options struct {
	Verbose bool
	Trace   bool
}

func options(c *cli.Context) Options {
	return Options{
		Verbose: c.Bool(FLAG_VERBOSE) || c.Bool(FLAG_TRACE),
		Trace:   c.Bool(FLAG_TRACE),
	}
}
