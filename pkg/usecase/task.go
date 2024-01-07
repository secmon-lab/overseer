package usecase

import (
	"context"

	"github.com/m-mizutani/overseer/pkg/domain/model"
	"github.com/m-mizutani/overseer/pkg/infra"
	"github.com/m-mizutani/overseer/pkg/utils/ctxutil"
)

func RunTask(ctx context.Context, clients *infra.Clients, task *model.Task) error {
	ctx = ctxutil.WithRequestID(ctx)
	ctxutil.Logger(ctx).Info("Start task",
		"title", task.Title,
		"description", task.Description,
	)

	rows, err := clients.BigQuery().Query(ctx, task.Query)
	if err != nil {
		return err
	}

	for i, row := range rows {
		if i >= int(task.Limit) {
			break
		}

		alert := &model.Alert{
			Title:       task.Title,
			Description: task.Description,
		}

		for key, value := range row {
			alert.Attrs = append(alert.Attrs, model.Attribute{
				Key:   key,
				Value: value,
			})
		}

		if err := clients.Emitter().Emit(ctx, alert); err != nil {
			return err
		}
	}

	return nil
}
