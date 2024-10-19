package cli

import (
	"context"
	"io"
	"os"

	"github.com/secmon-as-code/overseer/pkg/adaptor"
	"github.com/secmon-as-code/overseer/pkg/cli/config/policy"
	"github.com/secmon-as-code/overseer/pkg/cli/config/query"
	"github.com/secmon-as-code/overseer/pkg/logging"
	"github.com/secmon-as-code/overseer/pkg/usecase"

	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v3"
)

func cmdInspect() *cli.Command {
	var (
		policyCfg policy.Config
		queryCfg  query.Config
		output    string
	)

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "out",
			Usage:       "Output destination [stdout, stderr, <file_name>]",
			Destination: &output,
			Required:    false,
			Aliases:     []string{"o"},
			Value:       "stdout",
		},
	}
	flags = append(flags, policyCfg.Flags()...)
	flags = append(flags, queryCfg.Flags()...)

	return &cli.Command{
		Name:  "inspect",
		Usage: "Inspect policy and query",
		Flags: flags,
		Action: func(ctx context.Context, c *cli.Command) error {
			logging.Default().Info("Inspecting policy and query",
				"policy", policyCfg.FilePath(),
				"query", queryCfg.FilePath(),
			)

			policy, err := policyCfg.Build()
			if err != nil {
				return err
			}

			queries, err := queryCfg.Build()
			if err != nil {
				return err
			}

			uc := usecase.New(adaptor.New(adaptor.WithPolicy(policy)))

			var w io.Writer

			switch output {
			case "stdout":
				w = os.Stdout
			case "stderr":
				w = os.Stderr
			default:
				f, err := os.Create(output)
				if err != nil {
					return goerr.Wrap(err, "fail to open output file")
				}
				defer f.Close()
				w = f
			}

			return uc.Inspect(ctx, queries, w)
		},
	}
}
