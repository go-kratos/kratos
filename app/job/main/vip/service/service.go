package service

import (
	"context"
	"encoding/json"
	xlog "log"
	"strconv"
	"sync"
	"time"

	"go-common/app/job/main/vip/conf"
	"go-common/app/job/main/vip/dao"
	"go-common/app/job/main/vip/model"
	couponrpc "go-common/app/service/main/coupon/rpc/client"
	v1 "go-common/app/service/main/vip/api"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

const (
	_tableUserInfo = "vip_user_info"
	_tablePayOrder = "vip_pay_order"
	_insertAction  = "insert"
	_updateAction  = "update"
	_deleteAction  = "delete"

	notifyAction = "updateVip"

	_ps           = 50
	_defsleepmsec = 100
)

//Service vip service
type Service struct {
	dao *dao.Dao
	c   *conf.Config
	//vipRPC                   *client.Service
	reducePayOrder           map[string]*model.VipPayOrder
	appMap                   map[int64]*model.VipAppInfo
	confMap                  map[string]*model.VipConfig
	cleanVipCache            chan int64
	cleanAppCache            chan *model.AppCache
	ds                       *databus.Databus
	handlerFailPayOrder      chan *model.VipPayOrder
	handlerFailUserInfo      chan *model.VipUserInfo
	handlerFailRechargeOrder chan *model.VipPayOrder
	handlerFailVipbuy        chan *model.VipBuyResq
	handlerInsertOrder       chan *model.VipPayOrder
	handlerUpdateOrder       chan *model.VipPayOrder
	handlerRechargeOrder     chan *model.VipPayOrder
	handlerInsertUserInfo    chan *model.VipUserInfo
	handlerUpdateUserInfo    chan *model.VipUserInfo
	handlerStationActive     chan *model.VipPayOrder
	handlerAutoRenewLog      chan *model.VipUserInfo
	handlerAddVipHistory     chan *model.VipChangeHistoryMsg
	handlerAddBcoinSalary    chan *model.VipBcoinSalaryMsg
	handlerUpdateBcoinSalary chan *model.VipBcoinSalaryMsg
	handlerDelBcoinSalary    chan *model.VipBcoinSalaryMsg
	notifycouponchan         chan func()
	accLogin                 *databus.Databus
	frozenDate               time.Duration
	newVipDatabus            *databus.Databus
	salaryCoupnDatabus       *databus.Databus
	accountNoitfyDatabus     *databus.Databus
	couponNotifyDatabus      *databus.Databus
	autoRenewdDatabus        *databus.Databus
	// waiter
	waiter      sync.WaitGroup
	closed      bool
	sendmsgchan chan func()
	couponRPC   *couponrpc.Service
	// vip service
	vipgRPC v1.VipClient
}

//New new service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
		//vipRPC:                   client.New(c.VipRPC),
		cleanVipCache:            make(chan int64, 10240),
		confMap:                  make(map[string]*model.VipConfig),
		cleanAppCache:            make(chan *model.AppCache, 10240),
		handlerFailPayOrder:      make(chan *model.VipPayOrder, 10240),
		handlerFailUserInfo:      make(chan *model.VipUserInfo, 10240),
		handlerFailRechargeOrder: make(chan *model.VipPayOrder, 10240),
		handlerInsertOrder:       make(chan *model.VipPayOrder, 10240),
		handlerRechargeOrder:     make(chan *model.VipPayOrder, 10240),
		handlerUpdateOrder:       make(chan *model.VipPayOrder, 10240),
		handlerInsertUserInfo:    make(chan *model.VipUserInfo, 10240),
		handlerUpdateUserInfo:    make(chan *model.VipUserInfo, 10240),
		handlerStationActive:     make(chan *model.VipPayOrder, 10240),
		handlerAutoRenewLog:      make(chan *model.VipUserInfo, 10240),
		handlerAddVipHistory:     make(chan *model.VipChangeHistoryMsg, 10240),
		handlerAddBcoinSalary:    make(chan *model.VipBcoinSalaryMsg, 10240),
		handlerUpdateBcoinSalary: make(chan *model.VipBcoinSalaryMsg, 10240),
		handlerDelBcoinSalary:    make(chan *model.VipBcoinSalaryMsg, 10240),
		handlerFailVipbuy:        make(chan *model.VipBuyResq, 10240),
		notifycouponchan:         make(chan func(), 10240),
		ds:                       databus.New(c.Databus.OldVipBinLog),
		newVipDatabus:            databus.New(c.Databus.NewVipBinLog),
		accLogin:                 databus.New(c.Databus.AccLogin),
		frozenDate:               time.Duration(c.Property.FrozenDate),
		reducePayOrder:           make(map[string]*model.VipPayOrder),
		sendmsgchan:              make(chan func(), 10240),
		accountNoitfyDatabus:     databus.New(c.Databus.AccountNotify),
		couponRPC:                couponrpc.New(c.RPCClient2.Coupon),
	}
	vipgRPC, err := v1.NewClient(c.VipClient)
	if err != nil {
		panic(err)
	}
	s.vipgRPC = vipgRPC
	t := cron.New()
	go s.loadappinfoproc()
	go s.cleanappcacheretryproc()
	go s.cleanvipretryproc()
	go s.sendmessageproc()
	go s.handlerfailpayorderproc()
	go s.handlerfailrechargeorderproc()
	go s.handlerfailuserinfoproc()
	go s.handlerautorenewlogproc()
	go s.handlerdelbcoinproc()
	if c.Databus.SalaryCoupon != nil {
		s.salaryCoupnDatabus = databus.New(c.Databus.SalaryCoupon)
		s.waiter.Add(1)
		go s.salarycouponproc()
	}
	if c.Databus.CouponNotify != nil {
		s.couponNotifyDatabus = databus.New(c.Databus.CouponNotify)
		go s.couponnotifyproc()
		s.waiter.Add(1)
		go s.couponnotifybinlogproc()
	}
	for i := 0; i < s.c.Property.HandlerThread; i++ {
		go s.handlerinsertorderproc()
		go s.handlerupdateorderproc()
		go s.handlerinsertuserinfoproc()
		go s.handlerupdateuserinfoproc()
		go s.handleraddchangehistoryproc()
		go s.handleraddbcoinproc()
		go s.handlerupdatebcoinproc()
		go s.handlerupdaterechargeorderproc()
	}
	for i := 0; i < s.c.Property.ReadThread; i++ {
		go s.readdatabusproc()
	}
	go s.readnewvipdatabusproc()
	if c.Property.FrozenCron != "" {
		go s.accloginproc()
		t.AddFunc(c.Property.FrozenCron, s.unFrozenJob)
	}
	t.AddFunc(c.Property.UpdateUserInfoCron, s.updateUserInfoJob)
	t.AddFunc(c.Property.SalaryVideoCouponCron, s.salaryVideoCouponJob)
	t.AddFunc(c.Property.PushDataCron, s.pushDataJob)
	t.AddFunc(c.Property.EleEompensateCron, s.eleEompensateJob)
	//t.AddFunc(c.Property.HadExpiredMsgCron, s.hadExpiredMsgJob)
	//t.AddFunc(c.Property.WillExpireMsgCron, s.willExpiredMsgJob)
	//t.AddFunc(c.Property.SendMessageCron, s.sendMessageJob)
	//t.AddFunc(c.Property.AutoRenewCron, s.autoRenewJob)
	//t.AddFunc(c.Property.SendBcoinCron, s.sendBcoinJob)
	t.Start()
	go s.consumercheckproc()
	if c.Databus.AutoRenew != nil {
		s.autoRenewdDatabus = databus.New(c.Databus.AutoRenew)
		s.waiter.Add(1)
		go s.retryautorenewpayproc()
	}
	return
}

