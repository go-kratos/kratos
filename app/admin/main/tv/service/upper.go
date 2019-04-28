package service

import (
	"context"

	"go-common/app/admin/main/tv/model"
	accmdl "go-common/app/service/main/account/api"
	account "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

const (
	_removeArcs    = 2
	_deleted       = 1
	_toinit        = 1
	_requestedInit = 2
	_onlineAct     = "1"
	_offlineAct    = "0"
	_validOnline   = 1
	_validOffline  = 0
)

// AddMids CheckMids checks the mids, whether all the uppers exist
func (s *Service) AddMids(mids []int64) (res *model.AddResp, err error) {
	var (
		accsReply *accmdl.InfosReply
		accsInfo  map[int64]*account.Info
		midExist  bool
	)
	// init the response
	res = &model.AddResp{
		Succ:     []int64{},
		Exist:    []int64{},
		Invalids: []int64{},
	}
	// rpc get all the mids' info
	if accsReply, err = s.accClient.Infos3(ctx, &accmdl.MidsReq{
		Mids: mids,
	}); err != nil {
		log.Error("CheckMids Mids: %v, Error: %v", mids, err)
		return
	}
	accsInfo = accsReply.Infos
	for _, v := range mids {
		// if invalid account
		if _, ok := accsInfo[v]; !ok {
			res.Invalids = append(res.Invalids, v)
			continue
		}
		if midExist, err = s.existMid(v); err != nil {
			log.Error("AddMid %d, Error %v", v, err)
			return
		}
		// if the account is existing in our DB
		if midExist {
			res.Exist = append(res.Exist, v)
			continue
		}
		// add the upper
		if err = s.dao.UpAdd(v); err != nil {
			log.Error("AddMid %d, Error %v", v, err)
			return
		}
		res.Succ = append(res.Succ, v)
	}
	return
}

// existMid checks whether the mid is existing in our DB
func (s *Service) existMid(mid int64) (exist bool, err error) {
	if err = s.DB.Where("deleted = 0").Where("mid = ?", mid).First(&model.Upper{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		log.Error("existMid %d, Error %v", mid, err)
		return
	}
	return true, nil
}

// ImportMids is for updating the data to tell the tv-job to import the uppers' all videos
func (s *Service) ImportMids(mids []int64) (res *model.ImportResp, err error) {
	var midExist bool
	res = &model.ImportResp{
		NotExist: []int64{},
		Succ:     []int64{},
	}
	for _, v := range mids {
		if midExist, err = s.existMid(v); err != nil {
			log.Error("ImportMids %d, Error %v", v, err)
			return
		}
		if !midExist {
			res.NotExist = append(res.NotExist, v)
			continue
		}
		if err = s.DB.Model(&model.Upper{}).Where("mid = ?", v).Update(map[string]int{"toinit": _toinit, "state": _requestedInit}).Error; err != nil {
			log.Error("ImportMids %d, Error %v", v, err)
			return
		}
		res.Succ = append(res.Succ, v)
	}
	return
}

// DelMid is for updating remove one upper from the list
func (s *Service) DelMid(mid int64) (err error) {
	var exist bool
	if exist, err = s.existMid(mid); err != nil {
		return
	}
	if !exist {
		return ecode.TvUpperNotInList
	}
	if err = s.DB.Model(&model.Upper{}).Where("mid = ?", mid).Update(map[string]int{"deleted": _deleted, "toinit": _removeArcs}).Error; err != nil {
		log.Error("DelMid %d, Error %v", mid, err)
	}
	return
}

// UpList shows the upper list
func (s *Service) UpList(order int, page int, name string, id int) (pager *model.UpperPager, err error) {
	var (
		source                 []*model.Upper
		mids                   []int64
		match                  map[int64]*account.Info
		info                   *account.Info
		ok                     bool
		namesReply, infosReply *accmdl.InfosReply
		nameRes                map[int64]*account.Info
		ids                    []int64
	)
	pager = &model.UpperPager{}
	// id treatmnet
	if id != 0 {
		ids = append(ids, int64(id))
	}
	// name treatment
	if name != "" {
		if namesReply, err = s.accClient.InfosByName3(ctx, &accmdl.NamesReq{
			Names: []string{name},
		}); err != nil {
			log.Error("accRPC InfosByName3 %s, Err %v", name, err)
			return
		}
		nameRes = namesReply.Infos
		for k := range nameRes {
			ids = append(ids, k)
		}
	}
	if source, pager.Page, err = s.dao.UpList(order, page, ids); err != nil {
		return
	}
	// pick upper's name from AccRPC
	for _, v := range source {
		mids = append(mids, v.MID)
	}
	if infosReply, err = s.accClient.Infos3(ctx, &accmdl.MidsReq{
		Mids: mids,
	}); err != nil {
		log.Error("accRPC Infos3, %v, Error %v", mids, err)
		return
	}
	match = infosReply.Infos
	for _, v := range source { // arrange the data to output
		if info, ok = match[v.MID]; !ok {
			log.Error("Mid %v, AccRPC info is nil", v.MID)
			continue
		}
		pager.Items = append(pager.Items, &model.UpperR{
			MID:   v.MID,
			State: v.State,
			Name:  info.Name,
			Ctime: s.TimeFormat(v.Ctime),
			Mtime: s.TimeFormat(v.Mtime),
		})
	}
	return
}

// CmsList is the upper cms list service function
func (s *Service) CmsList(ctx context.Context, req *model.ReqUpCms) (pager *model.CmsUpperPager, err error) {
	ups, page, errList := s.dao.UpCmsList(req)
	if errList != nil {
		return nil, errList
	}
	for _, v := range ups {
		v.MtimeStr = v.Mtime.Time().Format("2006-01-02 15:04:05")
	}
	pager = &model.CmsUpperPager{
		Items: ups,
		Page:  page,
	}
	return
}

// CmsAudit updates the mids' valid status
func (s *Service) CmsAudit(ctx context.Context, mids []int64, action string) (resp *model.RespUpAudit, err error) {
	var (
		okMids   map[int64]*model.UpMC
		validAct int
	)
	resp = &model.RespUpAudit{
		Succ: mids,
	}
	if okMids, err = s.dao.VerifyIds(mids); err != nil {
		return
	}
	if len(okMids) != len(mids) {
		succ := []int64{}
		for _, v := range mids {
			if _, ok := okMids[v]; !ok {
				resp.Invalid = append(resp.Invalid, v)
			} else {
				succ = append(succ, v)
			}
		}
		resp.Succ = succ
	}
	if action == _onlineAct {
		validAct = _validOnline
	} else if action == _offlineAct {
		validAct = _validOffline
	}
	if err = s.dao.AuditIds(resp.Succ, validAct); err != nil {
		return
	}
	for _, v := range resp.Succ {
		s.dao.DelCache(ctx, v) // delete cache from MC because their status has been updated
	}
	return
}

// CmsEdit updates the upper's info in both DB and cache
func (s *Service) CmsEdit(ctx context.Context, req *model.ReqUpEdit) (err error) {
	var (
		okMids map[int64]*model.UpMC
		upInfo *model.UpMC
		ok     bool
	)
	if okMids, err = s.dao.VerifyIds([]int64{req.MID}); err != nil {
		return
	}
	if len(okMids) == 0 {
		return ecode.TvUpperNotInList
	}
	if upInfo, ok = okMids[req.MID]; !ok {
		log.Error("okMids %v, Mid %d, not found", okMids, req.MID)
		return ecode.NothingFound
	}
	if err = s.dao.SetUpper(req); err != nil { // update upper cms info in DB
		return
	}
	upInfo.CMSFace = req.Face
	upInfo.CMSName = req.Name
	err = s.dao.SetUpMetaCache(ctx, upInfo) // update Cache
	return
}
