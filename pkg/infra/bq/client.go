package bq

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/overseer/pkg/domain/model"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Client struct {
	bq *bigquery.Client
}

func New(ctx context.Context, projectID string, options ...option.ClientOption) (*Client, error) {
	bqClient, err := bigquery.NewClient(ctx, projectID, options...)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to create BigQuery client")
	}

	return &Client{
		bq: bqClient,
	}, nil
}

func (x *Client) Query(ctx context.Context, query string) ([]model.BigQueryRow, error) {
	it, err := x.bq.Query(query).Read(ctx)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to query")
	}

	var itRows []model.BigQueryRow
	for {
		var row map[string]bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, goerr.Wrap(err, "Fail to get next row")
		}

		itRows = append(itRows, row)
	}

	return itRows, nil
}
