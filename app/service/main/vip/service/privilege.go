package service

import (
	"context"

	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) loadPrivilege() (err error) {
	var (
		ps  []*model.Privilege
		prs []*model.PrivilegeResources
		c   = context.TODO()
	)
	if ps, err = s.dao.PrivilegeList(c); err != nil {
		return
	}
	if len(ps) == 0 {
		return
	}
	primap := make(map[int64][]*model.Privilege)
	for _, v := range ps {
		primap[v.LangType] = append(primap[v.LangType], v)
	}
	log.Info("loadPrivilege success %+v", primap)
	for k, v := range primap {
		s.vipPrivilege.Store(k, v)
	}

	if prs, err = s.dao.PrivilegeResourcesList(c); err != nil {
		return
	}
	prmap := map[int64]map[int8]*model.PrivilegeResources{}
	for _, v := range prs {
		if prmap[v.PID] == nil {
			prmap[v.PID] = map[int8]*model.PrivilegeResources{}
		}
		prmap[v.PID][v.Type] = v
	}
	s.vipPrivilegeResourcesMap = prmap
	return
}

// PrivilegesBySid by type query sid.
func (s *Service) PrivilegesBySid(c context.Context, arg *model.ArgPrivilegeBySid) (res *model.PrivilegesResp, err error) {
	var (
		p     *model.VipPriceConfig
		t     = model.AllPrivilege
		ps    = []*model.PrivilegeResp{}
		title = model.PrivilegeTitle
	)
	res = new(model.PrivilegesResp)
	if p = s.vipPriceMap[arg.Sid]; p == nil {
		err = ecode.VipSuitPirceNotFound
		return
	}
	if p.Month >= int16(_annualMonth) {
		t = model.OnlyAnnualPrivilege
		title = model.AnnualPrivilegeTitle
	}
	pris := s.getPrivileges(arg.Lang)
	for _, v := range pris {
		r := &model.PrivilegeResp{
			Name: v.Name,
			Type: v.Type,
		}
		if t == model.AllPrivilege && v.Type == model.OnlyAnnualPrivilege {
			r.IconURL = v.IconGrayURL
		} else {
			r.IconURL = v.IconURL
		}
		ps = append(ps, r)
	}
	res.List = ps
	res.Title = title
	return
}

// PrivilegesList by type query by type.
func (s *Service) PrivilegesList(c context.Context, t int8, lang string) (res *model.PrivilegesResp, err error) {
	pris := s.getPrivileges(lang)
	if len(pris) == 0 {
		return
	}
	res = new(model.PrivilegesResp)
	list := []*model.PrivilegeResp{}
	for _, v := range pris {
		r := &model.PrivilegeResp{
			Name: v.Name,
			Type: v.Type,
		}
		if t == model.AllPrivilege && v.Type == model.OnlyAnnualPrivilege {
			r.IconURL = v.IconGrayURL
		} else {
			r.IconURL = v.IconURL
		}
		list = append(list, r)
	}
	switch t {
	case model.OnlyAnnualPrivilege:
		res.Title = model.AnnualPrivilegeTitle
	default:
		res.Title = model.PrivilegeTitle
	}
	res.List = list
	return
}

// PrivilegesByType by type query privilege.
func (s *Service) PrivilegesByType(c context.Context, arg *model.ArgPrivilegeDetail) (res []*model.PrivilegeDetailResp, err error) {
	rt := model.ResourcesType(arg.Platform)
	pris := s.getPrivileges(arg.Lang)
	if len(pris) == 0 || len(s.vipPrivilegeResourcesMap) == 0 {
		log.Warn("privilege len(%d), map(%d)", len(pris), len(s.vipPrivilegeResourcesMap))
		return
	}
	for _, v := range pris {
		r := &model.PrivilegeDetailResp{
			Name:    v.Name,
			Title:   v.Title,
			Explain: v.Explain,
			Type:    v.Type,
			ID:      v.ID,
		}
		if arg.Type == model.AllPrivilege && v.Type == model.OnlyAnnualPrivilege {
			r.IconURL = v.IconGrayURL
		} else {
			r.IconURL = v.IconURL
		}
		if len(s.vipPrivilegeResourcesMap[v.ID]) != 0 && s.vipPrivilegeResourcesMap[v.ID][rt] != nil {
			r.Link = s.vipPrivilegeResourcesMap[v.ID][rt].Link
			r.ImageURL = s.vipPrivilegeResourcesMap[v.ID][rt].ImageURL
		}
		res = append(res, r)
	}
	return
}

func (s *Service) getPrivileges(lang string) (ps []*model.Privilege) {
	var lt int64
	switch lang {
	case "zh_TW", "zh_HK":
		lt = 1
	default:
		log.Warn("unknow lang(%s)", lang)
	}
	prisI, ok := s.vipPrivilege.Load(lt)
	if !ok {
		return
	}
	ps = prisI.(interface{}).([]*model.Privilege)
	return
}
