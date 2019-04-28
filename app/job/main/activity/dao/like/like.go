package like

import (
	"context"
	"database/sql"
	"net/url"
	"strconv"

	"go-common/app/job/main/activity/model/like"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_selLikeSQL        = "SELECT id,wid FROM likes WHERE state=1 AND sid=? ORDER BY type"
	_likeListSQL       = "SELECT id,wid FROM likes WHERE state= 1 AND sid = ? ORDER BY id LIMIT ?,?"
	_likesCntSQL       = "SELECT COUNT(1) AS cnt FROM likes WHERE state = 1 AND sid = ?"
	_setObjStatURI     = "/x/internal/activity/object/stat/set"
	_setViewRankURI    = "/x/internal/activity/view/rank/set"
	_setLikeContentURI = "/x/internal/activity/like/content/set"
)

// Like get like by sid
func (d *Dao) Like(c context.Context, sid int64) (ns []*like.Like, err error) {
	rows, err := d.db.Query(c, _selLikeSQL, sid)
	if err != nil {
		log.Error("notice.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		n := &like.Like{}
		if err = rows.Scan(&n.ID, &n.Wid); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		ns = append(ns, n)
	}
	return
}

// LikeList get like list by sid.
func (d *Dao) LikeList(c context.Context, sid int64, offset, limit int) (list []*like.Like, err error) {
	rows, err := d.db.Query(c, _likeListSQL, sid, offset, limit)
	if err != nil {
		err = errors.Wrapf(err, "LikeList:d.db.Query(%d,%d,%d)", sid, offset, limit)
		return
	}
	defer rows.Close()
	for rows.Next() {
		n := new(like.Like)
		if err = rows.Scan(&n.ID, &n.Wid); err != nil {
			err = errors.Wrapf(err, "LikeList:row.Scan row (%d,%d,%d)", sid, offset, limit)
			return
		}
		list = append(list, n)
	}
	if err = rows.Err(); err != nil {
		err = errors.Wrapf(err, "LikeList:rowsErr(%d,%d,%d)", sid, offset, limit)
	}
	return
}

// LikeCnt get like list total count by sid.
func (d *Dao) LikeCnt(c context.Context, sid int64) (count int, err error) {
	row := d.db.QueryRow(c, _likesCntSQL, sid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrapf(err, "LikeCnt:QueryRow(%d)", sid)
		}
	}
	return
}

// SetObjectStat .
func (d *Dao) SetObjectStat(c context.Context, sid int64, stat *like.SubjectTotalStat, count int) (err error) {
	params := url.Values{}
	params.Set("sid", strconv.FormatInt(sid, 10))
	params.Set("like", strconv.FormatInt(stat.SumLike, 10))
	params.Set("view", strconv.FormatInt(stat.SumView, 10))
	params.Set("fav", strconv.FormatInt(stat.SumFav, 10))
	params.Set("coin", strconv.FormatInt(stat.SumCoin, 10))
	params.Set("count", strconv.Itoa(count))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpClient.Get(c, d.setObjStatURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		err = errors.Wrapf(err, "SetObjectStat(%d,%v)", sid, stat)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrapf(ecode.Int(res.Code), "SetObjectStat Code(%d,%v)", sid, stat)
	}
	return
}

// SetViewRank set view rank list.
func (d *Dao) SetViewRank(c context.Context, sid int64, aids []int64) (err error) {
	params := url.Values{}
	params.Set("sid", strconv.FormatInt(sid, 10))
	params.Set("aids", xstr.JoinInts(aids))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpClient.Get(c, d.setViewRankURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		err = errors.Wrapf(err, "SetViewRank(%d,%v)", sid, aids)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrapf(ecode.Int(res.Code), "SetViewRank Code(%d,%v)", sid, aids)
	}
	return
}

// SetLikeContent .
func (d *Dao) SetLikeContent(c context.Context, id int64) (err error) {
	var res struct {
		Code int `json:"code"`
	}
	params := url.Values{}
	params.Set("lid", strconv.FormatInt(id, 10))
	if err = d.httpClient.Get(c, d.setLikeContentURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		err = errors.Wrapf(err, "SetLikeContent(%d)", id)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrapf(ecode.Int(res.Code), "SetLikeContent Code(%d)", id)
	}
	return
}
