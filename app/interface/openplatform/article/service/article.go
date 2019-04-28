package service

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"sync"

	"go-common/app/interface/openplatform/article/dao"
	artmdl "go-common/app/interface/openplatform/article/model"
	account "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"go-common/library/sync/errgroup"
)

var _moreNum = 3

// AddArticleCache adds artmdl.
func (s *Service) AddArticleCache(c context.Context, aid int64) (err error) {
	var a *artmdl.Article
	if a, err = s.dao.Article(c, aid); err != nil {
		dao.PromError("article:新增文章缓存获取文章")
		return
	}
	if a == nil {
		dao.PromError("article:新增文章缓存文章未找到")
		log.Error("s.Article(%d) is blank", aid)
		return
	}
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		return s.dao.AddArticlesMetaCache(errCtx, a.Meta)
	})
	group.Go(func() error {
		return s.dao.AddArticleContentCache(errCtx, aid, a.Content)
	})
	group.Go(func() error {
		if a.Keywords == "" {
			return s.dao.AddArticleKeywordsCache(errCtx, aid, a.Summary)
		}
		return s.dao.AddArticleKeywordsCache(errCtx, aid, a.Keywords)
	})
	group.Go(func() (err error) {
		return s.addUpperCache(errCtx, a.Author.Mid, aid, int64(a.PublishTime))
	})
	group.Go(func() (err error) {
		return s.addArtSortCache(errCtx, a.Meta)
	})
	group.Go(func() (err error) {
		return s.AddCacheHotspotArt(c, s.metaToSearch(c, a.Meta))
	})
	group.Go(func() (err error) {
		if lid, _ := s.rebuildArticleListCache(c, aid); lid > 0 {
			s.updateListInfo(c, lid)
			return s.RebuildListCache(c, lid)
		}
		return
	})
	if err = group.Wait(); err != nil {
		log.Errorv(c, log.KV("log", "AddArticleCache"), log.KV("error", err), log.KV("msg", "group.Wait()"))
		dao.PromError("article:添加文章缓存")
	}
	return
}

// UpdateArticleCache adds artmdl.
func (s *Service) UpdateArticleCache(c context.Context, aid, oldCid int64) (err error) {
	var a *artmdl.Article
	if a, err = s.dao.Article(c, aid); err != nil {
		dao.PromError("article:更新文章缓存获取文章")
		return
	}
	if a == nil {
		return
	}
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		return s.dao.AddArticlesMetaCache(errCtx, a.Meta)
	})
	group.Go(func() error {
		return s.dao.AddArticleContentCache(errCtx, aid, a.Content)
	})
	group.Go(func() error {
		if a.Keywords == "" {
			return s.dao.AddArticleKeywordsCache(errCtx, aid, a.Summary)
		}
		return s.dao.AddArticleKeywordsCache(errCtx, aid, a.Keywords)
	})
	group.Go(func() error {
		if a.AttrVal(artmdl.AttrBitNoDistribute) {
			return s.dao.DelUpperCache(errCtx, a.Meta.Author.Mid, aid)
		}
		return s.addUpperCache(errCtx, a.Meta.Author.Mid, aid, int64(a.Meta.PublishTime))
	})
	group.Go(func() (err error) {
		if err = s.DelCacheHotspotArt(c, aid); err != nil {
			return
		}
		if !a.AttrVal(artmdl.AttrBitNoDistribute) {
			return s.AddCacheHotspotArt(c, s.metaToSearch(c, a.Meta))
		}
		return
	})
	group.Go(func() (err error) {
		var lid int64
		if lid, err = s.rebuildArticleListCache(c, aid); lid > 0 {
			s.updateListInfo(c, lid)
			err = s.RebuildListCache(c, lid)
		}
		return
	})
	group.Go(func() (err error) {
		var root, oldRoot int64
		if root, err = s.CategoryToRoot(a.Category.ID); err != nil {
			dao.PromError("article:更新文章缓存查找分类")
			log.Error("s.CategoryToRoot(%d,%d) error(%+v)", aid, a.Category.ID, err)
			return
		}
		if oldCid == a.Category.ID {
			return nil
		}
		if oldRoot, err = s.CategoryToRoot(oldCid); err != nil {
			dao.PromError("article:更新文章缓存查找分类")
			log.Error("s.CategoryToRoot(%d,%d) error(%+v)", aid, oldCid, err)
			return
		}
		cids := []int64{oldCid}
		if root != oldRoot {
			cids = append(cids, oldRoot)
		}
		if err = s.delArtSortCacheFromCid(errCtx, aid, cids...); err != nil {
			return
		}
		return s.addArtSortCache(errCtx, a.Meta)
	})
	if err = group.Wait(); err != nil {
		log.Errorv(c, log.KV("log", "UpdateArticleCache"), log.KV("error", err), log.KV("msg", "group.Wait()"))
		dao.PromError("article:更新文章缓存")
	}
	return
}

