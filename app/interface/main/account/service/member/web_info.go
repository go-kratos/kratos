package member

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"go-common/app/interface/main/account/model"
	cmodel "go-common/app/service/main/coin/model"
	lmodel "go-common/app/service/main/location/model"
	mmodel "go-common/app/service/main/member/model"
	pmodel "go-common/app/service/main/passport/model"
	smodel "go-common/app/service/main/secure/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/time"
)

// UpdateSettings Update Settings
func (s *Service) UpdateSettings(c context.Context, mid int64, settings *model.Settings) (err error) {
	var (
		mb *mmodel.Member
		ip = metadata.String(c, metadata.RemoteIP)
	)
	if mb, err = s.memRPC.Member(c, &mmodel.ArgMemberMid{Mid: mid, RemoteIP: ip}); err != nil {
		log.Error("s.Member(%d) error(%v)", mid, err)
		return
	}
	if settings.Sex != sexName(mb.Sex) {
		if err = s.accDao.UpdateSex(c, mid, sexID(settings.Sex), ip); err != nil {
			return
		}
	}
	if settings.Uname != mb.Name {
		if err = s.updateName(c, mid, settings.Uname, ip); err != nil {
			return
		}
	}
	if settings.Sign != mb.Sign {
		if err = s.UpdateSign(c, mid, settings.Sign); err != nil {
			return
		}
	}
	if settings.Birthday != mb.Birthday.Time().Format("2006-01-02") {
		if err = s.accDao.UpdateBirthday(c, mid, ip, settings.Birthday); err != nil {
			return
		}
	}
	return
}

// SettingsInfo Settings Info
func (s *Service) SettingsInfo(c context.Context, mid int64) (user *model.User, err error) {
	var (
		mb       *mmodel.Member
		nickFree *model.NickFree
		userID   string
		ip       = metadata.String(c, metadata.RemoteIP)
	)
	user = &model.User{}
	if mb, err = s.memRPC.Member(c, &mmodel.ArgMemberMid{Mid: mid, RemoteIP: ip}); err != nil {
		log.Error("s.Member(%d) error(%v)", mid, err)
		return
	}
	if nickFree, err = s.NickFree(c, mid); err != nil {
		return
	}
	if userID, err = s.passDao.UserID(c, mid, ip); err != nil {
		return
	}
	user.Userid = userID
	user.Mid = mb.Mid
	user.Uname = mb.Name
	user.Birthday = mb.Birthday.Time().Format("2006-01-02")
	user.NickFree = nickFree.NickFree
	user.Sign = mb.Sign
	user.Sex = sexName(mb.Sex)
	// 尝试清理一次缓存减少用户反馈
	s.NotityPurgeCache(c, mid, "updateUname")
	return
}

// LogCoin Log Money
func (s *Service) LogCoin(c context.Context, mid int64) (logCoin *model.LogCoins, err error) {
	var (
		logs []*cmodel.Log
	)
	if logs, err = s.coinRPC.UserLog(c, &cmodel.ArgLog{Mid: mid, Recent: true, Translate: true}); err != nil {
		return
	}
	logCoin = &model.LogCoins{}
	logCoin.Count = len(logs)
	logCoin.List = make([]*model.LogCoin, 0)
	for _, l := range logs {
		s := fmt.Sprintf("%0.1f", l.To-l.From)
		delta, _ := strconv.ParseFloat(s, 64)
		model := &model.LogCoin{Time: time.Time(l.TimeStamp).Time().Format("2006-01-02 15:04:05"), Reason: l.Desc, Delta: delta}
		logCoin.List = append(logCoin.List, model)
	}
	return
}

// Coin coin.
func (s *Service) Coin(c context.Context, mid int64) (logCoin *model.Coin, err error) {
	var (
		coin float64
	)
	if coin, err = s.coinRPC.UserCoins(c, &cmodel.ArgCoinInfo{Mid: mid}); err != nil {
		return
	}
	logCoin = &model.Coin{Money: coin}
	return
}

