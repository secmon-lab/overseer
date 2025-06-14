package usecase

import (
	"context"
	"encoding/json"
	"runtime"

	"github.com/m-mizutani/goerr/v2"
	"github.com/m-mizutani/opaq"
	"github.com/secmon-lab/overseer/pkg/domain/interfaces"
	"github.com/secmon-lab/overseer/pkg/domain/model"
	"github.com/secmon-lab/overseer/pkg/logging"
	"github.com/secmon-lab/overseer/pkg/service"
)

func (x *UseCase) Eval(ctx context.Context, p *service.Policy, cache interfaces.CacheService, notify interfaces.NotifyService) error {
	for _, meta := range p.MetadataSet() {

		// Before evaluating policy, run GC to release unused memory
		runtime.GC()

		if err := evalPolicy(ctx, p.Client(), meta, cache, notify); err != nil {
			return err
		}
	}

	return nil
}

func evalPolicy(ctx context.Context, policy interfaces.PolicyClient, meta *model.PolicyMetadata, cache interfaces.CacheService, notify interfaces.NotifyService) error {
	input := model.QueryInput{}

	for _, queryID := range meta.Input {
		r, err := cache.NewReader(ctx, queryID)
		if err != nil {
			return err
		}
		defer r.Close()

		var body any
		decoder := json.NewDecoder(r)
		if err := decoder.Decode(&body); err != nil {
			return err
		}

		input[queryID] = body
	}

	logger := logging.FromCtx(ctx).With("package", meta.Package)
	hook := func(ctx context.Context, loc opaq.PrintLocation, msg string) error {
		logger.Info("[Rego] "+msg, "file", loc.File, "line", loc.Row)
		return nil
	}

	options := []opaq.QueryOption{
		opaq.WithPrintHook(hook),
	}

	var output model.QueryOutput
	if err := policy.Query(ctx, "data."+meta.Package, input, &output, options...); err != nil {
		return goerr.Wrap(err, "failed to evaluate policy")
	}

	logging.FromCtx(ctx).Info("Evaluated policy", "policy", meta.Package, "output", output)

	for _, body := range output.Alert {
		alert, err := model.NewAlert(ctx, body)
		if err != nil {
			return err
		}

		if err := notify.Publish(ctx, *alert); err != nil {
			return err
		}
	}

	return nil
}
