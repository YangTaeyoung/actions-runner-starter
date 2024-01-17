package main

import (
	"github.com/YangTaeyoung/actions-runner-starter/runner"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	r := runner.New()

	app := cli.App{
		Name:  "actions-runner-starter",
		Usage: "actions-runner-starter",
		Commands: []cli.Command{
			{
				Name:  "configure",
				Usage: "configure for self-hosted actions runners",
				Action: func(c *cli.Context) error {
					return r.Configure()
				},
			},
			{
				Name:  "register",
				Usage: "Register multiple self-Hosted actions runners into github",
				Action: func(c *cli.Context) error {
					return r.Register()
				},
			},
			{
				Name:  "unregister",
				Usage: "Unregister all self-hosted actions runners that were registered by actions-runner-starter",
				Action: func(c *cli.Context) error {
					return r.Unregister()
				},
			},
			{
				Name:  "serve",
				Usage: "Serve all self-hosted actions runners that were registered by actions-runner-starter",
				Action: func(c *cli.Context) error {
					return r.Serve()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
