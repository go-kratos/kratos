package unicom

import (
	"bytes"
	"context"
	"crypto/des"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go-common/app/interface/main/app-wall/model"
	"go-common/app/interface/main/app-wall/model/unicom"
	account "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	log "go-common/library/log"
	"go-common/library/queue/databus/report"
)

const (
	_unicomKey     = "unicom"
	_unicomPackKey = "unicom_pack"
)

// InOrdersSync insert OrdersSync
func (s *Service) InOrdersSync(c context.Context, usermob, ip string, u *unicom.UnicomJson, now time.Time) (err error) {
	if !s.iplimit(_unicomKey, ip) {
		err = ecode.AccessDenied
		return
	}
	var result int64
	if result, err = s.dao.InOrdersSync(c, usermob, u, now); err != nil || result == 0 {
		log.Error("unicom_s.dao.OrdersSync (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%v) error(%v) or result==0",
			usermob, u.Cpid, u.Spid, u.TypeInt, u.Ordertypes, u.Channelcode, u.Ordertime, u.Canceltime, u.Endtime, u.Province, u.Area, u.Videoid, err)
	}
	return
}

// InAdvanceSync insert AdvanceSync
func (s *Service) InAdvanceSync(c context.Context, usermob, ip string, u *unicom.UnicomJson, now time.Time) (err error) {
	if !s.iplimit(_unicomKey, ip) {
		err = ecode.AccessDenied
		return
	}
	var result int64
	if result, err = s.dao.InAdvanceSync(c, usermob, u, now); err != nil || result == 0 {
		log.Error("unicom_s.dao.InAdvanceSync (%s,%s,%s,%s,%s,%s,%s,%s,%v) error(%v) or result==0",
			usermob, u.Userphone, u.Cpid, u.Spid, u.Ordertypes, u.Channelcode, u.Province, u.Area, err)
	}
	return
}

// FlowSync update OrdersSync
func (s *Service) FlowSync(c context.Context, flowbyte int, usermob, time, ip string, now time.Time) (err error) {
	if !s.iplimit(_unicomKey, ip) {
		err = ecode.AccessDenied
		return
	}
	var result int64
	if result, err = s.dao.FlowSync(c, flowbyte, usermob, time, now); err != nil || result == 0 {
		log.Error("unicom_s.dao.OrdersSync(%s, %s, %s) error(%v) or result==0", usermob, time, flowbyte, err)
	}
	return
}

// InIPSync
func (s *Service) InIPSync(c context.Context, ip string, u *unicom.UnicomIpJson, now time.Time) (err error) {
	if !s.iplimit(_unicomKey, ip) {
		err = ecode.AccessDenied
		return
	}
	var result int64
	if result, err = s.dao.InIPSync(c, u, now); err != nil {
		log.Error("s.dao.InIpSync(%s,%s) error(%v)", u.Ipbegin, u.Ipend, err)
	} else if result == 0 {
		err = ecode.RequestErr
		log.Error("unicom_s.dao.InIpSync(%s,%s) error(%v) result==0", u.Ipbegin, u.Ipend, err)
	}
	return
}

// UserFlow
func (s *Service) UserFlow(c context.Context, usermob, mobiApp, ip string, build int, now time.Time) (res *unicom.Unicom, msg string, err error) {
	var (
		row       map[string]*unicom.Unicom
		orderType int
	)
	row = s.unicomInfo(c, usermob, now)
	res, msg, err = s.uState(row, usermob, now)
	switch err {
	case ecode.NothingFound:
		orderType = 1
	case ecode.NotModified:
		orderType = 3
	default:
		orderType = 2
	}
	s.unicomInfoc(mobiApp, usermob, ip, build, orderType, now)
	return
}

