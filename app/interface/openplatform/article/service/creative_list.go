package service

import (
	"context"
	"sort"
	"time"

	"go-common/app/interface/openplatform/article/dao"
	"go-common/app/interface/openplatform/article/model"
	filter "go-common/app/service/main/filter/model/rpc"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// const _novel = 16
const _positionStep = 1000

// CreativeUpLists up list
func (s *Service) CreativeUpLists(c context.Context, mid int64) (novel bool, lists []*model.CreativeList, err error) {
	if err = s.checkPrivilege(c, mid); err != nil {
		return
	}
	// var count int64
	// if count, err = s.dao.CreativeCountArticles(c, mid, s.novelCIDs()); err != nil {
	// 	return
	// }
	// novel = count > 0
	novel = true
	ls, err := s.dao.CreativeUpLists(c, mid)
	if err != nil {
		return
	}
	sortLists(ls)
	var ids []int64
	for _, l := range ls {
		ids = append(ids, l.ID)
	}
	arts, err := s.dao.CreativeListsArticles(c, ids)
	if err != nil {
		return
	}
	for _, l := range ls {
		lists = append(lists, &model.CreativeList{List: l, Total: len(arts[l.ID])})
	}
	return
}

func sortLists(lists []*model.List) {
	sort.Slice(lists, func(i, j int) bool {
		it := int64(lists[i].Ctime)
		jt := int64(lists[j].Ctime)
		if int64(lists[i].UpdateTime) > it {
			it = int64(lists[i].UpdateTime)
		}
		if int64(lists[j].UpdateTime) > jt {
			jt = int64(lists[j].UpdateTime)
		}
		return it > jt
	})
}

func (s *Service) filter(c context.Context, content string) (res string, err error) {
	arg := filter.ArgFilter{Area: "article", Message: content}
	var filterRes *filter.FilterRes
	if filterRes, err = s.filterRPC.FilterArea(c, &arg); err != nil {
		dao.PromError("creative:过滤服务")
		log.Errorv(c, log.KV("log", "s.filterRPC.Filter"), log.KV("content", content), log.KV("err", err))
		res = content
		return
	}
	res = filterRes.Result
	return
}

// CreativeAddList add list
func (s *Service) CreativeAddList(c context.Context, mid int64, name string, summary, imageURL string) (res *model.List, err error) {
	if err = s.checkPrivilege(c, mid); err != nil {
		return
	}
	if lists, _ := s.dao.CreativeUpLists(c, mid); len(lists) >= s.c.Article.ListLimit {
		err = ecode.ArtMaxListErr
		return
	}
	var newName, newSummary string
	newName, _ = s.filter(c, name)
	if newName != name {
		log.Infov(c, log.KV("log", "filter list"), log.KV("old", name), log.KV("new", newName))
	}
	if summary != "" {
		newSummary, _ = s.filter(c, summary)
		if newSummary != summary {
			log.Infov(c, log.KV("log", "filter list summary"), log.KV("old", summary), log.KV("new", newSummary))
		}
	}
	name = newName
	summary = newSummary
	var ok bool
	if name, ok = s.checkTitle(name); !ok || name == "" {
		log.Errorv(c, log.KV("log", "CreativeAddList"), log.KV("name", name), log.KV("mid", mid))
		err = ecode.ArtListNameErr
		return
	}
	id, err := s.dao.CreativeListAdd(c, mid, name, imageURL, summary, xtime.Time(0), 0)
	if err != nil {
		return
	}
	res, err = s.dao.RawList(c, id)
	cache.Save(func() {
		s.dao.RebuildUpListsCache(context.TODO(), mid)
	})
	return
}

// CreativeDelList del list
func (s *Service) CreativeDelList(c context.Context, mid int64, id int64) (err error) {
	if err = s.checkPrivilege(c, mid); err != nil {
		return
	}
	if _, err = s.checkList(c, mid, id); err != nil {
		return
	}
	arts, _ := s.dao.CreativeListArticles(c, id)
	err = s.dao.CreativeListDel(c, id)
	if err != nil {
		return
	}
	err = s.dao.CreativeListDelAllArticles(c, id)
	// del list cache/articles cache/article list cache
	cache.Save(func() {
		c := context.TODO()
		var aids []int64
		for _, a := range arts {
			aids = append(aids, a.ID)
		}
		s.deleteArtsListCache(c, aids...)
		s.deleteListArtsCache(c, id)
		s.deleteListCache(c, id)
		s.dao.RebuildUpListsCache(context.TODO(), mid)
	})
	return
}

// func (s *Service) novelCIDs() (cids []int64) {
// 	for _, a := range s.categoriesReverseMap[_novel] {
// 		cids = append(cids, a.ID)
// 	}
// 	return
// }

func (s *Service) creativeNotAddListArticles(c context.Context, mid int64) (res []*model.ListArtMeta, err error) {
	// cids := s.novelCIDs()
	arts, err := s.dao.CreativeCategoryArticles(c, mid)
	if err != nil {
		return
	}
	lists, err := s.dao.CreativeUpLists(c, mid)
	if err != nil {
		return
	}
	var ids []int64
	for _, l := range lists {
		ids = append(ids, l.ID)
	}
	listsArts, err := s.dao.CreativeListsArticles(c, ids)
	if err != nil {
		return
	}
	exists := make(map[int64]bool)
	for _, la := range listsArts {
		for _, a := range la {
			exists[a.ID] = true
		}
	}
	for _, art := range arts {
		if !exists[art.ID] {
			res = append(res, art)
		}
	}
	return
}

// CreativeCanAddArticles can added passed articles
func (s *Service) CreativeCanAddArticles(c context.Context, mid int64) (res []*model.ListArtMeta, err error) {
	if err = s.checkPrivilege(c, mid); err != nil {
		return
	}
	arts, err := s.creativeNotAddListArticles(c, mid)
	if err != nil {
		return
	}
	for _, art := range arts {
		if art.IsNormal() {
			res = append(res, art)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].PublishTime > res[j].PublishTime
	})
	return
}

