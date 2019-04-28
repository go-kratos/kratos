package service

import (
	"context"
	"github.com/pkg/errors"
	"go-common/app/admin/main/videoup/model/archive"
	"go-common/app/admin/main/videoup/model/manager"
	"go-common/app/admin/main/videoup/model/search"
	accApi "go-common/app/service/main/account/api"
	"go-common/library/database/elastic"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
	"strconv"
	"strings"
	"sync"
)

// SearchVideo search video
func (s *Service) SearchVideo(c context.Context, p *search.VideoParams) (result *search.VideoResultData, err error) {
	var (
		aids, cids,
		vids, mids,
		tids, tagIds,
		xcodes []int64
		fNames           []string
		ps               = 30
		sort             string
		order            string
		isMonitor        bool          //是否查看监控列表
		moniTotal        int           //监控列表的总数量
		moniMap          map[int64]int //监控列表的视频id map。键：vid；值：time（进入监控的时间）
		tags             map[int64]string
		aStates, vStates map[int64]int
		members          map[int64]*accApi.Info
		eReq             *elastic.Request
		wg               sync.WaitGroup
	)
	es := elastic.NewElastic(nil)
	//Page
	if p.Ps != 0 {
		ps = p.Ps
	}
	if p.Pn == 0 {
		p.Pn = 1
	}
	if p.OrderType == "1" {
		eReq = es.NewRequest("archive_video_score")
	} else {
		eReq = es.NewRequest("archive_video")
	}
	eReq.Index("archive_video")

	if p.Keywords != "" {
		eReq.WhereLike([]string{"arc_title", "arc_author"}, []string{p.Keywords}, true, elastic.LikeLevelLow)
	}
	if p.ArcTitle != "" {
		eReq.WhereLike([]string{"arc_title"}, []string{p.ArcTitle}, false, elastic.LikeLevelLow)
	}
	if p.Aids != "" {
		p.Aids = strings.Join(Slice2String(SliceUnique(Slice2Interface(strings.Split(p.Aids, ",")))), ",")
		if aids, err = s.SplitInts(p.Aids); err != nil {
			err = ecode.RequestErr
			return
		}
		if len(aids) > ps {
			ps = len(aids)
		}
		eReq.WhereIn("aid", aids)
	}

	//监控结果列表逻辑
	if p.MonitorList != "" {
		isMonitor = true
		var (
			rid int
		)
		moniP := strings.Split(p.MonitorList, "_")
		if len(moniP) != 3 {
			err = errors.New("监控列表标识字段格式错误")
			return
		}
		if rid, err = strconv.Atoi(moniP[2]); err != nil {
			return
		}
		if moniMap, err = s.MonitorStayOids(c, int64(rid)); err != nil {
			return
		}
		moniTotal = len(moniMap)
		for vid := range moniMap {
			vids = append(vids, vid)
		}
		if len(vids) == 0 {
			result = &search.VideoResultData{
				Result: []*search.Video{},
			}
			return
		}
		if len(vids) > ps {
			ps = len(vids)
		}
		eReq.WhereIn("vid", vids)
	}
	if p.Cids != "" {
		p.Cids = strings.Join(Slice2String(SliceUnique(Slice2Interface(strings.Split(p.Cids, ",")))), ",")
		if cids, err = s.SplitInts(p.Cids); err != nil {
			err = ecode.RequestErr
			return
		}
		if len(cids) > ps {
			ps = len(cids)
		}
		eReq.WhereIn("cid", cids)
	}
	if p.Vids != "" {
		if vids, err = s.SplitInts(p.Vids); err != nil {
			err = ecode.RequestErr
			return
		}
		if len(vids) > ps {
			ps = len(vids)
		}
		eReq.WhereIn("vid", vids)
	}
	if p.ArcMids != "" {
		if mids, err = s.SplitInts(p.ArcMids); err != nil {
			err = ecode.RequestErr
			return
		}
		if len(mids) > ps {
			ps = len(mids)
		}
		eReq.WhereIn("arc_mid", mids)
	}
	if p.Xcode != "" {
		if xcodes, err = s.SplitInts(p.Xcode); err != nil {
			err = ecode.RequestErr
			return
		}
		eReq.WhereIn("xcode_state", xcodes)
	}

	if p.TypeID != "" {
		if tids, err = s.SplitInts(p.TypeID); err != nil {
			err = ecode.RequestErr
			return
		}
		for _, tid := range tids {
			if ids, ok := s.typeCache2[int16(tid)]; ok {
				tids = append(tids, ids...)
			}
		}
		eReq.WhereIn("arc_typeid", tids)
	}
	if p.Filename != "" {
		if fNames = strings.Split(p.Filename, ","); err != nil {
			err = ecode.RequestErr
			return
		}
		eReq.WhereIn("filename", fNames)
	}
	if p.TagID != "" {
		eReq.WhereEq("tag_id", p.TagID)
	}

	if p.Status != "" {
		if p.Status == "-100" {
			eReq.WhereEq("relation_state", "-100")
			p.Status = ""
		} else {
			eReq.WhereEq("relation_state", "0")
			eReq.WhereEq("status", p.Status)
		}
	}
	if p.UserType != "" {
		eReq.WhereEq("user_type", p.UserType)
	}
	if p.DurationFrom != "" && p.DurationTo != "" {
		eReq.WhereRange("duration", p.DurationFrom, p.DurationTo, elastic.RangeScopeLcRc)
	} else if p.DurationFrom != "" && p.DurationTo == "" {
		eReq.WhereRange("duration", p.DurationFrom, "", elastic.RangeScopeLcRc)
	} else if p.DurationFrom == "" && p.DurationTo != "" {
		eReq.WhereRange("duration", "", p.DurationTo, elastic.RangeScopeLcRc)
	}

	//Order by
	if p.Order != "" {
		order = p.Order
	} else if p.Action == "trash" {
		order = "v_mtime"
	} else {
		order = "arc_senddate"
	}
	if p.Sort == 1 {
		sort = "asc"
	} else {
		sort = "desc"
	}
	eReq.Order(order, sort)

	//Page
	eReq.Ps(ps)
	eReq.Pn(p.Pn)

	//Execute
	log.Info("s.SearchVideo(%+v) SearchVideoQuery(%s)", p, eReq.Params())
	if err = eReq.Scan(c, &result); err != nil {
		log.Error("SearchVideoScan() error(%v)", err)
		return
	}
	if result == nil {
		log.Error("s.SearchVideo(%+v) search return nil result", p)
		result = &search.VideoResultData{
			Result: []*search.Video{},
		}
		return
	}
	if result.Result == nil {
		log.Error("s.SearchVideo(%+v) search return nil video", p)
		result.Result = []*search.Video{}
		return
	}
	if len(result.Result) == 0 {
		return
	}
	if isMonitor {
		result.Page.Total = moniTotal
	}
	aids = []int64{}
	vids = []int64{}
	mids = []int64{}
	for _, v := range result.Result {
		if v.TagID != 0 {
			tagIds = append(tagIds, v.TagID)
		}
		aids = append(aids, v.Aid)
		vids = append(vids, v.Vid)
		if v.ArcMid != 0 {
			mids = append(mids, v.ArcMid)
		}
	}
	//获取Tag名称
	wg.Add(1)
	go func(tags *map[int64]string) {
		defer wg.Done()
		if *tags, err = s.arc.TagNameMap(c, tagIds); err != nil {
			log.Error("s.arc.TagNameMap(%v) error(%v)", tagIds, err)
			err = nil
		}
		log.Info("s.arc.TagNameMap(%v) result(%v)", tagIds, tags)
	}(&tags)
	//获取数据库中的稿件状态
	wg.Add(1)
	go func(aStates *map[int64]int) {
		defer wg.Done()
		if *aStates, err = s.arc.ArcStateMap(c, aids); err != nil {
			log.Error("s.arc.ArcStateMap(%v) error(%v)", aids, err)
			err = nil
		}
	}(&aStates)
	//获取数据库中视频的状态
	wg.Add(1)
	go func(vStates *map[int64]int) {
		defer wg.Done()
		if *vStates, err = s.arc.VideoStateMap(c, vids); err != nil {
			log.Error("s.arc.VideoStateMap(%v) error(%v)", vids, err)
			err = nil
		}
	}(&vStates)
	//获取UP主信息
	wg.Add(1)
	go func(members *map[int64]*accApi.Info) {
		defer wg.Done()
		var infosReply *accApi.InfosReply
		if infosReply, err = s.accRPC.Infos3(c, &accApi.MidsReq{Mids: mids}); err != nil {
			log.Error("s.accRPC.Infos3(%v) error(%v)", mids, err)
			err = nil
			return
		}
		*members = infosReply.Infos
	}(&members)

	wg.Wait()
	sInt, _ := strconv.Atoi(p.Status)
	for i := 0; i < len(result.Result); i++ {
		v := result.Result[i]
		result.Result[i].ID = v.Vid
		if vs, ok := vStates[v.Vid]; ok {
			result.Result[i].Status = vs
		}
		//需要将status与查询条件中status不一致的视频剔除
		if p.Status != "" && v.Status != sInt {
			result.Result = append(result.Result[:i], result.Result[i+1:]...)
			i--
			continue
		}
		if tn, ok := tags[v.TagID]; ok {
			result.Result[i].TagName = tn
		}
		if as, ok := aStates[v.Aid]; ok {
			result.Result[i].ArcState = as
		}
		if m, ok := members[v.ArcMid]; ok {
			result.Result[i].ArcAuthor = m.Name
		}
		if v.UserType == nil {
			v.UserType = []int64{}
		}
		if v.UserGroup == nil {
			v.UserGroup = []*manager.UpGroup{}
		}
		for _, tp := range v.UserType {
			if up, ok := s.allUpGroupCache[tp]; ok {
				result.Result[i].UserGroup = append(result.Result[i].UserGroup, up)
			}
		}
	}
	return
}

