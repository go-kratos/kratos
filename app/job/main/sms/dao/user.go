package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/job/main/sms/model"
	"go-common/library/log"
)

// UserMobile get user mobile
func (d *Dao) UserMobile(c context.Context, mid int64) (*model.UserMobile, error) {
	res := struct {
		Code int              `json:"code"`
		Data model.UserMobile `json:"data"`
	}{}
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	if err := d.httpClient.Get(c, d.c.Sms.PassportMobileURL, "", params, &res); err != nil {
		log.Error("d.GetUserMobile(%d) error(%v)", mid, err)
		return nil, err
	}
	if res.Code != 0 {
		return nil, fmt.Errorf("GetUserMobile(%d) error, res(%+v)", mid, &res)
	}
	log.Info("GetUserMobile(%d) res(%+v)", mid, &res)
	return &res.Data, nil
}
