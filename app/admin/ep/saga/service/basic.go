package service

import (
	"context"

	"go-common/app/admin/ep/saga/conf"
	"go-common/app/admin/ep/saga/model"
	"go-common/library/log"
)

// QueryProjectStatus ...
func (s *Service) QueryProjectStatus(c context.Context, req *model.ProjectDataReq) (resp []string) {
	return conf.Conf.Property.DefaultProject.Status
}

// QueryProjectTypes ...
func (s *Service) QueryProjectTypes(c context.Context, req *model.ProjectDataReq) (resp []*model.QueryTypeItem) {
	queryTypes := conf.Conf.Property.DefaultProject.Types

	for _, queryType := range queryTypes {
		item := &model.QueryTypeItem{}
		switch queryType {
		case model.LastYearPerMonth:
			item.Name = queryType
			item.Value = model.LastYearPerMonthNote
		case model.LastMonthPerDay:
			item.Name = queryType
			item.Value = model.LastMonthPerDayNote
		case model.LastYearPerDay:
			item.Name = queryType
			item.Value = model.LastYearPerDayNote
		default:
			log.Warn("QueryProjectCommit Type is not in range")
			return
		}
		resp = append(resp, item)
	}

	return
}
