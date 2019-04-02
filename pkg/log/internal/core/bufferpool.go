package core

var (
	_pool = NewPool(_size)
	// GetPool retrieves a buffer from the pool, creating one if necessary.
	GetPool = _pool.Get
)
