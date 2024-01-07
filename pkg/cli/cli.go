package cli

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/overseer/pkg/cli/config"
	"github.com/m-mizutani/overseer/pkg/domain/types"
	"github.com/m-mizutani/overseer/pkg/utils"
	"github.com/urfave/cli/v2"
)

func Run(argv []string) error {
	var (
		logCfg    config.Logger
		logCloser func()
	)
	app := cli.App{
		Name:    "overseer",
		Version: types.AppVersion,
		Flags:   logCfg.Flags(),
		Commands: []*cli.Command{
			runCommand(),
		},
		Before: func(ctx *cli.Context) error {
			f, err := logCfg.Configure()
			if err != nil {
				return err
			}
			logCloser = f
			return nil
		},
		After: func(ctx *cli.Context) error {
			logCloser()
			return nil
		},
	}

	if err := app.Run(argv); err != nil {
		utils.Logger().Error("Fail to run overseer", "error", err)
		return goerr.Wrap(err, "Fail to run overseer")
	}

	return nil
}

func mergeFlags(flags ...[]cli.Flag) []cli.Flag {
	var merged []cli.Flag
	for _, f := range flags {
		merged = append(merged, f...)
	}
	return merged
}
