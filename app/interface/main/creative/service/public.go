package service

import (
	"context"
	"fmt"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/account"
	"go-common/app/interface/main/creative/dao/activity"
	"go-common/app/interface/main/creative/dao/archive"
	"go-common/app/interface/main/creative/dao/article"
	"go-common/app/interface/main/creative/dao/creative"
	"go-common/app/interface/main/creative/dao/pay"
	"go-common/app/interface/main/creative/dao/subtitle"
	"net/url"
	"os"

	"go-common/app/interface/main/creative/dao/tag"
	"go-common/app/interface/main/creative/dao/up"
	actmdl "go-common/app/interface/main/creative/model/activity"
	arcinter "go-common/app/interface/main/creative/model/archive"
	arcmdl "go-common/app/interface/main/creative/model/archive"
	"go-common/app/interface/main/creative/model/music"
	mMdl "go-common/app/interface/main/creative/model/music"
	"go-common/app/interface/main/creative/model/newcomer"
	tagMdl "go-common/app/interface/main/creative/model/tag"
	accMdl "go-common/app/service/main/account/model"
	mdlarc "go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"
	"go-common/library/xstr"
	"hash/crc32"
	"math"
	"strconv"
	"strings"
	"time"
)

//Public struct
type Public struct {
	c        *conf.Config
	creative *creative.Dao
	sub      *subtitle.Dao
	acc      *account.Dao
	act      *activity.Dao
	arc      *archive.Dao
	up       *up.Dao
	pay      *pay.Dao
	tag      *tag.Dao
	// type cache
	TypesCache       map[string][]*arcmdl.Type
	TopTypesCache    []*arcmdl.Type
	TypeMapCache     map[int16]*arcmdl.Type
	CTypesCache      map[string][]*arcmdl.Type
	AllMusics        map[int64]*mMdl.Music
	DescFmtsCache    map[int64]map[int8]map[int8]*arcmdl.DescFormat
	DescFmtsArrCache []*arcmdl.DescFormat
	// cache
	ActVideoAllCache []*actmdl.Activity
	TopActCache      []*actmdl.Activity
	ActMapCache      map[int64]*actmdl.Activity
	StaffTitlesCache []*tagMdl.StaffTitle
	//task
	taskPub              *databus.Databus
	AppWhiteMidsByGroups map[int64]map[int64]int64
}

//RPCDaos struct
type RPCDaos struct {
	Arc *archive.Dao
	Acc *account.Dao
	Art *article.Dao
	Up  *up.Dao
	Sub *subtitle.Dao
}

//New get service
func New(c *conf.Config, rpcdaos *RPCDaos) *Public {
	p := &Public{
		c:                    c,
		creative:             creative.New(c),
		arc:                  archive.New(c),
		sub:                  subtitle.New(c),
		act:                  activity.New(c),
		pay:                  pay.New(c),
		tag:                  tag.New(c),
		acc:                  rpcdaos.Acc,
		up:                   rpcdaos.Up,
		taskPub:              databus.New(c.TaskPub),
		AllMusics:            make(map[int64]*mMdl.Music),
		ActVideoAllCache:     make([]*actmdl.Activity, 0),
		TopActCache:          make([]*actmdl.Activity, 0),
		StaffTitlesCache:     make([]*tagMdl.StaffTitle, 0),
		ActMapCache:          make(map[int64]*actmdl.Activity),
		AppWhiteMidsByGroups: make(map[int64]map[int64]int64),
	}
	p.loadTypes()
	p.loadDescFormat()
	p.loadMusicTable()
	p.loadActivities()
	p.loadPortalGroups()
	p.loadStaffTitles()
	go p.loadproc()
	go p.tableproc()
	return p
}

//loadPortalGroups fn
func (p *Public) loadPortalGroups() {
	var (
		tmpGroupMaps    = make(map[int64]map[int64]int64)
		specialGroupIDs map[int64]int8
	)
	if os.Getenv("DEPLOY_ENV") == "uat" {
		specialGroupIDs = map[int64]int8{
			29: 1,
			12: 1,
		}
	} else {
		specialGroupIDs = map[int64]int8{
			22: 1, // 移动端新手任务白名单
			23: 1, // OPG用户组（内部人员名单
		}
	}
	c := context.TODO()
	gps := make([]int64, 0)
	for gpKey := range specialGroupIDs {
		gps = append(gps, gpKey)
	}
	type ChData struct {
		gp   int64
		gmap map[int64]int64
	}
	rechan := make(chan ChData, len(gps))
	g, ctx := errgroup.WithContext(c)
	for _, gpID := range gps {
		var gid = gpID
		g.Go(func() error {
			ret, e := p.up.UpSpecial(ctx, gid)
			if e != nil {
				log.Warn("p.up.UpSpecial gid (%d)", gid)
				return nil
			}
			if len(ret) > 0 {
				log.Warn("len of ret gid (%d)|(%d)", gid, len(ret))
			}
			rechan <- ChData{gid, ret}
			return nil
		})
	}
	g.Wait()
	close(rechan)
	for c := range rechan {
		tmpGroupMaps[c.gp] = c.gmap
	}
	p.AppWhiteMidsByGroups = tmpGroupMaps
}

