package interfaces

import (
	"context"
	"io"

	"github.com/secmon-lab/overseer/pkg/domain/model"
)

type CacheService interface {
	NewWriter(ctx context.Context, ID model.QueryID) (io.WriteCloser, error)
	NewReader(ctx context.Context, ID model.QueryID) (io.ReadCloser, error)
	String() string
}

type NotifyService interface {
	Publish(ctx context.Context, alert model.Alert) error
}
