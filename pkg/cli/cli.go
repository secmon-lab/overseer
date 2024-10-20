package cli

import (
	"context"

	"github.com/secmon-as-code/overseer/pkg/cli/config/logger"
	"github.com/secmon-as-code/overseer/pkg/domain/types"
	"github.com/secmon-as-code/overseer/pkg/logging"

	"github.com/urfave/cli/v3"
)

type CLI struct {
	app *cli.Command
}

func New() *CLI {
	var (
		loggerCfg logger.Config
	)

	var flags []cli.Flag
	flags = append(flags, loggerCfg.Flags()...)

	app := &cli.Command{
		Name:    "overseer",
		Version: types.AppVersion,
		Usage:   "Overseer is security data analysis framework",
		Flags:   flags,

		Before: func(ctx context.Context, c *cli.Command) error {
			logger, err := loggerCfg.Build()
			if err != nil {
				return err
			}

			logging.SetDefault(logger)
			return nil
		},

		Commands: []*cli.Command{
			cmdInspect(),
			cmdFetch(),
			cmdEval(),
		},
	}

	return &CLI{app: app}
}

func (x *CLI) Run(args []string) error {
	if err := x.app.Run(context.Background(), args); err != nil {
		logging.Default().Error("fail to run command", "err", err)
		return err
	}

	return nil
}
