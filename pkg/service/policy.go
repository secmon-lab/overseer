package service

import (
	"github.com/secmon-lab/overseer/pkg/domain/interfaces"
	"github.com/secmon-lab/overseer/pkg/domain/model"
)

type Policy struct {
	metadataSet model.PolicyMetadataSet
	client      interfaces.PolicyClient
	selector    model.PolicySelector
}

func NewPolicy(client interfaces.PolicyClient, selector model.PolicySelector) (*Policy, error) {
	meta, err := model.NewPolicyMetadataSetFromAnnotation(client.Metadata())
	if err != nil {
		return nil, err
	}

	return &Policy{
		client:      client,
		selector:    selector,
		metadataSet: meta,
	}, nil
}

func (x *Policy) Client() interfaces.PolicyClient {
	return x.client
}

func (x *Policy) MetadataSet() model.PolicyMetadataSet {
	var meta model.PolicyMetadataSet
	for _, m := range x.metadataSet {
		if x.selector(m) {
			meta = append(meta, m)
		}
	}

	return meta
}

func (x *Policy) SelectRequiredQueries(base model.Queries) model.Queries {
	if x.selector == nil {
		return base
	}

	baseQueries := map[model.QueryID]*model.Query{}
	for _, query := range base {
		baseQueries[query.ID()] = query
	}

	queries := map[model.QueryID]*model.Query{}
	for _, meta := range x.metadataSet {
		if x.selector(meta) {
			for _, queryID := range meta.Input {
				queries[queryID] = baseQueries[queryID]
			}
		}
	}

	var results model.Queries
	for _, query := range queries {
		results = append(results, query)
	}

	return results
}
