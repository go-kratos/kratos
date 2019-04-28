package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unsafe"

	recsys "go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/dao/parallel"
	"go-common/app/service/bbq/recsys/model"
	"go-common/library/cache/redis"
	"go-common/library/log"

	"github.com/Dai0522/workpool"
	"github.com/json-iterator/go"
)

//user const
const (
	TaskLastPage            = "TaskLastPage"
	TaskLastUpsPage         = "TaskLastUpsPage"
	TaskBiliUserProfile     = "TaskBiliUserProfile"
	TaskBBQUserProfile      = "TaskBBQUserProfile"
	TaskBBQDeviceProfile    = "TaskBBQDeviceProfile"
	TaskUserLike            = "TaskUserLike"
	TaskUserLikeYesterday   = "TaskUserLikeYesterday"
	TaskUserPlay            = "TaskUserPlay"
	TaskUserPlayYesterday   = "TaskUserPlayYesterday"
	TaskDevicePlay          = "TaskDevicePlay"
	TaskDevicePlayYesterday = "TaskDevicePlayYesterday"
	TaskUserFollow          = "TaskUserFollow"
	TaskUserFollowYesterday = "TaskUserFollowYesterday"

	//_BBQDeviceProfileKey = "bbq:device:profile:%s"
	_BBQDeviceProfileKey    = "bbq:device:profile:{buvid}:%s"
	_BBQUserProfileKey      = "bbq:user:profile:%d"
	_BiliUserProfileKey     = "bbq:user:basic:%d"
	_LastFewPageRecords1    = "bbq:last:v1:mid:%d"
	_LastFewPageRecords2    = "bbq:last:v1:buvid:%s"
	_LastFewUpsPageRecords1 = "bbq:last:v1:ups:mid:%d"
	_LastFewUpsPageRecords2 = "bbq:last:v1:ups:buvid:%s"

	_RealTimeUserLike      = "storm:v2:u:%d:like:%s"
	_RealTimeUserPlayMID   = "storm:v2:u:%d:%s:view:100"
	_RealTimeUserPlayBuvID = "storm:v2:u:%s:%s:view:100"
	_RealTimeUserFollow    = "storm:v2:u:%d:%s:follow:100"

	_ModelTest = "bbq:model:init"

	_Zone = "zone"
	_Tag  = "tag"
	_Up   = "up"
)

//LastPageRedisKey for main rec process
func (d *Dao) LastPageRedisKey(mid int64, buvid string) (key string) {
	if mid > 0 {
		key = fmt.Sprintf(_LastFewPageRecords1, mid)
	} else {
		key = fmt.Sprintf(_LastFewPageRecords2, buvid)
	}
	return
}

//LastUpsPageRedisKey for ups rec process
func (d *Dao) LastUpsPageRedisKey(mid int64, buvid string) (key string) {
	if mid > 0 {
		key = fmt.Sprintf(_LastFewUpsPageRecords1, mid)
	} else {
		key = fmt.Sprintf(_LastFewUpsPageRecords2, buvid)
	}
	return
}

//InitModel ...
func (d *Dao) InitModel(c context.Context, weights map[string]float64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	key := _ModelTest
	if result, err := redis.String(conn.Do("GET", key)); err == nil {
		for _, field := range strings.Split(result, ",") {
			featureWeightPair := strings.Split(field, ":")
			if len(featureWeightPair) >= 2 {
				feature := featureWeightPair[0]
				weight, _ := strconv.ParseFloat(featureWeightPair[1], 64)
				weights[feature] = weight
			}
		}
	}
	return
}

//StoreRecResults store rec or upsRec history according to getKeyFunc
func (d *Dao) StoreRecResults(c context.Context, u *model.UserProfile, mid int64, buvid string, response *recsys.RecsysResponse, getKeyFunc func(int64, string) string, lastRecords []model.Record4Dup) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	key := getKeyFunc(mid, buvid)
	maxPageNum := 10
	size := len(response.List)
	if len(lastRecords) > maxPageNum*size {
		lastRecords = lastRecords[size:]
	}

	for _, record := range response.List {
		svid := record.Svid
		mid, ok1 := record.Map[model.UperMid]
		tag, ok2 := record.Map[model.ScatterTag]
		if ok1 && ok2 {
			lastRecords = append(lastRecords, model.Record4Dup{
				SVID: svid,
				MID:  mid,
				Tag:  tag,
			})
		}
	}

	bytes, _ := jsoniter.Marshal(lastRecords)
	_, err = conn.Do("SETEX", key, 86400, bytes)
	if err != nil {
		log.Error("store last few records error: ", err)
	}

	////for test
	//if mid == 28272030 || mid == 390642849 {
	//	return
	//}
	// write bloomfilter for es
	svids := make([]uint64, len(response.List))
	for i, v := range response.List {
		svids[i] = uint64(v.Svid)
	}
	if _, bfErr := d.WriteBF(c, mid, buvid, svids); bfErr != nil {
		log.Errorv(c, log.KV("Write BF error: ", bfErr))
	}
	return
}

