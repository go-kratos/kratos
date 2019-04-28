package dao

import (
	"context"
	"net/http"

	"go-common/library/ecode"
)

// Live gets live dynamic count
func (d *Dao) Live(c context.Context) (count int, err error) {
	var req *http.Request
	if req, err = d.httpR.NewRequest("GET", d.liveURI, "", nil); err != nil {
		PromError("直播Live接口", "Live d.httpR.NewRequest(%s) error(%v)", d.liveURI, err)
		return
	}
	var res struct {
		Code int                 `json:"code"`
		Msg  string              `json:"msg"`
		Data struct{ Count int } `json:"data"`
	}
	err = d.httpR.Do(c, req, &res)
	if err != nil {
		PromError("直播Live接口", "Live d.httpR.Do(%s) error(%v)", d.liveURI, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		PromError("直播Live接口", "Live dao.httpR.Do(%s) error(%v)", d.liveURI, err)
		err = ecode.Int(res.Code)
		return
	}
	count = res.Data.Count
	return
}
