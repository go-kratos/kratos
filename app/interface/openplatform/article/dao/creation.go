package dao

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"net/http"
	"strconv"
	"strings"
	"time"

	"database/sql"
	artmdl "go-common/app/interface/openplatform/article/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

const (
	// article
	_addArticleMetaSQL              = "INSERT INTO articles (category_id,title,summary,banner_url,template_id,state,mid,reprint,image_urls,attributes,words,dynamic_intro,origin_image_urls,act_id,media_id,spoiler,apply_time) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_addArticleContentSQL           = "INSERT INTO article_contents_%s (article_id,content,tags) values (?,?,?)"
	_addArticleVersionSQL           = "INSERT INTO article_versions (article_id,category_id,title,state,content,summary,banner_url,template_id,mid,reprint,image_urls,attributes,words,dynamic_intro,origin_image_urls,act_id,media_id,spoiler,apply_time,ext_msg)values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_updateArticleVersionSQL        = "UPDATE article_versions SET category_id=?,title=?,state=?,content=?,summary=?,banner_url=?,template_id=?,mid=?,reprint=?,image_urls=?,attributes=?,words=?,dynamic_intro=?,origin_image_urls=?,spoiler=?,apply_time=?,ext_msg=? where article_id=? and deleted_time=0"
	_updateArticleMetaSQL           = "UPDATE articles SET category_id=?,title=?,summary=?,banner_url=?,template_id=?,state=?,mid=?,reprint=?,image_urls=?,attributes=?,words=?,dynamic_intro=?,origin_image_urls =?,spoiler=?,apply_time=? WHERE id=?"
	_updateArticleContentSQL        = "UPDATE article_contents_%s SET content=?, tags=? WHERE article_id=?"
	_deleteArticleMetaSQL           = "UPDATE articles SET deleted_time=? WHERE id=?"
	_deleteArticleContentSQL        = "UPDATE article_contents_%s SET deleted_time=? WHERE article_id=?"
	_deleteArticleVerionSQL         = "UPDATE article_versions SET deleted_time=? WHERE article_id=?"
	_updateArticleStateSQL          = "UPDATE articles SET state=? WHERE id=?"
	_updateArticleStateApplyTimeSQL = "UPDATE articles SET state=?,apply_time=? WHERE id=?"
	_upperArticlesMetaCreationSQL   = `SELECT id,category_id,title,summary,banner_url,template_id,state,mid,reprint,image_urls,publish_time,ctime,reason, attributes, dynamic_intro, origin_image_urls FROM articles WHERE mid=? and deleted_time=0 
	                                 and state in (%s)`
	_upperArticleCountCreationSQL = "SELECT state FROM articles WHERE mid=? and deleted_time=0"
	_articleMetaCreationSQL       = "SELECT id,category_id,title,summary,banner_url, template_id, state, mid, reprint, image_urls, publish_time,ctime, attributes, dynamic_intro, origin_image_urls, media_id, spoiler FROM articles WHERE id = ? and deleted_time = 0"
	_articleContentCreationSQL    = "SELECT content FROM article_contents_%s WHERE article_id=? AND deleted_time=0"
	_countEditTimesSQL            = "SELECT count(*) FROM article_histories_%s WHERE article_id=? AND deleted_time=0 AND state in (5,6,7)"
	_articleVersionSQL            = "SELECT article_id,category_id,title,state,content,summary,banner_url,template_id,reprint,image_urls,attributes,words,dynamic_intro,origin_image_urls,media_id,spoiler,apply_time,ext_msg FROM article_versions WHERE article_id=? AND deleted_time=0"
	_reasonOfVersion              = "SELECT reason FROM article_versions WHERE article_id=? AND state=? AND deleted_time=0 ORDER BY id DESC LIMIT 1"
)