//InitUserProfile ...
func (d *Dao) InitUserProfile(c context.Context, mid int64, buvid string) (u *model.UserProfile) {

	u = &model.UserProfile{
		Mid:        mid,
		Buvid:      buvid,
		Name:       "",
		Gender:     -1,
		ViewVideos: []int64{},
		Zones1:     map[string]float64{},
		BiliTags:   map[string]float64{}, //bili
		Zones2:     map[string]float64{}, //bili
		FollowUps:  map[int64]int64{},    //bili

		BBQTags:    map[string]float64{}, //bbq
		BBQZones:   map[string]float64{}, //bbq
		BBQPrefUps: map[int64]int64{},    //bbq

		BBQFollowAction: map[int64]int64{}, //bbq
		BBQFollow:       map[int64]int64{}, //bbq
		BBQBlack:        map[int64]int64{}, //bbq

		PosVideos:   map[int64]int64{},
		NegVideos:   map[int64]int64{},
		LikeVideos:  map[int64]int64{},
		LikeTags:    map[string]float64{},
		LikeTagIDs:  map[int64]int64{},
		LikeUPs:     map[int64]int64{},
		PosTagIDs:   map[int64]int64{},
		NegTagIDs:   map[int64]int64{},
		PosTags:     map[string]float64{},
		NegTags:     map[string]float64{},
		LastRecords: []model.Record4Dup{},
	}
	return
}

