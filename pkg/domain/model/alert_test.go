package model_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/overseer/pkg/domain/model"
)

func TestAlert(t *testing.T) {
	body := model.AlertBody{
		Title: "test alert",
	}
	alert, err := model.NewAlert(context.Background(), body)
	gt.NoError(t, err)

	gt.NotEqual(t, alert.ID, "")
	gt.False(t, alert.Timestamp.IsZero())
}

func TestParseAlert(t *testing.T) {
	testCases := map[string]struct {
		input string
		isErr bool
		eval  func(t *testing.T, alert *model.Alert)
	}{
		"valid alert": {
			input: `{
				"title": "alert-1",
				"timestamp": "2021-03-01T00:00:00Z",
				"attrs": {
					"key1": "value1"
				}
			}`,
			isErr: false,
			eval: func(t *testing.T, alert *model.Alert) {
				gt.Equal(t, alert.Title, "alert-1")
				gt.Equal(t, alert.Attrs["key1"], "value1")
				gt.True(t, alert.Timestamp.Equal(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)))
			},
		},
		"invalid timestamp": {
			input: `{
				"title": "alert-1",
				"timestamp": "invalid",
				"attrs": {
					"key1": "value1"
				}
			}`,
			isErr: true,
		},
		"float timestamp": {
			input: `{
				"title": "alert-1",
				"timestamp": 1614556800.0,
				"attrs": {
					"key1": "value1"
				}
			}`,
			isErr: false,
			eval: func(t *testing.T, alert *model.Alert) {
				gt.True(t, alert.Timestamp.Equal(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)))
			},
		},
		"no title": {
			input: `{
				"timestamp": "2021-03-01T00:00:00Z",
				"attrs": {
					"key1": "value1"
				}
			}`,
			isErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var body model.AlertBody
			gt.NoError(t, json.Unmarshal([]byte(tc.input), &body))

			ctx := context.Background()
			alert, err := model.NewAlert(ctx, body)
			if tc.isErr {
				gt.True(t, err != nil)
			} else {
				gt.NoError(t, err)
				tc.eval(t, alert)
			}
		})
	}
}