// LogMoral Log Moral
func (s *Service) LogMoral(c context.Context, mid int64) (logMorals *model.LogMorals, err error) {
	var (
		logs                       []*mmodel.UserLog
		moral                      *mmodel.Moral
		toMoral, fromMoral, origin int64
	)
	if logs, err = s.memRPC.MoralLog(c, &mmodel.ArgMemberMid{Mid: mid}); err != nil {
		log.Error("s.memRPC.MoralLog(%v) error(%v)", mid, err)
		return
	}
	if moral, err = s.memRPC.Moral(c, &mmodel.ArgMemberMid{Mid: mid}); err != nil {
		log.Error("s.memRPC.Moral(%v) error(%v)", mid, err)
		return
	}
	logMorals = &model.LogMorals{Count: len(logs), Moral: moral.Moral / 100}
	logMorals.List = make([]*model.LogMoral, 0)
	for _, l := range logs {
		ml := &model.LogMoral{}
		ml.Reason = l.Content["reason"]
		if origin, err = strconv.ParseInt(l.Content["origin"], 10, 64); err != nil {
			log.Error("strconv.ParseInt(%v) error(%v)", l.Content["origin"], err)
			continue
		}
		ml.Origin = model.Origin[origin]
		ml.Time = time.Time(l.TS).Time().Format("2006-01-02 15:04:05")
		if toMoral, err = strconv.ParseInt(l.Content["to_moral"], 10, 64); err != nil {
			log.Error("strconv.ParseInt(%v) error(%v)", l.Content["to_moral"], err)
			continue
		}
		if fromMoral, err = strconv.ParseInt(l.Content["from_moral"], 10, 64); err != nil {
			log.Error("strconv.ParseInt(%v) error(%v)", l.Content["from_moral"], err)
			continue
		}
		delta := float64(toMoral-fromMoral) / float64(100)
		if ml.Delta, err = strconv.ParseFloat(fmt.Sprintf("%0.2f", delta), 64); err != nil {
			log.Error("strconv.ParseFloat(%v) error(%v)", delta, err)
			continue
		}
		logMorals.List = append(logMorals.List, ml)
	}
	return
}

// LogExp Log Exp
func (s *Service) LogExp(c context.Context, mid int64) (logExp *model.LogExps, err error) {
	var (
		logs           []*mmodel.UserLog
		toExp, fromExp float64
		ip             = metadata.String(c, metadata.RemoteIP)
	)
	logExp = &model.LogExps{}
	if logs, err = s.memRPC.Log(c, &mmodel.ArgMid2{Mid: mid, RealIP: ip}); err != nil {
		log.Error("s.memRPC.Log(%v) error(%v)", mid, err)
		return
	}
	logExp.Count = len(logs)
	logExp.List = make([]*model.LogExp, 0)
	for _, l := range logs {
		expLog := &model.LogExp{}
		expLog.Time = time.Time(l.TS).Time().Format("2006-01-02 15:04:05")
		expLog.Reason = l.Content["reason"]
		if toExp, err = strconv.ParseFloat(l.Content["to_exp"], 10); err != nil {
			log.Error("strconv.ParseFloat(%v) error(%v)", l.Content["to_exp"], err)
			continue
		}
		if fromExp, err = strconv.ParseFloat(l.Content["from_exp"], 10); err != nil {
			log.Error("strconv.ParseFloat(%v) error(%v)", l.Content["from_exp"], err)
			continue
		}
		expLog.Delta = toExp - fromExp
		logExp.List = append(logExp.List, expLog)
	}
	return
}

