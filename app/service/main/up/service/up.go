package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	accgrpc "go-common/app/service/main/account/api"
	upgrpc "go-common/app/service/main/up/api/v1"
	"go-common/app/service/main/up/dao"
	"go-common/app/service/main/up/dao/global"
	"go-common/app/service/main/up/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

//Edit 处理通过databus注册的up主信息。如果该up曾经注册过则更新attr并更新db，否则第一次计算attr并插入db。
func (s *Service) Edit(c context.Context, mid int64, isAuthor int, from uint8) (row int64, err error) {
	var res, mdlUp *model.Up
	res, err = s.up.RawUp(c, mid) //查询db
	if err != nil {
		return
	}
	if res != nil {
		mdlUp = &model.Up{MID: res.MID, Attribute: res.Attribute}
	} else {
		mdlUp = &model.Up{MID: mid}
	}
	mdlUp.AttrSet(isAuthor, from)
	if row, err = s.up.AddUp(c, mdlUp); err != nil {
		return
	}
	if row > 0 {
		s.up.DelCpUp(c, mid)
	}
	return
}

//Info get auth info by mid.
func (s *Service) Info(c context.Context, mid int64, from uint8) (isAuthor int, err error) {
	var res *model.Up
	if res, err = s.up.Up(c, mid); err != nil {
		return
	}
	if res == nil {
		return
	}
	isAuthor = res.AttrVal(from)
	return
}

// UpAttr .
func (s *Service) UpAttr(c context.Context, req *upgrpc.UpAttrReq) (res *upgrpc.UpAttrReply, err error) {
	res = new(upgrpc.UpAttrReply)
	var up *model.Up
	if up, err = s.up.Up(c, req.Mid); err != nil {
		return
	}
	if up == nil {
		return
	}
	res.IsAuthor = uint8(up.AttrVal(req.From))
	return
}

// IdentifyAll get all type of uper auth.
func (s *Service) IdentifyAll(c context.Context, mid int64, ip string) (ia *model.IdentifyAll, err error) {
	cache := true
	if ia, err = s.up.IdentityAllCache(c, mid); err != nil {
		s.pCacheMiss.Incr("upinfo_cache")
		err = nil
		cache = false
	} else if ia != nil {
		s.pCacheHit.Incr("upinfo_cache")
		return
	}
	s.pCacheMiss.Incr("upinfo_cache")
	var (
		isArc, isArticle, isPic, isBlink int
	)
	ia = &model.IdentifyAll{}
	if isArc, err = s.Info(c, mid, model.AttrArchiveUp); err == nil && isArc > 0 {
		ia.Archive = model.IsUp
	}
	if isArticle, err = s.up.IsAuthor(c, mid, ip); err == nil && isArticle > 0 {
		ia.Article = model.IsUp
	}
	if isPic, err = s.up.Pic(c, mid, ip); err == nil && isPic > 0 {
		ia.Pic = model.IsUp
	}
	if isBlink, err = s.up.Blink(c, mid, ip); err == nil && isBlink > 0 {
		ia.Blink = model.IsUp
	}
	if cache {
		s.cacheWorker.Do(c, func(c context.Context) {
			s.up.AddIdentityAllCache(c, mid, ia)
		})
	}
	return
}

// UpsByGroup Flows get group ups list.
func (s *Service) UpsByGroup(c context.Context, group int64) (ups []*model.UpSpecial) {
	var mids []int64
	if group <= 0 {
		for _, v := range s.spGroupsCache {
			var ok bool
			if mids, ok = s.spGroupsMidsCache[v.ID]; !ok {
				continue
			}
			for _, mid := range mids {
				ups = append(ups, &model.UpSpecial{GroupID: v.ID, GroupName: v.Name, GroupTag: v.Tag, Mid: mid, Note: v.Note, FontColor: v.FontColor, BgColor: v.BgColor})
			}
		}
		return
	}
	var (
		mok, gok bool
		g        *upgrpc.UpGroup
	)
	if g, gok = s.spGroupsCache[group]; !gok {
		return
	}
	if mids, mok = s.spGroupsMidsCache[group]; !mok {
		return
	}
	for _, mid := range mids {
		ups = append(ups, &model.UpSpecial{GroupID: g.ID, GroupName: g.Name, GroupTag: g.Tag, Mid: mid, Note: g.Note, FontColor: g.FontColor, BgColor: g.BgColor})
	}
	return
}

