package ctxutil

import (
	"context"
	"log/slog"

	"github.com/m-mizutani/overseer/pkg/utils"
)

var (
	loggerKey = contextKey{}
)

func Logger(ctx context.Context) *slog.Logger {
	if v, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return v
	}
	return utils.Logger()
}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
