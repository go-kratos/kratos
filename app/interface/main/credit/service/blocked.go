package service

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/credit/model"
	acmdl "go-common/app/service/main/account/api"
	blkmdl "go-common/app/service/main/member/model/block"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

var (
	_emptyAnnounce  = []*model.BlockedAnnouncement{}
	_emptyBlockInfo = []*model.BlockedInfo{}
)

// BlockedUserCard get user blocked card info.
func (s *Service) BlockedUserCard(c context.Context, mid int64) (card *model.BlockedUserCard, err error) {
	var (
		count   int
		profile *acmdl.ProfileReply
	)
	accArg := &acmdl.MidReq{Mid: mid}
	if profile, err = s.accountClient.Profile3(c, accArg); err != nil {
		err = errors.Wrap(err, "accountClient.Profile3")
		return
	}
	card = &model.BlockedUserCard{UID: mid, Uname: profile.Profile.Name, Face: profile.Profile.Face}
	if count, err = s.dao.BlockedCount(c, mid); err != nil {
		return
	}
	card.BlockedSum = count
	bms, err := s.memRPC.BlockInfo(c, &blkmdl.RPCArgInfo{MID: mid})
	if err != nil {
		return
	}
	card.BlockedStatus = int(bms.BlockStatus)
	status := int8(bms.BlockStatus)
	if status == model.BlockStatusForever || status == model.BlockStatusOn {
		card.BlockedStatus = 1
		// TODO nothing record in credit.
		if count <= 0 {
			card.BlockedSum = 1
		}
	}
	if status == model.BlockStatusForever {
		card.BlockedForever = model.BlockedStateForever
	}
	card.BlockedEndTime = bms.EndTime
	card.MoralNum = int(profile.Profile.Moral)
	if status != 0 {
		delta := time.Until(time.Unix(bms.EndTime, 0))
		if int(delta) > 0 {
			card.BlockedRestDay = int64(delta / (time.Hour * 24))
			if delta%(time.Hour*24) > 0 {
				card.BlockedRestDay++
			}
		}
		if card.AnsWerStatus, err = s.dao.AnswerStatus(c, mid, time.Unix(bms.StartTime, 0)); err != nil {
			return
		}
	}
	return
}

// BlockedUserList get user blocked info list
func (s *Service) BlockedUserList(c context.Context, mid int64) (r []*model.BlockedInfo, err error) {
	var mc = true
	if r, err = s.dao.BlockedUserListCache(c, mid); err != nil {
		err = nil
		mc = false
	}
	if len(r) > 0 {
		return
	}
	if r, err = s.dao.BlockedUserList(c, mid); err != nil {
		return
	}
	if mc && len(r) > 0 {
		s.addCache(func() {
			s.dao.SetBlockedUserListCache(context.TODO(), mid, r)
		})
	}
	return
}

// BlockedInfo blocked info
func (s *Service) BlockedInfo(c context.Context, id int64) (info *model.BlockedInfo, err error) {
	var mc = true
	if info, err = s.dao.BlockedInfoCache(c, id); err != nil {
		err = nil
		mc = false
	}
	if info != nil {
		if int8(info.PublishStatus) == model.PublishStatusClose {
			err = ecode.NothingFound
		}
		return
	}
	if info, err = s.dao.BlockedInfoByID(c, id); err != nil {
		err = errors.Wrapf(err, "BlockedInfoByID(%d)", id)
		return
	}
	if info == nil || (int8(info.PublishStatus) == model.PublishStatusClose) {
		err = ecode.NothingFound
		return
	}
	if mc {
		s.addBlockedCache(c, info)
	}
	return
}

// BlockedInfoAppeal get blocked info for appeal .
func (s *Service) BlockedInfoAppeal(c context.Context, id, mid int64) (info *model.BlockedInfo, err error) {
	defer func() {
		if err == nil && info != nil && info.ID != 0 {
			if mid != info.UID {
				err = ecode.NothingFound
			}
		}
	}()
	var mc = true
	if info, err = s.dao.BlockedInfoCache(c, id); err != nil {
		err = nil
		mc = false
	}
	if info != nil {
		return
	}
	if info, err = s.dao.BlockedInfoByID(c, id); err != nil {
		err = errors.Wrapf(err, "BlockedInfoByID(%d)", id)
		return
	}
	if info == nil {
		err = ecode.NothingFound
		return
	}
	if mc {
		s.addBlockedCache(c, info)
	}
	return
}