func (s *Service) addUpperCache(c context.Context, mid, aid, ptime int64) (err error) {
	var (
		exists map[int64]bool
		arts   map[int64][][2]int64
	)
	if exists, err = s.dao.ExpireUppersCache(c, []int64{mid}); err != nil {
		return
	}
	if exists[mid] {
		return s.dao.AddUpperCache(c, mid, aid, ptime)
	}
	if arts, err = s.dao.UppersPassed(c, []int64{mid}); err != nil {
		dao.PromError("article:新增文章缓存获取up过审")
		return
	}
	return s.dao.AddUpperCaches(c, arts)
}

//RootCategory 找到一级分区
func (s *Service) RootCategory(c context.Context, aid int64) (root int64, cid int64, err error) {
	var art *artmdl.Meta
	if art, err = s.dao.AllArticleMeta(c, aid); err != nil {
		return
	}
	if art == nil {
		err = ecode.NothingFound
		return
	}
	cid = art.Category.ID
	root, err = s.CategoryToRoot(art.Category.ID)
	return
}

// CategoryToRoot 找到一级分区
func (s *Service) CategoryToRoot(cid int64) (res int64, err error) {
	for (s.categoriesMap[cid] != nil) && (s.categoriesMap[cid].ParentID != _recommendCategory) {
		cid = s.categoriesMap[cid].ParentID
	}
	if (s.categoriesMap[cid] == nil) || (s.categoriesMap[cid].ParentID != _recommendCategory) {
		err = ecode.ArtNoCategory
		return
	}
	res = cid
	return
}

// DelArticleCache deletes artmdl.
func (s *Service) DelArticleCache(c context.Context, mid, aid int64) (err error) {
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		return s.dao.DelUpperCache(errCtx, mid, aid)
	})
	group.Go(func() error {
		return s.dao.DelArticleMetaCache(errCtx, aid)
	})
	group.Go(func() error {
		return s.dao.DelArticleContentCache(errCtx, aid)
	})
	group.Go(func() error {
		return s.DelCacheHotspotArt(c, aid)
	})
	group.Go(func() error {
		return s.dao.DelArticleStatsCache(errCtx, aid)
	})
	group.Go(func() (err error) {
		err = s.delArtSortCache(errCtx, aid)
		return
	})
	group.Go(func() (err error) {
		var lid int64
		if lid, err = s.rebuildArticleListCache(c, aid); lid > 0 {
			s.updateListInfo(c, lid)
			err = s.RebuildListCache(c, lid)
		}
		return
	})
	if err = group.Wait(); err != nil {
		log.Errorv(c, log.KV("log", "DelArticleCache"), log.KV("error", err), log.KV("msg", "group.Wait()"))
		dao.PromError("article:删除文章缓存")
	}
	return
}

// Article  get article
func (s *Service) Article(c context.Context, id int64) (res *artmdl.Article, err error) {
	var am *artmdl.Meta
	if am, err = s.ArticleMeta(c, id); err != nil || am == nil {
		return
	}
	res = &artmdl.Article{Meta: am}
	if res.Content, err = s.content(c, id); err != nil {
		return
	}
	res.Keywords = s.keywords(c, id, res.Summary)
	if res.Content == "" {
		dao.PromError("article:文章内容为空")
		log.Error("s.Article(%v) content is blank", id)
	}
	s.media(c, am)
	log.Info("s.Article() aid(%d) title(%s) content length(%d)", res.ID, res.Title, len(res.Content))
	return
}

// MediaCategory .
func (s *Service) MediaCategory(c context.Context, mediaID int64, mid int64) (cg *artmdl.Category, err error) {
	var (
		res *artmdl.MediaResult
		cid int64
	)
	if mediaID == 0 {
		return
	}
	if res, err = s.dao.Media(c, mediaID, mid); err != nil {
		log.Error("s.MediaCategory(%d) get media info failed: %v", mediaID, err)
		dao.PromError("article:番剧信息获取失败")
		return
	}
	if res.Media.TypeID > 4 || res.Media.TypeID == 0 || len(s.c.Article.Media) < 5 {
		err = errors.New("番剧类型错误或者未配置番剧类别")
		return
	}
	cid = s.c.Article.Media[res.Media.TypeID]
	cg = s.categoriesMap[cid]
	return
}

