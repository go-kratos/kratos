package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/interface/main/activity/model/like"
	"go-common/app/job/main/activity/model/match"
	"go-common/app/service/main/coin/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_upMatchUserLimit = 100
	_reason           = "竞猜奖励"
)

func (s *Service) upMatchUser(c context.Context, newMsg, oldMsg json.RawMessage) {
	var (
		err         error
		newMatchObj = new(match.ActMatchObj)
		oldMatchObj = new(match.ActMatchObj)
	)
	if err = json.Unmarshal(newMsg, newMatchObj); err != nil {
		log.Error("upMatchUser json.Unmarshal(%s) error(%+v)", newMsg, err)
		return
	}
	if err = json.Unmarshal(oldMsg, oldMatchObj); err != nil {
		log.Error("upMatchUser json.Unmarshal(%s) error(%+v)", oldMsg, err)
		return
	}
	if oldMatchObj.Result == match.ResultNo {
		if newMatchObj.Result == match.ResultNo {
			return
		}
		s.upUsers(c, newMatchObj)
	}
}

// FinishMatch finish match.
func (s *Service) FinishMatch(c context.Context, moID int64) (err error) {
	var matchObj *match.ActMatchObj
	if matchObj, err = s.dao.MatchObjInfo(c, moID); err != nil || matchObj == nil {
		return
	}
	if matchObj.Result == match.ResultNo {
		log.Error("FinishMatch moID(%d) result error", moID)
		err = ecode.RequestErr
		return
	}
	go s.upUsers(context.Background(), matchObj)
	return
}

func (s *Service) upUsers(c context.Context, matchObj *match.ActMatchObj) {
	var (
		matchs []*like.Match
		list   []*match.ActMatchUser
		stake  int64
		err    error
	)
	if matchs, err = s.actRPC.Matchs(c, &like.ArgMatch{Sid: matchObj.SID}); err != nil {
		log.Error("upMatchUser s.actRPC.Matchs(%d) error(%v)", matchObj.SID, err)
		return
	}
	for _, v := range matchs {
		if v.ID == matchObj.MatchID {
			stake = v.Stake
		}
	}
	if stake == 0 {
		log.Error("upMatchUser match_id(%d) not found", matchObj.MatchID)
		return
	}
	for {
		if list, err = s.dao.UnDoMatchUsers(context.Background(), matchObj.ID, _upMatchUserLimit); err != nil {
			time.Sleep(time.Duration(s.c.Interval.QueryInterval))
			continue
		} else if len(list) == 0 {
			log.Info("upMatchUser finish m_o_id(%d)", matchObj.ID)
			return
		}
		var resultMids []int64
		for _, v := range list {
			if v.Result == matchObj.Result {
				count := float64(v.Stake * stake)
				if _, err = s.coinRPC.ModifyCoin(context.Background(), &model.ArgModifyCoin{Mid: v.Mid, Count: count, Reason: _reason}); err != nil {
					log.Error("upMatchUser coin error s.coinRPC.ModifyCoin mid(%d) coin(%v) error(%v)", v.Mid, count, err)
					continue
				}
				log.Info("upMatchUser s.coinRPC.ModifyCoin mid(%d) coin(%v)", v.Mid, count)
				resultMids = append(resultMids, v.Mid)
				time.Sleep(time.Duration(s.c.Interval.CoinInterval))
			} else {
				resultMids = append(resultMids, v.Mid)
			}
		}
		if err = s.dao.UpMatchUserResult(context.Background(), matchObj.ID, resultMids); err != nil {
			continue
		}
		log.Info("upMatchUser s.dao.UpMatchUserResult matchID(%d) mids(%+v)", matchObj.ID, resultMids)
		time.Sleep(time.Duration(s.c.Interval.QueryInterval))
	}
}
