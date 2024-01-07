package interfaces

import (
	"context"

	"github.com/m-mizutani/overseer/pkg/domain/model"
)

type BigQuery interface {
	Query(ctx context.Context, query string) ([]model.BigQueryRow, error)
}

type Emitter interface {
	Emit(ctx context.Context, alert *model.Alert) error
}
