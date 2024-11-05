package usecase_test

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"io"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/overseer/pkg/adaptor"
	"github.com/secmon-lab/overseer/pkg/domain/interfaces"
	"github.com/secmon-lab/overseer/pkg/domain/model"
	"github.com/secmon-lab/overseer/pkg/mock"
	"github.com/secmon-lab/overseer/pkg/usecase"
	"google.golang.org/api/iterator"
)

type buffer struct {
	bytes.Buffer
}

func (b *buffer) Close() error {
	return nil
}

type mockIterator struct {
	results []map[string]bigquery.Value
	index   int
}

func (it *mockIterator) Next(row interface{}) error {
	if it.index >= len(it.results) {
		return iterator.Done
	}

	data := it.results[it.index]
	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(rawData, row); err != nil {
		return err
	}

	it.index++
	return nil
}

//go:embed testdata/query1.sql
var query1 []byte

func TestExtract(t *testing.T) {
	mockBQ := mock.BigQueryClientMock{
		QueryFunc: func(ctx context.Context, query string) (interfaces.BigQueryIterator, *bigquery.JobStatistics, error) {
			return &mockIterator{
				results: []map[string]bigquery.Value{
					{"key1": "value1"},
					{"key1": "value2"},
				},
			}, &bigquery.JobStatistics{}, nil
		},
	}

	var buf buffer
	cache := &mock.CacheServiceMock{
		NewWriterFunc: func(ctx context.Context, ID model.QueryID) (io.WriteCloser, error) {
			return &buf, nil
		},
		StringFunc: func() string {
			return "mock"
		},
	}

	uc := usecase.New(adaptor.New(
		adaptor.WithBigQuery(&mockBQ),
	))

	gt.NoError(t, uc.Fetch(context.Background(), model.Queries{
		model.MustNewQuery("x", query1),
	}, cache))

	var result []map[string]string
	gt.NoError(t, json.Unmarshal(buf.Bytes(), &result))
	gt.Equal(t, len(result), 2)
	gt.Equal(t, result[0]["key1"], "value1")
	gt.Equal(t, result[1]["key1"], "value2")
}
