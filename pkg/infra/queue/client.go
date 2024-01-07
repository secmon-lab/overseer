package queue

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/overseer/pkg/domain/model"
	"google.golang.org/api/option"
)

type pubsubClient struct {
	client *pubsub.Client
	topic  *pubsub.Topic
}

func New(ctx context.Context, projectID, topicID string, options ...option.ClientOption) (*pubsubClient, error) {
	client, err := pubsub.NewClient(ctx, projectID, options...)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to create PubSub client")
	}
	topic := client.Topic(topicID)

	return &pubsubClient{
		client: client,
		topic:  topic,
	}, nil
}

func (x *pubsubClient) Publish(ctx context.Context, alert *model.Alert) error {
	data, err := json.Marshal(alert)
	if err != nil {
		return goerr.Wrap(err, "Fail to marshal alert")
	}
	result := x.topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	if _, err := result.Get(ctx); err != nil {
		return goerr.Wrap(err, "Fail to publish alert")
	}
	return nil
}
