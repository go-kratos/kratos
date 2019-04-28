package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/search/dao"
	"go-common/app/admin/main/search/model"
	"go-common/library/log"
)

func (s *Service) loadQueryConfproc() {
	for {
		if err := s.loadQueryConf(); err != nil {
			time.Sleep(time.Second)
			continue
		}
		time.Sleep(time.Minute)
	}
}

func (s *Service) loadQueryConf() (err error) {
	confs, err := s.dao.QueryConf(context.Background())
	if err != nil {
		return
	}
	if len(confs) > 0 {
		s.queryConf = confs
	}
	return
}

// CheckQueryConf check query conf
func (s *Service) CheckQueryConf(c context.Context, sp *model.QueryParams) (err error) {
	app, ok := s.queryConf[sp.Business]
	if app2, ok2 := model.QueryConf[sp.Business]; ok2 {
		app = app2
		ok = true
	}
	if !ok {
		err = fmt.Errorf("sp.Business(%s) not exist in queryConf", sp.Business)
		return
	}
	if app.ESCluster == "" {
		err = fmt.Errorf("app(%+v) escluster is empty", app)
		return
	}
	max := 1
	if app.MaxIndicesNum > 0 {
		max = app.MaxIndicesNum
	}
	indecies := strings.Split(sp.QueryBody.From, ",")
	if len(indecies) == 0 {
		err = fmt.Errorf("index name is required")
		return
	}
	if len(indecies) > max {
		err = fmt.Errorf("too many indecies(%v)", indecies)
		return
	}
	for _, index := range indecies {
		if !strings.Contains(index, app.IndexPrefix) {
			err = fmt.Errorf("invalid index name(%s)", index)
			return
		}
	}
	sp.AppIDConf = app
	return
}

// QueryBasic .
func (s *Service) QueryBasic(c context.Context, sp *model.QueryParams) (res *model.QueryResult, debug *model.QueryDebugResult, err error) {
	switch sp.Business {
	case "log_audit":
		t := strings.Split(sp.QueryBody.From, "_")
		if len(t) > 2 {
			logID, err := strconv.Atoi(t[2])
			if err != nil {
				log.Error("strconv.Atoi(%s) error(%v)", t[2], err)
			}
			logBusiness, ok := s.dao.GetLogInfo(sp.Business, logID)
			if ok {
				sp.AppIDConf.ESCluster = logBusiness.IndexCluster
			}
		}
	case "log_user_action":
		t := strings.Split(sp.QueryBody.From, "_")
		if len(t) > 3 {
			logID, err := strconv.Atoi(t[3])
			if err != nil {
				log.Error("strconv.Atoi(%s) error(%v)", t[3], err)
			}
			logBusiness, ok := s.dao.GetLogInfo(sp.Business, logID)
			if ok {
				sp.AppIDConf.ESCluster = logBusiness.IndexCluster
			}
		}
	}
	bQuery, qbDebug := s.dao.QueryBasic(c, sp)
	if res, debug, err = s.dao.QueryResult(c, bQuery, sp, qbDebug); err != nil {
		dao.PromError(fmt.Sprintf("es:%s 搜索失败", sp.Business), "s.dao.QueryBasic(%v) error(%v)", sp, err)
	}
	return
}

// QueryExtra .
func (s *Service) QueryExtra(c context.Context, sp *model.QueryParams) (res *model.QueryResult, debug *model.QueryDebugResult, err error) {
	switch sp.Business {
	case "archive_video_score":
		if res, debug, err = s.dao.ArchiveVideoScore(c, sp); err != nil {
			dao.PromError(fmt.Sprintf("es:%s 搜索失败", sp.Business), "s.dao.QueryExtra(%v) error(%v)", sp, err)
		}
	case "archive_score":
		if res, debug, err = s.dao.ArchiveScore(c, sp); err != nil {
			dao.PromError(fmt.Sprintf("es:%s 搜索失败", sp.Business), "s.dao.QueryExtra(%v) error(%v)", sp, err)
		}
	case "task_qa_random":
		if res, debug, err = s.dao.TaskQaRandom(c, sp); err != nil {
			dao.PromError(fmt.Sprintf("es:%s 搜索失败", sp.Business), "s.dao.QueryExtra.TaskQaRandom(%v) error(%v)", sp, err)
		}
	case "esports_contests_date":
		if res, debug, err = s.dao.EsportsContestsDate(c, sp); err != nil {
			dao.PromError(fmt.Sprintf("es:%s 搜索失败", sp.Business), "s.dao.QueryExtra.EsportsContestsDate(%v) error(%v)", sp, err)
		}
	case "creative_archive_search":
		if res, debug, err = s.dao.CreativeArchiveSearch(c, sp); err != nil {
			dao.PromError(fmt.Sprintf("es:%s 搜索失败", sp.Business), "s.dao.QueryExtra.CreativeArchiveSearch(%v) error(%v)", sp, err)
		}
	case "creative_archive_staff":
		if res, debug, err = s.dao.CreativeArchiveStaff(c, sp); err != nil {
			dao.PromError(fmt.Sprintf("es:%s 搜索失败", sp.Business), "s.dao.QueryExtra.CreativeArchiveStaff(%v) error(%v)", sp, err)
		}
	case "creative_archive_apply":
		if res, debug, err = s.dao.CreativeArchiveApply(c, sp); err != nil {
			dao.PromError(fmt.Sprintf("es:%s 搜索失败", sp.Business), "s.dao.QueryExtra.CreativeArchiveApply(%v) error(%v)", sp, err)
		}
	case "dm_history":
		if res, debug, err = s.dao.Scroll(c, sp); err != nil {
			dao.PromError(fmt.Sprintf("es:%s 搜索失败", sp.Business), "s.dao.QueryExtra.Scroll(%v) error(%v)", sp, err)
		}
	}
	return
}
