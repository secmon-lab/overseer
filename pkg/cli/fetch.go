package cli

import (
	"context"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-as-code/overseer/pkg/adaptor"
	"github.com/secmon-as-code/overseer/pkg/adaptor/bq"
	"github.com/secmon-as-code/overseer/pkg/cli/config/cache"
	"github.com/secmon-as-code/overseer/pkg/cli/config/query"
	"github.com/secmon-as-code/overseer/pkg/domain/model"
	"github.com/secmon-as-code/overseer/pkg/logging"
	"github.com/secmon-as-code/overseer/pkg/usecase"
	"github.com/urfave/cli/v3"
)

func cmdFetch() *cli.Command {
	var (
		queryCfg    query.Config
		cacheCfg    cache.Config
		bqProjectID string
	)

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "bigquery-project-id",
			Usage:       "BigQuery project ID",
			Category:    "fetch",
			Destination: &bqProjectID,
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_BIGQUERY_PROJECT_ID")),
			Required:    true,
		},
	}
	flags = append(flags, queryCfg.Flags()...)
	flags = append(flags, cacheCfg.Flags()...)

	return &cli.Command{
		Name:  "fetch",
		Usage: "Query data and save the result into cache",
		Flags: flags,
		Action: func(ctx context.Context, c *cli.Command) error {
			id := model.NewJobID()
			cacheSvc, err := cacheCfg.Build(ctx, id)
			if err != nil {
				return err
			}

			queries, err := queryCfg.Build()
			if err != nil {
				return err
			}

			bqClient, err := bq.New(ctx, bqProjectID)
			if err != nil {
				return err
			}

			uc := usecase.New(adaptor.New(adaptor.WithBigQuery(bqClient)))

			ctx = goerr.InjectValue(ctx, "job_id", id)
			ctx = logging.InjectCtx(ctx, logging.Default().With("job_id", id))

			return uc.Fetch(ctx, queries, cacheSvc)
		},
	}
}
