package model

import (
	"io"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/overseer/pkg/domain/types"
)

type TaskTest struct {
	Detectable bool
	YamlPath   string
}

type Task struct {
	ID          string
	Title       string
	Description string
	Tags        []string
	Query       string
	Dir         string
	Tests       []TaskTest
}

type Tasks []*Task

func (x Tasks) Validate() error {
	ids := map[string]struct{}{}
	for _, t := range x {
		if _, ok := ids[t.ID]; ok {
			return goerr.Wrap(types.ErrInvalidTask, "Duplicated task ID").With("id", t.ID)
		}
		ids[t.ID] = struct{}{}
	}

	return nil
}

func NewTask(r io.Reader) (*Task, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to read task file")
	}

	// parse body
	lines := strings.Split(string(body), "\n")

	var t Task
	var tags string
	params := []struct {
		fieldName string
		dst       *string
		psr       func(line string) error
	}{
		{
			fieldName: "id",
			dst:       &t.ID,
		},
		{
			fieldName: "title",
			dst:       &t.Title,
		},
		{
			fieldName: "description",
			dst:       &t.Description,
		},
		{
			fieldName: "tags",
			dst:       &tags,
		},
		{
			fieldName: "test",
			psr: func(line string) error {
				tt := TaskTest{}
				s := strings.TrimSpace(strings.TrimPrefix(line, "-- test:"))
				v := strings.Split(s, ",")
				if len(v) != 2 {
					return goerr.Wrap(types.ErrInvalidTask, "Invalid test format").With("line", line)
				}

				switch strings.TrimSpace(v[0]) {
				case "true":
					tt.Detectable = true
				case "false":
					tt.Detectable = false
				default:
					return goerr.Wrap(types.ErrInvalidTask, "Invalid test format").With("line", line)
				}
				tt.YamlPath = strings.TrimSpace(v[1])

				t.Tests = append(t.Tests, tt)
				return nil
			},
		},
	}

	for _, line := range lines {
		for _, param := range params {
			prefix := "-- " + param.fieldName + ":"
			if strings.HasPrefix(line, prefix) {
				value := strings.TrimSpace(strings.TrimPrefix(line, prefix))
				switch {
				case param.dst != nil:
					*param.dst = value
				case param.psr != nil:
					if err := param.psr(value); err != nil {
						return nil, err
					}
				}
				break
			}
		}
	}

	tagValues := strings.Split(tags, ",")
	for _, tag := range tagValues {
		t.Tags = append(t.Tags, strings.TrimSpace(tag))
	}

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

	return nil
}
