package manager

import (
	"context"

	"go-common/library/log"
	"net/url"
	"strconv"
)

// searchUpdate .
func (d *Dao) SearchUpdate(c context.Context, business, data string, aid int64) (err error) {
	params := url.Values{}
	params.Set("business", business)
	params.Set("data", data)
	params.Set("aid", strconv.FormatInt(aid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpW.Post(c, d.searchUpdateURL, "", params, &res); err != nil {
		log.Error("searchUpdate d.httpW.Post(%s) error(%v)", d.searchUpdateURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("searchUpdate url(%s) code(%d)", d.searchUpdateURL+"?"+params.Encode(), res.Code)
	}
	log.Info("aid (%d) SearchUpdate url(%s) code(%d)", aid, d.searchUpdateURL+"?"+params.Encode(), res.Code)
	return
}
