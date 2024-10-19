package policy

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/opac"
	"github.com/urfave/cli/v3"
)

type Config struct {
	filePath []string
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
	}
}

func (x *Config) FilePath() []string {
	return x.filePath[:]
}

func (x *Config) Build() (*opac.Client, error) {
	client, err := opac.New(opac.Files(x.filePath...))
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create policy client")
	}

	return client, nil
}
