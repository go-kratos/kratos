package service

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"go-common/app/service/main/seq-server/conf"
	"go-common/app/service/main/seq-server/dao"
	"go-common/app/service/main/seq-server/model"
	"go-common/library/log"
	"go-common/library/net/rpc"
	xtime "go-common/library/time"
)

// Service is service.
type Service struct {
	c             *conf.Config
	db            *dao.Dao
	bs            map[int64]*model.Business
	bl            sync.RWMutex
	idc           int64
	idcLen        uint
	svrNum        int64
	lastTimestamp int64
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
		// dao
		db:            dao.New(c),
		bs:            make(map[int64]*model.Business),
		lastTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}
	var (
		idc    string
		svrNum string
		err    error
	)
	idc = os.Getenv("SEQ_IDC")
	if s.idc, err = strconv.ParseInt(idc, 10, 64); err != nil || s.idc > 3 {
		panic(fmt.Sprintf("SEQ_IDC config(%s) error(%v)", idc, err))
	}
	svrNum = os.Getenv("SEQ_SERVERNUM")
	if s.svrNum, err = strconv.ParseInt(svrNum, 10, 64); err != nil || s.svrNum > 15 {
		panic(fmt.Sprintf("SEQ_SERVERNUM config(%s) or error(%v)", svrNum, err))
	}
	s.idcLen = uint(len(strconv.FormatInt(s.idc, 2)))
	svrLen := uint(len(strconv.FormatInt(s.svrNum, 2)))
	if s.idcLen+svrLen > 6 {
		panic("ids + svrNum must <= 6")
	}
	for k, addr := range s.c.SeqSvrs {
		client := rpc.Dial(addr, xtime.Duration(100*time.Millisecond), nil)
		var (
			_noArg = struct{}{}
			res    *model.SeqVersion
			err    error
		)
		if err = client.Call(context.TODO(), "RPC.CheckVersion", _noArg, &res); err != nil {
			log.Error("client.Call(RPC.CheckVersion) addr(%s) error(%v)", addr, err)
			continue
		}
		now := time.Now().Unix()
		if now-res.SvrTime > 2 || res.SvrTime-now < -2 {
			panic(fmt.Sprintf("myTime(%d) svrTime(%d) svrNum(%d),time is running out too much", now, res.SvrTime, res.SvrNum))
		}
		if res.IDC == s.idc && res.SvrNum == s.svrNum {
			panic(fmt.Sprintf("server(%d) idc and svrNum is same", k))
		}
		client.Close()
	}
	s.loadSeqs()
	go s.loadproc()
	return
}

func (s *Service) loadproc() {
	for {
		time.Sleep(time.Duration(s.c.Tick))
		s.loadSeqs()
	}
}

// Close resource.
func (s *Service) Close() {
	s.db.Close()
}

func (s *Service) loadSeqs() (err error) {
	var bs map[int64]*model.Business
	if bs, err = s.db.All(context.TODO()); err != nil {
		return
	}
	for id, b := range bs {
		s.bl.Lock()
		if _, ok := s.bs[id]; !ok {
			s.bs[id] = b
			s.bl.Unlock()
			continue
		}
		s.bl.Unlock()
	}
	return
}

// Ping ping service is ok.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.db.Ping(); err != nil {
		return
	}
	if len(s.bs) == 0 {
		err = model.ErrBusinessNotReady
	}
	return
}
