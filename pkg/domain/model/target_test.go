package model_test

import (
	"testing"

	"github.com/m-mizutani/overseer/pkg/domain/model"
)

func TestTarget_Contains(t *testing.T) {
	testCases := []struct {
		tgt      *model.Target
		name     string
		task     *model.Task
		expected bool
	}{
		{
			tgt: &model.Target{
				Tags: []string{"tag1", "tag2"},
				IDs:  []string{"id1", "id2"},
			},
			name: "Task with matching tag",
			task: &model.Task{
				Tags: []string{"tag1"},
				ID:   "id3",
			},
			expected: true,
		},
		{
			name: "Task with matching ID",
			tgt: &model.Target{
				Tags: []string{"tag1", "tag2"},
				IDs:  []string{"id1", "id2"},
			},
			task: &model.Task{
				Tags: []string{"tag3"},
				ID:   "id1",
			},
			expected: true,
		},
		{
			name: "Task with no matching tag or ID",
			tgt: &model.Target{
				Tags: []string{"tag1", "tag2"},
				IDs:  []string{"id1", "id2"},
			},
			task: &model.Task{
				Tags: []string{"tag3"},
				ID:   "id3",
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.tgt.Contains(tc.task)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}