func (s *Service) readnewvipdatabusproc() {
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("Runtime error caught: %+v", r)
			go s.readnewvipdatabusproc()
		}
	}()
	var err error

	for msg := range s.newVipDatabus.Messages() {
		val := msg.Value
		if err = msg.Commit(); err != nil {
			log.Error("readdatabusproc msg.commit() error(%v)", err)
			msg.Commit()
		}
		log.Info("cur consumer new vip db message(%v)", string(msg.Value))
		message := new(model.Message)
		if err = json.Unmarshal(val, message); err != nil {
			log.Error("readnewvipdatabusproc json.unmarshal val(%+v) error(%+v)", string(val), err)
			continue
		}
		if message.Table == "vip_user_info" {
			userInfo := new(model.VipUserInfoNewMsg)
			if err = json.Unmarshal(message.New, userInfo); err != nil {
				log.Error("readdatabusproc json.Unmarshal val(%v) error(%v)", string(message.New), err)
				continue
			}
			vipUser := convertUserInfoByNewMsg(userInfo)
			s.dao.DelInfoCache(context.Background(), vipUser.Mid)
			if message.Action == _insertAction {
				if vipUser.PayType == model.AutoRenew {
					select {
					case s.handlerAutoRenewLog <- vipUser:
					default:
						log.Error("s.handlerAutoRenewLog full!")
					}
				}
			} else if message.Action == _updateAction {
				oldUserMsg := new(model.VipUserInfoNewMsg)
				if err = json.Unmarshal(message.Old, oldUserMsg); err != nil {
					log.Error("readdatabusproc json.Unmarshal val(%v) error(%v)", string(message.Old), err)
					continue
				}
				oldUser := convertUserInfoByNewMsg(oldUserMsg)

				if oldUser.PayType != vipUser.PayType {
					select {
					case s.handlerAutoRenewLog <- vipUser:
					default:
						log.Error("s.handlerAutoRenewLog full update!")
					}
				}
			}
			s.pubAccountNotify(vipUser.Mid)
		}

	}
}

