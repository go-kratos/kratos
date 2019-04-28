package service

import (
	"encoding/binary"
	"go-common/app/service/bbq/recsys-recall/service/index"

	"github.com/Dai0522/go-hash/bloomfilter"
)

// Filter interface
type Filter interface {
	doFilter(uint64) bool
}

// FilterManager .
type FilterManager struct {
	filters map[string]*Filter
}

// NewFilterManager .
func NewFilterManager(args ...interface{}) *FilterManager {
	f := make(map[string]*Filter)
	return &FilterManager{
		filters: f,
	}
}

// SetFilter .
func (fm *FilterManager) SetFilter(name string, f Filter) {
	fm.filters[name] = &f
}

// DoFilter .
func (fm *FilterManager) DoFilter(svid uint64, names ...string) bool {
	res := false

	for _, n := range names {
		if _, ok := fm.filters[n]; !ok {
			continue
		}
		res = res || (*fm.filters[n]).doFilter(svid)
	}
	return res
}

// DefaultFilter .
type DefaultFilter struct{}

func (f *DefaultFilter) doFilter(svid uint64) bool {
	// 状态过滤 state > 0 可进推荐
	fi := index.Index.Get(svid)
	if fi == nil || fi.BasicInfo == nil || fi.BasicInfo.State < 0 || fi.BasicInfo.State == 2 {
		return true
	}
	return false
}

// BloomFilter struct
type BloomFilter struct {
	bf []*bloomfilter.BloomFilter
}

// doFilter .
func (f *BloomFilter) doFilter(svid uint64) bool {
	// 观看记录过滤
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, svid)
	for _, v := range f.bf {
		if v.MightContain(b) {
			return true
		}
	}

	return false
}
