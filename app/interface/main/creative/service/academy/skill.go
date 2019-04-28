package academy

import (
	"context"
	"sort"
	"time"

	"go-common/app/interface/main/creative/model/academy"
	"go-common/app/interface/main/creative/model/newcomer"
	"go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"
)

// Occupations get occ.
func (s *Service) Occupations(c context.Context) (res []*academy.Occupation, err error) {
	res = s.OccCache

	res = make([]*academy.Occupation, 0, len(res))
	for _, v := range s.OccCache {
		if v.ID == 7 || v.ID == s.NewbCourseID { //7-废弃 NewbCourseID-新人专区
			continue
		}
		res = append(res, v)
	}
	return
}

// NewbCourse get new up course.
func (s *Service) NewbCourse(c context.Context) (res []*academy.NewbCourseList, err error) {
	var (
		skids []int64
		pids  = []int64{}
		sids  = []int64{}
		tcs   *academy.ArcList
		skMap map[int64]string
		aids  []int64
	)

	skMap = make(map[int64]string)
	for _, v := range s.SkillCache {
		if v.OID == s.NewbCourseID {
			skids = append(skids, v.ID)
			skMap[v.ID] = v.Name
		}
	}

	tcs, err = s.ThemeCourse(c, pids, skids, sids, 1, 20, true)
	if err != nil {
		log.Error("NewbCourse s.ThemeCourse pids(%+v)|skids(%+v)|sids(%+v)|error(%v)", pids, skids, sids, err)
		return
	}
	if len(tcs.Items) == 0 {
		return
	}

	// add tags
	for _, v := range tcs.Items {
		aids = append(aids, v.AID)
	}
	tags, err := s.getTags(c, aids)
	if err != nil {
		log.Error("NewbCourse s.getTags err(%v)", err)
		return
	}
	s.setTags(tcs.Items, tags)

	newbCourseMap := make(map[int64][]*academy.ArcMeta)
	for _, v := range tcs.Items {
		if v == nil || v.Skill == nil {
			continue
		}
		newbCourseMap[v.Skill.SkID] = append(newbCourseMap[v.Skill.SkID], v)
	}

	res = make([]*academy.NewbCourseList, 0)
	for _, id := range skids {
		l := &academy.NewbCourseList{}
		if v, ok := newbCourseMap[id]; ok {
			l.Items = v
		}

		if sname, ok := skMap[id]; ok {
			l.Title = sname
			l.TID = id
		}
		res = append(res, l)
	}
	return
}

// ThemeCourse get theme course.
func (s *Service) ThemeCourse(c context.Context, pids, skids, sids []int64, pn, ps int, isnew bool) (res *academy.ArcList, err error) {
	var (
		skas    []*academy.SkillArc
		aids    []int64
		pidMap  map[int64]int64
		skidMap map[int64]int64
		sidMap  map[int64]int64
		total   int
	)

	res = &academy.ArcList{
		Items: []*academy.ArcMeta{},
		Page: &academy.ArchivePage{
			Pn: pn,
			Ps: ps,
		},
	}

	if !isnew && len(pids) == 0 { //如果不是新人课程，并且不传职业课程id则默认全部
		for _, v := range s.OccCache {
			pids = append(pids, v.ID)
		}
	}

	if !isnew { //不是新人课程，则去掉新人课程
		for i := 0; i < len(pids); i++ {
			if pids[i] == s.NewbCourseID {
				pids = append(pids[:i], pids[i+1:]...)
				i--
			}
		}
	}

	if skas, err = s.aca.SkillArcs(c, pids, skids, sids, (pn-1)*ps, ps); err != nil {
		log.Error("s.aca.SkillArcs pid(%+v)|skid(%+v)|sid(%+v)|error(%v)", pids, skids, sids, err)
		return
	}
	if len(skas) == 0 {
		log.Error("s.aca.SkillArcs has no data")
		return
	}
	pidMap = make(map[int64]int64)
	skidMap = make(map[int64]int64)
	sidMap = make(map[int64]int64)
	for _, v := range skas {
		aids = append(aids, v.AID)
		pidMap[v.AID] = v.PID
		skidMap[v.AID] = v.SkID
		sidMap[v.AID] = v.SID
	}

	if total, err = s.aca.SkillArcCount(c, pids, skids, sids); err != nil {
		log.Error("s.aca.SkillArcCount pids(%+v)|skids(%+v)|sids(%+v)|error(%v)", pids, skids, sids, err)
		return
	}
	res.Page.Total = total

	var (
		g, _    = errgroup.WithContext(c)
		arcInfo map[int64]*api.Arc
		as      map[int64]*api.Stat
	)

	g.Go(func() error {
		arcInfo, err = s.arc.Archives(c, aids, "")
		if err != nil {
			log.Error("s.arc.Archives aids(%+v)|error(%v)", aids, err)
		}
		return err
	})
	g.Go(func() error {
		as, err = s.arc.Stats(c, aids, "")
		if err != nil {
			log.Error("s.arc.Stats aids(%+v)|error(%v)", aids, err)
		}
		return err
	})
	if g.Wait() != nil {
		log.Error("s.aca.ThemeCourse g.Wait() error(%v)", err)
		return
	}

	items := make([]*academy.ArcMeta, 0, len(aids))
	for _, aid := range aids {
		v, ok := arcInfo[aid]
		if !ok || v == nil {
			log.Error("ThemeCourse bind ArcInfo aid(%d) error", aid)
			return
		}
		a := &academy.ArcMeta{
			AID:      aid,
			Cover:    v.Pic,
			Title:    v.Title,
			Type:     v.TypeName,
			MID:      v.Author.Mid,
			Duration: v.Duration,
			Skill:    &academy.SkillArc{},
			Business: 1, //技能树只有视频
		}
		if st, ok := as[aid]; ok {
			a.ArcStat = st
		} else {
			a.ArcStat = &api.Stat{}
		}
		if pid, ok := pidMap[aid]; ok {
			a.Skill.PID = pid
		}
		if skid, ok := skidMap[aid]; ok {
			a.Skill.SkID = skid
		}
		if sid, ok := sidMap[aid]; ok {
			a.Skill.SID = sid
		}
		items = append(items, a)
	}
	res.Items = items
	return
}

