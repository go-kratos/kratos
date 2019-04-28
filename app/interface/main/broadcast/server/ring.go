package server

import (
	"go-common/app/interface/main/broadcast/conf"
	"go-common/app/service/main/broadcast/model"
	"go-common/library/log"
)

// Ring .
type Ring struct {
	// read
	rp   uint64
	num  uint64
	mask uint64
	// TODO split cacheline, many cpu cache line size is 64
	// pad [40]byte
	// write
	wp   uint64
	data []model.Proto
}

// NewRing .
func NewRing(num int) *Ring {
	r := new(Ring)
	r.init(uint64(num))
	return r
}

// Init .
func (r *Ring) Init(num int) {
	r.init(uint64(num))
}

func (r *Ring) init(num uint64) {
	// 2^N
	if num&(num-1) != 0 {
		for num&(num-1) != 0 {
			num &= (num - 1)
		}
		num = num << 1
	}
	r.data = make([]model.Proto, num)
	r.num = num
	r.mask = r.num - 1
}

// Get .
func (r *Ring) Get() (proto *model.Proto, err error) {
	if r.rp == r.wp {
		return nil, ErrRingEmpty
	}
	proto = &r.data[r.rp&r.mask]
	return
}

// GetAdv .
func (r *Ring) GetAdv() {
	r.rp++
	if conf.Conf.Broadcast.Debug {
		log.Info("ring rp: %d, idx: %d", r.rp, r.rp&r.mask)
	}
}

// Set .
func (r *Ring) Set() (proto *model.Proto, err error) {
	if r.wp-r.rp >= r.num {
		return nil, ErrRingFull
	}
	proto = &r.data[r.wp&r.mask]
	return
}

// SetAdv .
func (r *Ring) SetAdv() {
	r.wp++
	if conf.Conf.Broadcast.Debug {
		log.Info("ring wp: %d, idx: %d", r.wp, r.wp&r.mask)
	}
}

// Reset .
func (r *Ring) Reset() {
	r.rp = 0
	r.wp = 0
	// prevent pad compiler optimization
	// r.pad = [40]byte{}
}