// SpecialDel id 是数据库主键
func (s *Service) SpecialDel(c *bm.Context, id int64) (affectedRow int64, err error) {

	oldUp, err := s.mng.GetSpecialByID(c, id)
	if err != nil {
		log.Warn("get old up fail from db, id=%d", id)
		return
	}
	if oldUp == nil {
		log.Warn("no old up found in db, id=%d", id)
		return
	}
	// 检查是否有特殊权限
	err = s.SpecialGroupPermit(c, oldUp.GroupID)
	if err != nil {
		return
	}

	res, err2 := s.mng.DelSpecialByID(c, id)
	err = err2
	if err != nil {
		log.Error("fail to delete id, id=%d, err=%v", id, err)
		return
	}
	affectedRow, _ = res.RowsAffected()
	if affectedRow > 0 {

		var up = *oldUp
		var uname, _ = BmGetStringOrDefault(c, "username", "unkown")
		var uid, _ = BmGetInt64OrDefault(c, "uid", 0)
		var oplog = UpSpecialLogInfo{
			Up:     &model.UpSpecialWithName{UpSpecial: up},
			UID:    uid,
			UName:  uname,
			CTime:  time.Now(),
			Action: ActionDelete,
		}
		s.sendUpSpecialLog(c, &oplog)
	}
	log.Info("delete id from ups, id=%d, affected row=%d", id, affectedRow)
	return
}

//SpecialAdd add special
func (s *Service) SpecialAdd(c context.Context, adminName string, special *model.UpSpecial, mids ...int64) (affectedRow int64, err error) {

	var insertMids []int64
	for _, mid := range mids {
		<-s.specialAddDBLimiter
		// 先获取是否有这个用户
		var oldUps, err1 = s.mng.GetSpecialByMidGroup(c, mid, special.GroupID)

		if err1 != nil {
			log.Error("get up from db fail, err=%v, mid=%d", err1, mid)
			continue
		}
		var up = *special
		up.Mid = mid
		var oplog = UpSpecialLogInfo{
			Up:     &model.UpSpecialWithName{UpSpecial: up},
			UID:    up.UID,
			UName:  adminName,
			CTime:  time.Now(),
			Action: ActionAdd,
		}
		// 如果是老用户，则记录为edit
		if oldUps != nil {
			oplog.Action = ActionEdit
			oplog.UpOld = &model.UpSpecialWithName{UpSpecial: *oldUps}
		}
		s.sendUpSpecialLog(c, &oplog)
		// 如果有，只是更新一下
		if oldUps != nil {
			s.mng.UpdateSpecialByID(c, oldUps.ID, special)
			log.Info("only update up's info, id=%d, mid=%d", oldUps.ID, mid)
			continue
		}
		// 如果没有，添加一下
		insertMids = append(insertMids, mid)
	}
	// 更新数据库
	if len(insertMids) == 0 {
		return
	}

	res, err2 := s.mng.InsertSpecial(c, special, insertMids...)
	err = err2
	if err != nil {
		log.Error("fail to add special ups, err=%v, midscount=%d", err, len(mids))
		return
	}
	affectedRow, _ = res.RowsAffected()

	log.Info("insert into ups, affected row=%d, midscount=%d", affectedRow, len(mids))
	return
}

//SpecialEdit edit up special
func (s *Service) SpecialEdit(c *bm.Context, special *model.UpSpecial, id int64) (affectedRow int64, err error) {
	oldUp, err := s.mng.GetSpecialByID(c, id)
	if err != nil {
		log.Warn("get old up fail from db, id=%d", id)
		return
	}
	if oldUp == nil {
		log.Warn("no old up found in db, id=%d", id)
		err = errors.New("up not found in db")
		return
	}

	special.Mid = oldUp.Mid

	res, err := s.mng.UpdateSpecialByID(c, id, special)
	if err != nil {
		log.Error("error update ups, id=%d, info=%v", id, special)
		return
	}
	affectedRow, _ = res.RowsAffected()
	log.Info("only update up's info, id=%d, info=%v, rowcount=%d", id, special, affectedRow)
	if affectedRow > 0 {
		var uname, _ = BmGetStringOrDefault(c, "username", "unkown")
		var oplog = UpSpecialLogInfo{
			UpOld:  &model.UpSpecialWithName{UpSpecial: *oldUp},
			Up:     &model.UpSpecialWithName{UpSpecial: *special},
			UID:    special.UID,
			UName:  uname,
			CTime:  time.Now(),
			Action: ActionEdit,
		}
		s.sendUpSpecialLog(c, &oplog)

	}
	return
}

