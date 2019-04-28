package like

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/interface/main/activity/conf"
	dao "go-common/app/interface/main/activity/dao/like"
	match "go-common/app/interface/main/activity/model/like"
	coinmdl "go-common/app/service/main/coin/model"
	suitmdl "go-common/app/service/main/usersuit/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

const (
	_matchTable   = "act_matchs"
	_objectTable  = "act_matchs_object"
	_userLogTable = "act_match_user_log"
	_reason       = "参与竞猜"
)

var (
	_emptyMatch   = make([]*match.Match, 0)
	_emptyObjects = make([]*match.Object, 0)
	_emptyUserLog = make([]*match.UserLog, 0)
	_emptyFollow  = make([]string, 0)
)

// Match get match.
func (s *Service) Match(c context.Context, sid int64) (rs []*match.Match, err error) {
	// get from  cache.
	if rs, err = s.dao.ActMatchCache(c, sid); err != nil || len(rs) == 0 {
		if rs, err = s.dao.ActMatch(c, sid); err != nil {
			log.Error("s.dao.Match sid(%d)  error(%v)", sid, err)
			return
		}
		if len(rs) == 0 {
			rs = _emptyMatch
			return
		}
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetActMatchCache(c, sid, rs)
		})
	}
	return
}

// AddGuess add match guess.
func (s *Service) AddGuess(c context.Context, mid int64, p *match.ParamAddGuess) (rs int64, err error) {
	var (
		object           *match.Object
		userGuess        []*match.UserLog
		group            *errgroup.Group
		coinErr, suitErr error
		count            float64
		ip               = metadata.String(c, metadata.RemoteIP)
	)
	if p.Stake > conf.Conf.Rule.MaxGuessCoin {
		err = ecode.ActivityOverCoin
		return
	}
	//check mid coin count
	if count, err = s.coin.UserCoins(c, &coinmdl.ArgCoinInfo{Mid: mid, RealIP: ip}); err != nil {
		dao.PromError("UserCoins接口错误", "s.coin.UserCoins(%d,%s) error(%v)", mid, ip, err)
		return
	}
	if count < float64(p.Stake) {
		err = ecode.ActivityNotEnoughCoin
		return
	}
	// get from  cache.
	if object, err = s.dao.ObjectCache(c, p.ObjID); err != nil || object == nil {
		if object, err = s.dao.Object(c, p.ObjID); err != nil {
			log.Error("s.dao.Match id(%d)  error(%v)", p.ObjID, err)
			return
		}
		if object == nil || object.ID == 0 {
			err = ecode.ActivityNotExist
			return
		}
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetObjectCache(c, p.ObjID, object)
		})
	}
	if time.Now().Unix() < object.Stime.Time().Unix() {
		err = ecode.ActivityNotStart
		return
	} else if object.Result > 0 || time.Now().Unix() > object.Etime.Time().Unix() {
		err = ecode.ActivityOverEnd
		return
	}
	sid := object.Sid
	if userGuess, err = s.ListGuess(c, sid, mid); err != nil {
		log.Error("s.ListGuess(%d,%d) error(%v)", sid, mid, err)
		return
	}
	for _, userLog := range userGuess {
		if userLog.MOId == p.ObjID {
			err = ecode.ActivityHaveGuess
			return
		}
	}
	if rs, err = s.dao.AddGuess(c, mid, object.MatchId, p.ObjID, sid, p.Result, p.Stake); err != nil || rs == 0 {
		log.Error("s.dao.AddGuess matchID(%d) objectID(%d) sid(%d) error(%v)", object.MatchId, p.ObjID, sid, err)
		return
	}
	s.dao.DelUserLogCache(context.Background(), sid, mid)
	group, errCtx := errgroup.WithContext(c)
	if len(s.c.Rule.SuitPids) > 0 && len(userGuess)+1 == s.c.Rule.GuessCount {
		for _, v := range s.c.Rule.SuitPids {
			pid := v
			group.Go(func() error {
				mids := []int64{mid}
				if suitErr = s.suit.GrantByMids(errCtx, &suitmdl.ArgGrantByMids{Mids: mids, Pid: pid, Expire: s.c.Rule.SuitExpire}); suitErr != nil {
					dao.PromError("GrantByMids接口错误", "s.suit.GrantByMids(%d,%d,%s) error(%v)", mid, p.Stake, ip, suitErr)
				}
				return nil
			})
		}
	}
	group.Go(func() error {
		loseCoin := float64(-p.Stake)
		if _, coinErr = s.coin.ModifyCoin(errCtx, &coinmdl.ArgModifyCoin{Mid: mid, Count: loseCoin, Reason: _reason, IP: ip}); coinErr != nil {
			dao.PromError("ModifyCoin接口错误", "s.coin.ModifyCoin(%d,%d,%s) error(%v)", mid, p.Stake, ip, coinErr)
		}
		return nil
	})
	if s.c.Rule.MatchLotteryID > 0 {
		group.Go(func() error {
			if lotteryErr := s.dao.AddLotteryTimes(errCtx, s.c.Rule.MatchLotteryID, mid); lotteryErr != nil {
				log.Error("s.dao.AddLotteryTimes(%d,%d) error(%+v)", s.c.Rule.MatchLotteryID, mid, lotteryErr)
			}
			return nil
		})
	}
	group.Wait()
	return
}

