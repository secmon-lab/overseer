package bq_test

import (
	"context"
	_ "embed"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/overseer/pkg/domain/model"
	"github.com/m-mizutani/overseer/pkg/infra/bq"
	"github.com/m-mizutani/overseer/pkg/utils"
)

//go:embed testdata/test.sql
var testQuery string

func TestBigQueryIntegration(t *testing.T) {
	var (
		projectID string
	)

	if err := utils.LoadEnv(
		utils.Env("TEST_BQ_PROJECT_ID", &projectID),
	); err != nil {
		t.Skipf("Skip integration test: %v", err)
	}

	ctx := context.Background()
	client := gt.R1(bq.New(ctx, projectID)).NoError(t)
	data := gt.R1(client.Query(ctx, testQuery)).NoError(t)
	gt.A(t, data).Length(1).At(0, func(t testing.TB, v model.BigQueryRow) {
		gt.Equal(t, v["Name"], "mizutani")
		gt.Equal(t, v["Age"].(int64), int64(13))
	})
}
