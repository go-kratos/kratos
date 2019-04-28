package dao

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-common/app/service/bbq/video/api/grpc/v1"
	"go-common/app/service/bbq/video/model"
	acc "go-common/app/service/main/account/api"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
	xhttp "net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	_BVCSubTableSize     = 100
	_queryVideo          = "SELECT svid FROM video WHERE svid = ?"
	_addVideo            = "INSERT INTO video(`cover_url`,`cover_width`,`cover_height`,`svid`,`title`,`mid`,`avid`,`cid`,`pubtime`,`from`,`tid`,`sub_tid`,`home_img_url`,`home_img_width`,`home_img_height`,`state`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_queryTagByName      = "SELECT `id` FROM tag WHERE name = ? and type = ?"
	_insertTag           = "INSERT INTO tag (`name`,`type`,`status`) VALUES %s "
	_insOrUpUserBase     = "INSERT IGNORE user_base (mid, uname, face, user_type) VALUES (?, ?, ?, ?) "
	_insOrUpUserSta      = "INSERT IGNORE user_statistics_hive (mid, uname) VALUES (?, ?)"
	_queryStatisticsList = "select `svid`, `play`, `subtitles`, `like`, `share`, `report` from video_statistics where svid in (%s)"
	_addBVCData          = "insert into %s (`svid`,`path`,`resolution_retio`,`code_rate`,`video_code`,`duration`,`file_size`) values (?,?,?,?,?,?,?)"
	_updateBVCData       = "update %s set path=?, resolution_retio=?, video_code=?, duration=?, file_size=? where svid = ? and code_rate = ?"
	_updateSvPIC         = "update video_repository set cover_url=?,cover_width=?,cover_height=? ,sync_status = sync_status|? where svid = ?"
	_addVideoViews       = "update `video_statistics` set `play` = `play` + ? where `svid` = ?"
	_existedStatistics   = "select `id` from `video_statistics` where `svid` = ?;"
	_insertStatistics    = "insert into `video_statistics`(`svid`, `play`, `subtitles`, `like`, `share`, `report`) values(?,?,?,?,?,?);"
	_queryVideoList      = "select `avid`, `cid`, `svid`, `title`, `mid`, `content`, `pubtime`,`duration`,`tid`,`sub_tid`,`cover_url`,`cover_width`,`cover_height`,`limits`, `state` from video where svid in (%s)"
	_updateVideoState    = "update `video` set `state` = ? where `svid`= ?;"
)

const (
	videoBaseCacheExpire = 600
	videoBaseCacheKey    = "video_base:%d"
)

func keyVideoBase(svid int64) string {
	return fmt.Sprintf(videoBaseCacheKey, svid)
}

// ModifyLimits .
func (d *Dao) ModifyLimits(c context.Context, svid int64, limitType uint64, limitOp uint64) (num int64, err error) {
	// 根据操作选择合适的limits update语句
	limitOpCond := fmt.Sprintf("|%d", 1<<limitType)
	if limitOp == 0 {
		limitOpCond = fmt.Sprintf("&~%d", 1<<limitType)
	}
	querySQL := fmt.Sprintf("update video set limits = limits%s where svid = %d", limitOpCond, svid)
	res, err := d.db.Exec(c, querySQL)
	if err != nil {
		log.Warnw(c, "log", "modify video limits fail", "sql", querySQL)
		return
	}
	num, _ = res.RowsAffected()
	log.V(1).Infow(c, "sql", querySQL, "affected_num", num)
	d.DelCacheVideoBase(c, svid)
	return
}