// CreativeListAllArticles  get read list articles
func (s *Service) CreativeListAllArticles(c context.Context, mid, id int64) (list *model.List, arts []*model.ListArtMeta, err error) {
	if err = s.checkPrivilege(c, mid); err != nil {
		return
	}
	list, err = s.checkList(c, mid, id)
	if err != nil {
		return
	}
	list.Read, _ = s.dao.CacheListReadCount(c, id)
	arts, err = s.rawListArticles(c, id)
	return
}

func (s *Service) checkList(c context.Context, mid, id int64) (list *model.List, err error) {
	if id == 0 {
		return
	}
	list, err = s.dao.RawList(c, id)
	if err != nil {
		return
	}
	if list == nil {
		err = ecode.NothingFound
		return
	}
	if list.Mid != mid {
		err = ecode.ArtCreationMIDErr
		return
	}
	return
}

// CreativeUpdateListArticles update list articles
func (s *Service) CreativeUpdateListArticles(c context.Context, listID int64, name, imageURL, summary string, onlyList bool, mid int64, aids []int64) (list *model.List, err error) {
	if err = s.checkPrivilege(c, mid); err != nil {
		return
	}
	list, err = s.checkList(c, mid, listID)
	if err != nil {
		return
	}
	list, err = s.CreativeUpdateList(c, listID, name, imageURL, summary, list.PublishTime, list.Words)
	if err != nil {
		return
	}
	if onlyList {
		cache.Save(func() {
			s.updateListCache(c, listID)
		})
		return
	}
	if len(aids) > s.c.Article.ListArtsLimit {
		err = ecode.ArtAddListLimitErr
		return
	}
	existArts, err := s.dao.CreativeListArticles(c, listID)
	if err != nil {
		return
	}
	existsMap := make(map[int64]bool)
	for _, a := range existArts {
		existsMap[a.ID] = true
	}
	// 过滤非自己的文章
	metas, err := s.CreativeCanAddArticles(c, mid)
	if err != nil {
		return
	}
	metasMap := make(map[int64]*model.ListArtMeta)
	for _, m := range metas {
		metasMap[m.ID] = m
	}
	var newAids []int64
	for _, aid := range aids {
		if existsMap[aid] || (metasMap[aid] != nil) {
			newAids = append(newAids, aid)
		}
	}
	aids = newAids
	// 计算排序
	updated, deleted := calculateListArtPosition(existArts, aids)
	log.Info("creative: update list update(%v) deleted(%v)", len(updated), len(deleted))
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
			dao.PromError("creative:修改文集")
			log.Error("tx.Commit() error(%+v)", err)
			return
		}
		err = s.updateListInfo(c, listID)
		cache.Save(func() {
			s.dao.CreativeListUpdateTime(context.TODO(), listID, time.Now())
			s.RebuildListCache(context.TODO(), listID)
			s.deleteArtsListCache(context.TODO(), deleted...)
		})
	}()
	for _, id := range deleted {
		if err = s.dao.TxDelListArticle(c, tx, listID, id); err != nil {
			return
		}
	}
	for _, a := range updated {
		if err = s.dao.TxAddListArticle(c, tx, listID, a.ID, a.Position); err != nil {
			return
		}
	}
	return
}

