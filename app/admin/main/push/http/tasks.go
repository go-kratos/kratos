package http

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/push/model"
	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"
)

func taskList(c *bm.Context) {
	var (
		items []*model.Task
		types []int
		err   error
		pager = new(model.Pager)
	)
	params := new(struct {
		Type  string `form:"type"`
		Stime string `form:"stime"`
		Etime string `form:"etime"`
	})
	if err = c.Bind(params); err != nil {
		return
	}
	if err = c.Bind(pager); err != nil {
		return
	}
	for _, v := range strings.Split(params.Type, ",") {
		i, e := strconv.Atoi(v)
		if e != nil {
			continue
		}
		types = append(types, i)
	}
	if params.Stime == "" {
		params.Stime = time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")
	}
	if params.Etime == "" {
		params.Etime = time.Now().Format("2006-01-02 15:04:05")
	}
	if err = pushSrv.DB.Model(&model.Task{}).Where("push_time between ? and ?", params.Stime, params.Etime).Where("type in(?)", types).Order("push_time desc").Find(&items).Error; err != nil {
		log.Error("taskList(%d,%s,%s) error(%v)", params.Type, params.Stime, params.Etime, err)
		c.JSON(nil, err)
		return
	}
	res := fmtTasks(items)
	total := len(res)
	data := map[string]interface{}{
		"pager": &model.Pager{
			Pn:    pager.Pn,
			Ps:    pager.Ps,
			Total: total,
		},
		"data": []*model.Task{},
	}
	start := pager.Ps * (pager.Pn - 1)
	if start >= total {
		c.JSONMap(data, nil)
		return
	}
	end := start + pager.Ps
	if len(res[start:]) < pager.Ps {
		end = start + len(res[start:])
	}
	data["data"] = res[start:end]
	c.JSONMap(data, nil)
}

func fmtTasks(items []*model.Task) (res []*model.Task) {
	var jobs []string
	tasks := make(map[string]*pushmdl.Task)
	for _, t := range items {
		p := &pushmdl.Progress{}
		if t.Progress != "" {
			if err := json.Unmarshal([]byte(t.Progress), &p); err != nil {
				log.Error("unmarshal task(%d) progress(%s) error(%v)", t.ID, t.Progress, err)
				continue
			}
		}
		extra := new(pushmdl.TaskExtra)
		if t.Extra != "" {
			if err := json.Unmarshal([]byte(t.Extra), &extra); err != nil {
				log.Error("unmarshal task(%d) extra(%s) error(%v)", t.ID, t.Extra, err)
				continue
			}
		}
		if v, ok := tasks[t.Job]; !ok {
			job, _ := strconv.ParseInt(t.Job, 10, 64)
			tasks[t.Job] = &pushmdl.Task{
				ID:         t.ID,
				Job:        job,
				Type:       t.Type,
				APPID:      t.AppID,
				BusinessID: t.BusinessID,
				Title:      t.Title,
				Summary:    t.Summary,
				LinkType:   int8(t.LinkType),
				LinkValue:  t.LinkValue,
				Progress:   p,
				PushTime:   xtime.Time(t.PushTime.Unix()),
				ExpireTime: xtime.Time(t.ExpireTime.Unix()),
				Status:     int8(t.Status),
				Extra:      extra,
			}
			jobs = append(jobs, t.Job)
		} else {
			v.Status = calStatus(v.Status, int8(t.Status))
			vp := v.Progress
			vp.MidTotal += p.MidTotal
			vp.MidValid += p.MidValid
			vp.TokenTotal += p.TokenTotal
			vp.TokenValid += p.TokenValid
			vp.TokenFailed += p.TokenFailed
		}
	}
	for _, j := range jobs {
		v := tasks[j]
		p, _ := json.Marshal(v.Progress)
		e, _ := json.Marshal(v.Extra)
		t := &model.Task{
			ID:             v.ID,
			Job:            strconv.FormatInt(v.Job, 10),
			Type:           v.Type,
			AppID:          v.APPID,
			BusinessID:     v.BusinessID,
			Title:          v.Title,
			Summary:        v.Summary,
			LinkType:       int(v.LinkType),
			LinkValue:      v.LinkValue,
			Progress:       string(p),
			PushTimeUnix:   int64(v.PushTime),
			ExpireTimeUnix: int64(v.ExpireTime),
			Status:         int(v.Status),
			Extra:          string(e),
		}
		res = append(res, t)
	}
	return
}

