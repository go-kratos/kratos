package service

import (
	"context"
	"strconv"
	"strings"

	"go-common/library/log"
)

// GenBloomFilter .
func (s *Service) GenBloomFilter() {
	log.Info("run [%s]", "GenBloomFilter")
	result, err := s.dao.FetchMidView(context.Background())
	if err != nil {
		log.Error("FetchMidView: %v", err)
		return
	}
	s.bloomFilter(result)

	result, err = s.dao.FetchBuvidView(context.Background())
	if err != nil {
		log.Error("FetchBuvidView: %v", err)
		return
	}
	s.bloomFilter(result)
	log.Info("finish [%s]", "GenBloomFilter")
}

func (s *Service) bloomFilter(result []string) {
	m := make(map[string][]uint64)
	for _, v := range result {
		items := strings.Split(v, "\u0001")
		if len(items) != 2 {
			continue
		}
		if _, ok := m[items[0]]; !ok {
			m[items[0]] = []uint64{}
		}
		svid, _ := strconv.Atoi(items[1])
		m[items[0]] = append(m[items[0]], uint64(svid))
	}
	for k, v := range m {
		if k == "" {
			continue
		}
		if err := s.dao.InsertBloomFilter(context.Background(), k, v); err != nil {
			log.Error("InsertBloomFilter: %v", err)
			continue
		}
	}
}
