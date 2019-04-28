package dao

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"go-common/app/job/bbq/recall/internal/model"
	"go-common/app/job/bbq/recall/proto"
	"go-common/app/job/bbq/recall/proto/quality"
	"go-common/library/log"

	"github.com/golang/snappy"
)

const (
	// _fetchVideo    = "select `id`, `title`, `content`, `mid`, `avid`, `cid`, `pubtime`, `ctime`, `mtime`, `duration`, `state`, `tid`, `sub_tid` from video where pubtime > ? limit ?, ?;"
	_fetchVideo          = "select `svid`, `title`, `content`, `mid`, `avid`, `cid`, `pubtime`, `ctime`, `mtime`, `duration`, `state`, `tid`, `sub_tid` from video limit ?, ?;"
	_fetchVideoTag       = "select `id`, `name`, `type` from `tag` where `id` = ? and `status` = 1;"
	_fetchVideoTagAll    = "select `id`, `name`, `type` from `tag` where `status` = 1;"
	_fetchVideoTextTag   = "select `tag` from `video_repository` where `svid` = ? limit 1;"
	_queryVideoQuality   = "select `stat_info` from `video_forward_index_stat_info` where `svid` = ?;"
	_fetchNewIncomeVideo = "select `svid` from `video` where ctime > ? and state in (%s);"
)

// FetchVideoInfo .
func (d *Dao) FetchVideoInfo(c context.Context, offset, size int) (result []*model.Video, err error) {
	// rows, err := d.db.Query(c, _fetchVideo, ptime.Format("2006-01-02 15:04:05"), offset, size)
	rows, err := d.db.Query(c, _fetchVideo, offset, size)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		tmp := &model.Video{}
		if err = rows.Scan(&tmp.SVID, &tmp.Title, &tmp.Content, &tmp.MID, &tmp.AVID, &tmp.CID, &tmp.PubTime, &tmp.CTime, &tmp.MTime, &tmp.Duration, &tmp.State, &tmp.TID, &tmp.SubTID); err != nil {
			log.Error("FetchVideoInfo: %v", err)
			return nil, err
		}
		result = append(result, tmp)
	}

	return result, nil
}

// FetchVideoTagAll .
func (d *Dao) FetchVideoTagAll(c context.Context) (result []*proto.Tag, err error) {
	result = make([]*proto.Tag, 0)
	rows, err := d.db.Query(c, _fetchVideoTagAll)
	if err != nil {
		return
	}

	for rows.Next() {
		tmp := new(proto.Tag)
		if err = rows.Scan(&tmp.TagID, &tmp.TagName, &tmp.TagType); err != nil {
			log.Error("FetchVideoTag: %v", err)
			continue
		}
		result = append(result, tmp)
	}

	return
}

// FetchVideoTag .
func (d *Dao) FetchVideoTag(c context.Context, tid int32) (result *proto.Tag, err error) {
	row := d.db.QueryRow(c, _fetchVideoTag, tid)

	result = new(proto.Tag)
	if err = row.Scan(&result.TagID, &result.TagName, &result.TagType); err != nil {
		log.Error("FetchVideoTag: %v", err)
		return
	}

	return
}

// FetchVideoTextTag .
func (d *Dao) FetchVideoTextTag(c context.Context, svid int64) (result []string, err error) {
	row := d.dbCms.QueryRow(c, _fetchVideoTextTag, svid)
	var tags string
	if err = row.Scan(&tags); err != nil {
		log.Errorv(c, log.KV("log", "_fetchVideoTextTag failed"), log.KV("error", err), log.KV("svid", svid))
		return
	}

	result = strings.Split(tags, ",")
	return
}

// FetchVideoQuality .
func (d *Dao) FetchVideoQuality(c context.Context, svid uint64) (result *quality.VideoQuality, err error) {
	var raw string
	row := d.dbOffline.QueryRow(c, _queryVideoQuality, svid)
	row.Scan(&raw)
	if raw == "" {
		return
	}

	trimed := strings.Trim(raw, "\n")
	hexDst, err := hex.DecodeString(trimed)
	if err != nil {
		log.Error("FetchVideoQuality: %v src[%s] raw[%s]", err, trimed, trimed)
		return
	}

	snappyDst, err := snappy.Decode(nil, hexDst)
	if err != nil {
		log.Error("FetchVideoQuality: %v src[%s] raw[%s]", err, string(hexDst), trimed)
		return
	}

	result = &quality.VideoQuality{}
	result.Unmarshal(snappyDst)
	if err != nil {
		log.Error("FetchVideoQuality: %v src[%s] raw[%s]", err, snappyDst, trimed)
	}
	return
}

// FetchNewincomeVideo .
func (d *Dao) FetchNewincomeVideo() (res []int64, err error) {
	duration, _ := time.ParseDuration("-24h")
	today := time.Now().Add(duration).Format("2006-01-02")
	_query := fmt.Sprintf(_fetchNewIncomeVideo, strings.Join(model.RecommendVideoState, ","))
	row, err := d.db.Query(context.Background(), _query, today)
	if err != nil {
		return
	}
	res = make([]int64, 0)
	for row.Next() {
		var tmp int64
		if err = row.Scan(&tmp); err != nil {
			return
		}
		res = append(res, tmp)
	}

	return
}
