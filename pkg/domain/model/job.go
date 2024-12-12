package model

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/secmon-lab/overseer/pkg/logging"
)

type JobID string
type ctxJobIDKey struct{}

func NewJobID(ctx context.Context) (context.Context, JobID) {
	id, err := uuid.NewV7()
	if err != nil {
		logging.Default().Error("fail to generate new JobID", "err", err)
		panic(err)
	}

	now := time.Now()
	jobID := JobID(now.Format("job200601021504_") + strings.ReplaceAll(id.String(), "-", ""))
	return context.WithValue(ctx, ctxJobIDKey{}, jobID), jobID
}

func JobIDFromCtx(ctx context.Context) JobID {
	if id, ok := ctx.Value(ctxJobIDKey{}).(JobID); ok {
		return id
	}
	return ""
}
