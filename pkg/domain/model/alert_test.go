package model_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/overseer/pkg/domain/model"
)

func TestAlert(t *testing.T) {
	a := model.Alert{}

	a.Finalize(context.TODO())
	gt.NotEqual(t, a.ID, "")
	gt.False(t, a.Timestamp.IsZero())
}
