package unicom

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/app-wall/model/unicom"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

const (
	_initIPUnicomKey = "ipunicom_%v_%v"
)

func (s *Service) clickConsumer() {
	defer s.waiter.Done()
	msgs := s.clickSub.Messages()
	for {
		msg, ok := <-msgs
		if !ok || s.closed {
			log.Info("s.clickSub.Cloesd")
			return
		}
		msg.Commit()
		var (
			sbs [][]byte
			err error
		)
		if err = json.Unmarshal(msg.Value, &sbs); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", msg.Value, err)
			continue
		}
		for _, bs := range sbs {
			var (
				err   error
				click *unicom.ClickMsg
			)
			if click, err = s.checkMsgIllegal(bs); err != nil {
				log.Error("s.checkMsgIllegal(%s) error(%v)", strings.Replace(string(bs), "\001", "|", -1), err)
				continue
			}
			log.Info("clickConsumer s.checkMsgIllegal(%s)", strings.Replace(string(bs), "\001", "|", -1))
			s.cliChan[click.AID%s.c.ChanNum] <- click
		}
	}
}

func (s *Service) checkMsgIllegal(msg []byte) (click *unicom.ClickMsg, err error) {
	var (
		aid        int64
		clickMsg   []string
		plat       int64
		bvID       string
		mid        int64
		lv         int64
		ctime      int64
		stime      int64
		epid       int64
		ip         string
		seasonType int
		userAgent  string
	)
	clickMsg = strings.Split(string(msg), "\001")
	if len(clickMsg) < 10 {
		err = errors.New("click msg error")
		return
	}
	if aid, err = strconv.ParseInt(clickMsg[1], 10, 64); err != nil {
		err = fmt.Errorf("aid(%s) error", clickMsg[1])
		return
	}
	if aid <= 0 {
		err = fmt.Errorf("wocao aid(%s) error", clickMsg[1])
		return
	}
	if plat, err = strconv.ParseInt(clickMsg[0], 10, 64); err != nil {
		err = fmt.Errorf("plat(%s) error", clickMsg[0])
		return
	}
	if plat != 3 && plat != 4 {
		err = fmt.Errorf("plat(%d) is not android or ios", plat)
		return
	}
	userAgent = clickMsg[10]
	bvID = clickMsg[8]
	if bvID == "" {
		err = fmt.Errorf("bvID(%s) is illegal", clickMsg[8])
		return
	}
	if clickMsg[4] != "" && clickMsg[4] != "0" {
		if mid, err = strconv.ParseInt(clickMsg[4], 10, 64); err != nil {
			err = fmt.Errorf("mid(%s) is illegal", clickMsg[4])
			return
		}
	}
	if clickMsg[5] != "" {
		if lv, err = strconv.ParseInt(clickMsg[5], 10, 64); err != nil {
			err = fmt.Errorf("lv(%s) is illegal", clickMsg[5])
			return
		}
	}
	if ctime, err = strconv.ParseInt(clickMsg[6], 10, 64); err != nil {
		err = fmt.Errorf("ctime(%s) is illegal", clickMsg[6])
		return
	}
	if stime, err = strconv.ParseInt(clickMsg[7], 10, 64); err != nil {
		err = fmt.Errorf("stime(%s) is illegal", clickMsg[7])
		return
	}
	if ip = clickMsg[9]; ip == "" {
		err = errors.New("ip is illegal")
		return
	}
	if clickMsg[17] != "" {
		if epid, err = strconv.ParseInt(clickMsg[17], 10, 64); err != nil {
			err = fmt.Errorf("epid(%s) is illegal", clickMsg[17])
			return
		}
		if clickMsg[15] != "null" {
			if seasonType, err = strconv.Atoi(clickMsg[15]); err != nil {
				err = fmt.Errorf("seasonType(%s) is illegal", clickMsg[15])
				return
			}
		}
	}
	click = &unicom.ClickMsg{
		Plat:       int8(plat),
		AID:        aid,
		MID:        mid,
		Lv:         int8(lv),
		CTime:      ctime,
		STime:      stime,
		BvID:       bvID,
		IP:         ip,
		KafkaBs:    msg,
		EpID:       epid,
		SeasonType: seasonType,
		UserAgent:  userAgent,
	}
	return
}

