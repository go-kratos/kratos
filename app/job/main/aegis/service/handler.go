package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"go-common/app/job/main/aegis/model"
	"go-common/library/log"
	"go-common/library/queue/databus"

	pkgerr "github.com/pkg/errors"
)

//RscHandler .
type RscHandler interface {
	CheckMessage(json.RawMessage) (interface{}, error)
	HandleMessage(context.Context, interface{}) error
}

//TaskHandler .
type TaskHandler interface {
	CheckMessage(*databus.Message) (interface{}, error)
	HandleMessage(context.Context, interface{}) error
}

var (
	_ TaskHandler = baseTaskHandler{}
	_ TaskHandler = dynamicTaskHandler{}
	_ RscHandler  = baseResourceAddHandler{}
	_ RscHandler  = mangaResourceAddHandler{}
	_ RscHandler  = baseResourceUpdateHandler{}
	_ RscHandler  = baseResourceCancelHandler{}
)

//单例
var (
	basehandleTask      *baseTaskHandler
	basehandleRscAdd    *baseResourceAddHandler
	basehandleRscUpdate *baseResourceUpdateHandler
	basehandleRscCancel *baseResourceCancelHandler
	dynamicHandleTask   *dynamicTaskHandler
	mangaHandelRscAdd   *mangaResourceAddHandler
	once                sync.Once
)

//ERROR
var (
	ErrTaskDuplicate       = errors.New("重复任务")
	ErrTaskFlowInvalid     = errors.New("流程失效")
	ErrTaskResourceInvalid = errors.New("资源失效")
	ErrInvalidMsg          = errors.New("无效消息")
	ErrHandlerMiss         = errors.New("handler NotFound")
)

//prefix
var (
	_prefixTask      = "task_"
	_prefixRscAdd    = "add_"
	_prefixRscUpdate = "update_"
	_prefixRscCancel = "cancel_"
)

//业务ID
var (
	_bizidDynamic = 1
	_bizidManga   = 2
)

func (s *Service) registerRscHandler(key string, handler RscHandler) {
	s.rschandle[key] = handler
}

func (s *Service) registerTaskHandler(key string, handler TaskHandler) {
	s.taskhandle[key] = handler
}

func (s *Service) findTaskHandler(key string) TaskHandler {
	if handler, ok := s.taskhandle[key]; ok {
		return handler
	}
	log.Warn("key(%s)没找到任务的处理器，根据类型使用默认handler", key)
	return s.getdynamicTaskHandler()
}

func (s *Service) findRscHandler(key string) RscHandler {
	if handler, ok := s.rschandle[key]; ok {
		return handler
	}

	log.Warn("key(%s)没找到业务的处理器，根据类型使用默认handler", key)
	switch {
	case strings.HasPrefix(key, _prefixRscAdd):
		return s.getbaseResourceAddHandler()
	case strings.HasPrefix(key, _prefixRscUpdate):
		return s.getbaseResourceUpdateHandler()
	case strings.HasPrefix(key, _prefixRscCancel):
		return s.getbaseResourceCancelHandler()
	default:
		return nil
	}
}

//TODO 先写死吧,之后可以根据配置里面的类名用反射实例化
func initHandler(s *Service) {
	var (
		dynamicTask      = fmt.Sprintf("%s%d", _prefixTask, _bizidDynamic)
		dynamicRscAdd    = fmt.Sprintf("%s%d", _prefixRscAdd, _bizidDynamic)
		dynamicRscUpdate = fmt.Sprintf("%s%d", _prefixRscUpdate, _bizidDynamic)
		dynamicRscCancel = fmt.Sprintf("%s%d", _prefixRscCancel, _bizidDynamic)
		managaTask       = fmt.Sprintf("%s%d", _prefixTask, _bizidManga)
		managaRscAdd     = fmt.Sprintf("%s%d", _prefixRscAdd, _bizidManga)
		managaRscUpdate  = fmt.Sprintf("%s%d", _prefixRscUpdate, _bizidManga)
		managaRscCancel  = fmt.Sprintf("%s%d", _prefixRscCancel, _bizidManga)
	)
	s.rschandle = make(map[string]RscHandler)
	s.taskhandle = make(map[string]TaskHandler)

	once.Do(func() {
		basehandleTask = &baseTaskHandler{Service: s}
		basehandleRscAdd = &baseResourceAddHandler{Service: s}
		basehandleRscUpdate = &baseResourceUpdateHandler{Service: s}
		basehandleRscCancel = &baseResourceCancelHandler{Service: s}
		dynamicHandleTask = &dynamicTaskHandler{baseTaskHandler: baseTaskHandler{Service: s}}
		mangaHandelRscAdd = &mangaResourceAddHandler{baseResourceAddHandler: baseResourceAddHandler{Service: s}}
	})

	s.registerRscHandler(dynamicRscAdd, s.getbaseResourceAddHandler())
	s.registerRscHandler(dynamicRscUpdate, s.getbaseResourceUpdateHandler())
	s.registerRscHandler(dynamicRscCancel, s.getbaseResourceCancelHandler())
	s.registerRscHandler(managaRscAdd, s.getmangaResourceAddHandler())
	s.registerRscHandler(managaRscUpdate, s.getbaseResourceUpdateHandler())
	s.registerRscHandler(managaRscCancel, s.getbaseResourceCancelHandler())

	s.registerTaskHandler(managaTask, s.getbaseTaskHandler())
	s.registerTaskHandler(dynamicTask, s.getdynamicTaskHandler())
}