func (s *Service) media(c context.Context, am *artmdl.Meta) {
	var (
		res *artmdl.MediaResult
		err error
	)
	if am.Media == nil || am.Media.MediaID == 0 {
		return
	}
	if res, err = s.dao.Media(c, am.Media.MediaID, am.Author.Mid); err != nil {
		log.Error("s.media(%d) get media info failed: %v", am.Media.MediaID, err)
		dao.PromError("article:番剧信息获取失败")
		return
	}
	am.Media.MediaID = res.Media.MediaID
	am.Media.Score = res.Score
	am.Media.Title = res.Media.Title
	am.Media.Cover = res.Media.Cover
	am.Media.Area = res.Media.Area
	am.Media.TypeID = res.Media.TypeID
	am.Media.TypeName = res.Media.TypeName
	return
}

// ArticleMeta gets article's meta.
func (s *Service) ArticleMeta(c context.Context, aid int64) (res *artmdl.Meta, err error) {
	var addCache = true
	if res, err = s.dao.ArticleMetaCache(c, aid); err != nil {
		addCache = false
		err = nil
	}
	if res == nil {
		if res, err = s.dao.ArticleMeta(c, aid); err != nil || res == nil {
			return
		}
	}
	if s.categoriesMap[res.Category.ID] != nil {
		res.Category = s.categoriesMap[res.Category.ID]
		res.Categories = s.categoryParents[res.Category.ID]
	}
	group := &errgroup.Group{}
	// get author
	group.Go(func() error {
		var author *artmdl.Author
		if author, _ = s.author(c, res.Author.Mid); author != nil {
			res.Author = author
		}
		return nil
	})
	// get stats
	group.Go(func() error {
		var stat *artmdl.Stats
		if stat, _ = s.stat(c, aid); stat != nil {
			res.Stats = stat
			return nil
		}
		if res.Stats == nil {
			res.Stats = new(artmdl.Stats)
		}
		return nil
	})
	// get tag
	group.Go(func() error {
		var tags []*artmdl.Tag
		if tags, _ = s.Tags(c, aid, false); len(tags) > 0 {
			res.Tags = tags
			return nil
		}
		if len(res.Tags) == 0 {
			res.Tags = []*artmdl.Tag{}
		}
		return nil
	})
	// get list
	group.Go(func() (err error) {
		lists, _ := s.dao.ArtsList(c, []int64{aid})
		res.List = lists[aid]
		return
	})
	group.Wait()
	if addCache {
		cache.Save(func() { s.dao.AddArticlesMetaCache(context.TODO(), res) })
	}
	return
}

func (s *Service) accountInfo(c context.Context, mid int64) (info *account.Card, err error) {
	var (
		arg = &account.ArgMid{Mid: mid}
	)
	if info, err = s.accountRPC.Card3(c, arg); err != nil {
		dao.PromError("article:获取作者信息")
		log.Error("s.accountRPC.Card3(%+v) error(%+v)", arg, err)
		return
	}
	return
}

func (s *Service) author(c context.Context, mid int64) (res *artmdl.Author, err error) {
	var (
		card *account.Card
		arg  = &account.ArgMid{Mid: mid}
	)
	if card, err = s.accountRPC.Card3(c, arg); err != nil {
		dao.PromError("article:获取作者信息")
		log.Error("s.accountRPC.Info(%+v) error(%+v)", arg, err)
		return
	}
	res = &artmdl.Author{
		Mid:  mid,
		Name: card.Name,
		Face: card.Face,
		Pendant: artmdl.Pendant{
			Pid:    int32(card.Pendant.Pid),
			Name:   card.Pendant.Name,
			Image:  card.Pendant.Image,
			Expire: int32(card.Pendant.Expire),
		},
		Nameplate: artmdl.Nameplate{
			Nid:        card.Nameplate.Nid,
			Name:       card.Nameplate.Name,
			Image:      card.Nameplate.Image,
			ImageSmall: card.Nameplate.ImageSmall,
			Level:      card.Nameplate.Level,
			Condition:  card.Nameplate.Condition,
		},
		Vip: artmdl.VipInfo{
			Type:   card.Vip.Type,
			Status: card.Vip.Status,
		},
	}
	if card.Official.Role == 0 {
		res.OfficialVerify.Type = -1
	} else {
		if card.Official.Role <= 2 {
			res.OfficialVerify.Type = 0
		} else {
			res.OfficialVerify.Type = 1
		}
		res.OfficialVerify.Desc = card.Official.Title
	}
	return
}

