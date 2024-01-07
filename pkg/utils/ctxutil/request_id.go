package ctxutil

import (
	"context"

	"github.com/m-mizutani/overseer/pkg/domain/types"
)

var (
	requestIDKey = contextKey{}
)

// WithRequestID returns context with request ID. If request ID is already set, it returns same context. If not, it generates new request ID and set it and logger to context.
func WithRequestID(ctx context.Context) context.Context {
	if _, ok := ctx.Value(requestIDKey).(string); ok {
		return ctx
	}

	id := types.NewRequestID()
	ctx = context.WithValue(ctx, requestIDKey, id)
	ctx = WithLogger(ctx, Logger(ctx).With("requestID", id))
	return ctx
}
