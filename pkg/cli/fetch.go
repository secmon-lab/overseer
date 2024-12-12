package cli

import (
	"context"

	"github.com/secmon-lab/overseer/pkg/adaptor"
	"github.com/secmon-lab/overseer/pkg/cli/config/bq"
	"github.com/secmon-lab/overseer/pkg/cli/config/cache"
	"github.com/secmon-lab/overseer/pkg/cli/config/policy"
	"github.com/secmon-lab/overseer/pkg/cli/config/query"
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
	)

	flags := []cli.Flag{}
	flags = append(flags, queryCfg.Flags()...)
	flags = append(flags, policyCfg.Flags()...)
	flags = append(flags, cacheCfg.Flags()...)
	flags = append(flags, bqCfg.Flags()...)

	action := func(ctx context.Context, c *cli.Command) error {
		logger := logging.FromCtx(ctx)
		ctx, jobID := types.NewJobID(ctx)
		logger = logger.With("job_id", jobID)
		ctx = logging.InjectCtx(ctx, logger)

		logger.Info("Start overseer(fetch)", "query", queryCfg, "policy", policyCfg, "cache", cacheCfg, "bq", bqCfg)

		cacheSvc, err := cacheCfg.Build(ctx, jobID)
		if err != nil {
			return err
		}

		queries, err := queryCfg.Build()
		if err != nil {
			return err
		}

		bqClient, err := bqCfg.Build(ctx)
		if err != nil {
			return err
		}

		if policySvc, err := policyCfg.Build(); err != nil {
			return err
		} else if policySvc != nil {
			filtered := policySvc.SelectRequiredQueries(queries)
			logging.Default().Info("Select required queries by policy",
				"before", len(queries),
				"after", len(filtered),
			)
			queries = filtered
		}

		uc := usecase.New(adaptor.New(adaptor.WithBigQuery(bqClient)))

		if err := uc.Fetch(ctx, queries, cacheSvc); err != nil {
			return err
		}
		return nil
	}

	return &cli.Command{
		Name:    "fetch",
		Aliases: []string{"f"},
		Usage:   "Query data and save the result into cache",
		Flags:   flags,
		Action:  action,
	}
}