// SearchCopyright search video copyright
func (s *Service) SearchCopyright(c context.Context, kw string) (result *search.CopyrightResultData, err error) {
	return s.search.SearchCopyright(c, kw)
}

// SearchArchive 稿件搜索列表
func (s *Service) SearchArchive(c *bm.Context, p *search.ArchiveParams) (result *search.ArchiveResultData, err error) {
	var (
		round                                 int
		aids                                  []int64
		froms                                 []int8
		forbid                                string
		tids                                  []int64
		pn                                    = 1
		ps                                    = 30
		isMonitor                             bool          //是否查看监控列表
		moniMap                               map[int64]int //监控列表的视频id map。键：vid；值：time（进入监控的时间）
		nPGCAids, sPGCAids, cPGCAids, misAids []int64
		additMap                              map[int64]*archive.Addit
		tips                                  string
		orders                                = map[string]string{
			"mtime":     "",
			"pubtime":   "",
			"ctime":     "",
			"dm_count":  "",
			"fav_count": "",
		}
		eReq *elastic.Request
	)
	es := elastic.NewElastic(nil)
	if p.OrderType != "" {
		eReq = es.NewRequest("archive_score")
	} else {
		eReq = es.NewRequest("archive")
	}
	eReq.Index("archive")
	if p.Ps != 0 {
		ps = p.Ps
	}
	if ps > 1000 {
		ps = 1000
	}
	//分区 逻辑
	if p.TypeID != "" {
		if tids, err = s.SplitInts(p.TypeID); err != nil {
			return
		}
		for _, tid := range tids {
			if ids, ok := s.typeCache2[int16(tid)]; ok {
				tids = append(tids, ids...)
			}
		}
	}
	//特殊分区逻辑
	if p.SpecialType != "" {
		for i := 0; i < len(tids); i++ { //剔除特殊/普通分区id
			_, ok := s.adtTpsCache[int16(tids[i])]
			if (ok && p.SpecialType == "0") || (!ok && p.SpecialType == "1") {
				tids = append(tids, tids...)
				i--
				continue
			}
		}
	}
	if len(tids) > 0 {
		eReq.WhereIn("typeid", tids)
		if len(tids) > ps {
			ps = len(tids)
		}
	}

	if p.UserType != "" && p.UserType != "0" {
		eReq.WhereEq("user_type", p.UserType)
	}
	//up_from 逻辑
	if p.UpFroms != "" {
		var upFroms []int64
		if upFroms, err = s.SplitInts(p.UpFroms); err != nil {
			return
		}
		eReq.WhereIn("up_from", upFroms)
	} else { //默认去除PGC稿件
		eReq.WhereIn("up_from", []int8{archive.UpFromPGC, archive.UpFromSecretPGC, archive.UpFromCoopera})
		eReq.WhereNot(elastic.NotTypeIn, "up_from")
	}
	//Round 逻辑
	if p.Round != "" {
		if round, err = strconv.Atoi(p.Round); err != nil {
			log.Warn("http.searchArchive() http.PGCListLogic() err(%v)", err)
			return
		}
		if round == 1 {
			eReq.WhereIn("round", []int8{1, archive.RoundReviewSecond})
		} else if round == 2 {
			eReq.WhereIn("round", []int8{2, archive.RoundEnd})
		} else {
			eReq.WhereIn("round", []int{round})
		}
	}
	//Aid 逻辑
	if p.Aids != "" {
		p.Aids = strings.Replace(p.Aids, "\n", ",", -1)
		p.Aids = strings.Join(Slice2String(SliceUnique(Slice2Interface(strings.Split(p.Aids, ",")))), ",")
		if aids, err = s.SplitInts(p.Aids); err != nil {
			return
		}
		eReq.WhereIn("id", aids)
	}
	//Mid 逻辑
	if p.Mids != "" {
		var mids []int64
		p.Mids = strings.Replace(p.Mids, "\n", ",", -1)
		if mids, err = s.SplitInts(p.Mids); err != nil {
			return
		}
		eReq.WhereIn("mid", mids)
		if len(mids) > ps {
			ps = len(mids)
		}
	}
	if p.State != "" {
		states := strings.Split(p.State, ",")
		eReq.WhereIn("state", states)
	}
	if p.Access != "" {
		eReq.WhereEq("access", p.Access)
	}
	if p.Copyright != "" {
		eReq.WhereEq("copyright", p.Copyright)
	}
	if p.IsFirst != "" {
		eReq.WhereEq("is_first", p.IsFirst)
	}

	//PGC 列表逻辑
	if p.PGCList != "" {
		if err = s.PGCListLogic(c, p.PGCList, p.State, froms, eReq); err != nil {
			return
		}
		if c.IsAborted() { //检查鉴权
			return
		}
	}

	//搜索关键字逻辑
	if p.Keywords != "" {
		fields := []string{"title", "content", "tag"}
		if p.KwFields != "" {
			fields = strings.Split(p.KwFields, ",")
			for i := 0; i < len(fields); i++ {
				if fields[i] == "channel" {
					fields[i] = "tid_names"
					break
				}
			}
		}
		eReq.WhereLike(fields, []string{p.Keywords}, true, elastic.LikeLevelLow)
	}
	if p.NoMission == "" { //默认去掉活动稿件
		p.NoMission = "1"
	}
	//禁止项逻辑
	if p.Attr == strconv.Itoa(archive.ForbidAttrChannel) {
		forbid = "channel"
		eReq.WhereEq("channel_group_id", archive.FlowGroupNoChannel)
		eReq.WhereEq("channel_pool", archive.PoolArcForbid)
		eReq.WhereEq("channel_state", archive.FlowOpen)
	} else if p.Attr == strconv.Itoa(archive.ForbidAttrHot) {
		forbid = "hot"
		eReq.WhereEq("hot_group_id", archive.FlowGroupNoHot)
		eReq.WhereEq("hot_pool", archive.PoolArcForbid)
		eReq.WhereEq("hot_state", archive.FlowOpen)
	} else if p.Attr != "" {
		attr := 0
		attr, err = strconv.Atoi(p.Attr)
		if err != nil {
			return
		}
		if (uint(attr-1) == archive.AttrBitIsPorder) && int8(round) != archive.RoundReviewFlow {
			s.auth.Permit("PRIVATE_ORDER_ALL")(c) //鉴权
			if c.IsAborted() {
				return
			}
			p.NoMission = "" //不做活动稿件剔除
		}
		eReq.WhereEq("attribute", attr)
	}

	//回查列表逻辑
	if p.Review != "" {
		if p.ReviewState == "" {
			p.ReviewState = strconv.Itoa(archive.RecheckStateWait)
		}
		switch p.Review {
		case strconv.Itoa(archive.TypeChannelRecheck):
			s.auth.Permit("CHANNEL_REVIEW")(c)
			eReq.WhereEq("recheck_ch_type", p.Review)
			eReq.WhereEq("recheck_ch_staten", p.ReviewState)
		case strconv.Itoa(archive.TypeHotRecheck):
			p.NoMission = ""
			s.auth.Permit("ARC_HOT_REVIEW")(c)
			eReq.WhereEq("recheck_hot_type", p.Review)
			eReq.WhereEq("recheck_hot_staten", p.ReviewState)

		case strconv.Itoa(archive.TypeInspireRecheck):
			p.NoMission = ""
			s.auth.Permit("ARC_INSPIRE_REVIEW")(c)
			eReq.WhereEq("recheck_inspire_type", p.Review)
			eReq.WhereEq("recheck_inspire_staten", p.ReviewState)
		}
		if c.IsAborted() {
			return
		}
		p.IsOrder = 0
	} else if p.Round == strconv.Itoa(int(archive.RoundAuditUGCPayFlow)) { //付费列表需要展示活动稿件
		p.NoMission = ""
	}
	//活动逻辑
	if p.MissionID != "" {
		eReq.WhereEq("mission_id", p.MissionID)
	}
	if p.MissionID != "" || p.NoMission == "0" { //需要活动稿件
		eReq.WhereRange("mission_id", 1, nil, elastic.RangeScopeLcRo)
	} else if p.NoMission == "1" { //不需要活动稿件
		eReq.WhereRange("mission_id", 0, nil, elastic.RangeScopeLoRo)
		eReq.WhereNot(elastic.NotTypeRange, "mission_id")
	}
	//商单逻辑
	if p.IsOrder == 1 { //商单 order_id 大于0
		eReq.WhereRange("order_id", 0, nil, elastic.RangeScopeLoRo)
		if p.OrderId != "" {
			eReq.WhereEq("order_id", p.OrderId)
		}
	} else if forbid == "" && p.Review == "" { //除了禁止列表和回查列表，都要去除商单order_id <= 0
		eReq.WhereRange("order_id", 0, nil, elastic.RangeScopeLoRo)
		eReq.WhereNot(elastic.NotTypeRange, "order_id")
	}

	//监控结果列表逻辑
	if p.MonitorList != "" {
		eReq = es.NewRequest("archive_score") //去掉其它条件
		eReq.Index("archive")
		isMonitor = true
		var (
			rid int
		)
		moniP := strings.Split(p.MonitorList, "_")
		if len(moniP) != 3 {
			err = errors.New("监控列表标识字段格式错误")
			return
		}
		if rid, err = strconv.Atoi(moniP[2]); err != nil {
			return
		}
		if moniMap, err = s.MonitorStayOids(c, int64(rid)); err != nil {
			return
		}
		for aid := range moniMap {
			aids = append(aids, aid)
		}
		if len(aids) == 0 {
			result = &search.ArchiveResultData{
				Result: []*search.Archive{},
			}
			return
		}
		/*if len(aids) > ps {
			ps = len(aids)
		}*/
		eReq.WhereIn("id", aids)
	}
	//分页、排序
	if p.Pn > 0 {
		pn = p.Pn
	}
	if _, ok := orders[p.Order]; !ok {
		p.Order = "ctime"
	}
	if p.Order == "" {
		p.Order = "ctime"
	}
	eReq.OrderScoreFirst(p.ScoreFirst != "0")

	if p.Sort == "" {
		p.Sort = "desc"
	}
	eReq.Order(p.Order, p.Sort)
	eReq.Pn(pn)
	eReq.Ps(ps)
	log.Info("SearchArchiveQuery(%s)", eReq.Params())
	if err = eReq.Scan(c, &result); err != nil {
		log.Error("SearchArchiveQuery() error(%v)", err)
		return
	}
	if isMonitor {
		result.MoniAids = moniMap
	}
	if len(aids) > 0 && !isMonitor {
		//获取稿件addit
		if additMap, err = s.arc.AdditBatch(c, aids); err != nil {
			log.Error("s.arc.AdditBatch(%v) error(%v) additMap(%v)", aids, err, additMap)
			err = nil
		}
		for _, v := range aids {
			if _, ok := additMap[v]; ok {
				switch additMap[v].UpFrom {
				case archive.UpFromPGC:
					nPGCAids = append(nPGCAids, v)
				case archive.UpFromSecretPGC:
					sPGCAids = append(sPGCAids, v)
				case archive.UpFromCoopera:
					cPGCAids = append(cPGCAids, v)
				}
				if additMap[v].MissionID > 0 {
					misAids = append(misAids, v)
				}
			}
		}

		if len(nPGCAids) > 0 {
			tips += "PGC稿件：" + xstr.JoinInts(nPGCAids)
		}
		if len(sPGCAids) > 0 {
			tips += "PGC机密：" + xstr.JoinInts(sPGCAids)
		}
		if len(cPGCAids) > 0 {
			tips += "PGC嵌套：" + xstr.JoinInts(cPGCAids)
		}
		if len(misAids) > 0 {
			tips += "活动稿件：" + xstr.JoinInts(misAids)
		}

		result.Tips = tips
	}
	if err = s.KneadArchiveResult(c, result, p, additMap); err != nil {
		return
	}
	return
}

