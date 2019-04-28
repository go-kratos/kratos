package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-common/app/admin/main/dm/conf"
	"go-common/app/admin/main/dm/dao"
	oplogDao "go-common/app/admin/main/dm/dao/oplog"
	"go-common/app/admin/main/dm/model"
	"go-common/app/admin/main/dm/model/oplog"
	accountApi "go-common/app/service/main/account/api"
	archive "go-common/app/service/main/archive/api/gorpc"
	"go-common/library/log"
	"go-common/library/log/infoc"
	"go-common/library/sync/pipeline/fanout"
)

// Service define Service struct
type Service struct {
	// dao
	dao      *dao.Dao
	oplogDao *oplogDao.Dao
	// rpc
	accountRPC        accountApi.AccountClient
	arcRPC            *archive.Service2
	dmOperationLogSvc *infoc.Infoc
	bakInfoc          *infoc.Infoc
	opsLogCh          chan *oplog.Infoc
	reduceMoralChan   chan *model.ReduceMoral
	blockUserChan     chan *model.BlockUser
	msgReporterChan   chan *model.ReportMsg
	msgPosterChan     chan *model.ReportMsg
	actionChan        chan *model.Action
	// async proc
	cache      *fanout.Fanout
	moniOidMap map[int64]struct{}
	oidLock    sync.Mutex
}

// New new a Service and return.
func New(c *conf.Config) *Service {
	s := &Service{
		// dao
		dao:      dao.New(c),
		oplogDao: oplogDao.New(c),
		// rpc
		arcRPC:            archive.New2(c.ArchiveRPC),
		dmOperationLogSvc: infoc.New(c.Infoc2),
		bakInfoc:          infoc.New(c.InfocBak),
		reduceMoralChan:   make(chan *model.ReduceMoral, 1024),
		blockUserChan:     make(chan *model.BlockUser, 1024),
		msgReporterChan:   make(chan *model.ReportMsg, 1024),
		msgPosterChan:     make(chan *model.ReportMsg, 1024),
		actionChan:        make(chan *model.Action, 1024),
		opsLogCh:          make(chan *oplog.Infoc, 1024),
		cache:             fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		moniOidMap:        make(map[int64]struct{}),
	}
	accountRPC, err := accountApi.NewClient(c.AccountRPC)
	if err != nil {
		panic(err)
	}
	s.accountRPC = accountRPC
	go s.changeReportStatProc()
	go s.actionproc()
	go s.oplogproc()
	go s.monitorproc()
	return s
}

// Ping check server ok
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.dao.Ping(c); err != nil {
		return
	}
	return
}

func (s *Service) addAction(action *model.Action) {
	select {
	case s.actionChan <- action:
	default:
		log.Error("action channel is full,action(%v) is discard", action)
	}
}

func (s *Service) actionproc() {
	for action := range s.actionChan {
		if err := s.dao.SendAction(context.TODO(), fmt.Sprint(action.Oid), action); err != nil {
			log.Error("dao.SendAction(%v) error(%v)", action, err)
		}
	}
}

func (s *Service) changeReportStatProc() {
	for {
		select {
		case msg := <-s.reduceMoralChan:
			if err := s.dao.ReduceMoral(context.TODO(), msg); err != nil {
				log.Error("s.dao.ReduceMoral(msg:%v) error(%v)", msg, err)
			}
		case msg := <-s.msgReporterChan:
			if err := s.dao.SendMsgToReporter(context.TODO(), msg); err != nil {
				log.Error("s.dao.SendMsgToReporter(msg:%v) error(%v)", msg, err)
			}
		case msg := <-s.msgPosterChan:
			if err := s.dao.SendMsgToPoster(context.TODO(), msg); err != nil {
				log.Error("s.dao.SendMsgToPoster(msg:%v) error(%v)", msg, err)
			}
		case msg := <-s.blockUserChan:
			if err := s.dao.BlockUser(context.TODO(), msg); err != nil {
				log.Error("s.dao.BlockUser(msg:%v) error(%v)", msg, err)
			}
		}
	}
}

