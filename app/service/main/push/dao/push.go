package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unicode"

	"go-common/app/service/main/push/dao/apns2"
	"go-common/app/service/main/push/dao/fcm"
	"go-common/app/service/main/push/dao/huawei"
	"go-common/app/service/main/push/dao/jpush"
	"go-common/app/service/main/push/dao/mi"
	"go-common/app/service/main/push/dao/oppo"
	"go-common/app/service/main/push/model"
	"go-common/library/log"
)

func fmtRoundIndex(appid int64, platform int) string {
	return fmt.Sprintf("%d_%d", appid, platform)
}

func (d *Dao) roundIndex(appid int64, platform int) (int, error) {
	i := fmtRoundIndex(appid, platform)
	l := d.clientsLen[i]
	if l == 0 {
		log.Error("no client app(%d) platform(%d)", appid, platform)
		PromError("push:no client")
		return 0, errNoClinets
	}
	n := atomic.AddUint32(d.clientsIndex[i], 1)
	if n%uint32(l) == 0 { // 把第一个client预留出来，做一些其它类型请求工作
		n = atomic.AddUint32(d.clientsIndex[i], 1)
	}
	return int(n % uint32(l)), nil
}

func logPushError(task string, platform int, tokens []string) {
	for _, t := range tokens {
		log.Error("push error, task(%s) platfrom(%d) token(%s)", task, platform, t)
	}
}

func buildAPNS(info *model.PushInfo, item *model.PushItem) *apns2.Payload {
	var aps apns2.Aps
	if info.PassThrough == model.SwitchOn {
		aps = apns2.Aps{
			ContentAvailable: 1, // 必带字段，让程序处于后台时也可以获取到推送内容
		}
	} else {
		aps = apns2.Aps{
			Alert: apns2.Alert{
				Title: info.Title,
				Body:  info.Summary,
			},
			Badge:          0,
			MutableContent: 1,
		}
		if info.Sound == model.SwitchOn {
			aps.Sound = "default" // 默认提示音
		}
	}
	scheme := model.Scheme(info.LinkType, info.LinkValue, item.Platform, item.Build)
	return &apns2.Payload{Aps: aps, URL: scheme, TaskID: info.TaskID, Token: item.Token, Image: info.ImageURL}
}

// PushIPhone .
func (d *Dao) PushIPhone(c context.Context, info *model.PushInfo, item *model.PushItem) (res *model.HTTPResponse, err error) {
	var (
		index    int
		response *apns2.Response
	)
	if index, err = d.roundIndex(info.APPID, item.Platform); err != nil {
		return
	}
	if response, err = d.clientsIPhone[info.APPID][index].Push(item.Token, buildAPNS(info, item), int64(info.ExpireTime)); err != nil {
		log.Error("push iPhone task(%s) mid(%d) token(%s) error", info.TaskID, item.Mid, item.Token)
		PromError("push: 推送iPhone")
		return
	}
	if response == nil {
		return
	}
	res = &model.HTTPResponse{Code: response.StatusCode, Msg: response.Reason}
	log.Info("push iPhone task(%s) mid(%d) token(%s) success", info.TaskID, item.Mid, item.Token)
	return
}

// PushIPad .
func (d *Dao) PushIPad(c context.Context, info *model.PushInfo, item *model.PushItem) (res *model.HTTPResponse, err error) {
	var (
		index    int
		response *apns2.Response
	)
	if index, err = d.roundIndex(info.APPID, item.Platform); err != nil {
		return
	}
	if response, err = d.clientsIPad[info.APPID][index].Push(item.Token, buildAPNS(info, item), int64(info.ExpireTime)); err != nil {
		log.Error("push iPad task(%s) mid(%d) token(%s) error", info.TaskID, item.Mid, item.Token)
		PromError("push:推送iPad")
		return
	}
	if response == nil {
		return
	}
	res = &model.HTTPResponse{Code: response.StatusCode, Msg: response.Reason}
	log.Info("push iPad task(%s) mid(%d) token(%s) success", info.TaskID, item.Mid, item.Token)
	return
}