// SpecialGet 这个接口支持更多参数，直接从数据库来查询
func (s *Service) SpecialGet(c *bm.Context, arg *model.GetSpecialArg) (res []*model.UpSpecialWithName, total int, err error) {
	var conditions []dao.Condition
	var con dao.Condition
	//先去查询管理员id
	if arg.AdminName != "" {
		r, err2 := s.mng.GetUIDByNames(c, []string{arg.AdminName})
		err = err2
		k, exist := r[arg.AdminName]
		if !exist {
			err = errors.New("admin name not found")
			log.Error("admin name not found, name=%s", arg.AdminName)
			return
		}
		log.Info("get uid by name, name=%s, id=%d", arg.AdminName, k)
		arg.UID = int(k)
	}

	if arg.GroupID != 0 {
		con = dao.Condition{
			Key:      "ups.type",
			Operator: "=",
			Value:    arg.GroupID,
		}
		conditions = append(conditions, con)
	}

	if arg.UID != 0 {
		con = dao.Condition{
			Key:      "ups.uid",
			Operator: "=",
			Value:    arg.UID,
		}
		conditions = append(conditions, con)
	}

	if arg.FromTime != "" {
		con = dao.Condition{
			Key:      "ups.mtime",
			Operator: ">=",
			Value:    arg.FromTime,
		}
		conditions = append(conditions, con)
	}

	if arg.ToTime != "" {
		con = dao.Condition{
			Key:      "ups.mtime",
			Operator: "<=",
			Value:    arg.ToTime,
		}
		conditions = append(conditions, con)
	}

	if arg.Mids != "" {
		var midstr = strings.Split(arg.Mids, ",")
		var mids []int64
		for _, v := range midstr {
			s, e := strconv.ParseInt(strings.Trim(v, " \n"), 10, 64)
			if e != nil {
				continue
			}
			mids = append(mids, s)
		}
		if len(mids) != 0 {
			con = dao.Condition{
				Key:      "ups.mid",
				Operator: " in (" + xstr.JoinInts(mids) + ")",
			}
			conditions = append(conditions, con)
		}

	}

	conditions = dao.AndCondition(conditions...)

	var order = "asc"
	if arg.Order == "desc" {
		order = arg.Order

	}
	con = dao.Condition{
		Key:   "order by ups.mtime",
		After: order,
	}
	conditions = append(conditions, con)

	var offset = (arg.Pn - 1) * arg.Ps
	var limit = arg.Ps

	// 非导出数据时做分页
	if arg.Export == "" {
		total, err = s.mng.GetSepcialCount(c, conditions...)
		if err != nil {
			log.Error("get ups count err, err=%v", err)
			return
		}

		conditions = append(conditions,
			dao.Condition{Key: fmt.Sprintf("LIMIT %d", limit)},
			dao.Condition{Key: fmt.Sprintf("OFFSET %d", offset)},
		)
	} else {
		if arg.Charset == "" {
			arg.Charset = "gbk"
		}
	}

	var ups []*model.UpSpecial
	ups, err = s.mng.GetSpecial(c, conditions...)

	if err != nil {
		log.Error("get ups err, err=%v", err)
		return
	}
	// 查询member name
	var (
		mids, uids  []int64
		infosReply  *accgrpc.InfosReply
		uid2NameMap = map[int64]string{}
	)

	for _, up := range ups {
		var upsWithName = &model.UpSpecialWithName{}
		upsWithName.Copy(up)
		res = append(res, upsWithName)
		uids = append(uids, up.UID)
		mids = append(mids, up.Mid)
	}

	if infosReply, err = global.GetAccClient().Infos3(c, &accgrpc.MidsReq{Mids: mids, RealIp: metadata.String(c, metadata.RemoteIP)}); err != nil {
		log.Error("global.GetAccClient().Infos3(%+d) error(%+v)", mids, err)
		err = nil
	}

	if len(uids) != 0 {
		if uid2NameMap, err = s.mng.GetUNamesByUids(c, uids); err != nil {
			log.Error("s.mng.GetUNamesByUids(%+v) err error(%+v)", uids, err)
			err = nil
		}
	}
	for _, v := range res {
		v.AdminName = uid2NameMap[v.UID]
		if infosReply != nil && infosReply.Infos != nil {
			if info := infosReply.Infos[v.Mid]; info != nil {
				v.UName = info.Name
			}
		}
	}

	log.Info("get special, arg=%+v, err=%v", arg, err)
	return
}

