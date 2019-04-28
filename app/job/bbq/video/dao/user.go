package dao

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-common/app/job/bbq/video/conf"
	"go-common/app/job/bbq/video/model"
	video "go-common/app/service/bbq/video/api/grpc/v1"
	acc "go-common/app/service/main/account/api"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	_userDmgCacheKey      = "bbq:user:profile:%s"
	_userDmgCacheKeyBuvid = "bbq:device:profile:{buvid}:%s"

	//_userDmgCacheKeyBbq="bbq:user:profile"
	_userDmgCacheTimeout = 86400
	_incrUpUserDmgSQL    = "insert into user_statistics_hive (`mid`,`uname`,`play_total`,`fan_total`,`av_total`,`like_total`) value (?,?,?,?,?,?)"
	_updateUpUserDmgSQL  = "update user_statistics_hive set `uname` = ? , `play_total` = ? , `fan_total` = ? , `av_total` = ? , `like_total` = ? , `mtime` = ? where `mid` = ?"
	_delUpUserDmgSQL     = "delete from user_statistics_hive where mtime < ?"
	_insertOnDupUpDmg    = "insert into user_statistics_hive (`mid`,`uname`,`play_total`,`fan_total`,`av_total`,`like_total`) value (?,?,?,?,?,?) on duplicate key update `uname`=?,`play_total`=?,`fan_total`=?,`av_total`=?,`like_total`=?"
	_selMidFromVideo     = "select distinct mid from video"
	_queryUsersByLast    = "select id,mid,uname from user_base where id > ? order by id ASC limit ?"
	_selMidFromUserBase  = "SELECT DISTINCT mid fROM user_base limit ?, 1000"
	_upUserBase          = "UPDATE user_base SET face = ? where mid = ?"
)

//getUserDmgKey .
func getUserDmgKey(mid string) (key string) {
	return fmt.Sprintf(_userDmgCacheKey, mid)
}

//getUserBuvidDmgKey .
func getUserBuvidDmgKey(buvid string) (key string) {
	return fmt.Sprintf(_userDmgCacheKeyBuvid, buvid)
}

// InsertOnDup ...
func (d *Dao) InsertOnDup(c context.Context, upUserDmg *model.UpUserDmg) (err error) {
	_, err = d.db.Exec(c, _insertOnDupUpDmg, upUserDmg.MID, upUserDmg.Uname, upUserDmg.Play, upUserDmg.Fans, upUserDmg.AVs, upUserDmg.Likes, upUserDmg.Uname, upUserDmg.Play, upUserDmg.Fans, upUserDmg.AVs, upUserDmg.Likes)
	return
}

//CacheUserDmg ...
func (d *Dao) CacheUserDmg(c context.Context, userDmg *model.UserDmg) (mid string, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var b []byte
	if b, err = json.Marshal(userDmg); err != nil {
		log.Error("cache user dmg marshal err(%v)", err)
		return
	}
	cacheKey := getUserDmgKey(userDmg.MID)
	fmt.Println(cacheKey)
	if _, err = conn.Do("SET", cacheKey, b, "EX", _userDmgCacheTimeout); err != nil {
		log.Error("cache user dmg redis set err(%v)", err)
		return
	}
	return
}

//CacheUserBbqDmg ...
func (d *Dao) CacheUserBbqDmg(c context.Context, userBbqDmg *model.UserBbqDmg) (mid string, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	tag2 := strings.Join(userBbqDmg.Tag2, ",")
	tag3 := strings.Join(userBbqDmg.Tag3, ",")
	up := strings.Join(userBbqDmg.Up, ",")
	cacheKey := getUserDmgKey(userBbqDmg.MID)

	if err = conn.Send("HSET", cacheKey, "zone", tag2); err != nil {
		log.Error("cache user bbq dmg redis set tag2 err(%v)", err)
		return
	}

	if err = conn.Send("HSET", cacheKey, "tag", tag3); err != nil {
		log.Error("cache user bbq dmg redis set tag3 err(%v)", err)
		return
	}

	if err = conn.Send("HSET", cacheKey, "up", up); err != nil {
		log.Error("cache user bbq dmg redis set up err(%v)", err)
		return
	}

	return
}

