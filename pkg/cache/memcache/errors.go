package memcache

import (
	"errors"
	"fmt"
	"strings"

	pkgerr "github.com/pkg/errors"
)

var (
	// ErrNotFound not found
	ErrNotFound = errors.New("memcache: key not found")
	// ErrExists exists
	ErrExists = errors.New("memcache: key exists")
	// ErrNotStored not stored
	ErrNotStored = errors.New("memcache: key not stored")
	// ErrCASConflict means that a CompareAndSwap call failed due to the
	// cached value being modified between the Get and the CompareAndSwap.
	// If the cached value was simply evicted rather than replaced,
	// ErrNotStored will be returned instead.
	ErrCASConflict = errors.New("memcache: compare-and-swap conflict")

	// ErrPoolExhausted is returned from a pool connection method (Store, Get,
	// Delete, IncrDecr, Err) when the maximum number of database connections
	// in the pool has been reached.
	ErrPoolExhausted = errors.New("memcache: connection pool exhausted")
	// ErrPoolClosed pool closed
	ErrPoolClosed = errors.New("memcache: connection pool closed")
	// ErrConnClosed conn closed
	ErrConnClosed = errors.New("memcache: connection closed")
	// ErrMalformedKey is returned when an invalid key is used.
	// Keys must be at maximum 250 bytes long and not
	// contain whitespace or control characters.
	ErrMalformedKey = errors.New("memcache: malformed key is too long or contains invalid characters")
	// ErrValueSize item value size must less than 1mb
	ErrValueSize = errors.New("memcache: item value size must not greater than 1mb")
	// ErrStat stat error for monitor
	ErrStat = errors.New("memcache unexpected errors")
	// ErrItem item nil.
	ErrItem = errors.New("memcache: item object nil")
	// ErrItemObject object type Assertion failed
	ErrItemObject = errors.New("memcache: item object protobuf type assertion failed")
)

type protocolError string

func (pe protocolError) Error() string {
	return fmt.Sprintf("memcache: %s (possible server error or unsupported concurrent read by application)", string(pe))
}

func formatErr(err error) string {
	e := pkgerr.Cause(err)
	switch e {
	case ErrNotFound, ErrExists, ErrNotStored, nil:
		return ""
	default:
		es := e.Error()
		switch {
		case strings.HasPrefix(es, "read"):
			return "read timeout"
		case strings.HasPrefix(es, "dial"):
			return "dial timeout"
		case strings.HasPrefix(es, "write"):
			return "write timeout"
		case strings.Contains(es, "EOF"):
			return "eof"
		case strings.Contains(es, "reset"):
			return "reset"
		case strings.Contains(es, "broken"):
			return "broken pipe"
		default:
			return "unexpected err"
		}
	}
}
