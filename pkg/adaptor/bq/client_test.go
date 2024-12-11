package bq_test

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/overseer/pkg/adaptor/bq"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func TestImpersonation(t *testing.T) {
	bqProjectID := os.Getenv("TEST_BQ_PROJECT_ID")
	bqQuery := os.Getenv("TEST_BQ_QUERY")
	bqImpersonateServiceAccount := os.Getenv("TEST_BQ_IMPERSONATE_SERVICE_ACCOUNT")

	if bqProjectID == "" || bqQuery == "" || bqImpersonateServiceAccount == "" {
		t.Skip("TEST_BQ_PROJECT_ID and TEST_BQ_QUERY are required")
	}

	ctx := context.Background()

	ts, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
		TargetPrincipal: bqImpersonateServiceAccount,
		Scopes: []string{
			"https://www.googleapis.com/auth/bigquery",
			"https://www.googleapis.com/auth/cloud-platform",
			"https://www.googleapis.com/auth/bigquery.readonly",
			"https://www.googleapis.com/auth/cloud-platform.read-only",
		},
	})
	gt.NoError(t, err)

	bqClient, err := bq.New(ctx, bqProjectID, option.WithTokenSource(ts))
	gt.NoError(t, err)

	it, _, err := bqClient.Query(ctx, bqQuery)
	gt.NoError(t, err)
	gt.NotEqual(t, it, nil)

	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err != nil {
			if err == iterator.Done {
				break
			}
			gt.NoError(t, err).Must()
		}
		t.Log(row)
	}
}
