package dao

import (
	"context"
	"fmt"
	"time"

	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/app/job/openplatform/article/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_sharding = 100

	// stat
	_statSQL   = "SELECT article_id,favorite,reply,share,likes,dislike,view,coin FROM article_stats_%s WHERE article_id=%d and deleted_time =0"
	_upStatSQL = `INSERT INTO article_stats_%s (article_id,favorite,reply,share,likes,dislike,view,coin) VALUES (?,?,?,?,?,?,?,?) 
					ON DUPLICATE KEY UPDATE favorite=?,reply=?,share=?,likes=?,dislike=?,view=?,coin=?`
	_updateSearch = `INSERT INTO search_articles(ctime, article_id, category_id, title, summary, template_id, mid, image_urls, publish_time, content, tags, stats_view, stats_favorite, stats_likes, stats_dislike, stats_reply, stats_share, stats_coin, origin_image_urls, attributes, keywords)
	     VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE ctime = ?, category_id =?, title=?, summary=?, template_id=?, mid=?, image_urls=?, publish_time=?, content=?, tags=?, stats_view=?, stats_favorite=?, stats_likes=?, stats_dislike=?, stats_reply=?, stats_share=?, stats_coin = ?, origin_image_urls = ?, attributes = ?, keywords = ?`
	_delSearch                 = "DELETE FROM search_articles where article_id = ?"
	_articleContentSQL         = "SELECT content FROM filtered_article_contents_%s WHERE article_id = ?"
	_updateSearchStats         = "UPDATE search_articles SET stats_view=?, stats_favorite=?, stats_likes=?, stats_dislike=?, stats_reply=?, stats_share=?, stats_coin =? where article_id = ?"
	_gameList                  = "SELECT mid FROM white_list_users WHERE group_id = 4 AND deleted = 0"
	_allCheat                  = "SELECT article_id, lv FROM stats_filters WHERE deleted_time = 0"
	_newestArtsMetaSQL         = "SELECT article_id, publish_time, attributes FROM filtered_articles ORDER BY publish_time DESC LIMIT ?"
	_newestArtCategorysMetaSQL = "SELECT article_id, publish_time, attributes FROM filtered_articles where category_id in (%s) ORDER BY publish_time DESC LIMIT %d"
	_searchArticles            = "select article_id, category_id, attributes, stats_view, stats_reply, stats_favorite, stats_likes, stats_coin from search_articles where publish_time >= ? and publish_time < ?"
	_checkStateSQL             = "SELECT publish_time,check_state FROM articles WHERE id = ? and deleted_time = 0"
	_updateCheckState          = "UPDATE articles SET check_state = 3 WHERE id = ?"
	_settingsSQL               = "SELECT name,value FROM article_settings WHERE deleted_time=0"
	_midsByPublishTimeSQL      = "SELECT mid FROM articles WHERE publish_time > ? and state = 0 and deleted_time = 0 group by mid"
	_statByMidSQL              = "select count(*), sum(words), category_id from articles where mid = ? and state = 0 and deleted_time = 0 group by category_id"
	_keywordsSQL               = "select tags from article_contents_%s where article_id = ?"
	_actIDSQL                  = "select act_id from articles where id = ?"
	_lastModsArtsSQL           = "select id from articles where state in (0,5,6,7) and deleted_time = 0 order by mtime desc limit ?"
)

var _searchInterval = int64(3 * 24 * 3600)

func (d *Dao) hit(id int64) string {
	return fmt.Sprintf("%02d", id%_sharding)
}

// Stat returns stat info.
func (d *Dao) Stat(c context.Context, aid int64) (stat *artmdl.StatMsg, err error) {
	stat = &artmdl.StatMsg{}
	err = d.db.QueryRow(c, fmt.Sprintf(_statSQL, d.hit(aid), aid)).Scan(&stat.Aid, &stat.Favorite, &stat.Reply, &stat.Share, &stat.Like, &stat.Dislike, &stat.View, &stat.Coin)
	if err == sql.ErrNoRows {
		err = nil
		stat = nil
	} else if err != nil {
		log.Error("Stat(%v) error(%+v)", aid, err)
		PromError("db:读取计数")
	}
	return
}

// Update updates stat in db.
func (d *Dao) Update(c context.Context, stat *artmdl.StatMsg) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_upStatSQL, d.hit(stat.Aid)), stat.Aid, *stat.Favorite, *stat.Reply, *stat.Share, *stat.Like, *stat.Dislike, *stat.View, *stat.Coin,
		*stat.Favorite, *stat.Reply, *stat.Share, *stat.Like, *stat.Dislike, *stat.View, *stat.Coin)
	if err != nil {
		log.Error("Update(%d,%+v) error(%+v)", stat.Aid, stat, err)
		PromError("db:更新计数")
		return
	}
	rows, err = res.RowsAffected()
	return
}

