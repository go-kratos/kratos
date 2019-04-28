package dao

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"go-common/app/service/live/recommend/internal/conf"
	"go-common/app/service/live/recommend/recconst"
	relation_api "go-common/app/service/live/relation/api/liverpc"
	room_api "go-common/app/service/live/room/api/liverpc"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/net/rpc/liverpc"
)

var _userRecCandidateKey = "rec_candidate_%d"
var _recommendOffsetKey = "rec_offset_%d"

// 已经推荐过的池子，用户+日期
var _recommendedKey = "recommended_%d_%s"

// RoomAPI room liverpc client
var RoomAPI *room_api.Client

// RelationAPI relation liverpc client
var RelationAPI *relation_api.Client

// Dao dao
type Dao struct {
	c     *conf.Config
	redis *redis.Pool
}

func init() {
	RoomAPI = room_api.New(getConf("room"))
	RelationAPI = relation_api.New(getConf("relation"))
}

func getConf(appName string) *liverpc.ClientConfig {
	c := conf.Conf.LiveRpc
	if c != nil {
		return c[appName]
	}
	return nil
}

// ClearRecommend 清空该用户相关的推荐缓存
func (d *Dao) ClearRecommend(ctx context.Context, uid int64) error {
	candidateKey := fmt.Sprintf(_userRecCandidateKey, uid)
	recommendedKey := fmt.Sprintf(_recommendedKey, uid, time.Now().Format("20060102"))
	offsetKey := fmt.Sprintf(_recommendOffsetKey, uid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err := conn.Do("DEL", candidateKey, recommendedKey, offsetKey)
	return errors.WithStack(err)
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:     c,
		redis: redis.NewPool(c.Redis),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
}

func (d *Dao) saveOffset(conn redis.Conn, uid int64, offset int) {
	conn.Do("SETEX", fmt.Sprintf(_recommendOffsetKey, uid), 86400, offset)
}

func (d *Dao) addToRecommended(conn redis.Conn, uid int64, ids []int64) {
	if len(ids) == 0 {
		return
	}
	day := time.Now().Format("20060102")

	key := fmt.Sprintf(_recommendedKey, uid, day)
	var is []interface{}
	is = append(is, key)
	for _, id := range ids {
		is = append(is, id)
	}

	conn.Send("EXPIRE", key, 86400)
	conn.Send("SADD", is...)
	conn.Flush()
	conn.Receive()
	_, err := conn.Receive()
	if err != nil {
		log.Info("addToRecommended error +%v", err)
	}
}

// GetRandomRoomIds 随机获取count个推荐
// 如果总数量total比count小，则返回total个
func (d *Dao) GetRandomRoomIds(ctx context.Context, uid int64, reqCount int, existRoomIDs []int64) (ret []int64, err error) {
	if reqCount == 0 {
		return
	}
	var (
		candidateLen int
	)
	r := d.redis.Get(ctx)
	defer r.Close()

	candidateKey := fmt.Sprintf(_userRecCandidateKey, uid)
	exists, err := redis.Int(r.Do("exists", candidateKey))
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	existMap := map[int64]struct{}{}
	for _, id := range existRoomIDs {
		existMap[id] = struct{}{}
	}
	if exists == 0 {
		var candidate []int64
		var currentOffset = 0
		candidate, err = d.generateLrCandidateList(r, uid, candidateKey)
		if err != nil {
			return
		}
	Loop:
		for len(ret) < reqCount && currentOffset < len(candidate) {
			var tmp []int64
			if len(candidate)-currentOffset < int(reqCount) {
				tmp = candidate[currentOffset:]
			} else {
				tmp = candidate[currentOffset : currentOffset+reqCount]
			}
			//去重
			for _, id := range tmp {
				_, ok := existMap[id]
				currentOffset += 1
				if !ok {
					ret = append(ret, id)
					if len(ret) >= int(reqCount) {
						break Loop
					}
				}
			}
		}

		d.addToRecommended(r, uid, ret)
		d.saveOffset(r, uid, currentOffset)
	} else {
		candidateLen, err = redis.Int(r.Do("LLEN", candidateKey))
		if err != nil {
			return
		}

		var offset int
		offset, _ = redis.Int(r.Do("GET", fmt.Sprintf(_recommendOffsetKey, uid)))
		if offset > (candidateLen - 1) {
			return
		}
		var currentOffset = offset
	Loop2:
		for len(ret) < reqCount && currentOffset < candidateLen {
			var ids []int64
			ids, err = redis.Int64s(r.Do("LRANGE", candidateKey, currentOffset, currentOffset+reqCount-1))
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			// 去重
			for _, id := range ids {
				currentOffset++
				_, ok := existMap[id]
				if !ok {
					ret = append(ret, id)
					if len(ret) >= int(reqCount) {
						break Loop2
					}
				}
			}
			if len(ids) == 0 {
				log.Error("Cannot get recommend candidate, key=%s, offset=%d, count=%d", candidateKey, offset, reqCount)
				break
			}
		}

		d.addToRecommended(r, uid, ret)
		d.saveOffset(r, uid, currentOffset)
	}
	return
}

// GetLrRecRoomIds 在GetRandomRoomIds的基础上进行LR计算并返回倒排的房间号列表
// 与GetRandomRoomIds有相同的输入输出结构
func (d *Dao) GetLrRecRoomIds(r redis.Conn, uid int64, candidateIds []int64) (ret []int64, err error) {
	var areas string
	areaIds := map[int64]struct{}{}
	areas, err = redis.String(r.Do("GET", fmt.Sprintf(recconst.UserAreaKey, uid)))
	if err != nil && err != redis.ErrNil {
		log.Error("redis GET error: %v", err)
		return
	}
	err = nil
	if areas != "" {
		split := strings.Split(areas, ";")
		for _, areaIdStr := range split {
			areaId, _ := strconv.ParseInt(areaIdStr, 10, 64)
			areaIds[areaId] = struct{}{}
		}
	}

	weightVector := makeWeightVec(d.c)
	roomFeatures, ok := roomFeatureValue.Load().(map[int64][]int64)
	if !ok {
		ret = candidateIds
		return
	}
	roomScoreSlice := ScoreSlice{}
	for _, roomId := range candidateIds {
		if fv, ok := roomFeatures[roomId]; ok {
			featureVector := make([]int64, len(fv))
			copy(featureVector, fv)
			areaId := featureVector[0]
			if _, ok := areaIds[areaId]; ok {
				featureVector[0] = 1
			} else {
				featureVector[0] = 0
			}
			counter := Counter{roomId: roomId, score: calcScore(weightVector, featureVector)}
			roomScoreSlice = append(roomScoreSlice, counter)
		}
	}

	sort.Sort(roomScoreSlice)
	for _, counter := range roomScoreSlice {
		ret = append(ret, counter.roomId)
	}
	return
}

// generateCandidateList 得到候选集
func (d *Dao) generateCandidateList(r redis.Conn, uid int64, candidateKey string) (ret []int64, err error) {

	// 第一步 itemcf，优先级最高。
	itemCFKey := fmt.Sprintf(recconst.UserItemCFRecKey, uid)
	var itemCFList []int64
	itemCFList, err = redis.Int64s(r.Do("ZREVRANGE", itemCFKey, 0, -1))
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	itemCFOnlineIds := d.FilterOnlineRoomIds(itemCFList)
	if len(itemCFOnlineIds) == 0 {
		log.Info("No item-cf room online for user, uid=%d, before online filter room ids: %+v", uid, itemCFList)
	}

	// 第二步 取兴趣分区的房间 人气超过100的房间
	var areas string
	areas, err = redis.String(r.Do("GET", fmt.Sprintf(recconst.UserAreaKey, uid)))
	if err != nil && err != redis.ErrNil {
		err = errors.WithStack(err)
		return
	}
	err = nil
	var areaRoomIDs []int64
	if areas != "" {
		split := strings.Split(areas, ";")
		for _, areaIdStr := range split {
			areaId, _ := strconv.ParseInt(areaIdStr, 10, 64)
			var ids = d.getAreaRoomIds(areaId)
			areaRoomIDs = append(areaRoomIDs, ids...)
		}
	}

	// 第三步 取兴趣分区大分区的100个 先不做
	// 第四步 减去已经推荐过的
	day := time.Now().Format("20060102")
	var recommendedList []int64
	edKey := fmt.Sprintf(_recommendedKey, uid, day)
	recommendedList, err = redis.Int64s(r.Do("SMEMBERS", edKey))
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	recommended := map[int64]struct{}{}
	for _, id := range recommendedList {
		recommended[id] = struct{}{}
	}
	var itemCFFinalIDs []int64
	for _, id := range itemCFOnlineIds {
		_, exist := recommended[id]
		if !exist {
			itemCFFinalIDs = append(itemCFFinalIDs, id)
		}
	}

	var areaRoomFinalIDs []int64
	for _, id := range areaRoomIDs {
		_, exist := recommended[id]
		if !exist {
			areaRoomFinalIDs = append(areaRoomFinalIDs, id)
		}
	}

	ret = mergeArr(itemCFFinalIDs, areaRoomFinalIDs)

	log.Info("UserRecommend : uid=%d total=%d, "+
		"itemcf.original=%d, itemcf.online=%d, itemcf.noviewd=%d, "+
		"areaRoom.original=%d, itemcf.noviewd=%d viewed=%d",
		uid, len(ret), len(itemCFList), len(itemCFOnlineIds), len(itemCFFinalIDs),
		len(areaRoomIDs), len(areaRoomFinalIDs), len(recommendedList))
	return
}

// generateCandidateList 得到进过LR的候选集
func (d *Dao) generateLrCandidateList(r redis.Conn, uid int64, candidateKey string) (ret []int64, err error) {
	roomIDs, err := d.generateCandidateList(r, uid, candidateKey)
	if err != nil {
		log.Error("generateLrCandidateList failed 1, error:%v", err)
		return
	}

	if len(ret) > 0 {
		ret, err = d.GetLrRecRoomIds(r, uid, roomIDs)
		if err != nil {
			log.Error("generateLrCandidateList failed 2, error:%v", err)
			return
		}
	}

	// 召回源不足的情况下补足推荐房间数
	if len(ret) < 150 {
		ids, ok := recDefaultRoomIds.Load().([]int64)
		if !ok {
			return
		}
		ret1, err1 := d.GetLrRecRoomIds(r, uid, ids)
		if err1 != nil {
			log.Error("generateLrCandidateList failed 3, error:%v", err1)
			return
		}
		ret = mergeArrWithOrder(ret, ret1, 150) // TODO:当前ret1的结果是没有过滤掉今天看过的房间的, 看后面是否需要优化
	}
	{
		for _, roomID := range ret {
			r.Send("RPUSH", candidateKey, roomID)
		}
		r.Send("EXPIRE", candidateKey, 60*2)
		err = r.Flush()
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		for i := 0; i < len(ret)+1; i++ {
			r.Receive()
		}
	}
	return
}

// Ping dao ping
func (d *Dao) Ping(ctx context.Context) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		err = errors.Wrap(err, "dao Ping err")
	}
	return err
}