func (p *Public) tableproc() {
	for {
		time.Sleep(time.Duration(10 * time.Second))
		p.loadMusicTable()
	}
}

func (p *Public) loadMusicTable() {
	var (
		err      error
		musicMap map[int64]*mMdl.Music
	)
	c := context.TODO()
	if musicMap, err = p.arc.AllMusics(c); err != nil {
		log.Error("p.music.MCategorys err(%+v)", err)
		return
	}
	if musicMap != nil {
		p.AllMusics = musicMap
	}
	log.Info("loadMusicTable (%d)", len(p.AllMusics))
}

// NewRPCDaos get all
func NewRPCDaos(c *conf.Config) *RPCDaos {
	rds := &RPCDaos{
		Arc: archive.New(c),
		Acc: account.New(c),
		Art: article.New(c),
		Up:  up.New(c),
		Sub: subtitle.New(c),
	}
	return rds
}

// loadproc
func (p *Public) loadproc() {
	for {
		time.Sleep(5 * time.Minute)
		p.loadTypes()
		p.loadDescFormat()
		p.loadActivities()
		p.loadPortalGroups()
		p.loadStaffTitles()
	}
}

// loadActivities fn
func (p *Public) loadActivities() {
	p.ActVideoAllCache = make([]*actmdl.Activity, 0)
	videoallActs, err := p.act.Activities(context.TODO())
	if err != nil {
		return
	}
	for _, act := range videoallActs {
		if len(act.Tags) == 0 {
			act.Tags = act.Name
		} else {
			act.Tags = strings.Split(act.Tags, ",")[0]
		}
		v := &actmdl.Activity{
			ID:       act.ID,
			Name:     act.Name,
			Tags:     act.Tags,
			ActURL:   act.ActURL,
			Protocol: act.Protocol,
			Type:     act.Type,
			Hot:      act.Hot,
			STime:    act.STime,
		}
		p.ActVideoAllCache = append(p.ActVideoAllCache, v)
		p.ActMapCache[act.ID] = v
	}
	topLen := 4
	multiplier := p.c.Coefficient.ActHeat
	if len(p.ActVideoAllCache) <= topLen {
		p.TopActCache = p.ActVideoAllCache
	} else {
		p.TopActCache = p.ActVideoAllCache[:topLen]
	}
	for _, topAct := range p.TopActCache {
		stime, _ := time.Parse("2006-01-02 15:04:05", topAct.STime)
		stimeAfter3Day := stime.AddDate(0, 0, 3).Unix()
		if time.Now().Unix() < stimeAfter3Day {
			topAct.New = 1
		}
		likeCnt, _ := p.act.Likes(context.Background(), topAct.ID)
		if likeCnt > 0 {
			topAct.Comment = fmt.Sprintf("%d人参与", int(math.Ceil(float64(likeCnt)*multiplier)))
		}
	}
}

//load types
func (p *Public) loadTypes() {
	tops, langs, typeMap, err := p.creative.Types(context.TODO())
	if err != nil {
		log.Error("p.creative.Types error(%v)", err)
		return
	}
	arcmdl.SortRulesForTopTypes(tops, arcmdl.WebType)
	p.TopTypesCache = tops
	for _, vals := range langs {
		arcmdl.SortRulesForTopTypes(vals, arcmdl.WebType)
	}
	p.TypesCache = langs
	p.CTypesCache = genCTypesCache(langs)
	p.TypeMapCache = typeMap
}

// 自动过滤不需要的二级分区，如果二级分区全部删除了，自动会删除对应的一级分区
func genCTypesCache(langs map[string][]*arcmdl.Type) (CTypesCache map[string][]*arcmdl.Type) {
	CTypesCache = make(map[string][]*arcmdl.Type)
	for lang, topTypes := range langs {
		CTypesCache[lang] = make([]*arcmdl.Type, 0)
		for _, topType := range topTypes {
			nt := &arcmdl.Type{
				ID:        topType.ID,
				Lang:      topType.Lang,
				Parent:    topType.Parent,
				Name:      topType.Name,
				Desc:      topType.Desc,
				Descapp:   topType.Descapp,
				Count:     topType.Count,
				Original:  topType.Original,
				IntroCopy: topType.IntroCopy,
				Notice:    topType.Notice,
				CopyRight: topType.CopyRight,
				Show:      topType.Show,
				Rank:      topType.Rank,
				Children:  []*arcmdl.Type{},
			}
			if arcmdl.ForbidTopTypesForAppAdd(topType.ID) {
				nt.Show = false
			}
			for _, child := range topType.Children {
				if arcmdl.ForbidSubTypesForAppAdd(child.ID) {
					continue
				}
				nt.Children = append(nt.Children, child)
			}
			if len(nt.Children) > 0 {
				CTypesCache[lang] = append(CTypesCache[lang], nt)
			}
		}
		arcmdl.SortRulesForTopTypes(CTypesCache[lang], arcmdl.AppType)
	}
	return
}

