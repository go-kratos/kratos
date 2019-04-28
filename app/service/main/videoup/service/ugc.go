package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/app/service/main/videoup/model/prom"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

// AddByUGC add ugc archive and videos.
func (s *Service) AddByUGC(c context.Context, ap *archive.ArcParam) (aid int64, err error) {
	var (
		a = &archive.Archive{
			Mid:       ap.Mid,
			Author:    ap.Author,
			TypeID:    ap.TypeID,
			Title:     ap.Title,
			Cover:     ap.Cover,
			Tag:       ap.Tag,
			Copyright: ap.Copyright,
			Desc:      ap.Desc,
			State:     archive.StateForbidSubmit,
		}
		fns        = make([]string, 0, len(ap.Videos))
		cids       = make([]int64, 0, len(ap.Videos)) // NOTE: for code mode
		nvs        = make([]*archive.Video, 0, len(ap.Videos))
		forbidList = make([]*archive.ForbidAttr, 0, 3)
		nvmIdx     = make(map[string]int)
		hasBGM     bool
	)
	if ap.Porder == nil {
		ap.Porder = &archive.Porder{}
	}
	if ap.Porder.IndustryID > 0 {
		flowForbid, _ := s.transferAttrByScope(c, a, ap)
		if flowForbid != nil {
			log.Info(" flowForbid(%+v)", flowForbid)
			forbidList = append(forbidList, flowForbid)
		}
	}
	// set attr
	if ap.Lang == "jp" {
		a.AttrSet(archive.AttrYes, archive.AttrBitIsJapan)
	}
	a.AttrSet(ap.NoReprint, archive.AttrBitNoReprint)
	a.AttrSet(ap.UGCPay, archive.AttrBitUGCPay)
	a.AttrSet(archive.AttrYes, archive.AttrBitAllowTag)
	// prom
	s.promSub.Incr(prom.RouteUpFrom + strconv.Itoa(int(ap.UpFrom)))
	s.promP.Incr(prom.FormatVideoKey(len(ap.Videos), "add"))
	s.promSub.Incr(prom.RouteNormalMode)
	// validate
	log.Info("aid() filename(%s) cid(%s) start tran add ugc", strings.Join(fns, ","), xstr.JoinInts(cids))
	if len(ap.Videos) == 0 {
		err = ecode.VideoAbnormalSubmit
		return
	}
	for idx, v := range ap.Videos {
		if !s.filenameCheck(v.Filename) {
			log.Error("s.arc.CheckNewVideosFilename (%+v)", fns)
			newErr := ecode.VideoupVideoFilenameErr
			err = ecode.Errorf(newErr, newErr.Message(), idx+1)
			return
		}
		if v.Editor != nil && v.Editor.Bgms != nil {
			hasBGM = true
		}
		nv := &archive.Video{
			Title:      v.Title,
			Desc:       v.Desc,
			Filename:   v.Filename,
			SrcType:    "vupload",
			Status:     archive.VideoStatusSubmit,
			XcodeState: archive.VideoUploadInfo,
			Sid:        v.Sid,
		}
		fns = append(fns, v.Filename)
		nvs = append(nvs, nv)
		nvmIdx[v.Filename] = idx + 1
	}
	if hasBGM {
		a.AttrSet(archive.AttrYes, archive.AttrBitHasBGM)
	}
	var avm map[int64]map[string]*archive.Video
	if hasTimeout, expiredFilename, _ := s.arc.CheckNewVideosTimeout(c, fns); hasTimeout {
		log.Error("s.arc.CheckNewVideosTimeout(%+v) timeout", fns)
		newErr := ecode.VideoupFilenameExpired
		err = ecode.Errorf(newErr, newErr.Message(), nvmIdx[expiredFilename])
		return
	}
	if avm, err = s.arc.NewVideosByFn(c, fns); err != nil {
		log.Error("s.arc.VideosByFn(%+v) error(%v)", fns, err)
		return
	}
	if len(avm) > 0 {
		err = ecode.VideoAbnormalSubmit // NOTE: when not code mode, fns must have not db record.
	}
	if err != nil {
		return
	}
	// begin tran
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	if s.hasForbid(c, ap.Mid) {
		s.transferAttr(c, a, ap)
	}
	if aid, err = s.txAddArchive(c, tx, a); err != nil {
		tx.Rollback()
		log.Error("s.txAddArchive(%+v) error(%v)", a, err)
		return
	}
	if err = s.assignCid(c, tx, nvs); err != nil {
		tx.Rollback()
		log.Error("s.assignCid(%d) error(%v)", aid, err)
		return
	}
	for i, v := range nvs {
		v.Aid = aid
		v.Index = i + 1
	}
	if err = s.txAddVideos(tx, nvs); err != nil {
		tx.Rollback()
		log.Error("s.txAddVideos() error(%v)", err)
		return
	}
	log.Info("aid(%d) attribute(%d) filename(%s) cid(%s) add videos", aid, a.Attribute, strings.Join(fns, ","), xstr.JoinInts(cids))
	if ap.DTime > 0 {
		if _, err = s.arc.TxAddDelay(tx, ap.Mid, aid, archive.StateOpen, archive.DelayTypeForUser, ap.DTime); err != nil {
			tx.Rollback()
			log.Error("s.arc.txAddDelay(%d, %d, %d) error(%v)", ap.Mid, aid, ap.DTime, err)
			return
		}
		log.Info("aid(%d) filename(%s) cid(%s) update archive_delay by mid(%d) dtime(%v)", aid, strings.Join(fns, ","), xstr.JoinInts(cids), ap.Mid, ap.DTime)
	}
	desc := ap.Desc
	if ap.DescFormatID <= 0 {
		desc = ""
	}
	if ap.POI != nil {
		data, _ := json.Marshal(ap.POI)
		if _, err = s.arc.TxUpArchiveBiz(tx, aid, archive.BIZPOI, string(data)); err != nil {
			tx.Rollback()
			log.Error("s.arc.TxUpArchiveBiz(%d ,%d,%+v) error(%v)", a.Aid, archive.BIZPOI, ap.POI, err)
			return
		}
	}
	if ap.Vote != nil {
		data, _ := json.Marshal(ap.Vote)
		if _, err = s.arc.TxUpArchiveBiz(tx, aid, archive.BIZVote, string(data)); err != nil {
			tx.Rollback()
			log.Error("s.arc.TxUpArchiveBiz(%d ,%d,%+v) error(%v)", a.Aid, archive.BIZVote, ap.Vote, err)
			return
		}
	}
	if _, err = s.arc.TxUpAddit(tx, aid, ap.MissionID, ap.OrderID, ap.Porder.FlowID, ap.DescFormatID, ap.IPv6, ap.Source, ap.Porder.Advertiser, ap.FlowRemark, desc, ap.Dynamic, ap.UpFrom); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpAddit(%d,%d,%q,%d,%d,%d,%s,%s,%s,%s,%s,%d) error(%v)", a.Aid, ap.MissionID, ap.IPv6, ap.OrderID, ap.Porder.FlowID, ap.DescFormatID, ap.Source, ap.Porder.Advertiser, ap.FlowRemark, desc, ap.Dynamic, ap.UpFrom, err)
		return
	}
	log.Info("aid(%d) filename(%s) cid(%s) add archive_addit missionid(%d) ipv6(%q) order_id(%d) source(%s) upfrom(%d) flow_id(%d)  advertiser(%s) flowRemark(%s) descFormation(%d) desc(%s) dynamic(%s)",
		aid, strings.Join(fns, ","), xstr.JoinInts(cids), ap.MissionID, ap.IPv6, ap.OrderID, ap.Source, ap.UpFrom, ap.Porder.FlowID, ap.Porder.Advertiser, ap.FlowRemark, ap.DescFormatID, ap.Desc, ap.Dynamic)
	if ap.Porder.IndustryID > 0 {
		//基于flow_scope场景生成 flow_design 对应流量限制方案
		if ap.Porder.FlowID, err = s.getGroupIDByScope(c, ap); err != nil {
			log.Error("s.getGroupIDByScope aid(%+v) error(%+v)", ap.Aid, err)
		}
		if ap.Porder.FlowID > 0 {
			//私单自动进flow_design业务
			if err = s.txAddFlowID(c, tx, aid, ap.Porder.FlowID, ap); err != nil {
				tx.Rollback()
				return
			}
		}
	}
	ap.Aid = a.Aid
	//天马白名单 自动进flow_design
	if err = s.txAddAppFeedID(c, tx, ap); err != nil {
		tx.Rollback()
		return
	}
	if s.hasForbid(c, ap.Mid) {
		midForbid, _ := s.getMidForbid(c, aid, ap)
		log.Info(" midForbid(%+v)", midForbid)
		if midForbid != nil {
			forbidList = append(forbidList, midForbid)
		}
	}
	if len(forbidList) > 0 {
		forbidEnd := s.mergeForbid(c, forbidList)
		forbidEnd.Aid = aid
		log.Info(" forbidEnd(%+v)", forbidEnd)
		if _, err = s.arc.TxUpForbidAttr(tx, forbidEnd); err != nil {
			log.Error("s.arc.TxUpForbidAttr forbidEnd(%+v) error(%+v)", forbidEnd, err)
			return
		}
	}
	log.Info("aid(%d) filename(%s) cid(%s) history", aid, strings.Join(fns, ","), xstr.JoinInts(cids))
	if err = s.txAddHistory(tx, a, nvs); err != nil {
		tx.Rollback()
		log.Error("s.txAddHistory(%+v) error(%v)", a, err)
		return
	}
	log.Info("aid(%d) filename(%s) cid(%s) commit", aid, strings.Join(fns, ","), xstr.JoinInts(cids))
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	log.Info("aid(%d) filename(%s) cid(%s) end tran add ugc", aid, strings.Join(fns, ","), xstr.JoinInts(cids))
	s.asyncCh <- func() {
		log.Info("aid(%d) filename(%s) cid(%s) start tran cid ", aid, strings.Join(fns, ","), xstr.JoinInts(cids))
		if ap.Staffs != nil && ap.HandleStaff {
			s.HandleArchiveApplys(context.TODO(), a.Aid, ap.Staffs, "add", false)
		}
		s.syncCid(aid, ap.Mid, nvs, true)
		s.busAddArchive(aid, ap.UpFrom)
		s.bindBgmWithVideo(aid, ap.Mid, nvs)
	}
	return
}

