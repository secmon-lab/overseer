package bq

import (
	"context"
	"log/slog"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/overseer/pkg/adaptor/bq"
	"github.com/urfave/cli/v3"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

type Config struct {
	projectID                 string
	impersonateServiceAccount string
}

func (x *Config) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "bq-project-id",
			Usage:       "BigQuery project ID",
			Category:    "bq",
			Destination: &x.projectID,
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_BQ_PROJECT_ID")),
		},
		&cli.StringFlag{
			Name:        "bq-impersonate-service-account",
			Usage:       "Impersonate service account",
			Category:    "bq",
			Destination: &x.impersonateServiceAccount,
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_BQ_IMPERSONATE_SERVICE_ACCOUNT")),
		},
	}
}

func (x *Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("bq-project-id", x.projectID),
		slog.String("bq-impersonate-service-account", x.impersonateServiceAccount),
	)
}

func (x *Config) ProjectID() string {
	return x.projectID
}

func (x *Config) ImpersonateServiceAccount() string {
	return x.impersonateServiceAccount
}

func (x *Config) Build(ctx context.Context) (*bq.Client, error) {
	var bqOptions []option.ClientOption
	if x.impersonateServiceAccount != "" {
		ts, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
			TargetPrincipal: x.impersonateServiceAccount,
			Scopes: []string{
				"https://www.googleapis.com/auth/bigquery",
				"https://www.googleapis.com/auth/cloud-platform",
				"https://www.googleapis.com/auth/bigquery.readonly",
				"https://www.googleapis.com/auth/cloud-platform.read-only",
			},
		})
		if err != nil {
			return nil, goerr.Wrap(err, "failed to create token source for impersonate")
		}
		bqOptions = append(bqOptions, option.WithTokenSource(ts))
	}

	bqClient, err := bq.New(ctx, x.projectID, bqOptions...)
	if err != nil {
		return nil, err
	}

	return bqClient, nil
}
