package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/space/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_arcSearchURI    = "/space/search/v2"
	_arcSearchType   = "sub_video"
	_additionalRanks = "-6"
)

// ArcSearchList archive search.
func (d *Dao) ArcSearchList(c context.Context, arg *model.SearchArg) (data *model.SearchRes, total int, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("search_type", _arcSearchType)
	params.Set("additional_ranks", _additionalRanks)
	if arg.Mid > 0 {
		params.Set("mid", strconv.FormatInt(arg.Mid, 10))
	}
	params.Set("page", strconv.Itoa(arg.Pn))
	params.Set("pagesize", strconv.Itoa(arg.Ps))
	params.Set("clientip", ip)
	if arg.Tid > 0 {
		params.Set("tid", strconv.FormatInt(arg.Tid, 10))
	}
	if arg.Order != "" {
		params.Set("order", arg.Order)
	}
	if arg.Keyword != "" {
		params.Set("keyword", arg.Keyword)
	}
	var res struct {
		Code   int              `json:"code"`
		Total  int              `json:"total"`
		Result *model.SearchRes `json:"result"`
	}
	if err = d.httpR.Get(c, d.arcSearchURL, ip, params, &res); err != nil {
		log.Error("d.httpR.Get(%s,%v) error(%v)", d.arcSearchURL, arg, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpR.Get(%s,%v) code error(%d)", d.arcSearchURL, arg, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	data = res.Result
	total = res.Total
	return
}
