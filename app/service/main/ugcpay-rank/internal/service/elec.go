package service

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"go-common/app/service/main/ugcpay-rank/internal/conf"
	"go-common/app/service/main/ugcpay-rank/internal/model"
	"go-common/app/service/main/ugcpay-rank/internal/service/rank"
	"go-common/library/log"

	"github.com/pkg/errors"
)

func sfElecRankAV(avID, ver int64, size int) string {
	return fmt.Sprintf("sf_era_%d_%d_%d", avID, ver, size)
}

func sfElecRankUP(upMID, ver int64, size int) string {
	return fmt.Sprintf("sf_eru_%d_%d_%d", upMID, ver, size)
}

// ElecTotalRankAV 历史av充电榜单
func (s *Service) ElecTotalRankAV(ctx context.Context, upMID int64, avID int64, rankSize int) (res *model.RankElecAVProto, err error) {
	if res, err = s.elecRankAV(ctx, upMID, avID, conf.Conf.Biz.ElecAVRankSize, 0); err != nil {
		return
	}
	if len(res.List) <= rankSize {
		return
	}
	res.List = res.List[:rankSize]
	return
}

// ElecMonthRankAV 本月av充电榜单
func (s *Service) ElecMonthRankAV(ctx context.Context, upMID int64, avID int64, rankSize int) (res *model.RankElecAVProto, err error) {
	if res, err = s.elecRankAV(ctx, upMID, avID, conf.Conf.Biz.ElecAVRankSize, rank.MonthVer(time.Now())); err != nil {
		return
	}
	if len(res.List) <= rankSize {
		return
	}
	res.List = res.List[:rankSize]
	return
}

func (s *Service) elecRankAV(ctx context.Context, upMID int64, avid int64, rankSize int, ver int64) (res *model.RankElecAVProto, err error) {
	var (
		pr   = rank.NewElecPrepAVRank(avid, rankSize, ver, s.Dao)
		r    = rank.NewElecAVRank(avid, rankSize, ver, s.elecRankAVStorage(avid, ver), s.userSetting(ctx, upMID), s.Dao)
		data interface{}
		ok   bool
	)
	if data, err = s.rank(ctx, sfElecRankAV(avid, ver, rankSize), r, pr); err != nil {
		return
	}
	res, ok = data.(*model.RankElecAVProto)
	if !ok {
		err = errors.Errorf("ElecRankAV convert data: %T %+v to *model.ElecAVRank failed", data, data)
		return
	}
	return
}

// ElecMonthRankUP 本月up充电榜单
func (s *Service) ElecMonthRankUP(ctx context.Context, upMID int64, rankSize int) (r *model.RankElecUPProto, err error) {
	if r, err = s.elecRankUP(ctx, upMID, conf.Conf.Biz.ElecUPRankSize, rank.MonthVer(time.Now())); err != nil {
		return
	}
	if len(r.List) <= rankSize {
		return
	}
	r.List = r.List[:rankSize]
	return
}

func (s *Service) elecRankUP(ctx context.Context, upMID int64, rankSize int, ver int64) (res *model.RankElecUPProto, err error) {
	var (
		pr   = rank.NewElecPrepUPRank(upMID, rankSize, ver, s.Dao)
		r    = rank.NewElecUPRank(upMID, rankSize, ver, s.elecRankUPStorage(upMID, ver), s.userSetting(ctx, upMID), s.Dao)
		data interface{}
		ok   bool
	)
	if data, err = s.rank(ctx, sfElecRankUP(upMID, ver, rankSize), r, pr); err != nil {
		return
	}
	res, ok = data.(*model.RankElecUPProto)
	if !ok {
		err = errors.Errorf("ElecRankUP convert data: %T %+v to *model.ElecUPRank failed", data, data)
		return
	}
	return
}

func (s *Service) elecRankUPStorage(upMID, ver int64) (storage rank.Storager) {
	ramFlag := false
	for _, id := range conf.Conf.Biz.RAMUPIDs {
		if id == upMID {
			ramFlag = true
			break
		}
	}
	if ramFlag {
		log.Info("elecRankUPStorage choose RAM_storage, upMID: %d,ver: %d", upMID, ver)
		return rank.NewElecUPRankRAMStorage(upMID, ver, s.Dao)
	}
	return rank.NewElecUPRankMCStorage(upMID, ver, s.Dao)
}

func (s *Service) elecRankAVStorage(avID int64, ver int64) (storage rank.Storager) {
	ramFlag := false
	for _, id := range conf.Conf.Biz.RAMAVIDs {
		if id == avID {
			ramFlag = true
			break
		}
	}
	if ramFlag {
		log.Info("elecRankAVStorage choose RAM_storage, avID: %d,ver: %d", avID, ver)
		return rank.NewElecAVRankRAMStorage(avID, ver, s.Dao)
	}
	return rank.NewElecAVRankMCStorage(avID, ver, s.Dao)
}

func (s *Service) userSetting(ctx context.Context, mid int64) (setting model.ElecUserSetting) {
	defer func() {
		log.Info("userSetting mid: %d, value: %d", mid, setting)
	}()
	// load from local cache, 类型转换报错依赖panic机制
	userSettings := s.ElecUserSettings.Load().(*sync.Map)
	if settingIF, ok := userSettings.Load(mid); ok {
		setting = settingIF.(model.ElecUserSetting)
		return
	}
	setting = model.ElecUserSetting(math.MaxInt32)
	return
}
