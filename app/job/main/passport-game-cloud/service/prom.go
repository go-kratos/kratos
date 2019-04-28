package service

import (
	"context"
	"fmt"
	"time"

	"go-common/library/log"
	"go-common/library/stat/prom"
)

const (
	_timeFormat = "2006-01-02 15:04:05"
)

var (
	_names = make(map[string]struct{})

	_loc = time.Now().Location()
)

// IntervalConfig interval config.
type IntervalConfig struct {
	Name string
	Rate int64
}

// Interval prom interval.
type Interval struct {
	name    string
	rate    int64
	counter int64
}

// NewInterval new a interval instance.
func NewInterval(c *IntervalConfig) *Interval {
	if _, ok := _names[c.Name]; ok {
		panic(fmt.Sprintf("%s already exists", c.Name))
	}
	_names[c.Name] = struct{}{}
	if c.Rate <= 0 {
		c.Rate = 1000
	}
	return &Interval{
		name:    c.Name,
		rate:    c.Rate,
		counter: 0,
	}
}

// MTS get mtime ts from mtime str with interval's counter.
func (s *Interval) MTS(c context.Context, mtStr string) int64 {
	s.counter++
	if s.counter%s.rate == 0 {
		return 0
	}
	if mtStr == "" {
		return 0
	}
	t, err := time.ParseInLocation(_timeFormat, mtStr, _loc)
	if err != nil {
		log.Error("failed to parse mtime str for %s, time.ParseInLocation(%s, %s, %v) error(%v)", s.name, _timeFormat, mtStr, _loc, err)
		return 0
	}
	return t.Unix()
}

// Prom prom interval if mts > 0.
func (s *Interval) Prom(c context.Context, mts int64) {
	if mts > 0 {
		prom.BusinessInfoCount.State(s.name, time.Now().Unix()-mts)
	}
}