func (s *Service) cliChanProc(i int64) {
	defer s.waiter.Done()
	var (
		cli     *unicom.ClickMsg
		cliChan = s.cliChan[i]
	)
	for {
		var (
			ub       *unicom.UserBind
			c        = context.TODO()
			err      error
			ok       bool
			count    int
			addFlow  int
			now      = time.Now()
			u        *unicom.Unicom
			cardType string
		)
		if cli, ok = <-cliChan; !ok || s.closed {
			return
		}
		if count, err = s.dao.UserPackReceiveCache(c, cli.MID); err != nil {
			log.Error("s.dao.UserBindCache error(%v) mid(%v) count(%v)", err, cli.MID, count)
			continue
		}
		if count > 0 {
			log.Info("s.dao.UserBindCache mid(%v) count(%v)", cli.MID, count)
			continue
		}
		if ub, err = s.dao.UserBindCache(c, cli.MID); err != nil {
			continue
		}
		if ub == nil || ub.Phone == 0 {
			continue
		}
		res := s.unicomInfo(c, ub.Usermob, now)
		if u, ok = res[ub.Usermob]; !ok || u == nil {
			continue
		}
		switch u.Spid {
		case 10019:
			cardType = "22卡"
		case 10020:
			cardType = "33卡"
		case 10021:
			cardType = "小电视卡"
		default:
			log.Info("unicom spid equal 979 (%v)", ub)
			continue
		}
		ub.Integral = ub.Integral + 10
		switch cli.Lv {
		case 0, 1, 2, 3:
			addFlow = 10
		case 4:
			addFlow = 15
		case 5:
			addFlow = 20
		case 6:
			addFlow = 30
		}
		ub.Flow = ub.Flow + addFlow
		if err = s.dao.AddUserBindCache(c, ub.Mid, ub); err != nil {
			log.Error("s.dao.AddUserBindCache error(%v)", err)
			continue
		}
		if err = s.dao.AddUserPackReceiveCache(c, ub.Mid, 1, now); err != nil {
			log.Error("s.dao.AddUserPackReceiveCache error(%v)", err)
			continue
		}
		s.dbcliChan[ub.Mid%s.c.ChanDBNum] <- ub
		log.Info("unicom mobile cliChanProc userbind(%v)", ub)
		s.unicomInfoc(ub.Usermob, ub.Phone, int(cli.Lv), 10, addFlow, cardType, ub.Mid, now)
		s.addUserIntegralLog(&unicom.UserIntegralLog{Phone: ub.Phone, Mid: ub.Mid, UnicomDesc: cardType, Type: 0, Integral: 10, Flow: addFlow, Desc: "每日礼包"})
	}
}

func (s *Service) dbcliChanProc(i int64) {
	defer s.waiter.Done()
	var (
		ub        *unicom.UserBind
		dbcliChan = s.dbcliChan[i]
	)
	for {
		var (
			c   = context.TODO()
			ok  bool
			row int64
			err error
		)
		if ub, ok = <-dbcliChan; !ok || s.closed {
			return
		}
		if row, err = s.dao.UpUserIntegral(c, ub); err != nil || row == 0 {
			log.Error("s.dao.UpUserIntegral ub(%v) error(%v) or result==0", ub, err)
			continue
		}
		log.Info("unicom mobile dbcliChanProc userbind(%v)", ub)
	}
}

