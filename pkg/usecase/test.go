package usecase

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/overseer/pkg/domain/model"
	"github.com/m-mizutani/overseer/pkg/domain/types"
	"github.com/m-mizutani/overseer/pkg/utils"
	"google.golang.org/api/option"
)

func RunTest(ctx context.Context, emulatorPath string, task *model.Task) error {
	utils.CtxLogger(ctx).Info("Start test", "task", task.Title)
	var hasFailed bool
	for _, tc := range task.Tests {
		if err := runTestCase(ctx, emulatorPath, task, tc); err != nil {
			if errors.Is(err, types.ErrTestFailed) {
				utils.Logger().Error("Test failed",
					slog.Any("task", task.Title),
					slog.Any("file", tc.YamlPath),
					slog.Any("message", err.Error()),
				)
				hasFailed = true
			} else {
				return err
			}
		} else {
			utils.Logger().Info("Test passed",
				slog.Any("task", task.Title),
				slog.Any("file", tc.YamlPath),
			)
		}
	}

	if hasFailed {
		return goerr.Wrap(types.ErrTestFailed, "Some test failed")
	}
	return nil
}

func runTestCase(ctx context.Context, emulatorPath string, task *model.Task, tc model.TaskTest) error {
	const (
		projectID = "test-project"
		testURL   = "http://localhost:9050"
	)

	filePath := filepath.Join(utils.CtxCWD(ctx), tc.YamlPath)
	tmpPath, err := replaceWithCurrentTime(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	// run bigquery-emulator as background process
	args := []string{
		"--project=" + projectID,
		"--data-from-yaml=" + tmpPath,
	}

	bq := exec.Command(emulatorPath, args...)
	if err := bq.Start(); err != nil {
		return goerr.Wrap(err, "Fail to start bigquery-emulator").With("args", args)
	}

	defer func() {
		if err := bq.Process.Kill(); err != nil {
			utils.Logger().Error("Fail to kill bigquery-emulator", "error", err)
		}
	}()

	bqClient, err := bigquery.NewClient(
		ctx,
		projectID,
		option.WithEndpoint(testURL),
		option.WithoutAuthentication(),
	)
	if err != nil {
		return goerr.Wrap(err, "Fail to create bigquery client")
	}

	it, err := bqClient.Query(task.Query).Read(ctx)
	if err != nil {
		return goerr.Wrap(err, "Fail to run query").With("query", task.Query)
	}

	var v []bigquery.Value
	if err := it.Next(&v); err != nil {
		return goerr.Wrap(err, "Fail to read query result")
	}

	if len(v) == 0 && tc.Detectable {
		return goerr.Wrap(types.ErrTestFailed, "Should detect something, but nothing detected")
	}
	if len(v) > 0 && !tc.Detectable {
		return goerr.Wrap(types.ErrTestFailed, "Should not detect anything, but detected")
	}

	return nil
}

func replaceWithCurrentTime(origPath string) (string, error) {
	origData, err := os.ReadFile(filepath.Clean(origPath))
	if err != nil {
		return "", goerr.Wrap(err, "Fail to open test yaml").With("path", origPath)
	}

	tmpData, err := os.CreateTemp("", "overseer-test-*.yaml")
	if err != nil {
		return "", goerr.Wrap(err, "Fail to create temp file")
	}
	replaced := bytes.ReplaceAll(
		origData,
		[]byte("0000-00-00T00:00:00Z"),
		[]byte(time.Now().Format("2006-01-02T15:04:05Z")),
	)

	if _, err := tmpData.Write(replaced); err != nil {
		return "", goerr.Wrap(err, "Fail to write temp file")
	}
	if err := tmpData.Close(); err != nil {
		return "", goerr.Wrap(err, "Fail to close temp file")
	}

	return tmpData.Name(), nil
}