// LogLogin logLogin
func (s *Service) LogLogin(c context.Context, mid int64) (logLogins *model.LogLogins, err error) {
	var (
		logs      []*pmodel.LoginLog
		ips       []string
		locations map[string]*lmodel.Info
		excLogs   []*smodel.Expection
	)
	logLogins = &model.LogLogins{}
	if logs, err = s.passRPC.LoginLogs(c, &pmodel.ArgLoginLogs{Mid: mid, Limit: 30}); err != nil {
		log.Error("s.passRPC.LoginLogs(%v) error(%v)", mid, err)
		return
	}
	if excLogs, err = s.secureRPC.ExpectionLoc(c, &smodel.ArgSecure{Mid: mid}); err != nil {
		log.Error("s.secureRPC.ExpectionLoc(%v) error(%v)", mid, err)
		return
	}
	logLogins.Count = len(logs)
	logLogins.List = make([]*model.LogLogin, 0)

	beRemoved := func(sip string) bool {
		ip := net.ParseIP(sip)
		for _, cidr := range s.removeLoginLogCIDR {
			if cidr.Contains(ip) {
				return true
			}
		}

		lip, ierr := s.locRPC.Info(c, &lmodel.ArgIP{IP: sip})
		if ierr != nil || lip == nil {
			log.Error("Failed to get ip info with ip: %s: %+v", sip, ierr)
			return false
		}
		// 过滤局域网登录
		if lip.Country == "局域网" {
			return true
		}

		return false
	}
	for _, l := range logs {
		sip := string(int64ToIP(l.LoginIP).String())
		if beRemoved(sip) {
			continue
		}
		nl := &model.LogLogin{}
		nl.Status = true
		nl.Time = l.Timestamp
		nl.TimeAt = time.Time(l.Timestamp).Time().Format("2006-01-02 15:04:05")
		if len(excLogs) != 0 {
			for _, e := range excLogs {
				if int64(e.IP) == l.LoginIP && time.Time(l.Timestamp) == e.Time {
					nl.Status = false
					nl.Type = int64(e.FeedBack)
				}
			}
		}
		nl.IP = sip
		ips = append(ips, nl.IP)
		logLogins.List = append(logLogins.List, nl)
	}
	if locations, err = s.locRPC.Infos(c, ips); err != nil {
		log.Error("s.locRPC.Infos(%v) error(%v)", ips, err)
		return
	}
	for _, log := range logLogins.List {
		if addr, ok := locations[log.IP]; ok {
			log.Geo = addr.Country + addr.Province + addr.City + addr.ISP
		}
		log.IP = vagueIP(log.IP)
	}
	return
}

// Reward exp reward.
func (s *Service) Reward(c context.Context, mid int64) (reward *model.Reward, err error) {
	var (
		expStat  *mmodel.ExpStat
		todayExp int64
		ip       = metadata.String(c, metadata.RemoteIP)
	)
	if expStat, err = s.memRPC.Stat(c, &mmodel.ArgMid2{Mid: mid, RealIP: ip}); err != nil {
		log.Error("s.s.memRPC.Stat(%d) error(%v)", mid, err)
		return
	}
	if todayExp, err = s.coinRPC.TodayExp(c, &cmodel.ArgMid{Mid: mid, RealIP: ip}); err != nil {
		log.Error("s.coinRPC.TodayExp(%d) error(%v)", mid, err)
		return
	}
	if todayExp > 50 {
		todayExp = 50
	}
	reward = &model.Reward{}
	reward.Login = expStat.Login
	reward.Share = expStat.Share
	reward.Watch = expStat.Watch
	reward.Coin = todayExp
	return
}

func sexID(sex string) int64 {
	switch sex {
	case "男":
		return 1
	case "女":
		return 2
	default:
		return 0
	}
}

func sexName(sex int64) string {
	switch sex {
	case 1:
		return "男"
	case 2:
		return "女"
	default:
		return "保密"
	}
}

func vagueIP(ip string) string {
	strs := strings.Split(ip, ".")
	if len(strs) != 4 {
		log.Error("error ip (%v)", ip)
		return ""
	}
	strs[2] = "*"
	strs[3] = "*"
	return strs[0] + "." + strs[1] + "." + strs[2] + "." + strs[3]
}

func int64ToIP(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}
