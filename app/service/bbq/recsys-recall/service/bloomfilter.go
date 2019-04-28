package service

import (
	"context"
	"fmt"
	"go-common/library/log"
	"strconv"
	"time"

	"github.com/Dai0522/go-hash/bloomfilter"
)

const (
	_baseBfKey      = "BBQ:BF:V1:%s"
	_baseBfKeyDaily = "BBQ:BF:V1:%s:%s"
)

func (s *Service) loadBloomFilter(ctx *context.Context, mid int64, buvid string) (result []*bloomfilter.BloomFilter) {
	dt := time.Now().Format("20060102")
	lastDt := time.Now().AddDate(0, 0, -1).Format("20060102")
	if mid != 0 {
		// history part
		key := fmt.Sprintf(_baseBfKey, strconv.FormatInt(mid, 10))
		bf, err := s.dao.LoadBloomFilter(ctx, key)
		if err != nil {
			log.Errorv(*ctx, log.KV("MID_HIS_BF", err))
		} else {
			result = append(result, bf)
		}
		// daily part
		key = fmt.Sprintf(_baseBfKeyDaily, strconv.FormatInt(mid, 10), dt)
		bf, err = s.dao.LoadBloomFilter(ctx, key)
		if err != nil {
			log.Errorv(*ctx, log.KV("MID_DAILY_BF", err))
		} else {
			result = append(result, bf)
		}
		// lastday part
		key = fmt.Sprintf(_baseBfKeyDaily, strconv.FormatInt(mid, 10), lastDt)
		bf, err = s.dao.LoadBloomFilter(ctx, key)
		if err != nil {
			log.Errorv(*ctx, log.KV("MID_DAILY_BF", err))
		} else {
			result = append(result, bf)
		}
	}

	if buvid != "" {
		// history part
		key := fmt.Sprintf(_baseBfKey, buvid)
		bf, err := s.dao.LoadBloomFilter(ctx, key)
		if err != nil {
			log.Errorv(*ctx, log.KV("MID_HIS_BF", err))
		} else {
			result = append(result, bf)
		}
		// daily part
		key = fmt.Sprintf(_baseBfKeyDaily, buvid, dt)
		bf, err = s.dao.LoadBloomFilter(ctx, key)
		if err != nil {
			log.Errorv(*ctx, log.KV("MID_DAILY_BF", err))
		} else {
			result = append(result, bf)
		}
		// lastday part
		key = fmt.Sprintf(_baseBfKeyDaily, buvid, lastDt)
		bf, err = s.dao.LoadBloomFilter(ctx, key)
		if err != nil {
			log.Errorv(*ctx, log.KV("MID_DAILY_BF", err))
		} else {
			result = append(result, bf)
		}
	}

	return
}