// unicomInfo
func (s *Service) unicomInfo(c context.Context, usermob string, now time.Time) (res map[string]*unicom.Unicom) {
	var (
		err error
		us  []*unicom.Unicom
	)
	res = map[string]*unicom.Unicom{}
	if us, err = s.dao.UnicomCache(c, usermob); err == nil && len(us) > 0 {
		s.pHit.Incr("unicoms_cache")
	} else {
		if us, err = s.dao.OrdersUserFlow(context.TODO(), usermob); err != nil {
			log.Error("unicom_s.dao.OrdersUserFlow error(%v)", err)
			return
		}
		s.pMiss.Incr("unicoms_cache")
		if len(us) > 0 {
			if err = s.dao.AddUnicomCache(c, usermob, us); err != nil {
				log.Error("s.dao.AddUnicomCache usermob(%v) error(%v)", usermob, err)
				return
			}
		}
	}
	if len(us) > 0 {
		row := &unicom.Unicom{}
		channel := &unicom.Unicom{}
		for _, u := range us {
			if u.TypeInt == 1 && now.Unix() <= int64(u.Endtime) {
				*row = *u
				continue
			} else if u.TypeInt == 0 {
				if int64(row.Ordertime) == 0 {
					*row = *u
				} else if int64(row.Ordertime) < int64(u.Ordertime) {
					continue
				}
				*row = *u
				continue
			} else if u.TypeInt == 1 {
				if int64(row.Ordertime) > int64(u.Ordertime) {
					continue
				}
				*channel = *u
			}
		}
		if row.Spid == 0 && channel.Spid == 0 {
			return
		} else if row.Spid == 0 && channel.Spid > 0 {
			row = channel
		}
		res[usermob] = row
	}
	return
}

// UserState
func (s *Service) UserState(c context.Context, usermob, mobiApp, ip string, build int, now time.Time) (res *unicom.Unicom, msg string, err error) {
	var (
		orderType int
	)
	row := s.unicomInfo(c, usermob, now)
	res, msg, err = s.uState(row, usermob, now)
	switch err {
	case ecode.NothingFound:
		orderType = 1
	case ecode.NotModified:
		orderType = 3
	default:
		orderType = 2
	}
	s.unicomInfoc(mobiApp, usermob, ip, build, orderType, now)
	return
}

// UnicomState
func (s *Service) UnicomState(c context.Context, usermob, mobiApp, ip string, build int, now time.Time) (res *unicom.Unicom, err error) {
	var (
		ok bool
	)
	row := s.unicomInfo(c, usermob, now)
	if res, ok = row[usermob]; !ok {
		res = &unicom.Unicom{Unicomtype: 1}
	} else if res.TypeInt == 1 && now.Unix() > int64(res.Endtime) {
		res.Unicomtype = 3
	} else if res.TypeInt == 1 {
		res.Unicomtype = 4
	} else if res.TypeInt == 0 {
		res.Unicomtype = 2
	}
	log.Info("unicomstate_type:%v unicomstate_type_usermob:%v", res.Unicomtype, usermob)
	s.unicomInfoc(mobiApp, usermob, ip, build, res.Unicomtype, now)
	return
}

// UserFlowState
func (s *Service) UserFlowState(c context.Context, usermob string, now time.Time) (res *unicom.Unicom, err error) {
	row := s.unicomInfo(c, usermob, now)
	var ok bool
	if res, ok = row[usermob]; !ok {
		res = &unicom.Unicom{Unicomtype: 1}
	} else if res.TypeInt == 1 {
		res.Unicomtype = 3
	} else if res.TypeInt == 0 {
		res.Unicomtype = 2
	}
	return
}

// uState
func (s *Service) uState(unicom map[string]*unicom.Unicom, usermob string, now time.Time) (res *unicom.Unicom, msg string, err error) {
	var ok bool
	if res, ok = unicom[usermob]; !ok {
		err = ecode.NothingFound
		msg = "该卡号尚未开通哔哩哔哩专属免流服务"
	} else if res.TypeInt == 1 && now.Unix() > int64(res.Endtime) {
		err = ecode.NotModified
		msg = "该卡号哔哩哔哩专属免流服务已退订且已过期"
	}
	return
}