func (s *Service) addBlockedCache(c context.Context, info *model.BlockedInfo) {
	var card *acmdl.CardReply
	if card, _ = s.userInfo(c, info.UID); card != nil {
		info.Uname = card.Card.Name
		info.Face = card.Card.Face
	}
	info.Build()
	s.addCache(func() {
		s.dao.SetBlockedInfoCache(context.TODO(), info.ID, info)
	})
}

// BlockedList blocked info list, public default.
func (s *Service) BlockedList(c context.Context, oType, bType int8, pn, ps int) (res []*model.BlockedInfo, err error) {
	var (
		start  = (pn - 1) * ps
		end    = pn * ps
		ok     bool
		ids    []int64
		missed []int64
		cache  = true
	)
	if ok, err = s.dao.ExpireBlockedIdx(c, oType, bType); ok && err == nil {
		if ids, err = s.dao.BlockedIdxCache(c, oType, bType, start, end-1); err != nil {
			return
		}
	} else {
		var ls, tmpls []*model.BlockedInfo
		if ls, err = s.dao.BlockedList(c, oType, bType); err != nil {
			return
		}
		switch {
		case len(ls) <= int(start):
			tmpls = _emptyBlockInfo
		case len(ls) <= int(end):
			tmpls = ls[start:]
		default:
			tmpls = ls[start:end]
		}
		s.addCache(func() {
			s.dao.LoadBlockedIdx(context.TODO(), oType, bType, ls)
		})
		for _, id := range tmpls {
			ids = append(ids, id.ID)
		}
	}
	if res, missed, err = s.dao.BlockedInfosCache(c, ids); err != nil {
		err = nil
		cache = false
		missed = ids
	}
	var missInfos []*model.BlockedInfo
	if len(missed) != 0 {
		missInfos, err = s.dao.BlockedInfos(c, missed)
		if err != nil {
			return
		}
		res = append(res, missInfos...)
	}
	var (
		mids []int64
		oids []int64
	)
	for _, i := range res {
		mids = append(mids, i.UID)
		oids = append(oids, i.ID)
	}
	arg := &acmdl.MidsReq{
		Mids: mids,
	}
	cards, err := s.accountClient.Infos3(c, arg)
	if err != nil {
		err = errors.Wrap(err, "Infos")
		return
	}
	reply, _ := s.dao.ReplysCount(c, oids)
	for _, i := range res {
		if card, ok := cards.Infos[i.UID]; ok {
			i.Uname = card.Name
			i.Face = card.Face
		}
		i.CommentSum = reply[strconv.FormatInt(i.ID, 10)]
		i.Build()
	}
	if cache {
		s.addCache(func() {
			s.dao.SetBlockedInfosCache(context.TODO(), missInfos)
		})
	}
	return
}

// AnnouncementInfo get announcement detail.
func (s *Service) AnnouncementInfo(c context.Context, aid int64) (res *model.BlockedAnnouncement, err error) {
	var ok bool
	if res, ok = s.announcement.amap[aid]; !ok {
		err = ecode.NothingFound
	}
	return
}

// AnnouncementList get announcement list.
func (s *Service) AnnouncementList(c context.Context, tp int8, pn, ps int64) (resp *model.AnnounceList, err error) {
	var (
		ok    bool
		start = (pn - 1) * ps
		end   = pn * ps
		alist []*model.BlockedAnnouncement
		count int64
	)
	resp = &model.AnnounceList{
		List: _emptyAnnounce,
	}
	if tp == model.PublishTypedef {
		resp.List = s.announcement.def
		resp.Count = int64(len(s.announcement.def))
		return
	}
	if alist, ok = s.announcement.alist[tp]; !ok {
		return
	}
	count = int64(len(alist))
	resp.Count = count
	switch {
	case count < start:
	case end >= count:
		resp.List = alist[start:]
	default:
		resp.List = alist[start:end]
	}
	return
}

