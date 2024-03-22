package utils

import (
	"context"
	"log/slog"
	"os"

	"github.com/m-mizutani/overseer/pkg/domain/types"
)

var ctxCwdKey struct{}

func CtxWithCWD(ctx context.Context, dir string) context.Context {
	return context.WithValue(ctx, ctxCwdKey, dir)
}

func CtxCWD(ctx context.Context) string {
	if v, ok := ctx.Value(ctxCwdKey).(string); ok {
		return v
	}

	if v, err := os.Getwd(); err == nil {
		return v
	}

	return ""
}

var ctxLoggerKey struct{}

func CtxLogger(ctx context.Context) *slog.Logger {
	if v, ok := ctx.Value(ctxLoggerKey).(*slog.Logger); ok {
		return v
	}
	return Logger()
}

func CtxWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey, logger)
}

var ctxRequestIDKey struct{}

// WithRequestID returns context with request ID. If request ID is already set, it returns same context. If not, it generates new request ID and set it and logger to context.
func CtxRequestID(ctx context.Context) (context.Context, types.RequestID) {
	if id, ok := ctx.Value(ctxRequestIDKey).(types.RequestID); ok {
		return ctx, id
	}

	id := types.NewRequestID()
	ctx = context.WithValue(ctx, ctxRequestIDKey, id)
	ctx = CtxWithLogger(ctx, CtxLogger(ctx).With("request_id", id))
	return ctx, id
}