func calStatus(currStatus, newStatus int8) int8 {
	// 聚合任务状态
	// 如果有进行中的任务，显示进行中
	// 如果没有，显示其中最小的状态（负数是异常状态）
	if currStatus == pushmdl.TaskStatusDoing || newStatus == pushmdl.TaskStatusDoing {
		return pushmdl.TaskStatusDoing
	}
	if newStatus == pushmdl.TaskStatusPrepared {
		return pushmdl.TaskStatusPrepared
	}
	if newStatus < currStatus {
		return newStatus
	}
	return currStatus
}

func addTask(c *bm.Context) {
	var (
		task     = &model.Task{}
		filename = c.Request.Form.Get("filename")
	)
	if err := c.Bind(task); err != nil {
		return
	}
	task.PushTime = time.Unix(task.PushTimeUnix, 0)
	task.ExpireTime = time.Unix(task.ExpireTimeUnix, 0)
	job := pushmdl.JobName(time.Now().UnixNano(), task.Summary, task.LinkValue, task.Group)
	task.Job = strconv.FormatInt(job, 10)
	extra, _ := json.Marshal(pushmdl.TaskExtra{
		Group:    task.Group,
		Filename: filename,
	})
	task.Extra = string(extra)
	task.Status = int(pushmdl.TaskStatusPending)
	c.JSON(nil, pushSrv.AddTask(context.Background(), task))
}

func taskInfo(c *bm.Context) {
	task := &model.Task{}
	id, _ := strconv.ParseInt(c.Request.Form.Get("id"), 10, 64)
	if id == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := pushSrv.DB.First(task, id).Error; err != nil {
		log.Error("taskInfo(%d) error(%v)", id, err)
		c.JSON(nil, err)
		return
	}
	task.PushTimeUnix = task.PushTime.Unix()
	task.ExpireTimeUnix = task.ExpireTime.Unix()
	c.JSON(task, nil)
}

