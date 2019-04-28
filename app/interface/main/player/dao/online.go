package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/net/metadata"
)

const _onlineCountURI = "/x/internal/chat/num/ol"

// OnlineCount get online count by aid and cid
func (d *Dao) OnlineCount(c context.Context, aid, cid int64) (count int64, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("cids", strconv.FormatInt(cid, 10))
	var res struct {
		Code int `json:"code"`
		Data []struct {
			Cid   int64 `json:"cid"`
			Count int64 `json:"count"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.onlineCountURL, ip, params, &res); err != nil {
		PromError("OnlineNum接口错误", "d.client.Get(%s) error(%v)", d.onlineCountURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		PromError("OnlineNum接口Code错误", "d.client.Get(%s) code(%d) error", d.onlineCountURL, res.Code)
		return
	}
	if len(res.Data) == 0 || res.Data[0].Cid != cid {
		PromError("OnlineNum接口数据错误", "d.client.Get(%s) data(%v) error", d.onlineCountURL, res.Data)
	}
	count = res.Data[0].Count
	return
}
