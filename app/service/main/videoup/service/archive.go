package service

import (
	"context"
	"time"

	"encoding/json"
	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

// SimpleArchive get SimpleArchive
func (s *Service) SimpleArchive(c context.Context, aid int64, mode int) (a *archive.SimpleArchive, err error) {
	if a, err = s.arc.SimpleArchive(c, aid); err != nil {
		log.Error("s.arc.SimpleArchive(%d) error(%v)", aid, err)
		return
	}
	if a == nil {
		log.Error("s.arc.SimpleArchive(%d) not found", aid)
		err = ecode.NothingFound
		return
	}
	if mode == 1 {
		// videos with simple info without deleted videos
		var vs []*archive.Video
		var svs []*archive.SimpleVideo
		if svs, err = s.arc.SimpleArcVideos(c, aid); err != nil {
			log.Error("s.arc.SimpleArchive(%d) error(%v)", aid, err)
			return
		}
		if svs == nil {
			err = ecode.NothingFound
			return
		}
		for _, sv := range svs {
			if sv.Status == archive.VideoStatusDelete {
				continue
			}
			v := &archive.Video{
				Cid:    sv.Cid,
				Title:  sv.Title,
				Index:  sv.Index,
				Status: sv.Status,
				MTime:  sv.MTime,
			}
			vs = append(vs, v)
		}
		a.Videos = vs
	} else if mode == 2 {
		// videos with full info, with deleted videos
		var vs []*archive.Video
		if vs, err = s.arc.NewVideos(c, aid); err != nil {
			log.Error("s.arc.SimpleArchive(%d) error(%v)", aid, err)
			return
		}
		if vs == nil {
			err = ecode.NothingFound
			return
		}
		a.Videos = vs
	}
	return
}

// SimpleVideos  get SimpleVideos
func (s *Service) SimpleVideos(c context.Context, aid int64) (svs []*archive.SimpleVideo, err error) {
	if svs, err = s.arc.SimpleArcVideos(c, aid); err != nil {
		log.Error("s.arc.SimpleVideos(%d) error(%v)", aid, err)
		return
	}
	if svs == nil {
		err = ecode.NothingFound
	}
	return
}

// ArchivePOI get poi info when add archive
func (s *Service) ArchivePOI(c context.Context, aid int64) (p *archive.PoiObj, err error) {
	var (
		a []byte
	)
	p = &archive.PoiObj{}
	// archive with POI
	if a, err = s.arc.POI(c, aid); err != nil || a == nil {
		log.Error("s.arc.Archive(%d) error(%v)", aid, err)
		return
	}
	if err = json.Unmarshal(a, &p); err != nil {
		log.Error("s.arc.ArchivePOI(%d) error(%v)", aid, err)
		return
	}
	return

}

// ArchiveVote get archive.
func (s *Service) ArchiveVote(c context.Context, aid int64) (p *archive.Vote, err error) {
	var (
		a []byte
	)
	p = &archive.Vote{}
	// archive with Vote
	if a, err = s.arc.Vote(c, aid); err != nil || a == nil {
		log.Error("s.arc.ArchiveVote(%d) error(%v)", aid, err)
		return
	}
	if err = json.Unmarshal(a, &p); err != nil {
		log.Error("s.arc.ArchiveVote(%d) error(%v)", aid, err)
		return
	}
	return

}

// Archive get archive.
func (s *Service) Archive(c context.Context, aid int64) (av *archive.ArcVideos, err error) {
	var (
		a    *archive.Archive
		dl   *archive.Delay
		avs  []*archive.Video
		rsnm map[int64]string
	)
	// archive
	if a, err = s.arc.Archive(c, aid); err != nil {
		log.Error("s.arc.Archive(%d) error(%v)", aid, err)
		return
	}
	if a == nil {
		err = ecode.NothingFound
		return
	}
	a.NoReprint = a.AttrVal(archive.AttrBitNoReprint)
	a.UGCPay = a.AttrVal(archive.AttrBitUGCPay)
	if dl, _ = s.arc.Delay(c, aid, archive.DelayTypeForUser); dl != nil {
		a.DTime = dl.DTime
	}
	if add, _ := s.arc.Addit(c, aid); add != nil {
		a.MissionID = add.MissionID
		a.Source = add.Source
		a.OrderID = add.OrderID
		a.UpFrom = add.UpFrom
		if add.DescFormatID > 0 {
			a.Desc = add.Desc
		}
		a.DescFormatID = add.DescFormatID
		a.Dynamic = add.Dynamic
	}
	if s.isPorder(a) {
		if porder, _ := s.arc.Porder(c, aid); porder != nil {
			a.Porder = &archive.Porder{
				FlowID:     1,
				IndustryID: porder.IndustryID,
				BrandID:    porder.BrandID,
				BrandName:  porder.BrandName,
				Official:   porder.Official,
				ShowType:   porder.ShowType,
				Advertiser: porder.Advertiser,
				Agent:      porder.Agent,
				State:      porder.State,
			}
		}
	}
	// videos
	if avs, err = s.arc.NewVideos(c, aid); err != nil {
		log.Error("s.arc.Videos(%d) error(%v)", aid, err)
		return
	}
	if rsnm, err = s.arc.NewVideosReason(c, aid); err != nil {
		log.Error("s.Archive() s.arc.VideosReason(%d) error(%v)", aid, err)
		rsnm = map[int64]string{}
	}
	var vs []*archive.Video
	for _, v := range avs {
		if v.Status == archive.VideoStatusDelete {
			continue
		}
		if v.Status == archive.VideoStatusRecicle || v.Status == archive.VideoStatusLock || v.Status == archive.VideoStatusXcodeFail {
			v.RejectReason = rsnm[v.ID]
		}
		vs = append(vs, v)
	}
	av = &archive.ArcVideos{
		Archive: a,
		Videos:  vs,
	}
	return
}

func (s *Service) isPorder(a *archive.Archive) bool {
	if a == nil {
		return false
	}
	return a.AttrVal(archive.AttrBitIsPorder) == archive.AttrYes
}

func (s *Service) isStaff(a *archive.Archive) bool {
	if a == nil {
		return false
	}
	return a.AttrVal(archive.AttrBitSTAFF) == archive.AttrYes
}

func (s *Service) isDynamicForbid(a *archive.Archive) bool {
	if a == nil {
		return false
	}
	return a.AttrVal(archive.AttrBitNoPushBplus) == archive.AttrYes
}

func (s *Service) isPGC(a *archive.Archive) bool {
	if a == nil {
		return false
	}
	return a.AttrVal(archive.AttrBitIsPGC) == archive.AttrYes
}

// Archives get archives.
func (s *Service) Archives(c context.Context, aids []int64) (avm map[int64]*archive.ArcVideos, err error) {
	avm = map[int64]*archive.ArcVideos{}
	for _, aid := range aids {
		av, err := s.Archive(c, aid)
		if err != nil {
			log.Error("s.Archive(%d) error(%v)", aid, err)
			continue
		}
		avm[aid] = av
	}
	return
}

// UpArchives  get uper archives list by mid.
func (s *Service) UpArchives(c context.Context, mid int64, pn, ps int, group int8) (ua *archive.UpArchives, err error) {
	var (
		count  int64
		aids   []int64
		offset = (pn - 1) * ps
	)
	ua = &archive.UpArchives{
		Count: 0,
	}
	switch group {
	case archive.UpArcAllIn:
		if count, err = s.arc.ArchiveAllUpCount(c, mid); count <= 0 {
			return
		}
		if aids, err = s.arc.ArchivesUpAll(c, mid, offset, ps); len(aids) < 0 {
			return
		}
	case archive.UpArcOpenIn:
		if count, err = s.arc.ArchiveOpenUpCount(c, mid); count <= 0 {
			return
		}
		if aids, err = s.arc.ArchivesUpOpen(c, mid, offset, ps); len(aids) < 0 {
			return
		}
	case archive.UpArcUnOpenIn:
		if count, err = s.arc.ArchiveUnOpenUpCount(c, mid); count <= 0 {
			return
		}
		if aids, err = s.arc.ArchivesUpUnOpen(c, mid, offset, ps); len(aids) < 0 {
			return
		}
	}
	ua = &archive.UpArchives{
		Count: count,
		Aids:  aids,
	}
	return
}

// ArchivesByCids get archives by cids
func (s *Service) ArchivesByCids(c context.Context, cids []int64) (caids map[int64][]int64) {
	var avm map[int64]map[int64]*archive.Video
	caids = make(map[int64][]int64)
	if avm, _ = s.arc.NewVideosByCID(c, cids); len(avm) <= 0 {
		return
	}
	for _, cv := range avm {
		for ck, vv := range cv {
			caids[ck] = append(caids[ck], vv.Aid)
		}
	}
	return
}

// Flows get flow list.
func (s *Service) Flows(c context.Context) (fs []*archive.Flow) {
	fs = s.flowsCache
	return
}

// UpsByGroup Flows get group ups list.
func (s *Service) UpsByGroup(c context.Context, group int64) (ups []*archive.Up) {
	if group <= 0 {
		ups = s.specialUpsCache
		return
	}
	for _, up := range s.specialUpsCache {
		if group == up.GroupID {
			ups = append(ups, up)
		}
	}
	return ups
}

//UpsForbid get list for forbid
func (s *Service) UpsForbid(c context.Context) map[int64]string {
	return s.forbidMidsCache
}

// VideoBycid get video by cid
func (s *Service) VideoBycid(c context.Context, cid int64) (v *archive.Video, err error) {
	if v, _ = s.arc.NewVideoByCID(c, cid); v == nil {
		err = ecode.NothingFound
		return
	}
	return
}

// ArcHistory select archive edit history
func (s *Service) ArcHistory(c context.Context, hid int64) (ah *archive.ArcHistory) {
	if ah, _ = s.arc.ArcHistory(c, hid); ah.ID <= 0 {
		return
	}
	var vhs []*archive.VideoHistory
	if vhs, _ = s.arc.VideoHistory(c, hid); len(vhs) == 0 {
		return
	}
	ah.VideoHistory = vhs
	return
}

// ArcHistorys select archive edit historys.
func (s *Service) ArcHistorys(c context.Context, aid int64) (ahs []*archive.ArcHistory) {
	// archive
	stime := time.Now().Add(-time.Hour * 720).Unix()
	if ahs, _ = s.arc.ArcHistorys(c, aid, time.Unix(stime, 0)); len(ahs) == 0 {
		ahs = []*archive.ArcHistory{}
		return
	}
	var vhs []*archive.VideoHistory
	for _, ah := range ahs {
		if vhs, _ = s.arc.VideoHistory(c, ah.ID); len(vhs) == 0 {
			continue
		}
		ah.VideoHistory = vhs
	}
	return
}

// AppFeedAids get aids info by appFeed flow type.
func (s *Service) AppFeedAids(c context.Context) (aids []int64, err error) {
	var (
		nowTime   = time.Now()
		monTime   = nowTime.Add(-time.Hour * 720)
		startTime = time.Date(monTime.Year(), monTime.Month(), monTime.Day(), 0, 0, 0, 0, monTime.Location())
		endTime   = time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 23, 59, 59, 0, nowTime.Location())
	)
	if aids, err = s.arc.AppFeedAids(c, startTime, endTime); err != nil {
		log.Error("s.arc.AppFeedAids() error(%v)", err)
	}
	return
}