// CoverURL convert cover url to full url.
func CoverURL(uri string) (cover string) {
	if uri == "" {
		//cover = "http://static.hdslb.com/images/transparent.gif"
		return
	}
	cover = uri
	if strings.Index(uri, "http://") == 0 {
		return
	}
	if len(uri) >= 10 && uri[:10] == "/templets/" {
		return
	}
	if strings.HasPrefix(uri, "group1") {
		cover = "http://i0.hdslb.com/" + uri
		return
	}
	if pos := strings.Index(uri, "/uploads/"); pos != -1 && (pos == 0 || pos == 3) {
		cover = uri[pos+8:]
	}
	cover = strings.Replace(cover, "{IMG}", "", -1)
	cover = "http://i" + strconv.FormatInt(int64(crc32.ChecksumIEEE([]byte(cover)))%3, 10) + ".hdslb.com" + cover
	return
}

//BatchArchives batch get archive info.
func (p *Public) BatchArchives(c context.Context, mid int64, aids []int64, ip string) (avm map[int64]*arcmdl.ArcVideo, err error) {
	avm, err = p.arc.Views(c, mid, aids, ip)
	if err != nil {
		log.Error("p.arc.Views aids (%v), ip(%s) err(%v)", aids, ip, err)
	}
	return
}

func (p *Public) loadDescFormat() {
	fmts, err := p.arc.DescFormat(context.TODO())
	if err != nil {
		return
	}
	fmtsArr := make([]*arcmdl.DescFormat, 0)
	tp := make(map[int64]map[int8]map[int8]*arcmdl.DescFormat)
	for _, d := range fmts {
		fmtsArr = append(fmtsArr, d)
		if _, okTp := tp[d.TypeID]; !okTp {
			tp[d.TypeID] = make(map[int8]map[int8]*arcmdl.DescFormat)
		}
		if _, okCp := tp[d.TypeID][d.Copyright]; !okCp {
			tp[d.TypeID][d.Copyright] = make(map[int8]*arcmdl.DescFormat)
		}
		if _, okCp := tp[d.TypeID][d.Copyright][d.Lang]; !okCp {
			tp[d.TypeID][d.Copyright][d.Lang] = &arcmdl.DescFormat{}
		}
		tp[d.TypeID][d.Copyright][d.Lang] = d
	}
	p.DescFmtsCache = tp
	p.DescFmtsArrCache = fmtsArr
}

//TaskPub fn pub task finished msg.
func (p *Public) TaskPub(mid int64, from, count int) (err error) {
	msg := &newcomer.TaskMsg{
		MID:       mid,
		From:      from,
		Count:     int64(count),
		TimeStamp: time.Now().Unix(),
	}
	log.Info("task Pub mid(%d) msg(%+v)", mid, msg)
	if err = p.taskPub.Send(context.TODO(), strconv.FormatInt(mid, 10), msg); err != nil {
		log.Error("s.taskPub.Send mid(%d) error(%v)", mid, err)
		return
	}
	return
}

// StaffList fn
func (p *Public) StaffList(c context.Context, aid int64, cache bool) (res []*arcinter.Staff, err error) {
	if cache {
		return p.arc.StaffData(c, aid)
	}
	if res, err = p.arc.RawStaffData(c, aid); err != nil {
		log.Error("s.StaffList(%d) error(%v)", aid, err)
		return
	}
	return
}

