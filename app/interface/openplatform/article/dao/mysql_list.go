package dao

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"sync"
	"time"

	"go-common/app/interface/openplatform/article/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

const (
	_creativeCountArticlesSQL     = "SELECT count(*) FROM articles WHERE mid = ? and deleted_time = 0 AND category_id in (%s)"
	_creativeListsSQL             = "SELECT id, image_url, name, update_time, ctime, summary, publish_time, words FROM lists WHERE deleted_time = 0 AND mid = ?"
	_creativeListAddSQL           = "INSERT INTO lists(name, image_url, summary, publish_time, words, mid) VALUES(?,?,?,?,?,?)"
	_creativeListDelSQL           = "update lists SET deleted_time = ? WHERE id = ? AND deleted_time = 0"
	_creativeListUpdateSQL        = "update lists SET name = ?, image_url = ?, summary = ?, words = ?, publish_time = ? where id = ? and deleted_time = 0"
	_creativeListUpdateTimeSQL    = "update lists SET update_time = ? where id = ? and update_time < ? and deleted_time = 0"
	_creativeCategoryArticlesSQL  = "SELECT id, title, publish_time, state FROM articles WHERE mid = ? AND deleted_time = 0"
	_creativeListArticlesSQL      = "SELECT article_id, position FROM article_lists WHERE list_id = ? AND deleted_time = 0 ORDER BY position ASC"
	_creativeListsArticlesSQL     = "SELECT article_id, position, list_id FROM article_lists WHERE list_id in (%s) AND deleted_time = 0"
	_creativeListAddArticleSQL    = "INSERT INTO article_lists(article_id, list_id, position) values(?,?,?) ON DUPLICATE KEY UPDATE deleted_time =0, position=?"
	_creativeListDelArticleSQL    = "UPDATE article_lists SET deleted_time = ? WHERE article_id = ? and list_id = ? and deleted_time = 0"
	_creativeDelArticleListSQL    = "UPDATE article_lists SET deleted_time = ? WHERE article_id = ? and deleted_time = 0"
	_creativeListDelAllArticleSQL = "UPDATE article_lists SET deleted_time = ? WHERE list_id = ? and deleted_time = 0"
	_creativeArticlesSQL          = "SELECT id, title, state, publish_time FROM articles WHERE id in (%s) and deleted_time = 0"
	_listSQL                      = "SELECT id, mid, image_url, name, update_time, ctime, summary, words, publish_time FROM lists WHERE deleted_time = 0 AND id = ?"
	_listsSQL                     = "SELECT id, mid, image_url, name, update_time, ctime, summary, words, publish_time FROM lists WHERE deleted_time = 0 AND id in (%s)"
	_allListsSQL                  = "SELECT id, mid, image_url, name, update_time, ctime, summary, words, publish_time FROM lists WHERE deleted_time = 0"
	_artslistSQL                  = "SELECT article_id, list_id FROM article_lists WHERE article_id IN (%s) AND deleted_time = 0"
	_allListsExSQL                = "SELECT id, mid, image_url, name, update_time, ctime, summary, words, publish_time FROM lists WHERE deleted_time = 0 ORDER BY id DESC LIMIT ?, ?"
)

