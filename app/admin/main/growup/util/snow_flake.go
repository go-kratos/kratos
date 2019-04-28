package util

import (
	"errors"
	"sync"
	"time"
)

const (
	// MaxSequence  max sequence in one time
	MaxSequence = 1000
	// WorkerBit worker bit
	WorkerBit = 10
	// SequenceBit sequence bit
	SequenceBit = 10
)

// SnowFlake snow flake
type SnowFlake struct {
	sync.Mutex
	lastTimestamp int64
	sequence      int64
	workerID      int64
}

// NewSnowFlake new
func NewSnowFlake() *SnowFlake {
	return &SnowFlake{
		workerID: time.Now().UnixNano() % 1000,
	}
}

// Generate generate
func (s *SnowFlake) Generate() (int64, error) {
	s.Lock()
	defer s.Unlock()

	now := time.Now().UnixNano() / 1e6
	if now == s.lastTimestamp {
		s.sequence = (s.sequence + 1) % MaxSequence
		if s.sequence == 0 {
			now = s.waitNextMill(now)
		}
	} else {
		s.sequence = 0
	}

	if now < s.lastTimestamp {
		return 0, errors.New("inner time error")
	}
	s.lastTimestamp = now
	return s.generate() % 1000000000, nil
}

func (s *SnowFlake) generate() int64 {
	return (s.lastTimestamp << (WorkerBit + SequenceBit)) | (s.workerID << SequenceBit) | s.sequence
}

func (s *SnowFlake) waitNextMill(t int64) int64 {
	for t == s.lastTimestamp {
		time.Sleep(100 * time.Microsecond)
		t = time.Now().UnixNano() / 1e6
	}
	return t
}
