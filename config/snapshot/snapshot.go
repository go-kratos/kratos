package snapshot

import "github.com/go-kratos/kratos/v2/config/provider"

// Store is config snapshot store.
type Store interface {
	Read() (*Snapshot, error)
	Write(*Snapshot) error
}

// Snapshot is config snapshot.
type Snapshot struct {
	Sources []provider.KeyValue
	Version string
}
