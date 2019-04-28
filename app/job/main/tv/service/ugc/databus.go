package ugc

import (
	"context"
	"encoding/json"

	appDao "go-common/app/job/main/tv/dao/app"
	ugcmdl "go-common/app/job/main/tv/model/ugc"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_tableArchive = "archive"
	_updateAction = "update"
	_insertAction = "insert"
	_deleted      = 1
)

// arcConsumeproc consumer archive
func (s *Service) arcConsumeproc() {
	var err error
	defer s.waiter.Done()
	for {
		msg, ok := <-s.archiveNotifySub.Messages()
		if !ok {
			log.Info("arc databus Consumer exit")
			break
		}
		var ms = &ugcmdl.ArcMsg{}
		log.Info("arcConsumeproc New message: %s", msg)
		if err = json.Unmarshal(msg.Value, ms); err != nil {
			msg.Commit()
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		switch ms.Table {
		case _tableArchive:
			s.ArcHandle(ms)
		}
		msg.Commit()
	}
}

// ArcHandle treats the archive notify-T message to update the DB if there is any change
func (s *Service) ArcHandle(arcMsg *ugcmdl.ArcMsg) {
	var (
		msgMID  = arcMsg.New.Mid
		trustUp = false
	)
	if len(s.activeUps) > 0 { // check whether it's our trust upper
		if _, ok := s.activeUps[msgMID]; ok {
			trustUp = true
		}
	} else { // when the memory is not ready, check upper from DB
		mid, _ := s.dao.UpInList(ctx, arcMsg.New.Mid)
		trustUp = mid > 0
	}
	if !trustUp { // if it's not our trust upper, ignore the message except the archive was imported manually
		if arcMsg.Action == _updateAction && s.arcExist(arcMsg.New.Aid) { // if it's update and the archive exists ( added manually ), we allow it
			s.arcUpdate(arcMsg.Old, arcMsg.New)
			return
		}
		log.Info("Message Aid %d, Mid %d, Not in List, Ignore", arcMsg.New.Aid, arcMsg.New.Mid)
		appDao.PromInfo("DsInsert:Ignore")
		return
	}
	// arc update
	if arcMsg.Action == _updateAction {
		s.arcUpdate(arcMsg.Old, arcMsg.New)
	}
	// arc insert
	if arcMsg.Action == _insertAction {
		s.arcInsert(arcMsg.New)
	}
}

// distinguishes whether an arc exist
func (s *Service) arcExist(aid int64) bool {
	var (
		res *ugcmdl.Archive
		err error
	)
	if res, err = s.dao.ParseArc(ctx, aid); err != nil || res == nil {
		return false
	}
	if res.Deleted == _deleted {
		return false
	}
	return true
}

// arcInsert inserts a new databus notified archive
func (s *Service) arcInsert(arc *ugcmdl.ArchDatabus) (err error) {
	if exist := s.arcExist(arc.Aid); exist {
		appDao.PromError("DsInsert:Exist")
		log.Error("Databus Insert Data Aid %d Exist", arc.Aid)
		return
	}
	if err = s.importArc(context.Background(), arc.Aid, false); err != nil {
		appDao.PromError("DsInsert:Err")
		log.Error("Databus Import Arc %d Error %v", arc.Aid, err)
		return
	}
	appDao.PromInfo("DsInsert:Succ")
	return
}

// arcUpdate updates a databus notified archive
func (s *Service) arcUpdate(old *ugcmdl.ArchDatabus, new *ugcmdl.ArchDatabus) (err error) {
	if !s.arcExist(new.Aid) { // if an archive is not existing yet in our DB, we insert it
		return s.arcInsert(new)
	}
	new.Cover = s.coverURL(new.Cover, s.c.UgcSync.Cfg.BFSPrefix)
	var (
		oldAllow = &ugcmdl.ArcAllow{}
		newAllow = &ugcmdl.ArcAllow{}
	)
	oldAllow.FromDatabus(old)
	newAllow.FromDatabus(new)
	if !oldAllow.CanPlay() && newAllow.CanPlay() { // if an archive is recovered, re-insert it
		log.Info("Aid %d is recovered, add it", new.Aid)
		return s.arcInsert(new)
	}
	if oldAllow.CanPlay() && !newAllow.CanPlay() { // if an archive is banned, delete it
		log.Info("Aid %d can't play, delete it", new.Aid)
		if err = s.dao.UpdateArc(ctx, new); err != nil {
			return
		}
		return s.delArc(new.Aid)
	}
	// if arc level changed or video level changed, treat and import data
	return s.arcCompare(old, new)
}

// arcCompare compares the archive & the videos of the old and the new, to update if needed
func (s *Service) arcCompare(old *ugcmdl.ArchDatabus, new *ugcmdl.ArchDatabus) (err error) {
	var (
		diff   *ugcmdl.VideoDiff
		hitPGC bool
	)
	if hitPGC, err = s.delPGC(new.TypeID, new.Aid); err != nil {
		return
	}
	if hitPGC { // if the archive hits PGC types, delete it
		log.Warn("arcCompare Del Aid %d, Because of its typeID %d", new.Aid, new.TypeID)
		return
	}
	if s.diffArc(old, new) { // archive level info update if different
		if err = s.dao.UpdateArc(ctx, new); err != nil {
			appDao.PromError("DsUpdArc:Err")
			return
		}
		s.modArcCh <- []int64{new.Aid} // add one archive to submit
		appDao.PromInfo("DsUpdArc:Succ")
	}
	// video level info update if different
	if diff, err = s.diffVideos(new.Aid); err != nil {
		appDao.PromError("DsUpdVideo:Err")
		return
	}
	log.Info("Diff Result For Aid %d, Equal %v, Updated %v, Removed %v, New %v", new.Aid, diff.Equal, diff.Updated, diff.Removed, diff.New)
	if err = s.treatDiffV(diff); err != nil {
		appDao.PromError("DsUpdVideo:Err")
		return
	}
	appDao.PromInfo("DsUpdVideo:Succ")
	return
}

// get first level of types name
func (s *Service) getPTypeName(typeID int32) (name string) {
	var (
		second, first *arcmdl.Tp
		ok            bool
	)
	if second, ok = s.arcTypes[typeID]; !ok {
		log.Error("can't find type for ID: %d ", typeID)
		return
	}
	if first, ok = s.arcTypes[second.Pid]; !ok {
		log.Error("can't find type for ID: %d, second Info: %v", second, second.Pid)
		return
	}
	return first.Name
}

// getPType first level of types name
func (s *Service) getPType(typeID int) (pid int) {
	var (
		second *arcmdl.Tp
		ok     bool
	)
	if second, ok = s.arcTypes[int32(typeID)]; !ok {
		log.Error("can't find type for ID: %d ", typeID)
		return
	}
	return int(second.Pid)
}

// diffArc distinguishes whether the key fields of an archive have been changed
func (s *Service) diffArc(old *ugcmdl.ArchDatabus, new *ugcmdl.ArchDatabus) (diff bool) {
	diff = (old.Title != new.Title)
	diff = diff || (old.Content != new.Content)
	diff = diff || (old.PubTime != new.PubTime)
	diff = diff || (old.TypeID != new.TypeID)
	diff = diff || (old.Cover != new.Cover)
	diff = diff || (s.getPTypeName(old.TypeID) != s.getPTypeName(new.TypeID))
	return
}

// diffVideos distinguishes whethe
func (s *Service) diffVideos(aid int64) (diff *ugcmdl.VideoDiff, err error) {
	var (
		rpcRes *arcmdl.ViewReply
		dbRes  map[int64]*ugcmdl.SimpleVideo
		video  *ugcmdl.SimpleVideo
		ok     bool
	)
	diff = &ugcmdl.VideoDiff{
		Aid: aid,
	}
	if rpcRes, err = s.videoPick(ctx, aid); err != nil {
		log.Error("rpc video pick %d, error %v", aid, err)
		return
	}
	if dbRes, err = s.dao.PickVideos(ctx, aid); err != nil {
		log.Error("db video pick %d, error %v", aid, err)
		return
	}
	for _, page := range rpcRes.Pages {
		if video, ok = dbRes[page.Cid]; !ok { // not found in DB, means it's new
			diff.New = append(diff.New, page.Cid)
			continue
		}
		if video.IndexOrder == int64(page.Page) && video.Eptitle == page.Part { // if title & index_order equal
			diff.Equal = append(diff.Equal, page.Cid)
		} else { // otherwise it's updated
			diff.Updated = append(diff.Updated, page)
		}
		delete(dbRes, page.Cid)
	}
	for _, v := range dbRes {
		diff.Removed = append(diff.Removed, v.CID)
	}
	return
}

// treatDiffV treats the result of diffVideos, like we add new ones, we deleted removed ones, and we updated the modified ones
func (s *Service) treatDiffV(diff *ugcmdl.VideoDiff) (err error) {
	var (
		newPages []*arcmdl.Page
		page     *arcmdl.Page
		aid      = diff.Aid
		tx       *sql.Tx
	)
	// all the operations about this archive's videos, will be in one transaction
	if tx, err = s.dao.BeginTran(ctx); err != nil {
		log.Error("BeginTran Error %v", err)
		return
	}
	// add new videos
	if len(diff.New) > 0 {
		for _, v := range diff.New {
			if page, err = s.pagePick(ctx, v, aid, ""); err != nil {
				continue
			}
			newPages = append(newPages, page)
		}
		if err = s.dao.TxAddVideos(tx, newPages, aid); err != nil {
			tx.Rollback()
			return
		}
	}
	if len(diff.Removed) > 0 {
		for _, v := range diff.Removed {
			if err = s.dao.TxDelVideo(tx, v); err != nil {
				tx.Rollback()
				return
			}
		}
	}
	if len(diff.Updated) > 0 {
		for _, v := range diff.Updated {
			if err = s.dao.TxUpdateVideo(tx, v); err != nil {
				tx.Rollback()
				return
			}
		}
	}
	tx.Commit()
	return
}
