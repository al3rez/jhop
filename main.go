package main

import (
	"os"

	"github.com/azbshiri/jhop/cmd"
)

type cat struct {
	*cat
}

func main() {
	cmd.Run(os.Args)
}
