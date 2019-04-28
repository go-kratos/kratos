package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	xtime "time"

	"go-common/app/job/bbq/video/model"
	"go-common/app/service/bbq/search/api/grpc/v1"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/time"

	jsoniter "github.com/json-iterator/go"
)

const (
	_limitSize = 2000
	// MaxSyncESNum 限制每次更新到es的数量
	MaxSyncESNum = 100
	// QueryVideoByMtime 根据mtime获取视频基础信息
	QueryVideoByMtime = "select `svid`,`mtime` from video where mtime >= ? order by mtime asc"
	// QueryVideoStatisticsByMtime 根据mtime获取视频相关播放信息
	QueryVideoStatisticsByMtime = "select `svid`,`mtime` from video_statistics where mtime >= ? order by mtime asc"
	// QueryVideoStatisticsHiveByMtime 根据mtime获取视频主站信息
	QueryVideoStatisticsHiveByMtime = "select `svid`,`mtime` from video_statistics_hive where mtime >= ? order by mtime asc"
	// QueryVideoTagByMtime 根据mtime获取视频tag信息
	QueryVideoTagByMtime           = "select `svid`,`mtime` from video_tag where mtime >= ? order by mtime asc"
	_recRecallOpVideoKey           = "job:bbq:rec:op"
	_syncOperVideoTagKey           = "job:bbq:video:syncvideotagkey"
	_syncOperVideoTimeKey          = "job:bbq:video:syncvideotimekey"
	_selectVideoInfo               = "select `svid`,`title`,`content`,`mid`,`cid`,`pubtime`,`ctime`,`mtime`,`duration`,`original`,`state`,`is_full_screen`,`ver_id`,`ver`,`from`,`avid`,`tid`,`sub_tid`,`score` from video where id > ? order by id asc limit 100"
	_selectVideoInfoByIDs          = "select `svid`,`title`,`content`,`mid`,`cid`,`pubtime`,`ctime`,`mtime`,`duration`,`original`,`state`,`is_full_screen`,`ver_id`,`ver`,`from`,`avid`,`tid`,`sub_tid`,`score` from video where svid in (%s)"
	_selectVideoStatisticsHiveInfo = "select `svid`,`play`,`fav`,`coin`,`subtitles`,`likes`,`share`,`report`,`duration_daily`,`duration_all`,`reply`,`share_daily`,`play_daily`,`subtitles_daily`,`likes_daily`,`fav_daily`,`access`,`reply_daily` from video_statistics_hive where svid in (%s)"
	_selectVideoStatisticsInfo     = "select `svid`,`play`,`subtitles`,`like`,`share`,`report` from video_statistics where svid in (%s)"
	_selectVideoTagsInfo           = "select v.svid,t.id,t.name,t.type from video_tag v inner join tag t on v.tag_id = t.id where v.svid in (%s)"
	_queryCheckTask                = "select `task_id`,`task_name`,`last_check` from check_task where `task_name` = ?"
	_updateTaskLastCheck           = "update check_task set last_check = ? where `task_name` = ?"
	_queryTagByMtime               = "select `id`,`mtime` from tag where mtime > ? order by mtime asc limit 10"
	_queryVideoTagByTagID          = "select `id`,`svid` from video_tag where tag_id in (%s) and id > ? order by id asc limit 100"
	_queryVideoBySVIDs             = "select `svid`,`title` from video where svid in (%s)"
	_queryIDs                      = "select `id`,`svid` from video where id > %d order by id asc limit %d"
	_queryOutPutVideos             = "select id,svid,title,pubtime from video where id > ? and state in (%s) order by id ASC limit ?"
	_updateUVSt                    = "update user_statistics set %s = %s+1 where mid = %d"
	_updateUVStDel                 = "update user_statistics set %s = %s-1 where mid = %d and %s >0"
	_updateSVTotal                 = "update user_statistics set av_total = av_total + 1 where mid = ?"
	_getSvidByCid                  = "select svid from video where cid = ?"
	_updateSVID                    = "update video_repository set svid = ? where id = ?"
	_updateSyncStatus              = "update video_repository set sync_status = sync_status|? where svid = ?"
	_queryCMSOne                   = "select `home_img_url`,`home_img_width`,`home_img_height`,`from`,`sync_status`, `tag`, `avid`, `cid`, `svid`, `title`, `mid`, `content`, `pubtime`,`duration`,`original`,`is_full_screen`,`tid`,`sub_tid`,`cover_url`,`cover_width`,`cover_height` from video_repository where svid = ?"
	_queryCMSOneByID               = "select `tag`, `avid`, `cid`, `svid`, `title`, `mid`, `content`, `pubtime`,`duration`,`original`,`is_full_screen`,`tid`,`sub_tid`,`cover_url`,`cover_width`,`cover_height` from video_repository where id = ?"
	_queryBbqVideo                 = "select title,mid,cid,state,tid,sub_tid,svid from video where svid in (%s)"
)

