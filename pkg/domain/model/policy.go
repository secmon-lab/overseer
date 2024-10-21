package model

import (
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/open-policy-agent/opa/ast"
)

type PolicyMetadataSet []*PolicyMetadata

func NewPolicyMetadataSetFromAnnotation(refs ast.FlatAnnotationsRefSet) (PolicyMetadataSet, error) {
	var metadataSet PolicyMetadataSet
	for _, ref := range refs {
		meta, err := NewPolicyMetadataFromAnnotation(ref)
		if err != nil {
			return nil, err
		}
		metadataSet = append(metadataSet, meta)
	}
	return metadataSet, nil
}

func (x PolicyMetadataSet) Filter(selector PolicySelector) PolicyMetadataSet {
	var filtered PolicyMetadataSet
	for _, meta := range x {
		if selector(meta) {
			filtered = append(filtered, meta)
		}
	}
	return filtered
}

func (x PolicyMetadataSet) RequiredQueries(base Queries) Queries {
	queryIDs := make(map[QueryID]struct{})
	for _, meta := range x {
		for _, id := range meta.Input {
			queryIDs[id] = struct{}{}
		}
	}

	var newQueries Queries
	for id := range queryIDs {
		newQueries = append(newQueries, base.FindByID(id))
	}

	return newQueries
}

type PolicyMetadata struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Input       []QueryID `json:"input"`
	Package     string    `json:"package"`
	Location    string    `json:"location"`
}

func (x *PolicyMetadata) HasTag(tag string) bool {
	for _, t := range x.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

type PolicySelector func(meta *PolicyMetadata) bool

func SelectPolicyByTag(tags ...string) PolicySelector {
	return func(meta *PolicyMetadata) bool {
		for _, tag := range tags {
			if meta.HasTag(tag) {
				return true
			}
		}
		return false
	}
}

func SelectPolicyAll(meta *PolicyMetadata) bool {
	return true
}

func NewPolicyMetadataFromAnnotation(ref *ast.AnnotationsRef) (*PolicyMetadata, error) {
	if ref == nil {
		return nil, goerr.New("AnnotationsRef is nil")
	}
	if ref.Annotations == nil {
		return nil, goerr.New("Annotations is nil")
	}

	eb := goerr.NewBuilder().With("ref", ref.Annotations)

	if ref.Annotations.Scope != "package" {
		return nil, eb.New("Annotations.Scope is not 'package'")
	}

	meta := &PolicyMetadata{
		Title:       ref.Annotations.Title,
		Description: ref.Annotations.Description,
		Location:    ref.Location.String(),
	}

	tags, ok := ref.Annotations.Custom["tags"].([]any)
	if !ok {
		return nil, eb.New("custom.tags is not found or invalid format")
	}
	if len(tags) == 0 {
		return nil, eb.New("custom.tags is empty")
	}
	for _, tag := range tags {
		v, ok := tag.(string)
		if !ok {
			return nil, eb.New("custom.tags contains non-string element")
		}
		if v == "" {
			return nil, eb.New("custom.tags contains empty element")
		}

		meta.Tags = append(meta.Tags, v)
	}

	input, ok := ref.Annotations.Custom["input"].([]any)
	if !ok {
		return nil, eb.New("custom.input is not found or invalid format")
	}
	if len(input) == 0 {
		return nil, eb.New("custom.input is empty")
	}
	for _, id := range input {
		v, ok := id.(string)
		if !ok {
			return nil, eb.New("custom.input contains non-string element")
		}

		meta.Input = append(meta.Input, QueryID(v))
	}

	// Extract path of metadata
	var path []string
	for _, p := range ref.Path {
		path = append(path, strings.Trim(p.Value.String(), `"`))
	}
	meta.Package = strings.Join(path[1:], ".")

	return meta, nil
}
