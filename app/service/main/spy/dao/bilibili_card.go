package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/service/main/spy/model"
	"go-common/library/log"
)

// UnicomGiftState get unicom gift state by mid from account service.
func (d *Dao) UnicomGiftState(c context.Context, mid int64) (state int, err error) {
	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%d", mid))
	var resp struct {
		Code int64                  `json:"code"`
		Data *model.UnicomGiftState `json:"data"`
	}
	if err = d.httpClient.Get(c, d.c.Property.UnicomGiftStateURL, "", params, &resp); err != nil {
		log.Error("message url(%s) error(%v)", d.c.Property.UnicomGiftStateURL+"?"+params.Encode(), err)
		return
	}
	if resp.Code != 0 {
		err = fmt.Errorf("GET UnicomGiftStateURL url resp(%v)", resp)
		return
	}
	log.Info("GET UnicomGiftStateURL suc url(%s) resp(%v)", d.c.Property.UnicomGiftStateURL+"?"+params.Encode(), resp)
	state = resp.Data.State
	return
}
