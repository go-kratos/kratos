package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/ep/saga/conf"
	"go-common/app/admin/ep/saga/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/xanzy/go-gitlab"
)

// QueryProject query commit info according to project id.
func (s *Service) QueryProject(c context.Context, object string, req *model.ProjectDataReq) (resp *model.ProjectDataResp, err error) {
	var (
		data     []*model.DataWithTime
		queryDes string
		total    int
	)

	log.Info("QuerySingleProjectData Type: %d", req.QueryType)
	switch req.QueryType {
	case model.LastYearPerMonth:
		queryDes = model.LastYearPerMonthNote
	case model.LastMonthPerDay:
		queryDes = model.LastMonthPerDayNote
	case model.LastYearPerDay:
		queryDes = model.LastYearPerDayNote
	default:
		log.Warn("QueryProjectCommit Type is not in range")
		return
	}
	queryDes = req.ProjectName + " " + object + queryDes

	if data, total, err = s.QueryProjectByTime(req.ProjectID, object, req.QueryType); err != nil {
		return
	}

	resp = &model.ProjectDataResp{
		ProjectName: req.ProjectName,
		QueryDes:    queryDes,
		Total:       total,
		Data:        data,
	}
	return
}

// QueryProjectByTime ...
func (s *Service) QueryProjectByTime(projectID int, object string, queryType int) (resp []*model.DataWithTime, allNum int, err error) {
	var (
		layout    = "2006-01-02"
		fmtLayout = `%d-%d-%d 00:00:00`
		//response *gitlab.Response
		since      time.Time
		until      time.Time
		count      int
		totalItems int
	)

	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)

	if queryType == model.LastYearPerMonth {
		count = model.MonthNumPerYear
	} else if queryType == model.LastMonthPerDay {
		_, _, count = thisMonth.AddDate(0, 0, -1).Date()
	} else if queryType == model.LastYearPerDay {
		count = model.DayNumPerYear
	}

	for i := 1; i <= count; i++ {
		if queryType == model.LastYearPerMonth {
			since = thisMonth.AddDate(0, -i, 0)
			until = thisMonth.AddDate(0, -i+1, 0)
		} else if queryType == model.LastMonthPerDay {
			since = thisMonth.AddDate(0, 0, -i)
			until = thisMonth.AddDate(0, 0, -i+1)
		} else if queryType == model.LastYearPerDay {
			since = thisMonth.AddDate(0, 0, -i)
			until = thisMonth.AddDate(0, 0, -i+1)
		}

		sinceStr := fmt.Sprintf(fmtLayout, since.Year(), since.Month(), since.Day())
		untilStr := fmt.Sprintf(fmtLayout, until.Year(), until.Month(), until.Day())
		if object == model.ObjectCommit {
			/*if _, response, err = s.gitlab.ListProjectCommit(projectID, 1, &since, &until); err != nil {
				return
			}*/
			if totalItems, err = s.dao.CountCommitByTime(projectID, sinceStr, untilStr); err != nil {
				return
			}
		} else if object == model.ObjectMR {
			/*if _, response, err = s.gitlab.ListProjectMergeRequests(projectID, &since, &until, -1); err != nil {
				return
			}*/
			if totalItems, err = s.dao.CountMRByTime(projectID, sinceStr, untilStr); err != nil {
				return
			}
		} else {
			log.Warn("QueryProjectByTime object(%s) is not support!", object)
			return
		}

		perData := &model.DataWithTime{
			//TotalItem: response.TotalItems,
			TotalItem: totalItems,
			StartTime: since.Format(layout),
			EndTime:   until.Format(layout),
		}
		resp = append(resp, perData)
		//allNum = allNum + response.TotalItems
		allNum = allNum + totalItems
	}
	return
}

