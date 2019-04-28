package ai

import (
	"context"
	"sync"
)

type Loader func(context.Context) (map[int64]int64, error)

type AI struct {
	whites map[int64]int64
	sync.RWMutex
}

func New() (a *AI) {
	return &AI{
		whites: make(map[int64]int64),
	}
}

func (a *AI) White(mid int64) (num int64, ok bool) {
	a.RLock()
	defer a.RUnlock()
	num, ok = a.whites[mid]
	return
}

func (a *AI) LoadWhite(c context.Context, loader Loader) (err error) {
	var (
		whites map[int64]int64
	)
	if whites, err = loader(c); err != nil {
		return
	}
	a.Lock()
	a.whites = whites
	a.Unlock()
	return
}
