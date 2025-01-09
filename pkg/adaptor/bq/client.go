package bq

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/m-mizutani/goerr/v2"
	"github.com/secmon-lab/overseer/pkg/domain/interfaces"
	"google.golang.org/api/option"
)

type Client struct {
	bqClient *bigquery.Client
}

// Query implements interfaces.BigQueryClient.
func (c *Client) Query(ctx context.Context, query string) (interfaces.BigQueryIterator, *bigquery.JobStatistics, error) {
	q := c.bqClient.Query(query)

	job, err := q.Run(ctx)
	if err != nil {
		return nil, nil, goerr.Wrap(err, "fail to run query")
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return nil, nil, goerr.Wrap(err, "fail to wait query job")
	}
	if err := status.Err(); err != nil {
		return nil, nil, goerr.Wrap(err, "query job failed")
	}

	it, err := job.Read(ctx)
	if err != nil {
		return nil, nil, goerr.Wrap(err, "fail to read query result")
	}

	return it, status.Statistics, nil
}

func New(ctx context.Context, projectID string, opts ...option.ClientOption) (*Client, error) {
	bqClient, err := bigquery.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create BigQuery client: %w", err)
	}
	return &Client{bqClient: bqClient}, nil
}

var _ interfaces.BigQueryClient = (*Client)(nil)