// ViewPlay view play archive by aid & mid & business.
func (s *Service) ViewPlay(c context.Context, mid, aid int64, bus int8) (play *academy.Play, err error) {
	p := &academy.Play{
		MID:      mid,
		AID:      aid,
		Business: bus,
	}
	play, err = s.aca.Play(c, p)
	if err != nil {
		log.Error("ViewPlay s.aca.Play error(%v)", err)
	}
	return
}

// PlayAdd add play archive by aid & mid.
func (s *Service) PlayAdd(c context.Context, mid, aid int64, bus, watch int8) (id int64, err error) {
	py := &academy.Play{
		MID:      mid,
		AID:      aid,
		Business: bus,
		Watch:    watch,
		CTime:    xtime.Time(time.Now().Unix()),
		MTime:    xtime.Time(time.Now().Unix()),
	}
	if id, err = s.aca.PlayAdd(c, py); err != nil {
		log.Error("s.aca.PlayAdd error(%v)", err)
		return
	}
	s.p.TaskPub(mid, newcomer.MsgForAcademyFavVideo, newcomer.MsgFinishedCount)
	return
}

// PlayDel del play archive by aid & mid & business.
func (s *Service) PlayDel(c context.Context, mid, aid int64, bus int8) (id int64, err error) {
	p := &academy.Play{
		MID:      mid,
		AID:      aid,
		Business: bus,
	}

	play, err := s.aca.Play(c, p)
	if err != nil {
		log.Error("PlayDel s.aca.Play error(%v)", err)
		return
	}
	if play == nil {
		err = ecode.NothingFound
		return
	}

	if id, err = s.aca.PlayDel(c, p); err != nil {
		log.Error("s.aca.PlayDel error(%v)", err)
	}
	return
}

// PlayList get play list.
func (s *Service) PlayList(c context.Context, mid int64, pn, ps int) (res *academy.ArcList, err error) {
	var (
		pls        []*academy.Play
		total      int
		aids, cids []int64
		playMap    map[int64]*academy.Play
	)

	res = &academy.ArcList{
		Items: []*academy.ArcMeta{},
		Page: &academy.ArchivePage{
			Pn: pn,
			Ps: ps,
		},
	}

	if pls, err = s.aca.Plays(c, mid, (pn-1)*ps, ps); err != nil {
		log.Error("s.aca.Plays mid(%d)|error(%v)", mid, err)
		return
	}
	if len(pls) == 0 {
		log.Error("s.aca.Plays has no mid(%d)", mid)
		return
	}
	playMap = make(map[int64]*academy.Play)

	for _, v := range pls {
		if v.Business == 1 {
			aids = append(aids, v.AID)
		} else if v.Business == 2 {
			cids = append(cids, v.AID)
		}
		playMap[v.AID] = v
	}

	if total, err = s.aca.PlayCount(c, mid); err != nil {
		log.Error("s.aca.PlayCount error(%v)", err)
		return
	}
	res.Page.Total = total

	var (
		arcs []*academy.ArcMeta
		arts []*academy.ArcMeta
		g, _ = errgroup.WithContext(c)
	)

	g.Go(func() error {
		arcs, err = s.getArcInfo(c, aids, playMap)
		return err
	})
	g.Go(func() error {
		arts, err = s.getArtInfo(c, cids, playMap)
		return err
	})
	if g.Wait() != nil {
		log.Error("s.PlayList g.Wait() error(%v)", err)
		return
	}

	tItems := make([]*academy.ArcMeta, 0, len(arcs)+len(arts))
	tItems = append(tItems, arcs...)
	tItems = append(tItems, arts...)

	sort.Slice(tItems, func(i, j int) bool { //按播放时间倒序
		return tItems[i].PlayTime > tItems[j].PlayTime
	})

	unReadItems := make([]*academy.ArcMeta, 0)
	readItems := make([]*academy.ArcMeta, 0)
	for _, v := range tItems {
		if v.Watch == 1 {
			unReadItems = append(unReadItems, v)
		} else if v.Watch == 2 {
			readItems = append(readItems, v)
		}
	}
	res.Items = append(res.Items, unReadItems...) //未观看最先展示
	res.Items = append(res.Items, readItems...)
	return
}