//ListUpBase list up base info
func (s *Service) ListUpBase(c *bm.Context, size int, lastID int64, activity []int64) (mids []int64, newLastID int64, err error) {
	where := ""
	if len(activity) > 0 {
		where = fmt.Sprintf("AND activity IN (%s)", xstr.JoinInts(activity))
	}

	var idMids map[int64]int64
	if idMids, err = s.card.ListUpBase(c, size, lastID, where); err != nil {
		log.Error("fail to list up base info, err=%v", err)
		return
	}

	for id, mid := range idMids {
		mids = append(mids, mid)
		if id > newLastID {
			newLastID = id
		}
	}

	return
}

// UpInfoActivitys up info activity by last_id.
func (s *Service) UpInfoActivitys(c context.Context, req *upgrpc.UpListByLastIDReq) (res *upgrpc.UpActivityListReply, err error) {
	var mup map[int64]*upgrpc.UpActivity
	res = new(upgrpc.UpActivityListReply)
	if mup, err = s.up.UpInfoActivitys(c, req.LastID, req.Ps); err != nil {
		return
	}
	for id, up := range mup {
		res.UpActivitys = append(res.UpActivitys, up)
		if id > res.LastID {
			res.LastID = id
		}
	}
	return
}

//SpecialGetByMid get special
func (s *Service) SpecialGetByMid(c *bm.Context, arg *model.GetSpecialByMidArg) (res []*model.UpSpecial, err error) {
	var condition = dao.Condition{
		Key:      "ups.mid",
		Operator: "=",
		Value:    arg.Mid,
	}

	res, err = s.mng.GetSpecial(c, condition)
	if err != nil {
		log.Error("error when get by mid, err=%v", err)
	}
	log.Info("get special, mid=%d, result=%v", arg.Mid, res)
	return
}

// UpSpecial .
func (s *Service) UpSpecial(c context.Context, req *upgrpc.UpSpecialReq) (res *upgrpc.UpSpecialReply, err error) {
	res = new(upgrpc.UpSpecialReply)
	res.UpSpecial, err = s.mng.UpSpecial(c, req.Mid)
	return
}

// UpsSpecial .
func (s *Service) UpsSpecial(c context.Context, req *upgrpc.UpsSpecialReq) (res *upgrpc.UpsSpecialReply, err error) {
	res = new(upgrpc.UpsSpecialReply)
	res.UpSpecials, err = s.mng.UpsSpecial(c, req.Mids)
	return
}

// UpGroupMids .
func (s *Service) UpGroupMids(c context.Context, req *upgrpc.UpGroupMidsReq) (res *upgrpc.UpGroupMidsReply, err error) {
	res = new(upgrpc.UpGroupMidsReply)
	var (
		ok    bool
		mids  []int64
		start = (req.Pn - 1) * req.Ps
		end   = req.Pn * req.Ps
	)
	if mids, ok = s.spGroupsMidsCache[req.GroupID]; !ok {
		return
	}
	res.Total = len(mids)
	switch {
	case res.Total <= int(start):
		res.Mids = []int64{}
	case res.Total <= int(end):
		res.Mids = mids[start:]
	default:
		res.Mids = mids[start:end]
	}
	return
}
