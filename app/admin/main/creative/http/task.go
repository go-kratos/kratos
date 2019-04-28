package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go-common/app/admin/main/creative/model/task"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
	"go-common/library/xstr"
)

func taskList(c *bm.Context) {
	var err error
	v := new(struct {
		Type int8 `form:"type"`
	})
	if err = c.Bind(v); err != nil {
		return
	}

	res, err := svc.TaskList(v.Type)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}

func taskPre(c *bm.Context) {
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data": map[string]interface{}{
			"target":     task.TargetMap,
			"root_type":  task.TaskRootNameMap,
			"group_type": task.TaskGroupNameMap,
		},
	}))
}

func viewGroup(c *bm.Context) {
	var err error
	v := new(struct {
		ID int64 `form:"id"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	res, err := svc.TaskGroup(v.ID)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}

func addGroup(c *bm.Context) {
	var err error
	v := new(task.TaskGroup)
	if err = c.Bind(v); err != nil {
		log.Error("addGroup c.Bind v(%+v) error(%v)", v, err)
		return
	}

	var rewardsIDs []int64
	if v.RewardIDs != "" {
		if rewardsIDs, err = xstr.SplitInts(v.RewardIDs); err != nil {
			log.Error("addGroup xstr.SplitInts task v.RewardIDs(%s) rewardsIDs(%+v)|error(%v)", v.RewardIDs, rewardsIDs, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}

	id, err := svc.AddTaskGroup(v, rewardsIDs)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendTaskLog(c, &task.LogParam{UID: uid, UName: uname, Action: "添加任务组", OID: id, Content: v})
	c.JSON(map[string]interface{}{
		"id": id,
	}, nil)
}

func editGroup(c *bm.Context) {
	var err error
	v := new(task.TaskGroup)
	if err = c.Bind(v); err != nil {
		log.Error("editGroup c.Bind v(%+v) error(%v)", v, err)
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var rewardsIDs []int64
	if v.RewardIDs != "" {
		if rewardsIDs, err = xstr.SplitInts(v.RewardIDs); err != nil {
			log.Error("editGroup xstr.SplitInts task v.RewardIDs(%s) rewardsIDs(%+v)|error(%v)", v.RewardIDs, rewardsIDs, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}

	id, err := svc.EditTaskGroup(v, rewardsIDs)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendTaskLog(c, &task.LogParam{UID: uid, UName: uname, Action: "编辑任务组", OID: id, Content: v})
	c.JSON(map[string]interface{}{
		"id": id,
	}, nil)
}

func orderGroup(c *bm.Context) {
	var err error
	v := new(task.OrderTask)
	if err = c.Bind(v); err != nil {
		log.Error("orderGroup c.Bind v(%+v) error(%v)", v, err)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendTaskLog(c, &task.LogParam{UID: uid, UName: uname, Action: "移动任务组顺序", OID: v.ID, Content: v})
	c.JSON(nil, svc.OrderTaskGroup(v))
}

func upStateGroup(c *bm.Context) {
	var err error
	v := new(struct {
		ID    int64 `form:"id" validate:"required"`
		State int8  `form:"state"`
	})
	if err = c.Bind(v); err != nil {
		log.Error("upStateGroup c.Bind v(%+v) error(%v)", v, err)
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendTaskLog(c, &task.LogParam{UID: uid, UName: uname, Action: "更新任务组状态", OID: v.ID, Content: v})
	c.JSON(nil, svc.UpStateGroup(v.ID, v.State))
}

func viewSubtask(c *bm.Context) {
	var err error
	v := new(struct {
		ID int64 `form:"id"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	res, err := svc.Task(v.ID)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}

func addSubtask(c *bm.Context) {
	var err error
	v := new(task.Task)
	if err = c.Bind(v); err != nil {
		log.Error("addSubtask c.Bind v(%+v) error(%v)", v, err)
		return
	}
	var rewardsIDs []int64
	if v.RewardIDs != "" {
		if rewardsIDs, err = xstr.SplitInts(v.RewardIDs); err != nil {
			log.Error("addSubtask xstr.SplitInts task v.RewardIDs(%s) rewardsIDs(%+v)|error(%v)", v.RewardIDs, rewardsIDs, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}

	id, err := svc.AddSubtask(v, rewardsIDs)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendTaskLog(c, &task.LogParam{UID: uid, UName: uname, Action: "添加子任务", OID: id, Content: v})
	c.JSON(map[string]interface{}{
		"id": id,
	}, nil)
}

func editSubtask(c *bm.Context) {
	var err error
	v := new(task.Task)
	if err = c.Bind(v); err != nil {
		log.Error("editSubtask c.Bind v(%+v) error(%v)", v, err)
		return
	}
	if v.ID == 0 {
		log.Error("editSubtask v.ID(%d)", v.ID)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var rewardsIDs []int64
	if v.RewardIDs != "" {
		if rewardsIDs, err = xstr.SplitInts(v.RewardIDs); err != nil {
			log.Error("addSubtask xstr.SplitInts task v.RewardIDs(%s) rewardsIDs(%+v)|error(%v)", v.RewardIDs, rewardsIDs, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}

	id, err := svc.EditSubtask(v, rewardsIDs)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendTaskLog(c, &task.LogParam{UID: uid, UName: uname, Action: "编辑子任务", OID: id, Content: v})
	c.JSON(map[string]interface{}{
		"id": id,
	}, nil)
}

func orderSubtask(c *bm.Context) {
	var err error
	v := new(task.OrderTask)
	if err = c.Bind(v); err != nil {
		log.Error("orderSubtask c.Bind v(%+v) error(%v)", v, err)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendTaskLog(c, &task.LogParam{UID: uid, UName: uname, Action: "移动子任务顺序", OID: v.ID, Content: v})
	c.JSON(nil, svc.OrderSubTask(v))
}

func upStateSubtask(c *bm.Context) {
	var err error
	v := new(struct {
		ID    int64 `form:"id" validate:"required"`
		State int8  `form:"state"`
	})
	if err = c.Bind(v); err != nil {
		log.Error("upStateSubtask c.Bind v(%+v) error(%v)", v, err)
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendTaskLog(c, &task.LogParam{UID: uid, UName: uname, Action: "更新子任务状态", OID: v.ID, Content: v})
	c.JSON(nil, svc.UpStateSubTask(v.ID, v.State))
}

func transferSubtask(c *bm.Context) {
	var err error
	v := new(struct {
		ID      int64 `form:"id" validate:"required"`
		GroupID int64 `form:"group_id" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		log.Error("transferSubtask c.Bind v(%+v) error(%v)", v, err)
		return
	}
	if v.ID == 0 || v.GroupID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendTaskLog(c, &task.LogParam{UID: uid, UName: uname, Action: "转移子任务到别的组", OID: v.ID, Content: v})
	c.JSON(nil, svc.TransferSubtask(v.ID, v.GroupID))
}

func rewardTree(c *bm.Context) {
	c.JSON(svc.RewardTree(), nil)
}

func viewReward(c *bm.Context) {
	var err error
	v := new(struct {
		ID int64 `form:"id"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	rd, err := svc.ViewReward(v.ID)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(rd, nil)
}

func addReward(c *bm.Context) {
	var err error
	v := new(task.Reward)
	if err = c.Bind(v); err != nil {
		log.Error("addReward c.Bind v(%+v) error(%v)", v, err)
		return
	}

	id, err := svc.AddReward(v)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendTaskLog(c, &task.LogParam{UID: uid, UName: uname, Action: "添加奖励", OID: id, Content: v})
	c.JSON(map[string]interface{}{
		"id": id,
	}, nil)
}

func editReward(c *bm.Context) {
	var err error
	v := new(task.Reward)
	if err = c.Bind(v); err != nil {
		log.Error("editReward c.Bind v(%+v) error(%v)", v, err)
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	id, err := svc.EditReward(v)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendTaskLog(c, &task.LogParam{UID: uid, UName: uname, Action: "编辑奖励", OID: id, Content: v})
	c.JSON(map[string]interface{}{
		"id": id,
	}, nil)
}

func upStateReward(c *bm.Context) {
	var err error
	v := new(struct {
		ID    int64 `form:"id" validate:"required"`
		State int8  `form:"state"`
	})
	if err = c.Bind(v); err != nil {
		log.Error("upStateReward c.Bind v(%+v) error(%v)", v, err)
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendTaskLog(c, &task.LogParam{UID: uid, UName: uname, Action: "更新奖励状态", OID: v.ID, Content: v})
	c.JSON(nil, svc.UpStateReward(v.ID, v.State))
}

func listGiftReward(c *bm.Context) {
	res, err := svc.ListGiftReward()
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}

func viewGiftReward(c *bm.Context) {
	var err error
	v := new(task.GiftReward)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.RootType == 0 || v.TaskType == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	res, err := svc.ViewGiftReward(v)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}

func addGiftReward(c *bm.Context) {
	var err error
	v := new(task.GiftReward)
	if err = c.Bind(v); err != nil {
		log.Error("addGiftReward c.Bind v(%+v) error(%v)", v, err)
		return
	}

	var rewardsIDs []int64
	if v.RewardIDs == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if rewardsIDs, err = xstr.SplitInts(v.RewardIDs); err != nil || len(rewardsIDs) == 0 { //注意礼包配置，奖励必传
		log.Error("addGiftReward xstr.SplitInts task v.RewardIDs(%s) rewardsIDs(%+v)|error(%v)", v.RewardIDs, rewardsIDs, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}

	rows, err := svc.AddGiftReward(v, rewardsIDs)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendTaskLog(c, &task.LogParam{UID: uid, UName: uname, Action: "添加礼包", OID: 0, Content: v})
	c.JSON(map[string]interface{}{
		"rows": rows,
	}, nil)
}

func editGiftReward(c *bm.Context) {
	var err error
	v := new(task.GiftReward)
	if err = c.Bind(v); err != nil {
		log.Error("editGiftReward c.Bind v(%+v) error(%v)", v, err)
		return
	}
	if v.RootType == 0 || v.TaskType == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	var rewardsIDs []int64
	if v.RewardIDs == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if rewardsIDs, err = xstr.SplitInts(v.RewardIDs); err != nil || len(rewardsIDs) == 0 { //注意礼包配置，奖励必传
		log.Error("editGiftReward xstr.SplitInts task v.RewardIDs(%s) rewardsIDs(%+v)|error(%v)", v.RewardIDs, rewardsIDs, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}

	rows, err := svc.EditGiftReward(v, rewardsIDs)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendTaskLog(c, &task.LogParam{UID: uid, UName: uname, Action: "编辑礼包", OID: 0, Content: v})
	c.JSON(map[string]interface{}{
		"rows": rows,
	}, nil)
}

func upStateGiftReward(c *bm.Context) {
	var err error
	v := new(task.GiftReward)
	if err = c.Bind(v); err != nil {
		log.Error("upStateGiftReward c.Bind v(%+v) error(%v)", v, err)
		return
	}
	if v.RootType == 0 || v.TaskType == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	rows, err := svc.UpGiftReward(v)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendTaskLog(c, &task.LogParam{UID: uid, UName: uname, Action: "更新礼包状态", OID: 0, Content: v})
	c.JSON(map[string]interface{}{
		"rows": rows,
	}, nil)
}

func batchOnline(c *bm.Context) {
	var err error
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error("batchOnline ioutil.ReadAll error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.Request.Body.Close()

	var v []*task.TaskGroup
	err = json.Unmarshal(bs, &v)
	if err != nil {
		log.Error("batchOnline json.Unmarshal v(%+v) error(%v)", v, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}

	err = svc.BatchOnline(v)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendTaskLog(c, &task.LogParam{UID: uid, UName: uname, Action: "批量刷新上线", OID: 0, Content: v})
	c.JSON(nil, nil)
}
