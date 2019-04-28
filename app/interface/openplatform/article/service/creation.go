package service

import (
	"context"
	"sort"

	"go-common/app/interface/openplatform/article/dao"
	"go-common/app/interface/openplatform/article/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"

	"go-common/library/sync/errgroup"
)

const (
	_sortDefault    = 0
	_sortByCtime    = 1
	_sortByLike     = 2
	_sortByReply    = 3
	_sortByView     = 4
	_sortByFavorite = 5
	_sortByCoin     = 6

	//_stateCancel 取消活动
	// _stateCancel = -1
	//_stateJoin 参加活动
	_stateJoin = 0

	_editTimes = 1
)

// checkPrivilege check that whether user has permission to write article.
func (s *Service) checkPrivilege(c context.Context, mid int64) (err error) {
	if res, _, e := s.IsAuthor(c, mid); !res {
		err = ecode.ArtCreationNoPrivilege
		if e != nil {
			dao.PromError("creation:检查作者权限")
		}
	}
	return
}

func (s *Service) checkArtAuthor(c context.Context, aid, mid int64) (am *model.Meta, err error) {
	if am, err = s.creationArticleMeta(c, aid); err != nil {
		return
	} else if am == nil {
		err = ecode.NothingFound
		return
	}
	if am.Author.Mid != mid {
		err = ecode.ArtCreationMIDErr
	}
	return
}

// CreationArticle get creation article
func (s *Service) creationArticleMeta(c context.Context, aid int64) (am *model.Meta, err error) {
	if am, err = s.dao.CreationArticleMeta(c, aid); err != nil {
		dao.PromError("creation:获取文章meta")
		return
	} else if am == nil {
		return
	}
	if s.categoriesMap[am.Category.ID] != nil {
		am.Category = s.categoriesMap[am.Category.ID]
	}
	var author *model.Author
	if author, _ = s.author(c, am.Author.Mid); author != nil {
		am.Author = author
	}
	return
}

// CreationArticle .
func (s *Service) CreationArticle(c context.Context, aid, mid int64) (a *model.Article, err error) {
	if err = s.checkPrivilege(c, mid); err != nil {
		return
	}
	var (
		content string
		am      *model.Meta
	)
	a = &model.Article{}
	if am, err = s.checkArtAuthor(c, aid, mid); err != nil {
		return
	}
	if am.State == model.StateRePending || am.State == model.StateReReject {
		if a, err = s.ArticleVersion(c, aid); err != nil {
			return
		}
		a.Author = am.Author
		a.Category = s.categoriesMap[a.Category.ID]
		a.List, _ = s.dao.ArtList(c, aid)
		a.Stats, _ = s.stat(c, aid)
		return
	}
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() (err error) {
		content, err = s.dao.CreationArticleContent(c, aid)
		return
	})
	group.Go(func() (err error) {
		am.Stats, _ = s.stat(errCtx, aid)
		return
	})
	group.Go(func() (err error) {
		am.Tags, _ = s.Tags(errCtx, aid, true)
		return
	})
	group.Go(func() (err error) {
		am.List, _ = s.dao.ArtList(errCtx, aid)
		return
	})
	if err = group.Wait(); err != nil {
		return
	}
	a.Meta = am
	a.Content = content
	log.Info("s.CreationArticle() aid(%d) title(%s) content length(%d)", a.ID, a.Title, len(a.Content))
	return
}

