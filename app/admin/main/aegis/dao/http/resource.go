package http

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/admin/main/aegis/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// SyncResource 同步到业务方
func (d *Dao) SyncResource(c context.Context, act *model.Action, ropt map[string]interface{}) (code int, err error) {
	params := url.Values{}

	for k, v := range ropt {
		params.Set(k, fmt.Sprint(v))
	}

	requestParams := fmt.Sprintf("%s?%s", act.URL, params.Encode())
	log.Info("SyncResource url(%v) params(%s)", act.URL, params.Encode())

	if d.c.Debug == "local" || d.c.Debug == "nobusiness" {
		return
	}
	res := new(struct {
		Code    int    `json:"code"`
		Msg     string `json:"msg"`
		Message string `json:"message"`
	})
	err = d.clientW.Post(c, act.URL, "", params, res)
	code = res.Code
	if err != nil || res.Code != 0 {
		log.Error("clientW.Post err(%v) response(%+v) request(%s)", err, res, requestParams)
		err = ecode.Errorf(ecode.AegisBusinessSyncErr, "业务回调错误 request(%s) httperror(%v) response(%+v)", requestParams, err, res)
	}
	return
}
