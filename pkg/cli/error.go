package cli

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/m-mizutani/goerr/v2"
	"github.com/secmon-lab/overseer/pkg/domain/types"
	"github.com/secmon-lab/overseer/pkg/logging"
)

func handleError(ctx context.Context, err error) {
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		if goErr := goerr.Unwrap(err); goErr != nil {
			for k, v := range goErr.Values() {
				scope.SetExtra(fmt.Sprintf("%v", k), v)
			}
		}

		scope.SetExtra("job_id", types.JobIDFromCtx(ctx))
	})
	evID := hub.CaptureException(err)

	logger := logging.FromCtx(ctx)
	logger.Error("Error", "error", err, "sentry.event_id", evID)
}
