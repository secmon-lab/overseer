package model

import (
	"strings"

	"github.com/google/uuid"
	"github.com/secmon-lab/overseer/pkg/logging"
)

type JobID string

func NewJobID() JobID {
	id, err := uuid.NewV7()
	if err != nil {
		logging.Default().Error("fail to generate new JobID", "err", err)
		panic(err)
	}

	return JobID("job" + strings.ReplaceAll(id.String(), "-", ""))
}
