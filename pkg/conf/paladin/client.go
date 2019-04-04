package paladin

import (
	"context"
)

const (
	// EventAdd config add event.
	EventAdd EventType = iota
	// EventUpdate config update event.
	EventUpdate
	// EventRemove config remove event.
	EventRemove
)

// EventType is config event.
type EventType int

// Event is watch event.
type Event struct {
	Event EventType
	Key   string
	Value string
}

// Watcher is config watcher.
type Watcher interface {
	WatchEvent(context.Context, ...string) <-chan Event
	Close() error
}

// Setter is value setter.
type Setter interface {
	Set(string) error
}

// Getter is value getter.
type Getter interface {
	// Get a config value by a config key(may be a sven filename).
	Get(string) *Value
	// GetAll return all config key->value map.
	GetAll() *Map
}

// Client is config client.
type Client interface {
	Watcher
	Getter
}
