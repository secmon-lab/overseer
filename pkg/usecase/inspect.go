package usecase

import (
	"context"
	"encoding/json"
	"io"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-as-code/overseer/pkg/domain/model"
	"github.com/secmon-as-code/overseer/pkg/logging"
)

func (x *UseCase) Inspect(ctx context.Context, queries model.Queries, w io.Writer) error {
	logger := logging.FromCtx(ctx)

	if err := queries.Validate(); err != nil {
		return err
	}

	queryIDs := map[model.QueryID]int{}
	for _, query := range queries {
		queryIDs[query.ID()] = 0
	}

	var metadataSet []*model.PolicyMetadata
	for _, ref := range x.clients.Policy().Metadata() {
		meta, err := model.NewPolicyMetadataFromAnnotation(ref)
		if err != nil {
			return err
		}

		metadataSet = append(metadataSet, meta)

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

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	result := map[string]any{
		"metadata": metadataSet,
		"queries":  queryIDs,
	}

	if err := encoder.Encode(result); err != nil {
		return goerr.Wrap(err, "fail to encode metadata set")
	}

	return nil
}
