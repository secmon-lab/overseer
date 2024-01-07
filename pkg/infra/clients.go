package infra

import (
	"github.com/m-mizutani/overseer/pkg/domain/interfaces"
)

type Clients struct {
	bq    interfaces.BigQuery
	queue interfaces.Queue
}

func New(bq interfaces.BigQuery, queue interfaces.Queue) *Clients {
	return &Clients{
		bq:    bq,
		queue: queue,
	}
}

func (x *Clients) BigQuery() interfaces.BigQuery { return x.bq }
func (x *Clients) Queue() interfaces.Queue       { return x.queue }
