package http

import (
	"go-common/library/cache/memcache"
	"strings"
	"time"

	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func searInterList(c *bm.Context) {
	var (
		req   = c.Request.Form
		items []*model.SearInter
		total int
		pn    = atoi(req.Get("pn"))
		ps    = atoi(req.Get("ps"))
		pubs  *model.PublishStatus
		err   error
	)
	if pn == 0 {
		pn = 1
	}
	if ps == 0 {
		ps = 20
	}
	if items, total, err = tvSrv.GetSearInterList(c, pn, ps); err != nil {
		log.Error("tvSrv.searInterList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	//rank
	for i := 0; i < len(items); i++ {
		items[i].Rank = int64(i) + 1
	}
	if pubs, err = tvSrv.GetPublishState(c); err != nil {
		if err == memcache.ErrNotFound {
			nowTime := time.Now()
			t := nowTime.Format("2006-01-02 15:04:05")
			pubs = &model.PublishStatus{
				Time:  t,
				State: 0,
			}
		} else {
			log.Error("tvSrv.searInterList GetHotPubState error(%v)", err)
			c.JSON("MC获取发布状态报错", ecode.RequestErr)
			return
		}
	}
	pager := &model.SearInterPager{
		TotalCount: total,
		Pn:         pn,
		Ps:         ps,
		Items:      items,
		PubState:   pubs.State,
		PubTime:    pubs.Time,
	}
	c.JSON(pager, nil)
}

func searInterAdd(c *bm.Context) {
	var (
		req        = c.Request.PostForm
		searchword = req.Get("searchword")
		err        error
		count      int64
		item       model.SearInter
		pubs       *model.PublishStatus
	)
	if err = tvSrv.DB.Model(&model.SearInter{}).Where("deleted!=?", _isDeleted).Count(&count).Error; err != nil {
		log.Error("tvSrv.searInterAdd err(%v)", err)
		c.JSON(nil, err)
		return
	}
	if count >= 20 {
		c.JSON("热词数最多只能添加20条数据", ecode.RequestErr)
		return
	}
	if searchword == "" {
		c.JSON("searchword can not null", ecode.RequestErr)
		return
	}
	exist := &model.SearInter{}
	if err = tvSrv.DB.Where("searchword=?", searchword).Where("deleted!=?", _isDeleted).First(exist).Error; err != nil && err != ecode.NothingFound {
		log.Error("tvSrv.searInterAdd error(%v)", err)
		c.JSON("查找搜索词Mysql报错", ecode.RequestErr)
		return
	}
	if exist.ID != 0 {
		log.Error("searchword is existed, error(%v)", err)
		c.JSON("搜索词已经存在", ecode.RequestErr)
		return
	}
	if item, err = tvSrv.GetMaxRank(c); err != nil && err != ecode.NothingFound {
		log.Error("tvSrv.searInterAdd GetMaxRank error(%v)", err)
		c.JSON("查找最大的排序报错", ecode.RequestErr)
		return
	}
	//default rank is last
	rank := item.Rank + 1
	searchInter := &model.SearInter{
		Searchword: searchword,
		Rank:       rank,
	}
	if err = tvSrv.AddSearInter(c, searchInter); err != nil {
		log.Error("tvSrv.searInterAdd error(%v)", err)
		c.JSON("添加搜索词报错", ecode.RequestErr)
		return
	}
	//get publish state
	if pubs, err = tvSrv.GetPublishState(c); err != nil {
		if err == memcache.ErrNotFound {
			nowTime := time.Now()
			t := nowTime.Format("2006-01-02 15:04:05")
			pubs = &model.PublishStatus{
				Time:  t,
				State: 0,
			}
		} else {
			log.Error("tvSrv.searInterList GetHotPubState error(%v)", err)
			c.JSON("MC获取发布状态报错", ecode.RequestErr)
			return
		}
	}
	//set publish state
	s := &model.PublishStatus{
		Time:  pubs.Time,
		State: 0,
	}
	if err = tvSrv.SetPublishState(c, s); err != nil {
		log.Error("tvSrv.SetPubStat error(%v)", err)
		c.JSON("设置发布状态到MC中报错", ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func searInterEdit(c *bm.Context) {
	var (
		req        = c.Request.PostForm
		id         = parseInt(req.Get("id"))
		searchword = req.Get("searchword")
		err        error
		pubs       *model.PublishStatus
	)
	if req.Get("id") == "" {
		c.JSON("id can no null", ecode.RequestErr)
		return
	}
	if req.Get("searchword") == "" {
		c.JSON("searchword can no null", ecode.RequestErr)
		return
	}
	exist := &model.SearInter{}
	if err = tvSrv.DB.Where("id=?", id).Where("deleted!=?", _isDeleted).First(exist).Error; err != nil {
		log.Error("tvSrv.searInterEdit error(%v)", err)
		c.JSON("can not find value", err)
		return
	}
	exist = &model.SearInter{}
	if err = tvSrv.DB.Where("searchword=?", searchword).Where("deleted!=?", _isDeleted).First(exist).Error; err != nil && err != ecode.NothingFound {
		log.Error("tvSrv.searInterEdit error(%v)", err)
		c.JSON(err, ecode.RequestErr)
		return
	}
	if exist.ID != 0 && exist.ID != id {
		c.JSON("searchword existed", nil)
		return
	}
	if err = tvSrv.UpdateSearInter(c, id, searchword); err != nil {
		log.Error("tvSrv.searInterEdit err(%v)", err)
		c.JSON(nil, err)
		return
	}
	//get publish state
	if pubs, err = tvSrv.GetPublishState(c); err != nil {
		log.Error("tvSrv.searInterEdit GetHotPubState error(%v)", err)
		c.JSON(nil, err)
		return
	}
	//set publish state
	s := &model.PublishStatus{
		Time:  pubs.Time,
		State: 0,
	}
	if err = tvSrv.SetPublishState(c, s); err != nil {
		log.Error("tvSrv.searInterEdit SetPubStat error(%v)", err)
		c.JSON(err, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func searInterDel(c *bm.Context) {
	var (
		req  = c.Request.PostForm
		id   = parseInt(req.Get("id"))
		err  error
		pubs *model.PublishStatus
	)
	if req.Get("id") == "" {
		c.JSON("id can not null", err)
		return
	}
	exist := &model.SearInter{}
	if err = tvSrv.DB.Where("id=?", id).Where("deleted!=?", _isDeleted).First(exist).Error; err != nil {
		c.JSON("can not find value", err)
		return
	}
	if err = tvSrv.DelSearInter(c, id); err != nil {
		log.Error("tvSrv.searInterDel err(%v)", err)
		c.JSON(nil, err)
		return
	}
	//get publish state
	if pubs, err = tvSrv.GetPublishState(c); err != nil {
		log.Error("tvSrv.searInterDel GetHotPubState error(%v)", err)
		c.JSON(nil, err)
		return
	}
	//set publish state
	s := &model.PublishStatus{
		Time:  pubs.Time,
		State: 0,
	}
	if err = tvSrv.SetPublishState(c, s); err != nil {
		log.Error("tvSrv.searInterDel error(%v)", err)
		c.JSON(err, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func searInterRank(c *bm.Context) {
	var (
		req   = c.Request.PostForm
		ids   = req.Get("ids")
		err   error
		pubs  *model.PublishStatus
		total int
	)
	if ids == "" {
		c.JSON("不能发布空数据", ecode.RequestErr)
		return
	}
	idsArr := strings.Split(ids, ",")
	if len(idsArr) <= 0 {
		c.JSON("不能发布空数据", ecode.RequestErr)
		return
	}
	if total, err = tvSrv.GetSearInterCount(c); err != nil {
		log.Error("tvSrv.GetSearInterCount err ", err)
		c.JSON("更新排序失败,GetSearInterCount error ", err)
		return
	}
	if len((idsArr)) != total {
		c.JSON("请提交全部数据", ecode.RequestErr)
		return
	}
	if err = tvSrv.RankSearInter(c, idsArr); err != nil {
		log.Error("tvSrv.searInterRank err(%v),idsArr(%v)", err, idsArr)
		c.JSON("更新排序失败, RankSearIntererror error ", err)
		return
	}
	//get publish state
	if pubs, err = tvSrv.GetPublishState(c); err != nil {
		log.Error("tvSrv.searInterDel GetHotPubState error(%v)", err)
		c.JSON(nil, err)
		return
	}
	//set publish state
	s := &model.PublishStatus{
		Time:  pubs.Time,
		State: 0,
	}
	if err = tvSrv.SetPublishState(c, s); err != nil {
		log.Error("tvSrv.searInterDel error(%v)", err)
		c.JSON(err, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func searInterPublish(c *bm.Context) {
	var (
		items []*model.SearInter
		err   error
	)
	if items, err = tvSrv.GetSearInterPublish(c); err != nil {
		log.Error("tvSrv.searInterPublish GetSearInterPublish error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(items) == 0 {
		c.JSON("不能发布空数据", ecode.RequestErr)
		return
	}
	var rank []*model.OutSearchInter
	for _, v := range items {
		out := &model.OutSearchInter{
			Keyword: v.Searchword,
			Status:  v.Tag,
		}
		rank = append(rank, out)
	}
	if err = tvSrv.SetSearInterRank(c, rank); err != nil {
		log.Error("tvSrv.searInterPublish SearInterRank error(%v)", err)
		c.JSON(nil, err)
		return
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	s := &model.PublishStatus{
		Time:  t,
		State: 1,
	}
	if err = tvSrv.SetPublishState(c, s); err != nil {
		log.Error("tvSrv.searInterPublish  SetPubStat error(%v)", err)
		c.JSON(err, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func searInterPubList(c *bm.Context) {
	var (
		items []*model.OutSearchInter
		err   error
	)
	if items, err = tvSrv.GetSearInterRank(c); err != nil {
		log.Error("tvSrv.searInterListOut error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(items, nil)
}
