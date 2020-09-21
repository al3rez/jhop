package main

import (
	"encoding/json"
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
		cli.StringFlag{
			Name:  "routes",
			Usage: "Set routes file",
		},
	}
	app.Name = "jhop"
	app.Usage = "Create fake REST API in one sec."
	app.Action = func(c *cli.Context) error {
		var filenames []string
		if c.NArg() > 0 {
			filenames = c.Args()
		}

		if len(filenames) == 0 {
			return errors.New("no files listed as arguments")
		}

		files := make([]io.ReadCloser, len(filenames))
		for i, filename := range filenames {
			f, err := os.Open(filename)
			if err != nil {
				return errors.Wrap(err, "failed to open file")
			}
			files[i] = f
		}

		routes := make(map[string]string)
		if c.String("routes") != "" {
			f, err := os.Open(c.String("routes"))
			defer f.Close()
			if err != nil {
				return errors.Wrap(err, "failed to open file")
			}
			if err := json.NewDecoder(f).Decode(&routes); err != nil {
				return errors.Wrap(err, "failed to unmarshal routes")
			}
		}

		handler, err := jhop.NewHandlerWithRoutes(routes, files...)
		if err != nil {
			return errors.Wrap(err, "failed to initialize handler")
		}

		addr := fmt.Sprintf("%s:%s", c.String("host"), c.String("port"))
		log.Printf("starting server on: %s\n", addr)
		return http.ListenAndServe(addr, handlers.LoggingHandler(os.Stdout, handler))
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("run app: %s", err)
	}
}
