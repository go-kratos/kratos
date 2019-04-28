package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_uri = "/danmaku/%d/rec"
)

func (d *Dao) danmakuURI(oid int64) string {
	return d.conf.Host.AI + fmt.Sprintf(_uri, oid%4)
}

// RecFlag get recommend flags from bigdata.
func (d *Dao) RecFlag(c context.Context, mid, aid, oid, limit, ps, pe int64, plat int32) (data []byte, err error) {
	var res struct {
		Code int64           `json:"code"`
		Msg  string          `json:"message"`
		Data json.RawMessage `json:"data"`
	}
	uri := d.danmakuURI(oid)
	params := url.Values{}
	realIP := metadata.String(c, metadata.RemoteIP)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("oid", strconv.FormatInt(oid, 10))
	params.Set("limit", strconv.FormatInt(limit, 10))
	params.Set("ps", strconv.FormatInt(ps, 10))
	params.Set("pe", strconv.FormatInt(pe, 10))
	params.Set("plat", strconv.FormatInt(int64(plat), 10))
	params.Set("ip", realIP)
	if err = d.httpCli.Get(c, uri, realIP, params, &res); err != nil {
		PromError(uri)
		log.Error("d.httpCli.Get(%s?%s) error(%v)", uri, params.Encode(), err)
		return
	}
	if res.Code != 0 {
		PromError(uri)
		log.Error("d.httpCli.Get(%s?%s) code:%d msg:%s", uri, params.Encode(), res.Code, res.Msg)
		err = fmt.Errorf("bigdata response code(%d) error", res.Code)
		return
	}
	data = res.Data
	return
}