// QueryTeam ...
func (s *Service) QueryTeam(c context.Context, object string, req *model.TeamDataRequest) (resp *model.TeamDataResp, err error) {
	var (
		projectInfo []*model.ProjectInfo
		reqProject  = &model.ProjectInfoRequest{}

		dataMap     = make(map[string]*model.TeamDataResp)
		data        []*model.DataWithTime
		queryDes    string
		total       int
		key         string
		keyNotExist bool
	)

	if len(req.Department) <= 0 && len(req.Business) <= 0 {
		log.Warn("query department and business are empty!")
		return
	}

	//log.Info("QueryTeamCommit Query Type: %d", req.QueryType)
	switch req.QueryType {
	case model.LastYearPerMonth:
		queryDes = model.LastYearPerMonthNote
	case model.LastMonthPerDay:
		queryDes = model.LastMonthPerDayNote
	case model.LastYearPerDay:
		queryDes = model.LastYearPerDayNote
	default:
		log.Warn("QueryTeamCommit Type is not in range")
		return
	}
	queryDes = req.Department + " " + req.Business + " " + object + queryDes

	//get value from mc
	key = "saga_admin_" + req.Department + "_" + req.Business + "_" + model.KeyTypeConst[req.QueryType]
	if err = s.dao.GetData(c, key, &dataMap); err != nil {
		if err == memcache.ErrNotFound {
			log.Warn("no such key (%s) in cache, err (%s)", key, err.Error())
			keyNotExist = true
		} else {
			return
		}
	}
	if _, ok := dataMap[object]; !ok {
		keyNotExist = true
	} else {
		resp = dataMap[object]
		return
	}

	log.Info("sync team %s start => type= %d, Department= %s, Business= %s", object, req.QueryType, req.Department, req.Business)

	reqProject.Department = req.Department
	reqProject.Business = req.Business
	if _, projectInfo, err = s.dao.QueryProjectInfo(false, reqProject); err != nil {
		return
	}

	if len(projectInfo) <= 0 {
		log.Warn("Found no project!")
		return
	}

	if data, total, err = s.QueryTeamByTime(object, req, req.QueryType, projectInfo); err != nil {
		return
	}

	resp = &model.TeamDataResp{
		Department: req.Department,
		Business:   req.Business,
		QueryDes:   queryDes,
		Total:      total,
		Data:       data,
	}

	//set value to mc
	if keyNotExist {
		dataMap[object] = resp
		if err = s.dao.SetData(c, key, dataMap); err != nil {
			return
		}
	}

	log.Info("sync team %s end", object)
	return
}

// QueryTeamByTime query commit info per month of last year and per day of last month according to team.
func (s *Service) QueryTeamByTime(object string, req *model.TeamDataRequest, queryType int, projectInfo []*model.ProjectInfo) (resp []*model.DataWithTime, allNum int, err error) {
	var (
		layout   = "2006-01-02"
		response *gitlab.Response
		since    time.Time
		until    time.Time
		count    int
		num      int
	)

	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)

	if queryType == model.LastYearPerMonth {
		count = model.MonthNumPerYear
	} else if queryType == model.LastMonthPerDay {
		_, _, count = thisMonth.AddDate(0, 0, -1).Date()
	} else if queryType == model.LastYearPerDay {
		count = model.DayNumPerYear
	}

	for i := 1; i <= count; i++ {
		if queryType == model.LastYearPerMonth {
			since = thisMonth.AddDate(0, -i, 0)
			until = thisMonth.AddDate(0, -i+1, 0)
		} else if queryType == model.LastMonthPerDay {
			since = thisMonth.AddDate(0, 0, -i)
			until = thisMonth.AddDate(0, 0, -i+1)
		} else if queryType == model.LastYearPerDay {
			since = thisMonth.AddDate(0, 0, -i)
			until = thisMonth.AddDate(0, 0, -i+1)
		}

		num = 0
		for _, project := range projectInfo {
			if object == model.ObjectCommit {
				if _, response, err = s.gitlab.ListProjectCommit(project.ProjectID, 1, &since, &until); err != nil {
					return
				}
			} else if object == model.ObjectMR {
				if _, response, err = s.gitlab.ListProjectMergeRequests(project.ProjectID, &since, &until, -1); err != nil {
					return
				}
			} else {
				log.Warn("QueryTeamByTime object(%s) is not support!", object)
				return
			}
			num = num + response.TotalItems
		}

		perData := &model.DataWithTime{
			TotalItem: num,
			StartTime: since.Format(layout),
			EndTime:   until.Format(layout),
		}
		resp = append(resp, perData)
		allNum = allNum + num
	}
	return
}

// SyncData ...
func (s *Service) SyncData(c context.Context) (err error) {
	log.Info("sync all data info start!")
	for _, de := range conf.Conf.Property.DeInfo {
		for _, bu := range conf.Conf.Property.BuInfo {
			for k, keyType := range model.KeyTypeConst {

				key := "saga_admin_" + de.Value + "_" + bu.Value + "_" + keyType
				if err = s.dao.DeleteData(c, key); err != nil {
					return
				}

				req := &model.TeamDataRequest{
					TeamParam: model.TeamParam{
						Department: de.Value,
						Business:   bu.Value,
					},
					QueryType: k,
				}

				if k == 3 {
					if _, err = s.QueryTeamPipeline(c, req); err != nil {
						return
					}
					continue
				}

				if _, err = s.QueryTeamCommit(c, req); err != nil {
					return
				}

				if _, err = s.QueryTeamMr(c, req); err != nil {
					return
				}
			}
		}
	}
	log.Info("sync all data info end!")
	return
}