func (s *Service) authors(c context.Context, mids []int64) (res map[int64]*artmdl.Author, err error) {
	res = make(map[int64]*artmdl.Author)
	if len(mids) == 0 {
		return
	}
	var (
		cards map[int64]*account.Card
		arg   = &account.ArgMids{Mids: mids}
	)
	if cards, err = s.accountRPC.Cards3(c, arg); err != nil {
		dao.PromError("article:批量获取作者信息")
		log.Error("s.accountRPC.Infos(%+v) error(%+v)", arg, err)
		return
	}
	for mid, card := range cards {
		au := &artmdl.Author{
			Mid:  mid,
			Name: card.Name,
			Face: card.Face,
			Pendant: artmdl.Pendant{
				Pid:    int32(card.Pendant.Pid),
				Name:   card.Pendant.Name,
				Image:  card.Pendant.Image,
				Expire: int32(card.Pendant.Expire),
			},
			Nameplate: artmdl.Nameplate{
				Nid:        card.Nameplate.Nid,
				Name:       card.Nameplate.Name,
				Image:      card.Nameplate.Image,
				ImageSmall: card.Nameplate.ImageSmall,
				Level:      card.Nameplate.Level,
				Condition:  card.Nameplate.Condition,
			},
			Vip: artmdl.VipInfo{
				Type:   card.Vip.Type,
				Status: card.Vip.Status,
			},
		}
		if card.Official.Role == 0 {
			au.OfficialVerify.Type = -1
		} else {
			if card.Official.Role <= 2 {
				au.OfficialVerify.Type = 0
			} else {
				au.OfficialVerify.Type = 1
			}
			au.OfficialVerify.Desc = card.Official.Title
		}
		res[mid] = au
	}
	return
}

func (s *Service) authorDetail(c context.Context, mid int64) (res *artmdl.Author, err error) {
	var (
		card *account.Card
		arg  = &account.ArgMid{Mid: mid}
	)
	if card, err = s.accountRPC.Card3(c, arg); err != nil {
		dao.PromError("article:单个获取作者信息")
		log.Error("s.accountRPC.Info(%+v) error(%+v)", arg, err)
		return
	}
	if card == nil {
		return
	}
	res = &artmdl.Author{
		Mid:  mid,
		Name: card.Name,
		Face: card.Face,
		Pendant: artmdl.Pendant{
			Pid:    int32(card.Pendant.Pid),
			Name:   card.Pendant.Name,
			Image:  card.Pendant.Image,
			Expire: int32(card.Pendant.Expire),
		},
		Nameplate: artmdl.Nameplate{
			Nid:        card.Nameplate.Nid,
			Name:       card.Nameplate.Name,
			Image:      card.Nameplate.Image,
			ImageSmall: card.Nameplate.ImageSmall,
			Level:      card.Nameplate.Level,
			Condition:  card.Nameplate.Condition,
		},
	}
	if card.Official.Role == 0 {
		res.OfficialVerify.Type = -1
	} else {
		if card.Official.Role <= 2 {
			res.OfficialVerify.Type = 0
		} else {
			res.OfficialVerify.Type = 1
		}
		res.OfficialVerify.Desc = card.Official.Title
	}
	return
}

func (s *Service) content(c context.Context, aid int64) (res string, err error) {
	var addCache = true
	if res, err = s.dao.ArticleContentCache(c, aid); err != nil {
		addCache = false
		err = nil
	} else if res != "" {
		return
	}
	if res, err = s.dao.ArticleContent(c, aid); err != nil {
		dao.PromError("article:稿件内容")
		return
	}
	if addCache && res != "" {
		cache.Save(func() {
			s.dao.AddArticleContentCache(context.TODO(), aid, res)
		})
	}
	return
}

// ListCategories list categories
func (s *Service) ListCategories(c context.Context, ip string) (res artmdl.Categories, err error) {
	if len(s.Categories) == 0 {
		err = ecode.NothingFound
		return
	}
	res = s.Categories
	return
}

// ListCategoriesMap list category map
func (s *Service) ListCategoriesMap(c context.Context, ip string) (res map[int64]*artmdl.Category, err error) {
	if len(s.Categories) == 0 {
		err = ecode.NothingFound
		return
	}
	res = s.categoriesMap
	return
}