// RawVideo 从数据库获取视频信息
func (d *Dao) RawVideo(ctx context.Context, SVID int64) (res *model.VideoRepRaw, err error) {
	res = new(model.VideoRepRaw)
	err = d.dbCms.QueryRow(ctx, _queryCMSOne, SVID).Scan(&res.HomeImgURL, &res.CoverWidth, &res.HomeImgHeight, &res.From, &res.SyncStatus, &res.Tag, &res.AVID, &res.CID, &res.SVID, &res.Title, &res.MID, &res.Content, &res.Pubtime, &res.Duration, &res.Original, &res.IsFull, &res.TID, &res.SubTID, &res.CoverURL, &res.CoverWidth, &res.CoverHeight)
	if err != nil {
		log.Error("RawVideo error(%v),svid:%d", err, SVID)
	}
	return
}

//RawBbqVideo ..
func (d *Dao) RawBbqVideo(ctx context.Context, SVID []int64) (res *model.VideoRaw, err error) {
	res = new(model.VideoRaw)

	svids := strings.Trim(strings.Join(strings.Split(fmt.Sprint(SVID), " "), ","), "[]")
	err = d.db.QueryRow(ctx, fmt.Sprintf(_queryBbqVideo, svids)).Scan(&res.Title, &res.MID, &res.CID, &res.State, &res.TID, &res.SubTID, &res.SVID)
	if err != nil {
		log.Errorw(ctx, "event", "RawBbqVideo queryrow scan err:[%v],SVID :[%v]", err, SVID)
	}
	return

}

//RawVideoByID ...get video info by id
func (d *Dao) RawVideoByID(ctx context.Context, ID int64) (res *model.VideoRepRaw, err error) {
	res = new(model.VideoRepRaw)
	err = d.dbCms.QueryRow(ctx, _queryCMSOneByID, ID).Scan(&res.Tag, &res.AVID, &res.CID, &res.SVID, &res.Title, &res.MID, &res.Content, &res.Pubtime, &res.Duration, &res.Original, &res.IsFull, &res.TID, &res.SubTID, &res.CoverURL, &res.CoverWidth, &res.CoverHeight)
	if err != nil {
		log.Error("RawVideoByID error(%v),id:%d", err, ID)
	}
	return
}

//UpdateSyncStatus update video_repository sync_status
func (d *Dao) UpdateSyncStatus(ctx context.Context, SVID int64, st int64) (err error) {
	if _, err = d.dbCms.Exec(ctx, _updateSyncStatus, st, SVID); err != nil {
		log.Error("UpdateSyncStatus err :%v,svid :", err, SVID)
	}
	return
}

//UpdateVideoUploadProcessStatus ...
func (d *Dao) UpdateVideoUploadProcessStatus(ctx context.Context, SVID int64, st int64) (err error) {
	if _, err = d.db.Exec(ctx, "update video_upload_process set status = ? where svid = ?", st, SVID); err != nil {
		log.Errorw(ctx, "errmsg", "UpdateVideoUploadProcessStatus update failed", "err", err)
	}
	return
}

//UpdateSvid ...
func (d *Dao) UpdateSvid(c context.Context, id int64, svid int64) (err error) {
	if _, err = d.dbCms.Exec(c, _updateSVID, svid, id); err != nil {
		log.Error("distribution svid err:%v,svid:%d", err, svid)
	}
	return
}

//AddSVTotal ...
func (d *Dao) AddSVTotal(mid int64) (err error) {
	_, err = d.db.Exec(context.Background(), _updateSVTotal, mid)
	if err != nil {
		log.Error("AddSVTotal ,mid:%s,err:%v", mid, err)
	}
	return
}

//UpdateUVSt 更新用户视频统计信息
func (d *Dao) UpdateUVSt(mid int64, field string) (err error) {
	_, err = d.db.Exec(context.Background(), fmt.Sprintf(_updateUVSt, field, field, mid))
	if err != nil {
		log.Error("UpdateUVSt ,mid:%s,field:%s,op:%d,err:%v", mid, field, err)
	}
	return
}

