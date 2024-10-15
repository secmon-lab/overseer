package usecase

import (
	"context"

	"github.com/secmon-as-code/overseer/pkg/domain/model"
)

func (x *UseCase) RunQueries(ctx context.Context, queries model.Queries, newWriter NewWriter) error {
	return x.runQueries(ctx, queries, newWriter)
}