// ArticleMetas get article meta
func (s *Service) ArticleMetas(c context.Context, ids []int64) (res map[int64]*artmdl.Meta, err error) {
	var (
		addCache                 = true
		group                    *errgroup.Group
		cachedMetas, missedMetas map[int64]*artmdl.Meta
		missedMetaIDs, resIDs    []int64
		mutex                    = &sync.Mutex{}
	)
	res = make(map[int64]*artmdl.Meta)
	// get meta
	if cachedMetas, missedMetaIDs, err = s.dao.ArticlesMetaCache(c, ids); err != nil {
		addCache = false
		err = nil
	}
	if len(missedMetaIDs) > 0 {
		missedMetas, _ = s.dao.ArticleMetas(c, missedMetaIDs)
	}
	// 合并缓存和回源的数据
	for id, artm := range cachedMetas {
		res[id] = artm
		resIDs = append(resIDs, id)
	}
	for id, artm := range missedMetas {
		res[id] = artm
		resIDs = append(resIDs, id)
	}
	// 更新分类
	for id, art := range res {
		if art.Category == nil {
			continue
		}
		if s.categoriesMap[art.Category.ID] != nil {
			res[id].Category = s.categoriesMap[art.Category.ID]
			res[id].Categories = s.categoryParents[art.Category.ID]
		}
	}
	group = &errgroup.Group{}
	// get author
	group.Go(func() (err error) {
		var (
			mids       []int64
			authors    map[int64]*artmdl.Author
			authorsMap = make(map[int64]bool)
		)
		for _, art := range res {
			authorsMap[art.Author.Mid] = true
		}
		for id := range authorsMap {
			mids = append(mids, id)
		}
		if authors, err = s.authors(c, mids); err != nil {
			dao.PromError("article:稿件获取作者信息")
			err = nil
			return
		}
		mutex.Lock()
		for _, art := range res {
			author := authors[art.Author.Mid]
			if author != nil {
				art.Author = author
			}
		}
		mutex.Unlock()
		return
	})
	//get stats
	group.Go(func() (err error) {
		stats, _ := s.stats(c, resIDs)
		mutex.Lock()
		for id := range res {
			s := stats[id]
			if s == nil {
				s = new(artmdl.Stats)
			}
			res[id].Stats = s
		}
		mutex.Unlock()
		return
	})
	// get list
	group.Go(func() (err error) {
		lists, _ := s.dao.ArtsList(c, resIDs)
		mutex.Lock()
		for id := range res {
			res[id].List = lists[id]
		}
		mutex.Unlock()
		return
	})
	group.Wait()
	if addCache && len(missedMetas) > 0 {
		cache.Save(func() {
			for _, art := range missedMetas {
				s.dao.AddArticlesMetaCache(context.TODO(), art)
			}
		})
	}
	return
}

func filterNoDistributeArts(as []*artmdl.Meta) (res []*artmdl.Meta) {
	for _, a := range as {
		if (a != nil) && !a.AttrVal(artmdl.AttrBitNoDistribute) {
			res = append(res, a)
		}
	}
	return
}

func filterNoDistributeArtsMap(as map[int64]*artmdl.Meta) {
	for id, a := range as {
		if (a != nil) && a.AttrVal(artmdl.AttrBitNoDistribute) {
			delete(as, id)
		}
	}
}

// AddArtContentCache add article content cache
func (s *Service) AddArtContentCache(c context.Context, aid int64, content string) (err error) {
	if content == "" {
		return
	}
	err = s.dao.AddArticleContentCache(c, aid, content)
	return
}

// ArticleRemainCount returns the number that user could be use to posting new articles.
func (s *Service) ArticleRemainCount(c context.Context, mid int64) (num int, err error) {
	if mid <= 0 {
		return
	}
	var count, limit int
	if count, err = s.dao.ArticleRemainCount(c, mid); err != nil {
		return
	}
	author, _ := s.dao.Author(c, mid)
	if author != nil {
		limit = author.Limit
	}
	if limit == 0 {
		limit = s.c.Article.UpperArticleLimit
	}
	if count > limit {
		return
	}
	num = limit - count
	return
}

