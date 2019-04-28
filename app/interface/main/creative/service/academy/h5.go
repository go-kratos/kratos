package academy

import (
	"context"
	"sort"
	"strconv"
	"time"

	"go-common/app/interface/main/creative/dao/tool"
	"go-common/app/interface/main/creative/model/academy"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"

	"go-common/app/service/main/archive/api"

	"github.com/davecgh/go-spew/spew"
)

// Tags get all h5 tags.
func (s *Service) Tags(c context.Context) (res []*academy.Tag) {
	if v, ok := s.TagsCache[academy.TagClassMap(academy.H5)]; ok {
		res = v
	}
	return
}

// Archives get all h5 archive.
func (s *Service) Archives(c context.Context, aca *academy.EsParam) (res *academy.ArchiveList, err error) {
	tids, err := s.webTags(aca.Tid)
	if err != nil {
		return
	}
	aca.Tid = tids

	res, err = s.ArchivesWithES(c, aca)
	return
}

func (s *Service) webTags(tids []int64) (webTIDs []int64, err error) {
	var (
		lts    []*academy.LinkTag
		h5TIDs []int64
	)

	if len(tids) > 0 {
		h5TIDs = tids
	} else {
		tgs := s.Tags(context.Background())
		for _, v := range tgs {
			h5TIDs = append(h5TIDs, v.ID)
		}
	}

	lts, err = s.aca.LinkTags(context.Background(), h5TIDs)
	if err != nil {
		return
	}

	for _, v := range lts {
		webTIDs = append(webTIDs, v.LinkID)
	}
	return
}

// RecommendV2 get recommend archive.
func (s *Service) RecommendV2(c context.Context, mid int64) (res []*academy.RecArcList, err error) {
	mainIDMap := make(map[int64]struct{}) //主题课程aid map for rm dup
	tgList, err := s.getRecTag(c, mid)
	if err != nil {
		return
	}
	log.Info("Recommend mid(%d)|tgList(%s)", mid, spew.Sdump(tgList))
	s.setSeed() //init Seed.
	res = make([]*academy.RecArcList, 0)
	var (
		g, _     = errgroup.WithContext(c)
		hotItems []*academy.RecArchive
	)

	g.Go(func() error {
		// get hot archives
		hotItems, err = s.hotArchives(c)
		if err != nil {
			log.Error("Recommend s.hotArchives mid(%d)", mid)
			return err
		}
		return nil
	})

	for _, t := range tgList {
		if t == nil {
			continue
		}
		pid, v := t.PID, t.TIDs

		if pid == 0 { //主题课程
			var ocid int64
			if len(v) > 0 {
				ocid = v[0]
			}
			rec := &academy.RecArcList{
				TID:   ocid,
				Items: []*academy.RecArchive{},
			}
			tg, o := s.OccMapCache[ocid]
			if !o || tg == nil {
				log.Error("s.OccMapCache ocid(%d) not exist", ocid)
				continue
			}
			rec.Name = tg.Name
			items, themeCourseErr := s.themeCourse(c, v)
			if themeCourseErr != nil {
				return nil, themeCourseErr
			}
			for _, v := range items {
				mainIDMap[v.OID] = struct{}{}
			}
			rec.Items = items
			res = append(res, rec)
		} else if tg, ok := s.TagMapCache[pid]; ok { //标签教程
			rec := &academy.RecArcList{
				TID:   pid,
				Name:  tg.Name,
				Items: []*academy.RecArchive{},
			}

			items, tagCourseErr := s.tagCourse(c, pid, v, mainIDMap)
			if tagCourseErr != nil {
				return nil, tagCourseErr
			}
			rec.Items = items
			res = append(res, rec)
		}
	}

	if g.Wait() != nil {
		log.Error("Recommend s.hotArchives g.Wait() error(%v)", err)
		return
	}
	// add host archives
	hotRec := &academy.RecArcList{
		TID:   0,
		Name:  "热门推荐",
		Items: hotItems,
	}
	res = append(res, hotRec)

	return
}

