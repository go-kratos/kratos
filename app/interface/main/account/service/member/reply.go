package member

import (
	"context"
	"fmt"

	"go-common/app/interface/main/account/model"
	artMdl "go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/api"
	arcMdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

var (
	_typeToURL = map[int64]string{
		1:  "https://www.bilibili.com/video/av%d/",                    // 稿件
		2:  "https://www.bilibili.com/topic/%d.html",                  // 话题
		3:  "https://h.bilibili.com/dy%d",                             // 画站
		5:  "https://vc.bilibili.com/video/%d",                        // 直播小视频
		6:  "https://www.bilibili.com/blackroom/ban/%d",               // 封禁信息
		7:  "https://www.bilibili.com/blackroom/notice/%d",            // 公告信息
		10: "https://link.bilibili.com/p/eden/news#/newsdetail?id=%d", // 直播公告
		11: "https://h.bilibili.com/%d",                               // 直播有文画
		12: "https://www.bilibili.com/read/cv%d",                      // 专栏
		13: "https://show.bilibili.com/platform/detail.html?id=%d",    // 票务
		15: "https://www.bilibili.com/judgement/case/%d",              // 风纪委
	}
)

// ReplyHistoryList reply history list
func (s *Service) ReplyHistoryList(c context.Context, mid int64, stime, etime, order, sort string, pn, ps int64, accessKey, cookie string) (rhl *model.ReplyHistory, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	if rhl, err = s.replyDao.ReplyHistoryList(c, mid, stime, etime, order, sort, pn, ps, accessKey, cookie, ip); err != nil {
		log.Error("s.replyDao.ReplyHistoryList error(%v)", err)
		return
	}
	idsMap := make(map[int64][]int64) // type -> ids
	unique := make(map[int64]struct{})
	for _, v := range rhl.Records {
		if _, ok := unique[v.Oid]; !ok {
			idsMap[v.Type] = append(idsMap[v.Type], v.Oid)
			unique[v.Oid] = struct{}{}
		}
	}
	rhlt, _ := s.fetchData(c, mid, idsMap, accessKey, cookie, ip)
	for _, v := range rhl.Records {
		if t, ok := rhlt[v.Type]; ok {
			for _, b := range t {
				if o, ok := b[v.Oid]; ok {
					v.Title = o.Title
					v.URL = o.URL
				}
			}
		}
	}
	return
}

// fetchData
func (s *Service) fetchData(c context.Context, mid int64, idsMap map[int64][]int64, accessKey, cookie, ip string) (rhlt map[int64][]map[int64]*model.RecordAppend, err error) {
	rhlt = make(map[int64][]map[int64]*model.RecordAppend) // type -> oid -> title/url
	for t, v := range idsMap {
		switch t {
		case 1:
			// 稿件
			if len(v) > 0 {
				var arcs map[int64]*api.Arc
				arcArg := &arcMdl.ArgAids2{Aids: v, RealIP: ip}
				if arcs, err = s.arcRPC.Archives3(c, arcArg); err != nil {
					log.Error("s.arcRPC.Archives3 error(%v)", err)
					return
				}
				for _, vv := range v {
					if arc, ok := arcs[vv]; ok {
						itu := &model.RecordAppend{
							Title: arc.Title,
						}
						if arc.RedirectURL != "" {
							itu.URL = arc.RedirectURL
						} else {
							itu.URL = fmt.Sprintf(_typeToURL[t], arc.Aid)
						}
						vitu := make(map[int64]*model.RecordAppend)
						vitu[vv] = itu
						rhlt[t] = append(rhlt[t], vitu)
					}
				}
			}
		case 4:
			// 活动
			if len(v) > 0 {
				var aps map[int64]*model.RecordAppend
				if aps, err = s.replyDao.ActivityPages(c, mid, v, accessKey, cookie, ip); err != nil {
					log.Error("s.replyDao.ActivityPages error(%v)", err)
					return
				}
				for _, vv := range v {
					if ap, ok := aps[vv]; ok {
						vitu := make(map[int64]*model.RecordAppend)
						vitu[vv] = ap
						rhlt[t] = append(rhlt[t], vitu)
					}
				}
			}
		case 12:
			// 专栏
			if len(v) > 0 {
				var arts map[int64]*artMdl.Meta
				artArg := &artMdl.ArgAids{Aids: v}
				if arts, err = s.artRPC.ArticleMetas(c, artArg); err != nil {
					log.Error("s.artRPC.ArticleMetas error(%v)", err)
					return
				}
				for _, vv := range v {
					if ap, ok := arts[vv]; ok {
						itu := &model.RecordAppend{
							Title: ap.Title,
							URL:   fmt.Sprintf(_typeToURL[t], ap.ID),
						}
						vitu := make(map[int64]*model.RecordAppend)
						vitu[vv] = itu
						rhlt[t] = append(rhlt[t], vitu)
					}
				}
			}
		default:
			if len(v) > 0 {
				for _, vv := range v {
					itu := &model.RecordAppend{
						Title: "",
						URL:   "",
					}
					if _, ok := _typeToURL[t]; ok {
						itu.URL = fmt.Sprintf(_typeToURL[t], vv)
					}
					vitu := make(map[int64]*model.RecordAppend)
					vitu[vv] = itu
					rhlt[t] = append(rhlt[t], vitu)
				}
			}
		}
	}
	return
}
