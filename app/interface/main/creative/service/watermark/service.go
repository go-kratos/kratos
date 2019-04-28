package watermark

import (
	"context"
	"time"

	"sync"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/account"
	"go-common/app/interface/main/creative/dao/bfs"
	"go-common/app/interface/main/creative/dao/drawimg"
	"go-common/app/interface/main/creative/dao/monitor"
	"go-common/app/interface/main/creative/dao/watermark"
	wmMDL "go-common/app/interface/main/creative/model/watermark"

	"go-common/app/interface/main/creative/service"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

//Service struct
type Service struct {
	c       *conf.Config
	wm      *watermark.Dao
	drawimg *drawimg.Dao
	acc     *account.Dao
	bfs     *bfs.Dao
	// wait group
	wg sync.WaitGroup
	// databus sub
	userInfoSub *databus.Databus
	// monitor
	monitor    *monitor.Dao
	userInfoMo int64
	// closed
	closed bool
	//async set watermark
	wmChan chan *wmMDL.WatermarkParam
	//task
	p *service.Public
}

//New get service
func New(c *conf.Config, rpcdaos *service.RPCDaos, p *service.Public) *Service {
	s := &Service{
		c:           c,
		wm:          watermark.New(c),
		drawimg:     drawimg.New(c),
		acc:         rpcdaos.Acc,
		bfs:         bfs.New(c),
		monitor:     monitor.New(c),
		userInfoSub: databus.New(c.UserInfoSub),
		wmChan:      make(chan *wmMDL.WatermarkParam, 1024),
		p:           p,
	}
	if c.WaterMark.Consume {
		s.wg.Add(1)
		go s.userInfoConsumer()
		go s.monitorConsume()
	}
	go s.asyncWmSetProc()
	return s
}

// Ping service
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.wm.Ping(c); err != nil {
		log.Error("s.watermark.Dao.PingDb err(%v)", err)
	}
	return
}

func (s *Service) monitorConsume() {
	var userinfo int64
	for {
		time.Sleep(1 * time.Minute)
		if s.userInfoMo-userinfo == 0 {
			s.monitor.Send(context.TODO(), "creative userinfo did not consume within a minute")
		}
		userinfo = s.userInfoMo
	}
}

// Close dao
func (s *Service) Close() {
	s.userInfoSub.Close()
	s.closed = true
	s.wm.Close()
	s.wg.Wait()
}
