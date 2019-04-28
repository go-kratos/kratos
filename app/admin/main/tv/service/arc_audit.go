package service

import (
	"go-common/app/admin/main/tv/model"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

func arcNormal(state int32) bool {
	if state >= 0 || state == -6 { // archive can play
		return true
	}
	return false
}

//AddArcs is used for adding archive
func (s *Service) AddArcs(aids []int64) (res *model.AddResp, err error) {
	var (
		valid  bool
		arc    *model.SimpleArc
		errFmt = "AddArcs %d, Error %v"
	)
	res = &model.AddResp{
		Succ:     []int64{},
		Invalids: []int64{},
		Exist:    []int64{},
	}
	for _, v := range aids {
		if valid, err = s.CheckArc(v); err != nil {
			log.Error(errFmt, v, err)
			return
		}
		// not valid aids
		if !valid {
			res.Invalids = append(res.Invalids, v)
			continue
		}
		if arc, err = s.ExistArc(v); err != nil {
			log.Error(errFmt, v, err)
			return
		}
		// in our DB, already exist aids
		if arc != nil {
			res.Exist = append(res.Exist, v)
			continue
		}
		if err = s.dao.NeedImport(v); err != nil {
			log.Error(errFmt, v, err)
			return
		}
		// added succesfully aids
		res.Succ = append(res.Succ, v)
	}
	return
}

// CheckArc checks whether the archive is able to play and existing in Archive DB
func (s *Service) CheckArc(aid int64) (ok bool, err error) {
	var (
		argAid2  = &arcmdl.ArcRequest{Aid: aid}
		arcReply *arcmdl.ArcReply
	)
	if arcReply, err = s.arcClient.Arc(ctx, argAid2); err != nil {
		if ecode.NothingFound.Equal(err) { // archive not found at all
			err = nil
			return
		}
		log.Error("s.arcRPC.Archive3(%v) error(%v)", argAid2, err)
		return
	}
	arc := arcReply.Arc
	if s.Contains(arc.TypeID) { // filter pgc types
		ok = false
		return
	}
	if arc.Copyright != 1 {
		ok = false
		return
	}
	if arc.Rights.UGCPay == archive.AttrYes {
		ok = false
		return
	}
	if arcNormal(arc.State) {
		ok = true
	}
	return
}

// ExistArc checks whether the archive is already in our TV DB, which means no need to import again
func (s *Service) ExistArc(aid int64) (res *model.SimpleArc, err error) {
	var arc = model.SimpleArc{}
	if err = s.DB.Where("aid = ?", aid).Where("deleted = ?", 0).First(&arc).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
			res = nil
			return
		}
		log.Error("ExistArc DB Error %v", err)
		return
	}
	return &arc, nil
}
