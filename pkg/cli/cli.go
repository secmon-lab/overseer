package cli

import (
	"context"

	"github.com/secmon-lab/overseer/pkg/cli/config/logger"
	"github.com/secmon-lab/overseer/pkg/cli/config/sentry"
	"github.com/secmon-lab/overseer/pkg/domain/types"
	"github.com/secmon-lab/overseer/pkg/logging"

	"github.com/urfave/cli/v3"
)

type CLI struct {
	app *cli.Command
}

func concat[A any](x ...[]A) []A {
	var ret []A
	for _, a := range x {
		ret = append(ret, a...)
	}
	return ret
}

func New() *CLI {
	var (
		loggerCfg logger.Config
		sentryCfg sentry.Config
	)

	app := &cli.Command{
		Name:    "overseer",
		Version: types.AppVersion,
		Usage:   "Overseer is security data analysis framework",
		Flags: concat(
			loggerCfg.Flags(),
			sentryCfg.Flags(),
		),

		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			logger, err := loggerCfg.Build()
			if err != nil {
				return nil, err
			}
			logging.SetDefault(logger)

			if err := sentryCfg.Build(ctx); err != nil {
				return nil, err
			}

			return ctx, nil
		},

		Commands: []*cli.Command{
			cmdInspect(),
			cmdFetch(),
			cmdEval(),
			cmdRun(),
		},
	}

	return &CLI{app: app}
}

func (x *CLI) Run(args []string) error {
	ctx, _ := types.NewRequestID(context.Background())
	if err := x.app.Run(ctx, args); err != nil {
		handleError(ctx, err)
		return err
	}

	return nil
}