func calculateListArtPosition(src []*model.ListArtMeta, dest []int64) (update []*model.ListArtMeta, delete []int64) {
	srcMap := make(map[int64]*model.ListArtMeta)
	for _, m := range src {
		srcMap[m.ID] = m
	}
	destMap := make(map[int64]bool)
	for _, id := range dest {
		destMap[id] = true
	}
	var newSrc []*model.ListArtMeta
	for _, m := range src {
		if !destMap[m.ID] {
			delete = append(delete, m.ID)
			continue
		}
		newSrc = append(newSrc, m)
	}
	if len(newSrc) == len(dest) {
		equal := true
		for i, d := range dest {
			if newSrc[i].ID != d {
				equal = false
				break
			}
		}
		if equal {
			return
		}
	}
	if len(dest) == 0 {
		return
	}
	for i, id := range dest {
		pos := (i + 1) * _positionStep
		if (srcMap[id] != nil) && (srcMap[id].Position == pos) {
			continue
		}
		update = append(update, &model.ListArtMeta{ID: id, Position: pos})
	}
	return
}

// CreativeUpdateList update list
func (s *Service) CreativeUpdateList(c context.Context, id int64, name, imageURL, summary string, publishTime xtime.Time, words int64) (res *model.List, err error) {
	var newName string
	newName, _ = s.filter(c, name)
	if newName != name {
		log.Infov(c, log.KV("log", "filter title"), log.KV("old", name), log.KV("new", newName))
	}
	name = newName
	var ok bool
	if name, ok = s.checkTitle(name); !ok || name == "" {
		log.Errorv(c, log.KV("log", "CreativeUpdateList"), log.KV("name", name), log.KV("id", id))
		err = ecode.ArtListNameErr
		return
	}
	if summary != "" {
		summary, _ = s.filter(c, summary)
	}
	if err = s.dao.CreativeListUpdate(c, id, name, imageURL, summary, publishTime, words); err != nil {
		return
	}
	res, err = s.dao.RawList(c, id)
	if err != nil {
		return
	}
	cache.Save(func() {
		s.dao.AddCacheList(context.TODO(), res.ID, res)
	})
	return
}

