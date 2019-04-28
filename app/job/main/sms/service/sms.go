package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/sms/dao"
	"go-common/app/job/main/sms/model"
	smsmdl "go-common/app/service/main/sms/model"
	"go-common/library/conf/env"
	"go-common/library/log"
)

const _retry = 3

var _contentRe = regexp.MustCompile(`\d`)

func (s *Service) subproc() {
	defer s.waiter.Done()
	for {
		item, ok := <-s.databus.Messages()
		if !ok {
			close(s.sms)
			close(s.actSms)
			close(s.batchSms)
			log.Info("databus: sms-job subproc consumer exit!")
			return
		}
		s.smsCount++
		msg := new(smsmdl.ModelSend)
		if err := json.Unmarshal(item.Value, &msg); err != nil {
			log.Error("json.Unmarshal (%v) error(%v)", string(item.Value), err)
			continue
		}
		log.Info("subproc topic(%s) key(%s) partition(%v) offset(%d) message(%s)", item.Topic, item.Key, item.Partition, item.Offset, item.Value)
		// 黑名单，用于压测
		if _, ok := s.blacklist[msg.Country+msg.Mobile]; ok {
			log.Info("country(%s) mobile(%s) in blacklist", msg.Country, msg.Mobile)
			item.Commit()
			continue
		}
		switch msg.Type {
		case smsmdl.TypeSms:
			s.sms <- msg
		case smsmdl.TypeActSms:
			s.actSms <- msg
		case smsmdl.TypeActBatch:
			s.batchSms <- msg
		}
		item.Commit()
	}
}

func (s *Service) smsproc() {
	defer s.waiter.Done()
	var (
		err   error
		msgid string
	)
	for {
		m, ok := <-s.sms
		if !ok {
			log.Info("smsproc exit!")
			return
		}
		if m.Mobile == "" {
			if m.Country, m.Mobile, err = s.userMobile(m.Mid); err != nil {
				continue
			}
		}
		if m.Country == "" || m.Mobile == "" {
			log.Error("invalid country or mobile, info(%+v)", m)
			continue
		}
		content := _contentRe.ReplaceAllString(m.Content, "*")
		l := &smsmdl.ModelUserActionLog{Mobile: m.Mobile, Content: content, Type: smsmdl.TypeSms, Action: smsmdl.UserActionTypeSend}
		if m.Country == smsmdl.CountryChina {
			for i := 0; i < s.providers; i++ {
				s.smsp.Lock()
				p := s.smsp.Value.(model.Provider)
				s.smsp.Ring = s.smsp.Next()
				s.smsp.Unlock()
				l.Provider = p.GetPid()
				if msgid, err = p.SendSms(context.Background(), m); err == nil {
					break
				}
				dao.PromInfo(fmt.Sprintf("service:retry %d", l.Provider))
				log.Error("retry send sms(%v) platform(%d) error(%v)", m, l.Provider, err)
			}
		} else {
			for i := 0; i < s.providers; i++ {
				s.intep.Lock()
				p := s.intep.Value.(model.Provider)
				s.intep.Ring = s.intep.Next()
				s.intep.Unlock()
				l.Provider = p.GetPid()
				if msgid, err = p.SendInternationalSms(context.Background(), m); err == nil {
					break
				}
				dao.PromInfo(fmt.Sprintf("service:retry international %d", l.Provider))
				log.Error("retry send international sms(%v) platform(%d) error(%v)", m, l.Provider, err)
			}
		}
		if err == nil {
			l.Status = smsmdl.UserActionSendSuccessStatus
			l.Desc = smsmdl.UserActionSendSuccessDesc
			dao.PromInfo(fmt.Sprintf("service:success %d", l.Provider))
			log.Info("send sms(%v) platform(%d) success", m, l.Provider)
		} else {
			l.Status = smsmdl.UserActionSendFailedStatus
			l.Desc = smsmdl.UserActionSendFailedDesc
			dao.PromError("service:sms")
			log.Error("send sms(%v) error(%v)", m, err)
			s.cache.Do(context.Background(), func(ctx context.Context) {
				s.dao.SendWechat(fmt.Sprintf("sms-job send msg(%d) error(%v)", m.ID, err))
			})

		}
		l.MsgID = msgid
		l.Ts = time.Now().Unix()
		s.sendUserActionLog(l)
	}
}

