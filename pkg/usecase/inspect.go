package usecase

import (
	"context"
	"encoding/json"
	"io"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/overseer/pkg/domain/model"
	"github.com/secmon-lab/overseer/pkg/logging"
	"github.com/secmon-lab/overseer/pkg/service"
)

func (x *UseCase) Inspect(ctx context.Context, queries model.Queries, policy *service.Policy, w io.Writer) error {
	logger := logging.FromCtx(ctx)

	if err := queries.Validate(); err != nil {
		return err
	}

	queryIDs := map[model.QueryID]int{}
	for _, query := range queries {
		queryIDs[query.ID()] = 0
	}

	for _, meta := range policy.MetadataSet() {
		for _, queryID := range meta.Input {
			if _, ok := queryIDs[queryID]; ok {
				queryIDs[queryID]++
			} else {
				return goerr.New("unknown query ID in policy metadata").With("queryID", queryID).With("policy", meta)
			}
		}
	}

	for queryID, count := range queryIDs {
		if count == 0 {
			logger.Warn("Unused query ID", "queryID", queryID)
		}
	}

	if w == nil {
		return nil
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	result := map[string]any{
		"policy": policy.MetadataSet(),
		"query":  queryIDs,
	}

	if err := encoder.Encode(result); err != nil {
		return goerr.Wrap(err, "fail to encode metadata set")
	}

	return nil
}
