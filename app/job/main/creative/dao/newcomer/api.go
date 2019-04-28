package newcomer

import (
	"context"
	"errors"
	"net/url"

	"go-common/library/log"
	"go-common/library/xstr"
)

// SendNotify send msg notify user
func (d *Dao) SendNotify(c context.Context, mids []int64, mc, title, context string) (err error) {
	var (
		params = url.Values{}
		res    struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
			Data *struct {
				TotalCount   int     `json:"total_count"`
				ErrorCount   int     `json:"error_count"`
				ErrorMidList []int64 `json:"error_mid_list"`
			} `json:"data"`
		}
	)
	params.Set("mc", mc)                        //消息码，用于识别消息类别
	params.Set("data_type", "4")                //消息类型：1、回复我的 2、@我 3、收到的爱 4、业务通知 5、系统公告
	params.Set("title", title)                  //消息标题
	params.Set("context", context)              //消息实体内容
	params.Set("mid_list", xstr.JoinInts(mids)) //用于接收该消息的用户mid列表，不超过1000个(半角逗号分割)

	log.Info("SendNotify params(%+v)|msgURI(%s)", params.Encode(), d.msgURI)
	if err = d.httpClient.Post(c, d.msgURI, "", params, &res); err != nil {
		log.Error("d.httpClient.Post(%s,%v,%d)", d.msgURI, params, err)
		return
	}
	if res.Code != 0 {
		err = errors.New("code != 0")
		log.Error("d.httpClient.Post(%s,%v,%v,%d)", d.msgURI, params, err, res.Code)
	}
	if res.Data != nil {
		log.Info("SendNotify log total_count(%d) error_count(%d) error_mid_list(%v)", res.Data.TotalCount, res.Data.ErrorCount, res.Data.ErrorMidList)
	}
	return
}
