package dao

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"go-common/app/admin/main/apm/model/ecode"
	"go-common/library/log"
)

// const (
// 	_codesLangsSQL = "select a.code,a.message,a.mtime,IFNULL(b.locale,''),IFNULL(b.msg,''),IFNULL(b.mtime,'') as bmtime from codes as a left join code_msg as b on a.id=b.code_id"
// )

// GetCodes ...
func (d *Dao) GetCodes(c context.Context, Interval1, Interval2 string) (data []*codes.Codes, err error) {
	var (
		req    *http.Request
		uri    = "http://sven.bilibili.co/x/admin/apm/ecode/get/ecodes"
		ret    = codes.ResultCodes{}
		params = url.Values{}
	)
	params.Set("interval1", Interval1)
	params.Set("interval2", Interval2)
	if req, err = d.client.NewRequest(http.MethodGet, uri, "", params); err != nil {
		log.Error("http.NewRequest error(%v)", err)
		return
	}
	if err = d.client.Do(c, req, &ret); err != nil {
		log.Error("client Do error(%v)", err)
		return
	}
	if ret.Code != 0 {
		err = fmt.Errorf("%s params(%s) response return_code(%d)", uri, params.Encode(), ret.Code)
		log.Error("error(%v)", err)
		return
	}
	data = ret.Data
	return
}
