package pgc

import (
	"context"
	"time"

	"go-common/app/job/main/tv/dao/lic"
	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/ecode"
	"go-common/library/log"
)

// pick the data from DB to audit and combine the XML for the license owner
// producer, content data => channel
func (s *Service) syncEPs() {
	defer s.waiter.Done()
	for {
		if s.daoClosed {
			log.Info("syncEPs DB closed!")
			return
		}
		readySids, err := s.dao.ReadySns(ctx)
		if err != nil || len(readySids) == 0 {
			time.Sleep(time.Duration(s.c.Sync.Frequency.ErrorWait))
			continue
		}
		for _, sid := range readySids {
			var contSlices [][]*model.Content
			if contSlices, err = s.dao.PickData(ctx, sid); err != nil || len(contSlices) == 0 {
				continue
			}
			for _, conts := range contSlices {
				if err = s.epsSync(sid, conts); err != nil {
					s.addRetryEps(conts)
				}
				s.dao.AuditingCont(ctx, conts) // update status to auditing
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (s *Service) epsSync(sid int64, conts []*model.Content) (err error) {
	var reqCall = &model.ReqEpLicCall{
		SID:   sid,
		Conts: conts,
	}
	if reqCall.EpLic, err = s.epLicCreate(ctx, sid, conts); err != nil {
		return
	}
	return s.epLicCall(ctx, reqCall)
}

// epLicCreate picks the sid and conts to create the license model
func (s *Service) epLicCreate(ctx context.Context, sid int64, conts []*model.Content) (epLic *model.License, err error) {
	var (
		season   *model.TVEpSeason
		prefix   = s.c.Sync.AuditPrefix
		programs []*model.Program
	)
	if season, err = s.dao.Season(ctx, int(sid)); err != nil {
		log.Error("Season ID %d, Err %v", sid, err)
		return
	}
	epLic = newLic(season, s.c.Sync)
	epLic.XMLData.Service.Head.Count = len(conts)
	for _, v := range conts {
		s.dao.WaitCall(ctx, v.EPID) // avoid always selecting the same data, give time to the caller
		url, _, errPlay := s.playurlDao.Playurl(ctx, v.CID)
		if errPlay != nil {
			log.Error("syncEPs EP Playurl EPID = %d, Error: %v", v.EPID, errPlay)
			s.addRetryEp(v)
			continue
		}
		ep, errEP := s.dao.EP(ctx, v.EPID)
		if errEP != nil {
			log.Error("EpContent EPID %d Can't found", v.EPID)
			continue
		}
		program := model.CreateProgram(prefix, ep)
		program.ProgramMediaList = &model.PMList{
			ProgramMedia: []*model.PMedia{model.CreatePMedia(s.c.Sync.AuditPrefix, v.EPID, url)},
		}
		programs = append(programs, program)
	}
	epLic.XMLData.Service.Body.ProgramSetList.ProgramSet[0].ProgramList.Program = programs
	return
}

// epLicCall picks the license and sync to audit
func (s *Service) epLicCall(ctx context.Context, req *model.ReqEpLicCall) (err error) {
	var cfg = s.c.Sync
	res, err := s.licDao.CallRetry(ctx, cfg.API.AddURL, lic.PrepareXML(req.EpLic))
	if res == nil {
		err = ecode.TvSyncErr
	}
	return
}