// LoadAnnouncement load AnnouncementList.
func (s *Service) LoadAnnouncement(c context.Context) {
	res, err := s.dao.AnnouncementList(c)
	if err != nil {
		return
	}
	var (
		def     []*model.BlockedAnnouncement
		new     = make([]*model.BlockedAnnouncement, 0, model.PublishInitLen)
		top     = make([]*model.BlockedAnnouncement, 0, model.PublishInitLen)
		alist   = make(map[int8][]*model.BlockedAnnouncement)
		topList = make(map[int8][]*model.BlockedAnnouncement)
		amap    = make(map[int64]*model.BlockedAnnouncement)
	)
	for _, ann := range res {
		if ann.StickStatus == 1 {
			top = append(top, ann)
			topList[ann.Ptype] = append(topList[ann.Ptype], ann)
		} else if len(new) < model.PublishInitLen {
			new = append(new, ann)
		}
		if ann.StickStatus != 1 {
			alist[ann.Ptype] = append(alist[ann.Ptype], ann)
		}
		amap[ann.ID] = ann
	}
	for t, p := range alist {
		alist[t] = append(topList[t], p...)
	}
	if len(top) < model.PublishInitLen {
		lack := model.PublishInitLen - len(top)
		if lack > len(new) {
			def = append(top, new...)
		} else {
			def = append(top, new[:lack]...)
		}
	} else {
		def = top[:model.PublishInitLen]
	}
	if len(def) == 0 {
		def = _emptyAnnounce
	}
	s.announcement.def = def
	s.announcement.alist = alist
	s.announcement.amap = amap
}

// BlockedNumUser get blocked user number.
func (s *Service) BlockedNumUser(c context.Context, mid int64) (blockedSum *model.ResBlockedNumUser, err error) {
	blockedSum = &model.ResBlockedNumUser{}
	blockedSum.BlockedSum, err = s.dao.BlockedNumUser(c, mid)
	return
}

// BatchPublishs get publish info list.
func (s *Service) BatchPublishs(c context.Context, ids []int64) (res map[int64]*model.BlockedAnnouncement, err error) {
	res, err = s.dao.BatchPublishs(c, ids)
	return
}

// AddBlockedInfo add blocked info.
func (s *Service) AddBlockedInfo(c context.Context, argJB *model.ArgJudgeBlocked) (err error) {
	if argJB.OID == 0 && argJB.OPName == "" {
		log.Error("origin_type(%d) oper_id(%d) && operator_name(%s) not both empty!", argJB.OType, argJB.OID, argJB.OPName)
		err = ecode.RequestErr
		return
	}
	if argJB.OID == 0 && argJB.OPName != "" {
		argJB.OID = s.managers[argJB.OPName]
	}
	if (argJB.PType == model.PunishTypeForever && argJB.BForever != model.InBlockedForever) ||
		(argJB.BForever == model.InBlockedForever && argJB.PType != model.PunishTypeForever) {
		argJB.PType = model.PunishTypeForever
		argJB.BForever = model.InBlockedForever
	}
	if argJB.PType == model.PunishTypeForever && argJB.BForever == model.InBlockedForever {
		argJB.BDays = 0
	}
	bi := &model.BlockedInfo{
		UID:            argJB.MID,
		OID:            argJB.OID,
		BlockedDays:    int64(argJB.BDays),
		BlockedForever: int64(argJB.BForever),
		BlockedRemark:  argJB.BRemark,
		MoralNum:       int64(argJB.MoralNum),
		OriginContent:  argJB.OContent,
		OriginTitle:    argJB.OTitle,
		OriginType:     int64(argJB.OType),
		OriginURL:      argJB.OURL,
		PunishTime:     xtime.Time(time.Now().Unix()),
		PunishType:     int64(argJB.PType),
		ReasonType:     int64(argJB.RType),
		BlockedType:    int64(model.PunishBlock),
		OperatorName:   argJB.OPName,
	}
	return s.dao.AddBlockedInfo(c, bi)
}

