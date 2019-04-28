package file

import (
	"sync"
	"time"
	"go-common/library/log"
)

// States handles list of FileState. One must use NewStates to instantiate a
// file states registry. Using the zero-value is not safe.
type States struct {
	sync.RWMutex

	// states store
	states map[uint64]State
}

// NewStates generates a new states registry.
func NewStates() *States {
	return &States{
		states: map[uint64]State{},
	}
}

// Update updates a state. If previous state didn't exist, new one is created
func (s *States) Update(newState State) {
	s.Lock()
	defer s.Unlock()

	id := newState.ID()

	if _, ok := s.states[id]; ok {
		s.states[id] = newState
		return
	}
	log.V(1).Info("New state added for %s", id)
	s.states[id] = newState
}

// Cleanup cleans up the state array. All states which are older then `older` are removed
// The number of states that were cleaned up is returned.
func (s *States) Cleanup() (int) {
	s.Lock()
	defer s.Unlock()

	currentTime := time.Now()
	statesBefore := len(s.states)

	for inode, state := range s.states {
		if state.Finished && state.TTL > 0 && currentTime.Sub(state.Timestamp) > state.TTL {
			delete(s.states, inode)
		}
	}

	return statesBefore - len(s.states)
}

// GetStates creates copy of the file states.
func (s *States) GetState(id uint64) State {
	s.RLock()
	defer s.RUnlock()

	if _, ok := s.states[id]; ok {
		return s.states[id]
	}
	return State{}
}

// FindPrevious lookups a registered state, that matching the new state.
// Returns a zero-state if no match is found.
func (s *States) FindPrevious(newState State) State {
	s.RLock()
	defer s.RUnlock()

	if s, ok := s.states[newState.ID()]; ok {
		return s
	}

	return State{}
}

// SetStates overwrites all internal states with the given states array
func (s *States) SetStates(states map[uint64]State) {
	s.Lock()
	defer s.Unlock()
	s.states = states
}

// UpdateWithTs updates a state, assigning the given timestamp.
// If previous state didn't exist, new one is created
func (s *States) UpdateWithTs(newState State, ts time.Time) {
	id := newState.ID()
	oldState := s.FindPrevious(newState)
	newState.Timestamp = ts
	s.Lock()
	defer s.Unlock()
	s.states[id] = newState
	if oldState.IsEmpty() {
		log.V(1).Info("New state added for %s", newState.Source)
	}
}
