package file

import (
	"time"
	"os"
	"fmt"
	"path/filepath"
	"syscall"
)

type State struct {
	Source    string            `json:"source"`
	Offset    int64             `json:"offset"`
	Inode     uint64            `json:"inode"`
	Fileinfo  os.FileInfo       `json:"-"` // the file info
	Timestamp time.Time         `json:"timestamp"`
	Finished  bool              `json:"finished"`
	Meta      map[string]string `json:"meta"`
	TTL       time.Duration     `json:"ttl"`
}

// NewState creates a new file state
func NewState(fileInfo os.FileInfo, path string) State {
	stat := fileInfo.Sys().(*syscall.Stat_t)
	return State{
		Fileinfo:  fileInfo,
		Inode:     stat.Ino,
		Source:    path,
		Finished:  false,
		Timestamp: time.Now(),
		TTL:       -1, // By default, state does have an infinite ttl
		Meta:      nil,
	}
}

func (s *State) ID() uint64 {
	return s.Inode
}

// IsEqual compares the state to an other state supporting stringer based on the unique string
func (s *State) IsEqual(c *State) bool {
	return s.ID() == c.ID()
}

// IsEmpty returns true if the state is empty
func (s *State) IsEmpty() bool {
	return s.Inode == 0 &&
		s.Source == "" &&
		len(s.Meta) == 0 &&
		s.Timestamp.IsZero()
}

func getFileState(path string, info os.FileInfo) (State, error) {
	var err error
	var absolutePath string
	absolutePath, err = filepath.Abs(path)
	if err != nil {
		return State{}, fmt.Errorf("could not fetch abs path for file %s: %s", absolutePath, err)
	}
	// Create new state for comparison
	newState := NewState(info, absolutePath)
	return newState, nil
}
