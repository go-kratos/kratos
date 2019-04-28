package dao

import (
	"context"
	"net/url"
	"strconv"
	"time"

	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_verifyAPI = "http://account.bilibili.co/api/internal/identify/info"
)

// Authors loads author list who are permitted to post articles.
func (d *Dao) Authors(c context.Context) (res map[int64]*artmdl.AuthorLimit, err error) {
	var rows *sql.Rows
	if rows, err = d.authorsStmt.Query(c); err != nil {
		PromError("db:作者列表查询")
		log.Error("db.authorsStmt.Query error(%+v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*artmdl.AuthorLimit)
	for rows.Next() {
		var (
			mid    int64
			author = &artmdl.AuthorLimit{}
		)
		if err = rows.Scan(&mid, &author.Limit, &author.State); err != nil {
			PromError("作者列表scan")
			log.Error("rows.Authors.Scan error(%+v)", err)
			return
		}
		res[mid] = author
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// ApplyCount get today apply count
func (d *Dao) ApplyCount(c context.Context) (count int64, err error) {
	var t = time.Now().Truncate(24 * time.Hour)
	if err = d.applyCountStmt.QueryRow(c, t).Scan(&count); err != nil {
		PromError("db:查询申请总数")
		log.Error("db.ApplyCountStmt.Query(%v) error(%+v)", t, err)
	}
	return
}

// AddApply add new apply
func (d *Dao) AddApply(c context.Context, mid int64, content, category string) (err error) {
	var t = time.Now()
	if _, err = d.applyStmt.Exec(c, mid, t, content, category, t, content, category); err != nil {
		PromError("db:申请作者权限")
		log.Error("db.applyStmt.Query(mid: %v, t: %v, category: %v) error(%+v)", mid, t, category, err)
	}
	return
}

// AddAuthor add author
func (d *Dao) AddAuthor(c context.Context, mid int64) (err error) {
	if _, err = d.addAuthorStmt.Exec(c, mid); err != nil {
		PromError("db:增加作者权限")
		log.Error("mysql: db.addAuthorStmt.Query(mid: %v) error(%+v)", mid, err)
	}
	return
}

// RawAuthor get author's info.
func (d *Dao) RawAuthor(c context.Context, mid int64) (res *artmdl.AuthorLimit, err error) {
	res = new(artmdl.AuthorLimit)
	if err = d.authorStmt.QueryRow(c, mid).Scan(&res.State, &res.Rtime, &res.Limit); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		PromError("db:RawAuthor scan")
		log.Error("RawAuthor row.Scan(%d) error(%+v)", mid, err)
	}
	if int64(res.Rtime) < 0 {
		res.Rtime = xtime.Time(0)
	}
	return
}

// Identify gets user verify info.
func (d *Dao) Identify(c context.Context, mid int64) (res *artmdl.Identify, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var resp struct {
		Code int              `json:"code"`
		Data *artmdl.Identify `json:"data"`
	}
	if err = d.httpClient.Get(c, _verifyAPI, "", params, &resp); err != nil {
		log.Error("d.httpClient.Get(%s) error(%+v)", _verifyAPI+"?"+params.Encode(), err)
		PromError("http:获取用户实名认证信息")
		return
	}
	if resp.Code != ecode.OK.Code() {
		log.Error("d.httpClient.Get(%s) code(%d)", _verifyAPI+"?"+params.Encode(), resp.Code)
		PromError("http:获取用户实名认证信息状态码异常")
		err = ecode.Int(resp.Code)
		return
	}
	res = resp.Data
	return
}