//UpdateUVStDel ...
func (d *Dao) UpdateUVStDel(mid int64, field string) (err error) {
	_, err = d.db.Exec(context.Background(), fmt.Sprintf(_updateUVStDel, field, field, mid, field))
	if err != nil {
		log.Error("UpdateUVStDel ,mid:%s,field:%s,op:%d,err:%v", mid, field, err)
	}
	return
}

//VideoList 获取视频基础信息
func (d *Dao) VideoList(c context.Context, id int64) (ids string, res []*v1.VideoESInfo, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _selectVideoInfo, id); err != nil {
		log.Error("select videos err(%v)", err)
		return
	}
	defer rows.Close()
	var (
		pubtime  time.Time
		ctime    time.Time
		mtime    time.Time
		idstring []string
	)
	for rows.Next() {
		tmp := new(v1.VideoESInfo)
		if err = rows.Scan(&tmp.SVID, &tmp.Title, &tmp.Content, &tmp.MID, &tmp.CID, &pubtime, &ctime, &mtime, &tmp.Duration, &tmp.Original, &tmp.State, &tmp.ISFullScreen, &tmp.VerID, &tmp.Ver, &tmp.From, &tmp.AVID, &tmp.Tid, &tmp.SubTid, &tmp.Score); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			}
			return
		}
		tmp.Pubtime = int64(pubtime)
		tmp.Ctime = int64(ctime)
		tmp.Mtime = int64(mtime)
		idstring = append(idstring, strconv.FormatInt(tmp.SVID, 10))
		res = append(res, tmp)
	}
	ids = strings.Join(idstring, ",")
	return
}

//VideoListByIDs 根据视频id获取视频基础信息
func (d *Dao) VideoListByIDs(c context.Context, ids string) (res []*v1.VideoESInfo, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_selectVideoInfoByIDs, ids)); err != nil {
		log.Error("select videos by ids err(%v)", err)
		return
	}
	defer rows.Close()
	var (
		pubtime time.Time
		ctime   time.Time
		mtime   time.Time
	)
	for rows.Next() {
		tmp := new(v1.VideoESInfo)
		if err = rows.Scan(&tmp.SVID, &tmp.Title, &tmp.Content, &tmp.MID, &tmp.CID, &pubtime, &ctime, &mtime, &tmp.Duration, &tmp.Original, &tmp.State, &tmp.ISFullScreen, &tmp.VerID, &tmp.Ver, &tmp.From, &tmp.AVID, &tmp.Tid, &tmp.SubTid, &tmp.Score); err != nil {
			log.Error("select videos scan err(%v)", err)
			return
		}
		tmp.Pubtime = int64(pubtime)
		tmp.Ctime = int64(ctime)
		tmp.Mtime = int64(mtime)
		res = append(res, tmp)
	}
	return
}

//VideoStatisticsHiveList 获取视频互动信息，hive表
func (d *Dao) VideoStatisticsHiveList(c context.Context, ids string) (res map[int64]*v1.VideoESInfo, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_selectVideoStatisticsHiveInfo, ids)); err != nil {
		log.Error("select video statistics hive err(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*v1.VideoESInfo)
	for rows.Next() {
		tmp := new(v1.VideoESInfo)
		if err = rows.Scan(&tmp.SVID, &tmp.PlayHive, &tmp.FavHive, &tmp.CoinHive, &tmp.SubtitlesHive, &tmp.LikesHive, &tmp.ShareHive, &tmp.ReportHive, &tmp.DurationDailyHive, &tmp.DurationAllHive, &tmp.ReplyHive, &tmp.ShareDailyHive, &tmp.PlayDailyHive, &tmp.SubtitlesDailyHive, &tmp.LikesDailyHive, &tmp.FavDailyHive, &tmp.AccessHive, &tmp.ReplyDailyHive); err != nil {
			log.Error("select video statistics hive scan err(%v)", err)
			return
		}
		res[tmp.SVID] = tmp
	}
	return
}

//VideoStatisticsList 获取视频互动信息
func (d *Dao) VideoStatisticsList(c context.Context, ids string) (res map[int64]*v1.VideoESInfo, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_selectVideoStatisticsInfo, ids)); err != nil {
		log.Error("select video statistics err(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*v1.VideoESInfo)
	for rows.Next() {
		tmp := new(v1.VideoESInfo)
		if err = rows.Scan(&tmp.SVID, &tmp.Play, &tmp.SubtitlesHive, &tmp.Like, &tmp.Share, &tmp.Report); err != nil {
			log.Error("select video statistics scan err(%v)", err)
			return
		}
		res[tmp.SVID] = tmp
	}
	return
}

//VideoTagsList 获取视频tags
func (d *Dao) VideoTagsList(c context.Context, ids string) (res map[int64][]*v1.VideoESTags, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_selectVideoTagsInfo, ids)); err != nil {
		fmt.Println(rows)
		log.Error("select video statistics err(%v)", err)
		return
	}
	defer rows.Close()
	var svid int64
	res = make(map[int64][]*v1.VideoESTags)
	for rows.Next() {
		tmp := new(v1.VideoESTags)
		if err = rows.Scan(&svid, &tmp.ID, &tmp.Name, &tmp.Type); err != nil {
			log.Error("select video statistics scan err(%v)", err)
			return
		}
		res[svid] = append(res[svid], tmp)
	}
	return
}

