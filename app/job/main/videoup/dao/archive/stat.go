package archive

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/library/log"
)

// Stat get archive stat.
func (d *Dao) Stat(c context.Context, aid int64) (click int, err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Click int `json:"click"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.statURI, "", params, &res); err != nil {
		log.Error("archive stat url(%s) error(%v)", d.statURI, err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("archive stat call failed")
		return
	}
	click = res.Data.Click
	return
}
