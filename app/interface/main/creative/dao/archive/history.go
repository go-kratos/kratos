package archive

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/creative/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_hList = "/videoup/history/list"
	_hView = "/videoup/history/view"
)

// HistoryList get the history of aid
func (d *Dao) HistoryList(c context.Context, mid, aid int64, ip string) (historys []*archive.ArcHistory, err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int                   `json:"code"`
		Data []*archive.ArcHistory `json:"data"`
	}
	if err = d.client.Get(c, d.hList, ip, params, &res); err != nil {
		log.Error("archive.HistoryList url(%s) mid(%d) error(%v)", d.hList+"?"+params.Encode(), mid, err)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("archive.HistoryList url(%s) mid(%d) res(%v)", d.hList+"?"+params.Encode(), mid, res)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	historys = res.Data
	return
}

// HistoryView get the history of hid
func (d *Dao) HistoryView(c context.Context, mid, hid int64, ip string) (history *archive.ArcHistory, err error) {
	params := url.Values{}
	params.Set("hid", strconv.FormatInt(hid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int                 `json:"code"`
		Data *archive.ArcHistory `json:"data"`
	}
	if err = d.client.Get(c, d.hView, ip, params, &res); err != nil {
		log.Error("archive.HistoryView url(%s) mid(%d) error(%v)", d.hView+"?"+params.Encode(), mid, err)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("archive.HistoryView url(%s) mid(%d) res(%v)", d.hView+"?"+params.Encode(), mid, res)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	history = res.Data
	return
}