// IsUnciomIP is unicom ip
func (s *Service) IsUnciomIP(ipUint uint32, ipStr, mobiApp string, build int, now time.Time) (err error) {
	if !model.IsIPv4(ipStr) {
		err = ecode.NothingFound
		return
	}
	isValide := s.unciomIPState(ipUint)
	s.ipInfoc(mobiApp, "", ipStr, build, isValide, now)
	if isValide {
		return
	}
	err = ecode.NothingFound
	return
}

// UserUnciomIP
func (s *Service) UserUnciomIP(ipUint uint32, ipStr, usermob, mobiApp string, build int, now time.Time) (res *unicom.UnicomUserIP) {
	res = &unicom.UnicomUserIP{
		IPStr:    ipStr,
		IsValide: false,
	}
	if !model.IsIPv4(ipStr) {
		return
	}
	if res.IsValide = s.unciomIPState(ipUint); !res.IsValide {
		log.Error("unicom_user_ip:%v unicom_ip_usermob:%v", ipStr, usermob)
	}
	s.ipInfoc(mobiApp, usermob, ipStr, build, res.IsValide, now)
	return
}

// Order unicom user order
func (s *Service) Order(c context.Context, usermobDes, channel string, ordertype int, now time.Time) (res *unicom.BroadbandOrder, msg string, err error) {
	var (
		usermob string
	)
	var (
		_aesKey = []byte("9ed226d9")
	)
	bs, err := base64.StdEncoding.DecodeString(usermobDes)
	if err != nil {
		log.Error("base64.StdEncoding.DecodeString(%s) error(%v)", usermobDes, err)
		return
	}
	bs, err = s.DesDecrypt(bs, _aesKey)
	if err != nil {
		log.Error("unicomSvc.DesDecrypt error(%v)", err)
		return
	}
	if len(bs) > 32 {
		usermob = string(bs[:32])
	} else {
		usermob = string(bs)
	}
	row := s.unicomInfo(c, usermob, now)
	if u, ok := row[usermob]; ok {
		if u.Spid != 979 && u.Spid != 0 && u.TypeInt == 0 {
			err = ecode.NotModified
			msg = "您当前是流量卡并且已生效无法再订购流量包"
			return
		}
	}
	if res, msg, err = s.dao.Order(c, usermobDes, channel, ordertype); err != nil {
		log.Error("s.dao.Order usermobDes(%v) error(%v)", usermobDes, err)
		return
	}
	return
}

// CancelOrder unicom user cancel order
func (s *Service) CancelOrder(c context.Context, usermob string) (res *unicom.BroadbandOrder, msg string, err error) {
	if res, msg, err = s.dao.CancelOrder(c, usermob); err != nil {
		log.Error("s.dao.CancelOrder usermob(%v) error(%v)", usermob, err)
		return
	}
	return
}

// UnicomSMSCode unicom sms code
func (s *Service) UnicomSMSCode(c context.Context, phone string, now time.Time) (msg string, err error) {
	if msg, err = s.dao.SendSmsCode(c, phone); err != nil {
		log.Error("s.dao.SendSmsCode phone(%v) error(%v)", phone, err)
		return
	}
	return
}

