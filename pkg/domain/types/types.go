package types

import "github.com/google/uuid"

type (
	RequestID string
)

func NewRequestID() RequestID {
	return RequestID(uuid.NewString())
}