// unicomInfo
func (s *Service) unicomInfo(c context.Context, usermob string, now time.Time) (res map[string]*unicom.Unicom) {
	var (
		err error
		u   []*unicom.Unicom
	)
	res = map[string]*unicom.Unicom{}
	if u, err = s.dao.UnicomCache(c, usermob); err == nil && len(u) > 0 {
		s.pHit.Incr("unicoms_cache")
	} else {
		if u, err = s.dao.OrdersUserFlow(context.TODO(), usermob); err != nil {
			log.Error("unicom_s.dao.OrdersUserFlow error(%v)", err)
			return
		}
		s.pMiss.Incr("unicoms_cache")
	}
	if len(u) > 0 {
		row := &unicom.Unicom{}
		for _, user := range u {
			if user.TypeInt == 1 && now.Unix() <= int64(user.Endtime) {
				*row = *user
				break
			} else if user.TypeInt == 0 {
				if user.Spid == 979 {
					continue
				}
				if int64(row.Ordertime) > int64(user.Ordertime) {
					continue
				}
				*row = *user
			}
		}
		if row.Spid == 0 {
			return
		}
		res[usermob] = row
	}
	return
}

func (s *Service) upBindAll() {
	var (
		orders []*unicom.UserBind
		err    error
		start  = 0
		end    = 1000
	)
	for {
		var tmp []*unicom.UserBind
		if tmp, err = s.dao.BindAll(context.TODO(), start, end); err != nil {
			log.Error("s.dao.BindAll error(%v)", err)
			return
		}
		start = end + start
		if len(tmp) == 0 {
			break
		}
		orders = append(orders, tmp...)
	}
	for _, b := range orders {
		var (
			c        = context.TODO()
			u        *unicom.Unicom
			ok       bool
			now      = time.Now()
			integral int
			ub       *unicom.UserBind
			err      error
			cardType string
		)
		if now.Month() == b.Monthly.Month() && now.Year() == b.Monthly.Year() {
			continue
		}
		res := s.unicomInfo(c, b.Usermob, now)
		if u, ok = res[b.Usermob]; !ok || u == nil {
			continue
		}
		switch u.Spid {
		case 10019:
			integral = 220
			cardType = "22卡"
		case 10020:
			integral = 330
			cardType = "33卡"
		case 10021:
			integral = 660
			cardType = "小电视卡"
		default:
			continue
		}
		if ub, err = s.dao.UserBindCache(c, b.Mid); err != nil {
			continue
		}
		if ub == nil || ub.Phone == 0 {
			continue
		}
		ub.Integral = ub.Integral + integral
		ub.Monthly = now
		if err = s.dao.AddUserBindCache(c, ub.Mid, ub); err != nil {
			log.Error("s.dao.AddUserBindCache error(%v)", err)
			continue
		}
		s.dbcliChan[ub.Mid%s.c.ChanDBNum] <- ub
		log.Info("unicom mobile upBindAll userbind(%v)", ub)
		s.unicomInfoc(ub.Usermob, ub.Phone, 0, integral, 0, cardType, ub.Mid, now)
		s.addUserIntegralLog(&unicom.UserIntegralLog{Phone: ub.Phone, Mid: ub.Mid, UnicomDesc: cardType, Type: 1, Integral: integral, Flow: 0, Desc: "每月礼包"})
	}
}

func (s *Service) updatemonth(now time.Time) {
	m := int(now.Month())
	if lmonth, ok := s.lastmonth[m]; !ok || !lmonth {
		if now.Day() == 1 {
			s.upBindAll()
			s.lastmonth[m] = true
			if m = m + 1; m > 12 {
				m = 1
			}
			s.lastmonth[m] = false
			log.Info("updatepro user monthly integral success")
		}
	}
}