// UpdateSearch update search article table
func (d *Dao) UpdateSearch(c context.Context, a *model.SearchArticle) (err error) {
	_, err = d.updateSearchStmt.Exec(c, a.CTime, a.ID, a.CategoryID, a.Title, a.Summary, a.TemplateID, a.Mid, a.ImageURLs, a.PublishTime, a.Content, a.Tags, a.StatsView, a.StatsFavorite, a.StatsLikes, a.StatsDisLike, a.StatsReply, a.StatsShare, a.StatsCoin, a.OriginImageURLs, a.Attributes, a.Keywords,
		a.CTime, a.CategoryID, a.Title, a.Summary, a.TemplateID, a.Mid, a.ImageURLs, a.PublishTime, a.Content, a.Tags, a.StatsView, a.StatsFavorite, a.StatsLikes, a.StatsDisLike, a.StatsReply, a.StatsShare, a.StatsCoin, a.OriginImageURLs, a.Attributes, a.Keywords)
	if err != nil {
		PromError("db:更新搜索表")
		log.Error("UpdateSearch(%+v) error(%+v)", a, err)
	}
	return
}

// DelSearch del search article table
func (d *Dao) DelSearch(c context.Context, aid int64) (err error) {
	_, err = d.delSearchStmt.Exec(c, aid)
	if err != nil {
		PromError("db:删除搜索表")
		log.Error("DelSearch(%v) error(%+v)", aid, err)
		return
	}
	return
}

// UpdateRecheck update recheck  table
func (d *Dao) UpdateRecheck(c context.Context, aid int64) (err error) {
	_, err = d.updateRecheckStmt.Exec(c, aid)
	if err != nil {
		PromError("db:修改回查状态")
		log.Error("UpdateRecheck(%v) error(%+v)", aid, err)
	}
	return
}

// ArticleContent get article content
func (d *Dao) ArticleContent(c context.Context, id int64) (res string, err error) {
	contentSQL := fmt.Sprintf(_articleContentSQL, d.hit(id))
	if err = d.db.QueryRow(c, contentSQL, id).Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		PromError("db:ArticleContent")
		log.Error("dao.ArticleContent(%s) error(%+v)", contentSQL, err)
	}
	return
}

// GetRecheckInfo get article recheck
func (d *Dao) GetRecheckInfo(c context.Context, id int64) (publishTime int64, checkState int, err error) {
	if err = d.getRecheckStmt.QueryRow(c, id).Scan(&publishTime, &checkState); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		PromError("db:ReCheckQuery")
		log.Error("dao.GetRecheckInfo(%d) error(%+v)", id, err)
	}
	return
}

// UpdateSearchStats update search stats
func (d *Dao) UpdateSearchStats(c context.Context, stat *artmdl.StatMsg) (err error) {
	_, err = d.updateSearchStatsStmt.Exec(c, *stat.View, *stat.Favorite, *stat.Like, *stat.Dislike, *stat.Reply, *stat.Share, *stat.Coin, stat.Aid)
	if err != nil {
		log.Error("updateSearchStatsStmt(%d,%+v) error(%+v)", stat.Aid, stat, err)
		PromError("db:更新搜索计数")
	}
	return
}

// RawGameList game list
func (d *Dao) RawGameList(c context.Context) (mids []int64, err error) {
	rows, err := d.gameStmt.Query(c)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if rows.Scan(&id); err != nil {
			PromError("db:GameList")
			log.Error("dao.GameList() error(%+v)", err)
			return
		}
		mids = append(mids, id)
	}
	return
}

// CheatArts cheat list
func (d *Dao) CheatArts(c context.Context) (res map[int64]int, err error) {
	rows, err := d.cheatStmt.Query(c)
	if err != nil {
		return
	}
	res = make(map[int64]int)
	defer rows.Close()
	defer func() {
		if err != nil {
			PromError("db:CheatArts")
			log.Error("dao.CheatArts() error(%+v)", err)
		}
	}()
	for rows.Next() {
		var id int64
		var lv int
		if rows.Scan(&id, &lv); err != nil {
			return
		}
		res[id] = lv
	}
	err = rows.Err()
	return
}

// NewestArtIDs find newest article's id
func (d *Dao) NewestArtIDs(c context.Context, limit int64) (res [][2]int64, err error) {
	rows, err := d.newestArtsMetaStmt.Query(c, limit)
	if err != nil {
		PromError("db:最新文章")
		log.Error("dao.newestArtsMetaStmt.Query error(%+v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id        int64
			ptime     int64
			attribute int32
		)
		if err = rows.Scan(&id, &ptime, &attribute); err != nil {
			PromError("db:最新文章scan")
			log.Error("dao.NewestArtIDs.rows.Scan error(%+v)", err)
			return
		}
		if artmdl.NoDistributeAttr(attribute) || artmdl.NoRegionAttr(attribute) {
			continue
		}
		res = append(res, [2]int64{id, ptime})
	}
	if err = rows.Err(); err != nil {
		PromError("db:最新文章")
		log.Error("dao.NewestArtIDs.rows error(%+v)", err)
	}
	return
}

