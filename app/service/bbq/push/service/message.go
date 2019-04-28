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
	jsoniter "github.com/json-iterator/go"
)

// Message .
func (s *Service) Message(ctx context.Context, in *v1.MessageRequest) (*v1.MessageResponse, error) {
	tasks := make([]workpool.Task, len(in.Dev))
	for i, dev := range in.Dev {
		tasks[i] = NewMessageTask(&ctx, s.c.JPush, s.dao, dev, in.Body)
	}

	ftasks, err := s.parallel(&ctx, tasks)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "parallel message error"), log.KV("error", err))
	}
	result := s.wait(ctx, ftasks)
	log.Infov(ctx, log.KV("log", result))

	return &v1.MessageResponse{
		Result: result,
	}, err
}

// AsyncMessage .
func (s *Service) AsyncMessage(ctx context.Context, in *v1.MessageRequest) (*v1.MessageResponse, error) {
	tasks := make([]workpool.Task, len(in.Dev))
	for i, dev := range in.Dev {
		tasks[i] = NewMessageTask(&ctx, s.c.JPush, s.dao, dev, in.Body)
	}

	_, err := s.parallel(&ctx, tasks)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "AsyncNotification parallel error"), log.KV("error", err))
	}

	return &v1.MessageResponse{}, err
}

// MessageTask .
type MessageTask struct {
	ctx  *context.Context
	c    *conf.JPushConfig
	d    *dao.Dao
	dev  *v1.Device
	body *v1.MessageBody
}

// NewMessageTask .
func NewMessageTask(ctx *context.Context, c *conf.JPushConfig, d *dao.Dao, dev *v1.Device, body *v1.MessageBody) *MessageTask {
	return &MessageTask{
		ctx:  ctx,
		c:    c,
		d:    d,
		dev:  dev,
		body: body,
	}
}

// Run .
func (pt *MessageTask) Run() *[]byte {
	ext := make(map[string]string)
	jsoniter.Unmarshal([]byte(pt.body.Extra), ext)

	tracer, _ := trace.FromContext(*pt.ctx)
	ext["callback"] = fmt.Sprintf("http://bbq.bilibili.com/bbq/app-bbq/push/callback?tid=%s&nid=%d", tracer, pt.dev.SendNo)

	var platform []string
	if pt.dev.Platform == model.PlatformAndroid {
		platform = append(platform, "android")
	} else {
		platform = append(platform, "ios")
	}
	msg := &jpush.Message{
		Title:       pt.body.Title,
		ContentType: pt.body.ContentType,
		MsgContent:  pt.body.Content,
		Extras:      ext,
	}
	payload := &jpush.Payload{
		Platform: platform,
		Audience: &jpush.Audience{
			RegID: []string{pt.dev.RegisterID},
		},
		Message: msg,
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
