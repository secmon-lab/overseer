package types_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/overseer/pkg/domain/types"
)

func TestNewJobID(t *testing.T) {
	ctx, jobID := types.NewJobID(context.Background())

	gt.NE(t, jobID, "")
	gt.S(t, string(jobID)).Contains("job")

	retrievedID := types.JobIDFromCtx(ctx)
	gt.EQ(t, jobID, retrievedID)
}
