package query

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/secmon-lab/overseer/pkg/domain/model"

	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v3"
)

type Config struct {
	filePath []string
}

func New() *Config {
	return &Config{}
}

func (x *Config) FilePath() []string {
	return x.filePath[:]
}

func (x *Config) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:        "query",
			Usage:       "Query file/directory",
			Category:    "query",
			Destination: &x.filePath,
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_QUERY")),
			Required:    true,
			Aliases:     []string{"q"},
		},
	}
}

func (x *Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("query", x.filePath),
	)
}

func (x *Config) Build() (model.Queries, error) {
	var queries model.Queries

	for _, path := range x.filePath {
		q, err := loadQueries(path)
		if err != nil {
			return nil, err
		}
		queries = append(queries, q...)
	}

	return queries, nil
}

func loadQueries(target string) (model.Queries, error) {
	var queries model.Queries

	if err := filepath.WalkDir(target, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		body, err := os.ReadFile(path)
		if err != nil {
			return goerr.Wrap(err, "fail to read query file").With("path", path)
		}

		name := strings.Trim(filepath.Base(path), filepath.Ext(path))

		q, err := model.NewQuery(name, body)
		if err != nil {
			return err
		}
		queries = append(queries, q)

		return nil
	}); err != nil {
		return nil, err
	}

	return queries, nil
}
