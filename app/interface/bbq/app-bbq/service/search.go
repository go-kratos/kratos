package service

import (
	"context"
	"encoding/json"
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/model"
	"go-common/app/service/bbq/common"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/json-iterator/go"
)

const (
	searchTypeVideo = "bbq_video"
	searchTypeUser  = "bbq_user"
	suggestType     = "bbq"
	suggestVer      = "v4"
	searchLen       = 20
)

// SearchCursor .
type SearchCursor struct {
	Offset int64 `json:"offset"` // 默认为0，数据从1开始
}

func parseSearchCursor(c context.Context, cursorPrev, cursorNext string) (cursor *SearchCursor, directionNext bool, err error) {
	cursor = new(SearchCursor)
	// 判断是向前还是向后查询
	directionNext = true
	cursorStr := cursorNext
	if len(cursorNext) == 0 && len(cursorPrev) > 0 {
		directionNext = false
		cursorStr = cursorPrev
	}
	// 解析
	if len(cursorStr) != 0 {
		var cursorData = []byte(cursorStr)
		err = json.Unmarshal(cursorData, &cursor)
		if err != nil {
			err = ecode.ReqParamErr
			return
		}
	}
	return
}

//HotWord .
func (s *Service) HotWord(c context.Context, req *v1.HotWordRequest) (res *v1.HotWordResponse, err error) {
	res = new(v1.HotWordResponse)
	res.List = append(res.List, "小姐姐")
	res.List = append(res.List, "舞蹈")
	res.List = append(res.List, "MMD")
	return
}

