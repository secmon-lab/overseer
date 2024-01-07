package cli

import (
	"os"
	"path/filepath"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/overseer/pkg/cli/config"
	"github.com/m-mizutani/overseer/pkg/domain/model"
	"github.com/m-mizutani/overseer/pkg/domain/types"
	"github.com/m-mizutani/overseer/pkg/infra"
	"github.com/m-mizutani/overseer/pkg/usecase"
	"github.com/m-mizutani/overseer/pkg/utils"
	"github.com/urfave/cli/v2"
)

func runCommand() *cli.Command {
	var (
		queryDir  string
		recursive bool
		bq        config.BigQuery
		pubsub    config.PubSub
	)
	return &cli.Command{
		Name:    "run",
		Usage:   `Run overseer task`,
		Aliases: []string{"r"},
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
			&cli.BoolFlag{
				Name:        "recursive",
				Usage:       "Recursively search query files",
				Category:    "query",
				Destination: &recursive,
				EnvVars:     []string{"OVERSEER_RECURSIVE"},
				Value:       false,
			},
		}, bq.Flags(), pubsub.Flags()),
		Action: func(ctx *cli.Context) error {
			queryFiles, err := listQueryFiles(queryDir, recursive)
			if err != nil {
				return err
			}

			if len(queryFiles) == 0 {
				return goerr.Wrap(types.ErrInvalidOption, "No query files")
			}

			bqClient, err := bq.Configure(ctx.Context)
			if err != nil {
				return err
			}
			pubsubClient, err := pubsub.Configure(ctx.Context)
			if err != nil {
				return err
			}

			clients := infra.New(bqClient, pubsubClient)

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

				if err := usecase.RunTask(ctx.Context, clients, task); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func listQueryFiles(queryDir string, recursive bool) ([]string, error) {
	entries, err := os.ReadDir(queryDir)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to read query directory")
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			if recursive {
				subFiles, err := listQueryFiles(filepath.Join(queryDir, e.Name()), recursive)
				if err != nil {
					return nil, err
				}
				files = append(files, subFiles...)
			}
		} else {
			files = append(files, filepath.Join(queryDir, e.Name()))
		}
	}

	return files, nil
}