// AddArticle adds article.
func (s *Service) AddArticle(c context.Context, a *model.Article, actID, listID int64, ip string) (id int64, err error) {
	log.Info("s.AddArticle() aid(%d) title(%s) actID(%d) content length(%d)", a.ID, a.Title, actID, len(a.Content))
	defer func() {
		if err != nil && err != ecode.CreativeArticleCanNotRepeat {
			s.dao.DelSubmitCache(c, a.Author.Mid, a.Title)
		}
	}()
	if err = s.checkPrivilege(c, a.Author.Mid); err != nil {
		return
	}
	a.Content = xssFilter(a.Content)
	if err = s.preArticleCheck(c, a); err != nil {
		return
	}
	var num int
	if num, err = s.ArticleRemainCount(c, a.Author.Mid); err != nil {
		return
	} else if num <= 0 {
		err = ecode.ArtCreationArticleFull
		return
	}
	if s.checkList(c, a.Author.Mid, listID); err != nil {
		return
	}
	a.State = model.StatePending
	var tx *sql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("tx.BeginTran() error(%+v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%+v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			dao.PromError("creation:添加文章")
			log.Error("tx.Commit() error(%+v)", err)
			return
		}
		// id is article id
		if e1 := s.creativeAddArticleList(c, a.Meta.Author.Mid, listID, id, false); e1 != nil {
			dao.PromError("creation:添加文章绑定文集")
			log.Errorv(c, log.KV("log", "creativeAddArticleList"), log.KV("mid", a.Meta.Author.Mid), log.KV("listID", listID), log.KV("article_id", id), log.KV("error", err))
		}
	}()
	// del draft
	if err = s.dao.TxDeleteArticleDraft(c, tx, a.Author.Mid, a.ID); err != nil {
		dao.PromError("creation:删除草稿")
		return
	}
	if id, err = s.dao.TxAddArticleMeta(c, tx, a.Meta, actID); err != nil {
		dao.PromError("creation:删除文章meta")
		return
	}
	keywords, _ := s.Segment(c, int32(id), a.Content, 1, "article")
	if err = s.dao.TxAddArticleContent(c, tx, id, a.Content, keywords); err != nil {
		dao.PromError("creation:添加文章content")
		return
	}
	// add version
	if err = s.dao.TxAddArticleVersion(c, tx, id, a, actID); err != nil {
		dao.PromError("creation:添加历史记录")
		return
	}
	if actID > 0 {
		if e := s.dao.HandleActivity(c, a.Author.Mid, id, actID, _stateJoin, ip); e != nil {
			log.Error("creation: s.act.HandleActivity mid(%d) aid(%d) actID(%d) ip(%s) error(%+v)", a.Author.Mid, id, actID, ip, e)
		}
	}
	var tags []string
	for _, t := range a.Tags {
		tags = append(tags, t.Name)
	}
	if err = s.BindTags(c, a.Meta.Author.Mid, id, tags, ip, actID); err != nil {
		dao.PromError("creation:发布文章绑定tag")
	}
	return
}

// UpdateArticle update article.
func (s *Service) UpdateArticle(c context.Context, a *model.Article, actID, listID int64, ip string) (err error) {
	log.Info("s.UpdateArticle() aid(%d) title(%s) content length(%d)", a.ID, a.Title, len(a.Content))
	if err = s.checkPrivilege(c, a.Author.Mid); err != nil {
		return
	}
	a.Content = xssFilter(a.Content)
	if err = s.preArticleCheck(c, a); err != nil {
		return
	}
	a.State = model.StatePending
	if a.ID <= 0 {
		err = ecode.ArtCreationIDErr
		return
	}
	if _, err = s.checkArtAuthor(c, a.ID, a.Author.Mid); err != nil {
		return
	}
	//这里转换成 -2 2 几种状态
	if err = s.convertState(c, a); err != nil {
		return
	}
	if a.State == model.StateRePending {
		err = s.updateArticleVersion(c, a, actID)
		return
	}
	if err = s.updateArticleDB(c, a); err != nil {
		return
	}
	var tags []string
	for _, t := range a.Tags {
		tags = append(tags, t.Name)
	}
	s.BindTags(c, a.Meta.Author.Mid, a.Meta.ID, tags, ip, actID)
	s.CreativeUpdateArticleList(c, a.Meta.Author.Mid, a.ID, listID, false)
	return
}

func (s *Service) updateArticleVersion(c context.Context, a *model.Article, actID int64) (err error) {
	var tx *sql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("updateArticleVersion.tx.BeginTran() error(%+v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("updateArticleVersion.tx.Rollback() error(%+v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("updateArticleVersion.tx.Commit() error(%+v)", err)
			return
		}
	}()
	if err = s.dao.TxUpdateArticleStateApplyTime(c, tx, a.ID, a.State); err != nil {
		dao.PromError("creation:修改文章状态")
		return
	}

	if x, e := s.ArticleVersion(c, a.ID); e == nil && x.ID != 0 {
		if err = s.dao.TxUpdateArticleVersion(c, tx, a.ID, a, actID); err != nil {
			dao.PromError("creation:更新版本")
		}
	} else {
		if err = s.dao.TxAddArticleVersion(c, tx, a.ID, a, actID); err != nil {
			dao.PromError("creation:创建版本")
		}
	}
	// if err = s.dao.TxUpdateArticleVersion(c, tx, a.ID, a, actID); err != nil {
	// 	dao.PromError("creation:更新版本")
	// }
	return
}

func (s *Service) updateArticleDB(c context.Context, a *model.Article) (err error) {
	var tx *sql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("tx.BeginTran() error(%+v)", err)
		return
	}
	if err = s.dao.TxUpdateArticleMeta(c, tx, a.Meta); err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			dao.PromError("creation:更新文章meta")
			log.Error("tx.Rollback() error(%+v)", err1)
		}
		return
	}
	keywords, _ := s.Segment(c, int32(a.ID), a.Content, 1, "article")
	if err = s.dao.TxUpdateArticleContent(c, tx, a.ID, a.Content, keywords); err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			dao.PromError("creation:更新文章content")
			log.Error("tx.Rollback() error(%+v)", err1)
		}
		return
	}
	// add version
	if err = s.dao.TxUpdateArticleVersion(c, tx, a.ID, a, 0); err != nil {
		dao.PromError("creation:添加历史记录")
		return
	}
	if err = tx.Commit(); err != nil {
		dao.PromError("creation:更新文章")
		log.Error("tx.Commit() error(%+v)", err)
	}
	return
}