func (s *Service) getbaseTaskHandler() *baseTaskHandler {
	return basehandleTask
}
func (s *Service) getbaseResourceAddHandler() *baseResourceAddHandler {
	return basehandleRscAdd
}
func (s *Service) getbaseResourceUpdateHandler() *baseResourceUpdateHandler {
	return basehandleRscUpdate
}
func (s *Service) getbaseResourceCancelHandler() *baseResourceCancelHandler {
	return basehandleRscCancel
}
func (s *Service) getdynamicTaskHandler() *dynamicTaskHandler {
	return dynamicHandleTask
}
func (s *Service) getmangaResourceAddHandler() *mangaResourceAddHandler {
	return mangaHandelRscAdd
}

//解析验证message
/*
TODO
根据DispatchLimit，动态设置分发数量
*/
func (s *Service) checkTaskMsg(msg *databus.Message) (*model.Task, error) {
	taskMsg := new(model.CreateTaskMsg)
	if err := json.Unmarshal(msg.Value, taskMsg); err != nil {
		log.Error("checkTaskMsg key(%s) value(%s)", msg.Key, string(msg.Value))
		return nil, err
	}

	if taskMsg.DispatchLimit == 0 || taskMsg.FlowID == 0 || taskMsg.RID == 0 {
		log.Error("checkTaskMsg key(%s) value(%s)", msg.Key, string(msg.Value))
		return nil, ErrTaskResourceInvalid
	}

	if s.dao.CheckTask(context.Background(), taskMsg.FlowID, taskMsg.RID) > 0 {
		return nil, ErrTaskDuplicate
	}

	ok, err := s.dao.CheckFlow(context.TODO(), taskMsg.RID, taskMsg.FlowID)
	if !ok || err != nil {
		return nil, ErrTaskFlowInvalid
	}

	//先兼容旧的task消息，没有传bizid
	if taskMsg.BizID == 0 {
		res, err := s.dao.Resource(context.Background(), taskMsg.RID)
		if err != nil || res == nil {
			return nil, ErrTaskResourceInvalid
		}
		taskMsg.BizID = res.BusinessID
	}

	return &model.Task{
		BusinessID: taskMsg.BizID,
		FlowID:     taskMsg.FlowID,
		RID:        taskMsg.RID,
	}, nil
}

func (s *Service) writeTaskToDB(c context.Context, task *model.Task) error {
	return s.dao.CreateTask(c, task)
}

func (s *Service) checkRscAddMsg(msg json.RawMessage) (*model.AddOption, error) {
	addMsg := new(model.AddOption)
	if err := json.Unmarshal(msg, addMsg); err != nil {
		return nil, err
	}

	if addMsg.BusinessID == 0 || len(addMsg.OID) == 0 {
		return nil, ErrInvalidMsg
	}
	return addMsg, nil
}

func (s *Service) writeRscAdd(c context.Context, opt *model.AddOption) error {
	//TODO 根据错误号重试
	return s.dao.RscAdd(c, opt)
}

func (s *Service) checkRscUpdateMsg(msg json.RawMessage) (*model.UpdateOption, error) {
	updateMsg := new(model.UpdateOption)
	if err := json.Unmarshal(msg, updateMsg); err != nil {
		return nil, err
	}

	if updateMsg.BusinessID == 0 || len(updateMsg.OID) == 0 || len(updateMsg.Update) == 0 {
		return nil, ErrInvalidMsg
	}
	return updateMsg, nil
}

func (s *Service) writeRscUpdate(c context.Context, opt *model.UpdateOption) error {
	return s.dao.RscUpdate(c, opt)
}

func (s *Service) checkRscCancelMsg(msg json.RawMessage) (*model.CancelOption, error) {
	cancelMsg := new(model.CancelOption)
	if err := json.Unmarshal(msg, cancelMsg); err != nil {
		return nil, err
	}

	if cancelMsg.BusinessID == 0 || len(cancelMsg.Oids) == 0 {
		return nil, ErrInvalidMsg
	}
	return cancelMsg, nil
}

func (s *Service) writeRscCancel(c context.Context, opt *model.CancelOption) error {
	return s.dao.RscCancel(c, opt)
}