//CacheUserBbqDmgBuvid ...
func (d *Dao) CacheUserBbqDmgBuvid(c context.Context, userBbqDmgBuvid *model.UserBbqBuvidDmg) (Buvid string, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	tag2 := strings.Join(userBbqDmgBuvid.Tag2, ",")
	tag3 := strings.Join(userBbqDmgBuvid.Tag3, ",")
	up := strings.Join(userBbqDmgBuvid.Up, ",")
	cacheKey := getUserBuvidDmgKey(userBbqDmgBuvid.Buvid)

	if err = conn.Send("HSET", cacheKey, "zone", tag2); err != nil {
		log.Error("cache user bbq buvid dmg redis set tag2 err(%v)", err)
		return
	}

	if err = conn.Send("HSET", cacheKey, "tag", tag3); err != nil {
		log.Error("cache user bbq buvid dmg redis set tag3 err(%v)", err)
		return
	}

	if err = conn.Send("HSET", cacheKey, "up", up); err != nil {
		log.Error("cache user bbq buvid dmg redis set up err(%v)", err)
		return
	}

	return
}

// AddUpUserDmg .
func (d *Dao) AddUpUserDmg(c context.Context, upUserDmg *model.UpUserDmg) (num int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _incrUpUserDmgSQL, upUserDmg.MID, upUserDmg.Uname, upUserDmg.Play, upUserDmg.Fans, upUserDmg.AVs, upUserDmg.Likes); err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateUpUserDmg .
func (d *Dao) UpdateUpUserDmg(c context.Context, upUserDmg *model.UpUserDmg) (num int64, err error) {
	t := time.Now().AddDate(0, 0, 0).Format("2006-01-02 15:04:05")
	var res sql.Result
	if res, err = d.db.Exec(c, _updateUpUserDmgSQL, upUserDmg.Uname, upUserDmg.Play, upUserDmg.Fans, upUserDmg.AVs, upUserDmg.Likes, t, upUserDmg.MID); err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// DelUpUserDmg .
func (d *Dao) DelUpUserDmg(c context.Context) (num int64, err error) {
	t := time.Unix(time.Now().Unix(), -int64(36*time.Hour)).Format("2006-01-02 15:04:05")
	var res sql.Result
	if res, err = d.db.Exec(c, _delUpUserDmgSQL, t); err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

//Download 下载文件
func (d *Dao) Download(url string, name string) (fpath string, err error) {
	if name == "" {
		u := strings.Split(url, "/")
		l := len(u)
		name = u[l-1]
	}
	t := time.Now().AddDate(0, 0, 0).Format("20060102")
	path := conf.Conf.Download.File + t
	err = d.CreateDir(path)
	if err != nil {
		log.Error("create dir(%s) err(%v)", path, err)
		return
	}
	fpath = path + "/" + name
	newFile, err := os.Create(fpath)

	if err != nil {
		log.Error("create path(%s) err(%v)", fpath, err)
		return
	}
	defer newFile.Close()

	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		log.Error("download url(%s) err(%v)", url, err)
		return
	}
	defer resp.Body.Close()
	_, err = io.Copy(newFile, resp.Body)
	if err != nil {
		log.Error("copy err(%v)", err)
		return
	}
	return
}

//CreateDir 创建文件夹
func (d *Dao) CreateDir(path string) (err error) {
	_, err = os.Stat(path)
	defer func() {
		if os.IsExist(err) {
			err = nil
		}
	}()
	if os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
	}
	return
}

// ReadLine 按行读取文件，hander回调
func (d *Dao) ReadLine(path string, handler func(string)) (err error) {
	f, err := os.Open(path)
	if err != nil {
		log.Error("open path(%s) err(%v)", path, err)
		return
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {

		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			log.Error("read path(%s) err(%v)", path, err)
			return nil
		}
		line = strings.TrimSpace(line)
		handler(line)
		time.Sleep(time.Duration(1) * time.Second)
	}
}

