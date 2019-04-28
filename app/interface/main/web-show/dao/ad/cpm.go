package ad

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	ad "go-common/app/interface/main/web-show/model/resource"
	"go-common/library/log"
	"go-common/library/xstr"
)

// Cpms get ads from cpm platform
func (d *Dao) Cpms(c context.Context, mid int64, ids []int64, sid, ip, country, province, city, buvid string) (advert *ad.Ad, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("sid", sid)
	params.Set("buvid", buvid)
	params.Set("resource", xstr.JoinInts(ids))
	params.Set("ip", ip)
	params.Set("country", country)
	params.Set("province", province)
	params.Set("city", city)
	var res struct {
		Code int    `json:"code"`
		Data *ad.Ad `json:"data"`
	}
	if err = d.httpClient.Get(c, d.cpmURL, "", params, &res); err != nil {
		log.Error("cpm url(%s) error(%v)", d.cpmURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("cpm api failed(%d)", res.Code)
		log.Error("url(%s) res code(%d) or res.data(%v)", d.cpmURL+"?"+params.Encode(), res.Code, res.Data)
		return
	}
	advert = res.Data
	return
}
