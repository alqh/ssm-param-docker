package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/alqh/ssm-param-docker/cmd"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			cmd.ExecCLI(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