// DescFormats get desc_format info.
func (s *Service) DescFormats(c context.Context) (dfs []*archive.DescFormat, err error) {
	if dfs, err = s.arc.DescFormats(c); err != nil {
		log.Error("s.arc.DescFormats() error(%v)", err)
	}
	return
}

// DelArchive delete archive.
func (s *Service) DelArchive(c context.Context, aid, mid int64) (err error) {
	var a *archive.Archive
	if a, err = s.arc.Archive(c, aid); err != nil {
		log.Error("s.arc.Archive(%d) error(%v)", aid, err)
		return
	}
	if a == nil {
		err = ecode.NothingFound
		return
	}
	if a.Mid != mid {
		err = ecode.ArchiveNotExist
		return
	}
	if a.State == archive.StateForbidUpDelete {
		err = ecode.ArchiveAlreadyDel
		return
	}
	// 网警锁定，不允许删除
	if a.State == archive.StateForbidPolice {
		err = ecode.TargetBlocked
		return
	}
	var vs, dvs []*archive.Video
	if vs, err = s.arc.NewVideos(c, aid); err != nil {
		log.Error("s.arc.Videos(%d) error(%v)", aid, err)
		return
	}
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	log.Info("aid(%d) begin tran archive deleted", aid)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	if _, err = s.arc.TxUpArchiveState(tx, aid, archive.StateForbidUpDelete); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpArchiveState(%d,%d) error(%v)", aid, archive.StateForbidUpDelete, err)
		return
	}
	log.Info("aid(%d) filename() archive deleted by type(%d)", aid, archive.DelayTypeForUser)
	for _, v := range vs {
		if v.Status != archive.VideoStatusDelete && v.XcodeState <= archive.VideoXcodeSDFinish {
			v.Status = archive.VideoStatusDelete
			dvs = append(dvs, v)
		}
	}
	//no need to update archive_video or archive_video_relation status to -100 when delete archive envent

	//if err = s.txDelVideos(tx, dvs); err != nil {
	//	tx.Rollback()
	//	log.Error("s.arc.txDelVideos() error(%v)", err)
	//	return
	//}
	if _, err = s.arc.TxDelDelay(tx, aid, archive.DelayTypeForUser); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxDelDelay(%d,%d) error(%v)", aid, archive.DelayTypeForUser, err)
		return
	}
	log.Info("aid(%d) filename() delete delay by type(%d)", aid, archive.DelayTypeForUser)
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	log.Info("aid(%d) end tran archive deleted", aid)
	// NOTE: send delete_archive for sync dede.
	s.busDeleteArchive(aid)
	// NOTE: send delete_video for sync dede
	s.busDeleteVideo(dvs)
	return
}

