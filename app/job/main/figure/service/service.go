package service

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"go-common/app/job/main/figure/conf"
	"go-common/app/job/main/figure/dao"
	"go-common/app/job/main/figure/model"
	coinm "go-common/app/service/main/coin/model"
	spym "go-common/app/service/main/spy/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	_insertAction = "insert"

	_vipTable     = "vip_user_info"
	_payTable     = "pay_pay_order"
	_blockedTable = "blocked_kpi"

	// pay success status.
	_paySuccess = 2
)

// Service biz service def.
type Service struct {
	c                *conf.Config
	figureDao        *dao.Dao
	accExpDatabus    *databus.Databus
	accRegDatabus    *databus.Databus
	vipDatabus       *databus.Databus
	spyDatabus       *databus.Databus
	coinDatabus      *databus.Databus
	replyInfoDatabus *databus.Databus
	payDatabus       *databus.Databus
	danmakuDatabus   *databus.Databus
	blockedDatabus   *databus.Databus
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:         c,
		figureDao: dao.New(c),
	}
	if c.DataSource.AccountExp != nil {
		s.accExpDatabus = databus.New(c.DataSource.AccountExp)
		go s.accexpproc()
	}
	if c.DataSource.AccountReg != nil {
		s.accRegDatabus = databus.New(c.DataSource.AccountReg)
		go s.accregproc()
	}
	if c.DataSource.Vip != nil {
		s.vipDatabus = databus.New(c.DataSource.Vip)
		go s.vipproc()
	}
	if c.DataSource.Spy != nil {
		s.spyDatabus = databus.New(c.DataSource.Spy)
		go s.spyproc()
	}
	if c.DataSource.Coin != nil {
		s.coinDatabus = databus.New(c.DataSource.Coin)
		go s.coinproc()
	}
	if c.DataSource.ReplyInfo != nil {
		s.replyInfoDatabus = databus.New(c.DataSource.ReplyInfo)
		go s.replyinfoproc()
	}
	if c.DataSource.Pay != nil {
		s.payDatabus = databus.New(c.DataSource.Pay)
		go s.payproc()
	}
	if c.DataSource.Danmaku != nil {
		s.danmakuDatabus = databus.New(c.DataSource.Danmaku)
		go s.danmakuproc()
	}
	if c.DataSource.Blocked != nil {
		s.blockedDatabus = databus.New(c.DataSource.Blocked)
		go s.blockedproc()
	}
	go s.syncproc()
	return s
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	s.figureDao.Ping(c)
	return
}

// Close close all dao.
func (s *Service) Close() {
	s.figureDao.Close()
}

func (s *Service) accexpproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("s.accexpproc panic(%v) %s", x, debug.Stack())
			go s.accexpproc()
			log.Info("s.accexpproc recover")
		}
	}()
	var (
		err     error
		msg     *databus.Message
		msgChan = s.accExpDatabus.Messages()
		ok      bool
		exp     *model.MsgAccountLog
	)
	for {
		select {
		case msg, ok = <-msgChan:
			if !ok {
				log.Info("accproc msgChan closed")
			}
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit err(%v)", err)
		}
		exp = &model.MsgAccountLog{}
		if err = json.Unmarshal([]byte(msg.Value), exp); err != nil {
			log.Error("json.Unmarshal(%v) err(%v)", msg, err)
			continue
		}
		log.Info("exp Info (%+v)", exp)
		if err = s.AccountExp(context.Background(), exp.Mid, int64(exp.ExpTo())); err != nil {
			log.Error("s.AccountExp(%v) err(%v)", exp, err)
			continue
		}
		s.figureDao.SetWaiteUserCache(context.Background(), exp.Mid, s.figureDao.Version(time.Now()))
		if exp.IsViewExp() {
			if err = s.AccountViewVideo(context.Background(), exp.Mid); err != nil {
				log.Error("s.AccountExp(%v) err(%v)", exp, err)
				continue
			}
		}
	}
}

func (s *Service) accregproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("s.accregproc panic(%v) %s", x, debug.Stack())
			go s.accregproc()
			log.Info("s.accregproc recover")
		}
	}()
	var (
		err     error
		msg     *databus.Message
		msgChan = s.accRegDatabus.Messages()
		ok      bool
	)
	for {
		select {
		case msg, ok = <-msgChan:
			if !ok {
				log.Info("accproc msgChan closed")
			}
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit err(%v)", err)
		}
		reg := &model.MsgCanal{}
		if err = json.Unmarshal([]byte(msg.Value), reg); err != nil {
			log.Error("json.Unmarshal(%v) err(%v)", msg, err)
			continue

		}
		log.Info("reg log %+v", reg)
		if reg.Action == _insertAction {
			var info struct {
				Mid int64 `json:"mid"`
			}
			if err = json.Unmarshal(reg.New, &info); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", string(reg.New), err)
				return
			}
			if err = s.AccountReg(context.Background(), info.Mid); err != nil {
				log.Error("s.AccountReg(%v) err(%v)", reg, err)
				continue
			}
			s.figureDao.SetWaiteUserCache(context.Background(), info.Mid, s.figureDao.Version(time.Now()))
		}
	}
}

