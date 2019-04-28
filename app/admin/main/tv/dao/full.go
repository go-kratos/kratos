package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

type msgReturn struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    []*model.APKInfo `json:"data"`
}

// FullImport .
func (d *Dao) FullImport(c context.Context, build int) (result []*model.APKInfo, err error) {
	var (
		res     = &msgReturn{}
		fullURL = d.fullURL
	)
	params := url.Values{}
	params.Set("version_code", fmt.Sprintf("%d", build))
	err = d.httpSearch.Get(c, fullURL, "", params, res)
	if err != nil {
		log.Error("d.httpSearch.Get(%s) error(%v)", fullURL+"?"+params.Encode(), err)
		return
	}
	result = res.Data
	if res.Code != ecode.OK.Code() {
		err = fmt.Errorf("return code:%d", res.Code)
		log.Error("d.httpSearch.Get(%s) error(%v)", fullURL+"?"+params.Encode(), err)
	}
	return
}
