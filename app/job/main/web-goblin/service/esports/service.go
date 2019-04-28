package esports

import (
	"context"
	"fmt"
	"math"
	"time"

	"go-common/app/job/main/web-goblin/conf"
	"go-common/app/job/main/web-goblin/dao/esports"
	mdlesp "go-common/app/job/main/web-goblin/model/esports"
	arcclient "go-common/app/service/main/archive/api"
	arcmdl "go-common/app/service/main/archive/api"
	favrpc "go-common/app/service/main/favorite/api/gorpc"
	"go-common/app/service/main/favorite/model"
	mdlfav "go-common/app/service/main/favorite/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_favUsers    = 1000
	_favTryTimes = 3
	_defContest  = 0
	_linkinfo    = "点击前往直播间"
	_pushinfo    = "进入直播>>"
	_msgSize     = 1000
	_tpMessage   = 0
	_tpPush      = 1
	_arcMaxLimit = 50
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *esports.Dao
	// rpc
	fav       *favrpc.Service
	arcClient arcclient.ArchiveClient
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: esports.New(c),
		fav: favrpc.New2(c.FavoriteRPC),
	}
	var err error
	if s.arcClient, err = arcclient.NewClient(c.ArcClient); err != nil {
		panic(err)
	}
	go s.contests()
	go s.contestsPush()
	go s.arcScore()
	return s
}

func (s *Service) contests() {
	var (
		err      error
		contests []*mdlesp.Contest
	)
	for {
		stime := time.Now().Add(time.Duration(s.c.Rule.Before))
		etime := stime.Add(time.Duration(s.c.Rule.SleepInterval))
		if contests, err = s.dao.Contests(context.Background(), stime.Unix(), etime.Unix()); err != nil {
			time.Sleep(time.Second)
			log.Error("contests s.dao.Contests stime(%d) error(%v)", stime.Unix(), err)
			continue
		}
		for _, contest := range contests {
			tmpContest := contest
			go s.sendContests(tmpContest)
		}
		time.Sleep(time.Duration(s.c.Rule.SleepInterval))
	}
}

func (s *Service) contestsPush() {
	var (
		err      error
		contests []*mdlesp.Contest
	)
	for {
		stime := time.Now()
		etime := stime.Add(time.Second)
		if contests, err = s.dao.Contests(context.Background(), stime.Unix(), etime.Unix()); err != nil {
			time.Sleep(time.Millisecond * 100)
			log.Error("contestsPush s.dao.Contests stime(%d) error(%v)", stime.Unix(), err)
			continue
		}
		for _, contest := range contests {
			tmpContest := contest
			go s.pubContests(tmpContest)
		}
		time.Sleep(time.Second)
	}
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) sendContests(contest *mdlesp.Contest) (err error) {
	var (
		mids []int64
		msg  string
	)
	link := fmt.Sprintf("#{%s}{\"https://live.bilibili.com/%d\"}", _linkinfo, contest.LiveRoom)
	if mids, msg, err = s.midsParams(contest, link, _tpMessage); err != nil {
		log.Error("sendContests s.midsParams(%+v) mids_total(%d) error(%v)", contest, len(mids), err)
		return
	}
	s.dao.Batch(mids, msg, contest, _msgSize, s.dao.SendMessage)
	return
}

func (s *Service) pubContests(contest *mdlesp.Contest) (err error) {
	var (
		mids []int64
		msg  string
	)
	if mids, msg, err = s.midsParams(contest, _pushinfo, _tpPush); err != nil {
		log.Error("pubContests s.midsParams(%+v) mids_total(%d) error(%v)", contest, len(mids), err)
		return
	}
	s.dao.Batch(mids, msg, contest, s.c.Push.PartSize, s.dao.NoticeUser)
	return
}