func (s *Service) newrsc(msg *databus.Message) (interface{}, error) {
	log.Info("databusgroup new msg key(%+v) partition(%d) offset(%d) value(%s) ", msg.Key, msg.Partition, msg.Offset, string(msg.Value))

	rscmsg := new(model.RscMsg)
	if err := json.Unmarshal(msg.Value, rscmsg); err != nil {
		log.Error("databusgroup  json.Unmarshal for msg(%+v)", string(msg.Value))
		return nil, ErrInvalidMsg
	}

	key := fmt.Sprintf("%s_%d", rscmsg.Action, rscmsg.BizID)
	handler := s.findRscHandler(key)
	if handler == nil {
		log.Error("databusgroup can not find handler for msg key(%+v)", key)
		return nil, ErrHandlerMiss
	}
	data, err := handler.CheckMessage(rscmsg.Raw)
	if err != nil {
		log.Error("databusgroup new msg key(%+v) partition(%d) offset(%d) value(%s) CheckMessage(%v)", msg.Key, msg.Partition, msg.Offset, string(msg.Value), pkgerr.WithStack(err))
	}
	return data, err
}

func (s *Service) splitrsc(msg *databus.Message, data interface{}) int {
	switch t := data.(type) {
	case *model.AddOption:
		return int(t.BusinessID)
	case *model.UpdateOption:
		return int(t.BusinessID)
	case *model.CancelOption:
		return int(t.BusinessID)
	default:
		return 0
	}
}

func (s *Service) dorsc(bmsgs []interface{}) {
	for _, msg := range bmsgs {
		log.Info("databusgroup do msg(%+v)", msg)
		var key string
		switch t := msg.(type) {
		case *model.AddOption:
			key = fmt.Sprintf("%s%d", _prefixRscAdd, t.BusinessID)
		case *model.UpdateOption:
			key = fmt.Sprintf("%s%d", _prefixRscUpdate, t.BusinessID)
		case *model.CancelOption:
			key = fmt.Sprintf("%s%d", _prefixRscCancel, t.BusinessID)
		default:
			log.Error("databusgroup unknow msg(%+v)", msg)
			continue
		}
		handler := s.findRscHandler(key)
		if handler == nil {
			log.Error("databusgroup msg(%+v) handler NotFound", msg)
			continue
		}
		if err := handler.HandleMessage(context.Background(), msg); err != nil {
			log.Error("databusgroup msg(%+v) handler err(%v)", msg, pkgerr.WithStack(err))
			continue
		}
	}
}

func (s *Service) newtask(msg *databus.Message) (interface{}, error) {
	log.Info("databusgroup newtask msg key(%+v) partition(%d) offset(%d) value(%s) ", msg.Key, msg.Partition, msg.Offset, string(msg.Value))

	taskmsg := new(model.CreateTaskMsg)
	if err := json.Unmarshal(msg.Value, taskmsg); err != nil {
		log.Error("databusgroup newtask json.Unmarshal for msg(%+v)", string(msg.Value))
		return nil, ErrInvalidMsg
	}

	key := fmt.Sprintf("%s%d", _prefixTask, taskmsg.BizID)
	handler := s.findTaskHandler(key)
	if handler == nil {
		log.Error("databusgroup can not find handler for msg key(%+v)", key)
		return nil, ErrHandlerMiss
	}
	data, err := handler.CheckMessage(msg)
	if err != nil {
		errmsg := fmt.Sprintf("databusgroup new msg key(%+v) partition(%d) offset(%d) value(%s) CheckMessage(%v)", msg.Key, msg.Partition, msg.Offset, string(msg.Value), pkgerr.WithStack(err))
		if err == ErrTaskDuplicate {
			log.Warn(errmsg)
		} else {
			log.Error(errmsg)
		}
	}
	return data, err
}

func (s *Service) splittask(msg *databus.Message, data interface{}) int {
	if t, ok := data.(*model.Task); ok {
		return int(t.BusinessID)
	}
	return 0
}

func (s *Service) dotask(bmsgs []interface{}) {
	for _, msg := range bmsgs {
		log.Info("databusgroup dotask msg(%+v)", msg)
		var key string
		if t, ok := msg.(*model.Task); ok {
			key = fmt.Sprintf("%s%d", _prefixTask, t.BusinessID)
		} else {
			log.Error("databusgroup dotask unknow msg(%+v)", msg)
			continue
		}

		handler := s.findTaskHandler(key)
		if handler == nil {
			log.Error("databusgroup dotask msg(%+v) handler NotFound", msg)
			continue
		}
		if err := handler.HandleMessage(context.Background(), msg); err != nil {
			log.Error("databusgroup dotask msg(%+v) handler err(%v)", msg, pkgerr.WithStack(err))
			continue
		}
	}
}
