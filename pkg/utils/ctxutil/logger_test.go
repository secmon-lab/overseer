package ctxutil_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/overseer/pkg/utils/ctxutil"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	ctx := context.Background()

	t.Run("default logger", func(t *testing.T) {
		ctxutil.Logger(ctx).Info("hello")
		gt.Equal(t, buf.String(), "")
	})

	t.Run("with logger", func(t *testing.T) {
		ctx = ctxutil.WithLogger(ctx, logger)
		ctxutil.Logger(ctx).Info("hello")
		gt.S(t, buf.String()).Contains("hello")
	})
}
