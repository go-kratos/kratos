package service

import (
	"context"

	"go-common/app/interface/main/web/model"
	resmdl "go-common/app/service/main/resource/model"
	"go-common/library/log"
)

const (
	_kvID = 2326
)

var _emptyAssi = make([]*model.Kv, 0)

// Kv get baidu kv
func (s *Service) Kv(c context.Context) (res []*model.Kv, err error) {
	var tmp *resmdl.Resource
	if tmp, err = s.res.Resource(c, &resmdl.ArgRes{ResID: _kvID}); err != nil {
		log.Error("s.res.Resource(%d) error(%v)", _kvID, err)
		return
	}
	if len(tmp.Assignments) == 0 {
		res = _emptyAssi
		return
	}
	for _, assi := range tmp.Assignments {
		res = append(res, &model.Kv{ID: assi.ID, Name: assi.Name, Pic: assi.Pic, URL: assi.URL, ResID: assi.ResID, STime: assi.STime, ETime: assi.STime})
	}
	return
}

// CmtBox get live dm box
func (s *Service) CmtBox(c context.Context, id int64) (res *resmdl.Cmtbox, err error) {
	if res, err = s.res.Cmtbox(c, &resmdl.ArgCmtbox{ID: id}); err != nil {
		log.Error("s.res.Cmtbox(%d) error(%v)", id, err)
	}
	return
}

// AbServer get ab server info.
func (s *Service) AbServer(c context.Context, mid int64, platform int, channel, buvid string) (data model.AbServer, err error) {
	return s.dao.AbServer(c, mid, platform, channel, buvid)
}
