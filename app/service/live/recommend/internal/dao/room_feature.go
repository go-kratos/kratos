package dao

import (
	"context"
	"errors"
	"regexp"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"go-common/app/service/live/recommend/internal/conf"
	relationV1 "go-common/app/service/live/relation/api/liverpc/v1"
	roomV1 "go-common/app/service/live/room/api/liverpc/v1"
	roomV2 "go-common/app/service/live/room/api/liverpc/v2"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

var roomFeatureValue atomic.Value
var recDefaultRoomIds atomic.Value

// StartRoomFeatureJob 更新在线房间的特征信息
func StartRoomFeatureJob(c *conf.Config) {
	t := time.Tick(time.Second * 30)
	refreshRoomFeature(context.Background(), c)
	for range t {
		refreshRoomFeature(context.Background(), c)
	}
}

func refreshRoomFeature(ctx context.Context, c *conf.Config) (err error) {
	n := 20
	currentIds, ok := onlineRoomIdValue.Load().(map[int64]struct{})
	if !ok {
		log.Warn("cannot load current online room ids")
		err = errors.New("cannot load current online room ids")
		return
	}
	keys := make([]int64, 0, len(currentIds))
	for k := range currentIds {
		keys = append(keys, k)
	}
	chunkIdsArray := sliceArray(keys, n)
	roomFeatures := map[int64][]int64{}
	var lock sync.Mutex
	var eg errgroup.Group
	for _, tmp := range chunkIdsArray {
		chunkIds := tmp
		eg.Go(func() (err error) {
			resp, err := RoomAPI.V2Room.GetByIds(ctx, &roomV2.RoomGetByIdsReq{Ids: chunkIds})
			if err != nil || resp.GetCode() != 0 {
				log.Error("dao.RoomAPI.V2Room.GetByIds (%v) error(%v) resp(%v)", chunkIds, err, resp)
				return
			}

			resp1, err1 := RoomAPI.V1RoomPendant.GetPendantByIds(ctx, &roomV1.RoomPendantGetPendantByIdsReq{Ids: chunkIds, Type: "mobile_index_badge", Position: 2})
			if err1 != nil || resp1.GetCode() != 0 {
				log.Error("dao.RoomAPI.V1Room.GetPendantByIds (%v) error(%v) resp(%v)", chunkIds, err1, resp1)
				return
			}

			uids := make([]int64, 0, n)
			for _, r := range resp.Data {
				uids = append(uids, r.Uid)
			}
			resp2, err2 := RelationAPI.V1Feed.GetUserFcBatch(ctx, &relationV1.FeedGetUserFcBatchReq{Uids: uids})
			if err2 != nil || resp.GetCode() != 0 {
				log.Error("dao.RelationAPI.V1Relation.GetUserFcBatch (%v) error(%v) resp(%v)", chunkIds, err2, resp2)
				return
			}

			roomPendantInfo := resp1.Data.Result
			fansCountInfo := resp2.Data

			for roomId, r := range resp.Data {
				cornerTag := ""
				fansNum := int64(0)
				if PendantInfo, ok := roomPendantInfo[roomId]; ok && PendantInfo != nil {
					cornerTag = PendantInfo.Value
				}
				if fans, ok := fansCountInfo[r.Uid]; ok {
					fansNum = fans.Fc
				}
				featureVector := createFeature(c, r.AreaV2Id, cornerTag, fansNum, r.Online)
				lock.Lock()
				roomFeatures[roomId] = featureVector
				lock.Unlock()
			}
			return
		})
	}
	eg.Wait()
	roomFeatureValue.Store(roomFeatures)
	//创建默认推荐房间列表
	roomScoreSlice := ScoreSlice{}
	for roomId, vec := range roomFeatures {
		featureVector := make([]int64, len(vec))
		copy(featureVector, vec)
		featureVector[0] = 0
		counter := Counter{roomId: roomId, score: calcScore(makeWeightVec(c), featureVector)}
		roomScoreSlice = append(roomScoreSlice, counter)
	}
	sort.Sort(roomScoreSlice)
	//默认的召回源
	limit := 400
	recDefault := make([]int64, 0, limit)
	for _, counter := range roomScoreSlice {
		limit = limit - 1
		if limit < 0 {
			break
		}
		recDefault = append(recDefault, counter.roomId)
	}
	recDefaultRoomIds.Store(recDefault)
	log.Info("refreshRoomFeature success, total num:%d recDefault_num:%d, recDefault:%+v", len(roomFeatures), len(recDefault), recDefault)
	return
}

//建立房间相关的特征向量
func createFeature(c *conf.Config, areaV2Id int64, cornerTag string, fansNum int64, onlineValue int64) (featureVector []int64) {
	fansMilestone := c.CommonFeature.FansNum.Values
	onlineMilestone := c.CommonFeature.Online.Values
	cornerSignList := c.CommonFeature.CornerSign.Values

	featureVector = append(featureVector, areaV2Id) //分区id, 留待在线计算的时候替换成0,1
	featureVector = append(featureVector, oneHotEncode(fansNum, fansMilestone)...)
	featureVector = append(featureVector, oneHotTextEncode(cornerTag, cornerSignList)...)
	featureVector = append(featureVector, oneHotEncode(onlineValue, onlineMilestone)...)
	return
}

//把slice按大小切成多个等大的小slice(除了最后一块)
func sliceArray(arr []int64, n int) (ret [][]int64) {
	remainder := len(arr) % n
	quotient := (len(arr) - remainder) / n
	num := int(quotient)
	if remainder > 0 {
		num = num + 1
	}
	ret = make([][]int64, 0, num)
	for i := 0; i < num; i++ {
		if i < num-1 {
			ret = append(ret, arr[n*i:n*(i+1)])
		} else {
			ret = append(ret, arr[n*i:])
		}
	}
	return
}

//构建0,1组成的特征向量; 如果x<0, 返回全为0的向量
func int2Slice(x int, n int) []int64 {
	p := make([]int64, n)
	if x < 0 {
		return p
	}
	p[x] = 1
	return p
}

func compAndSet(value int64, vList []int64) int {
	place := 0
	for _, v := range vList {
		if value < v {
			return place
		}
		place = place + 1
	}
	return place
}

func oneHotEncode(value int64, milestone []int64) []int64 {
	place := compAndSet(value, milestone)
	return int2Slice(place, len(milestone)+1)
}

// textList ["", A, B ]
// 如果targetText空或者没匹配到 ret[0] = 1
func oneHotTextEncode(targetText string, textList []string) (ret []int64) {
	place := 0
	ret = make([]int64, len(textList))
	if targetText == "" {
		ret[0] = 1
		return
	}
	for i, text := range textList {
		if text == "" {
			continue
		}
		match, err := regexp.MatchString(text, targetText)
		if err != nil {
			log.Error("oneHotTextEncode regex error " + text)
			place = 0
			break
		}
		if match {
			place = i
			break
		}
	}
	ret[place] = 1
	return
}
