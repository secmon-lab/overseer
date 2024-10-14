package model

import (
	"bufio"
	"bytes"

	"github.com/m-mizutani/goerr"
	"gopkg.in/yaml.v3"
)

type Query struct {
	metaData queryMetadata
	query    string
}

func (x *Query) ID() string {
	return x.metaData.ID
}

type queryMetadata struct {
	ID string `yaml:"id"`
}

func NewQuery(data []byte) (*Query, error) {
	// extract metadata in header qualified by "/*" and "*/"
	meta, err := extractMetaData(data)
	if err != nil {
		return nil, err
	}

	var q Query
	if err := yaml.Unmarshal(meta, &q.metaData); err != nil {
		return nil, goerr.Wrap(err, "fail to unmarshal metadata")
	}

	q.query = string(data)
	return &q, nil
}

func extractMetaData(data []byte) ([]byte, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var header [][]byte
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.HasPrefix(bytes.TrimSpace(line), []byte("/*")) {
			for scanner.Scan() {
				line = scanner.Bytes()
				if bytes.HasPrefix(bytes.TrimSpace(line), []byte("*/")) {
					return bytes.Join(header, []byte("\n")), nil
				}
				header = append(header, line)
			}
			break
		}
	}

	return nil, goerr.New("metadata not found")
}
