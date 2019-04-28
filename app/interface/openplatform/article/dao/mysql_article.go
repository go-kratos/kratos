package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	artmdl "go-common/app/interface/openplatform/article/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"

	"go-common/library/sync/errgroup"
)

// Article gets article's meta and content.
func (d *Dao) Article(c context.Context, aid int64) (res *artmdl.Article, err error) {
	res = &artmdl.Article{}
	if res.Meta, err = d.ArticleMeta(c, aid); err != nil {
		PromError("article:获取文章meta")
		return
	}
	if res.Meta == nil {
		res = nil
		return
	}
	if res.Content, err = d.ArticleContent(c, aid); err != nil {
		PromError("article:获取文章content")
	}
	if res.Keywords, err = d.ArticleKeywords(c, aid); err != nil {
		PromError("article:获取文章keywords")
	}
	res.Strong()
	return
}

// ArticleContent get article content
func (d *Dao) ArticleContent(c context.Context, id int64) (res string, err error) {
	contentSQL := fmt.Sprintf(_articleContentSQL, d.hit(id))
	if err = d.articleDB.QueryRow(c, contentSQL, id).Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		PromError("db:ArticleContent")
		log.Error("dao.ArticleContent(%s) error(%+v)", contentSQL, err)
	}
	return
}

// ArticleKeywords get article keywords
func (d *Dao) ArticleKeywords(c context.Context, id int64) (res string, err error) {
	keywordsSQL := fmt.Sprintf(_articleKeywordsSQL, d.hit(id))
	if err = d.articleDB.QueryRow(c, keywordsSQL, id).Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		PromError("db:ArticleKeywords")
		log.Error("dao.ArticleKeywords(%s) error(%+v)", keywordsSQL, err)
	}
	res = strings.Replace(res, "\001", ",", -1)
	return
}

// ArticleMeta get article meta
func (d *Dao) ArticleMeta(c context.Context, id int64) (res *artmdl.Meta, err error) {
	var (
		row                        *xsql.Row
		imageURLs, originImageURLs string
		category                   = &artmdl.Category{}
		author                     = &artmdl.Author{}
		t                          int64
		ct                         time.Time
	)
	res = &artmdl.Meta{Media: &artmdl.Media{}}
	row = d.articleMetaStmt.QueryRow(c, id)
	if err = row.Scan(&res.ID, &category.ID, &res.Title, &res.Summary, &res.BannerURL, &res.TemplateID, &res.State, &author.Mid, &res.Reprint, &imageURLs, &t, &ct, &res.Attributes, &res.Words, &res.Dynamic, &originImageURLs, &res.Media.MediaID, &res.Media.Spoiler); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		PromError("db:ArticleMeta")
		log.Error("dao.ArticleMeta.Scan error(%+v)", err)
		return
	}
	res.PublishTime = xtime.Time(t)
	res.Category = category
	res.Author = author
	res.Ctime = xtime.Time(ct.Unix())
	res.ImageURLs = strings.Split(imageURLs, ",")
	res.OriginImageURLs = strings.Split(originImageURLs, ",")
	res.BannerURL = artmdl.CompleteURL(res.BannerURL)
	res.ImageURLs = artmdl.CompleteURLs(res.ImageURLs)
	res.OriginImageURLs = artmdl.CompleteURLs(res.OriginImageURLs)
	res.Strong()
	return
}

// ArticleMetas get article metats
func (d *Dao) ArticleMetas(c context.Context, aids []int64) (res map[int64]*artmdl.Meta, err error) {
	var (
		group, errCtx = errgroup.WithContext(c)
		mutex         = &sync.Mutex{}
	)
	if len(aids) == 0 {
		return
	}
	res = make(map[int64]*artmdl.Meta)
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
			metasSQL := fmt.Sprintf(_articlesMetaSQL, xstr.JoinInts(keys))
			if rows, err = d.articleDB.Query(errCtx, metasSQL); err != nil {
				PromError("db:ArticleMetas")
				return
			}
			defer rows.Close()
			for rows.Next() {
				var (
					imageURLs, originImageURLs string
					t                          int64
					ct                         time.Time
					a                          = &artmdl.Meta{Category: &artmdl.Category{}, Author: &artmdl.Author{}, Media: &artmdl.Media{}}
				)
				err = rows.Scan(&a.ID, &a.Category.ID, &a.Title, &a.Summary, &a.BannerURL, &a.TemplateID, &a.State, &a.Author.Mid, &a.Reprint, &imageURLs, &t, &ct, &a.Attributes, &a.Words, &a.Dynamic, &originImageURLs, &a.Media.MediaID, &a.Media.Spoiler)
				if err != nil {
					return
				}
				a.ImageURLs = strings.Split(imageURLs, ",")
				a.OriginImageURLs = strings.Split(originImageURLs, ",")
				a.PublishTime = xtime.Time(t)
				a.Ctime = xtime.Time(ct.Unix())
				a.BannerURL = artmdl.CompleteURL(a.BannerURL)
				a.ImageURLs = artmdl.CompleteURLs(a.ImageURLs)
				a.OriginImageURLs = artmdl.CompleteURLs(a.OriginImageURLs)
				a.Strong()
				mutex.Lock()
				res[a.ID] = a
				mutex.Unlock()
			}
			err = rows.Err()
			return err
		})
	}
	if err = group.Wait(); err != nil {
		PromError("db:ArticleMetas")
		log.Error("dao.ArticleMetas error(%+v)", err)
		return
	}
	if len(res) == 0 {
		res = nil
	}
	return
}