func saveTask(c *bm.Context) {
	var (
		task     = new(model.Task)
		filename = c.Request.Form.Get("filename")
	)
	if err := c.Bind(task); err != nil {
		return
	}
	if task.Job == "" {
		log.Warn("job is empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	extra, _ := json.Marshal(pushmdl.TaskExtra{
		Group:    task.Group,
		Filename: filename,
	})
	data := map[string]interface{}{
		"app_id":      task.AppID,
		"type":        task.Type,
		"business_id": task.BusinessID,
		"title":       task.Title,
		"summary":     task.Summary,
		"link_type":   task.LinkType,
		"link_value":  task.LinkValue,
		"build":       task.Build,
		"sound":       task.Sound,
		"vibration":   task.LinkValue,
		"push_time":   time.Unix(task.PushTimeUnix, 0),
		"expire_time": time.Unix(task.ExpireTimeUnix, 0),
		"group":       task.Group,
		"image_url":   task.ImageURL,
		"extra":       string(extra),
	}
	if err := pushSrv.DB.Model(&model.Task{}).Where("job=?", task.Job).Updates(data).Error; err != nil {
		log.Error("saveTask(%+v) error(%v)", data, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func delTask(c *bm.Context) {
	job := c.Request.Form.Get("job")
	if job == "" {
		log.Error("job is empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := pushSrv.DB.Model(&model.Task{}).Where("job=?", job).Update("dtime", time.Now().Unix()).Error; err != nil {
		log.Error("delTask(%s) error(%v)", job, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func stopTask(c *bm.Context) {
	job := c.Request.Form.Get("job")
	if job == "" {
		log.Error("job is empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := pushSrv.DB.Model(&model.Task{}).Where("job=?", job).Update("status", pushmdl.TaskStatusStop).Error; err != nil {
		log.Error("stopTask(%s) error(%v)", job, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func confirmTask(c *bm.Context) {
	job := c.Request.Form.Get("job")
	if job == "" {
		log.Error("job is empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	task := &model.Task{}
	if err := pushSrv.DB.Model(&model.Task{}).Where("job=?", job).First(task).Error; err != nil {
		log.Error("confirmTask(%s) query task error(%v)", job, err)
		c.JSON(nil, err)
		return
	}
	status := pushmdl.TaskStatusPretreatmentPrepared
	if task.Type == pushmdl.TaskTypeDataPlatformMid || task.Type == pushmdl.TaskTypeDataPlatformToken {
		status = pushmdl.TaskStatusWaitDataPlatform
		if err := pushSrv.UpdateDpCondtionStatus(c, job, pushmdl.DpCondStatusPrepared); err != nil {
			log.Error("confirmTask(%s) update data platform conditions error(%v)", job, err)
			return
		}
	}
	if err := pushSrv.DB.Model(&model.Task{}).Where("job=?", job).Update("status", status).Error; err != nil {
		log.Error("confirmTask(%s) error(%v)", job, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func testPushMid(c *bm.Context) {
	task := &model.Task{}
	if err := c.Bind(task); err != nil {
		return
	}
	builds := make(map[int]*pushmdl.Build)
	if task.Build != "" {
		if err := json.Unmarshal([]byte(task.Build), &builds); err != nil {
			c.JSON(nil, ecode.RequestErr)
			log.Warn("buildformat task(%+v) error(%v)", task, err)
			return
		}
	}
	var (
		plats []string
		mid   = c.Request.Form.Get("mid")
	)
	for plat := range builds {
		plats = append(plats, strconv.Itoa(plat))
	}
	task.Platform = strings.Join(plats, ",")
	task.PushTime = time.Unix(task.PushTimeUnix, 0)
	task.ExpireTime = time.Unix(task.ExpireTimeUnix, 0)
	job := pushmdl.JobName(time.Now().UnixNano(), task.Summary, task.LinkValue, task.Group)
	task.Job = strconv.FormatInt(job, 10)
	extra, _ := json.Marshal(pushmdl.TaskExtra{Group: task.Group})
	task.Extra = string(extra)
	c.JSON(nil, pushSrv.TestPushMid(c, mid, task))
}

func testPushToken(c *bm.Context) {
	params := c.Request.Form
	appID, _ := strconv.ParseInt(params.Get("app_id"), 10, 64)
	if appID < 1 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("app_id is wrong: %s", params.Get("app_id"))
		return
	}
	platform, _ := strconv.Atoi(params.Get("platform"))
	if platform < 1 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("platform is wrong: %s", params.Get("platform"))
		return
	}
	alertTitle := params.Get("alert_title")
	if alertTitle == "" {
		alertTitle = pushmdl.DefaultMessageTitle
	}
	alertBody := params.Get("alert_body")
	if alertBody == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("alert_body is empty")
		return
	}
	token := params.Get("token")
	if token == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("token is empty")
		return
	}
	linkType, _ := strconv.Atoi(params.Get("link_type"))
	if linkType < 1 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("link_type is wrong: %s", params.Get("link_type"))
		return
	}
	linkValue := params.Get("link_value")
	expireTime, _ := strconv.ParseInt(params.Get("expire_time"), 10, 64)
	if expireTime == 0 {
		expireTime = time.Now().Add(7 * 24 * time.Hour).Unix()
	}
	sound, vibration := pushmdl.SwitchOn, pushmdl.SwitchOn
	if params.Get("sound") != "" {
		if sd, _ := strconv.Atoi(params.Get("sound")); sd == pushmdl.SwitchOff {
			sound = pushmdl.SwitchOff
		}
	}
	if params.Get("vibration") != "" {
		if vr, _ := strconv.Atoi(params.Get("vibration")); vr == pushmdl.SwitchOff {
			vibration = pushmdl.SwitchOff
		}
	}
	passThrough, _ := strconv.Atoi(params.Get("pass_through"))
	if passThrough != pushmdl.SwitchOn {
		passThrough = pushmdl.SwitchOff
	}
	info := &pushmdl.PushInfo{
		TaskID:      pushmdl.TempTaskID(),
		APPID:       appID,
		Title:       alertTitle,
		Summary:     alertBody,
		LinkType:    int8(linkType),
		LinkValue:   linkValue,
		ExpireTime:  xtime.Time(expireTime),
		PassThrough: passThrough,
		Sound:       sound,
		Vibration:   vibration,
	}
	c.JSON(nil, pushSrv.TestPushToken(c, info, platform, token))
}
