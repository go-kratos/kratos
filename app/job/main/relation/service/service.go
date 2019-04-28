package service

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go-common/app/job/main/relation/conf"
	"go-common/app/job/main/relation/dao"
	"go-common/app/job/main/relation/model"
	sml "go-common/app/service/main/relation/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"
)

const (
	_relationFidTable     = "user_relation_fid_"
	_relationMidTable     = "user_relation_mid_"
	_relationStatTable    = "user_relation_stat_"
	_relationTagUserTable = "user_relation_tag_user_"
	_retry                = 5
	_retrySleep           = time.Second * 1
)

// Service struct of service.
type Service struct {
	c      *conf.Config
	dao    *dao.Dao
	ds     *databus.Databus
	waiter *sync.WaitGroup
	// monitor
	mo int64
	// moni *monitor.Service
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:      c,
		dao:    dao.New(c),
		ds:     databus.New(c.DataBus),
		waiter: new(sync.WaitGroup),
		// moni:   monitor.New(),
	}
	for i := 0; i < 50; i++ {
		s.waiter.Add(1)
		go s.subproc()
	}
	go s.checkConsume()
	return
}

func (s *Service) subproc() {
	defer s.waiter.Done()
	for {
		var (
			ok  bool
			err error
			res *databus.Message
		)
		if res, ok = <-s.ds.Messages(); !ok {
			log.Error("s.ds.Messages() failed")
			return
		}
		log.Info("received message: res: %+v", res)
		mu := &model.Message{}
		if err = json.Unmarshal(res.Value, mu); err != nil {
			log.Error("go-common/app/job/main/relation,json.Unmarshal (%v) error(%v)", res.Value, err)
			continue
		}
		atomic.AddInt64(&s.mo, 1)
		for i := 0; ; i++ {
			switch {
			case strings.HasPrefix(mu.Table, _relationStatTable):
				err = s.stat(mu.Action, mu.New, mu.Old)
			case strings.HasPrefix(mu.Table, _relationMidTable):
				err = s.relationMid(mu.Action, mu.New, mu.Old)
			case strings.HasPrefix(mu.Table, _relationFidTable):
				err = s.relationFid(mu.Action, mu.New, mu.Old)
			case strings.HasPrefix(mu.Table, _relationTagUserTable):
				err = s.tagUser(mu.New)
			}
			if err != nil {
				i++
				log.Error("s.flush data(%s) error(%+v)", mu.New, err)
				time.Sleep(_retrySleep)
				if i > _retry {
					// if s.c.Env == "prod" {
					// 	s.moni.Sms(context.TODO(), s.c.Sms.Phone, s.c.Sms.Token, fmt.Sprintf("relation-job update cache fail: %v", err))
					// }
					break
				}
				continue
			}
			break
		}
		log.Info("consumer action:%v, table:%v,new :%s", mu.Action, mu.Table, mu.New)
		res.Commit()
	}
}
func (s *Service) tagUser(newMsg []byte) (err error) {
	var tags struct {
		Fid int64 `json:"fid"`
		Mid int64 `json:"mid"`
	}
	if err = json.Unmarshal(newMsg, &tags); err != nil {
		log.Error("json.Unmarshal err(%v)", err)
		return
	}
	f, err := s.dao.UserRelation(context.TODO(), tags.Mid, tags.Fid)
	if err != nil || f == nil {
		return
	}
	s.dao.DelTagsCache(context.TODO(), tags.Mid)
	s.upFollowingCache(context.TODO(), tags.Mid, f)
	return
}

func (s *Service) upFollowingCache(c context.Context, mid int64, f *sml.Following) (err error) {
	if f.Attribute == 0 {
		s.dao.DelFollowing(c, mid, f)
	} else {
		if err = s.dao.AddFollowingCache(c, mid, f); err != nil {
			return
		}
	}
	if err = s.dao.DelFollowingCache(c, mid); err != nil {
		return
	}
	return s.dao.DelTagCountCache(c, mid)
}

// relationFid
func (s *Service) relationFid(action string, nwMsg []byte, oldMsg []byte) error {
	var or *model.Relation
	mr := &model.Relation{}
	if err := json.Unmarshal(nwMsg, mr); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", nwMsg, err)
		return err
	}

	if len(oldMsg) > 0 {
		or = new(model.Relation)
		if err := json.Unmarshal(oldMsg, or); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", oldMsg, err)
		}
	}

	// step 1: add notify
	// if err := s.dao.AddNotify(context.TODO(), mr.Fid); err != nil {
	// 	log.Error("Failed to s.dao.AddNotify(%v): %+v", mr.Fid, err)
	// }

	// step 2: handle recent followers
	// s.RecentFollowers(action, mr, or)

	// step 3: delete cache
	return s.dao.DelFollowerCache(mr.Fid)
}

// relationMid
func (s *Service) relationMid(action string, nwMsg []byte, oldMsg []byte) (err error) {
	mr := &model.Relation{}
	if err = json.Unmarshal(nwMsg, mr); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", nwMsg, err)
		return
	}
	f := &sml.Following{
		Mid:       mr.Fid,
		Attribute: mr.Attribute,
		MTime:     xtime.Time(time.Now().Unix()),
	}
	if err = s.upFollowingCache(context.TODO(), mr.Mid, f); err != nil {
		return
	}
	log.Info("Succeed to update following cache: mid: %d: mr: %+v", mr.Mid, mr)

	// TODO: del this, just for special attention del old cache.
	s.dao.DelTagsCache(context.TODO(), mr.Mid)
	log.Info("Succeed to delete tags cache: mid: %d", mr.Mid)

	// s.dao.AddNotify(context.TODO(), mr.Mid)
	return
}

// stat
func (s *Service) stat(action string, nwMsg []byte, oldMsg []byte) (err error) {
	ms := &model.Stat{}
	if err = json.Unmarshal(nwMsg, ms); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", nwMsg, err)
		return
	}
	mo := &model.Stat{}
	if len(oldMsg) > 0 {
		if err = json.Unmarshal(oldMsg, mo); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", oldMsg, err)
			err = nil
		}
	}

	if ms.Follower > mo.Follower {
		s.dao.FollowerAchieve(context.TODO(), ms.Mid, ms.Follower)
	}
	return s.dao.DelStatCache(ms.Mid)
}

// check consumer stat
func (s *Service) checkConsume() {
	if s.c.Env != "prod" {
		return
	}
	var reMo int64
	for {
		time.Sleep(1 * time.Minute)
		atomic.AddInt64(&s.mo, -reMo)
		// if atomic.AddInt64(&s.mo, -reMo) == 0 {
		// s.moni.Sms(context.TODO(), s.c.Sms.Phone, s.c.Sms.Token, "relation-job stat did not consume within a minute")
		// }
		reMo = atomic.LoadInt64(&s.mo)
	}
}

// Ping check server ok
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close kafka consumer close.
func (s *Service) Close() (err error) {
	s.ds.Close()
	return
}

// Wait wait for service exit.
func (s *Service) Wait() {
	s.waiter.Wait()
}