// AddUnicomBind unicom user bind
func (s *Service) AddUnicomBind(c context.Context, phone string, code int, mid int64, now time.Time) (msg string, err error) {
	var (
		usermobDes string
		usermob    string
		ub         *unicom.UserBind
		result     int64
		res        *unicom.Unicom
	)
	if usermobDes, msg, err = s.dao.SmsNumber(c, phone, code); err != nil {
		log.Error("s.dao.SmsNumber error(%v)", err)
		return
	}
	if usermobDes == "" {
		err = ecode.NotModified
		msg = "激活失败，请重新输入验证码激活"
		return
	}
	var (
		_aesKey = []byte("9ed226d9")
	)
	bs, err := base64.StdEncoding.DecodeString(usermobDes)
	if err != nil {
		log.Error("base64.StdEncoding.DecodeString(%s) error(%v)", usermobDes, err)
		return
	}
	bs, err = s.DesDecrypt(bs, _aesKey)
	if err != nil {
		log.Error("unicomSvc.DesDecrypt error(%v)", err)
		return
	}
	if len(bs) > 32 {
		usermob = string(bs[:32])
	} else {
		usermob = string(bs)
	}
	row := s.unicomInfo(c, usermob, now)
	if res, msg, err = s.uState(row, usermob, now); err != nil {
		return
	}
	if res.Spid == 979 {
		err = ecode.NotModified
		msg = "该业务只支持哔哩哔哩免流卡"
		return
	}
	if _, err = s.unicomBindInfo(c, mid); err == nil {
		err = ecode.NotModified
		msg = "该账户已绑定过手机号"
		return
	}
	if midtmp := s.unicomBindMIdByPhone(c, phone); midtmp > 0 {
		err = ecode.NotModified
		msg = "该手机号已被注册"
		return
	}
	if ub, err = s.dao.UserBindOld(c, phone); err != nil || ub == nil {
		phoneInt, _ := strconv.Atoi(phone)
		ub = &unicom.UserBind{
			Usermob: usermob,
			Phone:   phoneInt,
			Mid:     mid,
			State:   1,
		}
	} else {
		ub.Mid = mid
		ub.State = 1
	}
	if result, err = s.dao.InUserBind(c, ub); err != nil || result == 0 {
		log.Error("s.dao.InUserBind ub(%v) error(%v) or result==0", ub, err)
		return
	}
	if err = s.dao.AddUserBindCache(c, mid, ub); err != nil {
		log.Error("s.dao.AddUserBindCache error(%v)", err)
		return
	}
	// databus
	s.addUserBindState(&unicom.UserBindInfo{MID: ub.Mid, Phone: ub.Phone, Action: "unicom_welfare_bind"})
	return
}

// ReleaseUnicomBind release unicom bind
func (s *Service) ReleaseUnicomBind(c context.Context, mid int64, phone int) (msg string, err error) {
	var (
		result int64
		ub     *unicom.UserBind
	)
	if ub, err = s.unicomBindInfo(c, mid); err != nil {
		msg = "用户未绑定手机号"
		return
	}
	ub.State = 0
	if err = s.dao.DeleteUserBindCache(c, mid); err != nil {
		log.Error("s.dao.DeleteUserBindCache error(%v)", err)
		return
	}
	if result, err = s.dao.InUserBind(c, ub); err != nil || result == 0 {
		log.Error("s.dao.InUserBind ub(%v) error(%v) or result==0", ub, err)
		return
	}
	// databus
	s.addUserBindState(&unicom.UserBindInfo{MID: ub.Mid, Phone: ub.Phone, Action: "unicom_welfare_untied"})
	return
}

// unicomBindInfo unicom bind info
func (s *Service) unicomBindInfo(c context.Context, mid int64) (res *unicom.UserBind, err error) {
	if res, err = s.dao.UserBindCache(c, mid); err == nil {
		s.pHit.Incr("unicoms_userbind_cache")
	} else {
		if res, err = s.dao.UserBind(c, mid); err != nil {
			log.Error("s.dao.UserBind error(%v)", err)
			return
		}
		s.pMiss.Incr("unicoms_userbind_cache")
		if res == nil {
			err = ecode.NothingFound
			return
		}
		if err = s.dao.AddUserBindCache(c, mid, res); err != nil {
			log.Error("s.dao.AddUserBindCache mid(%d) error(%v)", mid, err)
			return
		}
	}
	return
}

func (s *Service) unicomBindMIdByPhone(c context.Context, phone string) (mid int64) {
	var err error
	if mid, err = s.dao.UserBindPhoneMid(c, phone); err != nil {
		log.Error("s.dao.UserBindPhoneMid error(%v)", phone)
		return
	}
	return
}

// UserBind user bind
func (s *Service) UserBind(c context.Context, mid int64) (res *unicom.UserBind, msg string, err error) {
	var (
		acc *account.Info
		ub  *unicom.UserBind
	)
	if acc, err = s.accd.Info(c, mid); err != nil {
		log.Error("s.accd.info error(%v)", err)
		return
	}
	res = &unicom.UserBind{
		Name: acc.Name,
		Mid:  acc.Mid,
	}
	if ub, err = s.unicomBindInfo(c, mid); err != nil {
		log.Error("s.userBindInfo error(%v)", err)
		err = nil
	}
	if ub != nil {
		res.Phone = ub.Phone
		res.Integral = ub.Integral
		res.Flow = ub.Flow
	}
	return
}

