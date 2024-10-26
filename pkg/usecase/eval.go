package usecase

import (
	"context"
	"encoding/json"

	"github.com/secmon-lab/overseer/pkg/domain/interfaces"
	"github.com/secmon-lab/overseer/pkg/domain/model"
	"github.com/secmon-lab/overseer/pkg/logging"
	"github.com/secmon-lab/overseer/pkg/service"
)

func (x *UseCase) Eval(ctx context.Context, p *service.Policy, cache interfaces.CacheService, notify interfaces.NotifyService) error {
	for _, meta := range p.MetadataSet() {
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

	var output model.QueryOutput
	if err := policy.Query(ctx, "data."+meta.Package, input, &output); err != nil {
		return err
	}

	logging.FromCtx(ctx).Info("Evaluated policy", "policy", meta.Package, "output", output)

	for _, alert := range output.Alert {
		if err := notify.Publish(ctx, alert); err != nil {
			return err
		}
	}

	return nil
}
