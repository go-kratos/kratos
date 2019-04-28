package dao

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"

	"go-common/library/sync/errgroup"
)

const (
	_sharding      = 100
	_mysqlBulkSize = 50
	// article
	_articleMetaSQL           = "SELECT article_id,category_id,title,summary,banner_url, template_id, state, mid, reprint, image_urls, publish_time, ctime, attributes,words,dynamic_intro, origin_image_urls, media_id, spoiler FROM filtered_articles WHERE article_id = ?"
	_allArticleMetaSQL        = "SELECT id,category_id,title,summary,banner_url, template_id, state, mid, reprint, image_urls, publish_time, ctime, attributes,words,dynamic_intro, origin_image_urls, media_id, spoiler FROM articles WHERE id = ?"
	_articlesMetaSQL          = "SELECT article_id,category_id,title,summary,banner_url, template_id, state, mid, reprint, image_urls, publish_time, ctime, attributes,words,dynamic_intro, origin_image_urls, media_id, spoiler FROM filtered_articles WHERE article_id in (%s)"
	_upperPassedSQL           = "SELECT article_id, publish_time, attributes FROM filtered_articles WHERE mid = ? ORDER BY publish_time desc"
	_uppersPassedSQL          = "SELECT article_id, mid, publish_time, attributes FROM filtered_articles WHERE mid in (%s) ORDER BY publish_time desc"
	_articleContentSQL        = "SELECT content FROM filtered_article_contents_%s WHERE article_id = ?"
	_articleKeywordsSQL       = "SELECT tags FROM article_contents_%s WHERE article_id = ?"
	_articleUpperCountSQL     = "SELECT count(*) FROM filtered_articles WHERE mid = ?"
	_articleUpCntTodaySQL     = "SELECT count(*) FROM articles WHERE mid = ? and ctime >= ?"
	_delFilteredArtMetaSQL    = "DELETE FROM filtered_articles where article_id = ?"
	_delFilteredArtContentSQL = "DELETE FROM filtered_article_contents_%s where article_id = ?"
	// stat
	_statSQL  = "SELECT view,favorite,likes,dislike,reply,share,coin,dynamic FROM article_stats_%s WHERE article_id = ? and deleted_time = 0"
	_statsSQL = "SELECT article_id, view,favorite,likes,dislike,reply,share,coin,dynamic FROM article_stats_%s WHERE article_id in (%s) and deleted_time = 0"
	// category
	_categoriesSQL = "SELECT id,parent_id,name,position,banner_url FROM article_categories WHERE state = 1 and deleted_time = 0"
	// authors
	_authorsSQL    = "SELECT mid, daily_limit, state FROM article_authors WHERE deleted_time=0"
	_authorSQL     = "SELECT state,rtime, daily_limit FROM article_authors WHERE mid=? AND deleted_time=0"
	_applyCountSQL = "SELECT count(*) FROM article_authors WHERE atime >= ?"
	_applySQL      = "INSERT INTO article_authors (mid,atime,count,content,category) VALUES (?,?,1,?,?) ON DUPLICATE KEY UPDATE atime=?,rtime=0,state=0,count=count+1,content=?,category=?,deleted_time=0"
	_addAuthorSQL  = "INSERT INTO article_authors (mid,state,type) VALUES (?,1,5) ON DUPLICATE KEY UPDATE state=1,deleted_time=0"
	// recommends
	_recommendCategorySQL = "SELECT article_id, big_banner_url, show_recommend, position, end_time, big_banner_start_time, big_banner_end_time  FROM article_recommends WHERE start_time <= ? and (end_time >= ? or end_time = 0) and category_id = ? and deleted_time = 0 ORDER BY position ASC"
	_allRecommendSQL      = "SELECT article_id FROM article_recommends WHERE start_time <= ? and (end_time >= ? or end_time = 0) and deleted_time = 0 and category_id = 0 ORDER BY mtime DESC LIMIT ?,?"
	_allRecommendCountSQL = "SELECT COUNT(*) FROM article_recommends WHERE start_time <= ? and (end_time >= ? or end_time = 0) and deleted_time = 0 and category_id = 0"
	_deleteRecommendSQL   = "UPDATE article_recommends SET deleted_time=? WHERE article_id=? and deleted_time = 0"
	// setting
	_settingsSQL = "SELECT name,value FROM article_settings WHERE deleted_time=0"
	// sort
	_newestArtsMetaSQL = "SELECT article_id, publish_time, attributes FROM filtered_articles ORDER BY publish_time DESC LIMIT ?"
	//notice
	_noticeSQL = "SELECT id, title, url, plat, condi, build from article_notices where state = 1 and stime <= ? and etime > ?"
	// users
	_userNoticeSQL       = "SELECT notice_state from users where mid = ?"
	_updateUserNoticeSQL = "INSERT INTO users (mid,notice_state) VALUES (?,?) ON DUPLICATE KEY UPDATE notice_state=?"
	// hotspots
	_hotspotsSQL = "select id, title, tag, icon, top_articles from hotspots where deleted_time = 0 and `order` != 0 order by `order` asc"
	// search articles
	_searchArticles  = "select article_id, publish_time, tags, stats_view, stats_reply from search_articles where publish_time >= ? and publish_time < ?"
	_addCheatSQL     = "INSERT INTO stats_filters(article_id, lv) VALUES(?,?) ON DUPLICATE KEY UPDATE lv=?, deleted_time = 0"
	_delCheatSQL     = "UPDATE stats_filters SET deleted_time = ? WHERE article_id = ? and deleted_time = 0"
	_tagArticlesSQL  = "select tid, oid, log_date FROM article_tags where tid in (%s) and is_deleted = 0"
	_mediaArticleSQL = "select id from articles where mid = ? and media_id = ? and deleted_time = 0 and state > -10"
	_mediaByIDSQL    = "select media_id from articles where id = ?"
)

