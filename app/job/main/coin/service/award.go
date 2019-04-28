package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/coin/dao"
	"go-common/app/job/main/coin/model"
	accmdl "go-common/app/service/main/account/api"
	coinmdl "go-common/app/service/main/coin/model"
	mmdl "go-common/app/service/main/member/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"
)

var _passportLog = 53

func (s *Service) awardDo(ms []interface{}) {
	for _, m := range ms {
		mu, ok := m.(*model.LoginLog)
		if !ok {
			continue
		}
		if mu.Business != _passportLog {
			continue
		}
		prom.BusinessInfoCount.Incr("award-event")
		err := s.award(context.TODO(), mu.Mid, mu.Timestamp, mu.IP)
		if err != nil {
			log.Error("s.award mid %v err %v", mu.Mid, err)
		}
		log.Info("conmsumer login log, mid:%v,time %d, ip: %s err: %v", mu.Mid, mu.Timestamp, mu.IP, err)
	}
}

func split(msg *databus.Message, data interface{}) int {
	t, ok := data.(*model.LoginLog)
	if !ok {
		return 0
	}
	return int(t.Mid)
}

func newMsg(msg *databus.Message) (res interface{}, err error) {
	loginlog := new(model.LoginLog)
	if err = json.Unmarshal(msg.Value, &loginlog); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
		dao.PromError("loginlog:Unmarshal")
		return
	}
	var t time.Time
	if t, err = time.ParseInLocation("2006-01-02 15:04:05", loginlog.CTime, time.Local); err != nil {
		log.Error("time.parse(%s) error(%v)", msg.Value, err)
		dao.PromError("loginlog:Timeparse")
		return
	}
	loginlog.Timestamp = t.Unix()
	loginlog.RawData = string(msg.Value)
	res = loginlog
	return
}

func newExpMsg(msg *databus.Message) (res interface{}, err error) {
	loginlog := new(model.LoginLog)
	explog := new(model.AddExp)
	if err = json.Unmarshal(msg.Value, &explog); err != nil {
		log.Error("newExpMsg json.Unmarshal(%s) error(%v)", msg.Value, err)
		dao.PromError("loginlog:ExpUnmarshal")
		return
	}
	loginlog.Mid = explog.Mid
	loginlog.Timestamp = explog.Ts
	loginlog.IP = explog.IP
	loginlog.RawData = string(msg.Value)
	res = loginlog
	return
}

func (s *Service) award(c context.Context, mid, ts int64, ip string) (err error) {
	if !s.c.CoinJob.Start || ts < s.c.CoinJob.StartTime {
		return
	}
	if (mid == 0) || (ts == 0) || (ip == "") {
		return
	}
	day := int64(time.Unix(ts, 0).Day())
	var login bool
	for {
		if login, err = s.coinDao.Logined(c, mid, day); err == nil {
			break
		}
		dao.PromError("redis-logined-retry")
		time.Sleep(time.Millisecond * 500)
	}
	if login {
		return
	}
	// false mean first login,
	var base *mmdl.BaseInfoReply
	base, err = s.memRPC.Base(c, &mmdl.MemberMidReq{Mid: mid})
	if err != nil {
		if err == ecode.MemberNotExist {
			return
		}
		log.Errorv(c, log.KV("log", "memRPC"), log.KV("err", err), log.KV("mid", mid))
		dao.PromError("登录奖励member")
		return
	}
	if (base != nil) && (base.Rank == 5000) {
		log.Infov(c, log.KV("log", "add coin but user not member"), log.KV("mid", mid))
		return
	}
	var profile *accmdl.ProfileReply
	profile, _ = s.profile(c, mid)
	if (profile != nil) && (profile.Profile != nil) && (profile.Profile.TelStatus == 0) {
		log.Infov(c, log.KV("log", "login award failed. no telphone"), log.KV("mid", mid))
		return
	}
	if _, err = s.coinRPC.ModifyCoin(c, &coinmdl.ArgModifyCoin{Mid: mid, Count: 1, Reason: "登录奖励", IP: ip, CheckZero: 1}); err != nil {
		dao.PromError("登录奖励RPC")
		return
	}
	for {
		if err = s.coinDao.SetLogin(c, mid, day); err == nil {
			break
		}
		dao.PromError("登录奖励setLogin-retry")
		time.Sleep(time.Millisecond * 500)
	}
	prom.BusinessInfoCount.Incr("award-event-success")
	log.Info("add coin success mid: %+v", mid)
	return
}

func (s *Service) profile(c context.Context, mid int64) (res *accmdl.ProfileReply, err error) {
	arg := &accmdl.MidReq{Mid: mid}
	if res, err = s.accRPC.Profile3(c, arg); err != nil {
		dao.PromError("award:Profile3")
		log.Errorv(c, log.KV("log", "Profile3"), log.KV("err", err))
	}
	return
}