// UnicomPackList unicom pack list
func (s *Service) UnicomPackList() (res []*unicom.UserPack) {
	res = s.unicomPackCache
	return
}

// UnicomPackReceive unicom pack receive
func (s *Service) UnicomPackReceive(c context.Context, mid int64, packID int64, now time.Time) (msg string, err error) {
	var (
		pack             *unicom.UserPack
		userbind         *unicom.UserBind
		requestNo        int64
		unicomOrderID    string
		unicomOutorderID string
		u                *unicom.Unicom
	)
	if userbind, err = s.unicomBindInfo(c, mid); err != nil {
		err = ecode.NotModified
		msg = "用户未绑定手机号"
		return
	}
	row := s.unicomInfo(c, userbind.Usermob, now)
	if u, msg, err = s.uState(row, userbind.Usermob, now); err != nil {
		return
	}
	if u.Spid == 979 {
		err = ecode.NotModified
		msg = "该业务只支持哔哩哔哩免流卡"
		return
	}
	if pack, err = s.unicomPackInfos(c, packID); err != nil {
		err = ecode.NotModified
		msg = "该礼包不存在"
		return
	}
	if pack.Capped != 0 && pack.Amount == 0 {
		err = ecode.NotModified
		msg = "该礼包不存在"
		return
	}
	if userbind.Integral < pack.Integral {
		err = ecode.NotModified
		msg = "福利点不足"
		return
	}
	if requestNo, err = s.seqdao.SeqID(c); err != nil {
		log.Error("unicom_s.seqdao.SeqID error (%v)", err)
		return
	}
	switch pack.Type {
	case 0:
		if err = s.dao.UserFlowWaitCache(c, userbind.Phone); err == nil {
			err = ecode.NotModified
			msg = "请间隔一分钟之后再领取流量包"
			return
		}
		if msg, err = s.dao.FlowPre(c, userbind.Phone, requestNo, now); err != nil {
			log.Error("s.dao.FlowPre error(%v)", err)
			return
		}
		if unicomOrderID, unicomOutorderID, msg, err = s.dao.FlowExchange(c, userbind.Phone, pack.Param, requestNo, now); err != nil {
			log.Error("s.dao.FlowExchange error(%v)", err)
			return
		}
		uf := &unicom.UnicomUserFlow{
			Phone:      userbind.Phone,
			Mid:        mid,
			Integral:   pack.Integral,
			Flow:       0,
			Orderid:    unicomOrderID,
			Outorderid: unicomOutorderID,
			Desc:       pack.Desc,
		}
		key := strconv.Itoa(userbind.Phone) + unicomOutorderID
		if err = s.addUserFlowCache(c, key, uf); err != nil {
			log.Error("s.addUserFlowCache error(%v)", err)
			return
		}
		if err = s.dao.AddUserFlowWaitCache(c, userbind.Phone); err != nil {
			log.Error("s.dao.AddUserFlowWaitCache error(%v)", err)
			return
		}
	case 1:
		var batchID int
		if batchID, err = strconv.Atoi(pack.Param); err != nil {
			log.Error("batchID(%v) strconv.Atoi error(%v)", pack.Param, err)
			msg = "礼包参数错误"
			err = ecode.RequestErr
			return msg, err
		}
		if msg, err = s.accd.AddVIP(c, mid, requestNo, batchID, pack.Desc); err != nil {
			log.Error("s.accd.AddVIP error(%v)", err)
			return msg, err
		}
	case 2:
		var day int
		if day, err = strconv.Atoi(pack.Param); err != nil {
			log.Error("day(%v) strconv.Atoi error(%v)", pack.Param, err)
			msg = "礼包参数错误"
			err = ecode.RequestErr
			return msg, err
		}
		if msg, err = s.live.AddVip(c, mid, day); err != nil {
			log.Error("s.live.AddVip error(%v)", err)
			return "", err
		}
	case 3:
		var acc *account.Info
		if acc, err = s.accd.Info(c, mid); err != nil {
			log.Error("s.accd.info error(%v)", err)
			return
		}
		if msg, err = s.shop.Coupon(c, pack.Param, mid, acc.Name); err != nil {
			log.Error("s.shop.Coupon error(%v)", err)
			return
		}
	}
	var (
		p            = &unicom.UserPack{}
		ub           = &unicom.UserBind{}
		requestNoStr string
		result       int64
	)
	if pack.Capped != 0 {
		*p = *pack
		p.Amount = p.Amount - 1
		if p.Amount == 0 {
			if err = s.dao.DeleteUserPackCache(c, p.ID); err != nil {
				log.Error("s.dao.DeleteUserPackCache error(%v)", err)
				return
			}
			p.State = 0
		} else {
			if err = s.dao.AddUserPackCache(c, p.ID, p); err != nil {
				log.Error("s.dao.AddUserPackCache error(%v)", err)
				return
			}
		}
		if result, err = s.dao.UpUserPacks(c, p, p.ID); err != nil || result == 0 {
			log.Error("s.dao.UpUserPacks error(%v) or result==0", err)
			return
		}
	}
	*ub = *userbind
	ub.Integral = ub.Integral - pack.Integral
	if err = s.updateUserIntegral(c, mid, ub); err != nil {
		log.Error("s.updateUserIntegral error(%v)", err)
		return
	}
	msg = pack.Desc + ",领取成功"
	if unicomOrderID != "" {
		requestNoStr = unicomOrderID
	} else {
		requestNoStr = strconv.FormatInt(requestNo, 10)
	}
	log.Info("unicom_pack(%v) mid(%v)", pack.Desc+",领取成功", userbind.Mid)
	s.unicomPackInfoc(userbind.Usermob, pack.Desc, requestNoStr, userbind.Phone, pack.Integral, int(pack.Capped), userbind.Mid, now)
	ul := &unicom.UserPackLog{
		Phone:     userbind.Phone,
		Usermob:   userbind.Usermob,
		Mid:       userbind.Mid,
		RequestNo: requestNoStr,
		Type:      pack.Type,
		Desc:      pack.Desc,
		UserDesc:  ("您当前已领取" + pack.Desc + "，扣除" + strconv.Itoa(pack.Integral) + "福利点"),
		Integral:  pack.Integral,
	}
	s.addUserPackLog(ul)
	return
}

