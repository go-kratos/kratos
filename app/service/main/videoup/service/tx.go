package service

import (
	"context"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
	"time"
)

// txAddArchive add archive.
func (s *Service) txAddArchive(c context.Context, tx *sql.Tx, a *archive.Archive) (aid int64, err error) {
	if aid, err = s.arc.TxAddArchive(tx, a); err != nil {
		log.Error("s.arc.TxAddArchive(%v) error(%v)", a, err)
		return
	}
	a.Aid = aid
	log.Info("aid(%d) create archive(%+v)", aid, a)
	return
}

// txEditArchvie edit archvie.
func (s *Service) txEditArchvie(tx *sql.Tx, a *archive.Archive) (err error) {
	if _, err = s.arc.TxUpArchive(tx, a); err != nil {
		log.Error("s.arc.TxUpArchive(%d) error(%v)", a.Aid, err)
		return
	}
	log.Info("aid(%d) edit archive(%+v)", a.Aid, a)
	return
}

// txAddVideos add videos.
func (s *Service) txAddVideos(tx *sql.Tx, nvs []*archive.Video) (err error) {
	for _, v := range nvs {
		if _, err = s.arc.TxUpNewVideo(tx, v); err != nil {
			log.Error("s.arc.TxUpNewVideo(%+v) error(%v)", v, err)
			return
		}
		if v.ID, err = s.arc.TxAddRelation(tx, v); err != nil {
			log.Error("s.arc.TxAddArcVideoRelation(%+v) error(%v)", v, err)
			return
		}
		if _, err = s.arc.TxAddVideo(tx, v); err != nil {
			log.Error("s.arc.TxAddVideo(%+v) error(%v) or rows==0", v, err)
			return
		}
		log.Info("aid(%d) filename(%s) create video(%+v)", v.Aid, v.Filename, v)
	}
	if len(nvs) != 0 {
		if _, err = s.arc.TxAddAudit(tx, nvs); err != nil {
			log.Error("s.arc.TxAddAudit(%v) error(%v)", nvs, err)
			return
		}
	}
	return
}

// txEditVideos edit video.
func (s *Service) txEditVideos(tx *sql.Tx, evs []*archive.Video) (err error) {
	for _, v := range evs {
		if _, err = s.arc.TxUpVideo(tx, v); err != nil {
			log.Error("s.arc.TxUpVideo(%v) error(%v)", v, err)
			return
		}
		if _, err = s.arc.TxUpVideoRelation(tx, v); err != nil {
			log.Error("s.arc.TxUpVideoRelation(%v) error(%v)", v, err)
			return
		}
		log.Info("aid(%d) filename(%s) cid(%d) video edit", v.Aid, v.Filename, v.Cid)
		if _, err = s.arc.TxUpVideoStatus(tx, v.Aid, v.Filename, v.Status); err != nil {
			log.Error("s.arc.TxUpVideoStatus(%d,%s,%d) error(%v)", v.Aid, v.Filename, v.Status, err)
			return
		}
		if _, err = s.arc.TxUpVdoStatus(tx, v.Cid, v.Status); err != nil {
			log.Error("s.arc.TxUpVdoStatus(%d,%d) error(%v)", v.Cid, v.Status, err)
			return
		}
		log.Info("aid(%d) filename(%s) cid(%d) edit video change status(%d)", v.Aid, v.Filename, v.Cid, v.Status)
	}
	return
}

// txEditVideosAsync edit video.
func (s *Service) txEditVideosAsync(tx *sql.Tx, evs []*archive.Video) (err error) {
	for i, v := range evs {
		log.Info("aid(%d) i(%d) edit video(%+v)", v.Aid, i, v)
		time.Sleep(5 * time.Millisecond)
		if _, err = s.arc.TxUpVideo(tx, v); err != nil {
			log.Error("s.arc.TxUpVideo(%v) error(%v)", v, err)
			return
		}
		if _, err = s.arc.TxUpVideoRelation(tx, v); err != nil {
			log.Error("s.arc.TxUpVideoRelation(%v) error(%v)", v, err)
			return
		}
		log.Info("aid(%d) filename(%s) cid(%d) video edit", v.Aid, v.Filename, v.Cid)
		if _, err = s.arc.TxUpVideoStatus(tx, v.Aid, v.Filename, v.Status); err != nil {
			log.Error("s.arc.TxUpVideoStatus(%d,%s,%d) error(%v)", v.Aid, v.Filename, v.Status, err)
			return
		}
		if _, err = s.arc.TxUpVdoStatus(tx, v.Cid, v.Status); err != nil {
			log.Error("s.arc.TxUpVdoStatus(%d,%d) error(%v)", v.Cid, v.Status, err)
			return
		}
	}
	return
}

// txDelVideos del videos.
func (s *Service) txDelVideos(tx *sql.Tx, dvs []*archive.Video) (err error) {
	for _, v := range dvs {
		if _, err = s.arc.TxUpVideoStatus(tx, v.Aid, v.Filename, v.Status); err != nil {
			log.Error("s.arc.TxUpVideoStatus((%d,%s,%d) error(%v)", v.Aid, v.Filename, v.Status, err)
			return
		}
		if _, err = s.arc.TxUpRelationState(tx, v.Aid, v.Cid, v.Status); err != nil {
			log.Error("s.arc.TxUpRelationState((%d,%d,%d) error(%v)", v.Aid, v.Cid, v.Status, err)
			return
		}
		log.Info("aid(%d) filename(%s) cid(%d) xcodeState(%d) video delete", v.Aid, v.Filename, v.Cid)
	}
	return
}

// txDelayTime add or del delay time.
func (s *Service) txDelayTime(tx *sql.Tx, mid, aid int64, dTime xtime.Time) (err error) {
	if dTime == 0 {
		if _, err = s.arc.TxDelDelay(tx, aid, archive.DelayTypeForUser); err != nil {
			log.Error("s.arc.TxUpDelay(%d,%d) error(%v)", aid, archive.DelayTypeForUser, err)
			return
		}
		log.Info("aid(%d) delete delay", aid)
	} else {
		if _, err = s.arc.TxUpDelay(tx, mid, aid, archive.StateOpen, archive.DelayTypeForUser, dTime); err != nil {
			log.Error("s.arc.TxUpDelay(%d,%d,%v) error(%v)", aid, archive.DelayTypeForUser, dTime, err)
			return
		}
		log.Info("aid(%d) update delay(%s)", aid, dTime.Time().Format("2006-01-02 15:04:05"))
	}
	return
}

// txAddHistory  insert Archive edit history
func (s *Service) txAddHistory(tx *sql.Tx, a *archive.Archive, vs []*archive.Video) (err error) {
	var hid int64
	if hid, err = s.arc.TxAddArcHistory(tx, a.Aid, a.Mid, a.Title, a.Desc, a.Cover, a.Tag); err != nil {
		log.Error("s.arc.TxAddArcHistory(%d,%d,%s,%s,%s,%s error(%v)", a.Aid, a.Mid, a.Title, a.Desc, a.Cover, a.Tag, err)
		return
	}
	log.Info("aid(%d) mid(%d) title(%s) desc(%s) cover(%s) tag(%s) insert archive edit history and get hid(%d)", a.Aid, a.Mid, a.Title, a.Desc, a.Cover, a.Tag, hid)
	if len(vs) > 0 {
		if err = s.arc.TxAddVideoHistorys(tx, hid, vs); err != nil {
			log.Error("s.arc.TxAddVideoHistorys hid(%d)|vs(%+v)|vsCnt(%d)|error(%v)", hid, vs, len(vs), err)
			return
		}
	}
	return
}