// ReadLines 50条发起一次grpc请求
func (d *Dao) ReadLines(path string, handler func([]int64)) (err error) {
	f, err := os.Open(path)
	if err != nil {
		log.Error("ReadLine open path(%s) err(%v)", path, err)
		return
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	mids := make([]int64, 0, 50)
	i := 0
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			log.Error("read path(%s) err(%v)", path, err)
			break
		}
		mid, _ := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
		mids = append(mids, mid)
		i++
		if i == 50 {
			handler(mids)
			mids = make([]int64, 0, 50)
			i = 0
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
	if len(mids) != 0 {
		handler(mids)
	}
	return
}

//HandlerUserDmg mid, gender, age, geo, content_tag, viewed_video, content_zone, content_count, follow_ups
func (d *Dao) HandlerUserDmg(user string) {
	u := strings.Split(user, "\u0001")
	userDmg := &model.UserDmg{
		MID:          u[0],
		Gender:       u[1],
		Age:          u[2],
		Geo:          u[3],
		ContentTag:   u[4],
		ViewedVideo:  d.HandlerViewedVideo(u[5]),
		ContentZone:  u[6],
		ContentCount: u[7],
		FollowUps:    u[8],
	}
	d.CacheUserDmg(context.Background(), userDmg)
}

//HandlerUserBbqDmg ..
func (d *Dao) HandlerUserBbqDmg(user string) {
	u := strings.Split(user, ",")
	userBbqDmg := &model.UserBbqDmg{
		MID:  u[0],
		Tag2: strings.Split(u[1], "\u0002"),
		Tag3: strings.Split(u[2], "\u0002"),
		Up:   strings.Split(u[3], "\u0002"),
	}
	d.CacheUserBbqDmg(context.Background(), userBbqDmg)
}

//HandlerUserBbqDmgBuvid ..
func (d *Dao) HandlerUserBbqDmgBuvid(user string) {
	u := strings.Split(user, ",")
	UserBbqBuvidDmg := &model.UserBbqBuvidDmg{
		Buvid: u[0],
		Tag2:  strings.Split(u[1], "\u0002"),
		Tag3:  strings.Split(u[2], "\u0002"),
		Up:    strings.Split(u[3], "\u0002"),
	}
	d.CacheUserBbqDmgBuvid(context.Background(), UserBbqBuvidDmg)
}

// HandlerMids update userbase by mids
func (d *Dao) HandlerMids(mids []int64) {
	res, err := d.VideoClient.SyncUserStas(context.Background(), &video.SyncMidsRequset{MIDS: mids})
	if err != nil {
		log.Error("userbases update failes, mids(%v), err(%v)", mids, err)
		return
	}
	log.Info("userbases update success, affected %v rows", res.Affc)
}

// HandlerMid update userbase by mid
func (d *Dao) HandlerMid(s string) {
	mid, _ := strconv.ParseInt(s, 10, 64)
	res, err := d.VideoClient.SyncUserSta(context.Background(), &video.SyncMidRequset{MID: mid})
	if err != nil {
		log.Error("userbase update failes, mid(%v), err(%v)", mid, err)
		return
	}
	if res.Affc == 1 {
		log.Info("userbase insert success ,mid(%v)", mid)
	} else if res.Affc == 2 {
		log.Info("userbase update success , mid(%v)", mid)
	}
}

//HandlerViewedVideo 处理看过的视频，保存最近看过的100个
func (d *Dao) HandlerViewedVideo(v string) (res map[int64]string) {
	res = make(map[int64]string)
	var vv [][]interface{}
	var dd string

	err := json.Unmarshal([]byte(v), &vv)
	if err != nil {
		return
	}
	l := len(vv)
	n := 1
	for i := l - 1; i >= 0; i-- {
		for _, a := range vv[i] {
			switch b := a.(type) {
			case string:
				dd = b
			case []interface{}:
				ll := len(b)
				for j := ll - 1; j >= 0; j-- {
					switch c := b[j].(type) {
					case float64:
						k := int64(c)
						if _, ok := res[k]; !ok {
							res[k] = dd
							n++
						}
					}
					if n > 100 {
						return
					}
				}
			}
		}
	}
	return
}

// SelMidFromVideo get distinct mid list from table video
func (d *Dao) SelMidFromVideo() (mids []int64, err error) {
	rows, err := d.db.Query(context.Background(), _selMidFromVideo)
	if err != nil {
		log.Error("SelMidFromVideo failed, err(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var s string
		if err = rows.Scan(&s); err != nil {
			panic(err.Error())
		}
		var mid int64
		if mid, err = strconv.ParseInt(s, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", s, err)
			return
		}
		mids = append(mids, mid)
	}
	return
}

//MergeUpInfo merge up info
func (d *Dao) MergeUpInfo(mid int64) (err error) {
	var (
		ctx    = context.Background()
		params = url.Values{}
		req    = &http.Request{}
		id     int64
		res    struct {
			Code int
			Data model.UpUserInfoRes
		}
	)
	err = d.db.QueryRow(ctx, "select mid from user_base where mid = ?", mid).Scan(&id)
	if err == nil {
		log.Infow(ctx, "log", "already has mid in user_base", "mid", mid)
		return
	}
	if err == sql.ErrNoRows {
		params.Set("mid", strconv.FormatInt(mid, 10))
		req, err = d.HTTPClient.NewRequest("GET", d.c.URLs["account"], "", params)
		if err != nil {
			log.Error("MergeUpInfo error(%v)", err)
			return
		}
		if err = d.HTTPClient.Do(ctx, req, &res); err != nil {
			log.Error("MergeUpInfo http req failed ,err:%v", err)
			return
		}
		res := res.Data
		var sex int
		switch res.Sex {
		case "男":
			sex = 1
		case "女":
			sex = 2
		default:
			sex = 3
		}
		_, err = d.db.Exec(ctx,
			"insert into user_base (mid,uname,face,sex,user_type,complete_degree)values(?,?,?,?,?,?)",
			res.MID,
			res.Name,
			res.Face,
			sex,
			model.UserTypeUp,
			0)
		if err != nil {
			log.Error("MergeUpInfo insert upinfo failed,err:%v", err)
			return
		}
	} else {
		log.Error("MergeUpInfo query sql failed,err:%v", err)
	}
	if err = d.db.QueryRow(ctx, "select id from user_statistics where mid = ?", mid).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			if _, err = d.db.Exec(ctx, "insert into user_statistics (mid) values (?)", mid); err != nil {
				log.Error("init insert user_statistics failed,err:%v", err)
			}
		} else {
			log.Error("init query user_statistics failed,err:%v", err)
		}
	}
	return
}