func (s *Service) getArcInfo(c context.Context, aids []int64, playMap map[int64]*academy.Play) (items []*academy.ArcMeta, err error) {
	arcs, err := s.arc.Archives(c, aids, "")
	if err != nil {
		log.Error("s.arc.Archives aids(%+v)|error(%v)", aids, err)
		return
	}

	items = make([]*academy.ArcMeta, 0, len(aids))
	for _, aid := range aids {
		v, ok := arcs[aid]
		if !ok || v == nil {
			log.Error("PlayList bind ArcInfo aid(%d) error", aid)
			return
		}
		a := &academy.ArcMeta{
			AID:      aid,
			Cover:    v.Pic,
			Title:    v.Title,
			Type:     v.TypeName,
			MID:      v.Author.Mid,
			Duration: v.Duration,
		}
		if p, ok := playMap[aid]; ok {
			a.PlayTime = p.MTime
			a.Watch = p.Watch
			a.Business = p.Business
		}
		items = append(items, a)
	}

	return
}

func (s *Service) getArtInfo(c context.Context, cids []int64, playMap map[int64]*academy.Play) (items []*academy.ArcMeta, err error) {
	arts, err := s.art.ArticleMetas(c, cids, "")
	if err != nil {
		log.Error("s.arc.ArticleMetas cids(%+v) error(%v)", cids, err)
		return
	}

	items = make([]*academy.ArcMeta, 0, len(cids))
	for _, cid := range cids {
		v, ok := arts[cid]
		if !ok || v == nil {
			log.Error("PlayList bind ArtInfo cid(%d) error", cid)
			return
		}
		a := &academy.ArcMeta{
			AID:   cid,
			Title: v.Title,
			MID:   v.Author.Mid,
		}
		if v.Category != nil {
			a.Type = v.Category.Name
		}
		if len(v.ImageURLs) > 0 {
			a.Cover = v.ImageURLs[0]
		}

		if p, ok := playMap[cid]; ok {
			a.PlayTime = p.MTime
			a.Watch = p.Watch
			a.Business = p.Business
		}
		items = append(items, a)
	}

	return
}

// ProfessionSkill get theme course.
func (s *Service) ProfessionSkill(c context.Context, pids, skids, sids []int64, pn, ps int, isnew bool) (res []*academy.NewbCourseList, err error) {
	var (
		tcs   *academy.ArcList
		skMap map[int64]string
		aids  []int64
	)
	// 取消分页，默认最大获取100条
	tcs, err = s.ThemeCourse(c, pids, skids, sids, 1, 100, false)
	if err != nil {
		log.Error("ProfessionSkill s.ThemeCourse pids(%+v)|skids(%+v)|sids(%+v)|error(%v)", pids, skids, sids, err)
		return
	}
	if len(tcs.Items) == 0 {
		return
	}

	skMap = make(map[int64]string)
	for _, v := range s.SkillCache {
		skids = append(skids, v.ID)
		skMap[v.ID] = v.Name
	}

	// add tags
	for _, v := range tcs.Items {
		if v == nil {
			continue
		}
		aids = append(aids, v.AID)
	}
	tags, err := s.getTags(c, aids)
	if err != nil {
		log.Error("ProfessionSkill s.getTags err(%v)", err)
		return
	}
	s.setTags(tcs.Items, tags)

	newbCourseMap := make(map[int64][]*academy.ArcMeta)
	for _, v := range tcs.Items {
		if v == nil || v.Skill == nil {
			continue
		}
		newbCourseMap[v.Skill.SkID] = append(newbCourseMap[v.Skill.SkID], v)
	}

	res = make([]*academy.NewbCourseList, 0)
	for _, id := range skids {
		l := &academy.NewbCourseList{}
		if v, ok := newbCourseMap[id]; ok {
			l.Items = v
		}

		if sname, ok := skMap[id]; ok {
			l.Title = sname
			l.TID = id
		}

		if len(l.Items) == 0 {
			continue
		}
		res = append(res, l)
	}
	return
}

// Keywords get keywords.
func (s *Service) Keywords(c context.Context) (res []interface{}) {
	res = s.KWsCache
	return
}