// NewestArtIDByCategory find newest article's id
func (d *Dao) NewestArtIDByCategory(c context.Context, cids []int64, limit int64) (res [][2]int64, err error) {
	sql := fmt.Sprintf(_newestArtCategorysMetaSQL, xstr.JoinInts(cids), limit)
	rows, err := d.db.Query(c, sql)
	if err != nil {
		PromError("db:最新文章")
		log.Error("dao.NewestArtIDByCategorys.Query error(%+v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id        int64
			ptime     int64
			attribute int32
		)
		if err = rows.Scan(&id, &ptime, &attribute); err != nil {
			PromError("db:最新文章scan")
			log.Error("dao.NewestArtIDByCategory.rows.Scan error(%+v)", err)
			return
		}
		if artmdl.NoDistributeAttr(attribute) || artmdl.NoRegionAttr(attribute) {
			continue
		}
		res = append(res, [2]int64{id, ptime})
	}
	if err = rows.Err(); err != nil {
		PromError("db:最新文章")
		log.Error("dao.NewestArtIDByCategory.Scan error(%+v)", err)
	}
	return
}

// SearchArts get articles publish time after ptime
func (d *Dao) SearchArts(c context.Context, ptime int64) (res []*model.SearchArticle, err error) {
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
			ba := &model.SearchArticle{}
			if err = rows.Scan(&ba.ID, &ba.CategoryID, &ba.Attributes, &ba.StatsView, &ba.StatsReply, &ba.StatsFavorite, &ba.StatsLikes, &ba.StatsCoin); err != nil {
				PromError("db:searchArts")
				log.Error("mysql: search arts Scan() error(%+v)", err)
				return
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
	if err = rows.Err(); err != nil {
		PromError("db:loadSettings")
		log.Error("mysql: load settings Query() error(%+v)", err)
	}
	return
}

//MidsByPublishTime get mids by publish time
func (d *Dao) MidsByPublishTime(c context.Context, pubTime int64) (mids []int64, err error) {
	var rows *sql.Rows
	if rows, err = d.midByPubtimeStmt.Query(c, pubTime); err != nil {
		PromError("db:查询近7天mid")
		log.Error("mysql: db.MidsByPublishTime.Query error(%+v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid int64
		if err = rows.Scan(&mid); err != nil {
			PromError("查询近7天mid scan")
			log.Error("mysql: rows.Scan error(%+v)", err)
			return
		}
		mids = append(mids, mid)
	}
	err = rows.Err()
	if err = rows.Err(); err != nil {
		PromError("db:MidsByPublishTime")
		log.Error("mysql: get mids by publish time Query() error(%+v)", err)
	}
	return
}

//StatByMid get author info by mid
func (d *Dao) StatByMid(c context.Context, mid int64) (res map[int64][2]int64, err error) {
	res = make(map[int64][2]int64)
	var rows *sql.Rows
	if rows, err = d.statByMidStmt.Query(c, mid); err != nil {
		PromError("db:作者分区数据")
		log.Error("mysql: db.statByMidStmt.Query error(%+v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			data [2]int64
			cate int64
		)
		if err = rows.Scan(&data[0], &data[1], &cate); err != nil {
			PromError("作者分区数据 scan")
			log.Error("mysql: rows.Scan error(%+v)", err)
			return
		}
		res[cate] = data
	}
	err = rows.Err()
	if err = rows.Err(); err != nil {
		PromError("db:StatByMid")
		log.Error("mysql: get stat infos by mid Query() error(%+v)", err)
	}
	return
}

// Keywords .
func (d *Dao) Keywords(c context.Context, id int64) (keywords string, err error) {
	var sqlStr = fmt.Sprintf(_keywordsSQL, d.hit(id))
	row := d.db.QueryRow(c, sqlStr, id)
	if err = row.Scan(&keywords); err != nil {
		PromError("db:Keywords")
		log.Error("d.keywords scan error(%+v)", err)
	}
	return
}

// IsAct .
func (d *Dao) IsAct(c context.Context, id int64) (res bool) {
	var actID int64
	if err := d.db.QueryRow(c, _actIDSQL, id).Scan(&actID); err != nil {
		PromError("db:IsAct")
		log.Error("d.IsAct scan error(%+v)", err)
		return
	}
	if actID > 0 {
		res = true
	}
	return
}

// LastModIDs .
func (d *Dao) LastModIDs(c context.Context, size int) (aids []int64, err error) {
	var (
		id   int64
		rows *sql.Rows
	)
	if rows, err = d.db.Query(c, _lastModsArtsSQL, size); err != nil {
		PromError("db:LastModIDs")
		log.Error("d.LastModIDs query error(%+v) size(%d)", err, size)
		return
	}
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			PromError("db:LastModIDs")
			log.Error("d.LastModIDs scan error(%+v) size(%d)", err, size)
			return
		}
		aids = append(aids, id)
	}
	return
}