var _searchInterval = int64(3 * 24 * 3600)

func (d *Dao) hit(id int64) string {
	return fmt.Sprintf("%02d", id%_sharding)
}

// Categories get Categories
func (d *Dao) Categories(c context.Context) (res map[int64]*model.Category, err error) {
	var rows *sql.Rows
	if rows, err = d.categoriesStmt.Query(c); err != nil {
		PromError("db:分区查询")
		log.Error("mysql: db.Categories.Query error(%+v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.Category)
	for rows.Next() {
		ca := &model.Category{}
		if err = rows.Scan(&ca.ID, &ca.ParentID, &ca.Name, &ca.Position, &ca.BannerURL); err != nil {
			PromError("分区Scan")
			log.Error("mysql: rows.Categories.Scan error(%+v)", err)
			return
		}
		res[ca.ID] = ca
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// ArticleStats get article stats
func (d *Dao) ArticleStats(c context.Context, id int64) (res *model.Stats, err error) {
	res = new(model.Stats)
	row := d.articleDB.QueryRow(c, fmt.Sprintf(_statSQL, d.hit(id)), id)
	if err = row.Scan(&res.View, &res.Favorite, &res.Like, &res.Dislike, &res.Reply, &res.Share, &res.Coin, &res.Dynamic); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
		} else {
			PromError("Stat scan")
			log.Error("mysql: ArticleStats row.Scan(%d) error(%+v)", id, err)
		}
	}
	return
}

// ArticlesStats get articles stats
func (d *Dao) ArticlesStats(c context.Context, ids []int64) (res map[int64]*model.Stats, err error) {
	var (
		shardings = make(map[int64][]int64)
		group     = &errgroup.Group{}
		mutex     = &sync.Mutex{}
	)
	res = make(map[int64]*model.Stats)
	for _, id := range ids {
		shardings[id%_sharding] = append(shardings[id%_sharding], id)
	}

	for sharding, subIDs := range shardings {
		keysLen := len(subIDs)
		sharding := sharding
		subIDs := subIDs
		for i := 0; i < keysLen; i += _mysqlBulkSize {
			var keys []int64
			if (i + _mysqlBulkSize) > keysLen {
				keys = subIDs[i:]
			} else {
				keys = subIDs[i : i+_mysqlBulkSize]
			}
			group.Go(func() error {
				statsSQL := fmt.Sprintf(_statsSQL, d.hit(sharding), xstr.JoinInts(keys))
				rows, e := d.articleDB.Query(c, statsSQL)
				if e != nil {
					return e
				}
				defer rows.Close()
				for rows.Next() {
					s := &model.Stats{}
					var aid int64
					e = rows.Scan(&aid, &s.View, &s.Favorite, &s.Like, &s.Dislike, &s.Reply, &s.Share, &s.Coin, &s.Dynamic)
					if e != nil {
						return e
					}
					mutex.Lock()
					res[aid] = s
					mutex.Unlock()
				}
				return rows.Err()
			})
		}
	}
	err = group.Wait()
	if err != nil {
		PromError("stats Scan")
		log.Error("mysql: rows.ArticleStats.Scan error(%+v)", err)
	}
	if len(res) == 0 {
		res = nil
	}
	return
}

// Settings gets article settings.
func (d *Dao) Settings(c context.Context) (res map[string]string, err error) {
	var rows *sql.Rows
	if rows, err = d.settingsStmt.Query(c); err != nil {
		PromError("db:文章配置查询")
		log.Error("mysql: db.settingsStmt.Query error(%+v)", err)
		return
	}
	defer rows.Close()
	res = make(map[string]string)
	for rows.Next() {
		var name, value string
		if err = rows.Scan(&name, &value); err != nil {
			PromError("文章配置scan")
			log.Error("mysql: rows.Scan error(%+v)", err)
			return
		}
		res[name] = value
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// Notices notice .
func (d *Dao) Notices(c context.Context, t time.Time) (res []*model.Notice, err error) {
	var rows *sql.Rows
	if rows, err = d.noticeStmt.Query(c, t, t); err != nil {
		PromError("db:notice")
		log.Error("mysql: notice Query() error(%+v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ba := &model.Notice{}
		if err = rows.Scan(&ba.ID, &ba.Title, &ba.URL, &ba.Plat, &ba.Condition, &ba.Build); err != nil {
			PromError("db:notice")
			log.Error("mysql: notice Scan() error(%+v)", err)
			return
		}
		res = append(res, ba)
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// NoticeState .
func (d *Dao) NoticeState(c context.Context, mid int64) (res int64, err error) {
	if err = d.userNoticeStmt.QueryRow(c, mid).Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			PromError("db:notice_state")
			log.Error("mysql: notice state row.Scan error(%+v)", err)
		}
	}
	return
}

// UpdateNoticeState update notice state
func (d *Dao) UpdateNoticeState(c context.Context, mid int64, state int64) (err error) {
	if _, err = d.updateUserNoticeStmt.Exec(c, mid, state, state); err != nil {
		PromError("db:修改用户引导状态")
		log.Error("mysql: update_notice state(mid: %v, state: %v) error(%+v)", mid, state, err)
	}
	return
}

// Hotspots .
func (d *Dao) Hotspots(c context.Context) (res []*model.Hotspot, err error) {
	var rows *sql.Rows
	if rows, err = d.hotspotsStmt.Query(c); err != nil {
		PromError("db:hotspots")
		log.Error("mysql: hotspot Query() error(%+v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ba := &model.Hotspot{}
		var ic int
		var arts string
		if err = rows.Scan(&ba.ID, &ba.Title, &ba.Tag, &ic, &arts); err != nil {
			PromError("db:hostspot")
			log.Error("mysql: hotspot Scan() error(%+v)", err)
			return
		}
		if ic != 0 {
			ba.Icon = true
		}
		ba.TopArticles, _ = xstr.SplitInts(arts)
		res = append(res, ba)
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// SearchArts get articles publish time after ptime
func (d *Dao) SearchArts(c context.Context, ptime int64) (res []*model.SearchArt, err error) {
	var rows *sql.Rows
	now := time.Now().Unix()
	for ; ptime < now; ptime += _searchInterval {
		if rows, err = d.searchArtsStmt.Query(c, ptime, ptime+_searchInterval); err != nil {
			PromError("db:searchArts")
			log.Error("mysql: search arts Query() error(%+v)", err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			ba := &model.SearchArt{}
			var t string
			if err = rows.Scan(&ba.ID, &ba.PublishTime, &t, &ba.StatsView, &ba.StatsReply); err != nil {
				PromError("db:searchArts")
				log.Error("mysql: search arts Scan() error(%+v)", err)
				return
			}
			if t != "" {
				ba.Tags = strings.Split(t, ",")
			}
			res = append(res, ba)
		}
		if err = rows.Err(); err != nil {
			PromError("db:searchArts")
			log.Error("mysql: search arts Query() error(%+v)", err)
			return
		}
	}
	return
}

// AddCheatFilter .
func (d *Dao) AddCheatFilter(c context.Context, aid int64, lv int) (err error) {
	if _, err = d.addCheatStmt.Exec(c, aid, lv, lv); err != nil {
		PromError("db:新增防刷过滤")
		log.Error("mysql: addCheatFilter state(aid: %v, lv: %v) error(%+v)", aid, lv, err)
		return
	}
	log.Info("mysql: addCheatFilter state(aid: %v, lv: %v)", aid, lv)
	return
}

// DelCheatFilter .
func (d *Dao) DelCheatFilter(c context.Context, aid int64) (err error) {
	if _, err = d.delCheatStmt.Exec(c, time.Now().Unix(), aid); err != nil {
		PromError("db:删除防刷过滤")
		log.Error("mysql: delCheatFilter state(aid: %v) error(%+v)", aid, err)
		return
	}
	log.Info("mysql: delCheatFilter state(aid: %v)", aid)
	return
}