func (s *Service) loadUnicomFlow() {
	var (
		list map[string]*unicom.UnicomUserFlow
		err  error
	)
	if list, err = s.dao.UserFlowListCache(context.TODO()); err != nil {
		log.Error("load unicom s.dao.UserFlowListCache error(%v)", err)
		return
	}
	log.Info("load unicom flow total len(%v)", len(list))
	for key, u := range list {
		var (
			c           = context.TODO()
			requestNo   int64
			orderstatus string
			msg         string
		)
		if err = s.dao.UserFlowCache(c, key); err != nil {
			if err == memcache.ErrNotFound {
				if err = s.returnPoints(c, u); err != nil {
					if err != ecode.NothingFound {
						log.Error("load unicom s.returnPoints error(%v)", err)
						continue
					}
					err = nil
				}
				log.Info("load unicom userbind timeout flow(%v)", u)
			} else {
				log.Error("load unicom s.dao.UserFlowCache error(%v)", err)
				continue
			}
		} else {
			if requestNo, err = s.seqdao.SeqID(c); err != nil {
				log.Error("load unicom s.seqdao.SeqID error(%v)", err)
				continue
			}
			if orderstatus, msg, err = s.dao.FlowQry(c, u.Phone, requestNo, u.Outorderid, u.Orderid, time.Now()); err != nil {
				log.Error("load unicom s.dao.FlowQry error(%v) msg(%s)", err, msg)
				continue
			}
			log.Info("load unicom userbind flow(%v) orderstatus(%s)", u, orderstatus)
			if orderstatus == "00" {
				continue
			} else if orderstatus != "01" {
				if err = s.returnPoints(c, u); err != nil {
					if err != ecode.NothingFound {
						log.Error("load unicom s.returnPoints error(%v)", err)
						continue
					}
					err = nil
				}
			}
		}
		delete(list, key)
		if err = s.dao.DeleteUserFlowCache(c, key); err != nil {
			log.Error("load unicom s.dao.DeleteUserFlowCache error(%v)", err)
			continue
		}
	}
	if err = s.dao.AddUserFlowListCache(context.TODO(), list); err != nil {
		log.Error("load unicom s.dao.AddUserFlowListCache error(%v)", err)
		return
	}
	log.Info("load unicom flow last len(%v) success", len(list))
}

// returnPoints retutn user integral and flow
func (s *Service) returnPoints(c context.Context, u *unicom.UnicomUserFlow) (err error) {
	var (
		userbind *unicom.UserBind
		result   int64
	)
	if userbind, err = s.unicomBindInfo(c, u.Mid); err != nil {
		return
	}
	ub := &unicom.UserBind{}
	*ub = *userbind
	ub.Flow = ub.Flow + u.Flow
	ub.Integral = ub.Integral + u.Integral
	if err = s.dao.AddUserBindCache(c, ub.Mid, ub); err != nil {
		log.Error("unicom s.dao.AddUserBindCache error(%v)", err)
		return
	}
	if result, err = s.dao.UpUserIntegral(c, ub); err != nil || result == 0 {
		log.Error("unicom s.dao.UpUserIntegral error(%v) or result==0", err)
		return
	}
	var packInt int
	if u.Integral > 0 {
		packInt = u.Integral
	} else {
		packInt = u.Flow
	}
	ul := &unicom.UserPackLog{
		Phone:     u.Phone,
		Usermob:   ub.Usermob,
		Mid:       u.Mid,
		RequestNo: u.Outorderid,
		Type:      0,
		Desc:      u.Desc + ",领取失败并返还",
		Integral:  packInt,
	}
	s.addUserPackLog(ul)
	s.addUserIntegralLog(&unicom.UserIntegralLog{Phone: u.Phone, Mid: u.Mid, UnicomDesc: "", Type: 2, Integral: u.Integral, Flow: u.Flow, Desc: u.Desc + ",领取失败并返还"})
	log.Info("unicom_pack(%v) mid(%v)", u.Desc+",领取失败并返还", userbind.Mid)
	s.unicomPackInfoc(userbind.Usermob, u.Desc+",领取失败并返还", u.Orderid, userbind.Phone, packInt, 0, userbind.Mid, time.Now())
	return
}

