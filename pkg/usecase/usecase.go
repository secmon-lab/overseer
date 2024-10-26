package usecase

import (
	"github.com/secmon-lab/overseer/pkg/adaptor"
)

type UseCase struct {
	clients *adaptor.Clients
}

func New(clients *adaptor.Clients) *UseCase {
	return &UseCase{
		clients: clients,
	}
}