// TxAddArticleMeta adds article's meta via transaction.
func (d *Dao) TxAddArticleMeta(c context.Context, tx *xsql.Tx, a *artmdl.Meta, actID int64) (id int64, err error) {
	a.ImageURLs = artmdl.CleanURLs(a.ImageURLs)
	a.OriginImageURLs = artmdl.CleanURLs(a.OriginImageURLs)
	a.BannerURL = artmdl.CleanURL(a.BannerURL)
	var (
		res             sql.Result
		imageUrls       = strings.Join(a.ImageURLs, ",")
		originImageUrls = strings.Join(a.OriginImageURLs, ",")
		applyTime       = time.Now().Format("2006-01-02 15:04:05")
	)
	if res, err = tx.Exec(_addArticleMetaSQL, a.Category.ID, a.Title, a.Summary, a.BannerURL, a.TemplateID, a.State, a.Author.Mid, a.Reprint, imageUrls, a.Attributes, a.Words, a.Dynamic, originImageUrls, actID, a.Media.MediaID, a.Media.Spoiler, applyTime); err != nil {
		PromError("db:新增文章meta")
		log.Error("tx.Exec() error(%+v)", err)
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		log.Error("res.LastInsertId() error(%+v)", err)
	}
	return
}

// TxAddArticleContent adds article's body via transaction.
func (d *Dao) TxAddArticleContent(c context.Context, tx *xsql.Tx, aid int64, content string, tags []string) (err error) {

	var sqlStr = fmt.Sprintf(_addArticleContentSQL, d.hit(aid))
	if _, err = tx.Exec(sqlStr, aid, content, strings.Join(tags, "\001")); err != nil {
		PromError("db:新增文章content")
		log.Error("tx.Exec(%s,%d) error(%+v)", sqlStr, aid, err)
	}
	return
}

// TxAddArticleVersion adds article version.
func (d *Dao) TxAddArticleVersion(c context.Context, tx *xsql.Tx, id int64, a *artmdl.Article, actID int64) (err error) {
	a.ImageURLs = artmdl.CleanURLs(a.ImageURLs)
	a.OriginImageURLs = artmdl.CleanURLs(a.OriginImageURLs)
	a.BannerURL = artmdl.CleanURL(a.BannerURL)
	var (
		applyTime       = time.Now().Format("2006-01-02 15:04:05")
		imageUrls       = strings.Join(a.ImageURLs, ",")
		originImageUrls = strings.Join(a.OriginImageURLs, ",")
		extMsg          = &artmdl.ExtMsg{
			Tags: a.Tags,
		}
		extStr []byte
	)
	if extStr, err = json.Marshal(extMsg); err != nil {
		log.Error("json.Marshal error(%+v)", err)
		return
	}
	if _, err = tx.Exec(_addArticleVersionSQL, id, a.Category.ID, a.Title, a.State, a.Content, a.Summary, a.BannerURL, a.TemplateID, a.Author.Mid, a.Reprint, imageUrls, a.Attributes, a.Words, a.Dynamic, originImageUrls, actID, a.Media.MediaID, a.Media.Spoiler, applyTime, string(extStr)); err != nil {
		PromError("db:新增版本")
		log.Error("tx.Exec() error(%+v)", err)
	}
	return
}

// TxUpdateArticleVersion updates article version.
func (d *Dao) TxUpdateArticleVersion(c context.Context, tx *xsql.Tx, id int64, a *artmdl.Article, actID int64) (err error) {
	a.ImageURLs = artmdl.CleanURLs(a.ImageURLs)
	a.OriginImageURLs = artmdl.CleanURLs(a.OriginImageURLs)
	a.BannerURL = artmdl.CleanURL(a.BannerURL)
	var (
		applyTime       = time.Now().Format("2006-01-02 15:04:05")
		imageUrls       = strings.Join(a.ImageURLs, ",")
		originImageUrls = strings.Join(a.OriginImageURLs, ",")
		extMsg          = &artmdl.ExtMsg{
			Tags: a.Tags,
		}
		extStr []byte
	)
	if extStr, err = json.Marshal(extMsg); err != nil {
		log.Error("json.Marshal error(%+v)", err)
		return
	}
	if _, err = tx.Exec(_updateArticleVersionSQL, a.Category.ID, a.Title, a.State, a.Content, a.Summary, a.BannerURL, a.TemplateID, a.Author.Mid, a.Reprint, imageUrls, a.Attributes, a.Words, a.Dynamic, originImageUrls, a.Media.Spoiler, applyTime, string(extStr), id); err != nil {
		PromError("db:更新版本")
		log.Error("tx.Exec() error(%+v)", err)
	}
	return
}

