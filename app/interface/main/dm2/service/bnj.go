package service

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"go-common/library/log"
)

const (
	_bnjShieldCsvURL = "http://i0.hdslb.com/bfs/dm/bnj_shield.csv"
)

func (s *Service) shieldProc() {
	s.shieldFromCsv()
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()
	for range ticker.C {
		s.shieldFromCsv()
	}
}

func (s *Service) shieldFromCsv() {
	resp, err := http.Get(_bnjShieldCsvURL)
	if err != nil {
		log.Error("shieldFromCsv(url:%v) error(%v)", _bnjShieldCsvURL, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Error("shieldFromCsv(url:%v) status(%v)", _bnjShieldCsvURL, resp.StatusCode)
		return
	}
	r := csv.NewReader(resp.Body)
	// ignore first record
	r.Read()
	aids := make([]int64, 0, 100)
	mids := make([]int64, 0, 100)
	for {
		records, err := r.Read()
		if err != nil {
			break
		}
		if len(records) != 2 {
			continue
		}
		// ignore error
		aid, _ := strconv.ParseInt(records[0], 10, 64)
		mid, _ := strconv.ParseInt(records[1], 10, 64)
		if aid > 0 {
			aids = append(aids, aid)
		}
		if mid > 0 {
			mids = append(mids, mid)
		}
	}
	aidMap := make(map[int64]struct{})
	midMap := make(map[int64]struct{})
	for _, aid := range aids {
		aidMap[aid] = struct{}{}
	}
	for _, mid := range mids {
		midMap[mid] = struct{}{}
	}
	s.aidSheild = aidMap
	s.midsSheild = midMap
}