// ListGuess get match guess list.
func (s *Service) ListGuess(c context.Context, sid, mid int64) (rs []*match.UserLog, err error) {
	// get from  cache.
	if rs, err = s.dao.UserLogCache(c, sid, mid); err != nil || len(rs) == 0 {
		if rs, err = s.dao.ListGuess(c, sid, mid); err != nil {
			log.Error("s.dao.ListGuess sid(%d) mid(%d)  error(%v)", sid, mid, err)
			return
		}
		if len(rs) == 0 {
			rs = _emptyUserLog
			return
		}
	}
	var (
		moIDs   []int64
		objects map[int64]*match.Object
	)
	for _, v := range rs {
		moIDs = append(moIDs, v.MOId)
	}
	if len(moIDs) == 0 {
		return
	}
	if objects, err = s.dao.MatchSubjects(c, moIDs); err == nil {
		for _, v := range rs {
			if obj, ok := objects[v.MOId]; ok {
				v.HomeName = obj.HomeName
				v.AwayName = obj.AwayName
				v.ObjResult = obj.Result
				v.GameStime = obj.GameStime
			}
		}
	}
	s.cache.Do(c, func(c context.Context) {
		s.dao.SetUserLogCache(c, sid, mid, rs)
	})
	return
}

// Guess user guess
func (s *Service) Guess(c context.Context, mid int64, p *match.ParamSid) (rs *match.UserGuess, err error) {
	var (
		userGuess           []*match.UserLog
		totalCont, winCount int64
	)
	if userGuess, err = s.ListGuess(c, p.Sid, mid); err != nil {
		log.Error("s.ListGuess(%d,%d) error(%v)", p.Sid, mid, err)
		return
	}
	for _, guess := range userGuess {
		if guess.ObjResult > 0 {
			if guess.Result == guess.ObjResult {
				winCount++
			}
			totalCont++
		}
	}
	rs = new(match.UserGuess)
	rs.Total = totalCont
	rs.Win = winCount
	return
}

// ClearCache del match and object cache
func (s *Service) ClearCache(c context.Context, msg string) (err error) {
	var m struct {
		Table string `json:"table"`
		New   struct {
			ID    int64 `json:"id"`
			Sid   int64 `json:"sid"`
			MatID int64 `json:"match_id"`
			Mid   int64 `json:"mid"`
			MOId  int64 `json:"m_o_id"`
		} `json:"new,omitempty"`
	}
	if err = json.Unmarshal([]byte(msg), &m); err != nil {
		log.Error("ClearCache json.Unmarshal msg(%s) error(%v)", msg, err)
		return
	}
	log.Info("ClearCache json.Unmarshal msg(%s)", msg)
	if m.Table == _matchTable {
		if err = s.dao.DelActMatchCache(c, m.New.Sid, m.New.ID); err != nil {
			log.Error("s.dao.DelActMatchCache sid(%d) matchID(%d)  error(%v)", m.New.Sid, m.New.ID, err)
		}
	} else if m.Table == _objectTable {
		if err = s.dao.DelObjectCache(c, m.New.ID, m.New.Sid); err != nil {
			log.Error("s.dao.DelObjectCache objID(%d)  Sid(%d)  error(%v)", m.New.ID, m.New.Sid, err)
		}
	} else if m.Table == _userLogTable {
		if err = s.dao.DelUserLogCache(c, m.New.Sid, m.New.Mid); err != nil {
			log.Error("s.dao.DelUserLogCache mid(%d) error(%v)", m.New.Mid, err)
		}
	}
	return
}

// AddFollow add match follow
func (s *Service) AddFollow(c context.Context, mid int64, teams []string) (err error) {
	if err = s.dao.AddFollow(c, mid, teams); err != nil {
		log.Error("s.dao.AddFollow mid(%d) teams(%v)  error(%v)", mid, teams, err)
	}
	return
}

// Follow get match follow
func (s *Service) Follow(c context.Context, mid int64) (res []string, err error) {
	if res, err = s.dao.Follow(c, mid); err != nil {
		log.Error("s.dao.Follow mid(%d)  error(%v)", mid, err)
	}
	if len(res) == 0 {
		res = _emptyFollow
	}
	return
}

// ObjectsUnStart get unstart object list.
func (s *Service) ObjectsUnStart(c context.Context, mid int64, p *match.ParamObject) (rs []*match.Object, count int, err error) {
	var (
		userGuess []*match.UserLog
		objects   []*match.Object
		start     = (p.Pn - 1) * p.Ps
		end       = start + p.Ps - 1
	)
	// get from  cache.
	if rs, count, err = s.dao.ObjectsCache(c, p.Sid, start, end); err != nil || len(rs) == 0 {
		if objects, err = s.dao.ObjectsUnStart(c, p.Sid); err != nil {
			log.Error("s.dao.ObjectsUnStart id(%d)  error(%v)", p.Sid, err)
			return
		}
		count = len(objects)
		if count == 0 || count < start {
			rs = _emptyObjects
			return
		}
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetObjectsCache(c, p.Sid, objects, count)
		})
		if count > end+1 {
			rs = objects[start : end+1]
		} else {
			rs = objects[start:]
		}
	}
	if mid > 0 {
		if userGuess, err = s.ListGuess(c, p.Sid, mid); err != nil {
			log.Error("s.ListGuess(%d,%d) error(%v)", p.Sid, mid, err)
			err = nil
		}
		for _, rsObj := range rs {
			for _, guess := range userGuess {
				if rsObj.ID == guess.MOId {
					rsObj.UserResult = guess.Result
					break
				}
			}
		}
	}
	return
}
