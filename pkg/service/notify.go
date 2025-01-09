package service

import (
	"context"
	"encoding/json"
	"io"

	"github.com/m-mizutani/goerr/v2"
	"github.com/secmon-lab/overseer/pkg/domain/interfaces"
	"github.com/secmon-lab/overseer/pkg/domain/model"
)

type NotifyPubSub struct {
	client interfaces.PubSubClient
	topic  string
}

func NewNotifyPubSub(client interfaces.PubSubClient, topic string) *NotifyPubSub {
	return &NotifyPubSub{client: client, topic: topic}
}

func (x *NotifyPubSub) Publish(ctx context.Context, alert model.Alert) error {
	raw, err := json.Marshal(alert)
	if err != nil {
		return goerr.Wrap(err, "fail to marshal alert fot pubsub notification")
	}

	return x.client.Publish(ctx, x.topic, raw)
}

type NotifyWriter struct {
	w       io.Writer
	encoder *json.Encoder
}

func NewNotifyWriter(w io.Writer) *NotifyWriter {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return &NotifyWriter{w: w, encoder: encoder}
}

func (x *NotifyWriter) Publish(ctx context.Context, alert model.Alert) error {
	if err := x.encoder.Encode(alert); err != nil {
		return goerr.Wrap(err, "fail to encode alert")
	}

	return nil
}