// AddComplaint add complaint.
func (s *Service) AddComplaint(c context.Context, aid, mid, ctype int64, reason, imageUrls, ip string) (err error) {
	var exist, protected bool
	if exist, err = s.dao.ComplaintExist(c, aid, mid); (err != nil) || exist {
		return
	}
	if err = s.dao.AddComplaint(c, aid, mid, ctype, reason, imageUrls); err != nil {
		return
	}
	if protected, err = s.dao.ComplaintProtected(c, aid); err != nil || protected {
		return
	}
	err = s.dao.AddComplaintCount(c, aid)
	return
}

// MoreArts get author's more articles.
func (s *Service) MoreArts(c context.Context, aid int64) (res []*artmdl.Meta, err error) {
	var am *artmdl.Meta
	if am, err = s.ArticleMeta(c, aid); err != nil {
		dao.PromError("article:获取文章meta")
		return
	}
	if am == nil || am.Author == nil {
		return
	}
	var (
		exist                          bool
		beforeAids, afterAids, tmpAids []int64
		aidTimes                       [][2]int64
		addCache                       = true
		tmpRes                         map[int64]*artmdl.Meta
		mid                            = am.Author.Mid
	)
	if exist, err = s.dao.ExpireUpperCache(c, mid); err != nil {
		addCache = false
		err = nil
	} else if exist {
		if beforeAids, afterAids, err = s.dao.MoreArtsCaches(c, mid, int64(am.PublishTime), _moreNum+4); err != nil {
			addCache = false
			exist = false
		}
	}
	if !exist {
		if aidTimes, err = s.dao.UpperPassed(c, mid); err != nil {
			return
		}
		if addCache {
			cache.Save(func() {
				s.dao.AddUpperCaches(context.TODO(), map[int64][][2]int64{mid: aidTimes})
			})
		}
		for _, aidTime := range aidTimes {
			tmpAids = append(tmpAids, aidTime[0])
		}
		beforeAids, afterAids = splitAids(tmpAids, aid)
	}
	if len(beforeAids)+len(afterAids) == 0 {
		return
	}
	tmpAids = append([]int64{}, beforeAids...)
	tmpAids = append(tmpAids, afterAids...)
	if tmpRes, err = s.ArticleMetas(c, tmpAids); err != nil {
		return
	}
	filterNoDistributeArtsMap(tmpRes)
	res = fmtMoreArts(beforeAids, afterAids, tmpRes)
	return
}

func splitAids(aids []int64, aid int64) (beforeAids []int64, afterAids []int64) {
	position := -1
	for i, a := range aids {
		if a == aid {
			position = i
			break
		}
	}
	if position == -1 {
		return
	}
	l := len(aids)
	if position+_moreNum > l {
		beforeAids = aids[position+1 : l]
	} else {
		beforeAids = aids[position+1 : position+_moreNum]
	}
	if position-_moreNum < 0 {
		afterAids = aids[0:position]
	} else {
		afterAids = aids[position-_moreNum : position]
	}
	return
}

// 位置说明: 按照ptime逆序: afterAids -> (aid) -> beforeAids, after取一个 before取2个
func fmtMoreArts(beforeAids, afterAids []int64, tmpRes map[int64]*artmdl.Meta) (res []*artmdl.Meta) {
	var before, after []*artmdl.Meta
	for _, aid := range beforeAids {
		if v := tmpRes[aid]; v != nil {
			before = append(before, v)
		}
	}
	sort.Sort(artmdl.Metas(before))
	for _, aid := range afterAids {
		if v := tmpRes[aid]; v != nil {
			after = append(after, v)
		}
	}
	sort.Sort(artmdl.Metas(after))
	if len(after) > _moreNum {
		after = after[len(after)-_moreNum:]
	}
	if len(before) > _moreNum {
		before = before[:_moreNum]
	}
	lenAfter := len(after)
	lenBefore := len(before)
	// 正常逻辑: 前一个后三个
	if (lenAfter > 0) && (lenBefore > 1) {
		res = []*artmdl.Meta{after[lenAfter-1]}
		res = append(res, before[:2]...)
	} else if lenAfter == 0 {
		// 前面不够
		res = before
	} else if lenBefore == 0 {
		// 后面不够
		res = after
	} else {
		//前面补全后面缺失的
		if lenAfter-(_moreNum-lenBefore) > 0 {
			res = append(res, after[lenAfter-(_moreNum-lenBefore):]...)
		} else {
			res = append(res, after...)
		}
		res = append(res, before...)
	}
	return
}

// NewArticleCount get new article count
func (s *Service) NewArticleCount(c context.Context, ptime int64) (res int64, err error) {
	res, err = s.dao.NewArticleCount(c, ptime)
	return
}

