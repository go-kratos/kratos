package ugc

import (
	"encoding/json"

	appDao "go-common/app/job/main/tv/dao/app"
	ugcmdl "go-common/app/job/main/tv/model/ugc"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// arcDatabus refreshes the mc cache for archive media info
func (s *Service) arcDatabus(jsonstr json.RawMessage) (err error) {
	var (
		arc     = &ugcmdl.DatabusArc{}
		pubtime int64
	)
	if err = json.Unmarshal(jsonstr, arc); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", jsonstr, err)
		return
	}
	arcMark := arc.New
	if pubtime, err = appDao.TimeTrans(arcMark.Pubtime); err != nil {
		log.Warn("arcDatabus Pubtime AVID: %d, Err %v", arcMark.AID, err)
	}
	// we prepare the cms cache
	if err = s.dao.SetArcCMS(ctx, &ugcmdl.ArcCMS{
		// Media Info
		Title:   arcMark.Title,
		AID:     arcMark.AID,
		Content: arcMark.Content,
		Cover:   arcMark.Cover,
		TypeID:  arcMark.TypeID,
		Pubtime: xtime.Time(pubtime),
		Videos:  arcMark.Videos,
		Valid:   arcMark.Valid,
		Deleted: arcMark.Deleted,
		Result:  arcMark.Result,
	}); err != nil {
		log.Error("arcDatabus setArcCMS AVID: %d, Err %v", arcMark.AID, err)
	}
	// we prepare the rpc cache for the ugc view page if the archive is able to play
	if arcMark.IsPass() {
		s.viewCache(int64(arcMark.AID))
		appDao.PromInfo("ArcRPC-AddCache")
	}
	s.listMtn(arc.Old, arc.New)
	return
}

// videoDatabus refreshes the mc cache for video media info
func (s *Service) videoDatabus(jsonstr json.RawMessage) (err error) {
	var (
		video  = &ugcmdl.DatabusVideo{}
		criCID = s.c.UgcSync.Cfg.CriticalCid
	)
	if err = json.Unmarshal(jsonstr, video); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", jsonstr, err)
		return
	}
	vm := video.New
	if vm.ToReport(criCID) { // if the video has not been reported yet, we do it and update the mark field from 0 to 1
		s.repCidCh <- vm.CID
	}
	if vm.ToAudit(criCID) {
		log.Info("videoDatabus addAudCid cAid %d", vm.AID)
		s.audAidCh <- []int64{vm.AID} // add aid into channel to treat
	}
	if video.Old == nil { // if the brand new episode can play
		if vm.CanPlay() {
			log.Info("videoDatabus reshelfAid cAid %d", vm.AID)
			s.reshelfAidCh <- vm.AID
		}
	} else { // or it couldn't play and it passes now
		if !video.Old.CanPlay() && vm.CanPlay() {
			log.Info("videoDatabus reshelfAid cAid %d", vm.AID)
			s.reshelfAidCh <- vm.AID
		}
	}
	if err = s.dao.SetVideoCMS(ctx, vm.ToCMS()); err != nil { // we prepare the cms cache
		log.Warn("videoDatabus setVideoCMS CID: %d, Err %v", vm.CID, err)
	}
	return
}