// RawVideoBase mysql获取video_base
func (d *Dao) RawVideoBase(c context.Context, svids []int64) (res map[int64]*v1.VideoBase, err error) {
	res = make(map[int64]*v1.VideoBase)
	if len(svids) == 0 {
		return
	}
	querySQL := fmt.Sprintf(_queryVideoList, xstr.JoinInts(svids))
	rows, err := d.db.Query(c, querySQL)
	if err != nil {
		log.Errorw(c, "log", "get video base from mysql fail", "sql", querySQL, "err", err)
		return
	}
	defer rows.Close()
	log.V(1).Infow(c, "log", "raw get video base from mysql", "sql", querySQL)

	for rows.Next() {
		sv := new(v1.VideoBase)
		if err = rows.Scan(&sv.Avid, &sv.Cid, &sv.Svid, &sv.Title, &sv.Mid, &sv.Content, &sv.Pubtime, &sv.Duration, &sv.Tid, &sv.SubTid, &sv.CoverUrl, &sv.CoverWidth, &sv.CoverHeight, &sv.Limits, &sv.State); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		res[sv.Svid] = sv
	}
	if len(svids) > len(res) {
		var rspID []int64
		for k := range res {
			rspID = append(rspID, k)
		}
		log.Warnw(c, "log", fmt.Sprintf("video req and rsp size not equal: req=%v, rsp=%v", svids, rspID))
	}
	log.V(1).Infow(c, "req_size", len(svids), "rsp_size", len(res))

	return
}

// CacheVideoBase cache video base
func (d *Dao) CacheVideoBase(c context.Context, svids []int64) (res map[int64]*v1.VideoBase, err error) {
	res = make(map[int64]*v1.VideoBase)
	keys := make([]string, 0, len(svids))
	keyMidMap := make(map[int64]bool, len(svids))
	for _, svid := range svids {
		key := keyVideoBase(svid)
		if _, exist := keyMidMap[svid]; !exist {
			// duplicate svid
			keyMidMap[svid] = true
			keys = append(keys, key)
		}
	}

	conn := d.redis.Get(c)
	defer conn.Close()
	for _, key := range keys {
		conn.Send("GET", key)
	}
	conn.Flush()
	var data []byte
	for i := 0; i < len(keys); i++ {
		if data, err = redis.Bytes(conn.Receive()); err != nil {
			if err == redis.ErrNil {
				err = nil
			} else {
				log.Errorv(c, log.KV("event", "redis_get"), log.KV("key", keys[i]))
			}
			continue
		}
		baseItem := new(v1.VideoBase)
		json.Unmarshal(data, baseItem)
		res[baseItem.Svid] = baseItem
	}
	log.Infov(c, log.KV("event", "redis_get"), log.KV("row_num", len(res)))

	return
}

// AddCacheVideoBase 添加缓存
func (d *Dao) AddCacheVideoBase(c context.Context, videoBases map[int64]*v1.VideoBase) (err error) {
	keyValueMap := make(map[string][]byte, len(videoBases))
	for mid, videoBase := range videoBases {
		key := keyVideoBase(mid)
		if _, exist := keyValueMap[key]; !exist {
			data, _ := json.Marshal(videoBase)
			keyValueMap[key] = data
		}
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	for key, value := range keyValueMap {
		conn.Send("SET", key, value, "EX", videoBaseCacheExpire)
	}
	conn.Flush()
	for i := 0; i < len(keyValueMap); i++ {
		conn.Receive()
	}
	log.Infov(c, log.KV("event", "redis_set"), log.KV("row_num", len(videoBases)))
	return
}

// DelCacheVideoBase 删除缓存
func (d *Dao) DelCacheVideoBase(c context.Context, svid int64) {
	var key = keyVideoBase(svid)
	conn := d.redis.Get(c)
	defer conn.Close()
	conn.Do("DEL", key)
}

// AddOrUpdateVideo 添加或更新视频记录
func (d *Dao) AddOrUpdateVideo(c context.Context, vh *v1.ImportVideoInfo) (err error) {
	var (
		svid int64
	)
	tx, err := d.BeginTran(c)
	if err != nil {
		log.Error("begin transaction err :%v", err)
		return
	}
	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.Error("tx.Rollback() error(%v)", err)
			}
		} else {
			if err = tx.Commit(); err != nil {
				log.Error("tx.Commit() error(%v)", err)
			}
		}
	}()
	p := &model.VideoInfo{
		CoverURL:      vh.CoverUrl,
		CoverWidth:    vh.CoverWidth,
		CoverHeight:   vh.CoverHeight,
		SVID:          vh.Svid,
		Title:         vh.Title,
		MID:           vh.MID,
		AVID:          vh.AVID,
		CID:           vh.CID,
		Pubtime:       vh.Pubtime,
		From:          int16(vh.From),
		State:         int16(vh.State),
		TID:           vh.TID,
		SubTID:        vh.SubTID,
		HomeImgURL:    vh.HomeImgUrl,
		HomeImgWidth:  vh.HomeImgWidth,
		HomeImgHeight: vh.HomeImgHeight,
	}
	if err = tx.QueryRow(_queryVideo, vh.Svid).Scan(&svid); err == sql.ErrNoRows {
		if err = d.txInsertVideo(c, tx, p); err != nil {
			log.Warn("insert video err:%v,svid:%v", err, vh.Svid)
			return
		}
	} else if err != nil {
		log.Error("video queryrow scan err:[%v], svid[%v]", err, vh.Svid)
		return
	}
	//sync video_upload_process status
	if err = d.txUpdateVideoUploadProcessStatus(c, tx, vh.Svid, model.VideoUploadProcessStatusSuccessed); err != nil {
		log.Errorw(c, "event", "d.UpdateVideoUploadProcessStatus err", "err", err)
	}
	return
}