// CreativeUpLists get article lists
func (d *Dao) CreativeUpLists(c context.Context, mid int64) (res []*model.List, err error) {
	rows, err := d.creativeListsStmt.Query(c, mid)
	if err != nil {
		PromError("db:up主文集")
		log.Errorv(c, log.KV("log", "CreativeUplists"), log.KV("error", err))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			t, ctime time.Time
			r        = &model.List{Mid: mid}
			pt       int64
		)
		if err = rows.Scan(&r.ID, &r.ImageURL, &r.Name, &t, &ctime, &r.Summary, &pt, &r.Words); err != nil {
			PromError("db:up主文集scan")
			log.Error("dao.CreativeUpLists.rows.Scan error(%+v)", err)
			return
		}
		if t.Unix() > 0 {
			r.UpdateTime = xtime.Time(t.Unix())
		}
		r.Ctime = xtime.Time(ctime.Unix())
		r.PublishTime = xtime.Time(pt)
		res = append(res, r)
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// RawUpLists .
func (d *Dao) RawUpLists(c context.Context, mid int64) (res []int64, err error) {
	lists, err := d.CreativeUpLists(c, mid)
	for _, list := range lists {
		res = append(res, list.ID)
	}
	return
}

// CreativeListUpdate update list
func (d *Dao) CreativeListUpdate(c context.Context, id int64, name, imageURL, summary string, publishTime xtime.Time, words int64) (err error) {
	if _, err := d.creativeListUpdateStmt.Exec(c, name, imageURL, summary, words, int64(publishTime), id); err != nil {
		PromError("db:修改文集")
		log.Errorv(c, log.KV("dao.CreativeListUpdate.Exec", id), log.KV("name", name), log.KV("image_url", imageURL), log.KV("error", err), log.KV("summary", summary))
	}
	return
}

// CreativeListDelAllArticles del list
func (d *Dao) CreativeListDelAllArticles(c context.Context, id int64) (err error) {
	if _, err := d.creativeListDelAllArticleStmt.Exec(c, time.Now(), id); err != nil {
		PromError("db:删除文集下的文章")
		log.Errorv(c, log.KV("dao.CreativeListDelAllArticles.Exec", id), log.KV("error", err))
	}
	return
}

// CreativeListDel del list
func (d *Dao) CreativeListDel(c context.Context, id int64) (err error) {
	if _, err := d.creativeListDelStmt.Exec(c, time.Now().Unix(), id); err != nil {
		PromError("db:删除文集")
		log.Errorv(c, log.KV("dao.CreativeListDel.Exec", id), log.KV("error", err))
	}
	return
}

// CreativeListUpdateTime update list time
func (d *Dao) CreativeListUpdateTime(c context.Context, id int64, t time.Time) (err error) {
	if _, err := d.creativeListUpdateTimeStmt.Exec(c, t, id, t); err != nil {
		PromError("db:修改文集更新时间")
		log.Errorv(c, log.KV("dao.CreativeListUpdateTime.Exec", id), log.KV("time", t), log.KV("error", err))
	}
	return
}

// CreativeListAdd add list
func (d *Dao) CreativeListAdd(c context.Context, mid int64, name, imageURL, summary string, publishTime xtime.Time, words int64) (res int64, err error) {
	r, err := d.creativeListAddStmt.Exec(c, name, imageURL, summary, int64(publishTime), words, mid)
	if err != nil {
		PromError("db:增加文集")
		log.Errorv(c, log.KV("dao.CreativeListAdd.Exec", mid), log.KV("name", name), log.KV("image_url", imageURL), log.KV("error", err), log.KV(summary, summary))
		return
	}
	if res, err = r.LastInsertId(); err != nil {
		PromError("db:增加文集ID")
		log.Errorv(c, log.KV("log", "res.LastInsertId"), log.KV("error", err))
	}
	return
}

// CreativeCountArticles novel count
func (d *Dao) CreativeCountArticles(c context.Context, mid int64, cids []int64) (res int64, err error) {
	s := fmt.Sprintf(_creativeCountArticlesSQL, xstr.JoinInts(cids))
	if err = d.articleDB.QueryRow(c, s, mid).Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		PromError("db:小说计数")
		log.Errorv(c, log.KV("log", "dao.CreativeCountArticles"), log.KV("error", err))
	}
	return
}