func (s *Service) vipproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("s.vipproc panic(%v) %s", x, debug.Stack())
			go s.vipproc()
			log.Info("s.vipproc recover")
		}
	}()
	var (
		err     error
		msg     *databus.Message
		msgChan = s.vipDatabus.Messages()
		ok      bool
	)
	for {
		select {
		case msg, ok = <-msgChan:
			if !ok {
				log.Info("vipproc msgChan closed")
			}
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit err(%v)", err)
		}
		v := &model.MsgCanal{}
		if err = json.Unmarshal([]byte(msg.Value), v); err != nil {
			log.Error("json.Unmarshal(%v) err(%v)", v, err)
			continue
		}
		log.Info("vip log %+v", v)
		if v.Table == _vipTable {
			var vipInfo struct {
				Mid       int64 `json:"mid"`
				VipStatus int32 `json:"vip_status"`
			}
			if err = json.Unmarshal(v.New, &vipInfo); err != nil {
				log.Error("json.Unmarshal(%v) err(%v)", v.New, err)
				continue
			}
			if err = s.UpdateVipStatus(context.Background(), vipInfo.Mid, vipInfo.VipStatus); err != nil {
				log.Error("s.UpdateVipStatus(%v) err(%v)", v, err)
				continue
			}
			s.figureDao.SetWaiteUserCache(context.Background(), vipInfo.Mid, s.figureDao.Version(time.Now()))
		}
	}
}

func (s *Service) spyproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("s.spyproc panic(%v) %s", x, debug.Stack())
			go s.spyproc()
			log.Info("s.spyproc recover")
		}
	}()
	var (
		err     error
		msg     *databus.Message
		msgChan = s.spyDatabus.Messages()
		ok      bool
	)
	for {
		select {
		case msg, ok = <-msgChan:
			if !ok {
				log.Info("spyproc msgChan closed")
			}
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit err(%v)", err)
		}
		sc := &spym.ScoreChange{}
		if err = json.Unmarshal([]byte(msg.Value), sc); err != nil {
			log.Error("json.Unmarshal(%v) err(%v)", sc, err)
			continue
		}
		log.Info("spy log %+v", sc)
		if err = s.PutSpyScore(context.Background(), sc); err != nil {
			log.Error("s.PutSpyScore(%v) err(%v)", sc, err)
			continue
		}
		s.figureDao.SetWaiteUserCache(context.Background(), sc.Mid, s.figureDao.Version(time.Now()))
	}
}

func (s *Service) coinproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("s.coinproc panic(%v) %s", x, debug.Stack())
			go s.coinproc()
			log.Info("s.coinproc recover")
		}
	}()
	var (
		err     error
		msg     *databus.Message
		msgChan = s.coinDatabus.Messages()
		ok      bool
	)
	for {
		select {
		case msg, ok = <-msgChan:
			if !ok {
				log.Info("coinproc msgChan closed")
			}
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit err(%v)", err)
		}
		cd := &coinm.DataBus{}
		if err = json.Unmarshal([]byte(msg.Value), cd); err != nil {
			log.Error("json.Unmarshal(%v) err(%v)", cd, err)
			continue
		}
		log.Info("coin log %+v", cd)
		if err = s.PutCoinInfo(context.Background(), cd); err != nil {
			log.Error("s.PutCoinInfo(%v) err(%v)", cd, err)
			continue
		}
		s.figureDao.SetWaiteUserCache(context.Background(), cd.Mid, s.figureDao.Version(time.Now()))
	}
}

func (s *Service) replyinfoproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("s.replyinfoproc panic(%v) %s", x, debug.Stack())
			go s.replyinfoproc()
			log.Info("s.replyinfoproc recover")
		}
	}()
	var (
		err     error
		msg     *databus.Message
		msgChan = s.replyInfoDatabus.Messages()
		ok      bool
	)
	for {
		select {
		case msg, ok = <-msgChan:
			if !ok {
				log.Info("replyinfoproc msgChan closed")
			}
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit err(%v)", err)
		}
		m := &model.ReplyEvent{}
		if err = json.Unmarshal([]byte(msg.Value), m); err != nil {
			log.Error("json.Unmarshal(%v) err(%v)", m, err)
			continue
		}
		log.Info("reply log %+v", m)
		if err = s.PutReplyInfo(context.Background(), m); err != nil {
			log.Error("s.PutCoinInfo(%v) err(%v)", m, err)
			continue
		}
		if m.Action == model.EventAdd {
			s.figureDao.SetWaiteUserCache(context.Background(), m.Mid, s.figureDao.Version(time.Now()))
		} else {
			s.figureDao.SetWaiteUserCache(context.Background(), m.Reply.Mid, s.figureDao.Version(time.Now()))
		}
	}
}

