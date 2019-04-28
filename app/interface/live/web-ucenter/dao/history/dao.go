package history

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"go-common/library/ecode"

	"github.com/pkg/errors"

	"go-common/app/interface/live/web-ucenter/conf"
	"go-common/app/interface/live/web-ucenter/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

// Dao dao
type Dao struct {
	c      *conf.Config
	client *bm.Client
}

const (
	_historyResourceURI = "/x/internal/v2/history/resource"
	_historyDeleteURI   = "/x/internal/v2/history/clear"
	_historyPageSize    = 24
)

// New init
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		client: bm.NewClient(c.HTTPClient),
	}
	return
}

// GetMainHistory 获取直播历史记录
func (d *Dao) GetMainHistory(c context.Context, mid int32) (data []*model.HistoryData, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", fmt.Sprint(mid))
	params.Set("business", "live")
	params.Set("pn", fmt.Sprint(1))
	params.Set("ps", fmt.Sprint(_historyPageSize))
	params.Set("appkey", conf.APPKey)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))

	var res struct {
		Code int                  `json:"code"`
		Msg  string               `json:"msg"`
		Data []*model.HistoryData `json:"data"`
	}
	if d.client.Get(c, conf.MainInnerHostHTTP+_historyResourceURI, ip, params, &res); err != nil {
		err = errors.WithMessage(err, "调用主站获取观看历史错误")
		log.Error("call_history_resource_error:httpCode=%d params=%s", res.Code, params)
		return
	}
	if int(res.Code) != ecode.OK.Code() {
		log.Error("call_history_resource_error:%d", res.Code)
	}
	log.Info("call_history_resource_info:param:%v,%v", mid, res.Data)
	data = res.Data

	return
}

// DelHistory 删除直播历史记录
func (d *Dao) DelHistory(c context.Context, mid int64) (data int32, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("appkey", conf.APPKey)
	params.Set("mid", fmt.Sprint(mid))
	params.Set("business", "live")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	var res struct {
		Code int32 `json:"code"`
	}
	if d.client.Post(c, conf.MainInnerHostHTTP+_historyDeleteURI, ip, params, &res); err != nil {
		err = errors.WithMessage(err, "调用主站删除观看历史错误")
		log.Error("call_history_delete_error:%s", err)
		return
	}
	if int(res.Code) != ecode.OK.Code() {
		log.Error("call_history_delete_code_error:%d", res.Code)
	}
	data = res.Code
	return
}
