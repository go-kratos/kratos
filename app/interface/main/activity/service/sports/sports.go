package sports

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"go-common/app/interface/main/activity/conf"
	"go-common/app/interface/main/activity/dao/sports"
	mdlsp "go-common/app/interface/main/activity/model/sports"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_qqNews       = 1
	_qqMatch      = 2
	_qqMatchTid   = "14"
	_qqTeamRank   = 3
	_qqRankTid    = "34"
	_qqPlayerRank = 4
	_qqRoute      = "matchUnion/fetchData"
	_newsRoute    = "getQQNewsIndexAndItemsVerify"
)

// Service struct
type Service struct {
	dao *sports.Dao
}

// New Service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		dao: sports.New(c),
	}
	go s.qqNews()
	go s.qqMatch()
	go s.qqTeamRank()
	go s.qqPlayerRank()
	return
}

func (s *Service) qqNews() {
	var (
		params = url.Values{}
		rs     *mdlsp.QqRes
		err    error
		c      = context.Background()
	)
	for {
		for t := 0; t < conf.Conf.Rule.QqTryCount; t++ {
			if rs, err = s.dao.QqNews(c, params, _newsRoute); err != nil || rs == nil {
				continue
			}
			s.dao.SetQqCache(c, &rs.IDlist, _qqNews)
			break
		}
		time.Sleep(time.Duration(conf.Conf.Rule.TickQq))
	}
}

func (s *Service) qqMatch() {
	var (
		rs     *json.RawMessage
		err    error
		c      = context.Background()
		params = url.Values{}
	)
	params.Set("tid", _qqMatchTid)
	params.Set("indexName", "col_4")
	params.Set("startTime", conf.Conf.Rule.QqStartTime)
	params.Set("endTime", conf.Conf.Rule.QqEndTime)
	for {
		for t := 0; t < conf.Conf.Rule.QqTryCount; t++ {
			if rs, err = s.dao.Qq(c, params, _qqRoute); err != nil || rs == nil || len(*rs) == 0 {
				continue
			}
			s.dao.SetQqCache(c, rs, _qqMatch)
			break
		}
		time.Sleep(time.Duration(conf.Conf.Rule.TickQq))
	}
}

func (s *Service) qqTeamRank() {
	var (
		rs     *json.RawMessage
		err    error
		c      = context.Background()
		params = url.Values{}
	)
	params.Set("tid", _qqRankTid)
	params.Set("competitionId", "4")
	params.Set("seasonId", conf.Conf.Rule.QqYear)
	params.Set("valueType", "teamRank")
	params.Set("valueId", "teamRank")
	for {
		for t := 0; t < conf.Conf.Rule.QqTryCount; t++ {
			if rs, err = s.dao.Qq(c, params, _qqRoute); err != nil || rs == nil || len(*rs) == 0 {
				continue
			}
			s.dao.SetQqCache(c, rs, _qqTeamRank)
			break
		}
		time.Sleep(time.Duration(conf.Conf.Rule.TickQq))
	}
}

func (s *Service) qqPlayerRank() {
	var (
		rs     *json.RawMessage
		err    error
		c      = context.Background()
		params = url.Values{}
	)
	params.Set("tid", _qqRankTid)
	params.Set("competitionId", "4")
	params.Set("seasonId", conf.Conf.Rule.PlayerYear)
	params.Set("valueType", "playerGoalRank")
	params.Set("valueId", "playerGoalRank")
	for {
		for t := 0; t < conf.Conf.Rule.QqTryCount; t++ {
			if rs, err = s.dao.Qq(c, params, _qqRoute); err != nil || rs == nil || len(*rs) == 0 {
				continue
			}
			s.dao.SetQqCache(c, rs, _qqPlayerRank)
			break
		}
		time.Sleep(time.Duration(conf.Conf.Rule.TickQq))
	}
}

// Qq get qq.
func (s *Service) Qq(c context.Context, params url.Values, p *mdlsp.ParamQq) (rs *json.RawMessage, err error) {
	if p.Tp > 0 {
		if rs, err = s.dao.QqCache(c, p.Tp); err != nil {
			log.Error("s.dao.QqCache  tp(%d) error(%v) ", p.Tp, err)
		}
	} else if rs, err = s.dao.Qq(c, params, p.Route); err != nil {
		sports.PromError("QQ接口错误", "s.dao.Qq route(%s) error(%v)", p.Route, err)
	}
	if rs == nil {
		err = ecode.ActivityServerTimeout
	}
	return
}

// News get qq news.
func (s *Service) News(c context.Context, params url.Values, p *mdlsp.ParamNews) (rs *mdlsp.QqRes, err error) {
	if rs, err = s.dao.QqNews(c, params, p.Route); err != nil {
		sports.PromError("QQ接口错误", "s.dao.Qq  error(%v)", err)
	}
	if rs == nil {
		err = ecode.ActivityServerTimeout
	}
	return
}
