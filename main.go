package main

import (
	"os"

	"github.com/cooldrip/jhop/cmd"
)

type cat struct {
	*cat
}

func main() {
	cmd.Run(os.Args)
}
