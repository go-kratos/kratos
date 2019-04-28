package service

import (
	"context"
	"fmt"
	"go-common/library/log"
	"math/rand"
	"sort"
	"time"
)

const (
	_nblocks     = 5
	_redisPrefix = "RECALL:NEWPUB:%d"
)

// GenRealTimeInvertedIndex 实时倒排标签
func (s *Service) GenRealTimeInvertedIndex() {
	svids, err := s.dao.FetchNewincomeVideo()
	if err != nil || svids == nil || len(svids) == 0 {
		log.Error("GenRealTimeInvertedIndex FetchNewincomeVideo err[%v] svids[%v]", err, svids)
		return
	}
	// svid 乱序
	rand.Seed(time.Now().Unix())
	sort.Slice(svids, func(i int, j int) bool {
		return rand.Float32() > 0.5
	})

	// 平均分为5份
	offset := 0
	blocks := len(svids) / _nblocks
	invertedIndex := make([][]int64, _nblocks)
	for i := 0; i < _nblocks; i++ {
		invertedIndex[i] = make([]int64, 0)
		for j := 0; j < blocks; j++ {
			invertedIndex[i] = append(invertedIndex[i], svids[offset+j])
		}
		offset = offset + blocks
	}
	if blocks*_nblocks < len(svids) {
		invertedIndex[_nblocks-1] = append(invertedIndex[_nblocks-1], svids[len(svids)-1])
	}

	log.Info("GenRealTimeInvertedIndex invertedIndex[%v]", invertedIndex)

	// 序列化后写入redis
	for i, v := range invertedIndex {
		key := fmt.Sprintf(_redisPrefix, i)
		err = s.dao.SetInvertedIndex(context.Background(), key, v)
		if err != nil {
			log.Error("GenRealTimeInvertedIndex SetInvertedIndex err[%v]", err)
		}
	}

	log.Info("finish [GenRealTimeInvertedIndex]")
}
