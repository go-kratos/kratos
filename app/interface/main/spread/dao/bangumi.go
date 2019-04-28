package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/interface/main/spread/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// BangumiContent .
func (d *Dao) BangumiContent(c context.Context, pn, ps int, typ int8, appkey string) (resp model.BangumiResp, err error) {
	params := url.Values{}
	params.Set("page_no", strconv.Itoa(pn))
	params.Set("page_size", strconv.Itoa(ps))
	params.Set("season_type", strconv.Itoa(int(typ)))
	params.Set("bsource", appkey)
	err = d.httpClient.Get(c, d.c.Spread.BangumiContentURL, "", params, &resp)
	u := fmt.Sprintf("%s?%s", d.c.Spread.BangumiContentURL, params.Encode())
	if err != nil {
		PromError("bangumi:content接口")
		log.Errorv(c, log.KV("err", err), log.KV("url", u))
		return
	}
	if resp.Code != 0 {
		PromError("bangumi:content接口")
		log.Errorv(c, log.KV("res", resp), log.KV("url", u))
		err = ecode.Int(resp.Code)
		return
	}
	return
}

// BangumiOff .
func (d *Dao) BangumiOff(c context.Context, pn, ps int, typ int8, appkey string, ts int64) (resp model.BangumiOffResp, err error) {
	params := url.Values{}
	params.Set("page_no", strconv.Itoa(pn))
	params.Set("page_size", strconv.Itoa(ps))
	params.Set("timestamp", strconv.FormatInt(ts, 10))
	params.Set("season_type", strconv.Itoa(int(typ)))
	params.Set("bsource", appkey)
	err = d.httpClient.Get(c, d.c.Spread.BangumiOffURL, "", params, &resp)
	u := fmt.Sprintf("%s?%s", d.c.Spread.BangumiOffURL, params.Encode())
	if err != nil {
		PromError("bangumi:off接口")
		log.Errorv(c, log.KV("err", err), log.KV("url", u))
		return
	}
	if resp.Code != 0 {
		PromError("bangumi:off接口")
		log.Errorv(c, log.KV("res", resp), log.KV("url", u))
		err = ecode.Int(resp.Code)
		return
	}
	return
}