func (s *Service) actsmsproc() {
	defer s.waiter.Done()
	var (
		err   error
		msgid string
	)
	for {
		m, ok := <-s.actSms
		if !ok {
			log.Info("actsmsproc exit!")
			return
		}
		if m.Mobile == "" {
			if m.Country, m.Mobile, err = s.userMobile(m.Mid); err != nil {
				continue
			}
		}
		if m.Country == "" || m.Mobile == "" {
			log.Error("invalid country or mobile, info(%+v)", m)
			continue
		}
		content := _contentRe.ReplaceAllString(m.Content, "*")
		l := &smsmdl.ModelUserActionLog{Mobile: m.Mobile, Content: content, Type: smsmdl.TypeActSms, Action: smsmdl.UserActionTypeSend}
		if m.Country == smsmdl.CountryChina {
			for i := 0; i < s.providers; i++ {
				s.actp.Lock()
				p := s.actp.Value.(model.Provider)
				s.actp.Ring = s.actp.Next()
				s.actp.Unlock()
				l.Provider = p.GetPid()
				if msgid, err = p.SendActSms(context.Background(), m); err == nil {
					break
				}
				dao.PromInfo(fmt.Sprintf("service:retry act china %d", l.Provider))
				log.Error("retry send act sms(%v) platform(%d) error(%v)", m, l.Provider, err)
			}
		} else {
			for i := 0; i < s.providers; i++ {
				s.intep.Lock()
				p := s.intep.Value.(model.Provider)
				s.intep.Ring = s.intep.Next()
				s.intep.Unlock()
				l.Provider = p.GetPid()
				if msgid, err = p.SendInternationalSms(context.Background(), m); err == nil {
					break
				}
				dao.PromInfo(fmt.Sprintf("service:retry act international %d", l.Provider))
				log.Error("retry send act international sms(%v) platform(%d)", m, l.Provider, err)
			}
		}
		if err == nil {
			l.Status = smsmdl.UserActionSendSuccessStatus
			l.Desc = smsmdl.UserActionSendSuccessDesc
			dao.PromInfo(fmt.Sprintf("service:act china success %d", l.Provider))
			log.Info("send act sms(%v) platform(%d) success", m, l.Provider)
		} else {
			l.Status = smsmdl.UserActionSendFailedStatus
			l.Desc = smsmdl.UserActionSendFailedDesc
			dao.PromError("service:actSms")
			log.Error("send act sms(%v) error(%v)", m, err)
			s.cache.Do(context.Background(), func(ctx context.Context) {
				s.dao.SendWechat(fmt.Sprintf("sms-job send msg(%d) error(%v)", m.ID, err))
			})
		}
		l.MsgID = msgid
		l.Ts = time.Now().Unix()
		s.sendUserActionLog(l)
	}
}

func (s *Service) actbatchproc() {
	defer s.waiter.Done()
	var (
		err     error
		mids    []string
		country string
		mobile  string
		msgid   string
	)
	for {
		m, ok := <-s.batchSms
		if !ok {
			log.Info("actbatchproc exit!")
			return
		}
		if m.Mobile == "" && m.Mid != "" {
			mids = strings.Split(m.Mid, ",")
			var mobiles []string
			for _, midStr := range mids {
				if country, mobile, err = s.userMobile(midStr); err != nil {
					continue
				}
				if country == "" || mobile == "" {
					log.Error("invalid country or mobile, code(%s) mid(%s) country(%s) mobile(%s)", m.Code, midStr, country, mobile)
					continue
				}
				if country != smsmdl.CountryChina {
					continue
				}
				mobiles = append(mobiles, mobile)
			}
			m.Mobile = strings.Join(mobiles, ",")
		}
		if m.Mobile == "" {
			continue
		}
		content := _contentRe.ReplaceAllString(m.Content, "*")
		l := &smsmdl.ModelUserActionLog{Mobile: m.Mobile, Content: content, Type: smsmdl.TypeActSms, Action: smsmdl.UserActionTypeSend, Ts: time.Now().Unix()}
		send := &smsmdl.ModelSend{Mobile: m.Mobile, Content: m.Content, Type: smsmdl.TypeActSms}
		for i := 0; i < s.providers; i++ {
			s.batchp.Lock()
			p := s.batchp.Value.(model.Provider)
			s.batchp.Ring = s.batchp.Next()
			s.batchp.Unlock()
			l.Provider = p.GetPid()
			if msgid, err = p.SendBatchActSms(context.Background(), send); err == nil {
				break
			}
			dao.PromInfo(fmt.Sprintf("service:retry batch %d", l.Provider))
			log.Error("retry send act batch sms(%v) platform(%d)", m, l.Provider, err)
		}
		if err == nil {
			dao.PromInfo(fmt.Sprintf("service:batch success %d", l.Provider))
			log.Info("send act batch sms(%v) platform(%d) success", m, l.Provider)
			l.Status = smsmdl.UserActionSendSuccessStatus
			l.Desc = smsmdl.UserActionSendSuccessDesc
		} else {
			dao.PromError("service:actBatchSms")
			log.Error("send act batch sms(%v) error(%v)", m, err)
			s.cache.Do(context.Background(), func(ctx context.Context) {
				s.dao.SendWechat(fmt.Sprintf("sms-job send msg(%d) error(%v)", m.ID, err))
			})
			l.Status = smsmdl.UserActionSendFailedStatus
			l.Desc = smsmdl.UserActionSendFailedDesc
		}
		l.MsgID = msgid
		l.Ts = time.Now().Unix()
		s.sendUserActionLog(l)
	}
}

func (s *Service) userMobile(midStr string) (country, mobile string, err error) {
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("userMobile parse mid(%s) error(%v)", midStr, err)
		return
	}
	if mid <= 0 {
		log.Error("userMobile invalid mid(%s)", midStr)
		return
	}
	var um *model.UserMobile
	for i := 0; i < _retry; i++ {
		if um, err = s.dao.UserMobile(context.Background(), mid); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		log.Error("UserMobile mid(%d) error(%v)", mid, err)
		return
	}
	country = um.CountryCode
	mobile = um.Mobile
	return
}

func (s *Service) monitorproc() {
	if env.DeployEnv != env.DeployEnvProd {
		return
	}
	var smsCount int64
	for {
		time.Sleep(time.Duration(s.c.Sms.MonitorProcDuration))
		if s.smsCount-smsCount == 0 {
			msg := fmt.Sprintf("sms-job sms did not consume within %s seconds", time.Duration(s.c.Sms.MonitorProcDuration).String())
			s.dao.SendWechat(msg)
			log.Warn(msg)
		}
		smsCount = s.smsCount
	}
}