// DelArticle drops article.
func (s *Service) DelArticle(c context.Context, aid, mid int64) (err error) {
	if aid <= 0 {
		err = ecode.ArtCreationIDErr
		return
	}
	if err = s.checkPrivilege(c, mid); err != nil {
		return
	}
	var a *model.Meta
	if a, err = s.checkArtAuthor(c, aid, mid); err != nil {
		return
	}
	// can not delelte article which state is pending.
	if a.State == model.StatePending || a.State == model.StateOpenPending {
		err = ecode.ArtCreationDelPendingErr
		return
	}
	lists, _ := s.dao.RawArtsListID(c, []int64{aid})
	if err = s.delArticle(c, aid); err != nil {
		return
	}
	s.dao.DelActivity(c, aid, "")
	cache.Save(func() {
		c := context.TODO()
		s.dao.DelRecommend(c, aid)
		s.DelRecommendArtCache(c, aid, a.Category.ID)
		s.DelArticleCache(c, mid, aid)
		s.deleteArtsListCache(c, aid)
		if lists[aid] > 0 {
			s.updateListInfo(c, lists[aid])
			s.RebuildListCache(c, lists[aid])
		}
	})
	return
}

func (s *Service) delArticle(c context.Context, aid int64) (err error) {
	var tx *sql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("tx.BeginTran() error(%+v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%+v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			dao.PromError("creation:删除文章")
			log.Error("tx.Commit() error(%+v)", err)
		}
	}()
	if err = s.dao.TxDeleteArticleMeta(c, tx, aid); err != nil {
		dao.PromError("creation:删除文章meta")
		return
	}
	if err = s.dao.TxDeleteArticleContent(c, tx, aid); err != nil {
		dao.PromError("creation:delete文章content")
		return
	}
	if err = s.dao.TxDelFilteredArtMeta(c, tx, aid); err != nil {
		dao.PromError("creation:删除过滤文章meta")
		return
	}
	if err = s.dao.TxDelFilteredArtContent(c, tx, aid); err != nil {
		dao.PromError("creation:删除过滤文章content")
		return
	}
	if err = s.dao.TxDelArticleList(tx, aid); err != nil {
		return
	}
	if err = s.dao.TxDelArticleVersion(c, tx, aid); err != nil {
		return
	}
	return
}

// CreationUpperArticlesMeta gets article list by mid.
func (s *Service) CreationUpperArticlesMeta(c context.Context, mid int64, group, category, sortType, pn, ps int, ip string) (res *model.CreationArts, err error) {
	res = &model.CreationArts{}
	if err = s.checkPrivilege(c, mid); err != nil {
		return
	}
	var (
		aids    []int64
		ams, as []*model.Meta
		stats   = make(map[int64]*model.Stats)
		start   = (pn - 1) * ps
		end     = start + ps - 1
		total   int
	)
	res.Page = &model.ArtPage{
		Pn: pn,
		Ps: ps,
	}
	eg, errCtx := errgroup.WithContext(c)
	eg.Go(func() (err error) {
		res.Type, _ = s.dao.UpperArticlesTypeCount(errCtx, mid)
		return
	})
	eg.Go(func() (err error) {
		if ams, err = s.dao.UpperArticlesMeta(errCtx, mid, group, category); err != nil {
			return
		}
		total = len(ams)
		res.Page.Total = total
		return
	})
	if err = eg.Wait(); err != nil {
		log.Error("eg.Wait() error(%+v)", err)
		return
	}
	if total == 0 || total < start {
		return
	}
	for _, v := range ams {
		aids = append(aids, v.ID)
	}
	if stats, err = s.stats(c, aids); err != nil {
		dao.PromError("creation:获取计数信息")
		log.Error("s.stats(%v) err: %+v", aids, err)
		err = nil
	}
	for _, v := range ams {
		// stats
		if st := stats[v.ID]; st != nil {
			v.Stats = st
		} else {
			v.Stats = new(model.Stats)
		}
	}
	switch sortType {
	case _sortDefault, _sortByCtime:
		sort.Slice(ams, func(i, j int) bool { return ams[i].Ctime > ams[j].Ctime })
	case _sortByLike:
		sort.Slice(ams, func(i, j int) bool { return ams[i].Stats.Like > ams[j].Stats.Like })
	case _sortByReply:
		sort.Slice(ams, func(i, j int) bool { return ams[i].Stats.Reply > ams[j].Stats.Reply })
	case _sortByView:
		sort.Slice(ams, func(i, j int) bool { return ams[i].Stats.View > ams[j].Stats.View })
	case _sortByFavorite:
		sort.Slice(ams, func(i, j int) bool { return ams[i].Stats.Favorite > ams[j].Stats.Favorite })
	case _sortByCoin:
		sort.Slice(ams, func(i, j int) bool { return ams[i].Stats.Coin > ams[j].Stats.Coin })
	}
	if total > end {
		as = ams[start : end+1]
	} else {
		as = ams[start:]
	}
	for _, v := range as {
		var pid int64
		if pid, err = s.CategoryToRoot(v.Category.ID); err != nil {
			dao.PromError("creation:获取根分区")
			log.Error("s.CategoryToRoot(%d) error(%+v)", v.Category.ID, err)
			continue
		}
		v.Category = s.categoriesMap[pid]
		v.List, _ = s.dao.ArtList(c, v.ID)
	}
	res.Articles = as
	return
}

