package rank

import (
	"context"
	"fmt"

	accmdl "go-common/app/service/main/account/model"
	"go-common/app/service/main/ugcpay-rank/internal/dao"
	"go-common/app/service/main/ugcpay-rank/internal/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// NewElecPrepUPRank .
func NewElecPrepUPRank(upMID int64, size int, ver int64, d *dao.Dao) *ElecPrepUPRank {
	return &ElecPrepUPRank{
		upMID: upMID,
		size:  size,
		ver:   ver,
		dao:   d,
	}
}

// ElecPrepUPRank 充电up预备榜单
type ElecPrepUPRank struct {
	upMID int64
	size  int
	ver   int64
	dao   *dao.Dao
}

func (e *ElecPrepUPRank) String() string {
	return fmt.Sprintf("ElecPrepUPRank up_mid: %d, size: %d", e.upMID, e.size)
}

// Load 从cache中加载
func (e *ElecPrepUPRank) Load(ctx context.Context) (rank interface{}, err error) {
	var data *model.RankElecPrepUPProto
	if data, _, err = e.dao.CacheElecPrepUPRank(ctx, e.upMID, e.ver); err != nil {
		return
	}
	if data != nil {
		rank = data
	}
	return
}

// Save 存储到cache
func (e *ElecPrepUPRank) Save(ctx context.Context, rank interface{}) (err error) {
	r, ok := rank.(*model.RankElecPrepUPProto)
	if !ok {
		err = errors.Errorf("rank: %T %+v, can not convert to type: *model.ElecPrepUPRank", rank, rank)
		return
	}

	casFN := func() (fok bool, ferr error) {
		var (
			item *memcache.Item
			data *model.RankElecPrepUPProto
		)
		if fok, ferr = e.dao.AddCacheElecPrepUPRank(ctx, e.upMID, e.ver, r); ferr != nil {
			return
		}
		if fok {
			return
		}
		if data, item, ferr = e.dao.CacheElecPrepUPRank(ctx, e.upMID, e.ver); ferr != nil {
			return
		}
		if data == nil {
			fok = false
			return
		}
		if fok, ferr = e.dao.CASCacheElecPrepRank(ctx, data, item); ferr != nil {
			return
		}
		return
	}

	err = tryHard(casFN, "ElecPrepUPRank:CAS", 3)
	return
}

// Rebuild 从db重构
func (e *ElecPrepUPRank) Rebuild(ctx context.Context) (rank interface{}, err error) {
	var (
		theRank = &model.RankElecPrepUPProto{
			UPMID: e.upMID,
			Size_: e.size,
		}
		dbData  []*model.DBElecUPRank
		payMIDs []int64
	)
	if theRank.Count, err = e.dao.RawCountElecUPRank(ctx, e.upMID, e.ver); err != nil {
		return
	}
	if dbData, err = e.dao.RawElecUPRankList(ctx, e.upMID, e.ver, e.size); err != nil {
		return
	}
	for i, d := range dbData {
		ele := &model.RankElecPrepElementProto{
			MID:       d.PayMID,
			Rank:      i + 1,
			TrendType: model.TrendHold,
			Amount:    d.PayAmount,
		}
		payMIDs = append(payMIDs, d.PayMID)
		theRank.List = append(theRank.List, ele)
	}
	// 填充充电总人数
	if theRank.UPMID != 0 {
		if theRank.CountUPTotalElec, err = e.dao.RawCountUPTotalElec(ctx, theRank.UPMID); err != nil {
			err = nil
			theRank.CountUPTotalElec = 0
			log.Error("e.dao.RawCountUPTotalElec upMID: %d, err: %+v", theRank.UPMID, err)
		}
	}
	log.Info("Rebuild ElecPrepUPRank upMID: %d, ver: %d, count: %d, upTotalCount: %d", theRank.UPMID, e.ver, theRank.Count, theRank.CountUPTotalElec)
	// 填充留言信息
	if e.ver != 0 {
		var (
			messageMap map[int64]*model.DBElecMessage
		)
		if messageMap, err = e.dao.RawElecUPMessages(ctx, payMIDs, e.upMID, e.ver); err != nil {
			return
		}
		for _, r := range theRank.List {
			msg, ok := messageMap[r.MID]
			if ok {
				r.Message = &model.ElecMessageProto{
					Message: msg.Message,
					Hidden:  msg.Hidden,
				}
			}
		}
	}

	rank = theRank
	return
}

// UpdateOrder 在内存态通过订单更新并返回
func (e *ElecPrepUPRank) UpdateOrder(ctx context.Context, raw interface{}, payMID int64, fee int64) (res interface{}, err error) {
	r, ok := raw.(*model.RankElecPrepUPProto)
	if !ok {
		err = errors.Errorf("rank: %T %+v, can not convert to type: *model.ElecPrepUPRank", raw, raw)
		return
	}

	var (
		elecerRank *model.DBElecUPRank
		payAmount  = fee
	)
	if elecerRank, err = e.dao.RawElecUPRank(ctx, e.upMID, e.ver, payMID); err != nil {
		log.Error("e.dao.RawElecUPRank upMID: %d, ver: %d, payMID: %d, err: %+v", e.upMID, e.ver, payMID, err)
		err = nil
	}
	log.Info("ElecPrepUPRank: %s, update elecerRank: %+v, from payMID : %d, fee: %d", e, elecerRank, payMID, fee)
	if elecerRank != nil {
		payAmount = elecerRank.PayAmount
	}
	r.Charge(payMID, payAmount, payAmount == fee)
	log.Info("charge ElecPrepUPRank: payMID: %d, fee: %d, payAmount: %d, isNew: %t", payMID, fee, payAmount, payAmount == fee)
	res = r
	return
}

// UpdateMessage 在内存态通过留言更新并返回
func (e *ElecPrepUPRank) UpdateMessage(ctx context.Context, raw interface{}, payMID int64, message string, hidden bool) (res interface{}, err error) {

	r, ok := raw.(*model.RankElecPrepUPProto)
	if !ok {
		err = errors.Errorf("rank: %T %+v, can not convert to type: *model.ElecPrepUPRank", raw, raw)
		return
	}
	r.UpdateMessage(payMID, message, hidden)

	res = r
	return
}

// NewElecUPRank .
func NewElecUPRank(upMID int64, size int, ver int64, storage Storager, userSetting model.ElecUserSetting, d *dao.Dao) *ElecUPRank {
	ret := &ElecUPRank{
		upMID:       upMID,
		size:        size,
		ver:         ver,
		userSetting: userSetting,
		dao:         d,
	}
	ret.Storager = storage
	return ret
}

// ElecUPRank 充电up正式榜单
type ElecUPRank struct {
	Storager
	upMID       int64
	size        int
	ver         int64
	userSetting model.ElecUserSetting
	dao         *dao.Dao
}

func (e *ElecUPRank) String() string {
	return fmt.Sprintf("ElecUPRank up_mid: %d, size: %d", e.upMID, e.size)
}

// Rebuild 从预备榜单重构
func (e *ElecUPRank) Rebuild(ctx context.Context, prepRank interface{}) (rank interface{}, err error) {
	pr, ok := prepRank.(*model.RankElecPrepUPProto)
	if !ok {
		err = errors.Errorf("prepRank: %T %+v, can not convert to type: *model.ElecPrepUPRank", prepRank, prepRank)
		return
	}

	// 从 prepRank 填充数据
	theRank := &model.RankElecUPProto{
		CountUPTotalElec: pr.CountUPTotalElec,
		Count:            pr.Count,
		UPMID:            pr.UPMID,
		Size_:            e.size,
	}
	rank = theRank
	var (
		mids       = make([]int64, 0)
		accountMap map[int64]*accmdl.Card
	)
	for _, r := range pr.List {
		if r == nil {
			continue
		}
		mids = append(mids, r.MID)
		rankEle := &model.RankElecElementProto{
			RankElecPrepElementProto: *r,
		}
		// ^(用户设置允许展示留言 && top3用户留言)
		if !e.userSetting.ShowMessage() || r.Rank > 3 {
			log.Info("ElecUPRank add message, mid: %d, show_message: %t, rank: %d", e.upMID, e.userSetting.ShowMessage(), r.Rank)
			rankEle.Message = nil
		}
		theRank.List = append(theRank.List, rankEle)
	}
	if len(theRank.List) <= 0 {
		return
	}

	// 填充会员信息
	if accountMap, err = e.dao.AccountCards(ctx, mids); err != nil {
		log.Error("e.dao.AccountCards mids: %+v, err: %+v", mids, err)
		err = nil
	}
	for _, r := range theRank.List {
		card, ok := accountMap[r.MID]
		if ok {
			r.Nickname = card.Name
			r.Avatar = card.Face
			r.VIP = &model.VIPInfoProto{
				Type:   card.Vip.Type,
				Status: card.Vip.Status,
				// DueDate: card.Vip.DueDate,
			}
		}
	}
	return
}
