package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/library/log"
)

// BlockAccount block account
func (d *Dao) BlockAccount(c context.Context, mid int64, reason string) (err error) {
	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%d", mid))
	params.Set("admin_reason", reason)
	params.Set("blockType", "6")
	params.Set("operator", "anticheat")
	params.Set("type", "json")
	var resp struct {
		Code int64 `json:"code"`
	}
	// get
	if err = d.httpClient.Get(c, d.c.Property.BlockAccountURL, "", params, &resp); err != nil {
		log.Error("httpClient.Do() error(%v)", err)
		return
	}
	if resp.Code != 0 {
		err = fmt.Errorf("GET block account url resp(%v)", resp)
		return
	}
	log.Info("account user(%d) block suc(%v)", mid, resp)
	return
}
