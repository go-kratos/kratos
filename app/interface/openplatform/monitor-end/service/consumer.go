package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go-common/app/interface/openplatform/monitor-end/conf"
	"go-common/app/interface/openplatform/monitor-end/model/kafka"
	"go-common/library/log"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

const _group = "open-monitor-end"

// Consumer .
type Consumer struct {
	// c      sarama.Consumer
	c        *cluster.Consumer
	Config   *kafka.Config
	closed   bool
	paused   bool
	duration time.Duration
	messages chan *sarama.ConsumerMessage
}

var (
	errClosedMsgChannel    = errors.New("message channel is closed")
	errClosedNotifyChannel = errors.New("notification channel is closed")
	errConsumerOver        = errors.New("too many consumers")
)

// NewConsumer .
func NewConsumer(c *conf.Config) (con *Consumer, err error) {
	con = &Consumer{
		Config:   c.Kafka,
		messages: make(chan *sarama.ConsumerMessage, 1024),
	}
	cfg := cluster.NewConfig()
	cfg.Version = sarama.V0_8_2_0
	cfg.ClientID = fmt.Sprintf("%s-%s", _group, c.Kafka.Topic)
	cfg.Net.KeepAlive = 30 * time.Second
	// NOTE cluster auto commit offset interval
	cfg.Consumer.Offsets.CommitInterval = time.Second * 1
	// NOTE set fetch.wait.max.ms
	cfg.Consumer.MaxWaitTime = time.Millisecond * 250
	cfg.Consumer.MaxProcessingTime = 50 * time.Millisecond
	// NOTE errors that occur during offset management,if enabled, c.Errors channel must be read
	cfg.Consumer.Return.Errors = true
	// NOTE notifications that occur during consumer, if enabled, c.Notifications channel must be read
	cfg.Group.Return.Notifications = true
	// con.c = sarama.NewConsumer(c.Kafka.Addr, nil)
	// consumer.Partitions = consumer.c.Partitions(consumer.Config.Topic)
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	if con.c, err = cluster.NewConsumer(c.Kafka.Addr, _group, []string{c.Kafka.Topic}, cfg); err != nil {
		log.Error("s.NewConsumer group(%s) topic(%s) addr(%s) cluster.NewConsumer() error(%v)", _group, c.Kafka.Topic, strings.Join(c.Kafka.Addr, ","), err)
	} else {
		log.Info("s.NewConsumer group(%s) topic(%s) addr(%s) cluster.NewConsumer() ok", _group, c.Kafka.Topic, strings.Join(c.Kafka.Addr, ","))
	}
	return
}

func (s *Service) consume() {
	for {
		if s.consumer.closed {
			return
		}
		if err := s.consumer.message(); err != nil {
			time.Sleep(time.Minute)
		}
	}
}

func (s *Service) handleMsg() {
	for {
		select {
		case msg := <-s.consumer.messages:
			s.HandleMsg(msg.Value)
		}
		if s.consumer.closed {
			return
		}
	}
}

func (s *Consumer) message() (err error) {
	var (
		msg    *sarama.ConsumerMessage
		notify *cluster.Notification
		ok     bool
	)
	for {
		if s.closed {
			s.c.Close()
			err = nil
			return
		}
		if s.paused {
			s.paused = false
			time.Sleep(s.duration)
		}
		select {
		case err = <-s.c.Errors():
			log.Error("group(%s) topic(%s) addr(%s) catch error(%v)", _group, s.Config.Topic, s.Config.Addr, err)
			return
		case notify, ok = <-s.c.Notifications():
			if !ok {
				log.Info("notification notOk group(%s) topic(%s) addr(%v) catch error(%v)", _group, s.Config.Topic, s.Config.Addr, err)
				err = errClosedNotifyChannel
				return
			}
			switch notify.Type {
			case cluster.UnknownNotification, cluster.RebalanceError:
				log.Error("notification(%s) group(%s) topic(%s) addr(%v) catch error(%v)", notify.Type, _group, s.Config.Topic, s.Config.Addr, err)
				err = errClosedNotifyChannel
				return
			case cluster.RebalanceStart:
				log.Info("notification(%s) group(%s) topic(%s) addr(%v) catch error(%v)", notify.Type, _group, s.Config.Topic, s.Config.Addr, err)
				continue
			case cluster.RebalanceOK:
				log.Info("notification(%s) group(%s) topic(%s) addr(%v) catch error(%v)", notify.Type, _group, s.Config.Topic, s.Config.Addr, err)
			}
			if len(notify.Current[s.Config.Topic]) == 0 {
				log.Warn("notification(%s) no topic group(%s) topic(%s) addr(%v) catch error(%v)", notify.Type, _group, s.Config.Topic, s.Config.Addr, err)
				err = errConsumerOver
				return
			}
		case msg, ok = <-s.c.Messages():
			if !ok {
				log.Error("group(%s) topic(%s) addr(%v) message channel closed", _group, s.Config.Topic, s.Config.Addr)
				err = errClosedMsgChannel
				return
			}
			s.messages <- msg
		}
	}
}
