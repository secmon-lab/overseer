package sentry

import (
	"context"
	"log/slog"

	"github.com/getsentry/sentry-go"
	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/overseer/pkg/logging"
	"github.com/urfave/cli/v3"
)

type Config struct {
	// Sentry DSN
	dsn string
	env string
}

func (x *Config) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "sentry-dsn",
			Usage:       "Sentry DSN",
			Category:    "sentry",
			Destination: &x.dsn,
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_SENTRY_DSN")),
		},
		&cli.StringFlag{
			Name:        "sentry-env",
			Usage:       "Sentry environment",
			Category:    "sentry",
			Destination: &x.env,
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_SENTRY_ENV")),
		},
	}
}

func (x Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("sentry-dsn", x.dsn),
		slog.String("sentry-env", x.env),
	)
}

func (x *Config) Build(ctx context.Context) error {
	if x.dsn == "" {
		logging.FromCtx(ctx).Warn("Sentry DSN is not set")
		return nil
	}

	// Init sentry
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         x.dsn,
		Environment: x.env,
	}); err != nil {
		return goerr.Wrap(err, "Fail to initialize Sentry")
	}
	return nil
}
