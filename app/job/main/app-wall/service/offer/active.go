package offer

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/job/main/app-wall/model"
	"go-common/app/job/main/app-wall/model/offer"
	"go-common/library/log"
)

func (s *Service) activeConsumer() {
	defer s.waiter.Done()
LOOP:
	for {
		select {
		case err := <-s.consumer.Errors():
			log.Error("group(%s) topic(%s) addr(%s) catch error(%+v)", s.c.Consumer.Group, s.c.Consumer.Topic, s.c.Consumer.Brokers, err)
			continue
		case notify := <-s.consumer.Notifications():
			log.Info("notification(%v)", notify)
			continue
		case msg, ok := <-s.consumer.Messages():
			if !ok {
				log.Error("active consumer exit!")
				break LOOP
			}
			s.consumer.MarkOffset(msg, "")
			active, err := s.checkMsgIllegal(msg.Value)
			if err != nil {
				log.Error("s.checkMsgIllegal(%s) error(%v)", msg.Value, err)
				continue
			}
			if active == nil {
				continue
			}
			s.activeChan <- active
		}
	}
}

func (s *Service) checkMsgIllegal(msg []byte) (active *offer.ActiveMsg, err error) {
	var (
		msgs      []string
		pid       int64
		os        string
		androidid string
		imei      string
	)
	msgs = strings.Split(string(msg), "|")
	if len(msgs) < 9 {
		err = fmt.Errorf("active msg(%s) split len(%d)<9", msg, len(msgs))
		return
	}
	if pid, err = strconv.ParseInt(msgs[8], 10, 64); err != nil {
		return
	}
	if pid%10 == 3 {
		os = model.TypeAndriod
		if len(msgs) > 22 {
			androidid = msgs[22]
		}
		if len(msgs) > 23 {
			imei = msgs[23]
		}
		if imei == "" {
			log.Warn("active msg(%s) imei(%s) is illegal", msg, imei)
		} else {
			log.Warn("active msg(%s) imei(%s) is legal", msg, imei)
		}
		if androidid == "" {
			log.Warn("active msg(%s) androidid(%s) is illegal", msg, androidid)
		}
		if androidid == "" && imei == "" {
			err = fmt.Errorf("active msg(%s) androidid(%s) and imei(%s) is illegal", msg, androidid, imei)
			return
		}
	} else {
		err = fmt.Errorf("active msg(%s) pid(%d) platform not android", msg, pid)
		return
	}
	active = &offer.ActiveMsg{OS: os, IMEI: imei, Androidid: androidid, Mac: ""}
	return
}

func (s *Service) activeproc() {
	defer s.waiter.Done()
	for {
		msg, ok := <-s.activeChan
		if !ok {
			log.Error("active chan id closed")
			break
		}
		s.active(msg)
	}
}

func (s *Service) active(msg *offer.ActiveMsg) {
	var err error
	c := context.TODO()
	if err = retry(func() (err error) {
		return s.dao.Active(c, msg.OS, msg.IMEI, msg.Androidid, msg.Mac, "")
	}, _upActiveRetry, _sleep); err != nil {
		log.Error("%+v", err)
		if err = s.syncRetry(c, offer.ActionActive, msg.OS, msg.IMEI, msg.Androidid, msg.Mac); err != nil {
			log.Error("%+v", err)
		}
		return
	}
	log.Info("active device os(%s) imei(%s) androidid(%s) mac(%s) success", msg.OS, msg.IMEI, msg.Androidid, msg.Mac)
}