//RawCheckTask 查询脚本任务
func (d *Dao) RawCheckTask(c context.Context, taskName string) (res *model.CheckTask, err error) {
	res = new(model.CheckTask)
	raw := d.db.QueryRow(c, _queryCheckTask, taskName)
	if err = raw.Scan(&res.TaskID, &res.TaskName, &res.LastCheck); err != nil {
		log.Error("query check task name(%s) err(%v)", taskName, err)
	}
	return
}

//UpdateTaskLastCheck 更新上次执行时间
func (d *Dao) UpdateTaskLastCheck(c context.Context, taskName string, lastCheck int64) (num int64, err error) {
	res, err := d.db.Exec(c, _updateTaskLastCheck, lastCheck, taskName)
	if err != nil {
		log.Error("update task last check name(%s)(%d) err(%v)", taskName, lastCheck, err)
		return
	}
	return res.RowsAffected()
}

// RawGetIDByMtime 获取最近更新的那些svid，该函数可以用于多个表的查询，只需传入不同表的查询语句即可
func (d *Dao) RawGetIDByMtime(baseTableQuery string, mtime int64) (ids []int64, lastMtime int64, err error) {
	mtimeStr := xtime.Unix(mtime, 0).Format("2006-01-02 15:04:05")
	var rows *xsql.Rows
	if rows, err = d.db.Query(context.Background(), baseTableQuery, mtimeStr); err != nil {
		log.Error("select ids fail: err=%v, mtime=%s, sql=%s", err, mtimeStr, baseTableQuery)
		return
	}
	defer rows.Close()
	var (
		temp time.Time
		svid int64
	)
	for rows.Next() {
		if err = rows.Scan(&svid, &temp); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			}
			log.Error("select videos by mtime scan fail: err=%v, mtime=%s, sql=%s", err, mtimeStr, baseTableQuery)
			return
		}
		lastMtime = int64(temp)
		ids = append(ids, svid)
	}
	return
}

//RawTagByMtime 根据mtime获取tag信息
func (d *Dao) RawTagByMtime(c context.Context, mtime int64) (ids string, res int64, err error) {
	str := xtime.Unix(mtime, 0).Format("2006-01-02 15:04:05")
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _queryTagByMtime, str); err != nil {
		log.Error("select tag err(%v)", err)
		return
	}
	defer rows.Close()
	var (
		temp     time.Time
		id       int64
		idstring []string
	)
	for rows.Next() {
		if err = rows.Scan(&id, &temp); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			}
			log.Error("select tag by mtime scan err(%v)", err)
			return
		}
		res = int64(temp)
		idstring = append(idstring, strconv.FormatInt(id, 10))
	}
	ids = strings.Join(idstring, ",")
	return
}

//RawVideoTagByIDs .
func (d *Dao) RawVideoTagByIDs(c context.Context, ids string, id int64) (svids string, res int64, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_queryVideoTagByTagID, ids), id); err != nil {
		log.Error("select video tag err(%v)", err)
		return
	}
	defer rows.Close()
	var (
		temp     int64
		svid     int64
		idstring []string
	)
	for rows.Next() {
		if err = rows.Scan(&temp, &svid); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			} else {
				log.Error("select video tag by id scan err(%v)", err)
			}
			return
		}
		res = int64(temp)
		idstring = append(idstring, strconv.FormatInt(svid, 10))
	}
	svids = strings.Join(idstring, ",")
	return
}

