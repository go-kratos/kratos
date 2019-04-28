package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-common/app/interface/bbq/bullet/api"
	"go-common/app/interface/bbq/bullet/internal/model"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"

	"go-common/app/interface/bbq/bullet/internal/conf"
	user "go-common/app/service/bbq/user/api"
	video "go-common/app/service/bbq/video/api/grpc/v1"
	filter "go-common/app/service/main/filter/api/grpc/v1"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
)

// Dao dao
type Dao struct {
	c            *conf.Config
	redis        *redis.Pool
	db           *xsql.DB
	filterClient filter.FilterClient
	userClient   user.UserClient
	videoClient  video.VideoClient
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:            c,
		redis:        redis.NewPool(c.Redis),
		db:           xsql.NewMySQL(c.MySQL),
		filterClient: newFilterClient(c.GRPCClient["filter"]),
		userClient:   newUserClient(c.GRPCClient["user"]),
		videoClient:  newVideoClient(c.GRPCClient["video"]),
	}
	return
}

// newVideoClient .
func newVideoClient(cfg *conf.GRPCConf) video.VideoClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return video.NewVideoClient(cc)
}

// newUserClient .
func newUserClient(cfg *conf.GRPCConf) user.UserClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return user.NewUserClient(cc)
}

// newUserClient .
func newFilterClient(cfg *conf.GRPCConf) filter.FilterClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return filter.NewFilterClient(cc)
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(ctx context.Context) error {
	// TODO: add mc,redis... if you use
	return d.db.Ping(ctx)
}

// ContentPost .
func (d *Dao) ContentPost(ctx context.Context, req *api.Bullet) (dmid int64, err error) {
	result, err := d.db.Exec(ctx,
		"insert into bullet_content (oid, mid, offset_ms, offset, content) values (?, ?, ?, ?, ?)",
		req.Oid, req.Mid, req.OffsetMs, req.OffsetMs/1000, req.Content)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "insert bullet fail: req=%s"+req.String()))
		return
	}
	dmid, err = result.LastInsertId()
	return
}

// ContentGet .
func (d *Dao) ContentGet(ctx context.Context, req *api.ListBulletReq) (res []*api.Bullet, err error) {
	res = []*api.Bullet{}
	mid := req.Mid

	querySQL := fmt.Sprintf("select id, mid, offset, content from bullet_content where "+
		"oid=%d and state=0 and offset>=%d and offset<%d order by offset, id desc",
		req.Oid, req.StartMs/1000, req.EndMs/1000)
	rows, err := d.db.Query(ctx, querySQL)
	if err != nil {
		return
	}
	defer rows.Close()
	log.V(1).Infow(ctx, "sql", querySQL)
	// 获取时间范围内的全量视频
	var allBullet []*api.Bullet
	midBullets := make(map[int32]*[]*api.Bullet)
	for rows.Next() {
		bullet := new(api.Bullet)
		if err = rows.Scan(&bullet.Id, &bullet.Mid, &bullet.Offset, &bullet.Content); err != nil {
			log.Errorv(ctx, log.KV("log", "scan mysql fail: sql="+querySQL))
			return
		}
		bullet.OffsetMs = bullet.Offset * 1000
		allBullet = append(allBullet, bullet)
		// 先把访问者发过的弹幕按照秒级别进行汇总
		if mid == bullet.Mid {
			v, exists := midBullets[bullet.Offset]
			if !exists {
				v = new([]*api.Bullet)
				midBullets[bullet.Offset] = v
			}
			*v = append(*v, bullet)
		}
	}

	// 根据全量数据，选择满足条件的弹幕
	currSecond := int32(-1)
	currSecondCount := 0
	for _, bullet := range allBullet {
		if currSecond != bullet.Offset {
			currSecond = bullet.Offset
			currSecondCount = 0
			if midBulletArray, exists := midBullets[currSecond]; exists {
				log.V(10).Infow(ctx, "log", "current second user have published danmu", "offset", currSecond, "len", len(*midBulletArray))
				for _, midBullet := range *midBulletArray {
					currSecondCount++
					res = append(res, midBullet)
					if currSecondCount >= model.SecondMaxNum {
						break
					}
				}
			}
		}
		if currSecondCount >= model.SecondMaxNum {
			continue
		}
		if bullet.Mid != mid {
			currSecondCount++
			bullet.OffsetMs = bullet.Offset * 1000
			res = append(res, bullet)
		}
	}

	if len(res) > 0 {
		var cursor CursorValue
		cursor.Offset = res[len(res)-1].Offset
		b, _ := json.Marshal(cursor)
		res[len(res)-1].CursorValue = string(b)
	}

	return
}

// ContentList 用于返回弹幕列表
/*
 */
