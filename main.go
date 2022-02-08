package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "sshm",
		Usage: "Simple but powerful cli manager for your ssh connections.",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List all connections",
				Action:  ListEntries,
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Add a new connection",
				Action:  NewConnection,
			},
			{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "Deletes an existing connection",
				Action:  DeletConnection,
			},
			{
				Name:    "connect",
				Aliases: []string{"c"},
				Usage:   "Connect",
				Action:  Connect,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
