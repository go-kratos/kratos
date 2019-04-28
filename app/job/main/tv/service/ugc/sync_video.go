package ugc

import (
	"fmt"

	"go-common/app/job/main/tv/dao/lic"
	model "go-common/app/job/main/tv/model/pgc"
	ugcmdl "go-common/app/job/main/tv/model/ugc"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_crEnd      = "1970-01-01" // copyright end date
	_definition = "SD"
)

// syncVideoErr: it logs the error and postpone the videos for the next submit
func (s *Service) syncVideoErr(funcName string, cids []int64, aid int64, err error) {
	s.dao.PpVideos(ctx, cids)
	errArcVideos("syncLic:"+funcName, aid, cids, err)
}

// syncLic: sync our arc data to License owner
func (s *Service) syncLic(cAid int64, arc *ugcmdl.SimpleArc) (err error) {
	var (
		skeleton  = &ugcmdl.LicSke{}
		videoPces [][]*ugcmdl.SimpleVideo
		ps        = s.c.UgcSync.Batch.SyncPS // sync page size
		licData   *model.License
		xmlBody   string
		errCall   error
	)
	skeleton.Arc = arc
	if videoPces, err = s.dao.ParseVideos(ctx, cAid, ps); err != nil {
		log.Error("ParseVideos %d Error %v", cAid, err)
		return
	}
	if len(videoPces) == 0 { // no to audit cids
		return
	}
	for _, videos := range videoPces {
		skeleton.Videos = videos
		var cids = []int64{}
		for _, v := range videos {
			cids = append(cids, v.CID)
		}
		if licData, errCall = s.auditMsg(skeleton); errCall != nil { // build the license data and transform to xml
			s.syncVideoErr("AuditMsg ", cids, cAid, errCall)
			continue
		}
		xmlBody = lic.PrepareXML(licData)
		if _, errCall = s.licDao.CallRetry(ctx, s.c.Sync.API.AddURL, xmlBody); errCall != nil { // call api
			s.syncVideoErr("XmlCall ", cids, cAid, errCall)
			continue
		}
		if errCall = s.dao.FinishVideos(ctx, skeleton.Videos, cAid); errCall != nil { // update the arc & videos' submit status to finish
			s.syncVideoErr("FinishVideos ", cids, cAid, errCall)
			continue
		}
		infoArcVideos("syncLic", cAid, cids, "Succ Apply For Audit")
	}
	return
}

// auditMsg transforms a skeleton to license audit message struct for UGC
func (s *Service) auditMsg(skeleton *ugcmdl.LicSke) (licData *model.License, err error) {
	var (
		programSets []*model.PS
		programs    []*model.Program
		sign        = s.c.Sync.Sign
	)
	if len(skeleton.Videos) > 0 {
		if programs, err = s.videoProgram(skeleton.Arc.AID, skeleton.Videos); err != nil {
			log.Error("auditMsg videoProgram Aid %d, Err %v", skeleton.Arc.AID, err)
			return
		}
	}
	if programSets, err = s.arcPSet(skeleton.Arc, programs); err != nil {
		log.Error("arcPSet Error %v", err)
		return
	}
	licData = lic.BuildLic(sign, programSets, len(programs))
	return
}

// videoProgram transforms the videos to license defined program models
func (s *Service) videoProgram(aid int64, videos []*ugcmdl.SimpleVideo) (programs []*model.Program, err error) {
	var (
		ugcPrefix = s.c.Sync.UGCPrefix
		deadCIDs  = []int64{}
		arcValid  bool
	)
	for _, v := range videos {
		playurl, hitDead, errCall := s.playurlDao.Playurl(ctx, int(v.CID))
		if errCall != nil {
			log.Error("Playurl CID %d, Error %v", v.CID, errCall)
			continue
		}
		if hitDead { // hit playurl dead codes
			deadCIDs = append(deadCIDs, v.CID)
			continue
		}
		media := model.MakePMedia(ugcPrefix, playurl, v.CID)
		program := &model.Program{
			ProgramID:      fmt.Sprintf("%s%d", ugcPrefix, v.CID),
			ProgramName:    v.Eptitle,
			ProgramLength:  int(v.Duration),
			ProgramDesc:    v.Description,
			PublishDate:    _crEnd,
			Number:         fmt.Sprintf("%d", v.IndexOrder),
			DefinitionType: "SD",
			ProgramMediaList: &model.PMList{
				ProgramMedia: []*model.PMedia{
					media,
				},
			},
		}
		programs = append(programs, program)
	}
	if len(deadCIDs) > 0 { // treat deadCIDs, delete them
		if arcValid, err = s.dao.DelVideoArc(ctx, &ugcmdl.DelVideos{
			AID:  aid,
			CIDs: deadCIDs,
		}); err != nil {
			log.Error("VideoProgram DelVideos Aid %d, Cids %v, Err %v", aid, deadCIDs, err)
			return
		}
		if !arcValid {
			log.Info("VideoProgram DelVideos Aid %d is empty, delete it also", aid)
			err = ecode.NothingFound
			return
		}
		if len(deadCIDs) == len(videos) {
			log.Info("VideoProgram Passed Videos Aid %d are dead, Cids %v", aid, deadCIDs)
			err = ecode.NothingFound
			return
		}
		log.Info("VideoProgram DelVideos Aid %d, Playurl DeadCids %v", aid, deadCIDs)
	}
	return
}

// arcPSet transforms an archive model to a license programSet model
func (s *Service) arcPSet(arc *ugcmdl.SimpleArc, programs []*model.Program) (ps []*model.PS, err error) {
	var (
		secondType string
		firstType  string
		copyright  = s.c.UgcSync.Cfg.Copyright
		upper      *ugcmdl.Upper
	)
	// get second type name
	if tp, ok := s.arcTypes[arc.TypeID]; !ok {
		log.Error("For Aid %d, Can't find Second TypeID %d Name", arc.AID, arc.TypeID)
	} else {
		secondType = tp.Name
	}
	// get first type name
	firstType = s.getPTypeName(arc.TypeID)
	// build the programSet structure
	var program = &model.PS{
		ProgramSetID:     fmt.Sprintf("%s%d", s.c.Sync.UGCPrefix, arc.AID),
		ProgramSetName:   arc.Title,
		ProgramSetClass:  secondType,
		ProgramSetType:   firstType,
		PublishDate:      arc.Pubtime,
		Copyright:        copyright,
		ProgramCount:     int(arc.Videos),
		CREndData:        _crEnd,
		DefinitionType:   _definition,
		CpCode:           s.c.Sync.LConf.CPCode,
		PayStatus:        0,
		ProgramSetDesc:   arc.Content,
		ProgramSetPoster: arc.Cover,
		ProgramList: &model.ProgramList{
			Program: programs,
		},
	}
	// upper info
	if upper, err = s.upDao.LoadUpMeta(ctx, arc.MID); err != nil { // get upper meta info
		log.Error("modLic LoadUpMeta Aid %d, Mid %d, Err %v", arc.AID, arc.MID, err)
		err = nil
	}
	if upper != nil {
		program.Producer = upper.OriName
		program.Portrait = upper.OriFace
	}
	ps = append(ps, program)
	return
}
