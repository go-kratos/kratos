package service

import (
	"context"
	"encoding/json"
	"net"
	"strings"
	"time"

	location "go-common/app/service/main/location/rpc/client"
	"go-common/app/service/main/secure/conf"
	"go-common/app/service/main/secure/dao"
	"go-common/app/service/main/secure/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	_loginLog   = "aso_login_log"
	_pwdLog     = "aso_pwd_log"
	_retry      = 3
	_retrySleep = time.Second * 1
)

// Service struct of service.
type Service struct {
	c      *conf.Config
	dao    *dao.Dao
	ds     *databus.Databus
	locRPC *location.Service
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:      c,
		dao:    dao.New(c),
		ds:     databus.New(c.DataBus),
		locRPC: location.New(c.LocationRPC),
	}
	go s.subproc()
	//	go s.tableproc()
	return
}

// func (s *Service) tableproc() {
// 	for {
// 		now := time.Now()
// 		ts := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 0, 0, time.Local).Sub(time.Now())
// 		time.Sleep(ts)
// 		err := s.dao.CreateTable(context.TODO(), now.AddDate(0, 0, 1))
// 		fmt.Println("create table", err)
// 		log.Info("createtable err", err)
// 		err = s.dao.DropTable(context.TODO(), now.AddDate(0, 0, -8))
// 		log.Info("drop table err(%v)", err)
// 		time.Sleep(time.Hour)
// 	}
// }

func (s *Service) subproc() {
	var err error
	for {
		res, ok := <-s.ds.Messages()
		if !ok {
			return
		}
		mu := &model.Message{}
		if err = json.Unmarshal(res.Value, mu); err != nil {
			log.Error("credit-job,json.Unmarshal (%v) error(%v)", string(res.Value), err)
			continue
		}
		for i := 0; i < _retry; i++ {
			switch {
			case strings.HasPrefix(mu.Table, _loginLog):
				err = s.loginLog(context.TODO(), mu.Action, mu.New)
			case strings.HasPrefix(mu.Table, _pwdLog):
				if mu.Action == "insert" {
					err = s.changePWDRecord(context.TODO(), mu.New)
				}
			}
			if err != nil {
				log.Error("s.flush error(%v)", err)
				time.Sleep(_retrySleep)
				continue
			}
			break
		}
		log.Info("subproc key:%v,topic: %v, part:%v offset:%v,message %s,", res.Key, res.Topic, res.Partition, res.Offset, res.Value)
		res.Commit()
	}
}

// Close kafka consumer close.
func (s *Service) Close() (err error) {
	return s.ds.Close()
}

// Ping check service health.
func (s *Service) Ping(c context.Context) (err error) {
	return
}

func inetNtoA(sum uint32) string {
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte((sum >> 24) & 0xFF)
	ip[1] = byte((sum >> 16) & 0xFF)
	ip[2] = byte((sum >> 8) & 0xFF)
	ip[3] = byte(sum & 0xFF)
	return ip.String()
}

func inetAtoN(s string) (sum uint32) {
	ip := net.ParseIP(s)
	if ip == nil {
		return
	}
	ip = ip.To4()
	if ip == nil {
		return
	}
	sum += uint32(ip[0]) << 24
	sum += uint32(ip[1]) << 16
	sum += uint32(ip[2]) << 8
	sum += uint32(ip[3])
	return sum
}
