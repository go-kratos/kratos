package service

import (
	"context"

	"go-common/app/admin/ep/melloi/model"
	"go-common/library/log"
)

// QueryDependServiceAdmins query depend Service admin
func (s *Service) QueryDependServiceAdmins(c context.Context, serviceName string, sessionValue string) (map[string][]string, error) {

	var (
		biliPrex      = "bilibili."
		err           error
		roles         []*model.TreeRole
		dependService []string
		olesMap       = make(map[string][]string)
	)

	if dependService, err = s.dao.QueryServiceDepend(c, serviceName); err != nil {
		log.Error("query service depend error(%v)", err)
		return nil, err
	}
	dependService = append(dependService, serviceName)

	for _, service := range dependService {
		var userService []string
		if roles, err = s.QueryTreeAdmin(c, biliPrex+service, sessionValue); err != nil {
			log.Warn("query tree admin of service  (%s) depend error (%v)", biliPrex+service, err)
			continue
		}

		for _, role := range roles {
			// 增加多个service
			userService = append(olesMap[role.UserName], biliPrex+service)
			olesMap[role.UserName] = userService
		}

	}

	return olesMap, err
}
