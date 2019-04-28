package dao

import (
	"context"
	"time"

	"go-common/app/service/main/sms/model"
	"go-common/library/log"
)

const _retry = 3

// PubSingle pub single sms to databus
func (d *Dao) PubSingle(ctx context.Context, l *model.ModelSend) (err error) {
	for i := 0; i < _retry; i++ {
		if err = d.databus.Send(ctx, l.Mobile, l); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		log.Error("PubSingle(%+v) error(%v)", l, err)
		return
	}
	log.Info("PubSingle(%+v) success.", l)
	return
}

// PubBatch pub batch sms to databus
func (d *Dao) PubBatch(ctx context.Context, l *model.ModelSend) (err error) {
	for i := 0; i < _retry; i++ {
		if err = d.databus.Send(ctx, l.Code, l); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		log.Error("PubBatch(%+v) error(%v)", l, err)
		return
	}
	log.Info("PubBatch(%+v) success.", l)
	return
}