// PushMi .
func (d *Dao) PushMi(c context.Context, info *model.PushInfo, scheme, tokens string) (res *model.HTTPResponse, err error) {
	res = &model.HTTPResponse{}
	var index int
	if index, err = d.roundIndex(info.APPID, model.PlatformXiaomi); err != nil {
		return
	}
	passThrough := mi.NotPassThrough
	if info.PassThrough == model.SwitchOn {
		passThrough = mi.PassThrough
	}
	xmm := &mi.XMMessage{
		Payload:               scheme,
		RestrictedPackageName: d.clientsMi[info.APPID][0].Package,
		PassThrough:           passThrough,
		Title:                 info.Title,
		Description:           info.Summary,
		NotifyType:            mi.NotifyTypeDefaultNone,
		TaskID:                info.TaskID,
	}
	xmm.SetRegID(tokens)
	xmm.SetNotifyID(info.TaskID)
	xmm.SetTimeToLive(int64(info.ExpireTime))
	xmm.SetCallbackParam(strconv.FormatInt(info.APPID, 10))
	if info.Sound == model.SwitchOn {
		xmm.SetNotifyType(mi.NotifyTypeDefaultSound)
	}
	if info.Vibration == model.SwitchOn {
		xmm.SetNotifyType(mi.NotifyTypeDefaultVibration)
	}
	if info.Sound == model.SwitchOn && info.Vibration == model.SwitchOn {
		xmm.SetNotifyType(mi.NotifyTypeDefaultAll)
	}
	var response *mi.Response
	if response, err = d.clientsMi[info.APPID][index].Push(xmm); err != nil {
		log.Error("push mi task(%s) resp(%+v) error(%v)", info.TaskID, response, err)
		logPushError(info.TaskID, model.PlatformXiaomi, strings.Split(tokens, ","))
		PromError("push:推送Xiaomi")
		return
	}
	res.Code = response.Code
	res.Msg = response.Reason
	if response.Code == mi.ResultCodeOk {
		res.Code = model.HTTPCodeOk
		res.Msg = response.Info
	}
	log.Info("push mi task(%s) tokens(%d) result(%+v) traceid(%s) success", info.TaskID, len(tokens), res, response.TraceID)
	return
}

// PushMiByMids .
func (d *Dao) PushMiByMids(c context.Context, info *model.PushInfo, scheme, mids string) (res *model.HTTPResponse, err error) {
	if d.clientMiByMids[info.APPID] == nil {
		return
	}
	res = &model.HTTPResponse{}
	passThrough := mi.NotPassThrough
	if info.PassThrough == model.SwitchOn {
		passThrough = mi.PassThrough
	}
	xmm := &mi.XMMessage{
		Payload:               scheme,
		RestrictedPackageName: d.clientMiByMids[info.APPID].Package,
		PassThrough:           passThrough,
		Title:                 info.Title,
		Description:           info.Summary,
		NotifyType:            mi.NotifyTypeDefaultNone,
		TaskID:                info.TaskID,
	}
	xmm.SetUserAccount(mids)
	xmm.SetNotifyID(info.TaskID)
	xmm.SetTimeToLive(int64(info.ExpireTime))
	xmm.SetCallbackParam(strconv.FormatInt(info.APPID, 10))
	if info.Sound == model.SwitchOn {
		xmm.SetNotifyType(mi.NotifyTypeDefaultSound)
	}
	if info.Vibration == model.SwitchOn {
		xmm.SetNotifyType(mi.NotifyTypeDefaultVibration)
	}
	if info.Sound == model.SwitchOn && info.Vibration == model.SwitchOn {
		xmm.SetNotifyType(mi.NotifyTypeDefaultAll)
	}
	var response *mi.Response
	if response, err = d.clientMiByMids[info.APPID].Push(xmm); err != nil {
		log.Error("d.PushMi(%s,%s,%s) error(%v)", info.TaskID, scheme, mids, err)
		PromError("push:推送miByMids")
		return
	}
	res.Code = response.Code
	res.Msg = response.Reason
	if response.Code == mi.ResultCodeOk {
		res.Code = model.HTTPCodeOk
		res.Msg = response.Info
	}
	return
}

// PushHuawei push huawei notifications.
func (d *Dao) PushHuawei(c context.Context, info *model.PushInfo, scheme string, tokens []string) (res *huawei.Response, err error) {
	var index int
	if index, err = d.roundIndex(info.APPID, model.PlatformHuawei); err != nil {
		return
	}
	payload := huawei.NewMessage().SetTitle(info.Title).SetContent(info.Summary).SetCustomize("task_id", info.TaskID).SetCustomize("scheme", scheme).SetBiTag(info.TaskID).SetIcon(info.ImageURL)
	if info.PassThrough == model.SwitchOn {
		payload.SetMsgType(huawei.MsgTypePassthrough)
	}
	expire := time.Unix(int64(info.ExpireTime), 0)
	if res, err = d.clientsHuawei[info.APPID][index].Push(payload, tokens, expire); err != nil {
		if err == huawei.ErrLimit {
			return
		}
		log.Error("push huawei task(%s) resp(%+v) tokens(%v) error(%v)", info.TaskID, res, tokens, err)
		logPushError(info.TaskID, model.PlatformHuawei, tokens)
		return
	}
	log.Info("push huawei task(%s) tokens(%d) result(%+v) success", info.TaskID, len(tokens), res)
	return
}

// OppoMessage saves oppo message content.
func (d *Dao) OppoMessage(c context.Context, info *model.PushInfo, m *oppo.Message) (res *oppo.Response, err error) {
	var index int
	if index, err = d.roundIndex(info.APPID, model.PlatformOppo); err != nil {
		return
	}
	if res, err = d.clientsOppo[info.APPID][index].Message(m); err != nil {
		log.Error("save oppo message task(%s) result(%+v) error(%v)", info.TaskID, res, err)
		return
	}
	log.Info("save oppo message task(%s) result(%+v) success", info.TaskID, res)
	return
}

