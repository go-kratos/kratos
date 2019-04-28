package http

import (
	"strconv"
	"time"

	"go-common/app/admin/main/reply/conf"
	"go-common/app/admin/main/reply/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"
)

//panigate notice list
func listNotice(c *bm.Context) {
	var (
		err      error
		page     = int64(1)
		pageSize = int64(conf.Conf.Reply.PageSize)
	)

	params := c.Request.Form
	pageStr := params.Get("page")
	pageSizeStr := params.Get("pagesize")
	if pageStr != "" {
		if page, err = strconv.ParseInt(pageStr, 10, 64); err != nil {
			log.Warn("strconv.ParseInt(page:%s) error(%v)", pageStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if page < 1 {
			page = 1
		}
	}
	if pageSizeStr != "" {
		if pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64); err != nil || pageSize < 1 {
			log.Warn("strconv.ParseInt(pagesize:%s) error(%v)", pageSizeStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	data, total, err := rpSvc.ListNotice(c, page, pageSize)
	if err != nil {
		log.Error("svc.ListNotcie(%d,%d) error(%v)", page, pageSize, err)
		c.JSON(nil, err)
		return
	}
	res := map[string]interface{}{}
	res["data"] = data
	res["pager"] = map[string]interface{}{
		"num":   page,
		"size":  pageSize,
		"total": total,
	}
	c.JSONMap(res, nil)
	return
}

//panigate notice list
func listNotice2(c *bm.Context) {
	var (
		err      error
		page     int64 = 1
		pageSize       = int64(conf.Conf.Reply.PageSize)
	)

	params := c.Request.Form
	pageStr := params.Get("page")
	pageSizeStr := params.Get("pagesize")
	if pageStr != "" {
		if page, err = strconv.ParseInt(pageStr, 10, 64); err != nil {
			log.Warn("strconv.ParseInt(page:%s) error(%v)", pageStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if page < 1 {
			page = 1
		}
	}
	if pageSizeStr != "" {
		if pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64); err != nil || pageSize < 1 {
			log.Warn("strconv.ParseInt(pagesize:%s) error(%v)", pageSizeStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	data, total, err := rpSvc.ListNotice(c, page, pageSize)
	if err != nil {
		log.Error("svc.ListNotcie(%d,%d) error(%v)", page, pageSize, err)
		c.JSON(nil, err)
		return
	}
	res := map[string]interface{}{}
	res["data"] = data
	res["pager"] = model.Pager{Page: page, PageSize: pageSize, Total: total}
	c.JSONMap(res, nil)
	return
}

//return a notice detail
func getNotice(c *bm.Context) {
	var (
		err error
		id  uint64
	)

	params := c.Request.Form
	idStr := params.Get("id")
	if idStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if id, err = strconv.ParseUint(idStr, 10, 32); err != nil {
		log.Warn("strconv.ParseUint(id:%s) error(%v)", idStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := rpSvc.GetNotice(c, uint32(id))
	if err != nil {
		log.Error("svc.GetNotice(%d) error(%v)", id, err)
		c.JSON(nil, err)
		return
	}
	if data == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(data, nil)
	return
}

//update or create a notice
func editNotice(c *bm.Context) {
	var (
		err        error
		id         uint64
		plat       uint64
		condition  uint64
		version    string
		build      uint64
		title      string
		content    string
		link       string
		stime      int64
		etime      int64
		clientType string
	)

	params := c.Request.Form
	idStr := params.Get("id")
	platStr := params.Get("plat")
	conditionStr := params.Get("condi")
	version = params.Get("version")
	buildStr := params.Get("build")
	title = params.Get("title")
	content = params.Get("content")
	link = params.Get("link")
	stimeStr := params.Get("stime")
	etimeStr := params.Get("etime")
	clientType = params.Get("client_type")
	if platStr == "" || title == "" || content == "" || stimeStr == "" || etimeStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if idStr != "" {
		if id, err = strconv.ParseUint(idStr, 10, 32); err != nil {
			log.Warn("strconv.ParseUint(id:%s) error(%v)", idStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if plat, err = strconv.ParseUint(platStr, 10, 8); err != nil {
		log.Warn("strconv.ParseUint(plat:%s) error(%v)", platStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if conditionStr != "" {
		if condition, err = strconv.ParseUint(conditionStr, 10, 8); err != nil {
			log.Warn("strconv.ParseUint(condition:%s) error(%v)", conditionStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if buildStr != "" {
		if build, err = strconv.ParseUint(buildStr, 10, 32); err != nil {
			log.Warn("strconv.ParseUint(build:%s) error(%v)", buildStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	var tempTime time.Time
	if tempTime, err = time.Parse("2006-01-02 15:04:05", stimeStr); err != nil {
		//error,so try to parse as unix timestamp again
		stime, err = strconv.ParseInt(stimeStr, 10, 64)
		if err != nil {
			log.Warn("strconv.ParseUint(stime:%s) error(%v)", stimeStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		//标准库的time.parse默认以UTC为标准转换成unix时间戳，所以需要减去8小时
		stime = tempTime.Unix() - 8*3600
	}
	if tempTime, err = time.Parse("2006-01-02 15:04:05", etimeStr); err != nil {
		//error,so try to parse as unix timestamp again
		etime, err = strconv.ParseInt(etimeStr, 10, 64)
		if err != nil {
			log.Warn("strconv.ParseUint(etime:%s) error(%v)", etimeStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		//标准库的time.parse默认以UTC为标准转换成unix时间戳，所以需要减去8小时
		etime = tempTime.Unix() - 8*3600
	}
	//开始时间必须小于等于结束时间
	if stime > etime {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	nt := &model.Notice{
		ID:         uint32(id),
		Plat:       model.NoticePlat(plat),
		Version:    version,
		Condition:  model.NoticeCondition(condition),
		Build:      uint32(build),
		Title:      title,
		Content:    content,
		Link:       link,
		StartTime:  xtime.Time(stime),
		EndTime:    xtime.Time(etime),
		ClientType: clientType,
	}
	if idStr == "" {
		_, err = rpSvc.CreateNotice(c, nt)
		if err != nil {
			log.Error("svc.CreateNotice(%v) error(%v)", *nt, err)
			c.JSON(nil, err)
			return
		}
	} else {
		err = rpSvc.UpdateNotice(c, nt)
		if err != nil {
			log.Error("svc.UpdateNotice(%v) error(%v)", *nt, err)
			c.JSON(nil, err)
			return
		}
	}
	c.JSON(nil, nil)
	return
}

func deleteNotice(c *bm.Context) {
	var (
		err error
		id  uint64
	)

	params := c.Request.Form
	idStr := params.Get("id")
	if idStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if id, err = strconv.ParseUint(idStr, 10, 32); err != nil {
		log.Warn("strconv.ParseUint(id:%s) error(%v)", idStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = rpSvc.DeleteNotice(c, uint32(id))
	if err != nil {
		log.Error("svc.DeleteNotice(%d) error(%v)", id, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
	return
}

func offlineNotice(c *bm.Context) {

	var (
		err error
		id  uint64
	)

	params := c.Request.Form
	idStr := params.Get("id")
	if idStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if id, err = strconv.ParseUint(idStr, 10, 32); err != nil {
		log.Warn("strconv.ParseUint(id:%s) error(%v)", idStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = rpSvc.UpdateNoticeStatus(c, model.StatusOffline, uint32(id))
	if err != nil {
		log.Error("svc.UpdateNoticeStatus(%d) error(%v)", id, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
	return
}

func onlineNotice(c *bm.Context) {
	var (
		err error
		id  uint64
	)

	params := c.Request.Form
	idStr := params.Get("id")
	if idStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if id, err = strconv.ParseUint(idStr, 10, 32); err != nil {
		log.Warn("strconv.ParseUint(id:%s) error(%v)", idStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = rpSvc.UpdateNoticeStatus(c, model.StatusOnline, uint32(id))
	if err != nil {
		log.Error("svc.UpdateNoticeStatus(%d) error(%v)", id, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
	return
}
