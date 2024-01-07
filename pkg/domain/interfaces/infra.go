package interfaces

import (
	"context"

	"github.com/m-mizutani/overseer/pkg/domain/model"
)

type BigQuery interface {
	Query(ctx context.Context, query string) ([]model.BigQueryRow, error)
}

type Queue interface {
	Publish(ctx context.Context, alert *model.Alert) error
}