func (s *Service) tagCourse(c context.Context, pid int64, v []int64, aidMap map[int64]struct{}) (res []*academy.RecArchive, err error) {
	res = make([]*academy.RecArchive, 0)
	aca := &academy.EsParam{
		Tid: v,
		Pn:  1,
		Ps:  10,
	}

	if s.Seed > 0 { //取材创意/视频制作/个人运营 每日0点请求搜索更换时间种子
		aca.Seed = s.Seed
	}

	arcs, err := s.Archives(c, aca)
	if err != nil {
		log.Error("Recommend s.Archives EsParam(%+v)|error(%v)", aca, err)
		return nil, err
	}

	if arcs == nil {
		err = ecode.CreativeAcademyH5RecommendErr
		return nil, err
	}

	var aids []int64
	for _, i := range arcs.Items {
		// ignore if exist in resource service
		if _, exist := s.ResourceMapCache[i.OID]; exist {
			continue
		}
		// ignore if exist in main topic
		if _, exist := aidMap[i.OID]; exist {
			continue
		}
		ra := &academy.RecArchive{
			OID:      i.OID,
			MID:      i.MID,
			Cover:    i.Cover,
			Title:    i.Title,
			Business: i.Business,
			Duration: i.Duration,
			ArcStat:  i.ArcStat,
			ArtStat:  i.ArtStat,
		}
		res = append(res, ra)
		aids = append(aids, i.OID)
	}

	// add tags
	tags, err := s.getTags(c, aids)
	if err != nil {
		log.Error("tagCourse s.getTags err(%v)", err)
		return
	}
	s.setTags(res, tags)
	return
}

func (s *Service) themeCourse(c context.Context, v []int64) (res []*academy.RecArchive, err error) {
	res = make([]*academy.RecArchive, 0)
	arcs, err := s.ThemeCourse(c, v, []int64{}, []int64{}, 1, 10, false)
	if err != nil {
		log.Error("Recommend s.ThemeCourse v(%+v)|error(%v)", v, err)
		return nil, err
	}

	if arcs == nil {
		err = ecode.CreativeAcademyH5RecommendErr
		return nil, err
	}

	var aids []int64
	for _, i := range arcs.Items {
		// ignore if exist in resource service
		if _, exist := s.ResourceMapCache[i.AID]; exist {
			continue
		}
		ra := &academy.RecArchive{
			OID:      i.AID,
			MID:      i.MID,
			Cover:    i.Cover,
			Title:    i.Title,
			Duration: i.Duration,
			ArcStat:  i.ArcStat,
			Business: academy.BusinessForArchive,
		}
		res = append(res, ra)
		aids = append(aids, i.AID)
	}

	// add tags
	tags, err := s.getTags(c, aids)
	if err != nil {
		log.Error("themeCourse s.getTags err(%v)", err)
		return
	}
	s.setTags(res, tags)

	s.randomForMainCourse(res) //每日0点随机随机稿件列表
	if len(s.RecommendArcs) > 0 {
		res = s.RecommendArcs
	}
	return
}

