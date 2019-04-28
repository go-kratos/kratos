package service

import (
	"context"

	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
)

//RmToken Get PaaS token
func (s *Service) RmToken(c context.Context) (token string, err error) {
	return s.dao.RmToken(c)
}

//ClusterInfo get melloi server use
func (s *Service) ClusterInfo(c context.Context) (firstRetMap []*model.ClusterResponseItemsSon, err error) {
	var token string
	if token, err = s.RmToken(c); err != nil {
		//err = ecode.MelloiGetTreeTokenErr
		return
	}

	if firstRetMap, err = s.dao.NetInfo(c, token); err != nil {
		err = ecode.MerlinGetUserTreeFailed
		return
	}
	return
}
