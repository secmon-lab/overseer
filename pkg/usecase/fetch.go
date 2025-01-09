package usecase

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/dustin/go-humanize"
	"github.com/m-mizutani/goerr/v2"
	"github.com/secmon-lab/overseer/pkg/domain/interfaces"
	"github.com/secmon-lab/overseer/pkg/domain/model"
	"github.com/secmon-lab/overseer/pkg/logging"
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
	eb := goerr.NewBuilder(goerr.V("query_id", query.ID()))

	startTS := time.Now()

	it, stat, err := bq.Query(ctx, query.String())
	if err != nil {
		return eb.Wrap(err, "fail to query")
	}

	w, err := cache.NewWriter(ctx, query.ID())
	if err != nil {
		return eb.Wrap(err, "fail to create cache writer")
	}
	defer safeClose(ctx, w)

	dataSize := 0
	recordCount := 0
	if n, err := w.Write([]byte("[")); err != nil {
		return eb.Wrap(err, "fail to write header bracket")
	} else {
		dataSize += n
	}

	isFirst := true
	for {
		var row map[string]bigquery.Value

		if err := it.Next(&row); err == iterator.Done {
			break
		} else if err != nil {
			return eb.Wrap(err, "fail to get next row")
		}
		recordCount++

		if !isFirst {
			if _, err := w.Write([]byte(",")); err != nil {
				return eb.Wrap(err, "fail to write separator")
			}
		}
		isFirst = false

		data, err := json.Marshal(row)
		if err != nil {
			return eb.Wrap(err, "fail to marshal row")
		}

		if n, err := w.Write(data); err != nil {
			return err
		} else {
			dataSize += n
		}
	}

	if n, err := w.Write([]byte("]")); err != nil {
		return eb.Wrap(err, "fail to write footer bracket")
	} else {
		dataSize += n
	}

	logger.Info("Finalize query",
		"query", query.ID(),
		"data_size", dataSize,
		"record_count", recordCount,
		"duration", time.Since(startTS),
		"proceeded_bytes", stat.TotalBytesProcessed,
		"proceeded_bytes_formatted", humanize.Bytes(uint64(stat.TotalBytesProcessed)),
	)

	return nil
}
