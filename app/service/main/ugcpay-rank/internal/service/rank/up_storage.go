package rank

import (
	"context"

	"go-common/app/service/main/ugcpay-rank/internal/dao"
	"go-common/app/service/main/ugcpay-rank/internal/model"

	"github.com/pkg/errors"
)

// NewElecUPRankMCStorage .
func NewElecUPRankMCStorage(upMID int64, ver int64, dao *dao.Dao) (e *ElecUPRankMCStorage) {
	return &ElecUPRankMCStorage{
		upMID: upMID,
		ver:   ver,
		dao:   dao,
	}
}

// ElecUPRankMCStorage .
type ElecUPRankMCStorage struct {
	upMID int64
	ver   int64
	dao   *dao.Dao
}

// Load 从cache中加载
func (e *ElecUPRankMCStorage) Load(ctx context.Context) (rank interface{}, err error) {
	var data *model.RankElecUPProto
	if data, err = e.dao.CacheElecUPRank(ctx, e.upMID, e.ver); err != nil {
		return
	}
	if data != nil {
		rank = data
	}
	return
}

// Save 存储到cache
func (e *ElecUPRankMCStorage) Save(ctx context.Context, rank interface{}) (err error) {
	r, ok := rank.(*model.RankElecUPProto)
	if !ok {
		err = errors.Errorf("rank: %T %+v, can not convert to type: *model.ElecUPRank", rank, rank)
		return
	}
	err = e.dao.SetCacheElecUPRank(ctx, e.upMID, e.ver, r)
	return
}

// NewElecUPRankRAMStorage .
func NewElecUPRankRAMStorage(upMID int64, ver int64, dao *dao.Dao) (e *ElecUPRankRAMStorage) {
	return &ElecUPRankRAMStorage{
		upMID: upMID,
		ver:   ver,
		dao:   dao,
	}
}

// ElecUPRankRAMStorage .
type ElecUPRankRAMStorage struct {
	upMID int64
	ver   int64
	dao   *dao.Dao
}

// Load 从ram中加载
func (e *ElecUPRankRAMStorage) Load(ctx context.Context) (rank interface{}, err error) {
	var data *model.RankElecUPProto
	if data, err = e.dao.LCLoadElecUPRank(e.upMID, e.ver); err != nil {
		return
	}
	// log.Info("ElecUPRankRAMStorage load rank: %+v, upMID: %d, ver: %d", data, e.upMID, e.ver)
	if data != nil {
		rank = data
	}
	return
}

// Save 存储到ram
func (e *ElecUPRankRAMStorage) Save(ctx context.Context, rank interface{}) (err error) {
	r, ok := rank.(*model.RankElecUPProto)
	if !ok {
		err = errors.Errorf("rank: %T %+v, can not convert to type: *model.ElecUPRank", rank, rank)
		return
	}
	// log.Info("ElecUPRankRAMStorage save rank: %+v", rank)
	err = e.dao.LCStoreElecUPRank(e.upMID, e.ver, r)
	return
}
