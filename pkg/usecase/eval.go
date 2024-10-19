package usecase

import (
	"context"

	"github.com/secmon-as-code/overseer/pkg/domain/interfaces"
	"github.com/secmon-as-code/overseer/pkg/domain/model"
)

func (x *UseCase) Eval(ctx context.Context, id model.JobID, cache interfaces.CacheService) error {
	return nil
}
