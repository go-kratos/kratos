package dao

import (
	"time"

	"go-common/app/service/main/ugcpay-rank/internal/conf"
	"go-common/app/service/main/ugcpay-rank/internal/model"

	"github.com/bluele/gcache"
	"github.com/pkg/errors"
)

// LCStoreElecUPRank .
func (d *Dao) LCStoreElecUPRank(upMID, ver int64, rank *model.RankElecUPProto) (err error) {
	key := elecUPRankKey(upMID, ver)
	if err = d.elecUPRankLC.SetWithExpire(key, rank, time.Duration(conf.Conf.LocalCache.ElecUPRankTTL)); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// LCStoreElecAVRank .
func (d *Dao) LCStoreElecAVRank(avID, ver int64, rank *model.RankElecAVProto) (err error) {
	key := elecAVRankKey(avID, ver)
	if err = d.elecAVRankLC.SetWithExpire(key, rank, time.Duration(conf.Conf.LocalCache.ElecAVRankTTL)); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// LCLoadElecUPRank .
func (d *Dao) LCLoadElecUPRank(upMID, ver int64) (rank *model.RankElecUPProto, err error) {
	key := elecUPRankKey(upMID, ver)
	item, err := d.elecUPRankLC.Get(key)
	if err != nil {
		if err == gcache.KeyNotFoundError {
			err = nil
			rank = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	rank = item.(*model.RankElecUPProto)
	return
}

// LCLoadElecAVRank .
func (d *Dao) LCLoadElecAVRank(avID, ver int64) (rank *model.RankElecAVProto, err error) {
	key := elecAVRankKey(avID, ver)
	item, err := d.elecAVRankLC.Get(key)
	if err != nil {
		if err == gcache.KeyNotFoundError {
			err = nil
			rank = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	rank = item.(*model.RankElecAVProto)
	return
}
