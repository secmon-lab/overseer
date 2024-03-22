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
	gt.Equal(t, task.ID, "my-test1")
	gt.Equal(t, task.Title, "Example task")
	gt.Equal(t, task.Description, "This is an example task")
	gt.A(t, task.Tags).Length(2).Have("t1").Have("t2")
	gt.Equal(t, task.Query, testQuery)
	gt.A(t, task.Tests).Length(2).
		At(0, func(t testing.TB, v model.TaskTest) {
			gt.Equal(t, v.YamlPath, "/path/to/yaml/file1.yaml")
			gt.Equal(t, v.Detectable, true)
		}).
		At(1, func(t testing.TB, v model.TaskTest) {
			gt.Equal(t, v.YamlPath, "/path/to/yaml/file2.yaml")
			gt.Equal(t, v.Detectable, false)
		})
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
				Query:       "select * from test",
			},
			isErr: false,
		},
		"valid even if no description": {
			task: model.Task{
				Title: "test",
				Query: "select * from test",
			},
			isErr: false,
		},
		"invalid title": {
			task: model.Task{
				Query: "select * from test",
			},
			isErr: true,
		},
		"invalid query": {
			task: model.Task{
				Title: "test",
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
