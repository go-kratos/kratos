package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	artmdl "go-common/app/interface/openplatform/article/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_addArticleDraftSQL    = "REPLACE INTO article_draft_%s (id,category_id,title,summary,banner_url,template_id,mid,reprint,image_urls,tags,content, dynamic_intro, origin_image_urls, list_id, media_id, spoiler) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_articleDraftSQL       = "SELECT id,category_id,title,summary,banner_url,template_id,mid,reprint,image_urls,tags,content,mtime,dynamic_intro,origin_image_urls, list_id, media_id, spoiler FROM article_draft_%s WHERE id=? AND deleted_time=0"
	_checkDraftSQL         = "SELECT deleted_time FROM article_draft_%s WHERE id=?"
	_deleteArticleDraftSQL = "UPDATE article_draft_%s SET deleted_time=? WHERE id=?"
	_upperDraftsSQL        = "SELECT id,category_id,title,summary,template_id,reprint,image_urls,tags,mtime,dynamic_intro,origin_image_urls, list_id FROM article_draft_%s WHERE mid=? AND deleted_time=0 " +
		"ORDER BY mtime DESC LIMIT ?,?"
	_countUpperDraftSQL = "SELECT COUNT(*) FROM article_draft_%s WHERE mid=? AND deleted_time=0"
)

// ArtDraft get draft by article_id
func (d *Dao) ArtDraft(c context.Context, mid, aid int64) (res *artmdl.Draft, err error) {
	var (
		row                        *xsql.Row
		tags                       string
		imageURLs, originImageURLs string
		category                   = &artmdl.Category{}
		author                     = &artmdl.Author{}
		meta                       = &artmdl.Meta{Media: &artmdl.Media{}}
		mtime                      time.Time
		sqlStr                     = fmt.Sprintf(_articleDraftSQL, d.hit(mid))
	)
	res = &artmdl.Draft{Article: &artmdl.Article{}}
	row = d.articleDB.QueryRow(c, sqlStr, aid)
	if err = row.Scan(&meta.ID, &category.ID, &meta.Title, &meta.Summary, &meta.BannerURL, &meta.TemplateID, &author.Mid, &meta.Reprint, &imageURLs, &tags, &res.Content, &mtime, &meta.Dynamic, &originImageURLs, &res.ListID, &meta.Media.MediaID, &meta.Media.Spoiler); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		PromError("db:读取草稿")
		log.Error("ArtDraft.row.Scan() error(%d,%d,%v)", mid, aid, err)
		return
	}
	meta.Category = category
	meta.Author = author
	meta.Mtime = xtime.Time(mtime.Unix())
	if imageURLs == "" {
		meta.ImageURLs = []string{}
	} else {
		meta.ImageURLs = strings.Split(imageURLs, ",")
	}
	if originImageURLs == "" {
		meta.OriginImageURLs = []string{}
	} else {
		meta.OriginImageURLs = strings.Split(originImageURLs, ",")
	}
	if tags == "" {
		res.Tags = []string{}
	} else {
		res.Tags = strings.Split(tags, ",")
	}
	res.Meta = meta
	return
}

