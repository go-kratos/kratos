package notice

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/library/log"
)

// Drawyoo return link.
func (d *Dao) Drawyoo(c context.Context, hid int64) (title, link string, err error) {
	params := url.Values{}
	params.Set("hid", strconv.FormatInt(hid, 10))
	params.Set("act", "getHidInfo")
	var res struct {
		State int `json:"state"`
		Data  []*struct {
			Title string `json:"title"`
			Link  string `json:"link"`
		} `json:"data"`
	}
	if err = d.drawyooHTTPClient.Post(c, d.urlDrwayoo, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", d.urlDrwayoo, params.Encode(), err)
		return
	}
	if len(res.Data) == 0 {
		err = fmt.Errorf("url:%s code:%d", d.urlDrwayoo, res.State)
		return
	}
	title = res.Data[0].Title
	link = res.Data[0].Link
	return
}
