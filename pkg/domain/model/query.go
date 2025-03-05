package model

import (
	"bufio"
	"bytes"

	"github.com/m-mizutani/goerr/v2"
	"gopkg.in/yaml.v3"
)

type QueryID string

type Query struct {
	metaData queryMetadata
	query    string
}

func (x *Query) Validate() error {
	if x.metaData.ID == "" {
		return goerr.New("ID is required")
	}
	return nil
}

func (x *Query) ID() QueryID {
	return x.metaData.ID
}

func (x *Query) String() string {
	return x.query
}

type queryMetadata struct {
	ID QueryID `yaml:"id"`
}

func MustNewQuery(name string, data []byte) *Query {
	q, err := NewQuery(name, data)
	if err != nil {
		panic(err)
	}
	return q
}

var errMetadataNotFound = goerr.New("metadata not found")

func NewQuery(name string, data []byte) (*Query, error) {
	// extract metadata in header qualified by "/*" and "*/"
	meta, err := extractMetaData(data)
	if err != nil {
		if err != errMetadataNotFound {
			return nil, err
		}

		meta = []byte("id: " + name)
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

	return nil, errMetadataNotFound
}

type Queries []*Query

func (x Queries) IDs() []QueryID {
	ids := []QueryID{}
	for _, q := range x {
		ids = append(ids, q.ID())
	}
	return ids
}

func (x Queries) FindByID(id QueryID) *Query {
	for _, q := range x {
		if q.ID() == id {
			return q
		}
	}
	return nil
}

func (x Queries) Validate() error {
	ids := map[QueryID]struct{}{}
	for _, q := range x {
		if err := q.Validate(); err != nil {
			return goerr.Wrap(err, "invalid query")
		}

		if _, ok := ids[q.ID()]; ok {
			return goerr.New("duplicated query ID", goerr.V("id", q.ID()))
		}
		ids[q.ID()] = struct{}{}
	}

	return nil
}
