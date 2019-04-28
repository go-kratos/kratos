package notice

import (
	"context"
	"fmt"
	"net/url"

	"go-common/library/log"
)

const (
	_urlAudioPlayList = "https://m.bilibili.com/audio/am%d"
)

// AudioPlayList is show AudioPlay list
func (d *Dao) AudioPlayList(c context.Context, oid int64) (title, link string, err error) {
	params := url.Values{}
	uri := fmt.Sprintf(d.urlAudioPlaylist, oid)
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Title string `json:"title"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, uri, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", uri, params.Encode(), err)
		return
	}
	if res.Data == nil {
		err = fmt.Errorf("url:%s code:%d", uri, res.Code)
		return
	}
	title = res.Data.Title
	link = fmt.Sprintf(_urlAudioPlayList, oid)
	return
}
