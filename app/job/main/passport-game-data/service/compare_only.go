package service

import (
	"bufio"
	"context"
	"io"
	"os"
	"strconv"
	"time"

	"go-common/app/job/main/passport-game-data/conf"
	"go-common/app/job/main/passport-game-data/dao"
	"go-common/library/log"
)

// NewCompareOnly new a service for compare only.
func NewCompareOnly(c *conf.Config) (s *Service) {
	s = &Service{
		c:    c,
		d:    dao.New(c),
		l2cC: newCompareConfigFrom(c.Compare.Local2Cloud),
	}
	return
}

// CompareFromMidListFile load mid list from file and compare.
func (s *Service) CompareFromMidListFile(c context.Context, fn string) (err error) {
	var f *os.File
	if f, err = os.Open(fn); err != nil {
		log.Error("failed to open file %s, error(%v)", fn, err)
		return
	}
	defer f.Close()

	cc := s.l2cC

	rd := bufio.NewReader(f)
	skippedCount := 0
	var (
		mid      int64
		mids     = make([]int64, 0)
		line     []byte
		isPrefix bool
	)
	for {
		line, isPrefix, err = rd.ReadLine()

		if isPrefix || err != nil || err == io.EOF {
			break
		}
		mid, err = strconv.ParseInt(string(line), 10, 64)
		if err != nil {
			log.Error("failed to parse mid, strconv.ParseInt(%s, 10 ,64) error(%v), skip", line, err)
			skippedCount++
			continue
		}
		mids = append(mids, mid)
	}
	log.Info("mid list len: %d, total skipped count: %d", len(mids), skippedCount)
	if len(mids) == 0 {
		return
	}

	for {
		time.Sleep(cc.LoopDuration)
		if err = s.local2CloudCompare(context.TODO(), s.batchQueryLocalNonMiss(context.TODO(), mids, cc.BatchSize, cc.BatchMissRetryCount)); err == nil {
			break
		}
		log.Error("failed to compare mids %v, retrying", mids)
	}
	return
}
