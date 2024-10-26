package cli

import (
	"context"

	"github.com/secmon-lab/overseer/pkg/adaptor"
	"github.com/secmon-lab/overseer/pkg/cli/config/cache"
	"github.com/secmon-lab/overseer/pkg/cli/config/notify"
	"github.com/secmon-lab/overseer/pkg/cli/config/policy"
	"github.com/secmon-lab/overseer/pkg/domain/model"
	"github.com/secmon-lab/overseer/pkg/logging"
	"github.com/secmon-lab/overseer/pkg/usecase"
	"github.com/urfave/cli/v3"
)

func cmdEval() *cli.Command {
	var (
		policyCfg policy.Config
		cacheCfg  cache.Config
		notifyCfg notify.Config
		jobID     model.JobID
	)

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "job-id",
			Aliases:     []string{"i"},
			Usage:       "Job ID",
			Category:    "eval",
			Destination: (*string)(&jobID),
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_JOB_ID")),
			Required:    true,
		},
	}
	flags = append(flags, policyCfg.Flags()...)
	flags = append(flags, cacheCfg.Flags()...)
	flags = append(flags, notifyCfg.Flags()...)

	action := func(ctx context.Context, c *cli.Command) error {
		ctx = logging.InjectCtx(ctx, logging.Default().With("job_id", jobID))

		cacheSvc, err := cacheCfg.Build(ctx, jobID)
		if err != nil {
			return err
		}

		policySvc, err := policyCfg.Build()
		if err != nil {
			return err
		}

		notifySvc, err := notifyCfg.Build()
		if err != nil {
			return err
		}

		uc := usecase.New(adaptor.New())

		return uc.Eval(ctx, policySvc, cacheSvc, notifySvc)
	}

	return &cli.Command{
		Name:    "eval",
		Aliases: []string{"e"},
		Usage:   "Query data and save the result into cache",
		Flags:   flags,
		Action:  action,
	}
}