//BgmBindList fn
func (p *Public) BgmBindList(c context.Context, aid, cid, mType int64, cache bool) (resOk []*arcinter.ViewBGM, err error) {
	var (
		data         *creative.BgmData
		newIDS, sids []int64
		ret          map[int64]string
		musics       map[int64]*music.Music
		res          []*arcinter.ViewBGM
	)
	//无更新逻辑 注意空缓存
	if data, err = p.creative.BgmData(c, aid, cid, mType, cache); err != nil || data == nil {
		log.Error("s.GetMaterialData(%d,%d,%d) error(%v)", aid, cid, mType, err)
		return
	}
	if sids, err = xstr.SplitInts(data.Data); err != nil {
		log.Error("s.BgmBindList(%d,%d,%d) error(%v)", aid, cid, mType, err)
		return
	}
	for _, sid := range sids {
		if sid > 0 {
			newIDS = append(newIDS, sid)
		}
	}
	if len(newIDS) < 1 {
		return
	}
	//all localcache
	musics = p.AllMusics
	if musics == nil {
		return
	}
	res = make([]*arcinter.ViewBGM, 0)
	resOk = make([]*arcinter.ViewBGM, 0)
	var mids []int64
	for _, sid := range newIDS {
		if _, ok := musics[sid]; !ok {
			continue
		}
		musicData := musics[sid]
		newOne := &arcinter.ViewBGM{}
		newOne.SID = sid
		newOne.MID = musicData.UpMID
		newOne.Title = musicData.Name
		newOne.Author = musicData.Musicians
		if musicData.State == 0 {
			params := url.Values{}
			params.Set("bgm_id", strconv.FormatInt(sid, 10))
			params.Set("from_aid", strconv.FormatInt(aid, 10))
			params.Set("from_cid", strconv.FormatInt(cid, 10))
			params.Set("from_source", "player_page")
			newOne.JumpURL = p.c.H5Page.Cooperate + "?" + params.Encode()
		}
		mids = append(mids, musicData.UpMID)
		res = append(res, newOne)
	}
	if ret, err = p.getUpNames(c, mids); err != nil {
		log.Error("s.BgmBindList(%d,%d,%d) get mid(%v) name error(%v)", aid, cid, mType, mids, err)
		err = nil
	}
	for _, v := range res {
		if name, ok := ret[v.MID]; ok {
			v.Author = name
		}
		resOk = append(resOk, v)
	}
	if len(resOk) > 5 {
		resOk = resOk[:5]
	}
	return
}

// getUpNames fn
func (p *Public) getUpNames(c context.Context, mids []int64) (ret map[int64]string, err error) {
	var (
		minfos map[int64]*accMdl.Info
	)
	ret = make(map[int64]string)
	if len(mids) > 0 {
		minfos, err = p.acc.Infos(c, mids, "localhost")
		if err != nil {
			log.Info("minfos err mids (%+v)|err(%+v)", mids, err)
			return
		}
		for _, info := range minfos {
			ret[info.Mid] = info.Name
		}
	}
	return
}

// FillPayInfo fill pay
func (p *Public) FillPayInfo(c context.Context, a *arcmdl.Archive, ugcPayCfg *conf.UgcPay, ip string) (pay *arcmdl.UgcPayInfo) {
	var (
		err      error
		ass      *arcmdl.PayAsset
		registed bool
	)
	pay = &arcmdl.UgcPayInfo{
		Acts: make(map[string]*arcmdl.PayAct),
	}
	pay.Acts["edit"] = &arcmdl.PayAct{
		State: 1,
	}
	pay.Acts["delete"] = &arcmdl.PayAct{
		State: 1,
	}
	ass, registed, err = p.pay.Ass(c, a.Aid, ip)
	if err != nil {
		log.Error("p.pay.Ass aids (%v), ip(%s) err(%v)", a.Aid, ip, err)
	}
	pay.Asset = ass
	delDeadline := xtime.Time(a.PTime.Time().AddDate(0, 0, ugcPayCfg.AllowDeleteDays).Unix())
	editDeadline := xtime.Time(a.PTime.Time().AddDate(0, 0, ugcPayCfg.AllowEditDays).Unix())
	if !registed {
		pay.Acts["edit"] = &arcmdl.PayAct{
			Reason: "老稿件不允许参与UGC内容付费项目中，请重新投稿",
			State:  0,
		}
	} else {
		if a.UgcPay == 1 {
			if a.CTime != a.PTime &&
				xtime.Time(time.Now().Unix()) < delDeadline {
				pay.Acts["delete"] = &arcmdl.PayAct{
					Reason: fmt.Sprintf("付费稿件必须在开放之后的第%d天才能删除", ugcPayCfg.AllowDeleteDays),
					State:  0,
				}
			}
			if a.CTime != a.PTime &&
				a.State != mdlarc.StateForbidRecicle &&
				xtime.Time(time.Now().Unix()) < editDeadline {
				pay.Acts["edit"] = &arcmdl.PayAct{
					Reason: fmt.Sprintf("付费稿件必须在开放之后的第%d天才能编辑", ugcPayCfg.AllowDeleteDays),
					State:  0,
				}
			}
		}
		// 有注册过，但是现在已经被关闭付费标记的稿件，也允许编辑和删除
	}
	return
}

// loadStaffTitles 拉取联合投稿职能列表
func (p *Public) loadStaffTitles() {
	var (
		c   = context.TODO()
		err error
	)
	if p.StaffTitlesCache, err = p.tag.StaffTitleList(c); err != nil {
		log.Error("p.loadStaffTitles() error(%v)", err)
		return
	}
}