// AddBatchBlockedInfo add batch blocked info.
func (s *Service) AddBatchBlockedInfo(c context.Context, argJBs *model.ArgJudgeBatchBlocked) (err error) {
	if argJBs.OID == 0 && argJBs.OPName == "" {
		log.Error("origin_type(%d) oper_id(%d) && operator_name(%s) not both empty!", argJBs.OType, argJBs.OID, argJBs.OPName)
		err = ecode.RequestErr
		return
	}
	if argJBs.OID == 0 && argJBs.OPName != "" {
		argJBs.OID = s.managers[argJBs.OPName]
	}
	if (argJBs.PType == model.PunishTypeForever && argJBs.BForever != model.InBlockedForever) ||
		(argJBs.BForever == model.InBlockedForever && argJBs.PType != model.PunishTypeForever) {
		argJBs.PType = model.PunishTypeForever
		argJBs.BForever = model.InBlockedForever
	}
	if argJBs.PType == model.PunishTypeForever && argJBs.BForever == model.InBlockedForever {
		argJBs.BDays = 0
	}
	var bis []*model.BlockedInfo
	for _, mid := range argJBs.MID {
		bi := &model.BlockedInfo{
			UID:            mid,
			OID:            argJBs.OID,
			BlockedDays:    int64(argJBs.BDays),
			BlockedForever: int64(argJBs.BForever),
			BlockedRemark:  argJBs.BRemark,
			MoralNum:       int64(argJBs.MoralNum),
			OriginContent:  argJBs.OContent,
			OriginTitle:    argJBs.OTitle,
			OriginType:     int64(argJBs.OType),
			OriginURL:      argJBs.OURL,
			PunishTime:     xtime.Time(argJBs.PTime),
			PunishType:     int64(argJBs.PType),
			ReasonType:     int64(argJBs.RType),
			BlockedType:    int64(model.PunishBlock),
			OperatorName:   argJBs.OPName,
		}
		bis = append(bis, bi)
	}
	// begin tran
	var tx *xsql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		err = errors.Wrap(err, "s.dao.BeginTran()")
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err = s.dao.TxAddBlockedInfo(tx, bis)
	return
}

// BLKHistorys get blocked historys list.
func (s *Service) BLKHistorys(c context.Context, ah *model.ArgHistory) (rhs *model.ResBLKHistorys, err error) {
	count, err := s.dao.BLKHistoryCount(c, ah)
	if err != nil {
		err = errors.Wrap(err, "s.dao.BLKHistoryCount")
		return
	}
	rhs = &model.ResBLKHistorys{
		TotalCount: count,
		PN:         ah.PN,
		PS:         ah.PS,
		Items:      _emptyBlockInfo,
	}
	if count == 0 {
		return
	}
	rhs.Items, err = s.dao.BLKHistorys(c, ah)
	if err != nil {
		err = errors.Wrap(err, "s.dao.BLKHistorys")
		return
	}
	var uids []int64
	for _, item := range rhs.Items {
		uids = append(uids, item.UID)
	}
	infoMap, err := s.infoMap(c, uids)
	if err != nil {
		err = errors.Wrap(err, "s.infoMap")
		return
	}
	for _, item := range rhs.Items {
		if info, ok := infoMap[item.UID]; ok {
			item.Uname = info.Name
			item.Face = info.Face
		}
		item.Build()
	}
	return
}

// BatchBLKInfos mutli get blocked info by ids.
func (s *Service) BatchBLKInfos(c context.Context, ids []int64) (items map[int64]*model.BlockedInfo, err error) {
	items, err = s.dao.BlockedInfoIDs(c, ids)
	if err != nil {
		err = errors.Wrap(err, "s.dao.BLKHistorys")
		return
	}
	var uids []int64
	for _, item := range items {
		uids = append(uids, item.UID)
	}
	infoMap, err := s.infoMap(c, uids)
	if err != nil {
		err = errors.Wrap(err, "s.infoMap")
		return
	}
	for _, item := range items {
		if info, ok := infoMap[item.UID]; ok {
			item.Uname = info.Name
			item.Face = info.Face
		}
		item.Build()
	}
	return
}
