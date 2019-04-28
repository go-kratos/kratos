package service

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"go-common/app/job/main/passport-game-data/model"
	"go-common/library/log"
)

// cloud2localcompareproc compare aso accounts between cloud and local.
// select last modified from cloud
// load batchSize from origin
// compare:
// if cloud_mtime >= local_mtime, directly compare
// if cloud_mtime < local_mtime, sleep and reload from cloud, then do compare again
func (s *Service) cloud2localcompareproc() {
	var (
		err      error
		cloudRes []*model.AsoAccount

		ack = false
	)

	cc := s.c2lC

	delay := cc.DelayDuration

	cc.st = cc.StartTime
	cc.ed = cc.st.Add(cc.StepDuration)

	offsetFile, err := os.Create(cc.OffsetFilePath)
	if err != nil {
		log.Error("failed to create offset file, os.Create(%s) error(%v)", cc.OffsetFilePath, err)
		return
	}
	defer offsetFile.Close()
	log.Info("created offset file %s", cc.OffsetFilePath)

	for {
		time.Sleep(cc.LoopDuration)

		cc.sleeping = false

		if ack {
			cc.st = cc.st.Add(cc.StepDuration)
			cc.ed = cc.ed.Add(cc.StepDuration)
		}

		st, ed := cc.st, cc.ed

		if err = ioutil.WriteFile(cc.OffsetFilePath, []byte(st.Format(_timeFormat)), os.ModeAppend); err != nil {
			log.Error("failed to write offset, ioutil.WriteFile(%s, %s, os.ModeAppend), error(%v)", cc.OffsetFilePath, st.Format(_timeFormat), err)
			continue
		}

		if cc.Debug {
			log.Info("st: %s, ed: %s", st.Format(_timeFormat), ed.Format(_timeFormat))
		}

		if cc.End && st.After(cc.EndTime) {
			log.Info("st:%s is after endTime:%s, all data compares ok, cloud2localcompareproc exit", st.Format(_timeFormat), cc.EndTime.Format(_timeFormat))
			return
		}

		now := time.Now()

		if now.Sub(st) <= delay {
			delta := int64(delay/time.Second) - (now.Unix() - st.Unix())
			log.Info("now time is just after st by %d seconds, not greater than delay duration: %v, will sleep %d seconds", int64(delay/time.Second)-delta, delay, delta)

			cc.sleeping = true
			cc.sleepingSeconds = delta
			cc.sleepFromTs = now.Unix()

			time.Sleep(time.Duration(int64(time.Second) * delta))
			continue
		}

		if cloudRes, err = s.d.AsoAccountRangeCloud(context.TODO(), st, ed); err != nil {
			continue
		}
		cc.rangeCount = len(cloudRes)
		cc.totalCount += len(cloudRes)

		if err = s.cloud2LocalCompare(context.TODO(), cloudRes); err != nil {
			continue
		}

		ack = true
	}
}

func (s *Service) cloud2LocalCompare(c context.Context, cloudRes []*model.AsoAccount) (err error) {
	mids := make([]int64, 0)
	for _, item := range cloudRes {
		mids = append(mids, item.Mid)
	}

	cc := s.c2lC

	localRes := s.batchQueryLocalNonMiss(context.TODO(), mids, cc.BatchSize, cc.BatchMissRetryCount)

	m := make(map[int64]*model.OriginAsoAccount)
	for _, item := range localRes {
		m[item.Mid] = item
	}

	// compare
	pendingMids := make([]int64, 0)
	for _, item := range cloudRes {
		cloud := item
		local := m[item.Mid]
		status := doCompare(cloud, local, true)
		switch status {
		case _statusOK:
			// do nothing
		case _statusNo:
			cc.diffCount++
			s.doLog(cloud, local, false)
			if cc.Fix {
				s.fixCloudRecord(context.TODO(), model.Default(local), cloud)
			}
		case _statusPending:
			pendingMids = append(pendingMids, item.Mid)
		}
	}

	if len(pendingMids) == 0 {
		return
	}

	// reload pending mids from cloud
	var pendingRes []*model.AsoAccount
	if pendingRes, err = s.d.AsoAccountsCloud(context.TODO(), pendingMids); err != nil {
		return
	}
	// compare
	for _, item := range pendingRes {
		cloud := item
		local := m[item.Mid]
		status := doCompare(item, m[item.Mid], false)
		switch status {
		case _statusOK:
		case _statusNo:
			cc.diffCount++
			s.doLog(cloud, local, true)
			if cc.Fix {
				s.fixCloudRecord(context.TODO(), model.Default(local), cloud)
			}
		}
	}
	return
}
