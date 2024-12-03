package usecase

import (
	"context"
	"encoding/json"
	"log/slog"
	"runtime"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/opac"
	"github.com/open-policy-agent/opa/topdown/print"
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

type printHook struct {
	logger *slog.Logger
}

func (x *printHook) Print(ctx print.Context, msg string) error {
	x.logger.Info("[Rego] "+msg, "file", ctx.Location.File, "line", ctx.Location.Row)
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

	hook := &printHook{
		logger: logging.FromCtx(ctx).With("package", meta.Package),
	}

	options := []opac.QueryOption{
		opac.WithPrintHook(hook),
	}

	var output model.QueryOutput
	if err := policy.Query(ctx, "data."+meta.Package, input, &output, options...); err != nil {
		return goerr.Wrap(err, "failed to evaluate policy")
	}

	logging.FromCtx(ctx).Info("Evaluated policy", "policy", meta.Package, "output", output)

	for _, alert := range output.Alert {
		if err := alert.Validate(); err != nil {
			return goerr.Wrap(err, "validate evaluated alert").With("policy", meta.Package)
		}
		alert.Finalize(ctx)

		if err := notify.Publish(ctx, alert); err != nil {
			return err
		}
	}

	return nil
}