// UpperDrafts batch get draft by mid.
func (d *Dao) UpperDrafts(c context.Context, mid int64, start, ps int) (res []*artmdl.Draft, err error) {
	var (
		rows   *xsql.Rows
		sqlStr = fmt.Sprintf(_upperDraftsSQL, d.hit(mid))
	)
	if rows, err = d.articleDB.Query(c, sqlStr, mid, start, ps); err != nil {
		PromError("db:读取草稿")
		log.Error("d.articleDB.Query(%d,%d,%d) error(%+v)", mid, start, ps, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			tags                       string
			imageURLs, originImageURLs string
			category                   = &artmdl.Category{}
			author                     = &artmdl.Author{}
			art                        = &artmdl.Draft{Article: &artmdl.Article{}}
			meta                       = &artmdl.Meta{}
			mtime                      time.Time
			listID                     int64
		)
		if err = rows.Scan(&meta.ID, &category.ID, &meta.Title, &meta.Summary, &meta.TemplateID, &meta.Reprint, &imageURLs, &tags, &mtime, &meta.Dynamic, &originImageURLs, &listID); err != nil {
			log.Error("UpperDrafts.row.Scan() error(%d,%d,%d,%v)", mid, start, ps, err)
			return
		}
		meta.Category = category
		meta.Author = author
		meta.Mtime = xtime.Time(mtime.Unix())
		if imageURLs == "" {
			meta.ImageURLs = []string{}
		} else {
			meta.ImageURLs = strings.Split(imageURLs, ",")
		}
		if originImageURLs == "" {
			meta.OriginImageURLs = []string{}
		} else {
			meta.OriginImageURLs = strings.Split(originImageURLs, ",")
		}
		if tags == "" {
			art.Tags = []string{}
		} else {
			art.Tags = strings.Split(tags, ",")
		}
		art.Meta = meta
		art.ListID = listID
		res = append(res, art)
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// AddArtDraft add article draft .
func (d *Dao) AddArtDraft(c context.Context, a *artmdl.Draft) (id int64, err error) {
	var (
		deleted         bool
		res             sql.Result
		tags            = strings.Join(a.Tags, ",")
		imageUrls       = strings.Join(a.ImageURLs, ",")
		originImageUrls = strings.Join(a.OriginImageURLs, ",")
		sqlStr          = fmt.Sprintf(_addArticleDraftSQL, d.hit(a.Author.Mid))
	)
	if a.ID > 0 {
		if deleted, err = d.IsDraftDeleted(c, a.Author.Mid, a.ID); err != nil {
			return
		} else if deleted {
			err = ecode.ArtCreationDraftDeleted
			return
		}
	}
	if res, err = d.articleDB.Exec(c, sqlStr, a.ID, a.Category.ID, a.Title, a.Summary, a.BannerURL, a.TemplateID, a.Author.Mid, a.Reprint, imageUrls, tags, a.Content, a.Dynamic, originImageUrls, a.ListID, a.Media.MediaID, a.Media.Spoiler); err != nil {
		PromError("db:新增或更新草稿")
		log.Error("d.articleDB.Exec(%+v) error(%+v)", a, err)
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		PromError("db:获取新增草稿ID")
		log.Error("res.LastInsertId() error(%+v)", err)
	}
	return
}

// IsDraftDeleted judges is draft has been deleted.
func (d *Dao) IsDraftDeleted(c context.Context, mid, aid int64) (deleted bool, err error) {
	var (
		dt     int
		sqlStr = fmt.Sprintf(_checkDraftSQL, d.hit(mid))
	)
	if err = d.articleDB.QueryRow(c, sqlStr, aid).Scan(&dt); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		PromError("db:判断草稿是否被删除")
		log.Error("d.articleDB.QueryRow(%d,%d) error(%+v)", mid, aid, err)
		return
	}
	if dt > 0 {
		deleted = true
	}
	return
}

// TxDeleteArticleDraft deletes article draft via transaction.
func (d *Dao) TxDeleteArticleDraft(c context.Context, tx *xsql.Tx, mid, aid int64) (err error) {
	var (
		now    = time.Now().Unix()
		sqlStr = fmt.Sprintf(_deleteArticleDraftSQL, d.hit(mid))
	)
	if _, err = tx.Exec(sqlStr, now, aid); err != nil {
		PromError("db:删除草稿")
		log.Error("TxDeleteArticleDraft.Exec(%d,%d) error(%+v)", mid, aid, err)
	}
	return
}

// DelArtDraft deletes article draft.
func (d *Dao) DelArtDraft(c context.Context, mid, aid int64) (err error) {
	var (
		now    = time.Now().Unix()
		sqlStr = fmt.Sprintf(_deleteArticleDraftSQL, d.hit(mid))
	)
	if _, err = d.articleDB.Exec(c, sqlStr, now, aid); err != nil {
		PromError("db:删除草稿")
		log.Error("d.articleDB.Exec(%d,%d) error(%+v)", mid, aid, err)
	}
	return
}

// CountUpperDraft count upper's draft
func (d *Dao) CountUpperDraft(c context.Context, mid int64) (count int, err error) {
	var sqlStr = fmt.Sprintf(_countUpperDraftSQL, d.hit(mid))
	if err = d.articleDB.QueryRow(c, sqlStr, mid).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		PromError("db:读取草稿计数")
		log.Error("CountUpperDraft.row.Scan() error(%d,%v)", mid, err)
	}
	return
}
