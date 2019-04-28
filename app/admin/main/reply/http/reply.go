package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/reply/conf"
	"go-common/app/admin/main/reply/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/net/ip"
)

func replyExport(c *bm.Context) {
	var (
		err       error
		startTime time.Time
		endTime   time.Time
	)
	v := new(struct {
		Oid       int64   `form:"oid" validate:"required"`
		Type      int8    `form:"type" validate:"required"`
		StartTime string  `form:"start_time"`
		EndTime   string  `form:"end_time"`
		States    []int64 `form:"states,split"`
		Mid       int64   `form:"uid"`
	})
	err = c.Bind(v)
	if err != nil {
		return
	}
	if v.StartTime == "" {
		startTime = time.Now().AddDate(0, -3, 0)
		v.StartTime = startTime.Format("2006-01-02")
	} else {
		if startTime, err = time.Parse("2006-01-02", v.StartTime); err != nil {
			log.Warn("time.Parse(%s) error(%v)", v.StartTime, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if v.EndTime == "" {
		endTime = time.Now()
		v.EndTime = endTime.Format("2006-01-02")
	} else {
		if endTime, err = time.Parse("2006-01-02", v.EndTime); err != nil {
			log.Warn("time.Parse(%s) error(%v)", v.EndTime, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	states := strings.Trim(strings.Join(strings.Split(fmt.Sprint(v.States), " "), ","), "[]")
	data, err := rpSvc.ExportReply(c, v.Oid, v.Mid, v.Type, states, startTime, endTime)
	if err != nil {
		log.Error("Svc.ExportReply(%d,%d,%d,%s,%v,%v) error(%v)", v.Oid, v.Mid, v.Type, states, startTime, endTime, err)
		c.Render(http.StatusInternalServerError, model.CSV{
			Data:  nil,
			Title: "",
		})
		return
	}
	c.Render(http.StatusOK, model.CSV{
		Data:  data,
		Title: v.StartTime + "至" + v.EndTime,
	})
	return
}

func replySearch(c *bm.Context) {
	var (
		err                     error
		typ, startTime, endTime int64
		page                    = int64(conf.Conf.Reply.PageNum)
		pageSize                = int64(conf.Conf.Reply.PageSize)
		sp                      = &model.SearchParams{}
	)
	params := c.Request.Form
	typeStr := params.Get("type")
	oidStr := params.Get("oid")
	uidStr := params.Get("uid")
	adminIDStr := params.Get("adminid")
	startTimeStr := params.Get("start_time")
	endTimeStr := params.Get("end_time")
	ipStr := params.Get("ip")
	pageStr := params.Get("page")
	pageSizeStr := params.Get("pagesize")
	// parse params
	if typ, err = strconv.ParseInt(typeStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(type:%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if oidStr != "" {
		if sp.Oid, err = strconv.ParseInt(oidStr, 10, 64); err != nil {
			log.Warn("strconv.ParseInt(oid:%s) error(%v)", oidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if uidStr != "" {
		if sp.UID, err = strconv.ParseInt(uidStr, 10, 64); err != nil {
			log.Warn("strconv.ParseInt(uid:%s) error(%v)", uidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if adminIDStr != "" {
		if sp.AdminID, err = strconv.ParseInt(adminIDStr, 10, 64); err != nil {
			log.Warn("strconv.ParseInt(adminid:%s) error(%v)", adminIDStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if startTimeStr != "" {
		if startTime, err = strconv.ParseInt(startTimeStr, 10, 64); err == nil {
			sp.Begin = time.Unix(startTime, 0)
		} else {
			var t time.Time
			t, err = time.Parse("2006-01-02", startTimeStr)
			if err == nil {
				sp.Begin = t
			} else {
				c.JSON(nil, ecode.RequestErr)
				return
			}
		}
	}
	if endTimeStr != "" {
		if endTime, err = strconv.ParseInt(endTimeStr, 10, 64); err == nil {
			sp.End = time.Unix(endTime, 0)
		} else {
			var t time.Time
			t, err = time.Parse("2006-01-02", endTimeStr)
			if err == nil {
				sp.End = t
			} else {
				c.JSON(nil, ecode.RequestErr)
				return
			}
		}
	} else if startTimeStr != "" {
		sp.End = time.Now()
	}
	if ipStr != "" {
		if sp.IP, err = strconv.ParseInt(ipStr, 10, 64); err != nil {
			sp.IP = int64(ip.InetAtoN(ipStr))
		}
	}
	if pageStr != "" {
		if page, err = strconv.ParseInt(pageStr, 10, 64); err != nil || page < 1 {
			log.Warn("strconv.ParseInt(page:%s) error(%v)", pageStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if pageSizeStr != "" {
		if pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64); err != nil || pageSize < int64(1) || pageSize > int64(conf.Conf.Reply.PageSize) {
			log.Warn("strconv.ParseInt(pagesize:%s) error(%v)", pageSizeStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	sp.Type = int32(typ)
	sp.TypeIds = params.Get("typeids")
	sp.Keyword = params.Get("keyword")
	sp.KeywordHigh = params.Get("keyword_high")
	sp.Uname = params.Get("nickname")
	sp.States = params.Get("states")
	sp.Sort = params.Get("sort")
	sp.Order = params.Get("order")
	sp.Attr = params.Get("attr")
	sp.AdminName = params.Get("adminname")
	sp.Sort, sp.Order = sp.Order, sp.Sort
	if sp.Sort == "" {
		sp.Sort = "desc"
	}
	if sp.Order == "" {
		sp.Order = "ctime"
	}
	rs, err := rpSvc.ReplySearch(c, sp, page, pageSize)
	if err != nil {
		log.Error("svc.ReplySearch(%v,%d,%d) error(%v)", sp, page, pageSize, err)
		c.JSON(nil, err)
		return
	}
	res := map[string]interface{}{}
	res["data"] = rs
	c.JSONMap(res, nil)
	return
}

func updateReplyConfig(c *bm.Context) {
	var (
		err                                     error
		typ, oid, showEntry, showAdmin, adminID int64
		operator                                string
		config                                  = &model.Config{}
		bs                                      []byte
	)
	params := c.Request.Form
	typeStr := params.Get("type")
	oidStr := params.Get("oid")
	adminIDStr := params.Get("adminid")
	operator = params.Get("operator")
	showEntityStr := params.Get("showentry")
	showAdminStr := params.Get("showadmin")
	if oid, err = strconv.ParseInt(oidStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(oid:%s) error(%v)", oidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if typ, err = strconv.ParseInt(typeStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(type:%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if showEntry, err = strconv.ParseInt(showEntityStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(showEntry:%s) error(%v)", showEntityStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if showAdmin, err = strconv.ParseInt(showAdminStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(showAdmin:%s) error(%v)", showAdminStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	adminID, _ = strconv.ParseInt(adminIDStr, 10, 64)
	if uid, ok := c.Get("uid"); ok {
		adminID = uid.(int64)
	}
	config.Oid = oid
	config.Type = int32(typ)
	config.Category = 1
	config.AdminID = adminID
	config.Operator = operator
	configValue := map[string]int64{
		"showentry": showEntry,
		"showadmin": showAdmin,
	}
	config.ShowEntry = int32(showEntry)
	config.ShowAdmin = int32(showAdmin)
	if bs, err = json.Marshal(configValue); err == nil {
		config.Config = string(bs)
	}
	if _, err = rpSvc.AddReplyConfig(c, config); err != nil {
		log.Error("svc.AddReplyConfig(%v) error(%v)", config, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func renewReplyConfig(c *bm.Context) {
	var (
		err error
		id  int64
	)

	params := c.Request.Form
	idStr := params.Get("id")
	if id, err = strconv.ParseInt(idStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(type:%s) error(%v)", idStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if _, err = rpSvc.RenewReplyConfig(c, id); err != nil {
		log.Error("svc.RenewReplyConfig(%d) error(%v)", id, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func paginateReplyConfig(c *bm.Context) {
	var (
		err                error
		oid, typ, category int64
		page               = int64(conf.Conf.Reply.PageNum)
		pageSize           = int64(conf.Conf.Reply.PageSize)
	)

	params := c.Request.Form
	typeStr := params.Get("type")
	categoryStr := params.Get("category")
	oidStr := params.Get("oid")
	operator := params.Get("operator")
	pageStr := params.Get("pn")
	pageSizeStr := params.Get("ps")
	if oid, err = strconv.ParseInt(oidStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(oid:%s) error(%v)", oidStr, err)
	}
	if typ, err = strconv.ParseInt(typeStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(type:%s) error(%v)", typeStr, err)
	}
	if category, err = strconv.ParseInt(categoryStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(category:%s) error(%v)", categoryStr, err)
	}
	if pageStr != "" {
		if page, err = strconv.ParseInt(pageStr, 10, 64); err != nil || page < 1 {
			log.Warn("strconv.ParseInt(page:%s) error(%v)", pageStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if pageSizeStr != "" {
		if pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64); err != nil || pageSize < 1 || pageSize > int64(conf.Conf.Reply.PageSize) {
			log.Warn("strconv.ParseInt(pagesize:%s) error(%v)", pageSizeStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	configs, totalCount, pages, err := rpSvc.PaginateReplyConfig(c, int32(typ), int32(category), oid, operator, int(page-1)*int(pageSize), int(pageSize))
	if err != nil {
		log.Error("svc.PaginateReplyConfig(%d,%d,%d,%v,%d,%d) error(%v)", int32(typ), 1, oid, operator, int(page), int(pageSize), err)
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{
		"page":      page,
		"pagesize":  pageSize,
		"pagecount": pages,
		"total":     totalCount,
		"order":     "default",
		"result":    configs,
	}
	c.JSONMap(data, nil)
	return
}

func loadReplyConfig(c *bm.Context) {
	var (
		err                error
		oid, typ, category int64
	)

	params := c.Request.Form
	typeStr := params.Get("type")
	categoryStr := params.Get("category")
	oidStr := params.Get("oid")
	if oid, err = strconv.ParseInt(oidStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(oid:%s) error(%v)", oidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if typ, err = strconv.ParseInt(typeStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(type:%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if category, err = strconv.ParseInt(categoryStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(category:%s) error(%v)", categoryStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	config, err := rpSvc.LoadReplyConfig(c, int32(typ), int32(category), oid)
	if err != nil {
		log.Error("Svc.LoadReplyConfig(%d,%d,%d) error(%v)", int32(typ), 1, oid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(config, nil)
	return
}

// adminPassReply pass reply normal.
func adminPassReply(c *bm.Context) {
	v := new(struct {
		RpID    []int64 `form:"rpid,split" validate:"required"`
		Oid     []int64 `form:"oid,split" validate:"required"`
		Type    int32   `form:"type" validate:"required"`
		AdminID int64   `form:"adid"`
		Remark  string  `form:"remark"`
	})
	var err error
	err = c.Bind(v)
	if err != nil {
		return
	}
	if len(v.Oid) != len(v.RpID) {
		log.Warn("len(v.Oid) != len(v.RpID),%d != %d", len(v.Oid), len(v.RpID))
		c.JSON(nil, ecode.RequestErr)
	}
	adid := v.AdminID
	if uid, ok := c.Get("uid"); ok {
		adid = uid.(int64)
	}
	var adName string
	if uname, ok := c.Get("username"); ok {
		adName = uname.(string)
	}
	err = rpSvc.AdminPassReply(c, adid, adName, v.Oid, v.RpID, v.Type, v.Remark)
	c.JSON(nil, err)
	return
}

// adminRecoverReply recover reply normal.
func adminRecoverReply(c *bm.Context) {
	v := new(struct {
		RpID    int64  `form:"rpid" validate:"required"`
		Oid     int64  `form:"oid" validate:"required"`
		Type    int32  `form:"type" validate:"required"`
		AdminID int64  `form:"adid"`
		Remark  string `form:"remark"`
	})
	var err error
	err = c.Bind(v)
	if err != nil {
		return
	}
	adid := v.AdminID
	if uid, ok := c.Get("uid"); ok {
		adid = uid.(int64)
	}
	var adName string
	if uname, ok := c.Get("username"); ok {
		adName = uname.(string)
	}
	err = rpSvc.AdminRecoverReply(c, adid, adName, v.Oid, v.RpID, v.Type, v.Remark)
	c.JSON(nil, err)
	return
}

// adminEditReply edit reply normal.
func adminEditReply(c *bm.Context) {
	v := new(struct {
		RpID    int64  `form:"rpid" validate:"required"`
		Oid     int64  `form:"oid" validate:"required"`
		Type    int32  `form:"type" validate:"required"`
		AdminID int64  `form:"adid"`
		Message string `form:"message"`
		Remark  string `form:"remark"`
	})
	var err error
	err = c.Bind(v)
	if err != nil {
		return
	}
	adid := v.AdminID
	if uid, ok := c.Get("uid"); ok {
		adid = uid.(int64)
	}
	var adName string
	if uname, ok := c.Get("username"); ok {
		adName = uname.(string)
	}
	err = rpSvc.AdminEditReply(c, adid, adName, v.Oid, v.RpID, v.Type, v.Message, v.Remark)
	c.JSON(nil, err)
	return
}

// Compose compose struct
type Compose struct {
	RpIDs []int64
	Oids  []int64
}

// adminDelReply del reply normal.
func adminDelReply(c *bm.Context) {
	v := new(struct {
		RpID    []int64 `form:"rpid,split" validate:"required"`
		Oid     []int64 `form:"oid,split" validate:"required"`
		Type    []int32 `form:"type,split" validate:"required"`
		AdminID int64   `form:"adid"`
		Remark  string  `form:"remark"`
		Moral   int32   `form:"moral"`
		Notify  bool    `form:"notify"`
		FTime   int64   `form:"ftime"`
		AdName  string  `form:"adname"`
		Reason  int32   `form:"reason"`
		FReason int32   `form:"freason"`
	})
	var err error
	err = c.Bind(v)
	if err != nil {
		return
	}
	if len(v.RpID) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(v.RpID) != len(v.Oid) || len(v.RpID) != len(v.Type) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tMap := make(map[int32]*Compose)
	for i, tp := range v.Type {
		if c, ok := tMap[tp]; ok {
			c.Oids = append(c.Oids, v.Oid[i])
			c.RpIDs = append(c.RpIDs, v.RpID[i])
		} else {
			c = &Compose{
				Oids:  []int64{v.Oid[i]},
				RpIDs: []int64{v.RpID[i]},
			}
			tMap[tp] = c
		}
	}
	if uid, ok := c.Get("uid"); ok {
		v.AdminID = uid.(int64)
	}
	if uname, ok := c.Get("username"); ok {
		v.AdName = uname.(string)
	}
	for tp, com := range tMap {
		err = rpSvc.AdminDeleteReply(c, v.AdminID, com.Oids, com.RpIDs, v.FTime, tp, v.Moral, v.Notify, v.AdName, v.Remark, v.Reason, v.FReason)
		if !ecode.NothingFound.Equal(err) {
			err = nil
		}
	}
	c.JSON(nil, err)
	return
}

// callbackDelReply del reply normal.
func callbackDelReply(c *bm.Context) {
	v := new(struct {
		Actor struct {
			AdminID   int64  `json:"admin_id"`
			AdminName string `json:"admin_name"`
		} `json:"actor"`
		CTime  int64 `json:"ctime"`
		Object struct {
			TagID         int32  `json:"tag_id"`
			DecreaseMoral int32  `json:"decrease_moral"`
			BlockReasonID int32  `json:"block_reason_id"`
			BlockDay      int64  `json:"block_day"`
			Reason        string `json:"reason"`
			DisposeMode   int32  `json:"dispose_mode"`
			State         int32  `json:"state"`
			RID           int32  `json:"rid"`
		} `json:"object"`
		Target []struct {
			OID   int64 `json:"oid"`
			FID   int32 `json:"fid"`
			EID   int64 `json:"eid"`
			State int32 `json:"state"`
		} `json:"targets"`
	})
	var err error
	err = c.BindWith(v, binding.JSON)
	if err != nil {
		return
	}
	log.Info("callbackDelReply V:%+v", *v)
	if v.Object.State == 2 {
		for _, t := range v.Target {
			rpSvc.ReportIgnore(c, []int64{t.OID}, []int64{t.EID}, v.Actor.AdminID, v.Actor.AdminName, t.FID, v.Object.RID, v.Object.Reason, false)
		}
	} else if v.Object.State == 1 {
		for _, t := range v.Target {
			err = rpSvc.CallbackDeleteReply(c, v.Actor.AdminID, t.OID, t.EID, v.Object.BlockDay, t.FID, v.Object.DecreaseMoral, v.Actor.AdminName, v.Object.Reason, v.Object.TagID, v.Object.BlockReasonID)
		}
	} else if v.Object.State == 10 {
		for _, t := range v.Target {
			rpSvc.ReportStateSet(c, []int64{t.OID}, []int64{t.EID}, v.Actor.AdminID, v.Actor.AdminName, t.FID, 8, v.Object.Reason, false)
		}
	} else {
		log.Warn("callbackDelReply RequestErr(%d,%d)", v.Object.State, v.Object.DisposeMode)
	}

	c.JSON(nil, err)
	return
}

// adminTopReply top reply normal.
func adminTopReply(c *bm.Context) {
	v := new(struct {
		RpID    int64  `form:"rpid" validate:"required"`
		Oid     int64  `form:"oid" validate:"required"`
		Type    int32  `form:"type" validate:"required"`
		AdminID int64  `form:"adid"`
		Action  uint32 `form:"action"`
	})
	var err error
	err = c.Bind(v)
	if err != nil {
		return
	}
	adid := v.AdminID
	if uid, ok := c.Get("uid"); ok {
		adid = uid.(int64)
	}
	var adName string
	if uname, ok := c.Get("username"); ok {
		adName = uname.(string)
	}
	err = rpSvc.AddTop(c, adid, adName, v.Oid, v.RpID, v.Type, v.Action)
	c.JSON(nil, err)
	return
}

// adminTopReplyLog get reply top log.
func adminTopReplyLog(c *bm.Context) {
	v := new(struct {
		Oid       int64  `form:"oid" `
		Type      int32  `form:"type" `
		StartTime string `form:"start_time"`
		EndTime   string `form:"end_time"`
		Page      int64  `form:"pn"`
		PageSize  int64  `form:"ps"`
		Order     string `form:"order"`
		Sort      string `form:"sort"`
	})
	var err error
	err = c.Bind(v)
	if err != nil {
		return
	}

	sp := model.LogSearchParam{
		Oid:       v.Oid,
		Type:      v.Type,
		CtimeFrom: v.StartTime,
		CtimeTo:   v.EndTime,
		Pn:        v.Page,
		Ps:        v.PageSize,
		Order:     v.Order,
		Sort:      v.Sort,
	}
	data, err := rpSvc.ReplyTopLog(c, sp)
	res := map[string]interface{}{}
	res["data"] = data
	c.JSONMap(res, err)

	return
}

// adminMarkAsSpam make a reply as spam
func adminMarkAsSpam(c *bm.Context) {
	v := new(struct {
		RpIDs   []int64 `form:"rpid,split" validate:"required"`
		Oids    []int64 `form:"oid,split" validate:"required"`
		AdminID int64   `form:"adid"`
		Remark  string  `form:"remark"`
	})
	var err error
	err = c.Bind(v)
	if err != nil {
		return
	}
	adid := v.AdminID
	if uid, ok := c.Get("uid"); ok {
		adid = uid.(int64)
	}
	var adName string
	if uname, ok := c.Get("username"); ok {
		adName = uname.(string)
	}
	err = rpSvc.MarkAsSpam(c, v.Oids, v.RpIDs, adid, adName, v.Remark)
	c.JSON(nil, err)
	return
}

// adminMarkAsSpam make a reply as spam
func adminReplyList(c *bm.Context) {
	v := new(struct {
		RpIDs []int64 `form:"eids,split" validate:"required"`
		Oids  []int64 `form:"oids,split" validate:"required"`
		IDs   []int64 `form:"ids,split" validate:"required"`
	})
	var err error
	err = c.Bind(v)
	if err != nil {
		return
	}
	res := make(map[int64]*model.ReplyEx)

	if len(v.RpIDs) == len(v.IDs) && len(v.IDs) == len(v.Oids) && len(v.RpIDs) > 0 {
		for i := range v.IDs {
			ex := &model.ReplyEx{}
			ex.ID = v.RpIDs[i]
			ex.Oid = v.Oids[i]
			res[v.IDs[i]] = ex
		}
		var rpMap map[int64]*model.ReplyEx
		rpMap, err = rpSvc.ReplyList(c, v.Oids, v.RpIDs)
		for k, data := range res {
			res[k] = rpMap[data.ID]
		}
	} else {
		log.Error("adminReplyList request invalid(%v)", v)
		err = ecode.RequestErr
	}

	c.JSON(res, err)
	return
}

func topChildReply(c *bm.Context) {
	v := new(struct {
		Root  int64 `form:"root" validate:"required"`
		Child int64 `form:"child" validate:"required"`
		Oid   int64 `form:"oid" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, rpSvc.TopChildReply(c, v.Root, v.Child, v.Oid))
}

func foldReply(c *bm.Context) {
	v := new(struct {
		Oids  []int64 `form:"oids,split" validate:"required"`
		RpIDs []int64 `form:"rpids,split" validate:"required"`
		Types []int64 `form:"types,split" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if len(v.Oids) != len(v.RpIDs) || len(v.Oids) != len(v.Types) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, rpSvc.FoldReplies(c, v.Oids, v.Types, v.RpIDs))
}
