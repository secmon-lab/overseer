package cli

import "github.com/urfave/cli/v3"

func cmdFetch() *cli.Command {
	return &cli.Command{
		Name:  "fetch",
		Usage: "Fetch security data",
	}
}