// Counter 房间-分数结构体, 用于构建一个可排序的slice
type Counter struct {
	roomId int64
	score  float32
}

// ScoreSlice Counter对象的slice
type ScoreSlice []Counter

func (s ScoreSlice) Len() int {
	return len(s)
}

func (s ScoreSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ScoreSlice) Less(i, j int) bool {
	return s[j].score < s[i].score
}

func calcScore(weightVector []float32, featureVector []int64) (score float32) {
	if len(weightVector) != len(featureVector) {
		panic(fmt.Sprintf("权重数量和特征数量不匹配, 请检查配置或逻辑, weight: %+v, feature: %+v", weightVector, featureVector))
	}
	for i := 0; i < min(len(weightVector), len(featureVector)); i++ {
		score += weightVector[i] * float32(featureVector[i])
	}
	return
}

func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

// 合并两个集合
func mergeArr(x []int64, y []int64) (ret []int64) {
	tmpMap := map[int64]struct{}{}
	for _, id := range x {
		tmpMap[id] = struct{}{}
	}
	for _, id := range y {
		tmpMap[id] = struct{}{}
	}
	for id := range tmpMap {
		ret = append(ret, id)
	}
	return
}

// 按x, y的顺序合并两个集合, 当x的长度不小于limit则直接返回
func mergeArrWithOrder(x []int64, y []int64, limit int) (ret []int64) {
	if len(x) >= limit {
		ret = x
		return
	}
	tmpMap := map[int64]struct{}{}
	ret = append(ret, x...)
	num := len(ret)
	for _, id := range x {
		tmpMap[id] = struct{}{}
	}
	for _, id := range y {
		if _, ok := tmpMap[id]; ok {
			continue
		}
		num += 1
		tmpMap[id] = struct{}{}
		ret = append(ret, id)
		if num >= limit {
			break
		}
	}
	return
}

func makeWeightVec(c *conf.Config) (ret []float32) {
	ret = append(ret, c.CommonFeature.UserAreaInterest.Weights...)
	ret = append(ret, c.CommonFeature.FansNum.Weights...)
	ret = append(ret, c.CommonFeature.CornerSign.Weights...)
	ret = append(ret, c.CommonFeature.Online.Weights...)
	return
}
