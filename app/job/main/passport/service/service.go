package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"sync"

	"go-common/app/job/main/passport/conf"
	"go-common/app/job/main/passport/dao"
	igrpc "go-common/app/service/main/identify-game/rpc/client"
	"go-common/library/queue/databus"
)

const (
	_gameAppID      = int32(876)
	_telBindTable   = "aso_telephone_bind_log"
	_emailBindTable = "aso_email_bind_log"
)

// Service struct of service.
type Service struct {
	c *conf.Config
	d *dao.Dao
	// RPC
	igRPC *igrpc.Client
	// game app ids
	gameAppIDs []int32
	// token proc
	dsToken         *databus.Databus
	tokenMergeChans []chan *message
	tokenDoneChan   chan []*message
	head, last      *message
	mu              sync.Mutex
	// user proc
	dsUser             *databus.Databus
	userMergeChans     []chan *message
	userDoneChan       chan []*message
	userHead, userLast *message
	userMu             sync.Mutex
	// log proc
	dsLog            *databus.Databus
	logMergeChans    []chan *message
	logDoneChan      chan []*message
	logHead, logLast *message
	logMu            sync.Mutex
	// pwd log proc
	dsPwdLog               *databus.Databus
	pwdLogMergeChans       []chan *message
	pwdLogDoneChan         chan []*message
	pwdLogHead, pwdLogLast *message
	pwdLogMu               sync.Mutex

	// auth bin log proc
	authBinLog                     *databus.Databus
	authBinLogMergeChans           []chan *message
	authBinLogDoneChan             chan []*message
	authBinLogHead, authBinLogLast *message
	authBinLogMu                   sync.Mutex

	dsContactBindLog                       *databus.Databus
	contactBindLogMergeChans               []chan *message
	contactBindLogDoneChan                 chan []*message
	contactBindLogHead, contactBindLogLast *message
	contactBindLogMu                       sync.Mutex
	userLogPub                             *databus.Databus

	AESBlock cipher.Block
	hashSalt []byte
}

type message struct {
	next   *message
	data   *databus.Message
	object interface{}
	done   bool
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	gameAppIDs := make([]int32, 0)
	gameAppIDs = append(gameAppIDs, _gameAppID)
	for _, id := range c.Game.AppIDs {
		if id == _gameAppID {
			continue
		}
		gameAppIDs = append(gameAppIDs, id)
	}

	s = &Service{
		c:          c,
		d:          dao.New(c),
		gameAppIDs: gameAppIDs,
		// RPC
		igRPC: igrpc.New(c.RPC.IdentifyGame),
		// token
		dsToken:         databus.New(c.DataBus.AsoBinLog),
		tokenMergeChans: make([]chan *message, c.Group.AsoBinLog.Num),
		tokenDoneChan:   make(chan []*message, c.Group.AsoBinLog.Chan),
		// user
		dsUser:         databus.New(c.DataBus.User),
		userMergeChans: make([]chan *message, c.Group.User.Num),
		userDoneChan:   make(chan []*message, c.Group.User.Chan),
		// log
		dsLog:         databus.New(c.DataBus.Log),
		logMergeChans: make([]chan *message, c.Group.Log.Num),
		logDoneChan:   make(chan []*message, c.Group.Log.Chan),
		// log
		dsPwdLog:         databus.New(c.DataBus.PwdLog),
		pwdLogMergeChans: make([]chan *message, c.Group.PwdLog.Num),
		pwdLogDoneChan:   make(chan []*message, c.Group.PwdLog.Chan),

		// emial and tel log
		dsContactBindLog:         databus.New(c.DataBus.ContactBindLog),
		contactBindLogMergeChans: make([]chan *message, c.Group.ContactBindLog.Num),
		contactBindLogDoneChan:   make(chan []*message, c.Group.ContactBindLog.Chan),
		userLogPub:               databus.New(c.DataBus.UserLog),

		// auth bin log
		authBinLog:           databus.New(c.DataBus.AuthBinLog),
		authBinLogMergeChans: make([]chan *message, c.Group.AuthBinLog.Num),
		authBinLogDoneChan:   make(chan []*message, c.Group.AuthBinLog.Chan),

		hashSalt: []byte(c.Encode.Salt),
	}
	s.AESBlock, _ = aes.NewCipher([]byte(c.Encode.AesKey))

	// start token proc
	go s.tokencommitproc()
	for i := 0; i < c.Group.AsoBinLog.Num; i++ {
		ch := make(chan *message, c.Group.AsoBinLog.Chan)
		s.tokenMergeChans[i] = ch
		go s.tokenmergeproc(ch)
	}
	go s.tokenconsumeproc()
	// start user proc
	go s.usercommitproc()
	for i := 0; i < c.Group.User.Num; i++ {
		ch := make(chan *message, c.Group.User.Chan)
		s.userMergeChans[i] = ch
		go s.usermergeproc(ch)
	}
	go s.userconsumeproc()
	// start log proc
	go s.logcommitproc()
	for i := 0; i < c.Group.Log.Num; i++ {
		ch := make(chan *message, c.Group.Log.Chan)
		s.logMergeChans[i] = ch
		go s.logmergeproc(ch)
	}
	go s.logconsumeproc()
	// start pwd log proc
	go s.pwdlogcommitproc()
	for i := 0; i < c.Group.PwdLog.Num; i++ {
		ch := make(chan *message, c.Group.PwdLog.Chan)
		s.pwdLogMergeChans[i] = ch
		go s.pwdlogmergeproc(ch)
	}
	go s.pwdlogconsumeproc()

	go s.contactBindLogcommitproc()
	for i := 0; i < c.Group.ContactBindLog.Num; i++ {
		ch := make(chan *message, c.Group.ContactBindLog.Chan)
		s.contactBindLogMergeChans[i] = ch
		go s.contactBindLogMergeproc(ch)
	}
	go s.contactBindLogconsumeproc()

	// start auth bin log token proc
	go s.authBinLogcommitproc()
	for i := 0; i < c.Group.AuthBinLog.Num; i++ {
		ch := make(chan *message, c.Group.AuthBinLog.Chan)
		s.authBinLogMergeChans[i] = ch
		go s.authBinLogmergeproc(ch)
	}
	go s.authBinLogconsumeproc()
	// end auth bin log token proc

	go s.syncPwdLog()
	return
}

// Ping ping check service health.
func (s *Service) Ping(c context.Context) (err error) {
	err = s.d.Ping(c)
	return
}

// Close close service.
func (s *Service) Close() (err error) {
	if s.dsToken != nil {
		s.dsToken.Close()
	}
	if s.dsUser != nil {
		s.dsUser.Close()
	}
	if s.dsLog != nil {
		s.dsLog.Close()
	}
	if s.d != nil {
		s.d.Close()
	}
	if s.dsContactBindLog != nil {
		s.dsContactBindLog.Close()
	}
	if s.userLogPub != nil {
		s.userLogPub.Close()
	}
	if s.dsPwdLog != nil {
		s.dsPwdLog.Close()
	}
	return
}
