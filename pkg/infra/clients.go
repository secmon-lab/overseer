package infra

import (
	"github.com/m-mizutani/overseer/pkg/domain/interfaces"
)

type Clients struct {
	bq      interfaces.BigQuery
	emitter interfaces.Emitter
}

func New(bq interfaces.BigQuery, emitter interfaces.Emitter) *Clients {
	return &Clients{
		bq:      bq,
		emitter: emitter,
	}
}

func (x *Clients) BigQuery() interfaces.BigQuery { return x.bq }
func (x *Clients) Emitter() interfaces.Emitter   { return x.emitter }
