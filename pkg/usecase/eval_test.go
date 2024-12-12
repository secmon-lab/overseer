package usecase_test

import (
	"context"
	"embed"
	"io"
	"testing"
	"time"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/opac"
	"github.com/secmon-lab/overseer/pkg/adaptor"
	"github.com/secmon-lab/overseer/pkg/domain/model"
	"github.com/secmon-lab/overseer/pkg/mock"
	"github.com/secmon-lab/overseer/pkg/service"
	"github.com/secmon-lab/overseer/pkg/usecase"
)

//go:embed testdata/eval/*
var evalFiles embed.FS

func TestEval(t *testing.T) {
	policyList := []string{
		"testdata/eval/policy1.rego",
	}

	policies := make(map[string]string)
	for _, path := range policyList {
		policy, err := evalFiles.ReadFile(path)
		gt.NoError(t, err).Must()
		policies[path] = string(policy)
	}

	ctx := context.Background()
	policyClient, err := opac.New(opac.Data(policies))
	gt.NoError(t, err).Must()
	policySvc, err := service.NewPolicy(policyClient, model.SelectPolicyAll)
	gt.NoError(t, err).Must()

	cache := &mock.CacheServiceMock{
		NewReaderFunc: func(ctx context.Context, ID model.QueryID) (io.ReadCloser, error) {
			return evalFiles.Open("testdata/eval/" + string(ID) + ".json")
		},
	}
	notify := &mock.NotifyServiceMock{
		PublishFunc: func(ctx context.Context, alert model.Alert) error {
			gt.EQ(t, alert.Title, "Test Policy 1")
			gt.EQ(t, alert.Description, "Principal attempted to access data")
			gt.True(t, alert.Timestamp.Equal(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)))
			gt.M(t, alert.Attrs).HaveKeyValue("id", float64(3))
			return nil
		},
	}

	uc := usecase.New(adaptor.New())

	gt.NoError(t, uc.Eval(ctx, policySvc, cache, notify))
	gt.A(t, notify.PublishCalls()).Length(1)
}