// UnicomFlowPack unicom flow pack
func (s *Service) UnicomFlowPack(c context.Context, mid int64, flowID string, now time.Time) (msg string, err error) {
	var (
		userbind         *unicom.UserBind
		ub               = &unicom.UserBind{}
		requestNo        int64
		flowDesc         string
		flow             int
		unicomOrderID    string
		unicomOutorderID string
		u                *unicom.Unicom
	)
	if userbind, err = s.unicomBindInfo(c, mid); err != nil {
		err = ecode.NotModified
		msg = "用户未绑定手机号"
		return
	}
	row := s.unicomInfo(c, userbind.Usermob, now)
	if u, msg, err = s.uState(row, userbind.Usermob, now); err != nil {
		return
	}
	if u.Spid == 979 {
		err = ecode.NotModified
		msg = "该业务只支持哔哩哔哩免流卡"
		return
	}
	switch flowID {
	case "01":
		flow = 100
		flowDesc = "100MB流量包"
	case "02":
		flow = 200
		flowDesc = "200MB流量包"
	case "03":
		flow = 300
		flowDesc = "300MB流量包"
	case "04":
		flow = 500
		flowDesc = "500MB流量包"
	case "05":
		flow = 1024
		flowDesc = "1024MB流量包"
	case "06":
		flow = 2048
		flowDesc = "2048MB流量包"
	default:
		err = ecode.RequestErr
		msg = "流量包参数错误"
		return
	}
	if userbind.Flow < flow {
		err = ecode.NotModified
		msg = "可用流量不足"
		return
	}
	if err = s.dao.UserFlowWaitCache(c, userbind.Phone); err == nil {
		err = ecode.NotModified
		msg = "请间隔一分钟之后再领取流量包"
		return
	}
	if requestNo, err = s.seqdao.SeqID(c); err != nil {
		log.Error("unicom_s.seqdao.SeqID error (%v)", err)
		return
	}
	if unicomOrderID, unicomOutorderID, msg, err = s.dao.FlowExchange(c, userbind.Phone, flowID, requestNo, now); err != nil {
		log.Error("s.dao.FlowExchange error(%v)", err)
		return
	}
	uf := &unicom.UnicomUserFlow{
		Phone:      userbind.Phone,
		Mid:        mid,
		Integral:   0,
		Flow:       flow,
		Orderid:    unicomOrderID,
		Outorderid: unicomOutorderID,
		Desc:       flowDesc,
	}
	key := strconv.Itoa(userbind.Phone) + unicomOutorderID
	if err = s.addUserFlowCache(c, key, uf); err != nil {
		log.Error("s.addUserFlowCache error(%v)", err)
		return
	}
	*ub = *userbind
	ub.Flow = ub.Flow - flow
	if err = s.updateUserIntegral(c, mid, ub); err != nil {
		log.Error("s.updateUserIntegral error(%v)", err)
		return
	}
	msg = flowDesc + ",领取成功"
	log.Info("unicom_pack(%v) mid(%v)", flowDesc+",领取成功", userbind.Mid)
	s.unicomPackInfoc(userbind.Usermob, flowDesc, unicomOrderID, userbind.Phone, flow, 0, userbind.Mid, now)
	ul := &unicom.UserPackLog{
		Phone:     userbind.Phone,
		Usermob:   userbind.Usermob,
		Mid:       userbind.Mid,
		RequestNo: unicomOrderID,
		Type:      0,
		Desc:      flowDesc,
		UserDesc:  ("您当前已领取" + flowDesc + "，扣除" + strconv.Itoa(flow) + "MB流量"),
		Integral:  flow,
	}
	s.addUserPackLog(ul)
	if err = s.dao.AddUserFlowWaitCache(c, userbind.Phone); err != nil {
		log.Error("s.dao.AddUserFlowWaitCache error(%v)", err)
		return
	}
	return
}

