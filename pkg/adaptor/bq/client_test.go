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
	ctx := context.Background()

	bqImpersonateServiceAccount := os.Getenv("TEST_BQ_IMPERSONATE_SERVICE_ACCOUNT")
	var bqOptions []option.ClientOption
	if bqImpersonateServiceAccount != "" {
		t.Log("Impersonate service account:", bqImpersonateServiceAccount)
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
		bqOptions = append(bqOptions, option.WithTokenSource(ts))
	}

	bqProjectID := os.Getenv("TEST_BQ_PROJECT_ID")
	bqClient, err := bq.New(ctx, bqProjectID, bqOptions...)
	gt.NoError(t, err)

	bqQuery := os.Getenv("TEST_BQ_QUERY")
	println(bqQuery)
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
