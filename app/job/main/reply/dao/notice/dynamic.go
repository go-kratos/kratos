package notice

import (
	"context"
	"fmt"
	"net/url"

	"go-common/library/log"
)

const (
	_dynamicLink = "https://t.bilibili.com/%d"
)

// Dynamic return link and content.
func (d *Dao) Dynamic(c context.Context, oid int64) (content, link string, err error) {
	params := url.Values{}
	uri := fmt.Sprintf(d.urlDynamic, oid)
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Pairs []struct {
				DynamicID int64  `json:"dynamic_id"`
				Content   string `json:"rp_cont"`
				Type      int32  `json:"type"`
			} `json:"pairs"`
			TotalCount int64 `json:"total_count"`
		} `json:"data,omitempty"`
		Message string `json:"message"`
	}
	if err = d.httpClient.Get(c, uri, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", uri, params.Encode(), err)
		return
	}

	if res.Code != 0 || res.Data == nil || len(res.Data.Pairs) == 0 {
		err = fmt.Errorf("get dynamic failed!url:%s?%s code:%d message:%s pairs:%v", uri, params.Encode(), res.Code, res.Message, res.Data.Pairs)
		return
	}
	content = res.Data.Pairs[0].Content
	link = fmt.Sprintf(_dynamicLink, oid)
	return
}
