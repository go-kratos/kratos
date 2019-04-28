package http

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/credit/model"
	"go-common/app/admin/main/credit/model/blocked"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	xtime "go-common/library/time"
	"go-common/library/xstr"
)

// infos  get info list.
func infos(c *bm.Context) {
	v := new(blocked.ArgBlockedSearch)
	err := c.Bind(v)
	if err != nil {
		return
	}
	info, pager, err := creSvc.Infos(c, v)
	if err != nil {
		log.Error("creSvc.Infos(%+v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	if len(info) == 0 {
		httpData(c, nil, pager)
		return
	}
	httpData(c, info, pager)
}

// infoByID  get info by id.
func infoByID(c *bm.Context) {
	var err error
	v := new(struct {
		ID int64 `form:"id"  validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	item := &blocked.Info{}
	if err = creSvc.ReadDB.Where("id = ?", v.ID).Find(item).Error; err != nil {
		if err != ecode.NothingFound {
			log.Error("creSvc.infoByID error(%v)", err)
			httpCode(c, err)
			return
		}
		httpData(c, nil, nil)
		return
	}
	if item != nil {
		info, err := creSvc.UserInfo(c, item.UID)
		if err != nil {
			log.Error("creSvc.UserInfo error(%v)", err)
			err = nil
		}
		if info != nil {
			item.UName = info.Info.Name
		}
		item.ReasonTypeDesc = blocked.ReasonTypeDesc(item.ReasonType)
		item.BlockedDaysDesc = blocked.BDaysDesc(item.BlockedDays, item.MoralNum, item.PunishType, item.BlockedForever)
		item.PublishStatusDesc = blocked.PStatusDesc[item.PublishStatus]
		item.OriginTypeDesc = blocked.OriginTypeDesc[item.OriginType]
		item.BlockedTypeDesc = blocked.BTypeDesc[item.BlockedType]
	}
	httpData(c, item, nil)
}

// upInfo update info.
func upInfo(c *bm.Context) {
	var err error
	v := new(blocked.ArgUpInfo)
	if err = c.Bind(v); err != nil {
		return
	}
	item := &blocked.Info{}
	if err = creSvc.ReadDB.Where("id = ?", v.ID).Find(item).Error; err != nil {
		log.Error("http.blcokedInfoByID(%d) error(%v)", v.ID, err)
		httpCode(c, err)
		return
	}
	var data map[string]interface{}
	switch {
	case item.PublishStatus == blocked.StatusClose && v.Status == blocked.StatusOpen:
		data = map[string]interface{}{
			"origin_content_modify": v.Content,
			"publish_status":        v.Status,
			"oper_id":               v.OID,
			"publish_time":          xtime.Time(time.Now().Unix()),
		}
	default:
		data = map[string]interface{}{
			"origin_content_modify": v.Content,
			"publish_status":        v.Status,
			"oper_id":               v.OID,
		}
	}
	if err = creSvc.DB.Model(&blocked.Info{}).Where("id = ?", v.ID).Updates(data).Error; err != nil {
		log.Error("s.dao.UpInfo(%v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	log.Info("business_upInfo:%+v", v)
	var multiple []interface{}
	single := map[string]interface{}{
		"id":             item.ID,
		"oper_id":        v.OID,
		"status":         item.Status,
		"publish_status": v.Status,
	}
	multiple = append(multiple, single)
	creSvc.AddNotify(func() {
		creSvc.Search.SearchUpdate(context.TODO(), blocked.BusinessBlockedInfo, blocked.TableBlockedInfo, multiple)
	})
	httpCode(c, nil)
}

func delInfo(c *bm.Context) {
	var err error
	v := new(struct {
		IDS    []int64 `form:"ids,split" validate:"min=1,max=100"`
		OID    int64   `form:"op_id" validate:"required"`
		Status int8    `form:"status" default:"1"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	items := []*blocked.Info{}
	if err = creSvc.ReadDB.Where("id IN (?)", v.IDS).Find(&items).Error; err != nil {
		log.Error("http.blcokedInfoByIDs(%s) error(%v)", xstr.JoinInts(v.IDS), err)
		httpCode(c, err)
		return
	}
	if err = creSvc.DB.Model(&blocked.Info{}).Where("id IN (?)", v.IDS).Updates(map[string]interface{}{"status": v.Status, "oper_id": v.OID}).Error; err != nil {
		log.Error("s.dao.delInfo(%v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	log.Info("business_delInfo:%+v", v)
	var multiple []interface{}
	for _, item := range items {
		single := map[string]interface{}{
			"id":             item.ID,
			"oper_id":        v.OID,
			"status":         v.Status,
			"publish_status": item.PublishStatus,
		}
		multiple = append(multiple, single)
	}
	creSvc.AddNotify(func() {
		creSvc.Search.SearchUpdate(context.TODO(), blocked.BusinessBlockedInfo, blocked.TableBlockedInfo, multiple)
	})
	httpCode(c, nil)
}

// upInfoStatus  update info publish_status.
func upInfoStatus(c *bm.Context) {
	var err error
	v := new(blocked.ArgUpStatus)
	if err = c.Bind(v); err != nil {
		return
	}
	items := []*blocked.Info{}
	if err = creSvc.ReadDB.Where("id IN (?)", v.IDS).Find(&items).Error; err != nil {
		log.Error("http.blcokedInfoByIDs(%s) error(%v)", xstr.JoinInts(v.IDS), err)
		httpCode(c, err)
		return
	}
	if err = creSvc.DB.Model(&blocked.Info{}).Where("id IN (?)", v.IDS).Updates(
		map[string]interface{}{
			"oper_id":        v.OID,
			"publish_status": v.Status,
			"publish_time":   xtime.Time(time.Now().Unix()),
		}).Error; err != nil {
		log.Error("creSvc.upInfoStatus(%v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	log.Info("business_upInfoStatus:%+v", v)
	var multiple []interface{}
	for _, item := range items {
		single := map[string]interface{}{
			"id":             item.ID,
			"oper_id":        v.OID,
			"status":         item.Status,
			"publish_status": v.Status,
		}
		multiple = append(multiple, single)
	}
	creSvc.AddNotify(func() {
		creSvc.Search.SearchUpdate(context.TODO(), blocked.BusinessBlockedInfo, blocked.TableBlockedInfo, multiple)
	})
	httpCode(c, nil)
}

// infosEx  export info data.
func infosEx(c *bm.Context) {
	v := new(blocked.ArgBlockedSearch)
	if err := c.Bind(v); err != nil {
		return
	}
	list, err := creSvc.InfosEx(c, v)
	if err != nil {
		log.Error("creSvc.InfosEx  error(%v)", err)
		httpCode(c, err)
		return
	}
	re, err := blocked.DealInfo(list)
	if err != nil {
		log.Error("blocked.DealInfo  error(%v)", err)
		httpCode(c, err)
		return
	}
	outBuf(c, creSvc.FormatCSV(re), fmt.Sprintf("%s-%s", time.Now().Format(model.TimeFormatDay), "blocked-info"))
}

func publishs(c *bm.Context) {
	v := new(blocked.ArgPublishSearch)
	if err := c.Bind(v); err != nil {
		return
	}
	publish, pager, err := creSvc.Publishs(c, v)
	if err != nil {
		log.Error("creSvc.Publishs(%+v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	if len(publish) == 0 {
		httpData(c, nil, pager)
		return
	}
	httpData(c, publish, pager)
}

// publishByID  get publish by id.
func publishByID(c *bm.Context) {
	var err error
	v := new(struct {
		ID int64 `form:"id" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	item := &blocked.Publish{}
	if err = creSvc.ReadDB.Where("id = ?", v.ID).Find(item).Error; err != nil {
		if err != ecode.NothingFound {
			log.Error("creSvc.publishByID error(%v)", err)
			httpCode(c, err)
			return
		}
		httpData(c, nil, nil)
		return
	}
	if item != nil {
		item.PublishTypeDesc = blocked.PTypeDesc[item.Type]
		item.PublishStatusDesc = blocked.PStatusDesc[item.PublishStatus]
		item.StickStatusDesc = blocked.SStatusDesc[item.StickStatus]
	}
	httpData(c, item, nil)
}

func addPublish(c *bm.Context) {
	var err error
	v := new(blocked.ArgPublish)
	if err = c.Bind(v); err != nil {
		return
	}
	publish := &blocked.Publish{
		Title:         v.Title,
		SubTitle:      v.SubTitle,
		PublishStatus: v.PublishStatus,
		StickStatus:   v.StickStatus,
		Content:       v.Content,
		URL:           v.URL,
		Type:          v.PType,
		OPID:          v.OID,
		ShowTime:      xtime.Time(time.Now().Unix()),
	}
	if v.ShowTime != "" {
		var tm time.Time
		tm, err = time.ParseInLocation(model.TimeFormatSec, v.ShowTime, time.Local)
		if err == nil {
			publish.ShowTime = xtime.Time(tm.Unix())
		}
	}
	if err = creSvc.DB.Create(publish).Error; err != nil {
		log.Error("creSvc.AddPublish(%+v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	log.Info("business_addPublish:%+v", v)
	httpCode(c, nil)
}

func upPublish(c *bm.Context) {
	var err error
	v := new(blocked.ArgPublish)
	if err = c.Bind(v); err != nil {
		return
	}
	item := &blocked.Publish{}
	if err = creSvc.ReadDB.Where("id = ?", v.ID).Find(item).Error; err != nil {
		log.Error("http.publishByID(%d) error(%v)", v.ID, err)
		httpCode(c, err)
		return
	}
	var showTime xtime.Time
	tm, err := time.ParseInLocation(model.TimeFormatSec, v.ShowTime, time.Local)
	if err == nil {
		showTime = xtime.Time(tm.Unix())
	}
	if err = creSvc.DB.Model(&blocked.Publish{}).Where("id = ?", v.ID).
		Updates(
			map[string]interface{}{
				"title":          v.Title,
				"sub_title":      v.SubTitle,
				"publish_status": v.PublishStatus,
				"stick_status":   v.StickStatus,
				"content":        v.Content,
				"url":            v.URL,
				"ptype":          v.PType,
				"show_time":      showTime,
				"oper_id":        v.OID,
			}).Error; err != nil {
		log.Error("creSvc.UpPublish(%+v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	log.Info("business_upPublish:%+v", v)
	var multiple []interface{}
	single := map[string]interface{}{
		"id":        v.ID,
		"title":     v.Title,
		"sub_title": v.SubTitle,
		"oper_id":   v.OID,
		"ptype":     v.PType,
		"show_time": showTime.Time().Format(model.TimeFormatSec),
		"status":    item.Status,
	}
	multiple = append(multiple, single)
	creSvc.AddNotify(func() {
		creSvc.Search.SearchUpdate(context.TODO(), blocked.BusinessBlockedPublish, blocked.TableBlockedPublish, multiple)
	})
	httpCode(c, nil)
}

func delPublish(c *bm.Context) {
	var err error
	v := new(struct {
		IDS    []int64 `form:"ids,split" validate:"required"`
		Status int8    `form:"status" default:"1"`
		OID    int64   `form:"op_id"  validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if err = creSvc.DB.Model(&blocked.Publish{}).Where(v.IDS).Updates(map[string]interface{}{"status": v.Status, "oper_id": v.OID}).Error; err != nil {
		httpCode(c, err)
		return
	}
	log.Info("business_delPublish:%+v", v)
	var multiple []interface{}
	for _, id := range v.IDS {
		single := map[string]interface{}{
			"id":      id,
			"oper_id": v.OID,
			"status":  v.Status,
		}
		multiple = append(multiple, single)
	}
	creSvc.AddNotify(func() {
		creSvc.Search.SearchUpdate(context.TODO(), blocked.BusinessBlockedPublish, blocked.TableBlockedPublish, multiple)
	})
	httpCode(c, nil)
}

func notices(c *bm.Context) {
	var (
		err   error
		order = "id"
		sort  = "DESC"
	)
	v := new(struct {
		PN int `form:"pn" default:"1"`
		PS int `form:"ps" default:"20"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	count := int(0)
	items := []*blocked.Notice{}
	pager := &blocked.Pager{
		Total: count,
		PN:    v.PN,
		PS:    v.PS,
		Order: order,
		Sort:  sort,
	}
	if err = creSvc.ReadDB.Offset((v.PN - 1) * v.PS).Limit(v.PS).Order(fmt.Sprintf("%s %s", order, sort)).Find(&items).Error; err != nil {
		if err != ecode.NothingFound {
			log.Error("creSvc.notices error(%v)", err)
			httpData(c, nil, nil)
			return
		}
		httpData(c, nil, pager)
		return
	}
	if err = creSvc.ReadDB.Model(&blocked.Notice{}).Count(&count).Error; err != nil {
		if err != ecode.NothingFound {
			log.Error("creSvc.notices count error(%v)", err)
			httpData(c, nil, nil)
			return
		}
		httpData(c, nil, pager)
		return
	}
	pager.Total = count
	for _, v := range items {
		v.StatusDesc = blocked.NoticeStateDesc[v.Status]
		v.OPName = creSvc.Managers[v.OperID]
	}
	httpData(c, items, pager)
}

func addNotice(c *bm.Context) {
	var err error
	v := new(struct {
		Content string `form:"content" validate:"required"`
		URL     string `form:"url"  validate:"required"`
		OID     int64  `form:"op_id"  validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if err = creSvc.DB.Model(&blocked.Notice{}).Create(&blocked.Notice{Content: v.Content, URL: v.URL, OperID: v.OID}).Error; err != nil {
		log.Error("creSvc.addNotice(%v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	log.Info("business_addNotice:%+v", v)
	httpCode(c, nil)
}

func upNoticeStatus(c *bm.Context) {
	var err error
	v := new(struct {
		ID     int  `form:"id" validate:"required"`
		Status int8 `form:"status" validate:"min=0,max=1"`
		OID    int  `form:"op_id" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if err = creSvc.DB.Model(&blocked.Notice{}).Where("id=?", v.ID).Updates(map[string]interface{}{"status": v.Status, "oper_id": v.OID}).Error; err != nil {
		log.Error("creSvc.upNoticeStatus(%v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	log.Info("business_upNoticeStatus:%+v", v)
	httpCode(c, nil)
}