// PushOppo push oppo notifications.
func (d *Dao) PushOppo(c context.Context, info *model.PushInfo, msgID string, tokens []string) (res *oppo.Response, err error) {
	var index int
	if index, err = d.roundIndex(info.APPID, model.PlatformOppo); err != nil {
		return
	}
	if res, err = d.clientsOppo[info.APPID][index].Push(msgID, tokens); err != nil {
		log.Error("push oppo task(%s) resp(%+v) error(%v)", info.TaskID, res, err)
		logPushError(info.TaskID, model.PlatformOppo, tokens)
		return
	}
	log.Info("push oppo task(%s) tokens(%d) result(%+v) success", info.TaskID, len(tokens), res)
	return
}

// PushOppoOne push oppo notifications.
func (d *Dao) PushOppoOne(c context.Context, info *model.PushInfo, m *oppo.Message, token string) (res *oppo.Response, err error) {
	var index int
	if index, err = d.roundIndex(info.APPID, model.PlatformOppo); err != nil {
		return
	}
	if res, err = d.clientsOppo[info.APPID][index].PushOne(m, token); err != nil {
		log.Error("push oppo one task(%s) token(%s) result(%+v) error(%v)", info.TaskID, token, res, err)
		return
	}
	log.Info("push oppo one task(%s) token(%s) result(%+v) success", info.TaskID, token, res)
	return
}

// PushJpush push huawei notifications.
func (d *Dao) PushJpush(c context.Context, info *model.PushInfo, scheme string, tokens []string) (res *jpush.PushResponse, err error) {
	var (
		index   int
		ad      jpush.Audience
		notice  jpush.Notice
		plat    = jpush.NewPlatform(jpush.PlatformAndroid)
		payload = jpush.NewPayload()
		cbr     = jpush.NewCallbackReq()
		an      = &jpush.AndroidNotice{
			Title:     info.Title,
			Alert:     info.Summary,
			AlertType: jpush.AndroidAlertTypeNone,
			Extras: map[string]interface{}{
				"task_id": info.TaskID,
				"scheme":  scheme,
			},
		}
	)
	if index, err = d.roundIndex(info.APPID, model.PlatformJpush); err != nil {
		return
	}
	if info.Sound == model.SwitchOn {
		an.AlertType |= jpush.AndroidAlertTypeSound
	}
	if info.Vibration == model.SwitchOn {
		an.AlertType |= jpush.AndroidAlertTypeVibrate
	}
	if info.Sound == model.SwitchOn && info.Vibration == model.SwitchOn {
		an.AlertType = jpush.AndroidAlertTypeAll
	}
	ad.SetID(tokens)
	notice.SetAndroidNotice(an)
	payload.SetPlatform(plat)
	payload.SetAudience(&ad)
	payload.SetNotice(&notice)
	payload.Options.SetTimelive(int(int64(info.ExpireTime) - time.Now().Unix()))
	payload.Options.SetReturnInvalidToken(true)
	cbr.SetParam(map[string]string{"task": info.TaskID, "appid": strconv.FormatInt(info.APPID, 10)})
	payload.SetCallbackReq(cbr)
	if res, err = d.clientsJpush[info.APPID][index].Push(payload); err != nil {
		logPushError(info.TaskID, model.PlatformJpush, tokens)
		log.Error("push jpush task(%s) tokens(%d) result(%+v) error(%v)", info.TaskID, len(tokens), res, err)
		return
	}
	log.Info("push jpush task(%s) tokens(%d) result(%+v) success", info.TaskID, len(tokens), res)
	return
}

// PushFCM .
func (d *Dao) PushFCM(ctx context.Context, info *model.PushInfo, scheme string, tokens []string) (res *fcm.Response, err error) {
	var index int
	if index, err = d.roundIndex(info.APPID, model.PlatformFCM); err != nil {
		return
	}
	message := fcm.Message{
		Data: map[string]string{
			"task_id": info.TaskID,
			"scheme":  scheme,
		},
		RegistrationIDs: tokens,
		Priority:        fcm.PriorityHigh,
		DelayWhileIdle:  true,
		Notification: fcm.Notification{
			Title:       info.Title,
			Body:        info.Summary,
			ClickAction: "com.bilibili.app.in.com.bilibili.push.FCM_MESSAGE",
		},
		TimeToLive: int(int64(info.ExpireTime) - time.Now().Unix()),
		CollapseKey: strings.TrimFunc(info.TaskID, func(r rune) bool {
			return !unicode.IsNumber(r)
		}), // 应客户端要求，task_id 保证值转成 int 传到客户端
		Android: fcm.Android{Priority: fcm.PriorityHigh},
	}
	if res, err = d.clientsFCM[info.APPID][index].Send(&message); err != nil {
		log.Error("push fcm task(%s) tokens(%d) result(%+v) error(%v)", info.TaskID, len(tokens), res)
		PromError("push: 推送fcm")
		return
	}
	log.Info("push fcm task(%s) tokens(%d) result(%+v) error(%v)", info.TaskID, len(tokens), res)
	return
}