func (s *Service) readdatabusproc() {
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("Runtime error caught: %+v", r)
			go s.readdatabusproc()
		}
	}()
	var err error
	for msg := range s.ds.Messages() {
		val := msg.Value
		message := new(model.Message)

		if err = json.Unmarshal(val, message); err != nil {
			log.Error("readdatabusproc json.Unmarshal val(%v) error(%v)", string(val), err)
			if err = msg.Commit(); err != nil {
				log.Error("msg.commit() error(%v)", err)
			}
			continue
		}

		if message.Table == "vip_pay_order" {
			order := new(model.VipPayOrderOldMsg)
			if err = json.Unmarshal(message.New, order); err != nil {
				log.Error("readdatabusproc json.Unmarshal val(%v) error(%v)", string(message.New), err)
				if err = msg.Commit(); err != nil {
					log.Error("msg.commit() error(%v)", err)
				}
				continue
			}

			payOrder := s.convertPayOrder(order)
			if message.Action == "insert" {
				select {
				case s.handlerInsertOrder <- payOrder:
				default:
					xlog.Panic("s.handlerInsertOrder full!")
				}
			} else if message.Action == "update" {
				select {
				case s.handlerUpdateOrder <- payOrder:
				default:
					xlog.Panic("s.handlerUpdateOrder full!")
				}
			}

		} else if message.Table == "vip_recharge_order" {
			order := new(model.VipRechargeOrderMsg)
			if err = json.Unmarshal(message.New, order); err != nil {
				log.Error("readdatabusproc json.Unmarshal val(%v) error(%v)", string(message.New), err)
				if err = msg.Commit(); err != nil {
					log.Error("msg.commit() error(%v)", err)
				}
				continue
			}

			payOrder := s.convertPayOrderByMsg(order)
			if message.Action == "update" {
				select {
				case s.handlerRechargeOrder <- payOrder:
				default:
					xlog.Panic("s.handlerRechargeOrder full!")
				}
			} else if message.Action == "insert" {
				if len(payOrder.ThirdTradeNo) > 0 {
					select {
					case s.handlerRechargeOrder <- payOrder:
					default:
						xlog.Panic("s.handlerRechargeOrder full!")
					}
				}
			}

		} else if message.Table == "vip_user_info" {
			userInfo := new(model.VipUserInfoMsg)
			if err = json.Unmarshal(message.New, userInfo); err != nil {
				log.Error("readdatabusproc json.Unmarshal val(%v) error(%v)", string(message.New), err)
				if err = msg.Commit(); err != nil {
					log.Error("msg.commit() error(%v)", err)
				}
				continue
			}
			vipUser := convertMsgToUserInfo(userInfo)
			if message.Action == "insert" {
				select {
				case s.handlerInsertUserInfo <- vipUser:
				default:
					xlog.Panic("s.handlerInsertUserInfo full!")
				}
			} else if message.Action == "update" {
				oldUser := new(model.VipUserInfoMsg)
				if err = json.Unmarshal(message.Old, oldUser); err != nil {
					log.Error("readdatabusproc json.Unmarshal val(%v) error(%v)", string(message.Old), err)
					if err = msg.Commit(); err != nil {
						log.Error("msg.commit() error(%v)", err)
					}
					continue
				}
				vipUser.OldVer = oldUser.Ver
				select {
				case s.handlerUpdateUserInfo <- vipUser:
				default:
					xlog.Panic("s.handlerUpdateUserInfo full!")
				}
			}
			if !s.grayScope(userInfo.Mid) {
				s.cleanCache(userInfo.Mid)
			}

		} else if message.Table == "vip_change_history" {
			historyMsg := new(model.VipChangeHistoryMsg)
			if err = json.Unmarshal(message.New, historyMsg); err != nil {
				log.Error("readdatabusproc json.Unmarshal val(%v) error(%v)", string(message.New), err)
				if err = msg.Commit(); err != nil {
					log.Error("msg.commit() error(%v)", err)
				}
				continue
			}
			if message.Action == "insert" {
				select {
				case s.handlerAddVipHistory <- historyMsg:
				default:
					xlog.Panic("s.handlerAddVipHistory full!")
				}
			}
		} else if message.Table == "vip_bcoin_salary" {
			bcoinMsg := new(model.VipBcoinSalaryMsg)

			if err = json.Unmarshal(message.New, bcoinMsg); err != nil {
				log.Error("readdatabusproc json.Unmarshal val(%v) error(%v)", string(message.New), err)
				if err = msg.Commit(); err != nil {
					log.Error("msg.commit() error(%v)", err)
				}
				continue
			}
			if message.Action == _insertAction {
				select {
				case s.handlerAddBcoinSalary <- bcoinMsg:
				default:
					xlog.Panic("s.handlerAddBcoinSalary full!")
				}
			} else if message.Action == _updateAction {
				select {
				case s.handlerUpdateBcoinSalary <- bcoinMsg:
				default:
					xlog.Panic("s.handlerUpdateBcoinSalary full!")
				}
			} else if message.Action == _deleteAction {
				select {
				case s.handlerDelBcoinSalary <- bcoinMsg:
				default:
					xlog.Panic("s.handlerDelBcoinSalary full!")
				}
			}
		}

		if err = msg.Commit(); err != nil {
			log.Error("readdatabusproc msg.commit() error(%v)", err)
			msg.Commit()
		}
		log.Info("cur consumer message(%v)", string(msg.Value))
	}
}

