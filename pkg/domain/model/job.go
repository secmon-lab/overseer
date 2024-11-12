package model

import (
	"strings"
	"time"

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

	now := time.Now()
	return JobID(now.Format("job200601021504_") + strings.ReplaceAll(id.String(), "-", ""))
}