//EditMissionByUGC  ugc can update addit mission_id
func (s *Service) EditMissionByUGC(c context.Context, ap *archive.ArcMissionParam) (err error) {
	var (
		a            *archive.Archive
		cancelMisson int64
		addMisson    bool
	)
	if a, err = s.arc.Archive(c, ap.AID); err != nil {
		log.Error("s.arc.Archive(%d) error(%v)", ap.AID, err)
		return
	}
	if a == nil {
		err = ecode.NothingFound
		return
	}
	if a.NotAllowUp() {
		err = ecode.TargetBlocked
		return
	}
	if s.isPGC(a) {
		err = ecode.RequestErr
		return
	}
	if a.Mid != ap.MID {
		err = ecode.RequestErr
		return
	}
	//更换活动需求不存在  且不允许直接更换  -400
	//该接口仅处理  取消活动 ap.MissionID=0  以及 新增活动绑定活动Tag  ap.MissionID && ap.Tag
	if add, _ := s.arc.Addit(c, ap.AID); add != nil {
		if add.MissionID > 0 {
			if ap.MissionID == 0 {
				//取消活动
				cancelMisson = add.MissionID
			} else if add.MissionID != ap.MissionID {
				//不允许直接覆盖活动
				err = ecode.RequestErr
				return
			}
		} else {
			//add mission && tag
			addMisson = true
		}
	}
	_, err = s.arc.UpMissionID(c, ap)
	if err != nil {
		log.Error("aid(%d) Mission(%d) update error(%+v)", ap.AID, ap.MissionID, err)
		return
	}
	//新增活动及tag
	if addMisson && len(ap.Tag) > 0 {
		var tag string
		if len(a.Tag) > 0 {
			tag = fmt.Sprintf("%s,%s", a.Tag, ap.Tag)
		} else {
			tag = ap.Tag
		}
		tag = strings.Join(Slice2String(SliceUnique(Slice2Interface(strings.Split(tag, ",")))), ",")
		if err = s.UpTag(c, a.Aid, tag); err != nil {
			log.Error("aid(%d) Mission(%d) update tag(%s) error(%+v)", ap.AID, ap.MissionID, ap.Tag, err)
		}
	}
	s.asyncCh <- func() {
		s.busSecondRound(ap.AID, cancelMisson)
	}
	log.Info("aid(%d) Mission(%d) update success cancelMisson(%d)  addMisson(%s)", ap.AID, ap.MissionID, cancelMisson, addMisson)
	return
}

