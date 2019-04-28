package report

import (
	"bytes"
	"context"

	"go-common/app/interface/main/app-view/conf"
	arcdao "go-common/app/interface/main/app-view/dao/archive"
	elcdao "go-common/app/interface/main/app-view/dao/elec"
	reportdao "go-common/app/interface/main/app-view/dao/report"
	"go-common/app/interface/main/app-view/model"
	"go-common/app/interface/main/app-view/model/elec"
	"go-common/app/interface/main/app-view/model/report"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
)

const (
	_maxContentSize = 400
	_maxFileSize    = 1048576
)

var (
	_elecTypeIds = []int16{
		20, 154, 156, // dance
		31, 30, 59, 29, 28, // music
		26, 22, 126, 127, // guichu
		24, 25, 47, 27, // animae
		17, 18, 16, 65, 136, 19, 121, 171, 172, 173, // game
		37, 124, 122, 39, 96, 95, 98, // tech
		71, 137, 131, // yule
		157, 158, 159, 164, // fashion
		82, 128, // movie and tv
		138, 21, 75, 76, 161, 162, 163, 174, // life
		153, 168, // guo man
		85, 86, 182, 183, 184, // film and television
	}
)

// Service is appeal service .
type Service struct {
	arcDao    *arcdao.Dao
	reportDao *reportdao.Dao
	elcDao    *elcdao.Dao
	// elec
	allowTypeIds map[int16]struct{}
}

// New init appeal service .
func New(c *conf.Config) (s *Service) {
	s = &Service{
		reportDao: reportdao.New(c),
		arcDao:    arcdao.New(c),
		elcDao:    elcdao.New(c),
	}
	s.allowTypeIds = map[int16]struct{}{}
	for _, id := range _elecTypeIds {
		s.allowTypeIds[id] = struct{}{}
	}
	return
}

// CopyWriter archive appeal copy write .
func (s *Service) CopyWriter(c context.Context, aid int64, plat int8, lang string) (cps []report.CopyWriter, err error) {
	var a *api.Arc
	if a, err = s.arcDao.Archive(c, aid); err != nil {
		return
	}
	if a == nil {
		err = ecode.NothingFound
		return
	}
	// check elec
	var info *elec.Info
	if _, ok := s.allowTypeIds[int16(a.TypeID)]; ok && a.IsNormal() && a.Copyright == int32(archive.CopyrightOriginal) {
		if info, err = s.elcDao.Info(c, a.Author.Mid, 0); err != nil {
			return
		}
	}
	if info != nil {
		if model.IsOverseas(plat) {
			cps = make([]report.CopyWriter, 4)
			if lang == model.Hant {
				cps[0] = report.CopyWriter{Typ: report.ArchivePVRG, Reason: report.Reason7, Desc: report.Desc7, AllowAdd: true}
				cps[1] = report.CopyWriter{Typ: report.ArchiveCopyWriteG, Reason: report.Reason8, Desc: report.Desc8, AllowAdd: false}
				cps[2] = report.CopyWriter{Typ: report.ArchiveHarmG, Reason: report.Reason9, Desc: report.Desc9, AllowAdd: true}
				cps[3] = report.CopyWriter{Typ: report.OtherG, Reason: report.Reason10, Desc: report.Desc10, AllowAdd: true}
			} else {
				cps[0] = report.CopyWriter{Typ: report.ArchivePVR, Reason: report.Reason1, Desc: report.Desc1, AllowAdd: true}
				cps[1] = report.CopyWriter{Typ: report.ArchiveCopyWrite, Reason: report.Reason2, Desc: report.Desc2, AllowAdd: false}
				cps[2] = report.CopyWriter{Typ: report.ArchiveHarmG, Reason: report.Reason11, Desc: report.Desc11, AllowAdd: true}
				cps[3] = report.CopyWriter{Typ: report.Other, Reason: report.Reason6, Desc: report.Desc6, AllowAdd: true}
			}
		} else {
			cps = make([]report.CopyWriter, 6)
			cps[0] = report.CopyWriter{Typ: report.ArchivePVR, Reason: report.Reason1, Desc: report.Desc1, AllowAdd: true}
			cps[1] = report.CopyWriter{Typ: report.ArchiveCopyWrite, Reason: report.Reason2, Desc: report.Desc2, AllowAdd: false}
			cps[2] = report.CopyWriter{Typ: report.ArchiveCrash, Reason: report.Reason3, Desc: report.Desc3, AllowAdd: true}
			cps[3] = report.CopyWriter{Typ: report.ArchiveNotOwn, Reason: report.Reason4, Desc: report.Desc4, AllowAdd: true}
			cps[4] = report.CopyWriter{Typ: report.ArchiveBusiness, Reason: report.Reason5, Desc: report.Desc5, AllowAdd: true}
			cps[5] = report.CopyWriter{Typ: report.Other, Reason: report.Reason6, Desc: report.Desc6, AllowAdd: true}
		}
	} else {
		if model.IsOverseas(plat) && lang == model.Hant {
			cps = make([]report.CopyWriter, 4)
			cps[0] = report.CopyWriter{Typ: report.ArchivePVRG, Reason: report.Reason7, Desc: report.Desc7, AllowAdd: true}
			cps[1] = report.CopyWriter{Typ: report.ArchiveCopyWriteG, Reason: report.Reason8, Desc: report.Desc8, AllowAdd: false}
			cps[2] = report.CopyWriter{Typ: report.ArchiveHarmG, Reason: report.Reason9, Desc: report.Desc9, AllowAdd: true}
			cps[3] = report.CopyWriter{Typ: report.OtherG, Reason: report.Reason10, Desc: report.Desc10, AllowAdd: true}
		} else {
			cps = make([]report.CopyWriter, 4)
			cps[0] = report.CopyWriter{Typ: report.ArchivePVR, Reason: report.Reason1, Desc: report.Desc1, AllowAdd: true}
			cps[1] = report.CopyWriter{Typ: report.ArchiveCopyWrite, Reason: report.Reason2, Desc: report.Desc2, AllowAdd: false}
			cps[2] = report.CopyWriter{Typ: report.ArchiveCrash, Reason: report.Reason11, Desc: report.Desc11, AllowAdd: true}
			cps[3] = report.CopyWriter{Typ: report.Other, Reason: report.Reason6, Desc: report.Desc6, AllowAdd: true}
		}
	}
	return
}

// AddReport add a report .
func (s *Service) AddReport(c context.Context, mid, aid int64, mold int, ak, reason, pics string) (err error) {
	if len(reason) > _maxContentSize {
		err = ecode.FileTooLarge
		return
	}
	return s.reportDao.AddReport(c, mid, aid, mold, ak, reason, pics)
}

// Upload image upload .
func (s *Service) Upload(c context.Context, fileType string, body []byte) (url string, err error) {
	if len(body) == 0 {
		err = ecode.FileNotExists
		return
	}
	if len(body) > _maxFileSize {
		err = ecode.FileTooLarge
		return
	}
	return s.reportDao.Upload(c, fileType, bytes.NewReader(body))
}