// TxUpdateArticleMeta updates article's meta via transaction.
func (d *Dao) TxUpdateArticleMeta(c context.Context, tx *xsql.Tx, a *artmdl.Meta) (err error) {
	a.ImageURLs = artmdl.CleanURLs(a.ImageURLs)
	a.OriginImageURLs = artmdl.CleanURLs(a.OriginImageURLs)
	a.BannerURL = artmdl.CleanURL(a.BannerURL)
	var (
		imageURLs       = strings.Join(a.ImageURLs, ",")
		originImageURLs = strings.Join(a.OriginImageURLs, ",")
		applyTime       = time.Now().Format("2006-01-02 15:04:05")
	)
	if _, err = tx.Exec(_updateArticleMetaSQL, a.Category.ID, a.Title, a.Summary, a.BannerURL, a.TemplateID, a.State, a.Author.Mid, a.Reprint, imageURLs, a.Attributes, a.Words, a.Dynamic, originImageURLs, a.Media.Spoiler, applyTime, a.ID); err != nil {
		PromError("db:更新文章meta")
		log.Error("tx.Exec() error(%+v)", err)
	}
	return
}

// TxUpdateArticleContent updates article's body via transaction.
func (d *Dao) TxUpdateArticleContent(c context.Context, tx *xsql.Tx, aid int64, content string, tags []string) (err error) {
	var sqlStr = fmt.Sprintf(_updateArticleContentSQL, d.hit(aid))
	if _, err = tx.Exec(sqlStr, content, strings.Join(tags, "\001"), aid); err != nil {
		PromError("db:更新文章content")
		log.Error("tx.Exec(%s,%d) error(%+v)", sqlStr, aid, err)
	}
	return
}

// TxDeleteArticleMeta deletes article's meta via transaction.
func (d *Dao) TxDeleteArticleMeta(c context.Context, tx *xsql.Tx, aid int64) (err error) {
	var now = time.Now().Unix()
	if _, err = tx.Exec(_deleteArticleMetaSQL, now, aid); err != nil {
		PromError("db:删除文章meta")
		log.Error("tx.Exec() error(%+v)", err)
	}
	return
}

// TxDeleteArticleContent deletes article's meta via transaction.
func (d *Dao) TxDeleteArticleContent(c context.Context, tx *xsql.Tx, aid int64) (err error) {
	var (
		now    = time.Now().Unix()
		sqlStr = fmt.Sprintf(_deleteArticleContentSQL, d.hit(aid))
	)
	if _, err = tx.Exec(sqlStr, now, aid); err != nil {
		PromError("db:删除文章content")
		log.Error("tx.Exec(%s,%d,%d) error(%+v)", sqlStr, now, aid, err)
	}
	return
}

// TxDelArticleVersion deletes article version.
func (d *Dao) TxDelArticleVersion(c context.Context, tx *xsql.Tx, aid int64) (err error) {
	var now = time.Now().Unix()
	if _, err = tx.Exec(_deleteArticleVerionSQL, now, aid); err != nil {
		PromError("db:删除文章版本content")
		log.Error("tx.Exec(%s,%d,%d) error(%+v)", _deleteArticleVerionSQL, now, aid, err)
	}
	return
}

// TxDelFilteredArtMeta delete filetered article meta
func (d *Dao) TxDelFilteredArtMeta(c context.Context, tx *xsql.Tx, aid int64) (err error) {
	if _, err = tx.Exec(_delFilteredArtMetaSQL, aid); err != nil {
		PromError("db:删除过滤文章")
		log.Error("dao.DelFilteredArtMeta exec(%v) error(%+v)", aid, err)
	}
	return
}

