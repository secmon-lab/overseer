package model

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/overseer/pkg/domain/types"
)

type Target struct {
	Tags []string
	IDs  []string
}

func (x *Target) Contains(t *Task) bool {
	for _, tag := range x.Tags {
		for _, tgt := range x.Tags {
			if tag == tgt {
				return true
			}
		}
	}

	for _, id := range x.IDs {
		if id == t.ID {
			return true
		}
	}

	return false
}

func (x *Target) Validate() error {
	if len(x.Tags) == 0 && len(x.IDs) == 0 {
		return goerr.Wrap(types.ErrInvalidOption, "No target, specify tags or IDs")
	}

	return nil
}