// AllArticleMeta 所有状态/删除 的文章
func (d *Dao) AllArticleMeta(c context.Context, id int64) (res *artmdl.Meta, err error) {
	var (
		row                        *xsql.Row
		imageURLs, originImageURLs string
		category                   = &artmdl.Category{}
		author                     = &artmdl.Author{}
		t                          int64
		ct                         time.Time
	)
	res = &artmdl.Meta{Media: &artmdl.Media{}}
	row = d.allArticleMetaStmt.QueryRow(c, id)
	if err = row.Scan(&res.ID, &category.ID, &res.Title, &res.Summary, &res.BannerURL, &res.TemplateID, &res.State, &author.Mid, &res.Reprint, &imageURLs, &t, &ct, &res.Attributes, &res.Words, &res.Dynamic, &originImageURLs, &res.Media.MediaID, &res.Media.Spoiler); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		PromError("db:AllArticleMeta")
		log.Error("row.AllArticleMeta.Scan error(%+v)", err)
		return
	}
	res.PublishTime = xtime.Time(t)
	res.Category = category
	res.Author = author
	res.Ctime = xtime.Time(ct.Unix())
	res.ImageURLs = strings.Split(imageURLs, ",")
	res.OriginImageURLs = strings.Split(originImageURLs, ",")
	res.BannerURL = artmdl.CompleteURL(res.BannerURL)
	res.ImageURLs = artmdl.CompleteURLs(res.ImageURLs)
	res.OriginImageURLs = artmdl.CompleteURLs(res.OriginImageURLs)
	res.Strong()
	return
}

// UpperArticleCount get upper article count
func (d *Dao) UpperArticleCount(c context.Context, mid int64) (res int, err error) {
	row := d.articleUpperCountStmt.QueryRow(c, mid)
	if err = row.Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		PromError("db:UpperArticleCount")
		log.Error("dao.UpperArticleCount error(%+v)", err)
	}
	return
}

// ArticleRemainCount returns the number that user could be use to posting new articles.
func (d *Dao) ArticleRemainCount(c context.Context, mid int64) (count int, err error) {
	beginTime := time.Now().Format("2006-01-02") + " 00:00:00"
	if err = d.articleUpCntTodayStmt.QueryRow(c, mid, beginTime).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		PromError("db:ArticleRemainCount")
		log.Error("dao.ArticleRemainCount(%d,%s) error(%+v)", mid, beginTime, err)
	}
	return
}

// TagArticles .
func (d *Dao) TagArticles(c context.Context, tags []int64) (aids []int64, err error) {
	var (
		rows  *xsql.Rows
		query string
		tmps  = make(map[int64]bool)
	)
	query = fmt.Sprintf(_tagArticlesSQL, xstr.JoinInts(tags))
	if rows, err = d.articleDB.Query(c, query); err != nil {
		PromError("dao:TagArticles")
		log.Error("dao.TagArticles(%s) error(%+v)", query, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			tid     int64
			oid     string
			logDate xtime.Time
			ts      []int64
			aid     int64
			now     = time.Now()
		)
		rows.Scan(&tid, &oid, &logDate)
		if now.Sub(logDate.Time()) > time.Hour*60 {
			continue
		}
		ids := strings.Split(oid, "，")
		for _, id := range ids {
			if aid, err = strconv.ParseInt(id, 10, 64); err != nil {
				log.Error("dao.TagArticles.ParseInt(%s) error(%+v)", id, err)
				return
			}
			if !tmps[aid] {
				aids = append(aids, aid)
				tmps[aid] = true
			}
			ts = append(ts, aid)
		}
		d.AddCacheAidsByTag(c, tid, &artmdl.TagArts{Tid: tid, Aids: ts})
	}
	return
}

// MediaArticle .
func (d *Dao) MediaArticle(c context.Context, mediaID int64, mid int64) (id int64, err error) {
	var rows *xsql.Rows
	if rows, err = d.articleDB.Query(c, _mediaArticleSQL, mid, mediaID); err != nil {
		log.Error("dao.MediaArticle.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			log.Error("dao.MediaArticle.Scan error(%v)", err)
			return
		}
		if id > 0 {
			return
		}
	}
	return
}

// MediaIDByID .
func (d *Dao) MediaIDByID(c context.Context, aid int64) (id int64, err error) {
	row := d.articleDB.QueryRow(c, _mediaByIDSQL, aid)
	if err = row.Scan(&id); err != nil {
		log.Error("dao.MediaIDByID.Scan error(%v)", err)
	}
	return
}