//LoadUserProfile load user info from redis parallel
func (d *Dao) LoadUserProfile(c context.Context, mid int64, buvid string) (userProfile *model.UserProfile, err error) {
	tasks := make(map[string]workpool.Task)

	userProfile = d.InitUserProfile(c, mid, buvid)

	// lastPage
	if mid != 0 || buvid != "" {
		taskName := TaskLastPage
		key := fmt.Sprintf(_LastFewPageRecords2, buvid)
		if mid != 0 {
			key = fmt.Sprintf(_LastFewPageRecords1, mid)
		}
		task := parallel.NewRedisTaskWithName(&c, taskName, d.redis, "GET", key)
		tasks[taskName] = task
	}

	if mid != 0 || buvid != "" {
		taskName := TaskLastUpsPage
		key := fmt.Sprintf(_LastFewUpsPageRecords2, buvid)
		if mid != 0 {
			key = fmt.Sprintf(_LastFewUpsPageRecords1, mid)
		}
		task := parallel.NewRedisTaskWithName(&c, taskName, d.redis, "GET", key)
		tasks[taskName] = task
	}

	// user profile bili
	if mid != 0 {
		taskName := TaskBiliUserProfile
		key := fmt.Sprintf(_BiliUserProfileKey, mid)
		task := parallel.NewRedisTaskWithName(&c, taskName, d.redis, "HGETALL", key)
		tasks[taskName] = task
	}

	// user profile bbq: mid
	if mid != 0 {
		taskName := TaskBBQUserProfile
		key := fmt.Sprintf(_BBQUserProfileKey, mid)
		task := parallel.NewRedisTaskWithName(&c, taskName, d.redis, "HGETALL", key)
		tasks[taskName] = task
	}
	// user profile bbq: buvid
	if mid == 0 && buvid != "" {
		taskName := TaskBBQDeviceProfile
		key := fmt.Sprintf(_BBQDeviceProfileKey, buvid)
		task := parallel.NewRedisTask(&c, d.redis, "HGETALL", key)
		tasks[taskName] = task
	}
	// user real time like
	today := time.Now().Format("20060102")
	yesterday := time.Now().AddDate(0, 0, -1).Format("20060102")

	if mid != 0 {
		taskName := TaskUserLike
		key := fmt.Sprintf(_RealTimeUserLike, mid, today)
		task := parallel.NewRedisTask(&c, d.redis, "HGETALL", key)
		tasks[taskName] = task
	}
	if mid != 0 {
		taskName := TaskUserLikeYesterday
		key := fmt.Sprintf(_RealTimeUserLike, mid, yesterday)
		task := parallel.NewRedisTask(&c, d.redis, "HGETALL", key)
		tasks[taskName] = task
	}
	if mid != 0 {
		taskName := TaskUserFollow
		key := fmt.Sprintf(_RealTimeUserFollow, mid, today)
		task := parallel.NewRedisTask(&c, d.redis, "HGETALL", key)
		tasks[taskName] = task
	}
	if mid != 0 {
		taskName := TaskUserFollowYesterday
		key := fmt.Sprintf(_RealTimeUserFollow, mid, yesterday)
		task := parallel.NewRedisTask(&c, d.redis, "HGETALL", key)
		tasks[taskName] = task
	}
	if mid != 0 {
		taskName := TaskUserPlay
		key := fmt.Sprintf(_RealTimeUserPlayMID, mid, today)
		task := parallel.NewRedisTask(&c, d.redis, "HGETALL", key)
		tasks[taskName] = task
	}
	if mid != 0 {
		taskName := TaskUserPlayYesterday
		key := fmt.Sprintf(_RealTimeUserPlayMID, mid, yesterday)
		task := parallel.NewRedisTask(&c, d.redis, "HGETALL", key)
		tasks[taskName] = task
	}
	if mid == 0 && buvid != "" {
		taskName := TaskDevicePlay
		key := fmt.Sprintf(_RealTimeUserPlayBuvID, buvid, today)
		task := parallel.NewRedisTask(&c, d.redis, "HGETALL", key)
		tasks[taskName] = task
	}
	if mid == 0 && buvid != "" {
		taskName := TaskDevicePlayYesterday
		key := fmt.Sprintf(_RealTimeUserPlayBuvID, buvid, yesterday)
		task := parallel.NewRedisTask(&c, d.redis, "HGETALL", key)
		tasks[taskName] = task
	}

	ftTasks := d.parallelTask2(tasks)

	for name, task := range ftTasks {
		var raw *[]byte
		raw, err = task.Wait(100 * time.Millisecond)
		if err != nil && err != redis.ErrNil {
			log.Errorv(c, log.KV("REDIS_GET_ERROR", err))
			continue
		}
		if raw == nil {
			continue
		}

		switch name {
		case TaskLastPage:
			setLastPage(raw, userProfile, "lastRecords")
		case TaskLastUpsPage:
			setLastPage(raw, userProfile, "lastUpsRecords")

		case TaskBiliUserProfile:
			setUserProfileBili(raw, err, userProfile)
		case TaskBBQDeviceProfile:
			setUserProfileBBQ(raw, err, userProfile)
		case TaskBBQUserProfile:
			setUserProfileBBQ(raw, err, userProfile)

		case TaskUserLikeYesterday:
			setUserLikeInfo(raw, err, userProfile)
		case TaskUserLike:
			setUserLikeInfo(raw, err, userProfile)

		case TaskUserFollowYesterday:
			setUserFollowInfo(raw, err, userProfile)
		case TaskUserFollow:
			setUserFollowInfo(raw, err, userProfile)

		case TaskUserPlayYesterday:
			setUserPlayInfo(raw, err, userProfile)
		case TaskDevicePlayYesterday:
			setUserPlayInfo(raw, err, userProfile)
		case TaskUserPlay:
			setUserPlayInfo(raw, err, userProfile)
		case TaskDevicePlay:
			setUserPlayInfo(raw, err, userProfile)
		}
	}

	if err == redis.ErrNil {
		err = nil
	}

	return
}

func setUserProfileBBQ(bytes *[]byte, inErr error, u *model.UserProfile) (err error) {
	var res map[string]string
	if res, err = redis.StringMap(*(*interface{})(unsafe.Pointer(bytes)), inErr); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("redis HGETALL failed error(%v)", err)
		}
	}

	for key, value := range res {
		if key == _Zone {
			zone2s := strings.Split(value, ",")
			for _, zone2 := range zone2s {
				u.BBQZones[zone2] = 1.0
			}
		} else if key == _Tag {
			tags := strings.Split(value, ",")
			for _, tag := range tags {
				u.BBQTags[tag] = 1.0
			}
		} else if key == _Up {
			ups := strings.Split(value, ",")
			for _, upStr := range ups {
				upMID, _ := strconv.ParseInt(upStr, 10, 64)
				u.BBQPrefUps[upMID] = 1
			}
		}
	}
	return
}