//TxDelFilteredArtContent  delete filtered article content
func (d *Dao) TxDelFilteredArtContent(c context.Context, tx *xsql.Tx, aid int64) (err error) {
	contentSQL := fmt.Sprintf(_delFilteredArtContentSQL, d.hit(aid))
	if _, err = tx.Exec(contentSQL, aid); err != nil {
		PromError("db:删除过滤文章正文")
		log.Error("dao.DelFilteredArtContent exec(%v) error(%+v)", aid, err)
	}
	return
}

// UpdateArticleState updates article's state.
func (d *Dao) UpdateArticleState(c context.Context, aid int64, state int) (err error) {
	var res sql.Result
	if res, err = d.updateArticleStateStmt.Exec(c, state, aid); err != nil {
		PromError("db:更新文章状态")
		log.Error("s.dao.UpdateArticleState.Exec(aid: %v, state: %v) error(%+v)", aid, state, err)
		return
	}
	if count, _ := res.RowsAffected(); count == 0 {
		err = ecode.NothingFound
	}
	return
}

// TxUpdateArticleState updates article's state.
func (d *Dao) TxUpdateArticleState(c context.Context, tx *xsql.Tx, aid int64, state int32) (err error) {
	var res sql.Result
	if res, err = tx.Exec(_updateArticleStateSQL, state, aid); err != nil {
		PromError("db:更新文章状态")
		log.Error("s.dao.TxUpdateArticleState.Exec(aid: %v, state: %v) error(%+v)", aid, state, err)
		return
	}
	if count, _ := res.RowsAffected(); count == 0 {
		err = ecode.NothingFound
	}
	return
}

// TxUpdateArticleStateApplyTime updates article's state and apply time.
func (d *Dao) TxUpdateArticleStateApplyTime(c context.Context, tx *xsql.Tx, aid int64, state int32) (err error) {
	var (
		res       sql.Result
		applyTime = time.Now().Format("2006-01-02 15:03:04")
	)
	if res, err = tx.Exec(_updateArticleStateApplyTimeSQL, state, applyTime, aid); err != nil {
		PromError("db:更新文章状态和申请时间")
		log.Error("s.dao.TxUpdateArticleStateApplyTime.Exec(aid: %v, state: %v) error(%+v)", aid, state, err)
		return
	}
	if count, _ := res.RowsAffected(); count == 0 {
		err = ecode.NothingFound
	}
	return
}

