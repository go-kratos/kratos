package bws

import (
	"context"
	"encoding/json"

	bwsmdl "go-common/app/interface/main/activity/model/bws"
	"go-common/library/ecode"
	"go-common/library/log"
)

// RedisInfo .
func (s *Service) RedisInfo(c context.Context, loginMid, mid int64, key, day, typ string, del int) (data json.RawMessage, err error) {
	if !s.isAdmin(loginMid) {
		err = ecode.ActivityNotAdmin
		return
	}
	if key == "" {
		if key, err = s.midToKey(c, mid); err != nil {
			return
		}
	}
	var (
		bid int64 = 1
		bs  []byte
	)
	switch typ {
	case "point":
		if del == 1 {
			err = s.dao.DelCachePoints(c, bid)
			return
		}
		var points *bwsmdl.Points
		if points, err = s.dao.CachePoints(c, bid); err != nil || points == nil || len(points.Points) == 0 {
			log.Error("RedisInfo point error (%v)", err)
			return
		}
		if bs, err = json.Marshal(points.Points); err != nil {
			log.Error("RedisInfo point json error (%v)", err)
			return
		}
		data = json.RawMessage(bs)
	case "achieve":
		if del == 1 {
			err = s.dao.DelCacheAchievements(c, bid)
			return
		}
		var achieves *bwsmdl.Achievements
		if achieves, err = s.dao.CacheAchievements(c, bid); err != nil || achieves == nil || len(achieves.Achievements) == 0 {
			log.Error("RedisInfo achieve error (%v)", err)
			return
		}
		if bs, err = json.Marshal(achieves.Achievements); err != nil {
			log.Error("RedisInfo achieve json error (%v)", err)
			return
		}
		data = json.RawMessage(bs)
	case "user_point":
		if del == 1 {
			err = s.dao.DelCacheUserPoints(c, bid, key)
			return
		}
		var res []*bwsmdl.UserPoint
		if res, err = s.dao.CacheUserPoints(c, bid, key); err != nil {
			log.Error("RedisInfo user point key(%s) error (%v)", key, err)
			return
		}
		if bs, err = json.Marshal(res); err != nil {
			log.Error("RedisInfo user point json error (%v)", err)
			return
		}
		data = json.RawMessage(bs)
	case "user_achieve":
		if del == 1 {
			err = s.dao.DelCacheUserAchieves(c, bid, key)
			return
		}
		var res []*bwsmdl.UserAchieve
		if res, err = s.dao.CacheUserAchieves(c, bid, key); err != nil {
			log.Error("RedisInfo user achieve key(%s) error (%v)", key, err)
			return
		}
		if bs, err = json.Marshal(res); err != nil {
			log.Error("RedisInfo user achieve json error (%v)", err)
			return
		}
		data = json.RawMessage(bs)
	case "achieve_cnt":
		if day == "" {
			day = today()
		}
		if del == 1 {
			err = s.dao.DelCacheAchieveCounts(c, bid, day)
			return
		}
		var res []*bwsmdl.CountAchieves
		if res, err = s.dao.CacheAchieveCounts(c, bid, day); err != nil {
			log.Error("RedisInfo achieve_cnt day(%s) error (%v)", day, err)
			return
		}
		if bs, err = json.Marshal(res); err != nil {
			log.Error("RedisInfo achieve_cnt json error (%v)", err)
			return
		}
		data = json.RawMessage(bs)
	case "achieve_cnt_db":
		if day == "" {
			day = today()
		}
		var res []*bwsmdl.CountAchieves
		if res, err = s.dao.RawAchieveCounts(c, bid, day); err != nil {
			log.Error("RedisInfo achieve_cnt_db day(%s) error (%v)", day, err)
			return
		}
		if bs, err = json.Marshal(res); err != nil {
			log.Error("RedisInfo achieve_cnt_db json error (%v)", err)
			return
		}
		data = json.RawMessage(bs)
	default:
		err = ecode.RequestErr
	}
	return
}

// KeyInfo .
func (s *Service) KeyInfo(c context.Context, loginMid, keyID, mid int64, key, typ string, del int) (data json.RawMessage, err error) {
	if !s.isAdmin(loginMid) {
		err = ecode.ActivityNotAdmin
		return
	}
	var (
		bs []byte
	)
	switch typ {
	case "id":
		if keyID == 0 {
			err = ecode.RequestErr
			return
		}
		var user *bwsmdl.Users
		if user, err = s.dao.UserByID(c, keyID); err != nil {
			return
		}
		if bs, err = json.Marshal(user); err != nil {
			return
		}
		data = json.RawMessage(bs)
	case "mid":
		if mid == 0 {
			err = ecode.RequestErr
			return
		}
		if del == 1 {
			err = s.dao.DelCacheUsersMid(c, mid)
			return
		}
		var res *bwsmdl.Users
		if res, err = s.dao.CacheUsersMid(c, mid); err != nil {
			return
		}
		if bs, err = json.Marshal(res); err != nil {
			return
		}
		data = json.RawMessage(bs)
	case "key":
		if key == "" {
			err = ecode.RequestErr
			return
		}
		if del == 1 {
			err = s.dao.DelCacheUsersKey(c, key)
			return
		}
		var res *bwsmdl.Users
		if res, err = s.dao.CacheUsersKey(c, key); err != nil {
			return
		}
		if bs, err = json.Marshal(res); err != nil {
			return
		}
		data = json.RawMessage(bs)
	default:
		err = ecode.RequestErr
	}
	return
}

// AdminInfo get admin info.
func (s *Service) AdminInfo(c context.Context, bid, mid int64) (data *bwsmdl.AdminInfo, err error) {
	data = new(bwsmdl.AdminInfo)
	if s.isAdmin(mid) {
		data.IsAdmin = true
	}
	var points *bwsmdl.Points
	if points, err = s.dao.Points(c, bid); err != nil || points == nil || len(points.Points) == 0 {
		log.Error("s.dao.Points error(%v)", err)
		err = ecode.ActivityPointFail
		return
	}
	for _, v := range points.Points {
		if v.Ower == mid {
			data.Point = v
		}
	}
	if data.Point == nil {
		data.Point = struct{}{}
	}
	return
}
