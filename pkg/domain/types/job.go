package types

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type JobID string
type ctxJobIDKeyType struct{}

func NewJobID(ctx context.Context) (context.Context, JobID) {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	now := time.Now()
	jobID := JobID(now.Format("job200601021504_") + strings.ReplaceAll(id.String(), "-", ""))
	return context.WithValue(ctx, ctxJobIDKeyType{}, jobID), jobID
}

func JobIDFromCtx(ctx context.Context) JobID {
	if id, ok := ctx.Value(ctxJobIDKeyType{}).(JobID); ok {
		return id
	}
	return ""
}

func InjectJobID(ctx context.Context, id JobID) context.Context {
	return context.WithValue(ctx, ctxJobIDKeyType{}, id)
}