// convertState converts -1 to 1, -2 to 2 if the article had been published.
func (s *Service) convertState(c context.Context, a *model.Article) (err error) {
	var am *model.Meta
	if am, err = s.dao.CreationArticleMeta(c, a.ID); err != nil {
		return
	}
	if am.State == model.StateOpen || am.State == model.StateRePass || am.State == model.StateReReject {
		if s.EditTimes(c, a.ID) <= 0 {
			err = ecode.ArtUpdateFullErr
			return
		}
		a.State = model.StateRePending
		return
	}
	if (am.State != model.StateReject) && (am.State != model.StateOpenReject) {
		err = ecode.ArtCannotEditErr
		return
	}
	if am.State > 0 && a.State < 0 {
		a.State = -a.State
	}
	return
}

// CreationWithdrawArticle recall the article and add it to draft.
func (s *Service) CreationWithdrawArticle(c context.Context, mid, aid int64) (err error) {
	if err = s.checkPrivilege(c, mid); err != nil {
		return
	}
	var am *model.Meta
	if am, err = s.checkArtAuthor(c, aid, mid); err != nil {
		return
	}
	// 只有初次提交待审的文章才可以撤回，发布后的修改待审不允许撤回
	if am.State != model.StatePending {
		err = ecode.ArtCreationStateErr
		return
	}
	var content string
	if content, err = s.dao.CreationArticleContent(c, aid); err != nil {
		return
	}
	// add draft
	var tags []string
	if ts, e := s.Tags(c, aid, false); e == nil && len(ts) > 0 {
		for _, v := range ts {
			tags = append(tags, v.Name)
		}
	}
	am.ID = 0
	draft := &model.Draft{Article: &model.Article{Meta: am, Content: content}, Tags: tags}
	if _, err = s.AddArtDraft(c, draft); err != nil {
		return
	}
	// delete article
	err = s.delArticle(c, aid)
	// err = s.dao.UpdateArticleState(c, aid, model.StateOpen)
	return
}

// UpStat up stat
func (s *Service) UpStat(c context.Context, mid int64) (res model.UpStat, err error) {
	return s.dao.UpStat(c, mid)
}

// UpThirtyDayStat for 30 days stat.
func (s *Service) UpThirtyDayStat(c context.Context, mid int64) (res []*model.ThirtyDayArticle, err error) {
	res, err = s.dao.ThirtyDayArticle(c, mid)
	return
}

// ArticleUpCover article upload cover.
func (s *Service) ArticleUpCover(c context.Context, fileType string, body []byte) (url string, err error) {
	if len(body) == 0 {
		err = ecode.FileNotExists
		return
	}
	if len(body) > s.c.BFS.MaxFileSize {
		err = ecode.FileTooLarge
		return
	}
	url, err = s.dao.UploadImage(c, fileType, body)
	if err != nil {
		log.Error("creation: s.bfs.Upload error(%v)", err)
	}
	return
}

// ArticleVersion .
func (s *Service) ArticleVersion(c context.Context, aid int64) (a *model.Article, err error) {
	if a, err = s.dao.ArticleVersion(c, aid); err != nil {
		return
	}
	a.Category = s.categoriesMap[a.Category.ID]
	return
}

// EditTimes .
func (s *Service) EditTimes(c context.Context, id int64) (res int) {
	var (
		et    = _editTimes
		count = et
		err   error
	)
	if s.c.Article.EditTimes != 0 {
		et = s.c.Article.EditTimes
	}
	if count, err = s.dao.EditTimes(c, id); err != nil {
		return
	}
	res = et - count
	if res < 0 {
		res = 0
	}
	return
}

func (s *Service) lastReason(c context.Context, id int64, state int32) (res string, err error) {
	return s.dao.LastReason(c, id, state)
}
