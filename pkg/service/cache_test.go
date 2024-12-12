package service_test

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/overseer/pkg/adaptor/cs"
	"github.com/secmon-lab/overseer/pkg/domain/interfaces"
	"github.com/secmon-lab/overseer/pkg/domain/model"
	"github.com/secmon-lab/overseer/pkg/service"
)

func TestCacheFile(t *testing.T) {
	d := os.TempDir()
	ctx := context.Background()
	_, id1 := model.NewJobID(ctx)
	_, id2 := model.NewJobID(ctx)
	svc1, err := service.NewFileCache(id1, d)
	gt.NoError(t, err)
	svc2, err := service.NewFileCache(id2, d)
	gt.NoError(t, err)

	testCache(t, svc1, svc2)
}

func TestCacheFileGzip(t *testing.T) {
	d := os.TempDir()

	ctx := context.Background()
	_, id1 := model.NewJobID(ctx)
	_, id2 := model.NewJobID(ctx)

	svc1, err := service.NewFileCache(id1, d, service.WithGzip())
	gt.NoError(t, err)
	svc2, err := service.NewFileCache(id2, d, service.WithGzip())
	gt.NoError(t, err)

	testCache(t, svc1, svc2)
}

func TestCacheCloudStorage(t *testing.T) {
	bucketName, ok := os.LookupEnv("TEST_CLOUD_STORAGE_BUCKET_NAME")
	if !ok {
		t.Skip("TEST_CLOUD_STORAGE_BUCKET_NAME is not set")
	}

	client, err := cs.NewClient(context.Background())

	ctx := context.Background()
	_, id1 := model.NewJobID(ctx)
	_, id2 := model.NewJobID(ctx)

	svc1 := service.NewCloudStorageCache(id1, bucketName, "overseer-test", client)
	gt.NoError(t, err)
	svc2 := service.NewCloudStorageCache(id2, bucketName, "overseer-test", client)
	gt.NoError(t, err)

	testCache(t, svc1, svc2)
}

func TestCacheCloudStorageGzip(t *testing.T) {
	bucketName, ok := os.LookupEnv("TEST_CLOUD_STORAGE_BUCKET_NAME")
	if !ok {
		t.Skip("TEST_CLOUD_STORAGE_BUCKET_NAME is not set")
	}

	client, err := cs.NewClient(context.Background())

	ctx := context.Background()
	_, id1 := model.NewJobID(ctx)
	_, id2 := model.NewJobID(ctx)

	svc1 := service.NewCloudStorageCache(id1, bucketName, "overseer-test", client, service.WithGzip())
	gt.NoError(t, err)
	svc2 := service.NewCloudStorageCache(id2, bucketName, "overseer-test", client, service.WithGzip())
	gt.NoError(t, err)

	testCache(t, svc1, svc2)
}

func testCache(t *testing.T, svc1, svc2 interfaces.CacheService) {
	ctx := context.Background()
	t.Run("write data to cache", func(t *testing.T) {
		w, err := svc1.NewWriter(ctx, "test1")
		gt.NoError(t, err)

		_, err = w.Write([]byte("blue"))
		gt.NoError(t, err)
		gt.NoError(t, w.Close())
	})

	t.Run("read data from cache", func(t *testing.T) {
		r, err := svc1.NewReader(ctx, "test1")
		gt.NoError(t, err)

		buf, err := io.ReadAll(r)
		gt.NoError(t, err)
		gt.Equal(t, string(buf), "blue")
	})

	t.Run("can not read data from other id", func(t *testing.T) {
		r, err := svc1.NewReader(ctx, "test2")
		gt.Error(t, err)
		gt.Equal(t, r, nil)
	})

	t.Run("can not read data from other job", func(t *testing.T) {
		r, err := svc2.NewReader(ctx, "test1")
		gt.Error(t, err)
		gt.Equal(t, r, nil)
	})
}
