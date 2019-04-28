package rank

import (
	"context"

	"go-common/app/service/main/ugcpay-rank/internal/dao"
	"go-common/app/service/main/ugcpay-rank/internal/model"

	"github.com/pkg/errors"
)

// NewElecAVRankMCStorage .
func NewElecAVRankMCStorage(avID int64, ver int64, dao *dao.Dao) (e *ElecAVRankMCStorage) {
	return &ElecAVRankMCStorage{
		avID: avID,
		ver:  ver,
		dao:  dao,
	}
}

// ElecAVRankMCStorage memcache storage
type ElecAVRankMCStorage struct {
	avID int64
	ver  int64
	dao  *dao.Dao
}

// Load 从cache中加载
func (e *ElecAVRankMCStorage) Load(ctx context.Context) (rank interface{}, err error) {
	var data *model.RankElecAVProto
	if data, err = e.dao.CacheElecAVRank(ctx, e.avID, e.ver); err != nil {
		return
	}
	if data != nil {
		rank = data
	}
	return
}

// Save 存储到cache
func (e *ElecAVRankMCStorage) Save(ctx context.Context, rank interface{}) (err error) {
	r, ok := rank.(*model.RankElecAVProto)
	if !ok {
		err = errors.Errorf("rank: %T %+v, can not convert to type: *model.ElecAVRank", rank, rank)
		return
	}
	err = e.dao.SetCacheElecAVRank(ctx, e.avID, e.ver, r)
	return
}

// NewElecAVRankRAMStorage .
func NewElecAVRankRAMStorage(avID int64, ver int64, dao *dao.Dao) (e *ElecAVRankRAMStorage) {
	return &ElecAVRankRAMStorage{
		avID: avID,
		ver:  ver,
		dao:  dao,
	}
}

// ElecAVRankRAMStorage ram storage
type ElecAVRankRAMStorage struct {
	avID int64
	ver  int64
	dao  *dao.Dao
}

// Load 从RAM中加载
func (e *ElecAVRankRAMStorage) Load(ctx context.Context) (rank interface{}, err error) {
	var data *model.RankElecAVProto
	if data, err = e.dao.LCLoadElecAVRank(e.avID, e.ver); err != nil {
		return
	}
	// log.Info("ElecAVRankRAMStorage load rank: %+v, avID: %d, ver: %d", data, e.avID, e.ver)
	if data != nil {
		rank = data
	}
	return
}

// Save 存储到RAM
func (e *ElecAVRankRAMStorage) Save(ctx context.Context, rank interface{}) (err error) {
	r, ok := rank.(*model.RankElecAVProto)
	if !ok {
		err = errors.Errorf("rank: %T %+v, can not convert to type: *model.ElecAVRank", rank, rank)
		return
	}
	// log.Info("ElecAVRankRAMStorage save rank: %+v", rank)
	err = e.dao.LCStoreElecAVRank(e.avID, e.ver, r)
	return
}