// CreativeCategoryArticles can add articles
func (d *Dao) CreativeCategoryArticles(c context.Context, mid int64) (res []*model.ListArtMeta, err error) {
	rows, err := d.articleDB.Query(c, _creativeCategoryArticlesSQL, mid)
	if err != nil {
		PromError("db:up主可被加入文集的文章列表")
		log.Errorv(c, log.KV("log", "CreativeCategoryArticles"), log.KV("error", err))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			t int64
			m = &model.ListArtMeta{}
		)
		if err = rows.Scan(&m.ID, &m.Title, &t, &m.State); err != nil {
			PromError("db:up主可被加入文集的文章列表scan")
			log.Errorv(c, log.KV("log", "CreativeCategoryArticles"), log.KV("error", err))
			return
		}
		m.PublishTime = xtime.Time(t)
		res = append(res, m)
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// CreativeListArticles .
func (d *Dao) CreativeListArticles(c context.Context, listID int64) (res []*model.ListArtMeta, err error) {
	rows, err := d.creativeListArticlesStmt.Query(c, listID)
	if err != nil {
		PromError("db:文集的文章列表")
		log.Errorv(c, log.KV("log", "CreativeListArticles"), log.KV("error", err))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			r = &model.ListArtMeta{}
		)
		if err = rows.Scan(&r.ID, &r.Position); err != nil {
			PromError("db:文集的文章列表scan")
			log.Errorv(c, log.KV("log", "CreativeListArticles"), log.KV("error", err))
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// CreativeListsArticles .
func (d *Dao) CreativeListsArticles(c context.Context, listIDs []int64) (res map[int64][]*model.ListArtMeta, err error) {
	if len(listIDs) == 0 {
		return
	}
	s := fmt.Sprintf(_creativeListsArticlesSQL, xstr.JoinInts(listIDs))
	rows, err := d.articleDB.Query(c, s)
	if err != nil {
		PromError("db:多个文集的文章列表")
		log.Errorv(c, log.KV("log", "CreativeListsArticles"), log.KV("error", err))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			r   = &model.ListArtMeta{}
			lid int64
		)
		if err = rows.Scan(&r.ID, &r.Position, &lid); err != nil {
			PromError("db:多个文集的文章列表scan")
			log.Errorv(c, log.KV("log", "CreativeListsArticles"), log.KV("error", err))
			return
		}
		if res == nil {
			res = make(map[int64][]*model.ListArtMeta)
		}
		res[lid] = append(res[lid], r)
	}
	for lid, arts := range res {
		sort.Slice(arts, func(i, j int) bool { return arts[i].Position < arts[j].Position })
		res[lid] = arts
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// CreativeArticles get up all state articles
func (d *Dao) CreativeArticles(c context.Context, aids []int64) (res map[int64]*model.ListArtMeta, err error) {
	var (
		group, errCtx = errgroup.WithContext(c)
		mutex         = &sync.Mutex{}
	)
	if len(aids) == 0 {
		return
	}
	res = make(map[int64]*model.ListArtMeta)
	keysLen := len(aids)
	for i := 0; i < keysLen; i += _mysqlBulkSize {
		var keys []int64
		if (i + _mysqlBulkSize) > keysLen {
			keys = aids[i:]
		} else {
			keys = aids[i : i+_mysqlBulkSize]
		}
		group.Go(func() (err error) {
			var rows *xsql.Rows
			metasSQL := fmt.Sprintf(_creativeArticlesSQL, xstr.JoinInts(keys))
			if rows, err = d.articleDB.Query(errCtx, metasSQL); err != nil {
				PromError("db:CreativeArticles")
				log.Errorv(c, log.KV("log", "CreativeArticles"), log.KV("error", err))
				return
			}
			defer rows.Close()
			for rows.Next() {
				var (
					t int64
					a = &model.ListArtMeta{}
				)
				err = rows.Scan(&a.ID, &a.Title, &a.State, &t)
				if err != nil {
					return
				}
				a.PublishTime = xtime.Time(t)
				mutex.Lock()
				res[a.ID] = a
				mutex.Unlock()
			}
			err = rows.Err()
			return err
		})
	}
	if err = group.Wait(); err != nil {
		PromError("db:CreativeArticles")
		log.Errorv(c, log.KV("error", err))
		return
	}
	if len(res) == 0 {
		res = nil
	}
	return
}

// RawList get list from db
func (d *Dao) RawList(c context.Context, id int64) (res *model.List, err error) {
	var (
		t, ctime time.Time
		pt       int64
	)
	res = &model.List{ID: id}
	if err = d.listStmt.QueryRow(c, id).Scan(&res.ID, &res.Mid, &res.ImageURL, &res.Name, &t, &ctime, &res.Summary, &res.Words, &pt); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		PromError("db文集")
		log.Errorv(c, log.KV("log", err))
	}
	if t.Unix() > 0 {
		res.UpdateTime = xtime.Time(t.Unix())
	}
	res.Ctime = xtime.Time(ctime.Unix())
	res.PublishTime = xtime.Time(pt)
	return
}

// TxAddListArticle tx add list article
func (d *Dao) TxAddListArticle(c context.Context, tx *xsql.Tx, listID int64, aid int64, position int) (err error) {
	if _, err = tx.Exec(_creativeListAddArticleSQL, aid, listID, position, position); err != nil {
		PromError("db:新增文集文章tx")
		log.Error("tx.Exec() error(%+v)", err)
		return
	}
	return
}

// TxDelListArticle tx del list article
func (d *Dao) TxDelListArticle(c context.Context, tx *xsql.Tx, listID int64, aid int64) (err error) {
	t := time.Now().Unix()
	if _, err = tx.Exec(_creativeListDelArticleSQL, t, aid, listID); err != nil {
		PromError("db:删除文集文章tx")
		log.Error("tx.Exec() error(%+v)", err)
		return
	}
	return
}

// TxDelArticleList .
func (d *Dao) TxDelArticleList(tx *xsql.Tx, aid int64) (err error) {
	t := time.Now().Unix()
	if _, err = tx.Exec(_creativeDelArticleListSQL, t, aid); err != nil {
		PromError("db:tx删除文集文章")
		log.Error("tx.Exec() error(%+v)", err)
		return
	}
	return
}

// RawLists get lists from db
func (d *Dao) RawLists(c context.Context, ids []int64) (res map[int64]*model.List, err error) {
	s := fmt.Sprintf(_listsSQL, xstr.JoinInts(ids))
	rows, err := d.articleDB.Query(c, s)
	if err != nil {
		PromError("db:文集列表")
		log.Errorv(c, log.KV("log", "Lists"), log.KV("error", err))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			t, ctime time.Time
			r        = &model.List{}
			pt       int64
		)
		if err = rows.Scan(&r.ID, &r.Mid, &r.ImageURL, &r.Name, &t, &ctime, &r.Summary, &r.Words, &pt); err != nil {
			PromError("db:文集列表scan")
			log.Error("dao.Lists.rows.Scan error(%+v)", err)
			return
		}
		if t.Unix() > 0 {
			r.UpdateTime = xtime.Time(t.Unix())
		}
		r.Ctime = xtime.Time(ctime.Unix())
		r.PublishTime = xtime.Time(pt)
		if res == nil {
			res = make(map[int64]*model.List)
		}
		res[r.ID] = r
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// RawAllLists get lists from db
func (d *Dao) RawAllLists(c context.Context) (res []*model.List, err error) {
	rows, err := d.allListStmt.Query(c)
	if err != nil {
		PromError("db:全部文集列表")
		log.Errorv(c, log.KV("log", "AllLists"), log.KV("error", err))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			t, ctime time.Time
			r        = &model.List{}
			pt       int64
		)
		if err = rows.Scan(&r.ID, &r.Mid, &r.ImageURL, &r.Name, &t, &ctime, &r.Summary, &r.Words, &pt); err != nil {
			PromError("db:全部文集列表scan")
			log.Error("dao.AllLists.rows.Scan error(%+v)", err)
			return
		}
		if t.Unix() > 0 {
			r.UpdateTime = xtime.Time(t.Unix())
		}
		r.Ctime = xtime.Time(ctime.Unix())
		r.PublishTime = xtime.Time(pt)
		res = append(res, r)
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// RawAllListsEx get lists from db
func (d *Dao) RawAllListsEx(c context.Context, start int, size int) (res []*model.List, err error) {
	rows, err := d.articleDB.Query(c, _allListsExSQL, start, size)
	if err != nil {
		PromError("db:全部文集列表")
		log.Errorv(c, log.KV("log", "AllLists"), log.KV("error", err))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			t, ctime time.Time
			r        = &model.List{}
			pt       int64
		)
		if err = rows.Scan(&r.ID, &r.Mid, &r.ImageURL, &r.Name, &t, &ctime, &r.Summary, &r.Words, &pt); err != nil {
			PromError("db:全部文集列表scan")
			log.Error("dao.AllLists.rows.Scan error(%+v)", err)
			return
		}
		if t.Unix() > 0 {
			r.UpdateTime = xtime.Time(t.Unix())
		}
		r.Ctime = xtime.Time(ctime.Unix())
		r.PublishTime = xtime.Time(pt)
		res = append(res, r)
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// RawArtsListID get articles list from db
func (d *Dao) RawArtsListID(c context.Context, aids []int64) (res map[int64]int64, err error) {
	s := fmt.Sprintf(_artslistSQL, xstr.JoinInts(aids))
	rows, err := d.articleDB.Query(c, s)
	if err != nil {
		PromError("db:文章所属文集")
		log.Errorv(c, log.KV("log", "ArtsList"), log.KV("error", err))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			aid, listID int64
		)
		if err = rows.Scan(&aid, &listID); err != nil {
			PromError("db:文章所属文集scan")
			log.Error("dao.ArtsList.rows.Scan error(%+v)", err)
			return
		}
		if res == nil {
			res = make(map[int64]int64)
		}
		res[aid] = listID
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// AddListArticle add list article
func (d *Dao) AddListArticle(c context.Context, listID int64, aid int64, position int) (err error) {
	if _, err = d.creativeListAddArticleStmt.Exec(c, aid, listID, position, position); err != nil {
		PromError("db:新增文集文章")
		log.Error("d.creativeListAddArticleStmt(list: %v, aid: %v, position: %v) error(%+v)", listID, aid, position, err)
	}
	return
}

// DelListArticle delete list article
func (d *Dao) DelListArticle(c context.Context, listID int64, aid int64) (err error) {
	if _, err = d.creativeListDelArticleStmt.Exec(c, time.Now().Unix(), aid, listID); err != nil {
		PromError("db:删除文集文章")
		log.Error("d.DelListArticle(list: %v, aid: %v) error(%+v)", listID, aid, err)
	}
	return
}

// RawListArts .
func (d *Dao) RawListArts(c context.Context, listID int64) (res []*model.ListArtMeta, err error) {
	if listID <= 0 {
		return
	}
	arts, err := d.CreativeListArticles(c, listID)
	if err != nil {
		return
	}
	var ids []int64
	for _, art := range arts {
		ids = append(ids, art.ID)
	}
	metas, err := d.ArticleMetas(c, ids)
	if err != nil {
		return
	}
	for _, id := range ids {
		if metas[id] != nil {
			res = append(res, &model.ListArtMeta{
				ID:          id,
				Title:       metas[id].Title,
				PublishTime: metas[id].PublishTime,
				Words:       metas[id].Words,
				ImageURLs:   metas[id].ImageURLs,
				Category:    metas[id].Category,
				Categories:  metas[id].Categories,
				Summary:     metas[id].Summary,
			})
		}
	}
	return
}

// RawListsArts .
func (d *Dao) RawListsArts(c context.Context, ids []int64) (res map[int64][]*model.ListArtMeta, err error) {
	if len(ids) == 0 {
		return
	}
	for _, id := range ids {
		lists, err := d.RawListArts(c, id)
		if err != nil {
			return nil, err
		}
		if res == nil {
			res = make(map[int64][]*model.ListArtMeta)
		}
		res[id] = lists
	}
	return
}

// ArtsList get article's read list
func (d *Dao) ArtsList(c context.Context, aids []int64) (res map[int64]*model.List, err error) {
	if len(aids) == 0 {
		return
	}
	arts, err := d.ArtsListID(c, aids)
	if err != nil {
		return
	}
	listsMap := make(map[int64]bool)
	for _, list := range arts {
		listsMap[list] = true
	}
	var lids []int64
	for l := range listsMap {
		lids = append(lids, l)
	}
	lists, err := d.Lists(c, lids)
	if err != nil {
		return
	}
	res = make(map[int64]*model.List)
	for aid, lid := range arts {
		if lists[lid] != nil {
			res[aid] = lists[lid]
		}
	}
	return
}

// ArtList article list
func (d *Dao) ArtList(c context.Context, aid int64) (res *model.List, err error) {
	if aid <= 0 {
		return
	}
	lists, err := d.ArtsList(c, []int64{aid})
	res = lists[aid]
	return
}

// RawListReadCount .
func (d *Dao) RawListReadCount(c context.Context, id int64) (res int64, err error) {
	metas, err := d.RawListArts(c, id)
	if err != nil {
		return
	}
	var ids []int64
	for _, meta := range metas {
		if meta.IsNormal() {
			ids = append(ids, meta.ID)
		}
	}
	// get stats
	stats, err := d.ArticlesStats(c, ids)
	if err != nil {
		return
	}
	for _, aid := range ids {
		if stats[aid] != nil {
			res += stats[aid].View
		}
	}
	return
}