//UpdateVideoUploadProcessStatus ...
func (d *Dao) txUpdateVideoUploadProcessStatus(ctx context.Context, tx *xsql.Tx, SVID int64, st int64) (err error) {
	if _, err = tx.Exec("update video_upload_process set upload_status = ? where svid = ?", st, SVID); err != nil {
		log.Errorw(ctx, "errmsg", "UpdateVideoUploadProcessStatus update failed", "err", err)
	}
	return
}

//txInsertVideo insert video
func (d *Dao) txInsertVideo(c context.Context, tx *xsql.Tx, vh *model.VideoInfo) (err error) {
	if _, err = tx.Exec(_addVideo,
		vh.CoverURL,
		vh.CoverWidth,
		vh.CoverHeight,
		vh.SVID,
		vh.Title,
		vh.MID,
		vh.AVID,
		vh.CID,
		vh.Pubtime,
		vh.From,
		vh.TID,
		vh.SubTID,
		vh.HomeImgURL,
		vh.HomeImgWidth,
		vh.HomeImgHeight,
		vh.State,
	); err != nil {
		log.Errorw(c, "event", "insert video err", "err", err, "param", vh)
		return
	}
	return
}

// AddOrUpdateTag 更新或添加标签
func (d *Dao) AddOrUpdateTag(c context.Context, tmap []*v1.TagInfo) (tids []int64, err error) {
	// 检查已存在的tag
	for _, v := range tmap {
		row := d.db.QueryRow(c, _queryTagByName, v.TagName, v.TagType)
		t := &model.Tag{
			Type: v.TagType,
			Name: v.TagName,
		}
		err = row.Scan(&t.ID)
		if err == sql.ErrNoRows {
			var q string
			var id int64
			var res sql.Result
			n := strings.Replace(t.Name, "'", "\\'", -1)
			q = "('" + n + "'," + strconv.FormatInt(int64(t.Type), 10) + ",1)"
			res, _ = d.db.Exec(c, fmt.Sprintf(_insertTag, q))
			if res != nil {
				id, err = res.LastInsertId()
			}
			if id != 0 {
				tids = append(tids, id)
			}
		} else if t.ID != 0 {
			tids = append(tids, t.ID)
		} else {
			log.Error("d.db.QueryRow[%v],err:%v", v.TagName, err)
			return
		}
	}
	return
}

//根据mids批量查询用户基本信息
func (d *Dao) getUserInfos(c context.Context, mids []int64) (userBases []*model.UserBase, err error) {
	midsReq := &acc.MidsReq{
		Mids:   mids,
		RealIp: metadata.String(c, metadata.RemoteIP)}
	infosReply, err := d.AccountClient.Infos3(c, midsReq)
	if infosReply == nil {
		log.Error("query infos3 failed, err (%v)", err)
		return
	}
	userBases = make([]*model.UserBase, 0, 50)
	for _, info := range infosReply.Infos {
		if info.Mid != 0 {
			if len(info.Face) > 255 {
				info.Face = "http://i0.hdslb.com/bfs/bbq/video-image/userface/1558868601542006937.png"
				log.Info("the value of Face is too long, replace it as http://i0.hdslb.com/bfs/bbq/video-image/userface/1558868601542006937.png, mid(%v)", info.Mid)
			}
			userBase := &model.UserBase{
				Mid:  info.Mid,
				Name: info.Name,
				Sex:  info.Sex,
				Face: info.Face,
				Sign: info.Sign,
				Rank: info.Rank,
			}
			userBases = append(userBases, userBase)
		}
	}
	return
}

