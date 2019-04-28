package service

import (
	"context"

	"go-common/app/job/main/aegis/model"
	"go-common/library/queue/databus"
	"go-common/library/xstr"
)

type baseTaskHandler struct {
	*Service
}

type dynamicTaskHandler struct {
	baseTaskHandler
}

func (h baseTaskHandler) CheckMessage(msg *databus.Message) (taskObj interface{}, err error) {
	return h.checkTaskMsg(msg)
}

func (h baseTaskHandler) HandleMessage(c context.Context, taskObj interface{}) error {
	return h.writeTaskToDB(c, taskObj.(*model.Task))
}

func (h dynamicTaskHandler) CheckMessage(msg *databus.Message) (taskObj interface{}, err error) {
	var c = context.Background()

	if taskObj, err = h.baseTaskHandler.CheckMessage(msg); err != nil {
		return
	}

	//补充mid相关信息
	task := taskObj.(*model.Task)
	res, err := h.dao.Resource(c, task.RID)
	if err != nil || res == nil {
		return nil, ErrTaskResourceInvalid
	}

	task.MID = res.MID
	if task.MID > 0 {
		groupids, _ := h.dao.UpSpecial(c, task.MID)
		task.Group = xstr.JoinInts(groupids)
		task.Fans, _ = h.dao.FansCount(c, task.MID)
	}

	taskObj = task
	return
}

func (h dynamicTaskHandler) HandleMessage(c context.Context, obj interface{}) error {
	return h.baseTaskHandler.HandleMessage(c, obj.(*model.Task))
}
