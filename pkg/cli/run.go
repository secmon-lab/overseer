package cli

import (
	"context"

	"github.com/m-mizutani/goerr/v2"
	"github.com/secmon-lab/overseer/pkg/adaptor"
	"github.com/secmon-lab/overseer/pkg/cli/config/bq"
	"github.com/secmon-lab/overseer/pkg/cli/config/cache"
	"github.com/secmon-lab/overseer/pkg/cli/config/notify"
	"github.com/secmon-lab/overseer/pkg/cli/config/policy"
	"github.com/secmon-lab/overseer/pkg/cli/config/query"
	"github.com/secmon-lab/overseer/pkg/domain/types"
	"github.com/secmon-lab/overseer/pkg/logging"
	"github.com/secmon-lab/overseer/pkg/usecase"
	"github.com/urfave/cli/v3"
)

func cmdRun() *cli.Command {
	var (
		queryCfg  query.Config
		policyCfg policy.Config
		cacheCfg  cache.Config
		notifyCfg notify.Config
		bqCfg     bq.Config
	)

	flags := []cli.Flag{}
	flags = append(flags, queryCfg.Flags()...)
	flags = append(flags, policyCfg.Flags()...)
	flags = append(flags, cacheCfg.Flags()...)
	flags = append(flags, notifyCfg.Flags()...)
	flags = append(flags, bqCfg.Flags()...)

	action := func(ctx context.Context, c *cli.Command) error {
		logger := logging.FromCtx(ctx)
		ctx, jobID := types.NewJobID(ctx)
		logger = logger.With("job_id", jobID)
		ctx = logging.InjectCtx(ctx, logger)

		logger.Info("Start overseer(run)",
			"query", queryCfg,
			"policy", policyCfg,
			"cache", cacheCfg,
			"notify", notifyCfg,
			"bq", bqCfg,
		)

		cacheSvc, err := cacheCfg.Build(ctx, jobID)
		if err != nil {
			return err
		}

		allQueries, err := queryCfg.Build()
		if err != nil {
			return err
		}

		policySvc, err := policyCfg.Build()
		if err != nil {
			return err
		}
		if policySvc == nil {
			return goerr.New("policy config is required")
		}

		notifySvc, err := notifyCfg.Build()
		if err != nil {
			return err
		}

		bqClient, err := bqCfg.Build(ctx)
		if err != nil {
			return err
		}

		uc := usecase.New(adaptor.New(
			adaptor.WithBigQuery(bqClient),
		))

		queries, err := policySvc.SelectRequiredQueries(allQueries)
		if err != nil {
			return err
		}

		if err := uc.Fetch(ctx, queries, cacheSvc); err != nil {
			return err
		}

		if err := uc.Eval(ctx, policySvc, cacheSvc, notifySvc); err != nil {
			handleError(ctx, err)
			return err
		}

		return nil
	}

	return &cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "Run the overseer (fetch -> eval)",
		Flags:   flags,
		Action:  action,
	}
}