//根据mid查询用户基本信息
func (d *Dao) getUserInfo(c context.Context, mid int64) (userBase *model.UserBase, err error) {
	midReq := &acc.MidReq{
		Mid:    mid,
		RealIp: metadata.String(c, metadata.RemoteIP)}
	info, err := d.AccountClient.Info3(c, midReq)
	if err != nil {
		log.Error("query info3 failed,mid(%v), err(%v)", mid, err)
		return
	}
	if len(info.Info.Face) > 255 {
		info.Info.Face = "http://i0.hdslb.com/bfs/bbq/video-image/userface/1558868601542006937.png"
		log.Info("the value of Face is too long, replace it as http://i0.hdslb.com/bfs/bbq/video-image/userface/1558868601542006937.png,,mid(%v)", mid)
	}
	userBase = &model.UserBase{
		Mid:  info.Info.Mid,
		Name: info.Info.Name,
		Sex:  info.Info.Sex,
		Face: info.Info.Face,
		Sign: info.Info.Sign,
		Rank: info.Info.Rank,
	}
	log.Info("getUserInfo userbase (%v)", userBase)
	return
}

//InOrUpUserBase 更新用户基本信息
func (d *Dao) InOrUpUserBase(c context.Context, mid int64) (response *v1.SyncUserBaseResponse, err error) {
	var (
		retry = 3
		try   int
		tx    *xsql.Tx
		res   sql.Result
	)
	userBase, _ := d.getUserInfo(c, mid)
	response = &v1.SyncUserBaseResponse{Affc: -1}
	for try = 0; try <= retry; try++ {
		if tx, err = d.BeginTran(c); err != nil {
			time.Sleep(time.Duration(try) * time.Second)
			log.Warn("InOrUpUserBase  try  begin transaction failed ,err(%v)", err)
			continue
		}
		if res, err = tx.Exec(
			_insOrUpUserBase,
			userBase.Mid,
			userBase.Name,
			userBase.Face,
		); err != nil {
			if err = tx.Rollback(); err != nil {
				log.Warn("InOrUpUserBase  try rollback failed ,error(%v)", err)
			}
		} else {
			if err = tx.Commit(); err != nil {
				log.Warn("InOrUpUserBase try commit failed , error(%v)", err)
			} else {
				//提交成功,退出
				response.Affc, _ = res.RowsAffected()
				log.Info("InOrUpUserBase success, affected %v rows", response.Affc)
				break
			}
		}
	}
	if err != nil {
		log.Error("InOrUpUserBase failed, mid(%v), err(%v)", mid, err)
	}
	return
}

//InOrUpUserBases 批量更新用户基本信息
func (d *Dao) InOrUpUserBases(c context.Context, mids []int64) (response *v1.SyncUserBaseResponse, err error) {
	var (
		retry = 3
		try   int
		tx    *xsql.Tx
		res   sql.Result
	)
	userBases, _ := d.getUserInfos(c, mids)
	response = &v1.SyncUserBaseResponse{Affc: -1}
	for try = 0; try <= retry; try++ {
		if tx, err = d.BeginTran(c); err != nil {
			time.Sleep(time.Duration(try) * time.Second)
			log.Warn("InOrUpUserBases try  begin transaction failed failed ,error(%v)", err)
			continue
		}
		sql := "INSERT INTO user_base (mid, uname, face, user_type) VALUES "
		for _, userBase := range userBases {
			if userBase.Mid != 0 {
				sql = sql + "(" + strconv.FormatInt(userBase.Mid, 10) + ",'" + userBase.Name + "','" + userBase.Face + "', 1),"
			}
		}
		if sql == "INSERT INTO user_base (mid, uname, face) VALUES " {
			response.Affc = 0
			log.Info("InOrUpUserBases param mids are not exist")
			return
		}
		sql = sql[0:len(sql)-1] + " ON DUPLICATE KEY UPDATE uname=values(uname), face=values(face);"
		if res, err = tx.Exec(sql); err != nil {
			log.Info("InOrUpUserBases sql = (%s)", sql)
			if err = tx.Rollback(); err != nil {
				log.Warn("InOrUpUserBases  try rollback failed ,error(%v)", err)
			}
		} else {
			log.Info("InOrUpUserBases sql = (%s)", sql)
			if err = tx.Commit(); err != nil {
				log.Warn("InOrUpUserBases try commit failed , error(%v)", err)
			} else {
				//提交成功,退出
				response.Affc, _ = res.RowsAffected()
				log.Info("InOrUpUserBases commit success, affected %v rows", response.Affc)
				break
			}
		}
	}
	if err != nil {
		log.Error("InOrUpUserBases failed, err(%v)", err)
	}
	return
}

