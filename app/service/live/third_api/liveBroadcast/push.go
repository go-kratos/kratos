package liveBroadcast

import (
	"context"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xhttp "net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	_liveBroadCastURLAddr = "http://live-dm.bilibili.co:80/dm/1/push"
)

// type liveBroadcastReq struct {
// 	Ensure int64 `json:"ensure"`
// 	Cid    int64 `json:"cid"`
// }

type liveBroadcastResp struct {
	Code int64 `json:"ret"`
}

//PushBroadcast 广播
func (c *Client) PushBroadcast(ctx context.Context, cid int64, ensure int64, msg string) (err error) {
	if len([]rune(msg)) > 2000 {
		return
	}

	resp := &liveBroadcastResp{}
	cli := bm.NewClient(c.getConf())

	param := url.Values{}

	param.Set("cid", strconv.FormatInt(cid, 10))
	param.Set("ensure", strconv.FormatInt(0, 10))

	url := _liveBroadCastURLAddr + "?" + param.Encode()
	req, err := xhttp.NewRequest(xhttp.MethodPost, url, strings.NewReader(msg))
	if err != nil {
		log.Error("[BroadCastError]error:%+v=", err)
		return err
	}
	if err := cli.Do(ctx, req, resp); err != nil {
		log.Error("[BroadCastError]error:%+v=", err)
		return err
	}
	if resp.Code != 1 {
		log.Error("BroadCastError] errorcode:%d", resp.Code)
	}
	return
}
