package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cooldrip/jhop"
	"github.com/gorilla/handlers"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func main() {
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
		var filenames []string
		if c.NArg() > 0 {
			filenames = c.Args()
		}
		if len(filenames) > 0 {
			files := make([]io.Reader, len(filenames))
			for i, filename := range filenames {
				f, err := os.Open(filename)
				if err != nil {
					return errors.Wrapf(err, "failed to open file %s", filename)
				}
				files[i] = f
			}

			handler, err := jhop.NewHandler(files...)
			if err != nil {
				return errors.Wrap(err, "failed to initialize handler")
			}

			addr := fmt.Sprintf("%s:%s", c.String("host"), c.String("port"))
			return http.ListenAndServe(addr, handlers.LoggingHandler(os.Stdout, handler))
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("run app: %s", err)
	}
}