// Mores get more articles
func (s *Service) Mores(c context.Context, aid, loginMID int64) (res artmdl.MoreArts, err error) {
	var art *artmdl.Meta
	group := &errgroup.Group{}
	if art, err = s.ArticleMeta(c, aid); (err != nil) || (art == nil) || (!art.IsNormal()) {
		err = ecode.NothingFound
		return
	}
	res.Author = &artmdl.AccountCard{Mid: strconv.FormatInt(art.Author.Mid, 10), Face: art.Author.Face, Name: art.Author.Name}
	mid := art.Author.Mid
	group.Go(func() (e error) {
		var profile *account.ProfileStat
		if profile, e = s.accountRPC.ProfileWithStat3(c, &account.ArgMid{Mid: mid}); e != nil {
			dao.PromError("article:Card3")
			log.Error("s.acc.Card3(%d) error %v", mid, e)
			return
		}
		if profile != nil {
			res.Author.FromProfileStat(profile)
		}
		return
	})
	group.Go(func() (e error) {
		if res.Articles, e = s.similars(c, aid, art.Category.ID); e != nil {
			dao.PromError("article:similars")
		}
		if res.Articles != nil && len(res.Articles) > 0 {
			return
		}
		if res.Articles, e = s.MoreArts(c, aid); e != nil {
			dao.PromError("article:MoreArts")
		}
		if res.Articles == nil {
			res.Articles = []*artmdl.Meta{}
		}
		return
	})
	group.Go(func() (e error) {
		if res.Total, e = s.UpperArtsCount(c, mid); e != nil {
			dao.PromError("article:获取作者文章数")
		}
		return
	})
	group.Go(func() error {
		// read count
		if stat, e := s.dao.CacheUpStatDaily(c, mid); e != nil {
			dao.PromError("article:CacheUpStatDaily")
		} else if stat != nil {
			res.ReadCount = stat.View
			return nil
		}
		if stat, e := s.dao.UpStat(c, mid); e != nil {
			dao.PromError("article:获取作者文章数")
		} else {
			res.ReadCount = stat.View
		}
		return nil
	})
	group.Go(func() (e error) {
		if loginMID == 0 {
			return
		}
		if res.Attention, e = s.isAttention(c, loginMID, mid); e != nil {
			dao.PromError("article:获取作者文章数")
		}
		return
	})
	group.Wait()
	return
}

// FeedArticleMetas .
func (s *Service) FeedArticleMetas(c context.Context, ids []int64) (res map[int64]*artmdl.Meta, err error) {
	var (
		addCache                 = true
		group                    *errgroup.Group
		cachedMetas, missedMetas map[int64]*artmdl.Meta
		missedMetaIDs, resIDs    []int64
		mutex                    = &sync.Mutex{}
	)
	res = make(map[int64]*artmdl.Meta)
	// get meta
	if cachedMetas, missedMetaIDs, err = s.dao.ArticlesMetaCache(c, ids); err != nil {
		addCache = false
		err = nil
	}
	if len(missedMetaIDs) > 0 {
		missedMetas, _ = s.dao.ArticleMetas(c, missedMetaIDs)
	}
	// 合并缓存和回源的数据
	for id, artm := range cachedMetas {
		res[id] = artm
		resIDs = append(resIDs, id)
	}
	for id, artm := range missedMetas {
		res[id] = artm
		resIDs = append(resIDs, id)
	}
	// 更新分类
	for id, art := range res {
		if art.Category == nil {
			continue
		}
		if s.categoriesMap[art.Category.ID] != nil {
			res[id].Category = s.categoriesMap[art.Category.ID]
			res[id].Categories = s.categoryParents[art.Category.ID]
		}
	}
	group = &errgroup.Group{}
	// get author
	group.Go(func() (err error) {
		var (
			mids       []int64
			authors    map[int64]*artmdl.Author
			authorsMap = make(map[int64]bool)
		)
		for _, art := range missedMetas {
			authorsMap[art.Author.Mid] = true
		}
		for id := range authorsMap {
			mids = append(mids, id)
		}
		if authors, err = s.authors(c, mids); err != nil {
			dao.PromError("article:稿件获取作者信息")
			err = nil
			return
		}
		mutex.Lock()
		for _, art := range missedMetas {
			author := authors[art.Author.Mid]
			if author != nil {
				art.Author = author
			}
		}
		mutex.Unlock()
		return
	})
	//get stats
	group.Go(func() (err error) {
		stats, _ := s.stats(c, resIDs)
		mutex.Lock()
		for id := range res {
			s := stats[id]
			if s == nil {
				s = new(artmdl.Stats)
			}
			res[id].Stats = s
		}
		mutex.Unlock()
		return
	})
	group.Wait()
	if addCache && len(missedMetas) > 0 {
		cache.Save(func() {
			for _, art := range missedMetas {
				s.dao.AddArticlesMetaCache(context.TODO(), art)
			}
		})
	}
	return
}

