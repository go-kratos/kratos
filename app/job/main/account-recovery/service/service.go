package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"go-common/app/job/main/account-recovery/conf"
	"go-common/app/job/main/account-recovery/dao"
	"go-common/app/job/main/account-recovery/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	_recoveryLog = "account_recovery_info"
	//_retry       = 10
	_retrySleep = time.Second * 1
)

// Service struct
type Service struct {
	c               *conf.Config
	dao             *dao.Dao
	compareDatabus  *databus.Databus
	sendMailDatabus *databus.Databus
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:               c,
		dao:             dao.New(c),
		compareDatabus:  databus.New(c.DataBus.CompareDatabus),
		sendMailDatabus: databus.New(c.DataBus.SendMailDatabus),
	}
	go s.compareConsumeproc()
	go s.sendMailConsumeproc()
	return
}

func (s *Service) compareConsumeproc() {
	var (
		msg *databus.Message
		err error
		ok  bool
	)
	for {
		if msg, ok = <-s.compareDatabus.Messages(); !ok {
			log.Error("s.compareDatabus.Message err(%v)", err)
			return
		}
		log.Info("receive msg (%v)", msg.Key)
		mu := &model.Message{}
		if err = json.Unmarshal(msg.Value, &mu); err != nil {
			log.Error("s.compareDatabus.Message err(%v)", err)
			continue
		}
		for {
			switch {
			case strings.HasPrefix(mu.Table, _recoveryLog):
				if mu.Action == "insert" {
					log.Info("begin compare (%v)", msg.Key)
					err = s.compare(context.TODO(), mu.New)
					log.Info("end compare (%v)", msg.Key)
				}
			}
			log.Info("compare switch case (%v), (%v), (%v), (%v), (%v)", strings.HasPrefix(mu.Table, _recoveryLog), mu.Action == "insert", mu.Table, msg.Key, mu.Action)
			if err != nil {
				log.Error("s.flush error(%v)", err)
				time.Sleep(_retrySleep)
				continue
			}
			break
		}
		log.Info("subproc key:%v,topic: %v, part:%v offset:%v,message %s,", msg.Key, msg.Topic, msg.Partition, msg.Offset, msg.Value)
		err = msg.Commit()
		if err != nil {
			log.Error("commit err: %v", err)
		}
	}
}

func (s *Service) sendMailConsumeproc() {
	var (
		msg *databus.Message
		err error
		ok  bool
	)
	for {
		if msg, ok = <-s.sendMailDatabus.Messages(); !ok {
			log.Error("s.sendMailDatabus.Message err(%v)", err)
			return
		}
		log.Info("receive msg (%v)", msg.Key)
		mu := &model.Message{}
		if err = json.Unmarshal(msg.Value, &mu); err != nil {
			log.Error("s.sendMailDatabus.Message err(%v)", err)
			continue
		}
		for {
			switch {
			case strings.HasPrefix(mu.Table, _recoveryLog):
				if mu.Action == "update" {
					log.Info("begin sendMail (%v)", msg.Key)
					err = s.sendMail(context.TODO(), mu.New, mu.Old)
				}
			}
			log.Info("sendMail switch case (%v), (%v), (%v), (%v), (%v)", strings.HasPrefix(mu.Table, _recoveryLog), mu.Action == "update", mu.Table, msg.Key, mu.Action)
			if err != nil {
				log.Error("s.flush error(%v)", err)
				time.Sleep(_retrySleep)
				continue
			}
			break
		}
		log.Info("subproc key:%v,topic: %v, part:%v offset:%v,message %s,", msg.Key, msg.Topic, msg.Partition, msg.Offset, msg.Value)
		err = msg.Commit()
		if err != nil {
			log.Error("commit err: %v", err)
		}
	}
}

// compare compare
func (s *Service) compare(c context.Context, msg []byte) (err error) {
	r := &model.RecoveryInfo{}
	if err = json.Unmarshal(msg, r); err != nil {
		log.Error("s.compare err(%v)", err)
		return
	}
	err = s.dao.CompareInfo(c, r.Rid)
	return
}

// sendMail send mail
func (s *Service) sendMail(c context.Context, new []byte, old []byte) (err error) {
	oldInfo := &model.RecoveryInfo{}
	newInfo := &model.RecoveryInfo{}
	if err = json.Unmarshal(old, oldInfo); err != nil {
		log.Error("failed to oldInfo unmarshal err(%v)", err)
		return
	}
	if err = json.Unmarshal(new, newInfo); err != nil {
		log.Error("failed to newInfo unmarshal err(%v)", err)
		return
	}
	if oldInfo.Status == 0 && (newInfo.Status == 1 || newInfo.Status == 2) {
		err = s.dao.SendMail(c, newInfo.Rid, newInfo.Status)
	}
	return
}

// Close Service
func (s *Service) Close() {
	s.compareDatabus.Close()
	s.sendMailDatabus.Close()
}
