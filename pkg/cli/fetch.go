package cli

import (
	"context"

	"github.com/secmon-as-code/overseer/pkg/adaptor"
	"github.com/secmon-as-code/overseer/pkg/adaptor/bq"
	"github.com/secmon-as-code/overseer/pkg/cli/config/cache"
	"github.com/secmon-as-code/overseer/pkg/cli/config/policy"
	"github.com/secmon-as-code/overseer/pkg/cli/config/query"
	"github.com/secmon-as-code/overseer/pkg/domain/model"
	"github.com/secmon-as-code/overseer/pkg/logging"
	"github.com/secmon-as-code/overseer/pkg/usecase"
	"github.com/urfave/cli/v3"
)

func cmdFetch() *cli.Command {
	var (
		queryCfg    query.Config
		policyCfg   policy.Config
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
	flags = append(flags, policyCfg.Flags()...)
	flags = append(flags, cacheCfg.Flags()...)

	action := func(ctx context.Context, c *cli.Command) error {
		id := model.NewJobID()
		ctx = logging.InjectCtx(ctx, logging.Default().With("job_id", id))

		cacheSvc, err := cacheCfg.Build(ctx, id)
		if err != nil {
			return err
		}

		allQueries, err := queryCfg.Build()
		if err != nil {
			return err
		}

		bqClient, err := bq.New(ctx, bqProjectID)
		if err != nil {
			return err
		}

		policySvc, err := policyCfg.Build()
		if err != nil {
			return err
		}

		queries := policySvc.SelectRequiredQueries(allQueries)

		uc := usecase.New(adaptor.New(adaptor.WithBigQuery(bqClient)))

		return uc.Fetch(ctx, queries, cacheSvc)
	}

	return &cli.Command{
		Name:    "fetch",
		Aliases: []string{"f"},
		Usage:   "Query data and save the result into cache",
		Flags:   flags,
		Action:  action,
	}
}