// PGCListLogic PGC列表查询相关逻辑
func (s *Service) PGCListLogic(c *bm.Context, lnStr string, pState string, froms []int8, eReq *elastic.Request) (err error) {
	var (
		state     int ////临时存储PGC列表中查询的state状态，需要和pState字符串配合使用！
		lni       int
		ln        int8
		pgcConfig map[int8]*search.ArcPGCConfig
	)

	if lni, err = strconv.Atoi(lnStr); err != nil {
		log.Warn("s.PGCListLogic() err(%v)", err)
		return
	}
	ln = int8(lni)
	if state, err = strconv.Atoi(pState); err != nil {
		err = errors.New("PGC列表中的state参数错误")
		return
	}
	if len(froms) != 0 {
		for _, v := range froms {
			if v != archive.UpFromPGC && v != archive.UpFromSecretPGC && v != archive.UpFromCoopera {
				err = errors.New("PGC列表中的from参数错误")
				break
			}
		}
	}
	pgcConfig = make(map[int8]*search.ArcPGCConfig)

	//PGC常规二审列表   round=10 state=-1,-30,-40,-6 upfrom=1
	pgcConfig[1] = &search.ArcPGCConfig{
		UPFrom: []int8{archive.UpFromPGC},
		Rounds: []int8{
			archive.RoundBegin,
			archive.RoundAuditSecond,
		},
		States: []int8{
			archive.RecheckStateWait,
			archive.StateForbidSubmit,
			archive.StateForbidUserDelay,
			archive.StateForbidFixed,
		},
		InState: pState != "" && state < 0,
		Auth:    "PGC_NORMAL_2",
	}
	//PGC常规三审列表   round=20 state=-1,-30,-40,-6 upfrom=1
	pgcConfig[2] = &search.ArcPGCConfig{
		UPFrom: []int8{archive.UpFromPGC},
		Rounds: []int8{archive.RoundAuditThird},
		States: []int8{
			archive.RecheckStateWait,
			archive.StateForbidSubmit,
			archive.StateForbidUserDelay,
			archive.StateForbidFixed,
		},
		InState: pState != "" && state < 0,
		Auth:    "PGC_NORMAL_3",
	}
	//PGC机密待审列表   state=-1,-30,-40,-6 upfrom=5
	pgcConfig[3] = &search.ArcPGCConfig{
		UPFrom: []int8{archive.UpFromSecretPGC},
		Rounds: []int8{},
		States: []int8{
			archive.RecheckStateWait,
			archive.StateForbidSubmit,
			archive.StateForbidUserDelay,
			archive.StateForbidFixed,
		},
		InState: pState != "" && state < 0,
		Auth:    "PGC_SECRET_WAIT",
	}
	//PGC机密回查列表   round=90 state≥0 upfrom=5
	pgcConfig[4] = &search.ArcPGCConfig{
		UPFrom: []int8{archive.UpFromSecretPGC},
		Rounds: []int8{archive.RoundTriggerClick},
		States: []int8{
			archive.StateOpen,
			archive.StateOrange,
		},
		InState: pState != "" && state >= 0,
		Auth:    "PGC_SECRET_RECHECK",
	}

	//PGC全部已过审 round=99 state≥0 upfrom=1,5,6
	pgcConfig[5] = &search.ArcPGCConfig{
		UPFrom: []int8{archive.UpFromPGC, archive.UpFromSecretPGC, archive.UpFromCoopera},
		Rounds: []int8{archive.RoundEnd},
		States: []int8{
			archive.StateOpen,
			archive.StateOrange,
		},
		InState: pState != "" && state >= 0,
		Auth:    "PGC_OPEN",
	}

	//全部打回列表    state=-2,-3,-4,-7,-11,-16,-100 up_from=1,5,6
	pgcConfig[6] = &search.ArcPGCConfig{
		UPFrom: []int8{archive.UpFromPGC, archive.UpFromSecretPGC, archive.UpFromCoopera},
		Rounds: []int8{},
		States: []int8{
			archive.StateForbidRecycle,
			archive.StateForbidPolice,
			archive.StateForbidLock,
			archive.StateForbidLater,
			archive.StateForbidFixing,
			archive.StateForbidXcodeFail,
			archive.StateForbidUpDelete,
		},
		InState: pState != "" && state < 0,
		Auth:    "PGC_RECICLE",
	}

	//合作方嵌套列表   state=-1,-30,-6,-40  up_from=6
	pgcConfig[7] = &search.ArcPGCConfig{
		UPFrom: []int8{archive.UpFromCoopera},
		Rounds: []int8{},
		States: []int8{
			archive.StateForbidWait,
			archive.StateForbidSubmit,
			archive.StateForbidFixed,
			archive.StateForbidUserDelay,
		},
		InState: pState != "" && state < 0,
		Auth:    "PGC_PARTNER",
	}
	if _, ok := pgcConfig[ln]; !ok || pgcConfig[ln] == nil {
		err = errors.New("PGC列表不存在")
		return
	}
	if !pgcConfig[ln].InState { //如果前端传了合法的state，则加上state条件
		pgcConfig[ln].States = []int8{int8(state)}
	}
	if len(froms) != 0 { //如果前端传了up_from，则加上up_from条件
		pgcConfig[ln].UPFrom = froms
	}
	if len(pgcConfig[ln].UPFrom) == 0 {
		err = errors.New("PGC列表中的UPFrom不能设置为空")
		return
	}
	s.auth.Permit(pgcConfig[ln].Auth)(c) //鉴权
	eReq.WhereIn("up_from", pgcConfig[ln].UPFrom)
	eReq.WhereIn("state", pgcConfig[ln].States)
	eReq.WhereIn("round", pgcConfig[ln].Rounds)
	return
}