// VideoSearch 视频搜索
func (s *Service) VideoSearch(c context.Context, arg *v1.BaseSearchReq) (res *v1.VideoSearchRes, err error) {
	res = new(v1.VideoSearchRes)
	var (
		sres      []*model.VideoSearchResult
		searchRes *model.RawSearchRes
		relids    []int32
		svids     []int64
		svidMap   map[int64]int64
	)

	// 0. check

	// 1. 生成请求，会根据请求参数
	// v2接口
	reqPage := int64(1)
	firstItemOffset := int64(0)
	var cursor *SearchCursor
	var directionNext bool
	if arg.Page == 0 {
		cursor, directionNext, err = parseSearchCursor(c, arg.CursorPrev, arg.CursorNext)
		if err != nil {
			log.Warnw(c, "log", "parse search cursor fail", "prev", arg.CursorPrev, "next", arg.CursorNext)
			return
		}
		firstItemOffset = cursor.Offset - 1
		if directionNext {
			firstItemOffset = cursor.Offset + 1
		}
		reqPage = (firstItemOffset-1)/searchLen + 1
		log.V(1).Infow(c, "log", "v2 request", "next", arg.CursorNext, "prev", arg.CursorPrev)
	} else {
		reqPage = arg.Page
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	sp := &model.SearchBaseReq{
		KeyWord:   arg.Key,
		Type:      searchTypeVideo,
		Page:      reqPage,
		PageSize:  searchLen,
		Highlight: int64(arg.Highlight),
	}
	// 相信搜索结果返回的就是searchLen
	searchRes, err = s.dao.SearchBBQ(c, sp)
	if searchRes.Res == nil {
		return
	}
	if err = json.Unmarshal(searchRes.Res, &sres); err != nil {
		return
	}
	for _, v := range sres {
		relids = append(relids, v.ID)
	}
	if len(relids) == 0 {
		log.Infov(c, log.KV("log", "search ret empty"))
		return
	}
	ids := s.dao.ParseRel2ID(relids)
	svidMap, err = s.dao.ConvID2SVID(c, ids)
	if err != nil {
		log.Error("ConvID2SVID err [%v]", err)
		err = nil
		return
	}
	for _, svid := range svidMap {
		svids = append(svids, svid)
	}
	svMap, _ := s.svInfos(c, svids, 0, false)

	for index, sinfo := range sres {
		id := s.dao.ParseRel2ID([]int32{sinfo.ID})
		if id[0] == 0 {
			continue
		}
		if svid, ok := svidMap[id[0]]; ok {
			if sv, ok := svMap[svid]; ok {
				if common.IsSearchSvStateAvailable(int64(sv.State)) {
					r := new(v1.VideoSearchList)
					r.VideoResponse = *sv
					r.SVID = svid
					r.HitColumns = sinfo.HitColumns
					r.TitleHighlight = sinfo.Title
					r.Offset = (reqPage-1)*int64(searchLen) + int64(index) + 1
					res.List = append(res.List, r)
				} else {
					log.Warnw(c, "log", "get error svid in search list", "svid", svid, "sv", sv)
				}
			}
		}
	}

	res.NumPage = searchRes.PageNum
	res.Page = searchRes.Page

	res.HasMore = false
	if arg.Page == 0 {
		if (directionNext && searchRes.PageNum > searchRes.Page) || (!directionNext && searchRes.Page > 1) {
			res.HasMore = true
		}

		var list []*v1.VideoSearchList
		if directionNext {
			for _, item := range res.List {
				if item.Offset < firstItemOffset {
					continue
				}
				list = append(list, item)
			}
		} else {
			if firstItemOffset >= int64(len(res.List)) {
				firstItemOffset = int64(len(res.List)) - 1
			}
			for i := firstItemOffset - 1; i >= 0; i-- {
				list = append(list, res.List[i])
			}
		}
		res.List = list
	}

	var cursorValue SearchCursor
	for _, item := range res.List {
		cursorValue.Offset = item.Offset
		data, _ := json.Marshal(cursorValue)
		item.CursorValue = string(data)
	}

	return
}

// UserSearch 用户搜索
func (s *Service) UserSearch(c context.Context, mid int64, arg *v1.BaseSearchReq) (res *v1.UserSearchRes, err error) {
	res = new(v1.UserSearchRes)
	var (
		sres      []*model.UserSearchResult
		searchRes *model.RawSearchRes
		mids      []int64
	)

	// 1. 生成请求，会根据请求参数
	// v2接口
	reqPage := int64(1)
	firstItemOffset := int64(0)
	var cursor *SearchCursor
	var directionNext bool
	if arg.Page == 0 {
		cursor, directionNext, err = parseSearchCursor(c, arg.CursorPrev, arg.CursorNext)
		if err != nil {
			log.Warnw(c, "log", "parse search cursor fail", "prev", arg.CursorPrev, "next", arg.CursorNext)
			return
		}
		firstItemOffset = cursor.Offset - 1
		if directionNext {
			firstItemOffset = cursor.Offset + 1
		}
		reqPage = (firstItemOffset-1)/searchLen + 1
		log.V(1).Infow(c, "log", "v2 request", "next", arg.CursorNext, "prev", arg.CursorPrev)
	} else {
		reqPage = arg.Page
	}

	sp := &model.SearchBaseReq{
		KeyWord:   arg.Key,
		Type:      searchTypeUser,
		Page:      reqPage,
		PageSize:  searchLen,
		Highlight: int64(arg.Highlight),
	}
	searchRes, err = s.dao.SearchBBQ(c, sp)
	if searchRes.Res == nil {
		return
	}
	if err = json.Unmarshal(searchRes.Res, &sres); err != nil {
		return
	}
	for _, v := range sres {
		mids = append(mids, v.ID)
	}
	if len(mids) == 0 {
		log.Infov(c, log.KV("log", "search ret empty"))
		return
	}
	uMap, _ := s.dao.BatchUserInfo(c, mid, mids, false, true, true)
	for index, sinfo := range sres {
		if u, ok := uMap[sinfo.ID]; ok {
			r := &v1.UserSearchList{
				UserInfo: *u,
			}
			st := new(v1.UserStatic)
			st.Fan = u.Fan
			st.Follow = u.Follow
			st.Like = u.Like
			st.Liked = u.Liked
			st.FollowState = u.FollowState
			r.UserStatic = st
			r.HitColumns = sinfo.HitColumns
			r.UnameHighlight = sinfo.Uname
			r.Offset = (reqPage-1)*int64(searchLen) + int64(index) + 1
			res.List = append(res.List, r)
		}
	}
	res.NumPage = searchRes.PageNum
	res.Page = searchRes.Page

	res.HasMore = false
	if arg.Page == 0 {
		if (directionNext && searchRes.PageNum > searchRes.Page) || (!directionNext && searchRes.Page > 1) {
			res.HasMore = true
		}

		var list []*v1.UserSearchList
		if directionNext {
			for _, item := range res.List {
				if item.Offset < firstItemOffset {
					continue
				}
				list = append(list, item)
			}
		} else {
			if firstItemOffset >= int64(len(res.List)) {
				firstItemOffset = int64(len(res.List)) - 1
			}
			for i := firstItemOffset - 1; i >= 0; i-- {
				list = append(list, res.List[i])
			}
		}
		res.List = list
	}

	var cursorValue SearchCursor
	for _, item := range res.List {
		cursorValue.Offset = item.Offset
		data, _ := json.Marshal(cursorValue)
		item.CursorValue = string(data)
	}

	return
}

// BBQSug BBQSUG服务
func (s *Service) BBQSug(c context.Context, arg *v1.SugReq) (res []*v1.SugTag, err error) {
	var (
		sres []*model.RawSugTag
		data json.RawMessage
	)
	sp := &model.SugBaseReq{
		Term:        arg.KeyWord,
		SuggestType: suggestType,
		MainVer:     suggestVer,
		SugNum:      arg.PageSize,
		Highlight:   int64(arg.Highlight),
	}
	data, err = s.dao.SugBBQ(c, sp)
	if err != nil {
		log.Error("s.dao.SugBBQ err[%v] data [%s]", err, string(data))
		return
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err = json.Unmarshal(data, &sres); err != nil {
		log.Error("json.Unmarshal err[%v] data [%s]", err, string(data))
		return
	}
	for _, sinfo := range sres {
		r := new(v1.SugTag)
		name := sinfo.Name
		r.Name = name
		r.Value = sinfo.Value
		r.Type = sinfo.Type
		r.Ref = sinfo.Ref
		res = append(res, r)
	}
	return
}