// AddCheatFilter .
func (s *Service) AddCheatFilter(c context.Context, aid int64, lv int) (err error) {
	return s.dao.AddCheatFilter(c, aid, lv)
}

// DelCheatFilter .
func (s *Service) DelCheatFilter(c context.Context, aid int64) (err error) {
	return s.dao.DelCheatFilter(c, aid)
}

// similars .
func (s *Service) similars(c context.Context, aid int64, cid int64) (res []*artmdl.Meta, err error) {
	var (
		tags  []*artmdl.Tag
		aidst = make(map[int64]bool)
		aidsr []int64
		aidsa = make(map[int64]bool)
		tmps  []int64
		aids  []int64
		id    int64
		meta  *artmdl.Meta
		metas map[int64]*artmdl.Meta
		nils  []int64
	)
	if tags, err = s.Tags(c, aid, false); err != nil {
		return
	}
	for _, tag := range tags {
		var tagArts *artmdl.TagArts
		if tagArts, err = s.dao.CacheAidsByTag(c, tag.Tid); err != nil {
			return
		}
		if tagArts == nil {
			nils = append(nils, tag.Tid)
			continue
		}
		for _, id = range tagArts.Aids {
			aidst[id] = true
		}
	}
	if len(nils) > 0 {
		if tmps, err = s.dao.TagArticles(c, nils); err != nil {
			return
		}
		for _, id = range tmps {
			aidst[id] = true
		}
	}
	delete(aidst, aid)
	aidsr = s.getRecommentsGroups(c, cid, aid)
	if len(aidsr) == 0 && len(aidst) == 0 {
		return
	}
	if len(aidsr) == 0 {
		for id = range aidst {
			aids = append(aids, id)
			if len(aids) > 10 {
				break
			}
		}
		if metas, err = s.ArticleMetas(c, aids); err != nil {
			return
		}
		for _, meta := range metas {
			if artmdl.NoDistributeAttr(meta.Attributes) || artmdl.NoRegionAttr(meta.Attributes) {
				continue
			}
			res = append(res, meta)
			if len(res) == 3 {
				break
			}
		}
		return
	}
	for _, id = range aidsr {
		aids = append(aids, id)
		aidsa[id] = true
	}
	for id = range aidst {
		aidsa[id] = true
	}
	if metas, err = s.ArticleMetas(c, aids); err != nil {
		return
	}
	for _, meta = range metas {
		res = append(res, meta)
		delete(aidsa, meta.ID)
		break
	}
	tmps = []int64{}
	for id = range aidsa {
		tmps = append(tmps, id)

		if len(tmps) > 10 {
			break
		}
	}
	if metas, err = s.ArticleMetas(c, tmps); err != nil {
		return
	}
	for _, meta := range metas {
		if artmdl.NoDistributeAttr(meta.Attributes) || artmdl.NoRegionAttr(meta.Attributes) {
			continue
		}
		res = append(res, meta)
		if len(res) == 3 {
			break
		}
	}
	return
}

// MediaArticle .
func (s *Service) MediaArticle(c context.Context, mediaID, mid int64) (id int64, err error) {
	return s.dao.MediaArticle(c, mediaID, mid)
}

// MediaIDByID .
func (s *Service) MediaIDByID(c context.Context, aid int64) (id int64, err error) {
	return s.dao.MediaIDByID(c, aid)
}

func (s *Service) keywords(c context.Context, id int64, summary string) (keywords string) {
	var (
		err      error
		addCache = true
	)
	if keywords, err = s.dao.ArticleKeywordsCache(c, id); keywords != "" {
		return
	}
	if err != nil {
		addCache = false
	}
	if keywords, err = s.dao.ArticleKeywords(c, id); err != nil {
		dao.PromError("article:文章关键词")
	}
	if keywords == "" {
		keywords = summary
	}
	if addCache {
		cache.Save(func() {
			s.dao.AddArticleKeywordsCache(context.TODO(), id, keywords)
		})
	}
	return
}
