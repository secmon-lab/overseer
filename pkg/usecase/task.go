package usecase

import (
	"context"
	"sync"

	"github.com/hashicorp/go-multierror"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/overseer/pkg/domain/model"
	"github.com/m-mizutani/overseer/pkg/infra"
	"github.com/m-mizutani/overseer/pkg/utils"
)

const (
	concurrency = 10
)

func RunTasks(ctx context.Context, clients *infra.Clients, tasks []*model.Task, tgt *model.Target) error {
	ctx, _ = utils.CtxRequestID(ctx)

	errCh := make(chan error, len(tasks))
	taskCh := make(chan *model.Task, len(tasks))

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if err := runTask(ctx, clients, task); err != nil {
					err = goerr.Wrap(err).With("task", task)
					utils.HandleError(ctx, "fail to run task", err)
					errCh <- err
				}
			}
		}()
	}

	for _, task := range tasks {
		if !tgt.Contains(task) {
			continue
		}

		taskCh <- task
	}
	close(taskCh)
	wg.Wait()
	close(errCh)

	var resultErr error
	for err := range errCh {
		resultErr = multierror.Append(resultErr, err)
	}

	return resultErr
}

func runTask(ctx context.Context, clients *infra.Clients, task *model.Task) error {
	utils.CtxLogger(ctx).Info("Start task",
		"title", task.Title,
		"description", task.Description,
	)

	rows, err := clients.BigQuery().Query(ctx, task.Query)
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		return nil
	}

	alert := &model.Alert{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
	}

	for _, row := range rows {
		result := make(map[string]any)
		for key, value := range row {
			result[key] = value
		}
		alert.Results = append(alert.Results, result)
	}

	utils.CtxLogger(ctx).Info("detected alert", "alert", alert)

	if err := clients.Emitter().Emit(ctx, alert); err != nil {
		return err
	}

	return nil

}
