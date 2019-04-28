package search

import (
	"context"
	"fmt"
	"net/url"

	searchMdl "go-common/app/interface/main/tv/model/search"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// SearchAll gets the search all_tv result
func (d *Dao) SearchAll(ctx context.Context, req *searchMdl.ReqSearch) (result searchMdl.RespAll, common *searchMdl.ResultResponse, err error) {
	params := commonParam(req)
	if err = d.client.Get(ctx, d.resultURL, "", params, &result); err != nil {
		log.Error("[result] SearchAll URL(%s) error[%v]", d.resultURL, err)
		return
	}
	if result.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(result.Code), "Search API Error: "+result.Msg)
		log.Error("[result] SearchAll URL(%s) error[%v]", d.resultURL, err)
	}
	common = result.ResultResponse
	return
}

// SearchUgc gets the search tv_ugc result
func (d *Dao) SearchUgc(ctx context.Context, req *searchMdl.ReqSearch) (result searchMdl.RespUgc, common *searchMdl.ResultResponse, err error) {
	// common params
	params := commonParam(req)
	params.Set("category", fmt.Sprintf("%d", req.Category))
	if err = d.client.Get(ctx, d.resultURL, "", params, &result); err != nil {
		log.Error("[result] SearchUgc URL(%s) error[%v]", d.resultURL, err)
		return
	}
	if result.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(result.Code), "Search API Error: "+result.Msg)
		log.Error("[result] SearchUgc Code(%d) URL(%s) error[%v]", result.Code, d.resultURL, err)
	}
	common = result.ResultResponse
	return
}

// SearchPgc gets the search tv_pgc result
func (d *Dao) SearchPgc(ctx context.Context, req *searchMdl.ReqSearch) (result searchMdl.RespPgc, common *searchMdl.ResultResponse, err error) {
	params := commonParam(req)
	if err = d.client.Get(ctx, d.resultURL, "", params, &result); err != nil {
		log.Error("[result] SearchPgc URL(%s) error[%v]", d.resultURL, err)
		return
	}
	if result.Code != ecode.OK.Code() {
		log.Error("ClientGet Code Result Not OK [%v]", result)
		err = errors.Wrap(ecode.Int(result.Code), "Search API Error: "+result.Msg)
	}
	common = result.ResultResponse
	return
}

func commonParam(req *searchMdl.ReqSearch) (params url.Values) {
	params = url.Values{}
	params.Set("search_type", req.SearchType)
	params.Set("order", req.Order)
	params.Set("build", req.Build)
	params.Set("mobi_app", req.MobiAPP)
	params.Set("platform", req.Platform)
	params.Set("device", req.Device)
	params.Set("keyword", req.Keyword)
	params.Set("page", fmt.Sprintf("%d", req.Page))
	if req.Pagesize != 0 {
		params.Set("pagesize", fmt.Sprintf("%d", req.Pagesize))
	}
	return
}