func (s *Service) oplogproc() {
	for opLog := range s.opsLogCh {
		if len(opLog.Subject) == 0 || len(opLog.CurrentVal) == 0 || opLog.Source <= 0 ||
			opLog.Operator <= 0 || opLog.OperatorType <= 0 {
			log.Warn("oplogproc() it is an illegal log, warn(%v, %v, %v)", opLog.Subject, opLog.Subject, opLog.CurrentVal)
			continue
		} else {
			for _, dmid := range opLog.DMIds {
				s.dmOperationLogSvc.Info(opLog.Subject, strconv.FormatInt(opLog.Oid, 10), fmt.Sprint(opLog.Type),
					strconv.FormatInt(dmid, 10), opLog.Source.String(), opLog.OriginVal,
					opLog.CurrentVal, strconv.FormatInt(opLog.Operator, 10), opLog.OperatorType.String(),
					opLog.OperationTime, opLog.Remark)
				// 将管理员操作日志额外上报一份（用于验证数据报表完整性）
				s.bakInfoc.Info(opLog.Subject, strconv.FormatInt(opLog.Oid, 10), fmt.Sprint(opLog.Type),
					strconv.FormatInt(dmid, 10), opLog.Source.String(), opLog.OriginVal,
					opLog.CurrentVal, strconv.FormatInt(opLog.Operator, 10), opLog.OperatorType.String(),
					opLog.OperationTime, opLog.Remark)
				if strings.Contains(opLog.Subject, "\n") {
					log.Error("\n found in opLog.Subject(%s)", opLog.Source)
				}
			}
		}
	}
}

// OpLog put a new infoc format operation log into the channel
func (s *Service) OpLog(c context.Context, cid, operator int64, typ int32, dmids []int64, subject, originVal, currentVal, remark string, source oplog.Source, operatorType oplog.OperatorType) (err error) {
	infoLog := new(oplog.Infoc)
	infoLog.Oid = cid
	infoLog.Type = typ
	infoLog.DMIds = dmids
	infoLog.Subject = subject
	infoLog.OriginVal = originVal
	infoLog.CurrentVal = currentVal
	infoLog.OperationTime = strconv.FormatInt(time.Now().Unix(), 10)
	infoLog.Source = source
	infoLog.OperatorType = operatorType
	infoLog.Operator = operator
	infoLog.Remark = remark
	select {
	case s.opsLogCh <- infoLog:
	default:
		err = fmt.Errorf("opsLogCh full")
		log.Error("opsLogCh full (%v)", infoLog)
	}
	return
}

// QueryOpLogs query operation logs of damku equals dmid
func (s *Service) QueryOpLogs(c context.Context, dmid int64) (infos []*oplog.InfocResult, err error) {
	result, err := s.oplogDao.QueryOpLogs(c, dmid)
	if err != nil {
		return
	}
	for _, logVal := range result {
		var (
			tmp = &oplog.InfocResult{}
		)
		val, err := strconv.Atoi(logVal.CurrentVal)
		if err != nil {
			err = nil
			continue
		}
		switch logVal.Subject {
		case "status":
			tmp.Subject = model.StateDesc(int32(val))
		case "pool":
			if val == 0 {
				tmp.Subject = "普通弹幕池"
			} else if val == 1 {
				tmp.Subject = "字幕弹幕池"
			} else if val == 2 {
				tmp.Subject = "特殊弹幕池"
			}
		case "attribute":
			if val == 2 || (int32(val)>>model.AttrProtect)&int32(1) == 1 {
				tmp.Subject = "弹幕保护"
			} else if val == 3 || (int32(val)>>model.AttrProtect)&int32(1) == 0 {
				tmp.Subject = "取消弹幕保护"
			}
		default:
			tmp.Subject = logVal.Subject
		}
		tmp.CurrentVal = logVal.CurrentVal
		tmp.OperatorType = logVal.OperatorType
		if logVal.OperatorType == "用户" || logVal.OperatorType == "UP主" {
			mid, _ := strconv.ParseInt(logVal.Operator, 10, 64)
			arg3 := &accountApi.MidReq{Mid: mid}
			uInfo, err := s.accountRPC.Info3(c, arg3)
			if err != nil {
				tmp.Operator = logVal.Operator
				log.Error("s.accRPC.Info2(%v) error(%v)", arg3, err)
				err = nil
			} else {
				tmp.Operator = uInfo.GetInfo().GetName()
			}
		} else {
			tmp.Operator = logVal.Operator
		}
		OperationTimeStamp, _ := strconv.ParseInt(logVal.OperationTime, 10, 64)
		tmp.OperationTime = time.Unix(OperationTimeStamp, 0).Format("2006-01-02 15:04:05")
		tmp.Remark = "操作来源：" + logVal.Source + "；操作员身份：" + logVal.OperatorType + "；备注：" + logVal.Remark
		infos = append(infos, tmp)
	}
	return
}

func (s *Service) monitorproc() {
	for {
		time.Sleep(3 * time.Second)
		s.oidLock.Lock()
		oidMap := s.moniOidMap
		s.moniOidMap = make(map[int64]struct{})
		s.oidLock.Unlock()
		for oid := range oidMap {
			sub, err := s.dao.Subject(context.TODO(), model.SubTypeVideo, oid)
			if err != nil || sub == nil {
				continue
			}
			if err := s.updateMonitorCnt(context.TODO(), sub); err != nil {
				log.Error("s.updateMonitorCnt(%+v) error(%v)", sub, err)
			}
		}
	}
}
