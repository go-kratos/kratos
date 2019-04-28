package monitor

import (
	"context"
	"go-common/app/job/main/aegis/model/monitor"
	"go-common/library/ecode"
	"go-common/library/log"
	"net/url"
	"strconv"
)

const (
	_arcAdditURL = "/videoup/archive/addit"
)

// ArchiveAttr 获取稿件都附加属性
func (d *Dao) ArchiveAttr(c context.Context, aid int64) (addit *monitor.ArchiveAddit, err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	var res struct {
		Code int                   `json:"code"`
		Data *monitor.ArchiveAddit `json:"data"`
	}
	if err = d.http.Get(c, d.URLArcAddit, "", params, &res); err != nil {
		log.Error("d.ArchiveAttr(%s)  error(%v)", d.URLArcAddit+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.ArchiveAttr(%s) Code=(%d)", d.URLArcAddit+"?"+params.Encode(), res.Code)
		return
	}
	if res.Data == nil {
		err = ecode.NothingFound
		log.Warn("d.ArchiveAttr(%s) Code=(%d) data nil", d.URLArcAddit+"?"+params.Encode(), res.Code)
		return
	}
	addit = res.Data
	return
}
