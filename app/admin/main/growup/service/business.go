package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/admin/main/growup/dao/resource"
	"go-common/app/admin/main/growup/model"

	"go-common/library/log"
)

var (
	// 视频和专栏的权限点
	_videoPri  = 13
	_columnPri = 103
)

// BusPrivilege business privilege
func (s *Service) BusPrivilege(c context.Context, username string, ctypes string) (res []*model.BusRes, err error) {
	res = make([]*model.BusRes, 0)
	typs := strings.Split(ctypes, ",")
	if len(typs) == 0 {
		return
	}
	for _, typ := range typs {
		var ctype int
		ctype, err = strconv.Atoi(typ)
		if err != nil {
			log.Error("strconv.Atoi error(%v)", err)
			return
		}
		var r []*model.BusRes
		r, err = s.busPrivilege(c, username, ctype)
		if err != nil {
			log.Error("s.busPrivilege error(%v)", err)
			return
		}
		res = append(res, r...)
	}
	return
}

func (s *Service) busPrivilege(c context.Context, username string, ctype int) (res []*model.BusRes, err error) {
	category, err := s.getBusCategory(c, ctype)
	if err != nil {
		log.Error("s.getBusCategory error(%v)", err)
		return
	}
	userPri, err := s.GetUserPri(username)
	if err != nil {
		log.Error("s.GetUserPri error(%v)", err)
		return
	}

	// 获取数据源权限
	fatherID := 0
	switch ctype {
	case 1:
		fatherID = _videoPri
	case 2:
		fatherID = _columnPri
	}

	allPrivilege, err := s.dao.GetLevelPrivileges(fmt.Sprintf("level = 3 AND father_id = %d", fatherID))
	if err != nil {
		log.Error("s.dao.GetLevelPrivileges Error(%v)", err)
		return
	}
	res = make([]*model.BusRes, 0)
	for _, p := range allPrivilege {
		if !userPri[p.ID] {
			continue
		}
		if cid, ok := category[p.Title]; ok {
			res = append(res, &model.BusRes{
				PrivilegeID: p.ID,
				CategoryID:  cid,
				Name:        p.Title})
		}
	}
	return
}

func (s *Service) getBusCategory(c context.Context, ctype int) (categorys map[string]int64, err error) {
	categorys = make(map[string]int64)
	switch ctype {
	case 1:
		return resource.VideoCategoryNameToID(c)
	case 2:
		return resource.ColumnCategoryNameToID(c)
	}
	return
}
