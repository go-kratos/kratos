package service

import (
	"context"
	"strings"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AddByPGC add pgc archive and videos.
func (s *Service) AddByPGC(c context.Context, ap *archive.ArcParam) (aid int64, err error) {
	var (
		a   *archive.Archive
		nvs []*archive.Video
		fns = make([]string, 0, len(ap.Videos))
	)
	a = &archive.Archive{Mid: ap.Mid, Author: ap.Author, TypeID: ap.TypeID, Title: ap.Title, Cover: ap.Cover,
		Tag: ap.Tag, Copyright: ap.Copyright, Desc: ap.Desc, State: -30}
	// set attr
	a.AttrSet(ap.AllowTag, archive.AttrBitAllowTag)
	if len(ap.Videos) == 0 {
		err = ecode.VideoAbnormalSubmit
		return
	}
	for _, v := range ap.Videos {
		nv := &archive.Video{
			Title:    v.Title,
			Desc:     v.Desc,
			Filename: v.Filename,
			SrcType:  v.SrcType,
			Status:   archive.VideoStatusOpen,
		}
		if v.SrcType == "vupload" {
			nv.XcodeState = archive.VideoXcodeSDFinish
		} else {
			nv.XcodeState = archive.VideoDispatchFinish
		}
		fns = append(fns, v.Filename)
		nvs = append(nvs, nv)
	}
	log.Info("aid() filename(%s) start tran add pgc", strings.Join(fns, ","))
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
	if aid, err = s.txAddArchive(c, tx, a); err != nil {
		tx.Rollback()
		log.Error("s.add(%+v) error(%v)", a, err)
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
		return
	}
	log.Info("aid(%d) filename(%s) add videos", aid, strings.Join(fns, ","))
	if _, err = s.arc.TxUpAddit(tx, aid, 0, ap.OrderID, 0, 0, ap.IPv6, "", "", "", "", ap.Dynamic, ap.UpFrom); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpAddit(%d,%q,%d,%s,1) error(%v)", a.Aid, ap.IPv6, ap.OrderID, ap.Dynamic, err)
		return
	}
	log.Info("aid(%d) filename(%s) add archive_addit upfrom(%d) ipv6(%q)", aid, strings.Join(fns, ","), ap.UpFrom, ap.IPv6)
	if err = s.txAddHistory(tx, a, nvs); err != nil {
		tx.Rollback()
		log.Error("s.txAddHistory(%+v) error(%v)", a, err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	log.Info("aid(%d) filename(%s) end tran add pgc", aid, strings.Join(fns, ","))
	s.asyncCh <- func() {
		if ap.IsDRM == 1 {
			s.busDRMSubmit(nvs)
		} else {
			s.busPGCSubmit(nvs)
		}
		log.Info("aid(%d) send databus message pgc first round", a.Aid)
		s.syncCid(aid, ap.Mid, nvs, false)
		s.busAddArchive(aid, ap.UpFrom)
	}
	return
}

// EditByPGC edit pgc archive and videos.
func (s *Service) EditByPGC(c context.Context, ap *archive.ArcParam) (err error) {
	var (
		a                                         *archive.Archive
		aChange, vChange, tagChange, typeIDChange bool
		aReason, vReason                          string
		nvs, evs, dvs, hvs                        []*archive.Video
	)
	if a, aChange, tagChange, typeIDChange, aReason, err = s.checkArchive(c, ap); err != nil {
		log.Error("s.checkArchive(%+v) error(%v)", ap, err)
		return
	}
	if nvs, evs, dvs, hvs, vChange, vReason, err = s.checkVideo(c, ap, true); err != nil {
		log.Error("s.checkVideo(%+v) error(%v)", ap, err)
		return
	}
	log.Info("aid(%d) filename() cid(%d) start tran edit pgc aChange(%v) && vChange(%v)", ap.Aid, 0, aChange, vChange)
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
			log.Error("s.editArchvie(archive(%+v)) error(%v)", a, err)
			return
		}
	}
	for _, nv := range nvs {
		nv.Status = archive.VideoStatusOpen
		if nv.SrcType == "vupload" {
			nv.XcodeState = archive.VideoXcodeSDFinish
		} else {
			nv.XcodeState = archive.VideoDispatchFinish
		}
	}
	if err = s.assignCid(c, tx, nvs); err != nil {
		tx.Rollback()
		log.Error("s.assignCid(%d) error(%v)", ap.Aid, err)
		return
	}
	if err = s.txAddVideos(tx, nvs); err != nil {
		tx.Rollback()
		return
	}
	for _, ev := range evs {
		ev.Status = archive.VideoStatusOpen
	}
	if err = s.txEditVideos(tx, evs); err != nil {
		tx.Rollback()
		return
	}
	if err = s.txDelVideos(tx, dvs); err != nil {
		tx.Rollback()
		return
	}
	if _, err = s.arc.TxUpAddit(tx, a.Aid, 0, ap.OrderID, 0, 0, ap.IPv6, "", "", "", "", ap.Dynamic, ap.UpFrom); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpAddit(%d,%q,%d,%s,1) error(%v)", a.Aid, ap.IPv6, ap.OrderID, ap.Dynamic, err)
		return
	}
	log.Info("aid(%d) filename() cid(0) update archive_addit ipv6(%q) order_id(%d) upfrom(1) dynamic(%s)", ap.Aid, ap.IPv6, ap.OrderID, ap.Dynamic)
	if aChange || vChange {
		if _, err = s.arc.TxUpAdditReason(tx, a.Aid, aReason+vReason); err != nil {
			tx.Rollback()
			log.Error("s.arc.TxUpAdditReason(%d,%s) error(%v)", a.Aid, aReason+vReason, err)
			return
		}
		log.Info("aid(%d) filename() cid(0) update archive reason(%s)", ap.Aid, aReason+vReason)
	}
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
	s.asyncCh <- func() {
		if ap.IsDRM == 1 {
			s.busDRMSubmit(nvs)
		} else {
			s.busPGCSubmit(nvs)
		}
		log.Info("aid(%d) send databus message pgc first round", a.Aid)
		s.syncCid(ap.Aid, ap.Mid, nvs, false)
		s.busModifyArchive(ap.Aid, aChange, vChange, tagChange, len(nvs) > 0, typeIDChange)
		s.busDeleteVideo(dvs)
	}
	return
}
