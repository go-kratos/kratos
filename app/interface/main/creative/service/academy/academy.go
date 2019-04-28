package academy

import (
	"context"
	"time"

	"go-common/app/interface/main/creative/model/academy"
	"go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/api"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"

	"github.com/davecgh/go-spew/spew"
)

// TagList get all tag.
func (s *Service) TagList(c context.Context) (res map[string][]*academy.Tag, err error) {
	res = s.TagsCache
	return
}

// AddFeedBack add feedback.
func (s *Service) AddFeedBack(c context.Context, category, course, suggest string, mid int64) (id int64, err error) {
	fb := &academy.FeedBack{
		Category: category,
		Course:   course,
		Suggest:  suggest,
		CTime:    xtime.Time(time.Now().Unix()),
		MTime:    xtime.Time(time.Now().Unix()),
	}
	if id, err = s.aca.AddFeedBack(c, fb, mid); err != nil {
		log.Error("s.aca.AddFeedBack error(%v)", err)
	}
	return
}

// ArchivesWithES get all archive by es.
func (s *Service) ArchivesWithES(c context.Context, aca *academy.EsParam) (res *academy.ArchiveList, err error) {
	var sear *academy.SearchResult
	res = &academy.ArchiveList{
		Items: []*academy.ArchiveMeta{},
		Page:  &academy.ArchivePage{},
	}

	aca.TidsMap = s.filterTIDs(aca.Tid)
	if sear, err = s.aca.ArchivesWithES(c, aca); err != nil {
		log.Error("s.aca.ArchivesWithES sear(%+v)|param(%+v)|error(%v)", sear, aca, err)
		return
	}
	if sear == nil || len(sear.Result) == 0 {
		log.Error("s.aca.ArchivesWithES has no data sear(%+v)|param(%+v)|error(%v)", sear, aca, err)
		return
	}
	res.Page.Total = sear.Page.Total
	res.Page.Pn = sear.Page.Num
	res.Page.Ps = sear.Page.Size

	var searRes []*academy.EsArc
	if aca.Keyword != "" { //搜索关键词红点
		searRes = make([]*academy.EsArc, 0, len(sear.Result)/2)
		for i := 0; i < len(sear.Result)-1; i += 2 {
			sear.Result[i].Title = sear.Result[i+1].Title
			searRes = append(searRes, sear.Result[i])
		}
	} else {
		searRes = sear.Result
	}

	oids := make([]int64, 0, len(sear.Result))
	aidTIDsMap := make(map[int64][]int64)
	busAIDsMap := make(map[int][]int64)
	busAidMap := make(map[int64]int)
	highTitleMap := make(map[int64][]string)
	for _, v := range searRes {
		busAIDsMap[v.Business] = append(busAIDsMap[v.Business], v.OID)
		oids = append(oids, v.OID)
		aidTIDsMap[v.OID] = v.TID
		busAidMap[v.OID] = v.Business
		if aca.Keyword != "" { //搜索关键词红点
			highTitleMap[v.OID] = v.Title
		}
	}

	var (
		g, _    = errgroup.WithContext(c)
		tagInfo map[int64]map[string][]*academy.Tag
		arcs    map[int64]*api.Arc
		arts    map[int64]*model.Meta
		st      map[int64]*api.Stat
	)

	g.Go(func() error { //获取各种查询对象信息
		switch aca.Business {
		case academy.BusinessForAll: //查询所有
			if ids, ok := busAIDsMap[academy.BusinessForArchive]; ok { //稿件
				g.Go(func() error {
					arcs, err = s.arc.Archives(c, ids, aca.IP)
					if err != nil {
						log.Error("s.arc.Archives oids(%+v)|business(%d)|error(%v)", ids, aca.Business, err)
						return err
					}
					st, err = s.arc.Stats(c, ids, aca.IP)
					if err != nil {
						log.Error("s.arc.Stats oids(%+v)|business(%d)|error(%v)", ids, aca.Business, err)
					}
					return err
				})
			}
			if ids, ok := busAIDsMap[academy.BusinessForArticle]; ok { //文章
				g.Go(func() error {
					arts, err = s.art.ArticleMetas(context.Background(), ids, aca.IP)
					if err != nil {
						log.Error("s.arc.ArticleMetas oids(%+v)|business(%d)|error(%v)", ids, aca.Business, err)
					}
					return err
				})
			}
		case academy.BusinessForArchive: //稿件
			arcs, err = s.arc.Archives(context.Background(), oids, aca.IP)
			if err != nil {
				log.Error("s.arc.Archives oids(%+v)|business(%d)|error(%v)", oids, aca.Business, err)
				return err
			}
			st, err = s.arc.Stats(c, oids, aca.IP)
			if err != nil {
				log.Error("s.arc.Stats oids(%+v)|business(%d)|error(%v)", oids, aca.Business, err)
			}
			return err
		case academy.BusinessForArticle: //文章
			arts, err = s.art.ArticleMetas(context.Background(), oids, aca.IP)
			if err != nil {
				log.Error("s.arc.ArticleMetas oids(%+v)|business(%d)|error(%v)", oids, aca.Business, err)
			}
			return err
		}
		return nil
	})

	g.Go(func() error {
		tagInfo, err = s.bindTags(c, aidTIDsMap)
		return err
	})

	if g.Wait() != nil {
		log.Error("s.aca.ArchivesWithES g.Wait() error(%v)", err)
		return
	}

	items := make([]*academy.ArchiveMeta, 0, len(oids))
	for _, oid := range oids {

		a := &academy.ArchiveMeta{
			OID: oid,
		}

		if v, ok := tagInfo[oid]; ok {
			a.Tags = v
		}

		bs, ok := busAidMap[oid]
		if !ok {
			log.Error("s.aca.ArchivesWithES oid(%d) get invalid business", oid)
			return
		}
		a.Business = bs

		switch a.Business {
		case academy.BusinessForArchive: //稿件
			a = bindArchiveInfo(oid, arcs, a)
			if t, ok := st[oid]; ok {
				a.ArcStat = t
			} else {
				a.ArcStat = &api.Stat{}
			}
		case academy.BusinessForArticle: //文章
			a = bindArticleInfo(oid, arts, a)
		}

		if aca.Keyword != "" {
			if ht, ok := highTitleMap[oid]; ok && len(ht) > 0 {
				a.HighLightTitle = ht[0]
			}
		}
		items = append(items, a)
	}
	res.Items = items
	return
}

