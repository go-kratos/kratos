package collector

import (
	"fmt"
	"sync"

	"go-common/app/service/main/dapper/model"
)

type countBreaker struct {
	rmx  sync.RWMutex
	n    int
	slot map[string]struct{}
}

func (c *countBreaker) Break(key string) error {
	c.rmx.Lock()
	_, ok := c.slot[key]
	c.rmx.Unlock()
	if ok {
		return nil
	}
	c.rmx.Lock()
	c.slot[key] = struct{}{}
	l := len(c.slot)
	c.rmx.Unlock()
	if l <= c.n {
		return nil
	}
	return fmt.Errorf("%s reach limit number %d breaked", key, c.n)
}

func newCountBreaker(n int) *countBreaker {
	return &countBreaker{n: n, slot: make(map[string]struct{})}
}

type serviceBreaker struct {
	rmx  sync.RWMutex
	n    int
	slot map[string]*countBreaker
}

func (s *serviceBreaker) Process(span *model.Span) error {
	s.rmx.RLock()
	operationNameBreaker, ok1 := s.slot[span.ServiceName+"_o"]
	peerServiceBreaker, ok2 := s.slot[span.ServiceName+"_p"]
	s.rmx.RUnlock()
	if !ok1 || !ok2 {
		s.rmx.Lock()
		if !ok1 {
			operationNameBreaker = newCountBreaker(s.n)
			s.slot[span.ServiceName+"_o"] = operationNameBreaker
		}
		if !ok2 {
			peerServiceBreaker = newCountBreaker(s.n)
			s.slot[span.ServiceName+"_p"] = peerServiceBreaker
		}
		s.rmx.Unlock()
	}
	if err := operationNameBreaker.Break(span.OperationName); err != nil {
		return err
	}
	return peerServiceBreaker.Break(span.StringTag("peer.service"))
}

// NewServiceBreakerProcess .
func NewServiceBreakerProcess(n int) Processer {
	return &serviceBreaker{n: n, slot: make(map[string]*countBreaker)}
}