//InOrUpUserSta 更新用户up主主站画像
func (d *Dao) InOrUpUserSta(c context.Context, mid int64) (response *v1.SyncUserBaseResponse, err error) {
	var (
		retry = 3
		try   int
		tx    *xsql.Tx
		res   sql.Result
	)
	log.Info("InOrUpUserSta start")
	response = &v1.SyncUserBaseResponse{Affc: -1}
	userBase, _ := d.getUserInfo(c, mid)
	for try = 0; try <= retry; try++ {
		if tx, err = d.BeginTran(c); err != nil {
			time.Sleep(time.Duration(try) * time.Second)
			log.Info("InOrUpUserSta on mid(%v) try begin transaction failed failed ,error(%v)", userBase.Mid, err)
			continue
		}
		if res, err = tx.Exec(
			_insOrUpUserSta,
			userBase.Mid,
			userBase.Name,
		); err != nil {
			fmt.Printf("sql exec error,err(%v)", err)
			if err = tx.Rollback(); err != nil {
				log.Info("InOrUpUserSta on mid(%v) rollback failed ,error(%v)", userBase.Mid, err)
			} else {
				fmt.Println("rollbacked")
			}
		} else {
			if err = tx.Commit(); err != nil {
				log.Info("InOrUpUserSta on mid(%v) commit failed , error(%v)", userBase.Mid, err)
			} else {
				//提交成功,退出
				response.Affc, _ = res.RowsAffected()
				break
			}
		}
	}
	if err != nil {
		log.Error("InOrUpUserSta mid(%v) failed, err(%v)", mid, err)
	}
	return
}

//InOrUpUserStas 批量更新用户状态
func (d *Dao) InOrUpUserStas(c context.Context, mids []int64) (response *v1.SyncUserBaseResponse, err error) {
	var (
		retry = 3
		try   int
		tx    *xsql.Tx
		res   sql.Result
	)
	log.Info("InOrUpUserStas start")
	response = &v1.SyncUserBaseResponse{Affc: -1}
	userBases, _ := d.getUserInfos(c, mids)
	for try = 0; try <= retry; try++ {
		if tx, err = d.BeginTran(c); err != nil {
			time.Sleep(time.Duration(try) * time.Second)
			log.Warn("InOrUpUserStas try begin transaction failed failed ,error(%v)", err)
			continue
		}
		sql := "INSERT INTO user_statistics_hive (mid, uname) VALUES"
		for _, userBase := range userBases {
			sql = sql + "(" + strconv.FormatInt(userBase.Mid, 10) + ",'" + userBase.Name + "'),"
		}
		sql = sql[0:len(sql)-1] + "ON DUPLICATE KEY UPDATE uname=values(uname)"
		if res, err = tx.Exec(sql); err != nil {
			if err = tx.Rollback(); err != nil {
				log.Warn("InOrUpUserStas try rollback failed ,error(%v)", err)
			} else {
				log.Warn("InOrUpUserStas rollbacked")
			}
		} else {
			if err = tx.Commit(); err != nil {
				log.Warn("InOrUpUserStas try commit failed , error(%v)", err)
			} else {
				//提交成功,退出
				response.Affc, _ = res.RowsAffected()
				log.Info("InOrUpUserStas on commit success, affected %v rows", response.Affc)
				break
			}
		}
	}
	if err != nil {
		log.Error("InOrUpUserSta run failed, err(%v)", err)
	}
	return
}