// creativeAddArticleList set article list
func (s *Service) creativeAddArticleList(c context.Context, mid, listID, articleID int64, onlyPass bool) (err error) {
	if listID == 0 {
		return
	}
	log.Infov(c, log.KV("log", "creativeSetArticleList"), log.KV("list_id", listID), log.KV("article_id", articleID), log.KV("mid", mid))
	defer func() {
		if err != nil {
			log.Errorv(c, log.KV("log", "creativeSetArticleList"), log.KV("list_id", listID), log.KV("article_id", articleID), log.KV("mid", mid), log.KV("err", err))
		}
	}()
	if _, err = s.checkList(c, mid, listID); err != nil {
		return
	}
	var can bool
	if can, err = s.checkArticleCanAddList(c, mid, articleID, onlyPass); (err != nil) || !can {
		log.Errorv(c, log.KV("log", "checkArticleCanAddList"), log.KV("mid", mid), log.KV("aid", articleID), log.KV("err", err), log.KV("error", err))
		err = ecode.ArtArtAddListErr
		return
	}
	arts, _ := s.dao.CreativeListArticles(c, listID)
	if len(arts) >= s.c.Article.ListArtsLimit {
		err = ecode.ArtAddListLimitErr
		return
	}
	position := _positionStep
	if len(arts) > 0 {
		position = arts[len(arts)-1].Position + _positionStep
	}
	err = s.dao.AddListArticle(c, listID, articleID, position)
	if err != nil {
		return
	}
	err = s.updateListInfo(c, listID)
	if err != nil {
		return
	}
	cache.Save(func() {
		s.dao.CreativeListUpdateTime(context.TODO(), listID, time.Now())
		// only update passed art
		s.RebuildListCache(context.TODO(), listID)
	})
	return
}

// CreativeUpdateArticleList update article list
func (s *Service) CreativeUpdateArticleList(c context.Context, mid, aid, listID int64, onlyPass bool) (err error) {
	if err = s.checkPrivilege(c, mid); err != nil {
		return
	}
	if listID > 0 {
		_, err = s.checkList(c, mid, listID)
		if err != nil {
			return
		}
	}
	// check article
	meta, err := s.dao.AllArticleMeta(c, aid)
	if err != nil {
		return
	}
	if meta == nil {
		err = ecode.NothingFound
		return
	}
	if meta.Author.Mid != mid {
		err = ecode.ArtCreationMIDErr
		return
	}
	// get article list
	lists, err := s.dao.RawArtsListID(c, []int64{aid})
	if err != nil {
		return
	}
	oldListID := lists[aid]
	if oldListID == listID {
		return
	}
	defer func() {
		if err == nil {
			err = s.dao.CreativeListUpdateTime(c, listID, time.Now())
		}
	}()
	s.deleteArtsListCache(c, aid)
	if oldListID > 0 {
		err = s.dao.DelListArticle(c, oldListID, aid)
		if err != nil {
			return
		}
		err = s.updateListInfo(c, oldListID)
		if err != nil {
			return
		}
		cache.Save(func() {
			s.RebuildListCache(context.TODO(), oldListID)
		})
	}
	if listID > 0 {
		err = s.creativeAddArticleList(c, mid, listID, aid, onlyPass)
	}
	return
}

func (s *Service) checkArticleCanAddList(c context.Context, mid, aid int64, onlyPass bool) (res bool, err error) {
	metas, err := s.creativeNotAddListArticles(c, mid)
	if err != nil {
		log.Errorv(c, log.KV("log", "checkArticleCanAddList"), log.KV("mid", mid), log.KV("aid", aid), log.KV("err", err))
		return
	}
	metasMap := make(map[int64]*model.ListArtMeta)
	for _, m := range metas {
		metasMap[m.ID] = m
	}
	if onlyPass {
		res = (metasMap[aid] != nil) && metasMap[aid].IsNormal()
		return
	}
	res = metasMap[aid] != nil
	return
}

// updateListInfo update list words and publish_time
func (s *Service) updateListInfo(c context.Context, id int64) (err error) {
	list, err := s.dao.RawList(c, id)
	if err != nil {
		return
	}
	if list == nil {
		err = ecode.NothingFound
		return
	}
	metas, err := s.dao.RawListArts(c, id)
	if err != nil {
		return
	}
	var words, pt int64
	for _, meta := range metas {
		if meta.IsNormal() {
			words += meta.Words
			if int64(meta.PublishTime) > pt {
				pt = int64(meta.PublishTime)
			}
		}
	}
	err = s.dao.CreativeListUpdate(c, id, list.Name, list.ImageURL, list.Summary, xtime.Time(pt), words)
	return
}

// RefreshList refresh list info and cache use in api
func (s *Service) RefreshList(c context.Context, id int64) (err error) {
	if err = s.updateListInfo(c, id); err != nil {
		return
	}
	return s.RebuildListCache(c, id)
}
