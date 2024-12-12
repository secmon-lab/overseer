package types

import (
	"context"

	"github.com/google/uuid"
)

type RequestID string

type ctxRequestIDKeyType struct{}

func NewRequestID(ctx context.Context) (context.Context, RequestID) {
	requestID := RequestID(uuid.NewString())

	return context.WithValue(ctx, ctxRequestIDKeyType{}, requestID), requestID
}

func RequestIDFromCtx(ctx context.Context) RequestID {
	if id, ok := ctx.Value(ctxRequestIDKeyType{}).(RequestID); ok {
		return id
	}
	return ""
}
