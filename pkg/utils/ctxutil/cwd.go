package ctxutil

import (
	"context"
	"os"
)

var (
	cwdKey = contextKey{}
)

func WithCWD(ctx context.Context, dir string) context.Context {
	return context.WithValue(ctx, cwdKey, dir)
}

func CWD(ctx context.Context) string {
	if v, ok := ctx.Value(cwdKey).(string); ok {
		return v
	}

	if v, err := os.Getwd(); err == nil {
		return v
	}

	return ""
}