func (s *Service) getRecTag(c context.Context, mid int64) (res []*academy.RecConf, err error) {
	var tyID int64
	if mid > 0 {
		tyID, err = s.getFavType(c, mid) //获取推荐分区id
		if err != nil {
			log.Error("getFavType mid(%d)|error(%v)", mid, err)
		} else {
			log.Info("getFavType mid(%d)|tyID(%d)", mid, tyID)
		}
	}

	if s.c == nil || s.c.AcaRecommend == nil {
		log.Error("getRecTag get conf error mid(%d)", mid)
		return
	}
	rec := s.c.AcaRecommend.Recommend

	//按 主题课程-取材创意-视频制作-个人运营  排序
	var rec1, rec2, rec3, rec4 *academy.RecConf
	res = make([]*academy.RecConf, 0, 4)
	//主题课程
	if rec.Course != nil {
		course := rec.Course
		rec1 = &academy.RecConf{PID: course.ID}
		if tyID != 0 {
			if tool.ElementInSlice(tyID, course.Shoot.Val) { //如果最近投稿分区命中配置的分区，则设置当前一级分类下面的标签为最近投稿分区对应的二级标签目录
				rec1.TIDs = course.Shoot.Key

			} else if tool.ElementInSlice(tyID, course.Scene.Val) {
				rec1.TIDs = course.Scene.Key

			} else if tool.ElementInSlice(tyID, course.Edit.Val) {
				rec1.TIDs = course.Edit.Key

			} else if tool.ElementInSlice(tyID, course.Mmd.Val) {
				rec1.TIDs = course.Mmd.Key

			} else if tool.ElementInSlice(tyID, course.Sing.Val) {
				rec1.TIDs = course.Sing.Key

			} else if tool.ElementInSlice(tyID, course.Bang.Val) {
				rec1.TIDs = course.Bang.Key
			}
		} else {
			rec1.TIDs = course.Other.Key
		}
		res = append(res, rec1)
	} else {
		log.Error("getRecTag get cousre conf mid(%d)", mid)
	}

	//取材创意
	if rec.Drawn != nil {
		drawn := rec.Drawn
		rec2 = &academy.RecConf{PID: drawn.ID}
		if tyID != 0 {
			if tool.ElementInSlice(tyID, drawn.MobilePlan.Val) { //如果最近投稿分区命中配置的分区，则设置当前一级分类下面的标签为最近投稿分区对应的二级标签目录
				rec2.TIDs = drawn.MobilePlan.Key
			} else if tool.ElementInSlice(tyID, drawn.ScreenPlan.Val) {
				rec2.TIDs = drawn.ScreenPlan.Key
			} else if tool.ElementInSlice(tyID, drawn.RecordPlan.Val) {
				rec2.TIDs = drawn.RecordPlan.Key
			}
		} else {
			rec2.TIDs = drawn.Other.Key
		}
		res = append(res, rec2)
	} else {
		log.Error("getRecTag get drawn conf mid(%d)", mid)
	}

	//视频制作
	if rec.Video != nil {
		video := rec.Video
		rec3 = &academy.RecConf{PID: video.ID}
		if tyID != 0 {
			if tool.ElementInSlice(tyID, video.MobileMake.Val) { //如果最近投稿分区命中配置的分区，则设置当前一级分类下面的标签为最近投稿分区对应的二级标签目录
				rec3.TIDs = video.MobileMake.Key
			} else if tool.ElementInSlice(tyID, video.AudioEdit.Val) {
				rec3.TIDs = video.AudioEdit.Key
			} else if tool.ElementInSlice(tyID, video.EditCompose.Val) {
				rec3.TIDs = video.EditCompose.Key
			}
		} else {
			rec3.TIDs = video.Other.Key
		}
		res = append(res, rec3)
	} else {
		log.Error("getRecTag get video conf mid(%d)", mid)
	}

	//个人运营
	if rec.Person != nil {
		person := rec.Person
		rec4 = &academy.RecConf{PID: person.ID, TIDs: person.Other.Key}
		res = append(res, rec4)
	} else {
		log.Error("getRecTag get person conf mid(%d)", mid)
	}

	return
}

func (s *Service) getFavType(c context.Context, mid int64) (tyID int64, err error) { //获取最近投稿的一个分区
	tys, err := s.arc.FavTypes(c, mid)
	if err != nil {
		log.Error("s.arc.FavTypes mid(%d)|error(%v)", mid, err)
		return
	}

	if len(tys) == 0 {
		return 0, nil
	}

	type kv struct {
		id    int64
		ptime int64
	}
	var tps []*kv

	for id, t := range tys {
		tid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return 0, err
		}
		tps = append(tps, &kv{tid, t})
	}

	sort.Slice(tps, func(i, j int) bool {
		return tps[i].ptime > tps[j].ptime
	})
	if len(tps) > 0 && tps[0] != nil {
		tyID = tps[0].id
	}
	return
}

//randomForMainCourse 主题课程每日0点重新随机排序
func (s *Service) randomForMainCourse(arc []*academy.RecArchive) {
	count := len(arc)
	if count == 0 {
		return
	}
	keys := tool.RandomSliceKeys(0, count, count, s.Seed)

	res := make([]*academy.RecArchive, 0, count)
	for _, k := range keys {
		res = append(res, arc[k])
	}

	if len(res) > 0 { //获取随机排序的稿件列表
		s.RecommendArcs = res
	}
	log.Info("randomRecommend s.RecommendArcs (%s)", spew.Sdump(s.RecommendArcs))
}