func (d *Dao) ContentList(ctx context.Context, req *api.ListBulletReq) (res *api.ListBulletReply, err error) {
	res = new(api.ListBulletReply)

	// 0. 前期准备
	// 获取当前oid的最大offset弹幕的offset
	oidLastOffset, err := d.lastOffset(ctx, req.Oid)
	if err != nil {
		log.Warnv(ctx, log.KV("log", "get has more info fail"))
		return
	}
	// 解析cursor
	cursor, err := parseCursorValue(ctx, req.CursorNext)
	if err != nil {
		log.Warnv(ctx, log.KV("log", "parse cursor value fail"))
		return
	}
	// 当两者相等，则说明已经到列表的最后了
	if oidLastOffset <= cursor.Offset {
		res.HasMore = false
		log.Warnw(ctx, "log", "offset already end", "oid_last_offset", oidLastOffset, "cursor_offset", cursor.Offset)
		return
	}

	// 1. 按照条数取SecondMaxNum条，返回数据的offset范围start和end
	// 这步是为了保证该次返回至少有条数
	startS := cursor.Offset + 1
	startS, endS, err := d.getNumBulletTs(ctx, req.Oid, startS, model.SecondMaxNum)
	if err != nil {
		log.Warnv(ctx, log.KV("log", "get num start bullet fail"))
		return
	}
	log.V(1).Infow(ctx, "log", "get num bullet ts", "start_s", startS, "end_s", endS)
	endS += 1

	// 2. 根据选择的时间范围获取弹幕
	newReq := &api.ListBulletReq{StartMs: startS * 1000, EndMs: endS * 1000, Oid: req.Oid, Mid: req.Mid}
	bullets, err := d.ContentGet(ctx, newReq)
	if err != nil {
		log.Warnv(ctx, log.KV("log", "content get fail: req="+newReq.String()))
		return
	}
	res.List = bullets

	// 3. has_more设置，如果offset和最后时间offset相等，那么肯定没有更多弹幕了
	if len(bullets) > 0 && oidLastOffset > bullets[len(bullets)-1].Offset {
		res.HasMore = true
	} else {
		res.HasMore = false
	}

	return
}

func (d *Dao) getNumBulletTs(ctx context.Context, oid int64, startOffset, size int32) (startS, endS int32, err error) {
	querySQL := fmt.Sprintf(
		"select offset from bullet_content where oid=%d and state=0 and offset>=%d order by offset limit %d",
		oid, startOffset, size)
	rows, err := d.db.Query(ctx, querySQL)
	if err != nil {
		log.Errorv(ctx, log.KV("log", fmt.Sprintf("get num bullet from db fail: sql=%s", querySQL)))
		return
	}
	log.V(1).Infow(ctx, "sql", querySQL)
	var offset int32
	var index int32
	for rows.Next() {
		if err = rows.Scan(&offset); err != nil {
			log.Errorv(ctx, log.KV("log", "scan mysql fail: sql="+querySQL))
			return
		}
		if index == 0 {
			startS = offset
		}
		endS = offset
		index++
	}
	return
}

func (d *Dao) lastOffset(ctx context.Context, oid int64) (lastOffset int32, err error) {
	querySQL := fmt.Sprintf("select offset from bullet_content where oid=%d and state=0 order by offset desc limit 1", oid)
	row := d.db.QueryRow(ctx, querySQL)
	if err = row.Scan(&lastOffset); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			lastOffset = -1
		} else {
			log.Errorw(ctx, "log", "get has more from db fail", "sql", querySQL, "err", err)
			return
		}
	}
	return
}

// CursorValue .
type CursorValue struct {
	Offset int32 `json:"offset"`
	// level本来是想要用于避免一次选择太少的弹幕，但后面修改策略进行二次查找之后就没这个必要了
	//Level  int32 `json:"level"`
	//duration int32 `json:"duration"`
}

func parseCursorValue(ctx context.Context, cursorValue string) (cursor CursorValue, err error) {
	if len(cursorValue) == 0 {
		cursor.Offset = -1
		//cursor.Level = 1
		return
	}
	if err = json.Unmarshal([]byte(cursorValue), &cursor); err != nil {
		log.Errorw(ctx, "log", "unmarshal fail: str="+cursorValue, "err", err)
		return
	}
	return
}

//
//// 这里做了个优化，当对于弹幕数较少的视频，level等级定的高点，在弹幕列表页中就可以选取更长范围的弹幕
//func getCursorLevel(duration int32, num int32) (level int32) {
//	numPerSecond := num / duration
//	if numPerSecond < 1 {
//		level = 10
//	} else if numPerSecond < 2 {
//		level = 5
//	} else if numPerSecond < 5 {
//		level = 2
//	} else {
//		level = 1
//	}
//	return
//}
