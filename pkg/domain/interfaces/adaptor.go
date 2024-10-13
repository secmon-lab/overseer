package interfaces

import (
	"context"
	"io"
)

type CloudStorageClient interface {
	PutObject(ctx context.Context, bucketName, objectName string) (io.WriteCloser, error)
	GetObject(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
}

type BigQueryClient interface {
	Query(ctx context.Context, query string) (BigQueryIterator, error)
}

type BigQueryIterator interface {
	Next(dst interface{}) error
}
