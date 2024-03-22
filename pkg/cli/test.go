package cli

import (
	"os"
	"path/filepath"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/overseer/pkg/domain/model"
	"github.com/m-mizutani/overseer/pkg/domain/types"
	"github.com/m-mizutani/overseer/pkg/usecase"
	"github.com/m-mizutani/overseer/pkg/utils"
	"github.com/urfave/cli/v2"
)

func runTest() *cli.Command {
	var (
		queryDir     string
		emulatorPath string
	)
	return &cli.Command{
		Name:    "test",
		Usage:   `Run overseer test`,
		Aliases: []string{"t"},
		Flags: mergeFlags([]cli.Flag{
			&cli.StringFlag{
				Name:        "query-dir",
				Usage:       "Directory path of query files",
				Category:    "query",
				Destination: &queryDir,
				Aliases:     []string{"d"},
				EnvVars:     []string{"OVERSEER_QUERY_DIR"},
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "emulator-path",
				Aliases:     []string{"e"},
				Usage:       "Path of bigquery-emulator",
				Category:    "emulator",
				Destination: &emulatorPath,
				EnvVars:     []string{"OVERSEER_EMULATOR_PATH"},
				Value:       "bigquery-emulator",
			},
		}),
		Action: func(c *cli.Context) error {
			queryFiles, err := listQueryFiles(queryDir)
			if err != nil {
				return err
			}

			if len(queryFiles) == 0 {
				return goerr.Wrap(types.ErrInvalidOption, "No query files")
			}

			for _, queryFile := range queryFiles {
				fd, err := os.Open(queryFile)
				if err != nil {
					return goerr.Wrap(err, "Fail to open query file").With("file", queryFile)
				}
				defer utils.SafeClose(fd)

				task, err := model.NewTask(fd)
				if err != nil {
					return err
				}

				ctx := utils.CtxWithCWD(c.Context, filepath.Dir(queryFile))
				if err := usecase.RunTest(ctx, emulatorPath, task); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