func (s *Service) syncproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("s.syncproc panic(%v) %s", x, debug.Stack())
			go s.syncproc()
			log.Info("s.syncproc recover")
		}
	}()
	if s.c.Figure.Sync {
		log.Info("start import data after half hour.")
		time.Sleep(5 * time.Minute)
		log.Info("start import data.")
		var (
			vipPath = s.c.Figure.VipPath
			files   []os.FileInfo
			err     error
		)
		if files, err = ioutil.ReadDir(vipPath); err != nil {
			log.Error("ioutile.ReadDir(%s) err [%s]", vipPath, err)
		}
		for _, f := range files {
			s.SyncUserVIP(context.TODO(), vipPath+"/"+f.Name())
		}
		log.Info("end import data.")
	}
}

func (s *Service) payproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("s.payproc panic(%v) %s", x, debug.Stack())
			go s.payproc()
			log.Info("s.payproc recover")
		}
	}()
	var (
		err     error
		msg     *databus.Message
		msgChan = s.payDatabus.Messages()
		ok      bool
	)
	for {
		select {
		case msg, ok = <-msgChan:
			if !ok {
				log.Info("payproc msgChan closed")
			}
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit err(%v)", err)
		}
		v := &model.MsgCanal{}
		if err = json.Unmarshal([]byte(msg.Value), v); err != nil {
			log.Error("json.Unmarshal(%v) err(%v)", v, err)
			continue
		}
		log.Info("pay log %+v", v)
		if strings.HasPrefix(v.Table, _payTable) {
			var payOrder struct {
				Mid        int64   `json:"pay_mid"`
				Money      float64 `json:"bp"`
				MerchantID int8    `json:"merchant_id"`
				Status     int32   `json:"status"`
			}
			if err = json.Unmarshal(v.New, &payOrder); err != nil {
				log.Error("json.Unmarshal(%v) err(%v)", v.New, err)
				continue
			}
			if payOrder.Status != _paySuccess {
				continue
			}
			// update YUAN to Fen
			money := int64(payOrder.Money * 100)
			if err = s.PayOrderInfo(context.Background(), payOrder.Mid, money, payOrder.MerchantID); err != nil {
				log.Error("s.PayOrderInfo(%v) err(%v)", payOrder, err)
				continue
			}
			s.figureDao.SetWaiteUserCache(context.Background(), payOrder.Mid, s.figureDao.Version(time.Now()))
		}
	}
}

// 风纪委相关
func (s *Service) blockedproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("s.blockedproc panic(%v) %s", x, debug.Stack())
			go s.blockedproc()
			log.Info("s.blockedproc recover")
		}
	}()
	var (
		err     error
		msg     *databus.Message
		msgChan = s.blockedDatabus.Messages()
		ok      bool
	)
	for {
		select {
		case msg, ok = <-msgChan:
			if !ok {
				log.Info("blockedproc msgChan closed")
			}
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit err(%v)", err)
		}
		v := &model.MsgCanal{}
		if err = json.Unmarshal([]byte(msg.Value), v); err != nil {
			log.Error("json.Unmarshal(%v) err(%v)", v, err)
			continue
		}
		log.Info("blocked log %+v", v)
		if v.Table == _blockedTable {
			var KPIInfo struct {
				Mid  int64 `json:"mid"`
				Rate int16 `json:"rate"`
			}
			if err = json.Unmarshal(v.New, &KPIInfo); err != nil {
				log.Error("json.Unmarshal(%v) err(%v)", v.New, err)
				continue
			}
			if err = s.BlockedKPIInfo(context.Background(), KPIInfo.Mid, KPIInfo.Rate); err != nil {
				log.Error("s.BlockedKPIInfo(%v) err(%v)", KPIInfo, err)
				continue
			}
			s.figureDao.SetWaiteUserCache(context.Background(), KPIInfo.Mid, s.figureDao.Version(time.Now()))
		}
	}
}

func (s *Service) danmakuproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("s.danmakuproc panic(%v) %s", x, debug.Stack())
			go s.danmakuproc()
			log.Info("s.danmakuproc recover")
		}
	}()
	var (
		err     error
		msg     *databus.Message
		msgChan = s.danmakuDatabus.Messages()
		ok      bool
	)
	for {
		select {
		case msg, ok = <-msgChan:
			if !ok {
				log.Info("danmakuproc msgChan closed")
			}
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit err(%+v)", err)
		}
		m := &model.DMAction{}
		if err = json.Unmarshal([]byte(msg.Value), m); err != nil {
			log.Error("json.Unmarshal(%v) err(%+v)", m, err)
			continue
		}
		log.Info("danmaku msg %+v", m)
		if m.Action == _reportDel {
			if err = s.DanmakuReport(context.Background(), m); err != nil {
				log.Error("s.DanmakuReport(%v) err(%+v)", m, err)
				continue
			}
			s.figureDao.SetWaiteUserCache(context.Background(), m.Data.OwnerUID, s.figureDao.Version(time.Now()))
			s.figureDao.SetWaiteUserCache(context.Background(), m.Data.ReportUID, s.figureDao.Version(time.Now()))
		}
	}
}