func (s *Service) cleanCache(mid int64) {
	var (
		hv  = new(model.HandlerVip)
		err error
	)
	hv.Type = 2
	hv.Days = 0
	hv.Months = 0
	hv.Mid = mid
	if err = s.cleanCacheAndNotify(context.TODO(), hv); err != nil {
		select {
		case s.cleanVipCache <- hv.Mid:
		default:
			xlog.Panic("s.cleanVipCache full!")
		}
	}
	s.pubAccountNotify(mid)
}

func (s *Service) pubAccountNotify(mid int64) (err error) {

	data := new(struct {
		Mid    int64  `json:"mid"`
		Action string `json:"action"`
	})
	data.Mid = mid
	data.Action = notifyAction
	if err = s.accountNoitfyDatabus.Send(context.TODO(), strconv.FormatInt(mid, 10), data); err != nil {
		log.Error("send (%+v) error(%+v)", data, err)
	}
	log.Info("send(mid:%+v) data:%+v", mid, data)
	return
}

func (s *Service) loadappinfoproc() {
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("Runtime error caught: %+v", r)
			go s.loadappinfoproc()
		}
	}()
	for {
		s.loadAppInfo()
		time.Sleep(time.Minute * 2)
	}
}

func (s *Service) loadAppInfo() {
	var (
		res []*model.VipAppInfo
		err error
	)
	if res, err = s.dao.SelAppInfo(context.TODO()); err != nil {
		log.Error("loadAppInfo SelAppInfo error(%v)", err)
		return
	}
	aMap := make(map[int64]*model.VipAppInfo, len(res))
	for _, v := range res {
		aMap[v.ID] = v
	}
	s.appMap = aMap
	bytes, _ := json.Marshal(res)
	log.Info("load app success :%v", string(bytes))
}

func (s *Service) cleanvipretryproc() {
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("Runtime error caught: %+v", r)
			go s.cleanvipretryproc()
		}
	}()
	for {
		mid := <-s.cleanVipCache
		s.cleanVipRetry(mid)
	}

}

