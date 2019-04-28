package bnj

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"

	"go-common/app/interface/main/activity/conf"
	"go-common/app/interface/main/activity/dao/bnj"
	"go-common/app/interface/main/activity/dao/like"
	bnjmdl "go-common/app/interface/main/activity/model/bnj"
	arcclient "go-common/app/service/main/archive/api"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// Service .
type Service struct {
	c           *conf.Config
	arcClient   arcclient.ArchiveClient
	dao         *bnj.Dao
	likeDao     *like.Dao
	resetPub    *databus.Databus
	previewArcs map[int64]*arcclient.Arc
	bnjAdmins   map[int64]struct{}
	likeCount   int64
	timeReset   int64
	resetMid    int64
	timeFinish  int64
	resetCD     int32
}

// New init bnj service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:        c,
		dao:      bnj.New(c),
		likeDao:  like.New(c),
		resetPub: databus.New(c.Databus.Bnj),
	}
	var err error
	if s.arcClient, err = arcclient.NewClient(c.ArcClient); err != nil {
		panic(err)
	}
	if s.c.Bnj2019.AdminCheck != 0 {
		tmp := make(map[int64]struct{}, len(s.c.Bnj2019.Admins))
		for _, mid := range s.c.Bnj2019.Admins {
			tmp[mid] = struct{}{}
		}
		s.bnjAdmins = tmp
	}
	go s.bnjTimeproc()
	return s
}

// Close .
func (s *Service) Close() {
	s.dao.Close()
	s.resetPub.Close()
}

func (s *Service) timeResetproc() {
	for {
		time.Sleep(time.Second)
		if s.timeFinish != 0 {
			log.Info("timeResetproc finish")
			break
		}
		if s.timeReset == 1 {
			// sub databus
			msg := &bnjmdl.ResetMsg{Mid: s.resetMid, Ts: time.Now().Unix()}
			if err := s.resetPub.Send(context.Background(), strconv.FormatInt(s.resetMid, 10), msg); err != nil {
				log.Error("timeResetproc s.resetPub.Send(%+v) error(%v)", msg, err)
			}
			atomic.StoreInt64(&s.timeReset, 0)
			atomic.StoreInt64(&s.resetMid, 0)
		}
	}
}

func (s *Service) timeFinishproc() {
	for {
		time.Sleep(time.Second)
		if value, err := s.dao.CacheTimeFinish(context.Background()); err != nil {
			log.Error("timeFinishproc s.dao.CacheTimeFinish error(%v)")
		} else if value > 0 {
			log.Info("timeFinishproc cache value finish")
			atomic.StoreInt64(&s.timeFinish, value)
		}
	}
}

func (s *Service) bnjTimeproc() {
	for {
		time.Sleep(time.Second)
		if time.Now().Unix() > s.c.Bnj2019.Start.Unix() {
			go s.timeResetproc()
			go s.timeFinishproc()
			go s.bnjResetCDproc()
			go s.bnjArcproc()
			log.Info("bnjTimeproc start")
			break
		}
	}
}

func (s *Service) bnjResetCDproc() {
	for {
		time.Sleep(time.Second)
		lid := s.c.Bnj2019.SubID
		scoreMap, err := s.likeDao.LikeActLidCounts(context.Background(), []int64{lid})
		if err != nil || scoreMap == nil {
			log.Error("bnjScoreproc s.likeDao.LikeActLidCounts(%d) error(%v)", lid, err)
			continue
		}
		if score, ok := scoreMap[lid]; ok {
			if score >= s.c.Bnj2019.Reward[len(s.c.Bnj2019.Reward)-1].Condition && s.resetCD != _lastCD {
				atomic.StoreInt32(&s.resetCD, _lastCD)
				log.Info("bnjResetCDproc finish")
			}
			if score > s.likeCount {
				atomic.StoreInt64(&s.likeCount, score)
			}
		}
	}
}

func (s *Service) bnjArcproc() {
	for {
		time.Sleep(time.Second)
		now := time.Now().Unix()
		var aids []int64
		for _, v := range s.c.Bnj2019.Info {
			if v.Publish.Unix() < now {
				if v.Aid > 0 {
					aids = append(aids, v.Aid)
				}
			}
		}
		if len(aids) > 0 {
			if arcsReply, err := s.arcClient.Arcs(context.Background(), &arcclient.ArcsRequest{Aids: aids}); err != nil {
				log.Error("bnjArcproc s.arcClient.Arcs(%v) error(%v)", aids, err)
			} else if len(arcsReply.Arcs) > 0 {
				tmp := make(map[int64]*arcclient.Arc, len(aids))
				for _, aid := range aids {
					if arc, ok := arcsReply.Arcs[aid]; ok && arc != nil {
						tmp[aid] = arc
					} else {
						log.Error("bnjArcproc aid(%d) data(%v)", aid, arc)
						continue
					}
				}
				s.previewArcs = tmp
			}
		}
		log.Error("bnjArcproc aids(%v) conf error", aids)
	}
}
