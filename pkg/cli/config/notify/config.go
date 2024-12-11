package notify

import (
	"io"
	"log/slog"
	"os"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/overseer/pkg/adaptor/pubsub"
	"github.com/secmon-lab/overseer/pkg/domain/interfaces"
	"github.com/secmon-lab/overseer/pkg/service"
	"github.com/urfave/cli/v3"
)

type Config struct {
	pubsubProject string
	pubsubTopic   string
	output        string
}

func (x *Config) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "notify-pubsub-topic",
			Usage:       "Pub/Sub topic name for alert notification",
			Category:    "notify",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_NOTIFY_PUBSUB_TOPIC")),
			Destination: &x.pubsubTopic,
		},
		&cli.StringFlag{
			Name:        "notify-pubsub-project",
			Usage:       "Pub/Sub project ID for alert notification",
			Category:    "notify",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_NOTIFY_PUBSUB_PROJECT")),
			Destination: &x.pubsubProject,
		},
		&cli.StringFlag{
			Name:        "notify-out",
			Usage:       "Output destination ('-', stdout, stderr). Default is stdout",
			Value:       "stdout",
			Category:    "notify",
			Aliases:     []string{"o"},
			Sources:     cli.NewValueSourceChain(cli.EnvVar("OVERSEER_NOTIFY_OUT")),
			Destination: &x.output,
		},
	}
}

func (x *Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("notify-pubsub-topic", x.pubsubTopic),
		slog.String("notify-pubsub-project", x.pubsubProject),
		slog.String("notify-out", x.output),
	)
}

func (x *Config) Build() (interfaces.NotifyService, error) {
	if x.pubsubTopic != "" && x.pubsubProject != "" {
		client, err := pubsub.New(x.pubsubProject)
		if err != nil {
			return nil, err
		}

		return service.NewNotifyPubSub(client, x.pubsubTopic), nil
	}

	var w io.Writer
	switch x.output {
	case "stdout", "-":
		w = os.Stdout
	case "stderr":
		w = os.Stderr
	default:
		return nil, goerr.New("Invalid output destination: %s", x.output)
	}

	return service.NewNotifyWriter(w), nil
}
