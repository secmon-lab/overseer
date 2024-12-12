package types_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/overseer/pkg/domain/types"
)

func TestNewJobID(t *testing.T) {
	ctx := context.Background()
	ctx, jobID := types.NewJobID(ctx)

	gt.NE(t, jobID, "")
	gt.S(t, string(jobID)).Contains("job")
}