func (s *Service) midsParams(contest *mdlesp.Contest, link string, tp int) (mids []int64, msg string, err error) {
	var (
		userList       *mdlfav.UserList
		teams          []*mdlesp.Team
		homeID, awayID int64
		tMap           map[int64]string
		pageCount      int
	)
	if userList, err = s.favUsers(contest.ID, 1); err != nil || userList == nil || len(userList.List) == 0 {
		log.Error("s.favUsers  contestID(%v) first error(%+v)", contest.ID, err)
		return
	}
	ms := make(map[int64]struct{}, userList.Page.Total)
	for _, user := range userList.List {
		if _, ok := ms[user.Mid]; ok {
			continue
		}
		ms[user.Mid] = struct{}{}
		mids = append(mids, user.Mid)
	}
	if userList.Page.Size == 0 {
		pageCount = 0
	} else {
		pageCount = int(math.Ceil(float64(userList.Page.Total) / float64(userList.Page.Size)))
	}
	for i := 2; i <= pageCount; i++ {
		if userList, err = s.favUsers(contest.ID, i); err != nil || userList == nil || len(userList.List) == 0 {
			log.Error("s.favUsers  contestID(%v) pn(%d) error(%+v)", contest.ID, i, err)
			err = nil
			continue
		}
		for _, user := range userList.List {
			if _, ok := ms[user.Mid]; ok {
				continue
			}
			ms[user.Mid] = struct{}{}
			mids = append(mids, user.Mid)
		}
	}
	if len(mids) == 0 {
		err = ecode.RequestErr
		return
	}
	tm := time.Unix(contest.Stime, 0)
	stime := tm.Format("2006-01-02 15:04:05")
	if contest.Special == _defContest {
		homeID = contest.HomeID
		awayID = contest.AwayID
		if teams, err = s.dao.Teams(context.Background(), homeID, awayID); err != nil || len(teams) == 0 {
			log.Error("midsParams  s.dao.Teams homeID(%d) awayID(%d) error(%v)", homeID, awayID, err)
			return
		}
		tMap = make(map[int64]string, 2)
		for _, temp := range teams {
			tMap[temp.ID] = temp.Title
		}
		if tp == _tpMessage {
			msg = fmt.Sprintf(s.c.Rule.AlertBodyDefault, contest.SeasonTitle, stime, tMap[contest.HomeID], tMap[contest.AwayID], link)
		} else if tp == _tpPush {
			msg = fmt.Sprintf(s.c.Push.BodyDefault, contest.SeasonTitle, tMap[contest.HomeID], tMap[contest.AwayID], link)
		}
	} else {
		if tp == _tpMessage {
			msg = fmt.Sprintf(s.c.Rule.AlertBodySpecial, contest.SeasonTitle, stime, contest.SpecialName, link)
		} else if tp == _tpPush {
			msg = fmt.Sprintf(s.c.Push.BodySpecial, contest.SeasonTitle, contest.SpecialName, link)
		}
	}
	count := len(mids)
	log.Info("midsParams get contest cid(%d) users number(%d)", contest.ID, count)
	return
}

func (s *Service) favUsers(cid int64, pn int) (res *mdlfav.UserList, err error) {
	for i := 0; i < _favTryTimes; i++ {
		if res, err = s.fav.Users(context.Background(), &model.ArgUsers{Type: model.TypeEsports, Oid: cid, Pn: pn, Ps: _favUsers}); err == nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	if err != nil {
		err = errors.Wrapf(err, "favUsers s.fav.Users cid(%d) pn(%d)", cid, pn)
	}
	return
}

func (s *Service) arcScore() {
	var (
		id int64
		c  = context.Background()
	)
	for {
		av, err := s.dao.Arcs(c, id, _arcMaxLimit)
		if err != nil {
			log.Error("ArcScore  s.dao.Arcs ID(%d) Limit(%d) error(%v)", id, _arcMaxLimit, err)
			id = id + int64(_arcMaxLimit)
			time.Sleep(time.Second)
			continue
		}
		if len(av) == 0 {
			id = 0
			time.Sleep(time.Duration(s.c.Rule.ScoreSleep))
			continue
		}
		go s.upArcScore(c, av)
		id = av[len(av)-1].ID
		time.Sleep(time.Second)
	}
}

func (s *Service) upArcScore(c context.Context, partArcs []*mdlesp.Arc) (err error) {
	var (
		partAids  []int64
		arcsReply *arcmdl.ArcsReply
	)
	for _, arc := range partArcs {
		partAids = append(partAids, arc.Aid)
	}
	if len(partAids) == 0 {
		return
	}
	if arcsReply, err = s.arcClient.Arcs(c, &arcmdl.ArcsRequest{Aids: partAids}); err != nil || arcsReply == nil {
		log.Error("upArcScore  s.arcClient.Arcs(%v) error(%v)", partAids, err)
		return
	}
	if len(arcsReply.Arcs) > 0 {
		if err = s.dao.UpArcScore(c, partArcs, arcsReply.Arcs); err != nil {
			log.Error("upArcScore  s.dao.UpArcScore arcs(%+v) error(%v)", arcsReply, err)
		}
	}
	return
}