// UpperArticlesMeta gets article list by mid.
func (d *Dao) UpperArticlesMeta(c context.Context, mid int64, group, category int) (as []*artmdl.Meta, err error) {
	var (
		rows   *xsql.Rows
		sqlStr string
	)
	sqlStr = fmt.Sprintf(_upperArticlesMetaCreationSQL, xstr.JoinInts(artmdl.Group2State(group)))
	if category > 0 {
		sqlStr += " and category_id=" + strconv.Itoa(category)
	}
	if rows, err = d.articleDB.Query(c, sqlStr, mid); err != nil {
		PromError("db:获取文章meta")
		log.Error("d.articleDB.Query(%s,%s) error(%+v)", sqlStr, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &artmdl.Meta{Category: &artmdl.Category{}, Author: &artmdl.Author{}}
		var (
			ptime                      int64
			ctime                      time.Time
			imageURLs, originImageURLs string
		)
		if err = rows.Scan(&a.ID, &a.Category.ID, &a.Title, &a.Summary, &a.BannerURL, &a.TemplateID, &a.State, &a.Author.Mid, &a.Reprint, &imageURLs, &ptime, &ctime, &a.Reason, &a.Attributes, &a.Dynamic, &originImageURLs); err != nil {
			promErrorCheck(err)
			log.Error("rows.Scan error(%+v)", err)
			return
		}
		if imageURLs == "" {
			a.ImageURLs = []string{}
		} else {
			a.ImageURLs = strings.Split(imageURLs, ",")
		}
		if originImageURLs == "" {
			a.OriginImageURLs = []string{}
		} else {
			a.OriginImageURLs = strings.Split(originImageURLs, ",")
		}
		a.PublishTime = xtime.Time(ptime)
		a.Ctime = xtime.Time(ctime.Unix())
		a.BannerURL = artmdl.CompleteURL(a.BannerURL)
		a.ImageURLs = artmdl.CompleteURLs(a.ImageURLs)
		a.OriginImageURLs = artmdl.CompleteURLs(a.OriginImageURLs)
		as = append(as, a)
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// UpperArticlesTypeCount gets article count by type.
func (d *Dao) UpperArticlesTypeCount(c context.Context, mid int64) (res *artmdl.CreationArtsType, err error) {
	var rows *xsql.Rows
	res = &artmdl.CreationArtsType{}
	if rows, err = d.upperArtCntCreationStmt.Query(c, mid); err != nil {
		PromError("db:获取各种文章状态总数")
		log.Error("d.articleDB.Query(%d) error(%+v)", mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var state int
		if err = rows.Scan(&state); err != nil {
			promErrorCheck(err)
			return
		}
		switch state {
		case artmdl.StateAutoLock, artmdl.StateLock, artmdl.StateReject, artmdl.StateOpenReject:
			res.NotPassed++
		case artmdl.StateAutoPass, artmdl.StateOpen, artmdl.StateRePass, artmdl.StateReReject:
			res.Passed++
		case artmdl.StatePending, artmdl.StateOpenPending, artmdl.StateRePending:
			res.Audit++
		}
	}
	err = rows.Err()
	promErrorCheck(err)
	res.All = res.NotPassed + res.Passed + res.Audit
	return
}

// CreationArticleMeta querys article's meta info for creation center by aid.
func (d *Dao) CreationArticleMeta(c context.Context, id int64) (am *artmdl.Meta, err error) {
	var (
		imageURLs, originImageURLs string
		category                   = &artmdl.Category{}
		author                     = &artmdl.Author{}
		ptime                      int64
		ct                         time.Time
	)
	am = &artmdl.Meta{Media: &artmdl.Media{}}
	if err = d.articleMetaCreationStmt.QueryRow(c, id).Scan(&am.ID, &category.ID, &am.Title, &am.Summary,
		&am.BannerURL, &am.TemplateID, &am.State, &author.Mid, &am.Reprint, &imageURLs, &ptime, &ct, &am.Attributes, &am.Dynamic, &originImageURLs, &am.Media.MediaID, &am.Media.Spoiler); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			am = nil
			return
		}
		PromError("db:文章内容表")
		log.Error("row.ArticleContent.QueryRow error(%+v)", err)
		return
	}
	am.Category = category
	am.Author = author
	am.Ctime = xtime.Time(ct.Unix())
	if imageURLs == "" {
		am.ImageURLs = []string{}
	} else {
		am.ImageURLs = strings.Split(imageURLs, ",")
	}
	if originImageURLs == "" {
		am.OriginImageURLs = []string{}
	} else {
		am.OriginImageURLs = strings.Split(originImageURLs, ",")
	}
	am.PublishTime = xtime.Time(ptime)
	am.BannerURL = artmdl.CompleteURL(am.BannerURL)
	am.ImageURLs = artmdl.CompleteURLs(am.ImageURLs)
	am.OriginImageURLs = artmdl.CompleteURLs(am.OriginImageURLs)
	return
}

// CreationArticleContent gets article's content.
func (d *Dao) CreationArticleContent(c context.Context, aid int64) (res string, err error) {
	contentSQL := fmt.Sprintf(_articleContentCreationSQL, d.hit(aid))
	if err = d.articleDB.QueryRow(c, contentSQL, aid).Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		PromError("db:CreationArticleContent")
		log.Error("dao.CreationArticleContent(%s) error(%+v)", contentSQL, err)
	}
	return
}

// UploadImage upload bfs.
func (d *Dao) UploadImage(c context.Context, fileType string, bs []byte) (location string, err error) {
	req, err := http.NewRequest(d.c.BFS.Method, d.c.BFS.URL, bytes.NewBuffer(bs))
	if err != nil {
		PromError("creation:UploadImage")
		log.Error("creation: http.NewRequest error (%v) | fileType(%s)", err, fileType)
		return
	}
	expire := time.Now().Unix()
	authorization := authorize(d.c.BFS.Key, d.c.BFS.Secret, d.c.BFS.Method, d.c.BFS.Bucket, expire)
	req.Header.Set("Host", d.c.BFS.URL)
	req.Header.Add("Date", fmt.Sprint(expire))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", fileType)
	// timeout
	ctx, cancel := context.WithTimeout(c, time.Duration(d.c.BFS.Timeout))
	req = req.WithContext(ctx)
	defer cancel()
	resp, err := d.bfsClient.Do(req)
	if err != nil {
		PromError("creation:UploadImage")
		log.Error("creation: d.Client.Do error(%v) | url(%s)", err, d.c.BFS.URL)
		err = ecode.BfsUploadServiceUnavailable
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("creation: Upload http.StatusCode nq http.StatusOK (%d) | url(%s)", resp.StatusCode, d.c.BFS.URL)
		PromError("creation:UploadImage")
		err = errors.New("Upload failed")
		return
	}
	header := resp.Header
	code := header.Get("Code")
	if code != strconv.Itoa(http.StatusOK) {
		log.Error("creation: strconv.Itoa err, code(%s) | url(%s)", code, d.c.BFS.URL)
		PromError("creation:UploadImage")
		err = errors.New("Upload failed")
		return
	}
	location = header.Get("Location")
	return
}

// authorize returns authorization for upload file to bfs
func authorize(key, secret, method, bucket string, expire int64) (authorization string) {
	var (
		content   string
		mac       hash.Hash
		signature string
	)
	content = fmt.Sprintf("%s\n%s\n\n%d\n", method, bucket, expire)
	mac = hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(content))
	signature = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	authorization = fmt.Sprintf("%s:%s:%d", key, signature, expire)
	return
}

