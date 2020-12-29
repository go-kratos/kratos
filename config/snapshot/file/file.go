package file

import (
	"encoding/json"
	"io/ioutil"

	"github.com/go-kratos/kratos/v2/config/snapshot"
)

var _ snapshot.Store = (*fileStore)(nil)

type fileStore struct {
	path string
}

// NewSnapshot new a snapshot store.
func NewSnapshot(path string) snapshot.Store {
	return &fileStore{path: path}
}

func (s *fileStore) Read() (*snapshot.Snapshot, error) {
	var sst snapshot.Snapshot
	data, err := ioutil.ReadFile(s.path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &sst); err != nil {
		return nil, err
	}
	return &sst, nil
}

func (s *fileStore) Write(sst *snapshot.Snapshot) error {
	data, err := json.Marshal(sst)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.path, data, 0666)
}
