package service

import (
	"context"
	"encoding/json"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	location "go-common/app/service/main/location/model"
	"go-common/app/service/main/secure/model"
	"go-common/library/log"
	xtime "go-common/library/time"
)

var (
	_serverOutIP = map[string]struct{}{
		"183.131.11.57":   {},
		"58.220.29.45":    {},
		"58.220.29.46":    {},
		"58.220.29.47":    {},
		"121.43.189.63":   {},
		"58.220.29.48":    {},
		"192.73.241.22":   {},
		"58.220.29.41":    {},
		"122.228.103.138": {},
		"120.41.32.13":    {},
		"59.47.225.6":     {},
		"222.73.196.16":   {},
		"222.73.196.17":   {},
		"222.73.196.18":   {},
		"222.73.196.19":   {},
		"222.73.196.20":   {},
		"222.73.196.21":   {},
		"222.73.196.22":   {},
		"222.73.196.23":   {},
		"222.73.196.24":   {},
		"222.73.196.25":   {},
		"222.73.196.26":   {},
		"222.73.196.27":   {},
		"222.73.196.28":   {},
		"222.73.196.29":   {},
		"222.73.196.30":   {},
		"121.46.231.65":   {},
		"121.46.231.66":   {},
		"121.46.231.67":   {},
		"121.46.231.68":   {},
		"121.46.231.69":   {},
		"121.46.231.70":   {},
		"116.211.6.85":    {},
		"116.211.6.86":    {},
	}
)

func (s *Service) changePWDRecord(c context.Context, msg []byte) (err error) {
	l := &model.PWDlog{}
	if err = json.Unmarshal(msg, l); err != nil {
		log.Error("s.changePWDRecord err(%v)", err)
		return
	}
	if msg, _ := s.dao.ExpectionMsg(c, l.Mid); msg != nil {
		s.dao.AddChangePWDRecord(c, l.Mid)
	}
	return
}

func (s *Service) loginLog(c context.Context, action string, new []byte) (err error) {
	var (
		res *location.InfoComplete
		ms  = &model.Log{}
	)
	if err = json.Unmarshal(new, ms); err != nil {
		log.Error("json.Unmarshal err(%v)", err)
		return
	}
	switch action {
	case "insert":
		ipStr := inetNtoA(ms.IP)
		if _, ok := _serverOutIP[ipStr]; ok {
			return
		}
		if strings.HasPrefix(ipStr, "127") {
			return
		}
		if res, err = s.getIPZone(c, ipStr); err != nil || res == nil {
			return
		}
		if res.Country == "局域网" {
			return
		}
		ms.LocationID = res.ZoneID[2]
		ms.Location = res.Country + res.Province
		s.checkExpectLogin(c, ms)
		//_, err = s.dao.AddLoginLog(c, ms)
		err = s.dao.AddLocs(c, ms.Mid, ms.LocationID, int64(ms.Time))
		if err != nil {
			return
		}
	case "default":
	}
	return
}

func (s *Service) getIPZone(c context.Context, ip string) (res *location.InfoComplete, err error) {
	arg := &location.ArgIP{
		IP: ip,
	}
	if res, err = s.locRPC.InfoComplete(c, arg); err != nil {
		log.Error("s.locaRPC err(%v)", err)
	}
	return
}

// GetLoc get loc by mid
func (s *Service) GetLoc(c context.Context, mid int64) (res []int64, err error) {
	rs, err := s.dao.Locs(c, mid)
	if err != nil {
		return
	}
	if len(rs) == 0 {
		return
	}
	res = make([]int64, 0, len(rs))
	for k := range rs {
		res = append(res, k)
	}
	return
}

// Status get user expect login status.
func (s *Service) Status(c context.Context, mid int64, uuid string) (msg *model.Msg, err error) {
	var (
		count            int64
		unnotify, change bool
		lg               *model.Log
	)
	if lg, err = s.dao.ExpectionMsg(c, mid); err != nil {
		return
	}
	if lg == nil {
		return
	}
	msg = &model.Msg{
		Log:    lg,
		Notify: true,
	}
	if unnotify, err = s.dao.UnNotify(c, mid, uuid); err != nil {
		err = nil
		return
	}
	if unnotify {
		msg.Notify = false
		return
	}
	if change, err = s.dao.ChangePWDRecord(c, mid); err != nil {
		log.Error("s.dao.ChangePWDRecord err(%v)", err)
		err = nil
		return
	}
	if change {
		msg.Notify = false
		return
	}
	if count, err = s.dao.Count(c, mid, uuid); err != nil {
		err = nil
		return
	}
	if count >= s.c.Expect.CloseCount {
		msg.Notify = false
		return
	}
	if rand.Int63n(s.c.Expect.Rand) != 0 {
		msg.Notify = false
	}
	return
}

