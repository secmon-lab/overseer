package model

import (
	"context"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/open-policy-agent/opa/ast"
)

type PolicyMetadata struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Input       []QueryID `json:"input"`
	Package     string    `json:"package"`
	Location    string    `json:"location"`
}

func NewPolicyMetadataFromAnnotation(ref *ast.AnnotationsRef) (*PolicyMetadata, error) {
	ctx := context.Background()
	ctx = goerr.InjectValue(ctx, "ref", ref.Annotations)

	if ref == nil {
		return nil, goerr.New("AnnotationsRef is nil").WithContext(ctx)
	}
	if ref.Annotations == nil {
		return nil, goerr.New("Annotations is nil").WithContext(ctx)
	}
	if ref.Annotations.Scope != "package" {
		return nil, goerr.New("Annotations.Scope is not 'package'").WithContext(ctx)
	}

	meta := &PolicyMetadata{
		Title:       ref.Annotations.Title,
		Description: ref.Annotations.Description,
		Location:    ref.Location.String(),
	}

	tags, ok := ref.Annotations.Custom["tags"].([]any)
	if !ok {
		return nil, goerr.New("custom.tags is not found or invalid format").WithContext(ctx)
	}
	if len(tags) == 0 {
		return nil, goerr.New("custom.tags is empty").WithContext(ctx)
	}
	for _, tag := range tags {
		v, ok := tag.(string)
		if !ok {
			return nil, goerr.New("custom.tags contains non-string element").WithContext(ctx)
		}
		if v == "" {
			return nil, goerr.New("custom.tags contains empty element").WithContext(ctx)
		}

		meta.Tags = append(meta.Tags, v)
	}

	input, ok := ref.Annotations.Custom["input"].([]any)
	if !ok {
		return nil, goerr.New("custom.input is not found or invalid format").WithContext(ctx)
	}
	if len(input) == 0 {
		return nil, goerr.New("custom.input is empty").WithContext(ctx)
	}
	for _, id := range input {
		v, ok := id.(string)
		if !ok {
			return nil, goerr.New("custom.input contains non-string element").WithContext(ctx)
		}

		meta.Input = append(meta.Input, QueryID(v))
	}

	// Extract path of metadata
	var path []string
	for _, p := range ref.Path {
		path = append(path, p.Value.String())
	}
	meta.Package = strings.Join(path[1:], ".")

	return meta, nil
}
