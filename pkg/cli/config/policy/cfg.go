package policy

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/opac"
	"github.com/secmon-lab/overseer/pkg/domain/model"
	"github.com/secmon-lab/overseer/pkg/service"
	"github.com/urfave/cli/v3"
)

type Config struct {
	filePath   []string
	selectTags []string
}

func (x *Config) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:        "policy",
			Usage:       "Policy file/directory",
			Category:    "policy",
			Destination: &x.filePath,
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_POLICY")),
			Required:    true,
			Aliases:     []string{"p"},
		},
		&cli.StringSliceFlag{
			Name:        "tag",
			Usage:       "Target policy tag. If not specified, all policy is target",
			Category:    "policy",
			Destination: &x.selectTags,
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_POLICY_TAG")),
			Aliases:     []string{"t"},
		},
	}
}

func (x *Config) FilePath() []string {
	return x.filePath[:]
}

func (x *Config) Build() (*service.Policy, error) {
	client, err := opac.New(opac.Files(x.filePath...))
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create policy client")
	}

	var selector model.PolicySelector
	switch {
	case len(x.selectTags) > 0:
		selector = model.SelectPolicyByTag(x.selectTags...)
	default:
		selector = model.SelectPolicyAll
	}

	return service.NewPolicy(client, selector)
}