// GetVideoBvcTable 获取bvc分表名
func (d *Dao) getVideoBvcTable(svid int64) string {
	return fmt.Sprintf("video_bvc_%02d", svid%_BVCSubTableSize)
}

//RawVideoStatistic get video statistics
func (d *Dao) RawVideoStatistic(c context.Context, svids []int64) (res map[int64]*model.SvStInfo, err error) {
	const maxIDNum = 20
	var (
		idStr string
	)
	res = make(map[int64]*model.SvStInfo)
	if len(svids) > maxIDNum {
		svids = svids[:maxIDNum]
	}
	l := len(svids)
	for k, svid := range svids {
		if k < l-1 {
			idStr += strconv.FormatInt(svid, 10) + ","
		} else {
			idStr += strconv.FormatInt(svid, 10)
		}
		res[svid] = &model.SvStInfo{}
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_queryStatisticsList, idStr))
	if err != nil {
		log.Error("query error(%s)", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		ssv := new(model.SvStInfo)
		if err = rows.Scan(&ssv.SVID, &ssv.Play, &ssv.Subtitles, &ssv.Like, &ssv.Share, &ssv.Report); err != nil {
			log.Error("RawVideoStatistic rows.Scan() error(%v)", err)
			return
		}
		res[ssv.SVID] = ssv
	}
	cmtCount, _ := d.ReplyCounts(c, svids, DefaultCmType)
	for id, cmt := range cmtCount {
		if _, ok := res[id]; ok {
			res[id].Reply = cmt.Count
		}
	}
	return
}

// CommitTrans 提交转码
func (d *Dao) CommitTrans(c context.Context, arg *v1.BVideoTransRequset) error {
	path, ok := d.c.URLs["bvc_push"]
	if !ok {
		log.Warnv(c, log.KV("log", "bvc_push url not set"))
		return ecode.ReqParamErr
	}
	data, _ := json.Marshal(arg)
	b := string(data)
	req, err := xhttp.NewRequest("POST", path, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Error("bvc_push url(%s) req(%+v) body(%s) error(%v)", path, req, b, err)
		return err
	}
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("bvc_push url(%s) req(%+v) body(%s) ret (%+v) err[%v]", path, req, b, res, err)))
		return err
	}
	log.V(5).Infov(c, log.KV("log", fmt.Sprintf("bvc_push req(%+v) body(%s) ret (%+v)", req, b, res)))
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Errorv(c, log.KV("log", fmt.Sprintf("bvc_push url(%s) req(%+v) body(%s) ret(%+v) error(%v)", path, req, b, res, err)))
		return err
	}
	return nil
}

//AddOrUpdateBVCInfo 添加或更新BVC转码信息
func (d *Dao) AddOrUpdateBVCInfo(c context.Context, arg *model.VideoBVC) (err error) {
	err = d.AddBVCInfo(c, arg)
	if err != nil {
		if matched, _ := regexp.MatchString("Duplicate entry", err.Error()); matched {
			err = d.UpdataBVCInfo(c, arg)
			return
		}
		log.Errorv(c,
			log.KV("log", fmt.Sprintf("dao.db.Exec(AddOrUpdateBVCInfo[%+v]) err(%v)", arg, err)),
		)
	}
	return
}

//TxAddOrUpdateBVCInfo 事务添加或更新BVC转码信息
func (d *Dao) TxAddOrUpdateBVCInfo(c context.Context, tx *xsql.Tx, arg *model.VideoBVC) (err error) {
	err = d.TxAddBVCInfo(tx, arg)
	if err != nil {
		if matched, _ := regexp.MatchString("Duplicate entry", err.Error()); matched {
			err = d.TxUpdataBVCInfo(tx, arg)
			return
		}
		log.Errorv(c,
			log.KV("log", fmt.Sprintf("dao.db.Exec(AddOrUpdateBVCInfo[%+v]) err(%v)", arg, err)),
		)
	}
	return
}

