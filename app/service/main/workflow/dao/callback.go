package dao

import (
	"context"
	"errors"

	"go-common/app/service/main/workflow/model"
	"go-common/library/log"
)

// Callback callback message
func (d *Dao) Callback(c context.Context, chall *model.Challenge, businessID int8) (err error) {
	if URL, ok := d.callbackMap[businessID]; ok {
		if err = d.callback.Post(context.Background(), URL, "", nil, &chall); err != nil {
			log.Error("d.CallbackSetting(%s) error(%v)", chall, err)
			return
		}
		return
	}
	return errors.New("Callback cannot find businessID in callbackMap")
}
