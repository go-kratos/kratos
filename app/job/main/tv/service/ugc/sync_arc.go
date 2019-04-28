package ugc

import (
	"context"
	"time"

	"go-common/app/job/main/tv/dao/lic"
	model "go-common/app/job/main/tv/model/pgc"
	ugcmdl "go-common/app/job/main/tv/model/ugc"
	"go-common/library/log"
)

// syncLic: sync our arc data to License owner
func (s *Service) modArcproc() (err error) {
	defer s.waiter.Done()
	var cid int64
	for {
		cAids, ok := <-s.modArcCh
		if !ok {
			log.Warn("[modLic] channel quit")
			return
		}
		for _, cAid := range cAids {
			if cid, err = s.dao.VideoSubmit(ctx, cAid); err != nil {
				log.Warn("modArc Aid %d, Err %v, Jump", cAid, err)
				time.Sleep(time.Duration(s.c.UgcSync.Frequency.ErrorWait))
				continue
			}
			log.Info("modArc Aid %d, Can submit because CID %d already submitted", cAid, cid)
			if err = s.modArc(ctx, cAid); err != nil {
				log.Warn("modArc Aid %d Err %v", cAid, err)
				continue
			}
		}
	}
}

func (s *Service) modArc(ctx context.Context, cAid int64) (err error) {
	var (
		skeleton = &ugcmdl.LicSke{}
		licData  *model.License
		xmlBody  string
		arc      *ugcmdl.Archive
	)
	if arc, err = s.dao.ParseArc(ctx, cAid); err != nil {
		log.Warn("ParseArc Aid %d not found", cAid)
	}
	skeleton.Arc = arc.ToSimple()
	skeleton.Videos = []*ugcmdl.SimpleVideo{} // empty videos
	// build the license data and transform to xml
	if licData, err = s.auditMsg(skeleton); err != nil {
		log.Error("build lic msg %d error %v", cAid, err)
		return
	}
	xmlBody = lic.PrepareXML(licData)
	// call api
	if _, err = s.licDao.CallRetry(ctx, s.c.Sync.API.AddURL, xmlBody); err != nil {
		log.Error("xml call %d error %v", cAid, err)
		return
	}
	// update the arc & videos' submit status to finish
	if err = s.dao.FinishArc(ctx, cAid); err != nil {
		log.Error("finishArc %d Error %v", cAid, err)
	}
	return
}

// wrapSyncLic warps the syncLic method with aidMap
func (s *Service) wrapSyncLic(ctx context.Context, aids []int64) (err error) {
	var arc *ugcmdl.Archive
	for _, cAid := range aids {
		if arc, err = s.dao.ParseArc(ctx, cAid); err != nil {
			log.Warn("wrapSyncLic ParseArc Aid %d not found", cAid)
			continue
		}
		arcAllow := &ugcmdl.ArcAllow{}
		arcAllow.FromArchive(arc)
		if !s.arcAllowImport(arcAllow) {
			log.Warn("wrapSyncLic cAid %d Can't play", cAid)
			continue
		}
		if arc.Deleted == 1 {
			log.Warn("wrapSyncLic cAid %d Deleted", cAid)
			continue
		}
		if err = s.syncLic(cAid, arc.ToSimple()); err != nil {
			log.Error("wrapSyncLic cAid %d Err %v", cAid, err)
			continue
		}
	}
	return
}