// EditTimes count times of article edited.
func (d *Dao) EditTimes(c context.Context, id int64) (count int, err error) {
	var sqlStr = fmt.Sprintf(_countEditTimesSQL, d.hit(id))
	row := d.articleDB.QueryRow(c, sqlStr, id)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		PromError("db:EditTimes")
		log.Error("dao.EditTimes error(%+v)", err)
	}
	return
}

// ArticleVersion .
func (d *Dao) ArticleVersion(c context.Context, aid int64) (a *artmdl.Article, err error) {
	var (
		extStr                string
		extMsg                artmdl.ExtMsg
		imageURLs, originURLs string
	)
	row := d.articleDB.QueryRow(c, _articleVersionSQL, aid)
	a = &artmdl.Article{
		Meta: &artmdl.Meta{
			Media:    &artmdl.Media{},
			Category: &artmdl.Category{},
			List:     &artmdl.List{},
		},
	}
	if err = row.Scan(&a.ID, &a.Category.ID, &a.Title, &a.State, &a.Content, &a.Summary, &a.BannerURL, &a.TemplateID, &a.Reprint, &imageURLs, &a.Attributes, &a.Words, &a.Dynamic, &originURLs, &a.Media.MediaID, &a.Media.Spoiler, &a.ApplyTime, &extStr); err != nil {
		log.Error("dao.ArticleHistory.Scan error(%+v)", err)
		return
	}
	if err = json.Unmarshal([]byte(extStr), &extMsg); err != nil {
		log.Error("dao.ArticleHistory.Unmarshal error(%+v)", err)
		return
	}
	a.ImageURLs = strings.Split(imageURLs, ",")
	a.OriginImageURLs = strings.Split(originURLs, ",")
	a.BannerURL = artmdl.CompleteURL(a.BannerURL)
	a.ImageURLs = artmdl.CompleteURLs(a.ImageURLs)
	a.OriginImageURLs = artmdl.CompleteURLs(a.OriginImageURLs)
	a.Tags = extMsg.Tags
	return
}

// LastReason return last reason from article_versions by aid and state.
func (d *Dao) LastReason(c context.Context, id int64, state int32) (res string, err error) {
	if err = d.articleDB.QueryRow(c, _reasonOfVersion, id, state).Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		PromError("db:LastReason")
		log.Error("dao.LastReason(%s) error(%+v)", _reasonOfVersion, err)
	}
	return
}
