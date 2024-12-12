package types_test

import (
	"context"
	"testing"

	"github.com/secmon-lab/overseer/pkg/domain/types"
)

func TestNewRequestID(t *testing.T) {
	ctx := context.Background()
	newCtx, requestID := types.NewRequestID(ctx)

	if requestID == "" {
		t.Errorf("Expected a non-empty RequestID, got an empty string")
	}

	retrievedID := types.RequestIDFromCtx(newCtx)
	if retrievedID != requestID {
		t.Errorf("Expected RequestID %s, got %s", requestID, retrievedID)
	}
}