func (s *Service) cleanVipRetry(mid int64) {

	hv := new(model.HandlerVip)
	hv.Type = 2
	hv.Days = 0
	hv.Months = 0
	hv.Mid = mid
	s.dao.DelInfoCache(context.Background(), mid)
	for i := 0; i < s.c.Property.Retry; i++ {
		if err := s.cleanCacheAndNotify(context.TODO(), hv); err == nil {
			break
		}
		s.dao.DelVipInfoCache(context.TODO(), int64(hv.Mid))

	}
	log.Info("handler success cache fail mid(%v)", mid)
}

func (s *Service) cleanappcacheretryproc() {
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("Runtime error caught: %+v", r)
			go s.cleanappcacheretryproc()
		}
	}()
	for {
		ac := <-s.cleanAppCache
		s.cleanAppCacheRetry(ac)
	}
}

func (s *Service) cleanAppCacheRetry(ac *model.AppCache) {
	appInfo := s.appMap[ac.AppID]
	hv := new(model.HandlerVip)
	hv.Type = 2
	hv.Days = 0
	hv.Months = 0
	hv.Mid = ac.Mid
	for i := 0; i < s.c.Property.Retry; i++ {
		if err := s.dao.SendAppCleanCache(context.TODO(), hv, appInfo); err == nil {
			break
		}
	}
	log.Info("handler success cache app fail appInfo(%v)", ac)
}

func (s *Service) salarycouponproc() {
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("Runtime error caught: %+v", r)
			go s.salarycouponproc()
		}
	}()
	defer s.waiter.Done()
	var (
		err     error
		msg     *databus.Message
		msgChan = s.salaryCoupnDatabus.Messages()
		ok      bool
		c       = context.Background()
	)
	for {
		msg, ok = <-msgChan
		if !ok || s.closed {
			log.Info("salary coupon msgChan closed")
			return
		}
		msg.Commit()
		v := &model.Message{}
		if err = json.Unmarshal([]byte(msg.Value), v); err != nil {
			log.Error("json.Unmarshal(%v) err(%v)", v, err)
			continue
		}
		if v.Table != _tableUserInfo {
			continue
		}
		nvip := &model.VipUserInfoMsg{}
		if err = json.Unmarshal(v.New, &nvip); err != nil {
			log.Error("salary new json.Unmarshal values(%v),error(%v)", string(v.New), err)
			continue
		}
		ovip := &model.VipUserInfoMsg{}
		if v.Action != _insertAction {
			if err = json.Unmarshal(v.Old, &ovip); err != nil {
				log.Error("salary old json.Unmarshal values(%v),error(%v)", string(v.Old), err)
				continue
			}
		}
		log.Info("salary coupon start mid(%d)", nvip.Mid)
		if _, err = s.SalaryVideoCouponAtOnce(c, nvip, ovip, v.Action); err != nil {
			log.Error("SalaryVideoCouponAtOnce fail(%d) nvip(%v) ovip(%v) %s error(%v)", nvip.Mid, nvip, ovip, v.Action, err)
			continue
		}
		log.Info("salary coupon suc mid(%d)", nvip.Mid)
	}
}

func (s *Service) couponnotifybinlogproc() {
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("Runtime error couponnotifybinlogproc caught: %+v", r)
			go s.couponnotifybinlogproc()
		}
	}()
	defer s.waiter.Done()
	var (
		err     error
		msg     *databus.Message
		msgChan = s.couponNotifyDatabus.Messages()
		ok      bool
		c       = context.Background()
	)
	for {
		msg, ok = <-msgChan
		if !ok || s.closed {
			log.Info("coupon  notify couponnotifybinlogproc msgChan closed")
			return
		}
		if err = msg.Commit(); err != nil {
			log.Error("couponnotifybinlogproc msg.Commit err(%v)", err)
			continue
		}
		log.Info("cur consumer couponnotifybinlogproc(%v)", string(msg.Value))
		v := &model.Message{}
		if err = json.Unmarshal([]byte(msg.Value), v); err != nil {
			log.Error("couponnotifybinlogproc json.Unmarshal(%v) err(%v)", v, err)
			continue
		}
		if v.Table != _tablePayOrder || v.Action != _updateAction {
			continue
		}
		newo := new(model.VipPayOrderNewMsg)
		if err = json.Unmarshal(v.New, newo); err != nil {
			log.Error("couponnotifybinlogproc json.Unmarshal val(%v) error(%v)", string(v.New), err)
			continue
		}
		oldo := new(model.VipPayOrderNewMsg)
		if err = json.Unmarshal(v.Old, oldo); err != nil {
			log.Error("couponnotifybinlogproc json.Unmarshal val(%v) error(%v)", string(v.Old), err)
			continue
		}
		if newo == nil || oldo == nil {
			continue
		}
		if oldo.Status != model.PAYING {
			continue
		}
		if newo.Status != model.SUCCESS && newo.Status != model.FAILED {
			continue
		}
		if newo.CouponMoney <= 0 {
			continue
		}
		s.couponnotify(func() {
			s.CouponNotify(c, newo)
		})
	}
}

