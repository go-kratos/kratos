package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/app/admin/main/videoup/model/music"
	"go-common/app/admin/main/videoup/model/oversea"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus/report"
)

// send to log service
func (s *Service) sendVideoLog(c context.Context, vp *archive.VideoParam, others string) (err error) {
	var (
		v *archive.Video
		a *archive.Archive
	)
	if vp.Cid != 0 {
		v, err = s.arc.VideoByCID(c, vp.Cid)
	} else if vp.ID != 0 {
		v, err = s.arc.NewVideoByID(c, vp.ID)
	}
	if err != nil || v == nil {
		v = &archive.Video{} // ignore err
	}

	a, err = s.arc.Archive(c, vp.Aid)
	if err != nil || a == nil {
		a = &archive.Archive{} // ignore err
	}
	// send
	logData := &report.ManagerInfo{
		Uname:    vp.Oname,
		UID:      vp.UID,
		Business: archive.LogClientVideo,
		Type:     archive.LogClientTypeVideo,
		Oid:      vp.Cid,
		Action:   strconv.Itoa(int(vp.Status)),
		Ctime:    time.Now(),
		Index:    []interface{}{int64(vp.Attribute), v.CTime.Time().Unix(), vp.TagID, a.Title, vp.Note},
		Content: map[string]interface{}{
			"content": vp,
			"others":  others,
		},
	}
	report.Manager(logData)
	return
}

// send to log service
func (s *Service) sendArchiveLog(c context.Context, ap *archive.ArcParam, diff []string, a *archive.Archive) (err error) {
	// fmt
	ap.CTime = a.CTime
	if ap.Title == "" {
		ap.Title = a.Title
	}
	if ap.Attrs == nil {
		ap.Attrs = &archive.AttrParam{}
	}
	diffStr := strings.Join(diff, "\n")
	// log
	logData := &report.ManagerInfo{
		Uname:    ap.UName,
		UID:      ap.UID,
		Business: archive.LogClientArchive,
		Type:     archive.LogClientTypeArchive,
		Oid:      ap.Aid,
		Action:   strconv.Itoa(int(ap.State) + int(ap.Access)),
		Ctime:    time.Now(),
		Index:    []interface{}{a.Attribute, ap.CTime.Time().Unix(), ap.ReasonID, ap.Title, ap.Note},
		Content: map[string]interface{}{
			"content": ap,
			"diff":    diffStr,
		},
	}
	report.Manager(logData)
	extra, _ := json.Marshal(logData.Content)
	log.Info("sendArchiveLog json.Marshal(%s) logData(%+v) ap(%+v)", extra, logData, ap)
	return
}

//SendMusicLog send to log archive music
func (s *Service) SendMusicLog(c *bm.Context, clientType int, ap *music.LogParam) (err error) {
	if s.c.Env == "dev" {
		return
	}
	logData := &report.ManagerInfo{
		Uname:    ap.UName,
		UID:      ap.UID,
		Business: archive.LogClientArchiveMusic,
		Type:     clientType,
		Oid:      ap.ID,
		Action:   ap.Action,
		Ctime:    time.Now(),
		Index:    []interface{}{ap.ID},
		Content: map[string]interface{}{
			"object": ap,
		},
	}
	log.Info("sendMusicLog logData(%+v) ap(%+v)", logData, ap)
	report.Manager(logData)
	return
}

// sendPorderLog send porder modify log
func (s *Service) sendPorderLog(c context.Context, ap *archive.ArcParam, diff []string, porder *archive.Porder, a *archive.Archive) (err error) {
	if a.AttrVal(archive.AttrBitIsPorder) != 1 && ap.Attrs.IsPorder != 1 {
		log.Info("sendPorderLog ignore archive.is_porder(%d) ap.is_porder(%d) aid(%d)", a.AttrVal(archive.AttrBitIsPorder), ap.Attrs.IsPorder, a.Aid)
		return
	}
	// fmt
	var (
		oldP = map[string]interface{}{
			"is_porder":   a.AttrVal(archive.AttrBitIsPorder),
			"brand_id":    porder.BrandID,
			"brand_name":  porder.BrandName,
			"show_type":   porder.ShowType,
			"industry_id": porder.IndustryID,
			"official":    porder.Official,
			"allow_tag":   a.AttrVal(archive.AttrBitAllowTag),
		}
		newP = map[string]interface{}{
			"is_porder":   ap.Attrs.IsPorder,
			"brand_id":    ap.BrandID,
			"brand_name":  ap.BrandName,
			"show_type":   ap.ShowType,
			"industry_id": ap.IndustryID,
			"official":    ap.Official,
			"allow_tag":   ap.Attrs.AllowTag,
		}
	)
	ap.CTime = a.CTime
	if ap.Title == "" {
		ap.Title = a.Title
	}
	if ap.Attrs == nil {
		ap.Attrs = &archive.AttrParam{}
	}
	diffStr := strings.Join(diff, "\n")
	// log
	logData := &report.ManagerInfo{
		Uname:    ap.UName,
		UID:      ap.UID,
		Business: archive.LogClientPorder,
		Type:     archive.LogClientTypePorderLog,
		Oid:      ap.Aid,
		Action:   strconv.Itoa(int(ap.State) + int(ap.Access)),
		Ctime:    time.Now(),
		Index:    []interface{}{a.Attribute, ap.CTime.Time().Unix(), ap.ReasonID, ap.Title, ap.Note, ap.Porder.IndustryID, ap.Porder.Official, ap.Porder.GroupID},
		Content: map[string]interface{}{
			"content": ap,
			"diff":    diffStr,
			"old":     oldP,
			"new":     newP,
		},
	}
	report.Manager(logData)
	log.Info("sendPorderLog logData(%+v)", logData.Content)
	return
}

// sendConsumerLog send consumer log
func (s *Service) sendConsumerLog(c context.Context, cl *archive.ConsumerLog) (err error) {
	logData := &report.ManagerInfo{
		Uname:    cl.Uname,
		UID:      cl.UID,
		Business: archive.LogClientConsumer,
		Type:     archive.LogClientTypeConsumer,
		Oid:      cl.UID,
		Action:   strconv.Itoa(int(cl.Action)),
		Ctime:    time.Now(),
		Index:    []interface{}{cl.UID, cl.Action, cl.Ctime},
		Content: map[string]interface{}{
			"content": cl,
		},
	}
	report.Manager(logData)
	log.Info("sendConsumerLog logData(%+v)", cl)
	return
}

// sendPolicyLog send policy modify log
func (s *Service) sendPolicyLog(c context.Context, old, new *oversea.PolicyGroup) (err error) {
	var (
		action string
	)
	if new.ID == 0 {
		action = "add"
	} else if new.State == oversea.StateDeleted {
		action = "del"
	} else {
		action = "update"
	}
	// log
	logData := &report.ManagerInfo{
		Uname:    new.UserName,
		UID:      new.UID,
		Business: archive.LogClientPolicy,
		Type:     archive.LogClientTypePolicy,
		Oid:      new.ID,
		Action:   action,
		Ctime:    time.Now(),
		Index:    []interface{}{new.Type},
		Content: map[string]interface{}{
			"old": old,
			"new": new,
		},
	}
	report.Manager(logData)
	log.Info("sendPolicyLog logData(%+v)", logData)
	return
}
