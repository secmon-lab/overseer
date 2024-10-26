package service

import (
	"compress/gzip"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/overseer/pkg/domain/interfaces"
	"github.com/secmon-lab/overseer/pkg/domain/model"
	"github.com/secmon-lab/overseer/pkg/logging"
)

type gzipWriter struct {
	w  io.WriteCloser
	gz *gzip.Writer
}

func newGzipWriter(w io.WriteCloser) *gzipWriter {
	return &gzipWriter{w: w, gz: gzip.NewWriter(w)}
}

func (x *gzipWriter) Write(p []byte) (int, error) {
	return x.gz.Write(p)
}

func (x *gzipWriter) Close() error {
	rtErr := x.gz.Close()
	if rtErr != nil {
		logging.Default().Error("fail to close gzip writer", "err", rtErr)
	}

	if err := x.w.Close(); err != nil {
		rtErr = goerr.Wrap(err, "fail to close writer").With("err", rtErr)
	}

	return rtErr
}

type gzipReader struct {
	r  io.ReadCloser
	gz *gzip.Reader
}

func newGzipReader(r io.ReadCloser) (*gzipReader, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create gzip reader")
	}
	return &gzipReader{r: r, gz: gz}, nil
}

func (x *gzipReader) Read(p []byte) (int, error) {
	return x.gz.Read(p)
}

func (x *gzipReader) Close() error {
	rtErr := x.gz.Close()
	if rtErr != nil {
		logging.Default().Error("fail to close gzip reader", "err", rtErr)
	}

	if err := x.r.Close(); err != nil {
		rtErr = goerr.Wrap(err, "fail to close reader").With("err", rtErr)
	}

	return rtErr
}

type cacheOptions struct {
	enableGzip bool
}

type CacheOption func(opt *cacheOptions)

func WithGzip() CacheOption {
	return func(opt *cacheOptions) {
		opt.enableGzip = true
	}
}

type FileCache struct {
	id      model.JobID
	baseDir string
	options cacheOptions
}

func (x *FileCache) String() string {
	return "fs://" + x.baseDir
}

func NewFileCache(id model.JobID, baseDir string, options ...CacheOption) (*FileCache, error) {
	dirPath := filepath.Dir(fromIDtoFilePath(baseDir, id, ""))
	if err := os.MkdirAll(dirPath, 0700); err != nil {
		return nil, goerr.Wrap(err, "fail to create baseDir for cache").With("baseDir", baseDir)
	}

	var opt cacheOptions
	for _, o := range options {
		o(&opt)
	}

	return &FileCache{id: id, baseDir: baseDir, options: opt}, nil
}

func fromIDtoFilePath(baseDir string, jobID model.JobID, queryID model.QueryID) string {
	return filepath.Join(baseDir, string(jobID), string(queryID), "data.json")
}

func (x *FileCache) NewWriter(_ context.Context, ID model.QueryID) (io.WriteCloser, error) {
	fpath := fromIDtoFilePath(x.baseDir, x.id, ID)

	dirPath := filepath.Dir(fpath)
	if err := os.MkdirAll(dirPath, 0700); err != nil {
		return nil, goerr.Wrap(err, "fail to create directory for cache").With("dir_path", dirPath)
	}

	fd, err := os.Create(filepath.Clean(fpath))
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create file").With("path", fpath)
	}
	var w io.WriteCloser = fd

	if x.options.enableGzip {
		w = newGzipWriter(w)
	}

	return w, nil
}

func (x *FileCache) NewReader(_ context.Context, ID model.QueryID) (io.ReadCloser, error) {
	fpath := fromIDtoFilePath(x.baseDir, x.id, ID)

	fd, err := os.Open(filepath.Clean(fpath))
	if err != nil {
		return nil, goerr.Wrap(err, "fail to open file").With("path", fpath)
	}

	if x.options.enableGzip {
		return newGzipReader(fd)
	}

	return fd, nil
}

type CloudStorageCache struct {
	id      model.JobID
	bucket  string
	prefix  string
	client  interfaces.CloudStorageClient
	options cacheOptions
}

func (x *CloudStorageCache) String() string {
	return "gcs://" + x.bucket + "/" + x.prefix
}

func fromIDtoCloudStoragePath(prefix string, jobID model.JobID, queryID model.QueryID) string {
	prefix = strings.TrimRight(prefix, "/")
	return strings.Join([]string{prefix, string(jobID), string(queryID), "data.json"}, "/")
}

func NewCloudStorageCache(id model.JobID, bucket string, prefix string, client interfaces.CloudStorageClient, options ...CacheOption) *CloudStorageCache {
	var opt cacheOptions
	for _, o := range options {
		o(&opt)
	}

	return &CloudStorageCache{
		id:      id,
		bucket:  bucket,
		prefix:  prefix,
		client:  client,
		options: opt,
	}
}

func (x *CloudStorageCache) NewWriter(ctx context.Context, ID model.QueryID) (io.WriteCloser, error) {
	w, err := x.client.PutObject(ctx, x.bucket, fromIDtoCloudStoragePath(x.prefix, x.id, ID))
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create writer")
	}

	if x.options.enableGzip {
		w = newGzipWriter(w)
	}

	return w, nil
}

func (x *CloudStorageCache) NewReader(ctx context.Context, ID model.QueryID) (io.ReadCloser, error) {
	r, err := x.client.GetObject(ctx, x.bucket, fromIDtoCloudStoragePath(x.prefix, x.id, ID))
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create reader")
	}

	if x.options.enableGzip {
		return newGzipReader(r)
	}

	return r, nil
}
