package usecase

import (
	"context"

	"github.com/secmon-as-code/overseer/pkg/adaptor"
	"github.com/secmon-as-code/overseer/pkg/domain/model"
)

type UseCase struct {
	clients *adaptor.Clients
}

func New(clients *adaptor.Clients) *UseCase {
	return &UseCase{
		clients: clients,
	}
}

func (x *UseCase) QueryAndDump(ctx context.Context, queries []*model.Query) error {
	return nil
}
