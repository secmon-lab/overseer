package interfaces

import (
	"context"
	"io"

	"cloud.google.com/go/bigquery"
	"github.com/m-mizutani/opaq"
	"github.com/open-policy-agent/opa/v1/ast"
)

type CloudStorageClient interface {
	PutObject(ctx context.Context, bucketName, objectName string) (io.WriteCloser, error)
	GetObject(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
}

type BigQueryClient interface {
	Query(ctx context.Context, query string) (BigQueryIterator, *bigquery.JobStatistics, error)
}

type BigQueryIterator interface {
	Next(dst interface{}) error
}

type PubSubClient interface {
	Publish(ctx context.Context, topic string, data []byte) error
}

type PolicyClient interface {
	Query(ctx context.Context, query string, input, output any, options ...opaq.QueryOption) error
	Metadata() ast.FlatAnnotationsRefSet
}