func (s *Service) setSeed() {
	now := time.Now()
	last := now
	next := now.Add(time.Hour * 24)

	last = time.Date(last.Year(), last.Month(), last.Day(), 0, 0, 0, 0, last.Location()) //昨日0点
	next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location()) //明日0点

	if now.Unix() > last.Unix() && now.Unix() < next.Unix() {
		s.Seed = last.Unix() //set last seed
	} else {
		s.Seed = next.Unix() //set next seed
	}
	log.Info("setSeed s.Seed (%d)", s.Seed)
}

func (s *Service) getTags(c context.Context, aids []int64) (res map[int64]map[string][]*academy.Tag, err error) {
	if len(aids) == 0 {
		log.Error("getTags len(aids) == 0")
		return
	}

	aidTIDsMap, err := s.aca.ArchiveTagsByOids(c, aids)
	if err != nil {
		log.Error("getTags s.aca.ArchiveTagsByOids aids(%+v)", aids)
		return
	}
	if len(aidTIDsMap) == 0 {
		log.Error("getTags len(aidTIDsMap) == 0 | aids(%+v)", aids)
		return
	}
	res, err = s.bindTags(c, aidTIDsMap)
	if err != nil {
		log.Error("getTags s.bindTags | err(%v)", err)
		return
	}
	return
}

func (s *Service) setTags(x interface{}, tags map[int64]map[string][]*academy.Tag) {
	switch arcs := x.(type) {
	case []*academy.RecArchive:
		for _, v := range arcs {
			if v != nil {
				if tag, ok := tags[v.OID]; ok {
					v.Tags = tag
				}
			}
		}
	case []*academy.ArcMeta:
		for _, v := range arcs {
			if v != nil {
				if tag, ok := tags[v.AID]; ok {
					v.Tags = tag
				}
			}
		}
	}
}

// HotArchives get host archives
func (s *Service) HotArchives(c context.Context, oids []int64) (res []*academy.ArchiveMeta, err error) {
	res = make([]*academy.ArchiveMeta, 0)
	if len(oids) == 0 {
		log.Error("HotArchives len(oids) == 0")
		return
	}
	var (
		g, _ = errgroup.WithContext(c)
		arcs map[int64]*api.Arc
		st   map[int64]*api.Stat
	)

	g.Go(func() error {
		arcs, err = s.arc.Archives(c, oids, "")
		if err != nil {
			log.Error("HotArchives s.arc.Archives oids(%+v)| error(%v)", oids, err)
			return err
		}
		st, err = s.arc.Stats(c, oids, "")
		if err != nil {
			log.Error("HotArchives s.arc.Stats oids(%+v)| error(%v)", oids, err)
			return err
		}
		return nil
	})
	if g.Wait() != nil {
		log.Error("HotArchives g.Wait() error(%v)", err)
		return
	}

	for _, oid := range oids {
		a := &academy.ArchiveMeta{
			OID: oid,
		}
		a = bindArchiveInfo(oid, arcs, a)
		if t, ok := st[oid]; ok {
			a.ArcStat = t
		} else {
			a.ArcStat = &api.Stat{}
		}
		res = append(res, a)
	}
	return
}

func (s *Service) hotArchives(c context.Context) (res []*academy.RecArchive, err error) {
	if len(s.ResourceMapCache) == 0 {
		log.Error("hotArchives len(oids) == 0 | ResourceMapCache(%+v)", s.ResourceMapCache)
		return
	}
	res = make([]*academy.RecArchive, 0)
	var oids []int64

	for k := range s.ResourceMapCache {
		oids = append(oids, k)
	}

	arcs, err := s.HotArchives(c, oids)
	if err != nil {
		log.Error("hotArchives s.HotArchives oids(%+v) | error(%v)", oids, err)
		return
	}

	var aids []int64
	for _, v := range arcs {
		ra := &academy.RecArchive{
			OID:      v.OID,
			MID:      v.MID,
			Cover:    v.Cover,
			Title:    v.Title,
			Business: 1, //热门推荐默认为视频
			Duration: v.Duration,
			ArcStat:  v.ArcStat,
			ArtStat:  v.ArtStat,
		}
		res = append(res, ra)
		aids = append(aids, v.OID)
	}

	//add tags
	tags, err := s.getTags(c, aids)
	if err != nil {
		log.Error("hotArchives s.getTags err(%v)", err)
		return
	}
	s.setTags(res, tags)

	return
}
