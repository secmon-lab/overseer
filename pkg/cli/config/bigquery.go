package config

import (
	"context"
	"log/slog"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/overseer/pkg/domain/interfaces"
	"github.com/m-mizutani/overseer/pkg/domain/types"
	"github.com/m-mizutani/overseer/pkg/infra/bq"
	"github.com/m-mizutani/overseer/pkg/utils"
	"github.com/urfave/cli/v2"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

type BigQuery struct {
	projectID string
	saKeyData string `masq:"secret"`
	saKeyFile string

	impersonateServiceAccount string
}

func (x *BigQuery) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "bq-project-id",
			Usage:       "BigQuery project ID",
			Destination: &x.projectID,
			EnvVars:     []string{"OVERSEER_BIGQUERY_PROJECT_ID"},
		},
		&cli.StringFlag{
			Name:        "bq-impersonate-service-account",
			Usage:       "Impersonate service account for BigQuery",
			Destination: &x.impersonateServiceAccount,
			EnvVars:     []string{"OVERSEER_BIGQUERY_IMPERSONATE_SERVICE_ACCOUNT"},
		},

		&cli.StringFlag{
			Name:        "bq-sa-key-data",
			Usage:       "BigQuery service account key data",
			Destination: &x.saKeyData,
			EnvVars:     []string{"OVERSEER_BIGQUERY_SA_KEY_DATA"},
		},
		&cli.StringFlag{
			Name:        "bq-sa-key-file",
			Usage:       "BigQuery service account key file",
			Destination: &x.saKeyFile,
			EnvVars:     []string{"OVERSEER_BIGQUERY_SA_KEY_FILE"},
		},
	}
}

func (x *BigQuery) Configure(ctx context.Context) (interfaces.BigQuery, error) {
	utils.Logger().Info("Configure BigQuery",
		"projectID", x.projectID,
	)
	if x.projectID == "" {
		return nil, goerr.Wrap(types.ErrInvalidOption, "BigQuery project ID is empty")
	}

	var options []option.ClientOption
	if x.saKeyData != "" {
		options = append(options, option.WithCredentialsJSON([]byte(x.saKeyData)))
	}
	if x.saKeyFile != "" {
		options = append(options, option.WithCredentialsFile(x.saKeyFile))
	}
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

		options = append(options, option.WithTokenSource(ts))
	}

	return bq.New(ctx,
		x.projectID,
		options...,
	)
}

func (x *BigQuery) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("project_id", x.projectID),
		slog.String("impersonate_service_account", x.impersonateServiceAccount),
	)
}
