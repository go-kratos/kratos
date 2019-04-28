package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"go-common/app/interface/main/account/service/realname/crypto"
	"go-common/app/job/main/member/conf"
	"go-common/app/job/main/member/dao"
	"go-common/app/job/main/member/model"
	"go-common/app/job/main/member/model/queue"
	"go-common/app/job/main/member/service/block"
	memrpc "go-common/app/service/main/member/api/gorpc"
	"go-common/library/log"
	"go-common/library/log/infoc"
	"go-common/library/queue/databus"
	"go-common/library/queue/databus/databusutil"

	"golang.org/x/time/rate"
)

// Service struct of service.
type Service struct {
	c                   *conf.Config
	dao                 *dao.Dao
	block               *block.Service
	ds                  *databus.Databus
	accDs               *databus.Databus
	passortDs           *databus.Databus
	logDatabus          *databus.Databus
	expDatabus          *databus.Databus
	realnameDatabus     *databus.Databus
	shareMidDatabus     *databus.Databus
	loginGroup          *databusutil.Group
	awardGroup          *databusutil.Group
	memrpc              *memrpc.Service
	cachepq             *queue.PriorityQueue
	alipayCryptor       *crypto.Alipay
	limiter             *limiter
	ParsedRealnameInfoc *infoc.Infoc
}

type limiter struct {
	UpdateExp *rate.Limiter
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:               c,
		ds:              databus.New(c.DataBus),
		accDs:           databus.New(c.AccDataBus),
		passortDs:       databus.New(c.PassortDataBus),
		logDatabus:      databus.New(c.LogDatabus),
		expDatabus:      databus.New(c.ExpDatabus),
		realnameDatabus: databus.New(c.RealnameDatabus),
		shareMidDatabus: databus.New(c.ShareMidDatabus),
		alipayCryptor:   crypto.NewAlipay(string(c.RealnameAlipayPub), string(c.RealnameAlipayBiliPriv)),
		loginGroup:      databusutil.NewGroup(c.Databusutil, databus.New(c.LoginDatabus).Messages()),
		awardGroup:      databusutil.NewGroup(c.Databusutil, databus.New(c.AwardDatabus).Messages()),
		memrpc:          memrpc.New(nil),
		cachepq:         queue.NewPriorityQueue(1024, false),
		limiter: &limiter{
			UpdateExp: rate.NewLimiter(200, 10),
		},
		ParsedRealnameInfoc: infoc.New(c.ParsedRealnameInfoc),
	}
	s.dao = dao.New(c)
	s.block = block.New(c, s.dao.BlockImpl(), databus.New(c.BlockCreditDatabus))
	s.loginGroup.New = newMsg
	s.loginGroup.Split = split
	s.loginGroup.Do = s.awardDo
	s.awardGroup.New = newMsg
	s.awardGroup.Split = split
	s.awardGroup.Do = s.awardDo
	s.loginGroup.Start()
	s.awardGroup.Start()
	go s.passportSubproc()
	go s.realnameSubproc()
	go s.realnamealipaycheckproc()
	go s.cachedelayproc(context.Background())
	go s.shareMidproc()
	accproc := int32(10)
	expproc := int32(1)
	if c.Biz.AccprocCount > accproc {
		accproc = c.Biz.AccprocCount
	}
	if c.Biz.ExpprocCount > expproc {
		expproc = c.Biz.ExpprocCount
	}
	log.Info("Starting %d account sub proc", accproc)
	for i := 0; i < int(accproc); i++ {
		go s.subproc()
		go s.accSubproc()
		go s.logproc()
	}
	log.Info("Starting %d exp sub proc", expproc)
	for i := 0; i < int(expproc); i++ {
		go s.expproc()
	}
	if s.c.FeatureGates.DataFixer && s.dao.LeaderEleciton(context.Background()) {
		fmt.Println("Leader elected")
		// 数据检查
		s.makeChan(30)
		go s.dataCheckMids()
		for i := 0; i < 60; i++ {
			go s.dataFixer(csclice[i%30])
		}
	}
	return
}

func split(msg *databus.Message, data interface{}) int {
	t, ok := data.(*model.LoginLogIPString)
	if !ok {
		return 0
	}
	return int(t.Mid)
}

func newMsg(msg *databus.Message) (res interface{}, err error) {
	llm := new(model.LoginLogIPString)
	ll := new(model.LoginLog)
	if err = json.Unmarshal(msg.Value, &llm); err != nil {
		if err = json.Unmarshal(msg.Value, &ll); err != nil {
			msg.Commit()
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			return
		}
		llm.Mid = ll.Mid
		llm.Timestamp = ll.Timestamp
		llm.Loginip = inetNtoA(uint32(ll.Loginip))
	}
	res = llm
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

// Ping check service health.
func (s *Service) Ping(c context.Context) error {
	return s.dao.Ping(c)
}

// Close kafka consumer close.
func (s *Service) Close() (err error) {
	if err = s.ds.Close(); err != nil {
		log.Error("s.ds.Close(),err(%v)", err)
	}
	if err = s.accDs.Close(); err != nil {
		log.Error("s.accDs.Close(),err(%v)", err)
	}
	if err = s.passortDs.Close(); err != nil {
		log.Error("s.passportDs.Close(),err(%v)", err)
	}
	s.block.Close()
	return
}