// KneadArchiveResult 拼接稿件数据
func (s *Service) KneadArchiveResult(c *bm.Context, result *search.ArchiveResultData, p *search.ArchiveParams, additMap map[int64]*archive.Addit) (err error) {
	var (
		mids, aids []int64
		ups        map[int64]*accApi.Card
		rStates    map[int64]int8
		chNames    map[int64][]string
		archives   = result.Result
		dbArchives map[int64]*archive.Archive
		wg         sync.WaitGroup
	)
	for _, v := range archives {
		if v.Mid != 0 {
			mids = append(mids, v.Mid)
		}
		aids = append(aids, v.ID)
	}

	//获取频道信息
	wg.Add(1)
	go func(chNames *map[int64][]string) {
		defer wg.Done()
		*chNames = s.ChannelNamesByAids(c, aids)
	}(&chNames)

	wg.Add(1)
	go func(dbArchives *map[int64]*archive.Archive) {
		defer wg.Done()
		if *dbArchives, err = s.arc.Archives(c, aids); err != nil || dbArchives == nil {
			log.Error("s.arc.Archives(%v) error(%v) archives(%v)", mids, err, dbArchives)
			err = nil
		}
	}(&dbArchives)

	//获取UP主信息
	wg.Add(1)
	go func(ups *map[int64]*accApi.Card) {
		defer wg.Done()
		if *ups, err = s.upCards(c, mids); err != nil || ups == nil {
			log.Error("s.upCards(%v) error(%v) ups(%v)", mids, err, ups)
			err = nil
			*ups = make(map[int64]*accApi.Card)
		}
	}(&ups)
	//获取回查稿件实时状态
	if p.Review != "" && p.ReviewState != "" {
		wg.Add(1)
		go func(rStates *map[int64]int8) {
			defer wg.Done()
			tp, _ := strconv.Atoi(p.Review)
			if *rStates, err = s.arc.RecheckStateMap(c, tp, aids); err != nil || rStates == nil {
				log.Error("s.arc.RecheckStateMap(%v,%v) error(%v) rStates(%v)", tp, aids, err, rStates)
				err = nil
			}
		}(&rStates)
	}

	wg.Wait()
	for i := 0; i < len(archives); i++ {
		v := archives[i]
		if a, ok := dbArchives[v.ID]; ok && a.Cover != "" {
			archives[i].Cover = coverURL(a.Cover)
		}
		if m, ok := ups[v.Mid]; ok {
			archives[i].Official = m.Official
			archives[i].Author = m.Name
		}
		if p.Review != "" && p.ReviewState != "" {
			if state, ok := rStates[v.ID]; ok { //将与查询条件中回查状态不一致的稿件剔除
				rState, err := strconv.Atoi(p.ReviewState)
				if err != nil {
					log.Error("s.KneadArchive() error(%v) reviewState(%v)", err, p.ReviewState)
					err = nil
					continue
				}
				if state != int8(rState) {
					archives = append(archives[:i], archives[i+1:]...)
					i--
					continue
				}
			}
		}
		if _, ok := additMap[v.ID]; ok {
			archives[i].UpFrom = additMap[v.ID].UpFrom
			if additMap[v.ID].MissionID > 0 {
				archives[i].MissionID = additMap[v.ID].MissionID
			}
		}
		if names, ok := chNames[v.ID]; ok {
			archives[i].TagNames = names
		}
		if archives[i].TagNames == nil {
			archives[i].TagNames = []string{}
		}
		if archives[i].UserType == nil {
			archives[i].UserType = []int64{}
		}
		if v.UserGroup == nil {
			v.UserGroup = []*manager.UpGroup2{}
		}
		if v.Attribute == nil {
			v.Attribute = []int{}
		}
		v.Attrs = v.Attribute
		for _, tp := range v.UserType {
			if up, ok := s.allUpGroupCache[tp]; ok {
				//因为前端希望稿件列表返回的up_group与任务质检结构一致，所以在这做一次转换
				gp := &manager.UpGroup2{
					GroupID:   up.ID,
					GroupName: up.Name,
					GroupTag:  up.ShortTag,
				}
				result.Result[i].UserGroup = append(result.Result[i].UserGroup, gp)
			}
		}
	}
	result.Result = archives
	return
}
