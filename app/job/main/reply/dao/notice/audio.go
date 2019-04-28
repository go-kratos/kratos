package notice

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/library/log"
)

const (
	_urlAudio = "https://m.bilibili.com/audio/au%d"
)

type audio struct {
	Title string `json:"title"`
}

// Audio is Audio
func (d *Dao) Audio(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	uri := d.urlAudio
	params.Set("ids", strconv.FormatInt(oid, 10))
	var res struct {
		Code int              `json:"code"`
		Data map[int64]*audio `json:"data"`
	}
	if err = d.httpClient.Get(c, uri, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", uri, params.Encode(), err)
		return
	}
	if res.Data == nil {
		err = fmt.Errorf("url:%s code:%d", uri, res.Code)
		return
	}
	if r := res.Data[oid]; r != nil {
		title = r.Title
	}
	link = fmt.Sprintf(_urlAudio, oid)
	return
}
