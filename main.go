package main

import (
	"github.com/YangTaeyoung/actions-runner-starter/action"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.App{
		Name:  "actions-runner-starter",
		Usage: "actions-runner-starter",
		Commands: []cli.Command{
			{
				Name:   "register",
				Usage:  "Register multiple self-Hosted actions runners into github",
				Action: action.Register,
			},
			{
				Name:  "unregister",
				Usage: "Unregister all self-hosted actions runners that were registered by actions-runner-starter",
			},
			{
				Name:  "serve",
				Usage: "Serve all self-hosted actions runners that were registered by actions-runner-starter",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