//GetSyncOperVideoFlag 获取同步信号灯
func (d *Dao) GetSyncOperVideoFlag(c context.Context) (tag int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	val, err := conn.Do("get", _syncOperVideoTagKey)
	if err != nil {
		log.Error("cache sync oper video tag get err(%v)", err)
		return
	}
	if val == nil {
		if err = d.SetSyncOperVideoFlag(c, model.DenySyncOperVideoTag); err != nil {
			log.Error("set sync oper video flag faild")
			return
		}
		tag = model.DenySyncOperVideoTag
	} else {
		tag, err = redis.Int64(val, err)
	}
	return
}

//SetSyncOperVideoFlag 设置同步信号灯
func (d *Dao) SetSyncOperVideoFlag(c context.Context, v int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("set", _syncOperVideoTagKey, v); err != nil {
		log.Error("cache sync oper video tag set err(%v)", err)
		return
	}
	return
}

//GetSyncOperVideoExportTime ...
func (d *Dao) GetSyncOperVideoExportTime(c context.Context) (t string, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	val, err := conn.Do("get", _syncOperVideoTimeKey)
	if err == nil {
		if val == nil {
			conn.Do("set", _syncOperVideoTimeKey, "")
		}
		t, err = redis.String(val, err)
	} else {
		log.Error("get sync oper video export time err,errinfo:%v", err)
	}
	return
}

//RawVideoBySVIDS 根据svids获取视频
func (d *Dao) RawVideoBySVIDS(c context.Context, svids []string) (res map[int64]string, err error) {
	res = make(map[int64]string)
	str := strings.Join(svids, ",")
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_queryVideoBySVIDs, str)); err != nil {
		log.Error("select videos by svids err(%v)", err)
		return
	}
	defer rows.Close()
	var (
		svid  int64
		title string
	)
	for rows.Next() {
		if err = rows.Scan(&svid, &title); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			}
			log.Error("select videos by svids scan err(%v)", err)
			return
		}
		res[svid] = title
	}
	return
}

//GetVideoByLastID 获取所有SVID
func (d *Dao) GetVideoByLastID(c context.Context, last int64) (IDs []int64, lastRet int64, err error) {
	length := 1000 //分批大小
	var rows *xsql.Rows
	rows, err = d.db.Query(c, fmt.Sprintf(_queryIDs, last, length))
	if err != nil {
		log.Error("db _queryIDs err(%v)", err)
		return
	}
	for rows.Next() {
		var svid int64
		if err = rows.Scan(&lastRet, &svid); err != nil {
			log.Error("scan err(%v)", err)
			continue
		}
		IDs = append(IDs, svid)
	}
	return
}

//GetRecallOpVideo 获取精选视频
func (d *Dao) GetRecallOpVideo(c context.Context) (ids []int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	val, err := redis.Bytes(conn.Do("GET", _recRecallOpVideoKey))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("cache rec recall op video get redis err (%v)", err)
		}
		return ids, err
	}
	if err = jsoniter.Unmarshal(val, &ids); err != nil {
		log.Error("rec recall op video unmarshal err (%v)", err)
	}
	return
}

//SetRecallOpVideo 写入精选视频
func (d *Dao) SetRecallOpVideo(c context.Context, ids []int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	bytes, _ := jsoniter.Marshal(ids)
	_, err = conn.Do("SET", _recRecallOpVideoKey, bytes)
	if err != nil {
		log.Error("rec recall op video set redis error(%v) ", err)
	}
	return
}

//VideosByLast 使用lastid批量获取视屏
func (d *Dao) VideosByLast(c context.Context, lastid int64) (svinfo []*model.VideoDB, err error) {
	var rows *xsql.Rows
	query := fmt.Sprintf(_queryOutPutVideos, model.VideoStateOutPut)
	rows, err = d.db.Query(c, query, lastid, _limitSize)
	if err != nil {
		log.Error("db _queryVideos err(%v)", err)
		return
	}
	for rows.Next() {
		video := new(model.VideoDB)
		if err = rows.Scan(&video.AutoID, &video.ID, &video.Title, &video.Pubtime); err != nil {
			log.Error("scan err(%v)", err)
			continue
		}
		svinfo = append(svinfo, video)
	}
	return
}

// GetSvidByCid 根据cid获取svid
func (d *Dao) GetSvidByCid(c context.Context, cid int64) (svid int64, err error) {
	err = d.db.QueryRow(c, _getSvidByCid, cid).Scan(&svid)
	if err != nil {
		log.Warn("db _getSvidByCid err(%v)", err)
		return
	}
	return
}
