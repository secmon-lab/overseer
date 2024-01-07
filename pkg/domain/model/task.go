package model

import (
	"io"
	"strconv"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/overseer/pkg/domain/types"
)

type Task struct {
	Title       string
	Description string
	Limit       int
	Query       string
}

func NewTask(r io.Reader) (*Task, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to read task file")
	}

	// parse body
	lines := strings.Split(string(body), "\n")

	var t Task
	var limit string
	params := []struct {
		fieldName string
		dst       *string
	}{
		{
			fieldName: "title",
			dst:       &t.Title,
		},
		{
			fieldName: "description",
			dst:       &t.Description,
		},
		{
			fieldName: "limit",
			dst:       &limit,
		},
	}

	for _, line := range lines {
		for _, param := range params {
			prefix := "-- " + param.fieldName + ":"
			if strings.HasPrefix(line, prefix) {
				*param.dst = strings.TrimSpace(strings.TrimPrefix(line, prefix))
			}
		}
	}

	nLimit, err := strconv.ParseInt(limit, 10, 32)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to parse limit").With("limit", limit)
	}

	t.Limit = int(nLimit)
	t.Query = string(body)

	if err := t.Validate(); err != nil {
		return nil, err
	}

	return &t, nil
}

func (x *Task) Validate() error {
	if x.Title == "" {
		return goerr.Wrap(types.ErrInvalidTask, "Title must not be empty")
	}
	if x.Query == "" {
		return goerr.Wrap(types.ErrInvalidTask, "Query must not be empty")
	}
	if x.Limit <= 0 {
		return goerr.Wrap(types.ErrInvalidTask, "Limit must be positive number")
	}

	return nil
}
