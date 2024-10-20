package service

import (
	"context"
	"encoding/json"
	"io"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-as-code/overseer/pkg/domain/interfaces"
	"github.com/secmon-as-code/overseer/pkg/domain/model"
)

type NotifyPubSub struct {
	client interfaces.PubSubClient
	topic  string
}

func NewNotifyPubSub(client interfaces.PubSubClient, topic string) *NotifyPubSub {
	return &NotifyPubSub{client: client}
}

func (x *NotifyPubSub) Publish(ctx context.Context, alert model.Alert) error {
	raw, err := json.Marshal(alert)
	if err != nil {
		return goerr.Wrap(err, "fail to marshal alert fot pubsub notification")
	}

	return x.client.Publish(ctx, x.topic, raw)
}

type NotifyWriter struct {
	w io.Writer
}

func NewNotifyWriter(w io.Writer) *NotifyWriter {
	return &NotifyWriter{w: w}
}

func (x *NotifyWriter) Publish(ctx context.Context, alert model.Alert) error {
	encoder := json.NewEncoder(x.w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(alert)
}
