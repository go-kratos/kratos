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

// NewElecPrepAVRank .
func NewElecPrepAVRank(avID int64, size int, ver int64, d *dao.Dao) *ElecPrepAVRank {
	return &ElecPrepAVRank{
		avID: avID,
		size: size,
		ver:  ver,
		dao:  d,
	}
}

// ElecPrepAVRank 充电av预备榜单
type ElecPrepAVRank struct {
	avID int64
	size int
	ver  int64
	dao  *dao.Dao
}

func (e *ElecPrepAVRank) String() string {
	return fmt.Sprintf("ElecPrepAVRank avID: %d, size: %d, ver: %d", e.avID, e.size, e.ver)
}

// Load 从cache中加载
func (e *ElecPrepAVRank) Load(ctx context.Context) (rank interface{}, err error) {
	var data *model.RankElecPrepAVProto
	if data, _, err = e.dao.CacheElecPrepAVRank(ctx, e.avID, e.ver); err != nil {
		return
	}
	if data != nil {
		rank = data
	}
	return
}

// Save 存储到cache
func (e *ElecPrepAVRank) Save(ctx context.Context, rank interface{}) (err error) {
	r, ok := rank.(*model.RankElecPrepAVProto)
	if !ok {
		err = errors.Errorf("rank: %T %+v, can not convert to type: *model.ElecPrepAVRank", rank, rank)
		return
	}

	casFN := func() (fok bool, ferr error) {
		var (
			item *memcache.Item
			data *model.RankElecPrepAVProto
		)
		if fok, ferr = e.dao.AddCacheElecPrepAVRank(ctx, e.avID, e.ver, r); ferr != nil {
			return
		}
		if fok {
			return
		}
		if data, item, ferr = e.dao.CacheElecPrepAVRank(ctx, e.avID, e.ver); ferr != nil {
			return
		}
		if data == nil {
			fok = false
			return
		}
		if fok, ferr = e.dao.CASCacheElecPrepRank(ctx, r, item); ferr != nil {
			return
		}
		return
	}

	err = tryHard(casFN, "ElecPrepAVRank:CAS", 3)
	return
}

// Rebuild 从db重构
func (e *ElecPrepAVRank) Rebuild(ctx context.Context) (rank interface{}, err error) {
	var (
		theRank = &model.RankElecPrepAVProto{
			AVID: e.avID,
		}
		dbData  []*model.DBElecAVRank
		payMIDs []int64
	)
	theRank.Size_ = e.size
	if theRank.Count, err = e.dao.RawCountElecAVRank(ctx, e.avID, e.ver); err != nil {
		return
	}
	if dbData, err = e.dao.RawElecAVRankList(ctx, e.avID, e.ver, e.size); err != nil {
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
		theRank.UPMID = d.UPMID
	}
	// 填充充电总人数
	if theRank.UPMID != 0 {
		if theRank.CountUPTotalElec, err = e.dao.RawCountUPTotalElec(ctx, theRank.UPMID); err != nil {
			err = nil
			theRank.CountUPTotalElec = 0
			log.Error("e.dao.RawCountUPTotalElec upMID: %d, err: %+v", theRank.UPMID, err)
		}
	}
	log.Info("Rebuild ElecPrepAVRank avID: %d, upMID: %d, ver: %d, count: %d, upTotalCount: %d", e.avID, theRank.UPMID, e.ver, theRank.Count, theRank.CountUPTotalElec)
	var (
		messageMap map[int64]*model.DBElecMessage
	)
	if e.ver == 0 {
		if messageMap, err = e.dao.RawElecAVMessages(ctx, payMIDs, e.avID); err != nil {
			return
		}
	} else {
		if messageMap, err = e.dao.RawElecAVMessagesByVer(ctx, payMIDs, e.avID, e.ver); err != nil {
			return
		}
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

	rank = theRank
	return
}

// UpdateOrder 在内存态通过订单更新并返回
func (e *ElecPrepAVRank) UpdateOrder(ctx context.Context, raw interface{}, payMID int64, fee int64) (res interface{}, err error) {
	r, ok := raw.(*model.RankElecPrepAVProto)
	if !ok {
		err = errors.Errorf("rank: %T %+v, can not convert to type: *model.ElecPrepAVRank", raw, raw)
		return
	}

	var (
		elecerRank *model.DBElecAVRank
		payAmount  = fee
	)
	if elecerRank, err = e.dao.RawElecAVRank(ctx, e.avID, e.ver, payMID); err != nil {
		log.Error("e.dao.RawElecUPRank avID: %d, ver: %d, payMID: %d, err: %+v", e.avID, e.ver, payMID, err)
		err = nil
	}
	log.Info("ElecPrepAVRank: %s, update elecerRank: %+v, from pay_mid: %d, fee: %d", e, elecerRank, payMID, fee)
	if elecerRank != nil {
		payAmount = elecerRank.PayAmount
	}
	r.Charge(payMID, payAmount, payAmount == fee)
	log.Info("charge ElecPrepAVRank: payMID: %d, fee: %d, payAmount: %d, isNew: %t", payMID, fee, payAmount, payAmount == fee)

	res = r
	return
}

// UpdateMessage 在内存态通过留言更新并返回
func (e *ElecPrepAVRank) UpdateMessage(ctx context.Context, raw interface{}, payMID int64, message string, hidden bool) (res interface{}, err error) {
	r, ok := raw.(*model.RankElecPrepAVProto)
	if !ok {
		err = errors.Errorf("rank: %T %+v, can not convert to type: *model.ElecPrepAVRank", raw, raw)
		return
	}

	r.UpdateMessage(payMID, message, hidden)

	res = r
	return
}

// NewElecAVRank .
func NewElecAVRank(avID int64, size int, ver int64, storage Storager, userSetting model.ElecUserSetting, d *dao.Dao) *ElecAVRank {
	ret := &ElecAVRank{
		avID:        avID,
		size:        size,
		ver:         ver,
		userSetting: userSetting,
		dao:         d,
	}
	ret.Storager = storage
	return ret
}

// ElecAVRank 充电av正式榜单
type ElecAVRank struct {
	Storager
	avID        int64
	size        int
	ver         int64
	userSetting model.ElecUserSetting
	dao         *dao.Dao
}

func (e *ElecAVRank) String() string {
	return fmt.Sprintf("ElecAVRank avID: %d, size: %d", e.avID, e.size)
}

// Rebuild 从预备榜单重构
func (e *ElecAVRank) Rebuild(ctx context.Context, prepRank interface{}) (rank interface{}, err error) {
	pr, ok := prepRank.(*model.RankElecPrepAVProto)
	if !ok {
		err = errors.Errorf("prepRank: %T %+v, can not convert to type: *model.ElecAVRank", prepRank, prepRank)
		return
	}

	// 从 prepRank 填充数据
	theRank := &model.RankElecAVProto{
		CountUPTotalElec: pr.CountUPTotalElec,
		Count:            pr.Count,
		UPMID:            pr.UPMID,
		AVID:             pr.AVID,
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
			log.Info("ElecAVRank add message, show_message: %t, rank: %d", e.userSetting.ShowMessage(), r.Rank)
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
