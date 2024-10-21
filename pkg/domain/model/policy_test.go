package model_test

import (
	_ "embed"
	"testing"

	"github.com/secmon-as-code/overseer/pkg/domain/model"

	"github.com/m-mizutani/gt"
	"github.com/open-policy-agent/opa/ast"
)

//go:embed testdata/policy/valid.rego
var validPolicy string

func TestPolicyValidate(t *testing.T) {
	data := map[string]string{
		"policy.rego": validPolicy,
	}

	compiler, err := ast.CompileModulesWithOpt(data, ast.CompileOpts{
		EnablePrintStatements: true,
		ParserOptions: ast.ParserOptions{
			ProcessAnnotation: true,
		},
	})
	gt.NoError(t, err)

	as := compiler.GetAnnotationSet().Flatten()
	gt.A(t, as).Length(1)

	meta, err := model.NewPolicyMetadataFromAnnotation(as[0])
	gt.NoError(t, err)
	gt.Equal(t, meta.Tags, []string{"red"})
	gt.Equal(t, meta.Input, []model.QueryID{"stone"})
}