//UsersByLast 使用lastid批量获取用户
func (d *Dao) UsersByLast(c context.Context, lastid int64) (r []*model.UserBaseDB, err error) {
	var rows *xsql.Rows
	rows, err = d.db.Query(c, _queryUsersByLast, lastid, _limitSize)
	if err != nil {
		log.Error("db _queryVideos err(%v)", err)
		return
	}
	for rows.Next() {
		u := new(model.UserBaseDB)
		if err = rows.Scan(&u.ID, &u.MID, &u.Uname); err != nil {
			log.Error("scan err(%v)", err)
			continue
		}
		r = append(r, u)
	}
	return
}

// SelMidFromUserBase get distinct mid list from table user_base
func (d *Dao) SelMidFromUserBase(start int) (mids []int64, err error) {
	var mid int64
	rows, err := d.db.Query(context.Background(), _selMidFromUserBase, start)
	if err != nil {
		log.Error("SelMidFromUserBase failed, err(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var s string
		if err = rows.Scan(&s); err != nil {
			panic(err.Error())
		}
		if mid, err = strconv.ParseInt(s, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", s, err)
			return
		}
		mids = append(mids, mid)
	}
	return
}

// UpUserBases 根据mids更新用户基本信息
func (d *Dao) UpUserBases(c context.Context, mids []int64) (err error) {
	var (
		tx *xsql.Tx
	)
	midsReq := &acc.MidsReq{
		Mids:   mids,
		RealIp: metadata.String(c, metadata.RemoteIP)}
	infosReply, err := d.AccountClient.Infos3(c, midsReq)
	if infosReply == nil {
		log.Error("查询infos3失败,err（%v）", err)
		fmt.Printf("查询infos3失败,err（%v）", err)
		return
	}
	if tx, err = d.BeginTran(c); err != nil {
		log.Error("begin transaction error(%v)", err)
		return
	}
	for _, info := range infosReply.Infos {
		if info.Mid != 0 {
			if len(info.Face) > 255 {
				info.Face = "http://i0.hdslb.com/bfs/bbq/video-image/userface/1558868601542006937.png"
			}
			for try := 0; try < 3; try++ {
				if _, err = tx.Exec(_upUserBase, info.Face, info.Mid); err == nil {
					break
				}
			}
			if err != nil {
				log.Error("mid(%v) update failed", info.Mid)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("UpUserBases commit failed err(%v)", err)
	}
	return
}