func bindArchiveInfo(oid int64, arcs map[int64]*api.Arc, a *academy.ArchiveMeta) (res *academy.ArchiveMeta) {
	if v, ok := arcs[oid]; ok {
		a.Title = v.Title
		a.State = v.State
		a.Type = v.TypeName
		a.Cover = v.Pic
		a.UName = v.Author.Name
		a.Face = v.Author.Face
		a.MID = v.Author.Mid
		a.Duration = v.Duration
		a.Rights = v.Rights
	}
	res = a
	return
}

func bindArticleInfo(oid int64, arts map[int64]*model.Meta, a *academy.ArchiveMeta) (res *academy.ArchiveMeta) {
	if v, ok := arts[oid]; ok && v != nil {
		a.Title = v.Title
		a.State = v.State
		a.MID = v.Author.Mid
		a.Comment = v.Summary
		if v.Category != nil {
			a.Type = v.Category.Name
		}
		if len(v.ImageURLs) > 0 {
			a.Cover = v.ImageURLs[0]
		}
		if v.Author != nil {
			a.UName = v.Author.Name
			a.Face = v.Author.Face
		}
		if v.Stats != nil {
			a.ArtStat = v.Stats
		} else {
			a.ArtStat = &model.Stats{}
		}
	}
	res = a
	return
}

func (s *Service) bindTags(c context.Context, tidsMap map[int64][]int64) (res map[int64]map[string][]*academy.Tag, err error) {
	res = make(map[int64]map[string][]*academy.Tag)
	for oid, tids := range tidsMap {
		tgs := s.getTagsByTIDs(tids)
		if len(tgs) == 0 {
			continue
		}
		oidTag := make(map[string][]*academy.Tag)
		ctgs := make(map[int64][]*academy.Tag)
		for _, tg := range tgs {
			k := academy.TagClassMap(int(tg.Type))
			if tg.Type == academy.Classify { //获取多个分类标签
				ctgs[tg.ParentID] = append(ctgs[tg.ParentID], tg)
			} else {
				oidTag[k] = append(oidTag[k], tg)
			}
		}
		for pid, tgs := range ctgs {
			if p, ok := s.TagMapCache[pid]; ok {
				tp := *p
				tp.Children = tgs
				oidTag[academy.TagClassMap(academy.Classify)] = append(oidTag[academy.TagClassMap(academy.Classify)], &tp)
			}
		}
		res[oid] = oidTag
	}
	return
}

func (s *Service) filterTIDs(tids []int64) (res map[int][]int64) {
	if len(tids) == 0 {
		return
	}
	log.Info("s.filterTIDs origin tids(%+v)", tids)
	res = make(map[int][]int64)
	ochs := make([]int64, 0) //原始提交的二级标签
	ops := make([]int64, 0)  //原始提交的一级标签
	qchs := make([]int64, 0) //通过一级标签查询出来的二级标签
	for _, id := range tids {
		t, ok := s.parentChildMapCache[id]
		if !ok || t == nil {
			continue
		}
		if t.Type == academy.Classify {
			if t.ParentID != 0 { //原始提交的二级标签
				ochs = append(ochs, id)
			} else if t.ParentID == 0 && len(t.Children) > 0 { //通过一级标签查询出来的二级标签
				for _, v := range t.Children {
					qchs = append(qchs, v.ID)
				}
			} else if t.ParentID == 0 && len(t.Children) == 0 {
				ops = append(ops, id)
			}
		} else {
			res[int(t.Type)] = append(res[int(t.Type)], id)
		}
	}
	if len(ochs) > 0 { //如果分类标签中提交了原始的二级标签则认为按该二级标签进行筛选，如果可以查询到二级标签认为筛选全部二级，否则一级参与查询.
		res[academy.Classify] = ochs
	} else if len(qchs) > 0 {
		res[academy.Classify] = qchs
	} else if len(ops) > 0 {
		res[academy.Classify] = ops
	}
	log.Info("s.filterTIDs res(%s)", spew.Sdump(res))
	return
}

func (s *Service) getTagsByTIDs(tids []int64) (res []*academy.Tag) {
	res = make([]*academy.Tag, 0)
	if len(tids) == 0 {
		return
	}
	for _, tid := range tids {
		tag, ok := s.TagMapCache[tid]
		if !ok || tag == nil {
			continue
		}
		res = append(res, tag)
	}
	return
}
