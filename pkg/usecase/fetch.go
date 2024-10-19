package usecase

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/m-mizutani/goerr"
	"github.com/secmon-as-code/overseer/pkg/domain/interfaces"
	"github.com/secmon-as-code/overseer/pkg/domain/model"
	"github.com/secmon-as-code/overseer/pkg/logging"
	"google.golang.org/api/iterator"
)

type NewWriter func(ID model.QueryID) (io.WriteCloser, error)

func (x *UseCase) Fetch(ctx context.Context, queries model.Queries, cache interfaces.CacheService) error {
	logging.FromCtx(ctx).Info("Start fetching queries",
		"query_count", len(queries),
		"cache", cache.String(),
	)

	if err := queries.Validate(); err != nil {
		return err
	}

	for _, query := range queries {
		if err := queryAndDump(ctx, x.clients.BigQuery(), query, cache); err != nil {
			return goerr.Wrap(err, "fail to extract query")
		}
	}

	return nil
}

func queryAndDump(ctx context.Context, bq interfaces.BigQueryClient, query *model.Query, cache interfaces.CacheService) error {
	logger := logging.FromCtx(ctx)
	logger.Debug("Start fetching queries", "query", query.ID())
	ctx = goerr.InjectValue(ctx, "query_id", query.ID())

	startTS := time.Now()

	it, err := bq.Query(ctx, query.String())
	if err != nil {
		return goerr.Wrap(err).WithContext(ctx)
	}

	w, err := cache.NewWriter(ctx, query.ID())
	if err != nil {
		return goerr.Wrap(err).WithContext(ctx)
	}
	defer safeClose(ctx, w)

	dataSize := 0
	recordCount := 0
	if n, err := w.Write([]byte("[")); err != nil {
		return goerr.Wrap(err, "fail to write header bracket").WithContext(ctx)
	} else {
		dataSize += n
	}

	isFirst := true
	for {
		var row map[string]bigquery.Value

		if err := it.Next(&row); err == iterator.Done {
			break
		} else if err != nil {
			return goerr.Wrap(err, "fail to get next row").WithContext(ctx)
		}
		recordCount++

		if !isFirst {
			if _, err := w.Write([]byte(",")); err != nil {
				return goerr.Wrap(err, "fail to write separator").WithContext(ctx)
			}
		}
		isFirst = false

		data, err := json.Marshal(row)
		if err != nil {
			return goerr.Wrap(err, "fail to marshal row").WithContext(ctx)
		}

		if n, err := w.Write(data); err != nil {
			return err
		} else {
			dataSize += n
		}
	}

	if n, err := w.Write([]byte("]")); err != nil {
		return goerr.Wrap(err, "fail to write footer bracket").WithContext(ctx)
	} else {
		dataSize += n
	}

	logger.Info("complete query",
		"query", query.ID(),
		"data_size", dataSize,
		"record_count", recordCount,
		"duration", time.Since(startTS),
	)

	return nil
}