// unicomBindInfo unicom bind info
func (s *Service) unicomBindInfo(c context.Context, mid int64) (res *unicom.UserBind, err error) {
	if res, err = s.dao.UserBindCache(c, mid); err != nil {
		if res, err = s.dao.UserBind(c, mid); err != nil {
			log.Error("s.dao.UserBind error(%v)", err)
			return
		}
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

// loadUnicomIPOrder load unciom ip order update
func (s *Service) loadUnicomIPOrder(now time.Time) {
	var (
		dbips map[string]*unicom.UnicomIP
		err   error
	)
	if dbips, err = s.loadUnicomIP(context.TODO()); err != nil {
		log.Error("s.loadUnicomIP", err)
		return
	}
	if len(dbips) == 0 {
		log.Error("db cache ip len 0")
		return
	}
	unicomIP, err := s.dao.UnicomIP(context.TODO(), now)
	if err != nil {
		log.Error("s.dao.UnicomIP(%v)", err)
		return
	}
	if len(unicomIP) == 0 {
		log.Info("unicom ip orders is null")
		return
	}
	tx, err := s.dao.BeginTran(context.TODO())
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}
	for _, uip := range unicomIP {
		key := fmt.Sprintf(_initIPUnicomKey, uip.Ipbegin, uip.Ipend)
		if _, ok := dbips[key]; ok {
			delete(dbips, key)
			continue
		}
		var (
			result int64
		)
		if result, err = s.dao.InUnicomIPSync(tx, uip, time.Now()); err != nil || result == 0 {
			tx.Rollback()
			log.Error("s.dao.InUnicomIPSync error(%v)", err)
			return
		}
	}
	for _, uold := range dbips {
		var (
			result int64
		)
		if result, err = s.dao.UpUnicomIP(tx, uold.Ipbegin, uold.Ipend, 0, time.Now()); err != nil || result == 0 {
			tx.Rollback()
			log.Error("s.dao.UpUnicomIP error(%v)", err)
			return
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	log.Info("update unicom ip success")
}

// loadUnicomIP load unicom ip
func (s *Service) loadUnicomIP(c context.Context) (res map[string]*unicom.UnicomIP, err error) {
	var unicomIP []*unicom.UnicomIP
	unicomIP, err = s.dao.IPSync(c)
	if err != nil {
		log.Error("s.dao.IPSync error(%v)", err)
		return
	}
	tmp := map[string]*unicom.UnicomIP{}
	for _, u := range unicomIP {
		key := fmt.Sprintf(_initIPUnicomKey, u.Ipbegin, u.Ipend)
		tmp[key] = u
	}
	res = tmp
	log.Info("loadUnicomIPCache success")
	return
}

func (s *Service) addUserPackLog(u *unicom.UserPackLog) {
	select {
	case s.packLogCh <- u:
	default:
		log.Warn("user pack log buffer is full")
	}
}

func (s *Service) addUserIntegralLog(u *unicom.UserIntegralLog) {
	select {
	case s.integralLogCh[u.Mid%s.c.ChanDBNum] <- u:
	default:
		log.Warn("user add integral and flow log buffer is full")
	}
}

func (s *Service) addUserPackLogproc() {
	for {
		i, ok := <-s.packLogCh
		if !ok || s.closed {
			log.Warn("user pack log proc exit")
			return
		}
		var (
			c      = context.TODO()
			result int64
			err    error
		)
		switch v := i.(type) {
		case *unicom.UserPackLog:
			if result, err = s.dao.InUserPackLog(c, v); err != nil || result == 0 {
				log.Error("s.dao.UpUserIntegral error(%v) or result==0", err)
				continue
			}
			log.Info("unicom user flow or integral back mid(%d) phone(%d)", v.Mid, v.Phone)
		}
	}
}

func (s *Service) addUserIntegralLogproc(i int64) {
	var (
		dbcliChan = s.integralLogCh[i]
	)
	for {
		i, ok := <-dbcliChan
		if !ok || s.closed {
			log.Warn("user pack log proc exit")
			return
		}
		var (
			logID = 91
		)
		switch v := i.(type) {
		case *unicom.UserIntegralLog:
			// if result, err = s.dao.InUserIntegralLog(c, v); err != nil || result == 0 {
			// 	log.Error("s.dao.InUserIntegralLog error(%v) or result==0", err)
			// 	continue
			// }
			report.User(&report.UserInfo{
				Mid:      v.Mid,
				Business: logID,
				Action:   "unicom_userpack_add",
				Ctime:    time.Now(),
				Content: map[string]interface{}{
					"phone":     v.Phone,
					"pack_desc": v.Desc,
					"integral":  v.Integral,
				},
			})
		}
	}
}
