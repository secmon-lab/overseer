package model_test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/overseer/pkg/domain/model"
)

//go:embed testdata/test_task.sql
var testQuery string

func TestNewTask(t *testing.T) {
	r := strings.NewReader(testQuery)
	task := gt.R1(model.NewTask(r)).NoError(t)
	gt.Equal(t, task.Title, "Example task")
	gt.Equal(t, task.Description, "This is an example task")
	gt.Equal(t, task.Query, testQuery)
}

func Test_Task_Validate(t *testing.T) {
	testCases := map[string]struct {
		task  model.Task
		isErr bool
	}{
		"valid": {
			task: model.Task{
				Title:       "test",
				Description: "test",
				Limit:       10,
				Query:       "select * from test",
			},
			isErr: false,
		},
		"valid even if no description": {
			task: model.Task{
				Title: "test",
				Query: "select * from test",
				Limit: 1,
			},
			isErr: false,
		},
		"invalid title": {
			task: model.Task{
				Query: "select * from test",
				Limit: 1,
			},
			isErr: true,
		},
		"invalid query": {
			task: model.Task{
				Title: "test",
				Limit: 1,
			},
			isErr: true,
		},
		"invalid limit": {
			task: model.Task{
				Title: "test",
				Limit: 0,
				Query: "select * from test",
			},
			isErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.task.Validate()
			if tc.isErr {
				gt.True(t, err != nil)
			} else {
				gt.True(t, err == nil)
			}
		})
	}
}
