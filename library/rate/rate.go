package rate

import (
	"context"
)

// Op operations type.
type Op int

const (
	// Success opertion type: success
	Success Op = iota
	// Ignore opertion type: ignore
	Ignore
	// Drop opertion type: drop
	Drop
)

// Limiter limit interface.
type Limiter interface {
	Allow(ctx context.Context) (func(Op), error)
}