// AddBVCInfo 添加BVC转码信息
func (d *Dao) AddBVCInfo(c context.Context, arg *model.VideoBVC) (err error) {
	t := d.getVideoBvcTable(arg.SVID)
	sql := fmt.Sprintf(_addBVCData, t)
	_, err = d.db.Exec(c, sql, arg.SVID, arg.Path, arg.ResolutionRetio, arg.CodeRate, arg.VideoCode, arg.Duration, arg.FileSize)
	return
}

// TxAddBVCInfo 事务添加BVC转码信息
func (d *Dao) TxAddBVCInfo(tx *xsql.Tx, arg *model.VideoBVC) (err error) {
	t := d.getVideoBvcTable(arg.SVID)
	sql := fmt.Sprintf(_addBVCData, t)
	_, err = tx.Exec(sql, arg.SVID, arg.Path, arg.ResolutionRetio, arg.CodeRate, arg.VideoCode, arg.Duration, arg.FileSize)
	return
}

// TxUpdataBVCInfo 事务更新BVC转码信息
func (d *Dao) TxUpdataBVCInfo(tx *xsql.Tx, arg *model.VideoBVC) (err error) {
	t := d.getVideoBvcTable(arg.SVID)
	sql := fmt.Sprintf(_updateBVCData, t)
	_, err = tx.Exec(sql, arg.Path, arg.ResolutionRetio, arg.VideoCode, arg.Duration, arg.FileSize, arg.SVID, arg.CodeRate)
	return
}

// UpdataBVCInfo 更新BVC转码信息
func (d *Dao) UpdataBVCInfo(c context.Context, arg *model.VideoBVC) (err error) {
	t := d.getVideoBvcTable(arg.SVID)
	sql := fmt.Sprintf(_updateBVCData, t)
	_, err = d.db.Exec(c, sql, arg.Path, arg.ResolutionRetio, arg.VideoCode, arg.Duration, arg.FileSize, arg.SVID, arg.CodeRate)
	return
}

// UpdateCmsSvPIC 更新封面图
func (d *Dao) UpdateCmsSvPIC(c context.Context, svid int64, pic *v1.SvPic, st int64) error {
	_, err := d.cmsdb.Exec(c, _updateSvPIC, pic.PicURL, pic.PicWidth, pic.PicHeight, st, svid)
	return err
}

// HostnameRegister .
func (d *Dao) HostnameRegister(hostnameIndex int64) (succ bool) {
	conn := d.redis.Get(context.Background())
	defer conn.Close()

	redisKey := fmt.Sprintf("hostname:index:%d", hostnameIndex)
	exists, err := redis.Int(conn.Do("EXISTS", redisKey))
	if err != nil {
		log.Errorv(context.Background(), log.KV("event", "fatal"), log.KV("log", fmt.Sprintf("get hostname index from redis fail: key=%s", redisKey)))
		// 即使redis失败了，也给返回成功
		return true
	}
	if exists == 1 {
		return false
	}
	// 不去管返回结果，永远返回成功
	if _, err = conn.Do("SETEX", redisKey, 1000, 1); err != nil {
		log.Errorv(context.Background(), log.KV("event", "fatal"), log.KV("log", fmt.Sprintf("get hostname index from redis fail: key=%s", redisKey)))
	}

	return true
}

// AddVideoViews .
func (d *Dao) AddVideoViews(c context.Context, svid int64, views int) (affected int64, err error) {
	row := d.db.QueryRow(c, _existedStatistics, svid)
	tmp := 0
	if err = row.Scan(&tmp); err != nil || tmp == 0 {
		_, err = d.db.Exec(c, _insertStatistics, svid, 0, 0, 0, 0, 0)
		if err != nil {
			return
		}
	}

	result, err := d.db.Exec(c, _addVideoViews, views, svid)
	if err != nil {
		return
	}

	return result.RowsAffected()
}

// VideoStateUpdate .
func (d *Dao) VideoStateUpdate(c context.Context, svid int64, newState int) (aff int64, err error) {
	result, err := d.db.Exec(c, _updateVideoState, newState, svid)
	if err != nil {
		return
	}

	aff, err = result.RowsAffected()
	return
}