//QueryDynamicSetting  user can share as dynamic
func (s *Service) QueryDynamicSetting(c context.Context, ap *archive.ArcDynamicParam) (canShare bool, err error) {
	var (
		a *archive.Archive
	)
	if a, err = s.arc.Archive(c, ap.AID); err != nil {
		log.Error("s.arc.Archive(%d) error(%v)", ap.AID, err)
		return
	}
	if a == nil {
		err = ecode.NothingFound
		return
	}
	canShare = true
	if s.isDynamicForbid(a) && a.Mid == ap.MID {
		canShare = false
	}
	return
}

// EditByUGC edit ugc archive and videos.
func (s *Service) EditByUGC(c context.Context, ap *archive.ArcParam) (err error) {
	var (
		a                                                 *archive.Archive
		aChange, vChange, tagChange, typeIDChange, hasBGM bool
		aReason, vReason                                  string
		nvs, evs, dvs, hvs                                []*archive.Video
	)
	s.promP.Incr(prom.FormatVideoKey(len(ap.Videos), "edit"))
	if a, aChange, tagChange, typeIDChange, aReason, err = s.checkArchive(c, ap); err != nil {
		log.Error("s.checkArchive(%v) error(%v)", ap, err)
		return
	}
	//init Porder
	if ap.Porder == nil {
		ap.Porder = &archive.Porder{}
	}
	if len(ap.Videos) > 0 {
		s.promSub.Incr(prom.RouteNormalMode)
		if nvs, evs, dvs, hvs, vChange, vReason, err = s.checkVideo(c, ap, false); err != nil {
			log.Error("aid(%d) filename() s.checkVideo() error(%v)", ap.Aid, err)
			return
		}
		for _, v := range nvs {
			if !s.filenameCheck(v.Filename) {
				newErr := ecode.VideoupVideoFilenameErr
				err = ecode.Errorf(newErr, newErr.Message(), v.Index)
				return
			}
		}
		for _, v := range ap.Videos {
			if v.Editor != nil && v.Editor.Bgms != nil {
				//一次写入不可修改
				hasBGM = true
			}
		}
	}
	log.Info("aid(%d) start tran edit ugc aChange(%v) && vChange(%v)", ap.Aid, aChange, vChange)

	// begin tran
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	if aChange {
		if err = s.txEditArchvie(tx, a); err != nil {
			tx.Rollback()
			log.Error("s.txEditArchvie(archive(%+v)) error(%v)", a, err)
			return
		}
	}
	if _, err = s.arc.TxUpArcAttr(tx, a.Aid, archive.AttrBitNoReprint, ap.NoReprint); err != nil {
		tx.Rollback()
		return
	}
	if _, err = s.arc.TxUpArcAttr(tx, a.Aid, archive.AttrBitUGCPay, ap.UGCPay); err != nil {
		tx.Rollback()
		return
	}
	if hasBGM {
		if _, err = s.arc.TxUpArcAttr(tx, a.Aid, archive.AttrBitHasBGM, 1); err != nil {
			tx.Rollback()
			return
		}
	}
	log.Info("aid(%d) update no reprint attr(%d)", a.Aid, ap.NoReprint)
	if err = s.assignCid(c, tx, nvs); err != nil {
		tx.Rollback()
		log.Error("s.assignCid(%d) error(%v)", ap.Aid, err)
		return
	}
	if err = s.txAddVideos(tx, nvs); err != nil {
		tx.Rollback()
		return
	}
	// >= threshold async update
	if len(evs) < s.c.AsyncThreshold {
		if success := s.lockEdit(c, ap.Aid); !success {
			log.Error("EditByUGC(%d) error(%v)", ap.Aid, ecode.VideoupEditLocked)
			return ecode.VideoupEditLocked
		}
		defer s.unlockEdit(context.TODO(), ap.Aid)

		log.Info("sync edit aid(%d)", a.Aid)
		if err = s.txEditVideos(tx, evs); err != nil {
			tx.Rollback()
			return
		}
	}
	if err = s.txDelVideos(tx, dvs); err != nil {
		tx.Rollback()
		return
	}
	desc := ap.Desc
	if ap.DescFormatID <= 0 {
		desc = ""
	}
	if ap.POI != nil {
		data, _ := json.Marshal(ap.POI)
		if _, err = s.arc.TxUpArchiveBiz(tx, a.Aid, archive.BIZPOI, string(data)); err != nil {
			tx.Rollback()
			log.Error("s.arc.TxUpArchiveBiz(%d ,%d,%+v) error(%v)", a.Aid, archive.BIZPOI, ap.POI, err)
			return
		}
	}
	if _, err = s.arc.TxUpAddit(tx, a.Aid, ap.MissionID, ap.OrderID, ap.Porder.FlowID, ap.DescFormatID, ap.IPv6, ap.Source, ap.Porder.Advertiser, ap.FlowRemark, desc, ap.Dynamic, ap.UpFrom); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpAddit(%d,%d,%q,%d,%d,%s,%s,%s,%s,%d) error(%v)", a.Aid, ap.MissionID, ap.IPv6, ap.OrderID, ap.Porder.FlowID, ap.DescFormatID, ap.Source, ap.Porder.Advertiser, ap.FlowRemark, desc, ap.Dynamic, ap.UpFrom, err)
		return
	}
	log.Info("aid(%d) update archive_addit missionid(%d) ipv6(%q) order_id(%d) source(%s) upfrom(%d) flow_id(%d)  advertiser(%s) flowRemark(%s)  descFormation(%d) desc(%s) dynamic(%s)",
		ap.Aid, ap.MissionID, ap.IPv6, ap.OrderID, ap.Source, ap.UpFrom, ap.Porder.FlowID, ap.Porder.Advertiser, ap.FlowRemark, ap.DescFormatID, ap.Desc, ap.Dynamic)
	if aChange || vChange {
		if _, err = s.arc.TxUpAdditReason(tx, a.Aid, aReason+vReason); err != nil {
			tx.Rollback()
			log.Error("s.arc.TxUpAdditReason(%d,%s) error(%v)", a.Aid, aReason+vReason, err)
			return
		}
		log.Info("aid(%d) update archive reason(%s)", ap.Aid, aReason+vReason)
	}
	if err = s.txDelayTime(tx, ap.Mid, ap.Aid, ap.DTime); err != nil {
		tx.Rollback()
		return
	}
	log.Info("aid(%d) update archive_delay by mid(%d) dtime(%v)", ap.Aid, ap.Mid, ap.DTime)
	if err = s.txAddHistory(tx, a, hvs); err != nil {
		tx.Rollback()
		log.Error("s.txAddHistory(%+v) error(%v)", a, err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	log.Info("aid(%d) end tran edit ugc by nvs_num(%d) evs_num(%d) dvs_num(%d)", a.Aid, len(nvs), len(evs), len(dvs))

	afterhandler := func() {
		//不支持PGC
		if ap.HandleStaff {
			//up清空staffs 情况
			if ap.Staffs == nil {
				ap.Staffs = []*archive.StaffParam{}
			}
			s.HandleArchiveApplys(context.TODO(), a.Aid, ap.Staffs, "edit", false)
		}
		s.syncCid(ap.Aid, ap.Mid, nvs, true)
		s.busModifyArchive(ap.Aid, aChange, vChange, tagChange, len(nvs) > 0, typeIDChange)
		s.busDeleteVideo(dvs)
		s.busModifyVideo(evs)
		s.bindBgmWithVideo(ap.Aid, ap.Mid, nvs)
		log.Info("archive(%d) edit Success", ap.Aid)
	}

	if len(evs) <= s.c.SplitThreshold {
		if len(evs) >= s.c.AsyncThreshold {
			if success := s.lockEdit(c, ap.Aid); !success {
				log.Error("EditByUGC(%d) error(%v)", ap.Aid, ecode.VideoupEditLocked)
				return ecode.VideoupEditLocked
			}
		}
		s.asyncCh <- func() {
			// >=threshold async update
			if len(evs) >= s.c.AsyncThreshold {
				defer s.unlockEdit(context.TODO(), ap.Aid)
				s.asyncUpdate(context.TODO(), evs)
				log.Info("async edit aid(%d) success len(%d)", a.Aid, len(evs))
			}
			afterhandler()
		}
	}

	if len(evs) > s.c.SplitThreshold {
		if success := s.lockEdit(c, ap.Aid); !success {
			log.Error("EditByUGC(%d) error(%v)", ap.Aid, ecode.VideoupEditLocked)
			return ecode.VideoupEditLocked
		}

		log.Info("EditByUGC split(%d) len(%d)", ap.Aid, len(evs))
		s.veditor.Add(
			ap.Aid,
			// success callback
			func() {
				afterhandler()
				s.unlockEdit(context.TODO(), ap.Aid)
				log.Info("archive(%d) edit Success len(%d)", ap.Aid, len(evs))
			},
			// fail callback
			func() {
				s.unlockEdit(context.TODO(), ap.Aid)
				log.Error("archive(%d) edit Fail len(%d)", ap.Aid, len(evs))
			},
			time.Duration(s.c.EditTimeout),
			// retry funcs
			s.splitUpdate(ap.Aid, evs)...)
	}
	return
}

func (s *Service) asyncUpdate(c context.Context, evs []*archive.Video) (err error) {
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("asyncUpdate wocao jingran recover le error(%v)", r)
		}
	}()
	if err = s.txEditVideosAsync(tx, evs); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	return
}

func (s *Service) splitUpdate(aid int64, evs []*archive.Video) (retry []func() (int64, int, int, error)) {
	var (
		splitTH = s.c.SplitGroupCount
	)

	for i := 0; i < len(evs); i += splitTH {
		head := i
		tail := i + splitTH
		if tail > len(evs) {
			tail = len(evs)
		}
		log.Info("splitUpdate evs[%d:%d]", head, tail)
		retry = append(retry, func() (int64, int, int, error) {
			return aid, head, tail, s.asyncUpdate(context.TODO(), evs[head:tail])
		})
	}
	return
}

func (s *Service) lockEdit(c context.Context, aid int64) bool {
	res, _ := s.arc.SetNXLock(c, aid, int64(time.Duration(s.c.EditTimeout)/time.Second))
	return res
}

func (s *Service) unlockEdit(c context.Context, aid int64) {
	s.arc.DelLock(c, aid)
}
