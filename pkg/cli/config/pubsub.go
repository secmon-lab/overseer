package config

import (
	"context"

	"github.com/m-mizutani/overseer/pkg/domain/interfaces"
	"github.com/m-mizutani/overseer/pkg/infra/queue"
	"github.com/m-mizutani/overseer/pkg/utils"
	"github.com/urfave/cli/v2"
)

type PubSub struct {
	ProjectID string `toml:"project_id"`
	TopicID   string `toml:"topic_id"`
}

func (x *PubSub) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "pubsub-project-id",
			Usage:       "PubSub project ID",
			Destination: &x.ProjectID,
			EnvVars:     []string{"OVERSEER_PUBSUB_PROJECT_ID"},
		},
		&cli.StringFlag{
			Name:        "pubsub-topic-id",
			Usage:       "PubSub topic ID",
			Destination: &x.TopicID,
			EnvVars:     []string{"OVERSEER_PUBSUB_TOPIC_ID"},
		},
	}
}

func (x *PubSub) Configure(ctx context.Context) (interfaces.Emitter, error) {
	utils.Logger().Info("Configure PubSub",
		"projectID", x.ProjectID,
		"topicID", x.TopicID,
	)
	return queue.New(ctx, x.ProjectID, x.TopicID)
}
