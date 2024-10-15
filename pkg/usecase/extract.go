package usecase

import (
	"context"
	"encoding/json"
	"io"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-as-code/overseer/pkg/domain/model"
	"github.com/secmon-as-code/overseer/pkg/interfaces"
)

type NewWriter func(ID model.QueryID) (io.WriteCloser, error)

func (x *UseCase) Extract(ctx context.Context, queries model.Queries) error {
	if err := queries.Validate(); err != nil {
		return err
	}

	return nil
}

func (x *UseCase) runQueries(ctx context.Context, queries model.Queries, newWriter NewWriter) error {
	if err := queries.Validate(); err != nil {
		return err
	}

	for _, query := range queries {
		if err := queryAndDump(ctx, x.clients.BigQuery(), query, newWriter); err != nil {
			return goerr.Wrap(err, "fail to extract query")
		}
	}

	return nil
}

func queryAndDump(ctx context.Context, bq interfaces.BigQueryClient, query *model.Query, newWriter NewWriter) error {
	it, err := bq.Query(ctx, query.String())
	if err != nil {
		return err
	}

	w, err := newWriter(query.ID())
	if err != nil {
		return err
	}
	defer safeClose(ctx, w)

	if _, err := w.Write([]byte("[")); err != nil {
		return goerr.Wrap(err, "fail to write header bracket")
	}

	isFirst := true
	for {
		var row any
		if err := it.Next(&row); err == io.EOF {
			break
		} else if err != nil {
			return goerr.Wrap(err, "fail to get next row")
		}

		if !isFirst {
			if _, err := w.Write([]byte(",")); err != nil {
				return goerr.Wrap(err, "fail to write separator")
			}
		}
		isFirst = false

		data, err := json.Marshal(row)
		if err != nil {
			return goerr.Wrap(err, "fail to marshal row")
		}

		if _, err := w.Write(data); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte("]")); err != nil {
		return goerr.Wrap(err, "fail to write footer bracket")
	}

	return nil
}
