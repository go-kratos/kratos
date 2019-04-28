package extern

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"time"

	"go-common/library/log"
)

const (
	PathDeleteReplyByIds = "/x/internal/v2/reply/del"
)

type Reply struct {
	Id    int64 `json:"id"`
	OId   int64 `json:"oid"`
	OType int64 `json:"typ"`
}

var replySvrCli *ReplyServiceClient

type ReplyServiceClient struct {
	*commonClient
	host string
}

type ReplyServiceResp struct {
	Code    int         `json:"code"`
	Message string      `json:"messge"`
	Data    interface{} `json:"data"`
}

type Replys []*Reply

func (rs Replys) OIds() string {
	var s string
	for _, r := range rs {
		s += fmt.Sprintf("%d,", r.OId)
	}
	return s
}

func (rs Replys) Ids() string {
	var s string
	for _, r := range rs {
		s += fmt.Sprintf("%d,", r.Id)
	}
	return s[:len(s)-1]
}

func (rs Replys) OTypes() string {
	var s string
	for _, r := range rs {
		s += fmt.Sprintf("%d,", r.OType)
	}
	return s[:len(s)-1]
}

func (self *ReplyServiceClient) DeleteReply(ctx context.Context, adminId int64, rs []*Reply) error {
	val := url.Values{}
	val.Add("adid", fmt.Sprintf("%d", adminId))
	val.Add("adname", "antispam")
	val.Add("oid", Replys(rs).OIds())
	val.Add("rpid", Replys(rs).Ids())
	val.Add("type", Replys(rs).OTypes())
	val.Add("moral", "0")
	val.Add("notify", "false")
	val.Add("remark", "")
	val.Add("ftime", "")
	val.Add("reason", "delete by antispam")
	return self.do(ctx, PathDeleteReplyByIds, val, &ReplyServiceResp{}, replySvrCli.httpCli.Post)
}

func (rs *ReplyServiceClient) do(ctx context.Context,
	urlPath string, params url.Values, resp *ReplyServiceResp,
	fn func(ctx context.Context, uri string, ip string, params url.Values, resp interface{}) error,
) error {
	params.Set("appkey", rs.key)
	params.Set("appsecret", rs.secret)
	params.Set("ts", fmt.Sprintf("%d", time.Now().Unix()+int64(10)))

	urlAddr := path.Join(rs.host + urlPath)

	err := fn(ctx, urlAddr, "", params, resp)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		err = fmt.Errorf("Call reply service(%s), response code is not 0, resp:%v", urlAddr+"?"+params.Encode(), resp)
		log.Error("%v", err)
		return err
	}
	log.Info("Call reply service(%s) successful, resp: %v", urlAddr+"?"+params.Encode(), resp)
	return nil
}