func (s *Service) updateUserInfoJob() {
	log.Info("update user info job start ....................................")
	s.ScanUserInfo(context.TODO())
	log.Info("update user info job end ........................................")
}

func (s *Service) salaryVideoCouponJob() {
	log.Info("salary video coupon job start ....................................")
	var err error
	if ok := s.dao.AddTransferLock(context.TODO(), "_transferLock"); !ok {
		log.Info("salary video coupon job had run ....................................")
		return
	}
	if err = s.ScanSalaryVideoCoupon(context.TODO()); err != nil {
		log.Error("ScanSalaryVideoCoupon error(%v)", err)
		return
	}
	log.Info("salary video coupon job end ........................................")
}

//Ping check db live
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close all resource.
func (s *Service) Close() {
	defer s.waiter.Wait()
	s.closed = true
	s.salaryCoupnDatabus.Close()
	s.dao.Close()
	s.ds.Close()
	s.newVipDatabus.Close()
}

func (s *Service) sendmessageproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.sendmessageproc panic(%v)", x)
			go s.sendmessageproc()
			log.Info("service.sendmessageproc recover")
		}
	}()
	for {
		f := <-s.sendmsgchan
		f()
	}
}

func (s *Service) sendmessage(f func()) {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.sendmessage panic(%v)", x)
		}
	}()
	select {
	case s.sendmsgchan <- f:
	default:
		log.Error("service.sendmessage chan full")
	}
}

func (s *Service) consumercheckproc() {
	for {
		time.Sleep(time.Second)
		log.Info("consumercheckproc chan(cleanVipCache) size: %d", len(s.cleanVipCache))
		log.Info("consumercheckproc chan(cleanAppCache) size: %d", len(s.cleanAppCache))
		log.Info("consumercheckproc chan(handlerFailPayOrder) size: %d", len(s.handlerFailPayOrder))
		log.Info("consumercheckproc chan(handlerFailUserInfo) size: %d", len(s.handlerFailUserInfo))
		log.Info("consumercheckproc chan(handlerFailRechargeOrder) size: %d", len(s.handlerFailRechargeOrder))
		log.Info("consumercheckproc chan(handlerFailVipbuy) size: %d", len(s.handlerFailVipbuy))
		log.Info("consumercheckproc chan(handlerInsertOrder) size: %d", len(s.handlerInsertOrder))
		log.Info("consumercheckproc chan(handlerUpdateOrder) size: %d", len(s.handlerUpdateOrder))
		log.Info("consumercheckproc chan(handlerRechargeOrder) size: %d", len(s.handlerRechargeOrder))
		log.Info("consumercheckproc chan(handlerInsertUserInfo) size: %d", len(s.handlerInsertUserInfo))
		log.Info("consumercheckproc chan(handlerUpdateUserInfo) size: %d", len(s.handlerUpdateUserInfo))
		log.Info("consumercheckproc chan(handlerStationActive) size: %d", len(s.handlerStationActive))
		log.Info("consumercheckproc chan(handlerAutoRenewLog) size: %d", len(s.handlerAutoRenewLog))
		log.Info("consumercheckproc chan(handlerAddVipHistory) size: %d", len(s.handlerAddVipHistory))
		log.Info("consumercheckproc chan(handlerAddBcoinSalary) size: %d", len(s.handlerAddBcoinSalary))
		log.Info("consumercheckproc chan(handlerUpdateBcoinSalary) size: %d", len(s.handlerUpdateBcoinSalary))
		log.Info("consumercheckproc chan(handlerDelBcoinSalary) size: %d", len(s.handlerDelBcoinSalary))
		log.Info("consumercheckproc chan(sendmsgchan) size: %d", len(s.sendmsgchan))
		log.Info("consumercheckproc chan(notifycouponchan)  size: %d", len(s.notifycouponchan))
	}
}