// CloseNotify add unnotify of mid and uuid.
func (s *Service) CloseNotify(c context.Context, mid int64, uuid string) (err error) {
	err = s.dao.AddUnNotify(c, mid, uuid)
	if err != nil {
		return
	}
	return s.dao.AddCount(c, mid, uuid)
}
func (s *Service) checkExpectLogin(c context.Context, l *model.Log) (err error) {
	locs, err := s.commonLoc(c, l.Mid)
	if err != nil {
		log.Error("checkExpectLogin err(%v)", err)
		return
	}
	// no common location.
	if len(locs) < 3 {
		return
	}
	// expection login.
	expect := true
	for _, loc := range locs {
		if loc.LocID == l.LocationID {
			expect = false
			break
		}
	}
	if expect {
		if err = s.addExpectionLog(c, l); err != nil {
			return
		}
		err = s.addExpectionMsg(c, l)
		// del old unnotify record.
		s.delUnNotify(c, l.Mid)
	}
	return
}

func (s *Service) addExpectionLog(c context.Context, l *model.Log) (err error) {
	return s.dao.AddException(c, l)
}

func (s *Service) delUnNotify(c context.Context, mid int64) (err error) {
	return s.dao.DelUnNotify(c, mid)
}

func (s *Service) addExpectionMsg(c context.Context, l *model.Log) (err error) {
	return s.dao.AddExpectionMsg(c, l)
}

// AddFeedBack add user expection login feedback.
func (s *Service) AddFeedBack(c context.Context, mid, ts int64, tp int8, IP string) (err error) {
	var res *location.InfoComplete
	if res, err = s.getIPZone(c, IP); err != nil || res == nil {
		return
	}
	l := &model.Log{
		Location: res.Country + res.Province,
		IP:       inetAtoN(IP),
		Mid:      mid,
		Type:     tp,
		Time:     xtime.Time(ts),
	}
	return s.dao.AddFeedBack(c, l)
}

// ExpectionLoc get user expction login list.
func (s *Service) ExpectionLoc(c context.Context, mid int64) (res []*model.Expection, err error) {
	return s.dao.ExceptionLoc(c, mid)
}

func (s *Service) commonLoc(c context.Context, mid int64) (rcs []*model.Record, err error) {
	//locs, err := s.dao.CommonLoc(c, mid, now)
	locs, _ := s.dao.LocsCache(c, mid)
	if len(locs) == 0 {
		// cache miss ,get from db.
		locs, err = s.dao.Locs(c, mid)
		if err != nil {
			return
		}
		s.dao.AddLocsCache(c, mid, &model.Locs{LocsCount: locs})
	}
	if len(locs) >= int(s.c.Expect.DoubleCheck) {
		if ok, _ := s.dao.DCheckCache(c, mid); ok {
			return
		}
		if ok, _ := s.dao.ChangePWDRecord(c, mid); ok {
			return
		}
		s.dao.DoubleCheck(c, mid)
		s.dao.AddDCheckCache(c, mid)
	}
	for l, c := range locs {
		if c > s.c.Expect.Count {
			rc := &model.Record{LocID: l, Count: c}
			rcs = append(rcs, rc)
		}
	}
	sort.Slice(rcs, func(i, j int) bool { return rcs[i].Count > rcs[j].Count })
	if len(rcs) > int(s.c.Expect.Top) {
		rcs = rcs[:s.c.Expect.Top]
	}
	return
}

// AddLog for test.
func (s *Service) AddLog(c context.Context, mid, ips string) (err error) {
	log := &model.Log{}
	log.Mid, _ = strconv.ParseInt(mid, 10, 64)
	log.IP = inetAtoN(ips)
	log.Time = xtime.Time(time.Now().Unix())
	bs, _ := json.Marshal(log)
	err = s.loginLog(c, "insert", bs)
	return
}

// OftenCheck check user often use ipaddress
func (s *Service) OftenCheck(c context.Context, mid int64, ip string) (often *model.Often, err error) {
	var (
		lip  *location.InfoComplete
		locs []*model.Record
	)
	often = &model.Often{}
	if lip, err = s.getIPZone(c, ip); err != nil || lip == nil {
		return
	}
	if locs, err = s.commonLoc(c, mid); err != nil {
		log.Error("checkExpectLogin err(%v)", err)
		return
	}
	log.Info("user common len:%d", len(locs))
	for _, loc := range locs {
		log.Info("user curr(%+v), common log:(%+v)", lip, loc)
		if loc.LocID == lip.ZoneID[2] {
			often.Result = true
			break
		}
	}
	return
}
