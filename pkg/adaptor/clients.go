package adaptor

import "github.com/secmon-as-code/overseer/pkg/domain/interfaces"

type Clients struct {
	bigQuery interfaces.BigQueryClient
	pubSub   interfaces.PubSubClient
}

func (x *Clients) BigQuery() interfaces.BigQueryClient {
	return x.bigQuery
}

func (x *Clients) PubSub() interfaces.PubSubClient {
	return x.pubSub
}

func New(options ...Option) *Clients {
	c := &Clients{}
	for _, opt := range options {
		opt(c)
	}
	return c
}

type Option func(*Clients)

func WithBigQuery(client interfaces.BigQueryClient) Option {
	return func(c *Clients) {
		c.bigQuery = client
	}
}

func WithPubSub(client interfaces.PubSubClient) Option {
	return func(c *Clients) {
		c.pubSub = client
	}
}