func setUserProfileBili(bytes *[]byte, inErr error, u *model.UserProfile) {
	var res map[string]string
	var err error
	if res, err = redis.StringMap(*(*interface{})(unsafe.Pointer(bytes)), inErr); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("redis HGETALL failed error(%v)", err)
		}
	}

	for key, value := range res {
		if key == _Zone {
			zone2s := strings.Split(value, ",")
			for _, zone2 := range zone2s {
				u.Zones2[zone2] = 1.0
			}
		} else if key == _Tag {
			tags := strings.Split(value, ",")
			for _, tag := range tags {
				u.BiliTags[tag] = 1.0
			}
		} else if key == _Up {
			ups := strings.Split(value, ",")
			for _, upStr := range ups {
				upMID, _ := strconv.ParseInt(upStr, 10, 64)
				u.FollowUps[upMID] = 1
			}
		}
	}
}

func setUserLikeInfo(bytes *[]byte, inErr error, u *model.UserProfile) {

	var object struct {
		SVID  int64  `json:"svid"`
		CTime int64  `json:"ctime"`
		BuvID string `json:"buvid"`
	}
	var res map[string]string
	var err error
	if res, err = redis.StringMap(*(*interface{})(unsafe.Pointer(bytes)), inErr); err != nil {
		if err != redis.ErrNil {
			log.Error("redis HGETALL failed error(%v)", err)
		}
	}
	for _, value := range res {
		err = jsoniter.UnmarshalFromString(value, &object)
		if err != nil {
			log.Error("json parse error: %v", err)
		}
		u.LikeVideos[object.SVID] = object.CTime
	}
}

func setUserFollowInfo(bytes *[]byte, inErr error, u *model.UserProfile) {

	var object struct {
		UpID  int64 `json:"upid"`
		CTime int64 `json:"ctime"`
		MID   int64 `json:"mid"`
	}
	var res map[string]string
	var err error
	if res, err = redis.StringMap(*(*interface{})(unsafe.Pointer(bytes)), inErr); err != nil {
		if err != redis.ErrNil {
			log.Error("user real time follow redis HGETALL failed error(%v)", err)
		}
	}
	for _, value := range res {
		err = jsoniter.UnmarshalFromString(value, &object)
		if err != nil {
			log.Error("json parse error: %v", err)
		}
		u.BBQFollowAction[object.UpID] = object.CTime
	}
}

func setUserPlayInfo(bytes *[]byte, inErr error, u *model.UserProfile) {

	var object struct {
		Svid         int64 `json:"svid"`
		CTime        int64 `json:"ctime"`
		Duration     int64 `json:"duration"`
		ViewDuration int64 `json:"viewDuration"`
	}
	var res map[string]string
	var err error
	if res, err = redis.StringMap(*(*interface{})(unsafe.Pointer(bytes)), inErr); err != nil {
		if err != redis.ErrNil {
			log.Error("redis HGETALL failed error(%v)", err)
		} else {
			err = nil
		}
	}

	for _, value := range res {
		err = jsoniter.UnmarshalFromString(value, &object)
		if err != nil {
			log.Error("json parse error: %v", err)
			continue
		}
		u.ViewVideos = append(u.ViewVideos, object.Svid)
		if object.ViewDuration >= 15000 || (object.Duration >= 5000 && float64(object.ViewDuration) >= 0.95*float64(object.Duration)) {
			u.PosVideos[object.Svid] = object.CTime
		}
		if object.ViewDuration <= 500 {
			u.NegVideos[object.Svid] = object.CTime
		}
	}
}

func setLastPage(bytes *[]byte, u *model.UserProfile, lastRecordType string) {
	var results []model.Record4Dup
	if len(*bytes) == 0 {
		return
	}
	err := jsoniter.Unmarshal(*bytes, &results)
	if err != nil {
		log.Error("UnmarshalFromString value(%v) error(%v)", bytes, err)
	} else {
		if lastRecordType == "lastRecords" {
			u.LastRecords = results
		} else {
			u.LastUpsRecords = results
		}

	}
}
