package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/service/bbq/push/api/grpc/v1"
	"go-common/app/service/bbq/push/conf"
	"go-common/app/service/bbq/push/dao"
	"go-common/app/service/bbq/push/dao/jpush"
	"go-common/app/service/bbq/push/model"
	"go-common/library/log"
	"go-common/library/net/trace"

	"github.com/Dai0522/workpool"
	"github.com/json-iterator/go"
)

// Notification .
func (s *Service) Notification(ctx context.Context, in *v1.NotificationRequest) (*v1.NotificationResponse, error) {
	tasks := make([]workpool.Task, len(in.Dev))
	for i, dev := range in.Dev {
		tasks[i] = NewNotificationTask(&ctx, s.c.JPush, s.dao, dev, in.Body)
	}

	ftasks, err := s.parallel(&ctx, tasks)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "Notification parallel error"), log.KV("error", err))
	}
	result := s.wait(ctx, ftasks)
	log.Infov(ctx, log.KV("log", result))

	return &v1.NotificationResponse{
		Result: result,
	}, err
}

// AsyncNotification .
func (s *Service) AsyncNotification(ctx context.Context, in *v1.NotificationRequest) (*v1.NotificationResponse, error) {
	tasks := make([]workpool.Task, len(in.Dev))
	for i, dev := range in.Dev {
		tasks[i] = NewNotificationTask(&ctx, s.c.JPush, s.dao, dev, in.Body)
	}

	_, err := s.parallel(&ctx, tasks)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "AsyncNotification parallel error"), log.KV("error", err))
	}

	return &v1.NotificationResponse{}, err
}

// NotificationTask .
type NotificationTask struct {
	ctx  *context.Context
	c    *conf.JPushConfig
	d    *dao.Dao
	dev  *v1.Device
	body *v1.NotificationBody
}

// NewNotificationTask .
func NewNotificationTask(ctx *context.Context, c *conf.JPushConfig, d *dao.Dao, dev *v1.Device, body *v1.NotificationBody) *NotificationTask {
	return &NotificationTask{
		ctx:  ctx,
		c:    c,
		d:    d,
		dev:  dev,
		body: body,
	}
}

// Run .
func (pt *NotificationTask) Run() *[]byte {
	ext := make(map[string]string)
	err := jsoniter.UnmarshalFromString(pt.body.Extra, &ext)
	if err != nil {
		log.Error("NotificationTask json Unmarshal error [%+v]", err)
	}

	tracer, _ := trace.FromContext(*pt.ctx)
	ext["callback"] = fmt.Sprintf("http://bbq.bilibili.com/bbq/app-bbq/push/callback?tid=%s&nid=%d", tracer, pt.dev.SendNo)

	var platform []string
	notify := &jpush.Notification{}
	if pt.dev.Platform == model.PlatformAndroid {
		platform = append(platform, "android")
		notify.Android = &jpush.AndroidNotification{
			Alert:  pt.body.Content,
			Extras: ext,
		}
	} else {
		platform = append(platform, "ios")
		var alert interface{}
		if pt.body.Title == "" {
			alert = pt.body.Content
		} else {
			alert = &jpush.IOSAlert{
				Title: pt.body.Title,
				Body:  pt.body.Content,
			}
		}
		notify.IOS = &jpush.IOSNotification{
			Alert:  alert,
			Sound:  pt.body.Sound,
			Badge:  pt.body.Badge,
			Extras: ext,
		}
	}

	payload := &jpush.Payload{
		Platform: platform,
		Audience: &jpush.Audience{
			RegID: []string{pt.dev.RegisterID},
		},
		Notification: notify,
		Options: &jpush.Option{
			SendNo:         int(pt.dev.SendNo),
			ApnsProduction: pt.c.ApnsProduction,
		},
	}

	b, _ := jsoniter.Marshal(payload)
	res, _ := pt.d.JPush.Push(b)

	// 埋点
	result := v1.PushResult{}
	result.Unmarshal(res)
	Infoc.Info(tracer, "", "", pt.dev.SendNo, "", "", "", time.Now().Unix(), result)

	return &res
}
