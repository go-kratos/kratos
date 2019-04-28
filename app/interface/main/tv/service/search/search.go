package search

import (
	"context"
	"strconv"

	"go-common/app/interface/main/tv/model"
	searchMdl "go-common/app/interface/main/tv/model/search"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_searchAll = "all_tv"
	_searchPGC = "tv_pgc"
	_searchUGC = "tv_ugc"
	_typePGC   = "pgc"
)

// SearchSug returns the result of search sug
func (s *Service) SearchSug(ctx context.Context, req *searchMdl.ReqSug) (result searchMdl.SugResponse, err error) {
	if result, err = s.dao.SearchSug(ctx, req); err != nil {
		return
	}
	build, _ := strconv.Atoi(req.Build)
	if build != 0 && build <= s.conf.Search.SugPGCBuild && len(result.Result.Tag) > 0 {
		var filtered = []*searchMdl.STag{}
		for _, v := range result.Result.Tag {
			if v.Type == _typePGC {
				filtered = append(filtered, v)
			}
		}
		result.Result.Tag = filtered
	}
	return
}

func (s *Service) batchToCommonPgc(ctx context.Context, input []*searchMdl.PgcResult) (output []*searchMdl.CommonResult) {
	var (
		err    error
		cids   []int64
		cmsRes map[int64]*model.SeasonCMS
	)
	for _, v := range input {
		output = append(output, v.ToCommon())
		cids = append(cids, int64(v.ID))
	}
	if cmsRes, err = s.cmsDao.LoadSnsCMSMap(ctx, cids); err != nil {
		log.Error("[search.cornerMark] cids(%s) error(%v)", xstr.JoinInts(cids), err)
		return
	}
	for idx, v := range output {
		if r, ok := cmsRes[int64(v.ID)]; ok && r.NeedVip() {
			output[idx].CornerMark = &(*s.conf.Cfg.SnVipCorner)
		}
	}
	return
}

func batchToCommonUgc(input []*searchMdl.UgcResult) (output []*searchMdl.CommonResult) {
	for _, v := range input {
		output = append(output, v.ToCommon())
	}
	return
}

// SearchRes distinguishes the search type and pick the result
func (s *Service) SearchRes(ctx context.Context, req *searchMdl.ReqSearch) (data *searchMdl.RespForClient, err error) {
	var resCommon *searchMdl.ResultResponse
	data = &searchMdl.RespForClient{
		SearchType: req.SearchType,
	}
	switch req.SearchType {
	case _searchAll:
		var resAll searchMdl.RespAll
		if resAll, resCommon, err = s.dao.SearchAll(ctx, req); err != nil {
			return
		}
		if resAll.PageInfo != nil {
			data.PageInfo = resAll.PageInfo
		}
		if resAll.Result != nil {
			data.ResultAll = &searchMdl.AllForClient{
				Pgc: s.batchToCommonPgc(ctx, resAll.Result.Pgc),
				Ugc: batchToCommonUgc(resAll.Result.Ugc),
			}
		}
	case _searchPGC:
		var resPgc searchMdl.RespPgc
		if resPgc, resCommon, err = s.dao.SearchPgc(ctx, req); err != nil {
			return
		}
		data.PGC = s.batchToCommonPgc(ctx, resPgc.Result)
	case _searchUGC:
		if req.Category == 0 { // in case of ugc, must have category
			err = ecode.RequestErr
			return
		}
		var resUgc searchMdl.RespUgc
		if resUgc, resCommon, err = s.dao.SearchUgc(ctx, req); err != nil {
			return
		}
		data.UGC = batchToCommonUgc(resUgc.Result)
	default:
		data = nil
		err = ecode.TvDangbeiWrongType
		return
	}
	data.ResultResponse = resCommon
	return
}
