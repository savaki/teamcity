package main

import "github.com/codegangsta/cli"

var serverCommand = cli.Command{
	Name:  "server",
	Usage: "retrieve info on the server",
	Flags: []cli.Flag{
		FlagUrl,
		FlagUsername,
		FlagPassword,
		FlagVerbose,
	},
	Action: func(c *cli.Context) {
		client := Get80Client(c)
		server, err := client.Server()
		Print(server, err)
	},
}
