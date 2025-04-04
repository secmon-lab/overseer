package cli

import (
	"context"

	"github.com/secmon-lab/overseer/pkg/adaptor"
	"github.com/secmon-lab/overseer/pkg/cli/config/bq"
	"github.com/secmon-lab/overseer/pkg/cli/config/cache"
	"github.com/secmon-lab/overseer/pkg/cli/config/policy"
	"github.com/secmon-lab/overseer/pkg/cli/config/query"
	"github.com/secmon-lab/overseer/pkg/domain/model"
	"github.com/secmon-lab/overseer/pkg/domain/types"
	"github.com/secmon-lab/overseer/pkg/logging"
	"github.com/secmon-lab/overseer/pkg/usecase"
	"github.com/urfave/cli/v3"
)

func cmdFetch() *cli.Command {
	var (
		queryCfg  query.Config
		policyCfg policy.Config
		cacheCfg  cache.Config
		bqCfg     bq.Config
		queryIDs  []string
	)

	flags := []cli.Flag{}
	flags = append(flags, queryCfg.Flags()...)
	flags = append(flags, policyCfg.Flags()...)
	flags = append(flags, cacheCfg.Flags()...)
	flags = append(flags, bqCfg.Flags()...)
	flags = append(flags, &cli.StringSliceFlag{
		Name:        "query-id",
		Usage:       "Query ID",
		Destination: &queryIDs,
		Aliases:     []string{"i"},
	})

	return &cli.Command{
		Name:    "fetch",
		Aliases: []string{"f"},
		Usage:   "Query data and save the result into cache",
		Flags:   flags,
		Action: func(ctx context.Context, c *cli.Command) error {
			logger := logging.FromCtx(ctx)
			ctx, jobID := types.NewJobID(ctx)
			logger = logger.With("job_id", jobID)
			ctx = logging.InjectCtx(ctx, logger)

			logger.Info("Start overseer(fetch)", "query", queryCfg, "policy", policyCfg, "cache", cacheCfg, "bq", bqCfg, "query-id", queryIDs)

			cacheSvc, err := cacheCfg.Build(ctx, jobID)
			if err != nil {
				return err
			}

			queries, err := queryCfg.Build()
			if err != nil {
				return err
			}
			logger.Debug("Loaded queries", "queries", queries.IDs())

			bqClient, err := bqCfg.Build(ctx)
			if err != nil {
				return err
			}

			if policySvc, err := policyCfg.Build(); err != nil {
				return err
			} else if policySvc != nil {
				filtered, err := policySvc.SelectRequiredQueries(queries)
				if err != nil {
					return err
				}

				logging.Default().Info("Select required queries by policy",
					"before", len(queries),
					"after", len(filtered),
				)
				logger.Debug("Select required queries by policy", "before", queries.IDs(), "after", filtered.IDs())
				queries = filtered
			}

			if len(queryIDs) > 0 {
				filtered := model.Queries{}
				for _, id := range queryIDs {
					if q := queries.FindByID(model.QueryID(id)); q != nil {
						filtered = append(filtered, q)
					}
				}
				logger.Debug("Select queries by query-id", "before", queries.IDs(), "after", filtered.IDs())
				queries = filtered
			}

			uc := usecase.New(adaptor.New(adaptor.WithBigQuery(bqClient)))

			if err := uc.Fetch(ctx, queries, cacheSvc); err != nil {
				return err
			}
			return nil
		},
	}
}
