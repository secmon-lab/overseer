package main

import (
	"os"

	"github.com/secmon-as-code/overseer/pkg/cli"
)

func main() {
	if err := cli.New().Run(os.Args); err != nil {
		os.Exit(1)
	}
}
