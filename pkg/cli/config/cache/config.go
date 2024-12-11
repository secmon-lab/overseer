package cache

import (
	"context"
	"log/slog"

	"github.com/secmon-lab/overseer/pkg/adaptor/cs"
	"github.com/secmon-lab/overseer/pkg/domain/interfaces"
	"github.com/secmon-lab/overseer/pkg/domain/model"
	"github.com/secmon-lab/overseer/pkg/service"

	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v3"
)

type Config struct {
	// use file cache
	fsDir string

	// use Cloud Storage cache
	gcsBucket string
	gcsPrefix string
}

func (x *Config) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "cache-dir",
			Usage:       "Directory path to store cache files",
			Category:    "cache",
			Destination: &x.fsDir,
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_FS_DIR")),
		},
		&cli.StringFlag{
			Name:        "cache-gcs-bucket",
			Usage:       "Cloud Storage bucket name",
			Category:    "cache",
			Destination: &x.gcsBucket,
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_GCS_BUCKET")),
		},
		&cli.StringFlag{
			Name:        "cache-gcs-prefix",
			Usage:       "Cloud Storage prefix",
			Category:    "cache",
			Destination: &x.gcsPrefix,
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_GCS_PREFIX")),
		},
	}
}

func (x *Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("cache-dir", x.fsDir),
		slog.String("cache-gcs-bucket", x.gcsBucket),
		slog.String("cache-gcs-prefix", x.gcsPrefix),
	)
}

func (x *Config) Build(ctx context.Context, id model.JobID) (interfaces.CacheService, error) {
	if x.fsDir != "" && x.gcsBucket != "" {
		return nil, goerr.New("cache-dir and cache-bucket are exclusive, specify only one")
	}

	switch {
	case x.fsDir != "":
		return service.NewFileCache(id, x.fsDir)

	case x.gcsBucket != "":
		client, err := cs.NewClient(ctx)
		if err != nil {
			return nil, goerr.Wrap(err).With("gcsBucket", x.gcsBucket)
		}

		return service.NewCloudStorageCache(id, x.gcsBucket, x.gcsPrefix, client), nil

	default:
		return nil, goerr.New("No cache service is specified")
	}
}
