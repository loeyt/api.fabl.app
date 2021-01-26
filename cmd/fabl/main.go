package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                   "factorio-blueprints",
		Usage:                  "save and share blueprints",
		HideHelpCommand:        true,
		UseShortOptionHandling: true,

		Commands: []*cli.Command{
			serverCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