// UpTag up archive tag.
func (s *Service) UpTag(c context.Context, aid int64, tag string) (err error) {
	var a *archive.Archive
	if a, err = s.arc.Archive(c, aid); err != nil {
		log.Error("s.arc.Archive(%d) error(%v)", aid, err)
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
	var tx *sql.Tx
	// begin tran
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	log.Info("aid(%d) begin tran archive up tag", aid)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	if _, err = s.arc.TxUpTag(tx, aid, tag); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpTag(%d,%s) error(%v)", aid, tag, err)
		return
	}
	log.Info("aid(%d) archive up tag(%s)", aid, tag)
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	log.Info("aid(%d) end tran archive up tag", aid)
	return
}

// checkArchive check archive changed.
func (s *Service) checkArchive(c context.Context, ap *archive.ArcParam) (a *archive.Archive, change, tagChange, typeIDChange bool, reason string, err error) {
	if a, err = s.arc.Archive(c, ap.Aid); err != nil {
		log.Error("s.arc.Archive(%d) error(%v)", ap.Aid, err)
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
	if ap.TypeID != a.TypeID {
		a.TypeID = ap.TypeID
		change = true
		typeIDChange = true
		reason += "分区|"
	}
	if ap.Copyright != a.Copyright {
		a.Copyright = ap.Copyright
		change = true
		reason += "类型|"
	}
	if ap.Title != a.Title {
		a.Title = ap.Title
		change = true
		reason += "标题|"
	}
	if ap.Cover != a.Cover {
		a.Cover = ap.Cover
		change = true
		reason += "封面|"
	}
	if ap.Tag != a.Tag {
		a.Tag = ap.Tag
		change = true
		tagChange = true
		reason += "标签|"
	}
	if ap.Desc != a.Desc {
		a.Desc = ap.Desc
		change = true
		reason += "简介|"
	}
	if ap.Source != a.Source {
		a.Source = ap.Source
		change = true
		reason += "转载来源|"
	}
	var (
		ad *archive.Addit
		dl *archive.Delay
	)
	if ad, _ = s.arc.Addit(c, ap.Aid); ad != nil {
		if ap.MissionID != ad.MissionID {
			change = true
			reason += "活动|"
		}
		if ap.OrderID != ad.OrderID {
			change = true
			reason += "商单|"
		}
		if ap.Dynamic != ad.Dynamic {
			change = true
			reason += "动态描述|"
		}
	}
	if dl, _ = s.arc.Delay(c, ap.Aid, archive.DelayTypeForUser); dl != nil {
		if ap.DTime != dl.DTime {
			change = true
			reason += "定时发布|"
		}
	}
	if change {
		reason = "稿件：" + reason[:len(reason)-1] + "\n"
	}
	return
}

// checkVideo check video changed when edit.
func (s *Service) checkVideo(c context.Context, ap *archive.ArcParam, isPGC bool) (nvs, evs, dvs, hvs []*archive.Video, change bool, reason string, err error) {
	var (
		vm      map[string]*archive.Video
		newNum  []int64
		editNum []int64
		fns     []string
		fnsMap  = make(map[string]int)
	)
	if vm, _, err = s.arc.NewVideoMap(c, ap.Aid); err != nil {
		log.Error("s.arc.VideoMap(%d) error(%v)", ap.Aid, err)
		return
	}
	// check videos
	for i, v := range ap.Videos {
		var (
			ov, ok        = vm[v.Filename]
			nv            = &archive.Video{}
			ovChg, idxChg bool
		)
		if !ok { // NOTE: add video
			nv.Aid = ap.Aid
			nv.Title = v.Title
			nv.Sid = v.Sid
			nv.Desc = v.Desc
			nv.Filename = v.Filename
			nv.Status = archive.VideoStatusSubmit
			nv.XcodeState = archive.VideoUploadInfo
			nv.SrcType = v.SrcType
			nv.Index = i + 1
			fns = append(fns, v.Filename)
			fnsMap[nv.Filename] = nv.Index
			// add videos
			nvs = append(nvs, nv)
			hvs = append(hvs, nv) // NOTE: ADD HISTORY
			newNum = append(newNum, int64(nv.Index))
			change = true

		} else { // NOTE: edit video
			*nv = *ov
			if nv.Title != v.Title {
				nv.Title = v.Title
				ovChg = true
				change = true
			}
			if nv.Desc != v.Desc {
				nv.Desc = v.Desc
				ovChg = true
				change = true
			}
			if i+1 != ov.Index {
				nv.Index = i + 1
				idxChg = true
			}
			if nv.Status == archive.VideoStatusLock || nv.Status == archive.VideoStatusXcodeFail || nv.Status == archive.VideoStatusSubmit {
				// nv.Status = nv.Status // NOTE: no changes, but highest priority.
			} else if nv.Status == archive.VideoStatusRecicle {
				if nv.XcodeState == archive.VideoXcodeSDFinish {
					nv.Status = archive.VideoStatusWait
					change = true
				} else {
					nv.Status = archive.VideoStatusOpen
				}
				ovChg = true
			} else if nv.Status == archive.VideoStatusDelete {
				if nv.XcodeState != 6 {
					newErr := ecode.VideoupVideoFilenameErr
					err = ecode.Errorf(newErr, newErr.Message(), nv.Index)
					return
				}
				nv.Status = archive.VideoStatusOpen
				ovChg = true
				change = true
			}
			//else if nv.Status == archive.VideoStatusWait || nv.Status == archive.VideoStatusAccess || nv.Status == archive.VideoStatusOpen {
			// nv.Status = nv.Status // NOTE: no changes, but lowest priority
			//}
			// edit videos
			if ovChg || idxChg {
				evs = append(evs, nv)
				editNum = append(editNum, int64(nv.Index))
			}
			hvs = append(hvs, nv) // NOTE: ADD HISTORY
			delete(vm, nv.Filename)
		}
	}
	log.Info("aid(%d) nvs filenames(%+v) ", ap.Aid, fns)
	if len(fns) != 0 {
		if !isPGC {
			if hasTimeout, expiredFilename, _ := s.arc.CheckNewVideosTimeout(c, fns); hasTimeout {
				log.Error("s.arc.CheckNewVideosTimeout(%+v) timeout", fns)
				newErr := ecode.VideoupFilenameExpired
				err = ecode.Errorf(newErr, newErr.Message(), fnsMap[expiredFilename])
				return
			}
		}
		//check new add videos
		var avm map[int64]map[string]*archive.Video
		if avm, err = s.arc.NewVideosByFn(c, fns); err != nil {
			log.Error("s.arc.NewVideosByFn(%+v) error(%v)", fns, err)
			return
		}
		if len(avm) > 0 {
			err = ecode.VideoAbnormalSubmit // NOTE: when not code mode, fns must have not db record.
			return
		}
	}
	// delete videos
	for _, v := range vm {
		if v.Status != archive.VideoStatusDelete {
			v.Status = archive.VideoStatusDelete
			dvs = append(dvs, v)
			change = true
		}
	}
	if len(newNum) != 0 {
		reason = "新增分P：" + xstr.JoinInts(newNum) + "\n"
	}
	if len(editNum) != 0 {
		reason = reason + "编辑分P：" + xstr.JoinInts(editNum) + "\n"
	}
	return
}

// upVideoCid update video cid.
func (s *Service) upVideoCid(c context.Context, aid int64, filename string, cid int64) (err error) {
	log.Info("aid(%d) filename(%s) cid(%d) begin update archive_video cid", aid, filename, cid)
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("begin tran error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	if _, err = s.arc.TxUpVideoCid(tx, aid, filename, cid); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpVideoCid  aid(%d) filename(%s) cid(%d) error(%v) ", aid, filename, cid, err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	log.Info("aid(%d) filename(%s) cid(%d) end update archive_video cid", aid, filename, cid)
	return
}

// upVideoHistoryCid update archive_edit_video_history cid.
func (s *Service) upVideoHistoryCid(c context.Context, aid, cid int64, filename string) (err error) {
	log.Info("aid(%d) filename(%s) cid(%d) begin update archive_edit_video_history cid", aid, filename, cid)
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("begin tran error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	if _, err = s.arc.TxUpVideoHistory(tx, aid, cid, filename); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpVideoHistory  aid(%d) filename(%s) cid(%d) error(%v) ", aid, filename, cid, err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	log.Info("aid(%d) filename(%s) cid(%d) end update archive_edit_video_history cid", aid, filename, cid)
	return
}

// Recos fn
func (s *Service) Recos(c context.Context, aid int64) (aids []int64, err error) {
	if aids, err = s.arc.Recos(c, aid); err != nil {
		log.Error("s.arc.Recos error(%v)|aid(%v) ", err, aid)
		return
	}
	return
}

// RecoUpdate fn
func (s *Service) RecoUpdate(c context.Context, mid, aid int64, recoIds []int64) (err error) {
	if newEffCnt, newErr := s.arc.RecoUpdate(c, aid, recoIds); newErr != nil {
		log.Warn("s.arc.RecoUpdate: newEffCnt (%+v)|newErr (%+v)", newEffCnt, newErr)
		return newErr
	}
	return
}

// VideoJamLevel  get video jam level
func (s *Service) VideoJamLevel(c context.Context) (level int, err error) {
	sec, err := s.arc.GetVideoJam(c)
	if err != nil || sec == 0 {
		level = 0
	} else if sec <= 1800 {
		level = 1
	} else if sec <= 3600 {
		level = 2
	} else if sec <= 6000 {
		level = 3
	} else if sec <= 13200 {
		level = 4
	} else {
		level = 5
	}
	return
}

// ArchiveAddit get archive addit
func (s *Service) ArchiveAddit(c context.Context, aid int64) (addit *archive.Addit, err error) {
	if aid == 0 {
		err = ecode.RequestErr
		return
	}
	return s.arc.Addit(c, aid)
}

// RejectedArchives list rejected archives
func (s *Service) RejectedArchives(c context.Context, mid int64, state, pn, ps int32, start *time.Time) (arcs []*archive.Archive, count int32, err error) {
	offset := (pn - 1) * ps
	arcs, count, err = s.arc.RejectedArchives(c, mid, state, offset, ps, start)
	if err != nil {
		log.Error("failed to get reject arcs,err(%+v)", err)
	}
	return
}
