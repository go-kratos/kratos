package service

import (
	"context"

	"go-common/app/job/main/passport-game-data/model"
	"go-common/library/log"
)

func (s *Service) batchQueryLocalNonMiss(c context.Context, mids []int64, batchSize, batchMissRetryCount int) (res []*model.OriginAsoAccount) {
	as, miss := s.batchQueryLocalByMid(c, mids, batchSize)
	if len(miss) == 0 {
		return as
	}

	res = as
	for i := 0; i < batchMissRetryCount; i++ {
		log.Info("try for the %dth retry, miss mids: %v", miss)
		as, miss = s.batchQueryLocalByMid(c, miss, batchSize)

		res = append(res, as...)

		if len(miss) == 0 {
			return
		}

		if i == batchMissRetryCount-1 {
			log.Error("still miss those mids: %v after %d tries", miss, batchMissRetryCount)
		}
	}
	return
}

func (s *Service) batchQueryLocalByMid(c context.Context, mids []int64, batchSize int) (res []*model.OriginAsoAccount, miss []int64) {
	if len(mids) == 0 {
		return
	}
	res = make([]*model.OriginAsoAccount, 0)
	miss = make([]int64, 0)
	bc := len(mids)/batchSize + 1
	for i := 0; i < bc; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > len(mids) {
			end = len(mids)
		}
		partMids := mids[start:end]
		as, err := s.d.AsoAccountsLocal(c, partMids)
		if err != nil {
			miss = append(miss, partMids...)
			continue
		}
		res = append(res, as...)
	}
	return
}

func (s *Service) batchQueryCloudNonMiss(c context.Context, mids []int64, batchSize, batchMissRetryCount int) (res []*model.AsoAccount) {
	if len(mids) == 0 {
		return
	}
	as, miss := s.batchQueryCloudByMid(c, mids, batchSize)
	if len(miss) == 0 {
		return as
	}

	res = as
	for i := 0; i < batchMissRetryCount; i++ {
		log.Info("try for the %dth times, miss mids: %v", i, miss)
		as, miss = s.batchQueryCloudByMid(c, miss, batchSize)

		res = append(res, as...)
		if len(miss) == 0 {
			return
		}

		if i == batchMissRetryCount-1 {
			log.Error("still miss those mids: %v after %d tries", miss, batchMissRetryCount)
		}
	}
	return
}

func (s *Service) batchQueryCloudByMid(c context.Context, mids []int64, batchSize int) (res []*model.AsoAccount, miss []int64) {
	if len(mids) == 0 {
		return
	}
	res = make([]*model.AsoAccount, 0)
	miss = make([]int64, 0)
	bc := len(mids)/batchSize + 1
	for i := 0; i < bc; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > len(mids) {
			end = len(mids)
		}
		partMids := mids[start:end]
		as, err := s.d.AsoAccountsCloud(c, partMids)
		if err != nil {
			miss = append(miss, partMids...)
			continue
		}
		res = append(res, as...)
	}
	return
}
