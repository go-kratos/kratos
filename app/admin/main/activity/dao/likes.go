package dao

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go-common/app/admin/main/activity/model"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	_likeBatchSQL        = "INSERT INTO likes(`sid`,`wid`,`mid`,`type`,`state`,`stick_top`,`ctime`,`mtime`) VALUES %s"
	_likeContentBatchSQL = "INSERT INTO like_content(`id`,`ipv6`,`ctime`,`mtime`) VALUES %s"
)

func imgAddKey(uri string) (url string) {
	if strings.Contains(uri, "http://drawyoo.hdslb.com") {
		path := strings.Replace(uri, "http://drawyoo.hdslb.com", "", -1)
		expire := time.Now().Unix() + 3600
		md5Byte := md5.Sum([]byte(fmt.Sprintf("rjZOPr8w%s%d", path, expire)))
		md5Key := md5Byte[:]
		afterKey := strings.Replace(base64.StdEncoding.EncodeToString(md5Key), "+/", "-_", -1)
		return uri + "?key=" + strings.Replace(afterKey, "=", "", -1) + "&expires=" + fmt.Sprintf("%d", expire)
	}
	return uri
}

// GetLikeContent .
func (d *Dao) GetLikeContent(c context.Context, ids []int64) (outRes map[int64]*model.LikeContent, err error) {
	var likeContent []*model.LikeContent
	if err = d.DB.Where("id in (?)", ids).Find(&likeContent).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("s.DB.Where(id in %v).Find(),error(%v)", ids, err)
		return
	}
	outRes = make(map[int64]*model.LikeContent, len(likeContent))
	for _, item := range likeContent {
		outRes[item.ID] = item
		outRes[item.ID].Image = imgAddKey(item.Image)
	}
	return
}

// ActSubject get likesubject from db.
func (d *Dao) ActSubject(c context.Context, sid int64) (rp *model.ActSubject, err error) {
	rp = new(model.ActSubject)
	if err = d.DB.Where("id = ?", sid).First(rp).Error; err != nil {
		log.Error(" s.DB.Where(id ,%d).First() error(%v)", sid, err)
	}
	return
}

// Musics get music info .
func (d *Dao) Musics(c context.Context, aids []int64, ip string) (music *model.MusicRes, err error) {
	params := url.Values{}
	params.Set("songIds", xstr.JoinInts(aids))
	if err = d.client.Post(c, d.songsURL, ip, params, &music); err != nil {
		err = errors.Wrapf(err, "d.client.Post(%s)", d.songsURL)
	}
	if music.Code != 0 {
		err = errors.New("get music error")
	}
	return
}

// BatchLike .
func (d *Dao) BatchLike(c context.Context, item *model.Like, wids []int64, ipv6 []byte) (err error) {
	var (
		likesVal []*model.Like
	)
	if len(wids) == 0 {
		return
	}
	lidString := make([]string, 0, len(wids))
	rowStrings := make([]string, 0, len(wids))
	ctime := time.Now()
	for _, v := range wids {
		rowStrings = append(rowStrings, fmt.Sprintf("(%d,%d,%d,%d,%d,%d,'%s','%s')", item.Sid, v, item.Mid, item.Type, item.State, item.StickTop, ctime, ctime))
	}
	tx := d.DB.Begin()
	if err = tx.Model(&model.Like{}).Exec(fmt.Sprintf(_likeBatchSQL, strings.Join(rowStrings, ","))).Error; err != nil {
		err = errors.Wrapf(err, " d.DB.Model(&model.Like{}).Exec(%s)", _likeBatchSQL)
		tx.Rollback()
		return
	}
	if err = tx.Model(&model.Like{}).Where(fmt.Sprintf("sid = ? and wid in (%s)", xstr.JoinInts(wids)), item.Sid).Find(&likesVal).Error; err != nil {
		err = errors.Wrapf(err, " d.DB.Model(&model.Like{}).find()")
		tx.Rollback()
		return
	}
	for _, itm := range likesVal {
		lidString = append(lidString, fmt.Sprintf("(%d,'%s','%s','%s')", itm.ID, ipv6, ctime, ctime))
	}
	if err = tx.Model(&model.LikeContent{}).Exec(fmt.Sprintf(_likeContentBatchSQL, strings.Join(lidString, ","))).Error; err != nil {
		err = errors.Wrapf(err, " d.DB.Model(&model.LikeContent{}).Exec(%s)", _likeContentBatchSQL)
		tx.Rollback()
		return
	}
	err = tx.Commit().Error
	return
}
