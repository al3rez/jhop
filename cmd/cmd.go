package cmd

import (
	"log"

	"github.com/cooldrip/jhop/api"
	"github.com/urfave/cli"
)

func Run(args []string) {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port",
			Value: "6000",
			Usage: "Set port",
		},
		cli.StringFlag{
			Name:  "host",
			Value: "localhost",
			Usage: "Set host",
		},
	}
	app.Name = "jhop"
	app.Usage = "Create fake REST API in one sec."
	app.Action = func(c *cli.Context) error {
		var filename string
		if c.NArg() > 0 {
			filename = c.Args()[0]
		}
		if filename != "" {
			api.Create(filename, c.String("host"), c.String("port"))
		}
		return nil
	}

	err := app.Run(args)
	if err != nil {
		log.Fatal(err)
	}
}
