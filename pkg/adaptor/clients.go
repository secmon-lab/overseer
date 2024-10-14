package adaptor

import "github.com/secmon-as-code/overseer/pkg/interfaces"

type Clients struct {
	CloudStorage interfaces.CloudStorageClient
	BigQuery     interfaces.BigQueryClient
	PubSub       interfaces.PubSubClient
	Policy       interfaces.PolicyClient
}

func New(options ...Option) (*Clients, error) {
	c := &Clients{}
	for _, opt := range options {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

type Option func(*Clients) error

func WithCloudStorage(client interfaces.CloudStorageClient) Option {
	return func(c *Clients) error {
		c.CloudStorage = client
		return nil
	}
}

func WithBigQuery(client interfaces.BigQueryClient) Option {
	return func(c *Clients) error {
		c.BigQuery = client
		return nil
	}
}

func WithPubSub(client interfaces.PubSubClient) Option {
	return func(c *Clients) error {
		c.PubSub = client
		return nil
	}
}

func WithPolicy(client interfaces.PolicyClient) Option {
	return func(c *Clients) error {
		c.Policy = client
		return nil
	}
}