// UserBindLog user bind week log
func (s *Service) UserBindLog(c context.Context, mid int64, now time.Time) (res []*unicom.UserLog, err error) {
	if res, err = s.dao.SearchUserBindLog(c, mid, now); err != nil {
		log.Error("unicom s.dao.SearchUserBindLog error(%v)", err)
		return
	}
	return
}

// WelfareBindState welfare user bind state
func (s *Service) WelfareBindState(c context.Context, mid int64) (res int) {
	if ub, err := s.dao.UserBindCache(c, mid); err == nil && ub != nil {
		res = 1
	}
	return
}

func (s *Service) updateUserIntegral(c context.Context, mid int64, ub *unicom.UserBind) (err error) {
	var result int64
	if err = s.dao.AddUserBindCache(c, mid, ub); err != nil {
		log.Error("s.dao.AddUserBindCache error(%v)", err)
		return
	}
	if result, err = s.dao.UpUserIntegral(c, ub); err != nil {
		log.Error("s.dao.UpUserIntegral error(%v) ", err)
		return
	}
	if result == 0 {
		log.Error("s.dao.UpUserIntegral result==0")
		err = ecode.NotModified
		return
	}
	return
}

// unicomPackInfos unicom pack infos
func (s *Service) unicomPackInfos(c context.Context, id int64) (res *unicom.UserPack, err error) {
	var (
		ub  *unicom.UserPack
		row = map[int64]*unicom.UserPack{}
		ok  bool
	)
	if ub, err = s.dao.UserPackCache(c, id); err == nil {
		res = ub
		s.pHit.Incr("unicoms_pack_cache")
	} else {
		if row, err = s.dao.UserPackByID(c, id); err != nil {
			log.Error("s.dao.UserBind error(%v)", err)
			return
		}
		s.pMiss.Incr("unicoms_pack_cache")
		ub, ok = row[id]
		if !ok {
			err = ecode.NothingFound
			return
		}
		if err = s.dao.AddUserPackCache(c, id, ub); err != nil {
			log.Error("s.dao.AddUserPackCache id(%d) error(%v)", id, err)
			return
		}
		res = ub
	}
	return
}

// unciomIPState
func (s *Service) unciomIPState(ipUint uint32) (isValide bool) {
	for _, u := range s.unicomIpCache {
		if u.IPStartUint <= ipUint && u.IPEndUint >= ipUint {
			isValide = true
			return
		}
	}
	isValide = false
	return
}

// unicomIp ip limit
func (s *Service) iplimit(k, ip string) bool {
	key := fmt.Sprintf(_initIPlimitKey, k, ip)
	if _, ok := s.operationIPlimit[key]; ok {
		return true
	}
	return false
}

// DesDecrypt
func (s *Service) DesDecrypt(src, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	out = s.zeroUnPadding(out)
	return out, nil
}

// zeroUnPadding
func (s *Service) zeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}

func (s *Service) addUserFlowCache(c context.Context, key string, uf *unicom.UnicomUserFlow) (err error) {
	if err = s.dao.AddUserFlowCache(c, key); err != nil {
		log.Error("s.dao.AddUserFlowCache error(%v)", err)
		return
	}
	var flowList map[string]*unicom.UnicomUserFlow
	if flowList, err = s.dao.UserFlowListCache(c); err != nil {
		log.Error("s.dao.UserFlowListCache error(%v)", err)
		return
	}
	if flowList == nil {
		flowList = map[string]*unicom.UnicomUserFlow{
			key: uf,
		}
	} else {
		flowList[key] = uf
	}
	if err = s.dao.AddUserFlowListCache(c, flowList); err != nil {
		log.Error("s.dao.AddUserFlowListCache error(%v)", err)
		return
	}
	return
}

// UserPacksLog user pack logs
func (s *Service) UserPacksLog(c context.Context, starttime, now time.Time, start int, ip string) (res []*unicom.UserPackLog, err error) {
	if !s.iplimit(_unicomPackKey, ip) {
		err = ecode.AccessDenied
		return
	}
	var (
		endday time.Time
	)
	if starttime.Month() >= now.Month() && starttime.Year() >= now.Year() {
		res = []*unicom.UserPackLog{}
		return
	}
	if endInt := starttime.AddDate(0, 1, -1).Day(); start > endInt {
		res = []*unicom.UserPackLog{}
		return
	} else if start == endInt {
		endday = starttime.AddDate(0, 1, 0)
	} else {
		endday = starttime.AddDate(0, 0, start)
	}
	if res, err = s.dao.UserPacksLog(c, endday.AddDate(0, 0, -1), endday); err != nil {
		log.Error("user pack logs s.dao.UserPacksLog error(%v)", err)
		return
	}
	if len(res) == 0 {
		res = []*unicom.UserPackLog{}
	}
	return
}

func (s *Service) addUserPackLog(u *unicom.UserPackLog) {
	select {
	case s.packLogCh <- u:
	default:
		log.Warn("user pack log buffer is full")
	}
}

func (s *Service) addUserPackLogproc() {
	for {
		i, ok := <-s.packLogCh
		if !ok {
			log.Warn("user pack log proc exit")
			return
		}
		var (
			c      = context.TODO()
			result int64
			err    error
			logID  = 91
		)
		switch v := i.(type) {
		case *unicom.UserPackLog:
			if result, err = s.dao.InUserPackLog(c, v); err != nil || result == 0 {
				log.Error("s.dao.UpUserIntegral error(%v) or result==0", err)
				continue
			}
			report.User(&report.UserInfo{
				Mid:      v.Mid,
				Business: logID,
				Action:   "unicom_userpack_deduct",
				Ctime:    time.Now(),
				Content: map[string]interface{}{
					"phone":     v.Phone,
					"pack_desc": v.UserDesc,
					"integral":  (0 - v.Integral),
				},
			})
		}
	}
}
