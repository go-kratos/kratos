package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	banclient "go-common/app/service/live/banned_service/api/liverpc"
	ban "go-common/app/service/live/banned_service/api/liverpc/v1"
	danmuku "go-common/app/service/live/broadcast-proxy/api/v1"
	captchaclient "go-common/app/service/live/captcha/api/liverpc"
	captcha "go-common/app/service/live/captcha/api/liverpc/v0"
	giftclient "go-common/app/service/live/gift/api/liverpc"
	gift "go-common/app/service/live/gift/api/liverpc/v1"
	dm "go-common/app/service/live/live-dm/api/grpc/v1"
	v1 "go-common/app/service/live/xlottery/api/grpc/v1"
	filter "go-common/app/service/main/filter/api/grpc/v1"

	"go-common/app/service/live/xlottery/conf"
	"go-common/app/service/live/xlottery/dao"
	"go-common/app/service/live/xlottery/model"
	xuserclient "go-common/app/service/live/xuser/api/grpc/v1"

	"math/rand"
	"strconv"
	"strings"
	"time"

	accountApi "go-common/app/service/main/account/api"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/pipeline/fanout"

	"github.com/pkg/errors"
)

const (
	ttl                = 90                   // 抽奖持续时间 90s
	cacheTTL           = 100                  // 缓存过期时间
	awardNum           = 100                  // 单个节奏风暴发放奖励个数 100
	roomStormInfo      = "room:storm:info:%d" // 房间节奏风暴信息
	stormInfo          = "storm:info:%d"      // 节奏风暴信息
	stormAwardCount    = "storm:award:%d"     // 节奏风暴已经发放奖品数量
	joinFlag           = "storm:join:%d"      // 用户参加标识
	stormGif           = "http://static.hdslb.com/live-static/live-room/images/gift-section/mobilegift/2/jiezou.gif?2017011901"
	stormDanmuLimit    = "storm:danmu:limit:%d:%d" //lotteryid, uid
	stormDanmuInterval = 5                         //用户点击节奏风暴没中奖时发的弹幕的间隔时间

)

var errCodeMsg = map[errCode]string{
	UserNotAuthErr:      "用户未登陆",
	UserCaptchaFail:     "验证码没通过",
	StormExpireErr:      "节奏风暴抽奖过期",
	AlreadyPickErr:      "已经领取奖励",
	StormFullErr:        "领奖人数已经满了",
	StormContentEmpty:   "节奏风暴内容为空",
	JustSupportOneStorm: "同时仅支持一个节奏风暴，请稍后使用本功能",
	BeatNotExsit:        "尚未设置自定义节奏",
	BeatFailureAudit:    "你自定义的节奏未通过审核",
	BeatAuditing:        "你编辑的自定义节奏正在审核中,请稍等",
	WrongTryAgain:       "出错啦，再试试吧",
	BeatNotAudited:      "节奏内容未通过审核，请尝试修改为其他内容哦",
	BeatBanned:          "该自定义内容被屏蔽",
	BeatIsShield:        "看来你的弹幕里有直播间屏蔽词，这节奏怕是带不起来了",
	UserNotVip:          "你现在还不是年费老爷,无法使用自定义节奏",
	MissBeat:            "你错过了奖励，下次要更快一点哦~",
	InnerErr:            "内部错误",
	ParamErr:            "参数错误",
}

type errCode int32

const (
	//UserNotAuthErr 用户未登陆
	UserNotAuthErr errCode = 1005001
	//UserCaptchaFail 验证码没通过
	UserCaptchaFail errCode = 1005002
	//StormExpireErr 节奏风暴抽奖过期
	StormExpireErr errCode = 1005003
	//AlreadyPickErr 已经领取奖励
	AlreadyPickErr errCode = 1005004
	//StormFullErr 领奖人数已经满了
	StormFullErr errCode = 1005005
	//StormContentEmpty 节奏风暴内容为空
	StormContentEmpty errCode = 1005006
	//JustSupportOneStorm 一个房间只支持一个节奏风暴，请稍后
	JustSupportOneStorm errCode = 1005007
	//BeatNotExsit 尚未设置自定义节奏
	BeatNotExsit errCode = 1005008
	//BeatFailureAudit 你自定义的节奏未通过审核
	BeatFailureAudit errCode = 1005009
	//BeatAuditing 你编辑的自定义节奏正在审核中,请稍等
	BeatAuditing errCode = 1005010
	//WrongTryAgain 出错啦，再试试吧
	WrongTryAgain errCode = 1005011
	//BeatNotAudited 节奏内容未通过审核，请尝试修改为其他内容哦
	BeatNotAudited errCode = 1005012
	//BeatBanned 该自定义内容被屏蔽
	BeatBanned errCode = 1005013
	//BeatIsShield 看来你的弹幕里有直播间屏蔽词，这节奏怕是带不起来了
	BeatIsShield errCode = 1005014
	//UserNotVip 你现在还不是年费老爷,无法使用自定义节奏
	UserNotVip errCode = 1005015
	//MissBeat 你错过了奖励，下次要更快一点哦~
	MissBeat errCode = 1005016
	//InnerErr 内部错误
	InnerErr errCode = 1005017
	//ParamErr 参数错误
	ParamErr errCode = 1005018
)

var publicBeats = map[int64]string{
	1:     "前方高能预警，注意这不是演习",
	2:     "我从未见过如此厚颜无耻之人",
	3:     "那万一赢了呢",
	4:     "你们城里人真会玩",
	5:     "左舷弹幕太薄了",
	7:     "要优雅，不要污",
	8:     "我选择狗带",
	9:     "可爱即正义~~",
	10:    "糟了，是心动的感觉！",
	41000: "这个直播间已经被我们承包了！",
	41001: "妈妈问我为什么跪着看直播 w(ﾟДﾟ)w",
	41002: "你们对力量一无所知~(￣▽￣)~",
}

// StormService  节奏风暴服务
type StormService struct {
	c             *conf.Config
	dao           *dao.Dao
	DMClient      dm.DMClient
	AccountClient accountApi.AccountClient
	GiftClient    *giftclient.Client
	BanClient     *banclient.Client
	CaptchaClient *captchaclient.Client
	VipClient     xuserclient.VipClient
	DanmakuClient danmuku.DanmakuClient
	FilterClient  filter.FilterClient
	task          *fanout.Fanout
}

// NewStromService 获得StormService实例并初始化
func NewStromService(c *conf.Config) (s *StormService) {
	dmClient, err := dm.NewClient(nil)
	if err != nil {
		panic(err)
	}
	accountClient, err := accountApi.NewClient(nil)
	if err != nil {
		panic(err)
	}
	xuserClient, err := xuserclient.NewClient(nil)
	if err != nil {
		panic(err)
	}
	danmuku, err := danmuku.NewClient(nil)
	if err != nil {
		panic(err)
	}
	fliterClient, err := filter.NewClient(nil)
	if err != nil {
		panic(err)
	}
	s = &StormService{
		c:             c,
		dao:           dao.New(c),
		DMClient:      dmClient,
		AccountClient: accountClient,
		GiftClient:    giftclient.New(nil),
		BanClient:     banclient.New(nil),
		CaptchaClient: captchaclient.New(nil),
		VipClient:     xuserClient.VipClient,
		DanmakuClient: danmuku,
		FilterClient:  fliterClient,
		task:          fanout.New("storm", fanout.Worker(10), fanout.Buffer(10240)),
	}
	expireLocalCache()
	return s
}

// Ping Service
func (s *StormService) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *StormService) Close() {
	s.dao.Close()
}

// Start 开启节奏风暴抽奖
func (s *StormService) Start(ctx context.Context, req *v1.StartStormReq) (*v1.StartStormResp, error) {
	roomID := req.GetRoomid()
	uID := req.GetUid()
	beatID := req.GetBeatid()
	useShield := req.GetUseShield()
	ruID := req.GetRuid()
	skipExternalCheck := req.GetSkipExternalCheck()
	log.Info("#storm_start# beatid=%d uid=%d ruid=%d roomdid=%d num=%d", beatID, uID, ruID, roomID, req.GetNum())
	if roomID == 0 || uID == 0 || beatID == 0 || ruID == 0 {
		return &v1.StartStormResp{Code: int32(ParamErr), Msg: errCodeMsg[ParamErr]}, nil
	}

	content, ec := s.preCheck(ctx, roomID, uID, beatID, ruID, skipExternalCheck, useShield)
	if ec != 0 {
		return &v1.StartStormResp{Code: int32(ec), Msg: errCodeMsg[ec]}, nil
	}

	//创建 SpecialGift 得到 返回的id
	dbid, err := s.dao.InsertSpecialGift(&model.SpecialGift{
		UID:         req.GetUid(),
		RoomID:      req.GetRoomid(),
		GiftNum:     req.GetNum(),
		GiftID:      39,
		CreateTime:  time.Now(),
		CustomField: fmt.Sprintf(`{"content":"%s"}`, content),
	})
	if err != nil {
		log.Error("insert special gift err:%s", err.Error())
		return &v1.StartStormResp{Code: int32(InnerErr), Msg: errCodeMsg[InnerErr]}, nil
	}
	//节奏风暴对外id, 是数据库中的id拼接上一个随机6位数而成, 是个递增的值
	rand.Seed(time.Now().Unix())
	randid := rand.Intn(999999)
	id := dbid*1000000 + int64(randid)
	newCtx := metadata.WithContext(ctx)
	s.setCache(newCtx, id, req.GetRoomid(), req.GetNum(), content)
	//broadcast message
	m := map[string]interface{}{
		"cmd": "SPECIAL_GIFT",
		"data": map[string]interface{}{
			"39": map[string]interface{}{
				"id":        fmt.Sprintf("%d", id),
				"time":      ttl,
				"hadJoin":   0,
				"num":       req.GetNum(),
				"content":   content,
				"action":    "start",
				"storm_gif": stormGif,
			},
		},
	}
	str, err := json.Marshal(m)
	if err != nil {
		log.Error("json marshal err:%s", err.Error())
		return &v1.StartStormResp{Code: int32(InnerErr), Msg: errCodeMsg[InnerErr]}, nil
	}

	err = s.callBroadCastRoom(ctx, string(str), req.GetRoomid())
	if err != nil {
		return &v1.StartStormResp{Code: int32(InnerErr), Msg: errCodeMsg[InnerErr]}, nil
	}

	log.Info("#storm_start_success# beatid=%d uid=%d ruid=%d roomdid=%d num=%d", beatID, uID, ruID, roomID, req.GetNum())

	return &v1.StartStormResp{Start: &v1.StartData{Time: ttl, Id: id}}, nil
}

// CanStart 检查房间是否能开启节奏风暴抽奖
func (s *StormService) CanStart(ctx context.Context, req *v1.StartStormReq) (*v1.CanStartStormResp, error) {
	log.Info("#storm_canstart# beatid=%d uid=%d ruid=%d roomdid=%d num=%d", req.GetBeatid(), req.GetUid(), req.GetRuid(), req.GetRoomid(), req.GetNum())
	_, ec := s.preCheck(ctx, req.GetRoomid(), req.GetUid(), req.GetBeatid(), req.GetRuid(), req.GetSkipExternalCheck(), req.GetUseShield())
	if ec != 0 {
		return &v1.CanStartStormResp{Code: int32(ec), Msg: errCodeMsg[ec]}, nil
	}
	return &v1.CanStartStormResp{RetStatus: true}, nil
}

// Join 用户参加节奏风暴抽奖
func (s *StormService) Join(ctx context.Context, req *v1.JoinStormReq) (resp *v1.JoinStormResp, err error) {
	var m map[string]string
	// 本来defer后面的逻辑是在controller层中做的，go的话把这块代码移到网关层就不太好

	defer func() {
		//resp不为nil 说明抽奖成功 ，发送弹幕
		if m != nil {

			if resp.GetJoin() != nil {
				if err1 := s.task.Do(metadata.WithContext(ctx), func(ctx context.Context) {
					s.sendDamu(ctx, req.GetRoomid(), req.GetMid(), m["content"], req.GetPlatform())
				}); err1 == fanout.ErrFull {
					log.Info("sendDamu_task_is_full roomid= %d ", req.GetRoomid())
					s.sendDamu(ctx, req.GetRoomid(), req.GetMid(), m["content"], req.GetPlatform())
				}
			} else {
				//抽奖不成功也要发送弹幕，概率20% 造成一种很多人中奖的假象
				rand.Seed(time.Now().Unix())
				if s.canSendDamu(ctx, req.GetId(), req.GetMid()) && rand.Intn(5) == 2 {
					if err1 := s.task.Do(metadata.WithContext(ctx), func(ctx context.Context) {
						s.sendDamu(ctx, req.GetRoomid(), req.GetMid(), m["content"], req.GetPlatform())
					}); err1 == fanout.ErrFull {
						log.Info("sendDamu_task_is_full roomid= %d ", req.GetRoomid())
						s.sendDamu(ctx, req.GetRoomid(), req.GetMid(), m["content"], req.GetPlatform())
					}
				}
			}
		}

	}()

	mid := req.GetMid()
	platform := req.GetPlatform()

	if req.GetId() != 0 {
		m = s.getStormByID(ctx, req.GetId())
	} else {
		m = s.getByRoomID(ctx, req.GetRoomid(), mid)
	}
	if m == nil {
		return &v1.JoinStormResp{Code: int32(MissBeat), Msg: errCodeMsg[MissBeat]}, nil
	}
	//过期了
	if t, _ := strconv.ParseInt(m["time"], 10, 64); t < time.Now().Unix() {
		newCtx := metadata.WithContext(ctx)
		s.delCache(newCtx, req.GetId(), req.GetRoomid())
		m = nil
		return &v1.JoinStormResp{Code: int32(StormExpireErr), Msg: errCodeMsg[StormExpireErr]}, nil
	}
	roomid, _ := strconv.ParseInt(m["roomid"], 10, 64)
	req.Roomid = roomid

	//验证码
	notMobile := !strings.Contains("android,ios", strings.ToLower(platform))
	if !s.isVerity(ctx, mid) && notMobile {
		if isv, _ := s.callCaptcha(ctx, req.GetCaptchaToken(), req.GetCaptchaPhrase()); !isv {
			m = nil
			return &v1.JoinStormResp{Code: int32(UserCaptchaFail), Msg: errCodeMsg[UserCaptchaFail]}, nil
		}
	}

	stormID, err := strconv.ParseInt(m["id"], 10, 64)
	if err != nil {
		m = nil
		return &v1.JoinStormResp{Code: int32(MissBeat), Msg: errCodeMsg[MissBeat]}, nil
	}

	joinKey := fmt.Sprintf(joinFlag, stormID)
	exist, _ := s.dao.SIsMember(ctx, joinKey, fmt.Sprintf("%d", mid))
	if exist {
		m = nil
		return &v1.JoinStormResp{Code: int32(AlreadyPickErr), Msg: errCodeMsg[AlreadyPickErr]}, nil
	}
	awardKey := fmt.Sprintf(stormAwardCount, stormID)
	count, _ := s.dao.GetInt64(ctx, awardKey)
	num, _ := strconv.ParseInt(m["num"], 10, 64)
	if count > num {
		return &v1.JoinStormResp{Code: int32(MissBeat), Msg: errCodeMsg[MissBeat]}, nil
	}
	//10秒钟放出10个奖励，防止一下子就被抽完了
	coldNum := num / 10
	coldKey := fmt.Sprintf("storm_cold:%d:%d", stormID, time.Now().Unix())
	finalcoldNum := s.dao.Incr(ctx, coldKey)
	s.dao.Expire(ctx, coldKey, 10)
	if int64(finalcoldNum) > coldNum {
		return &v1.JoinStormResp{Code: int32(MissBeat), Msg: errCodeMsg[MissBeat]}, nil
	}

	countIncr := s.dao.Incr(ctx, awardKey)
	s.dao.SAdd(ctx, joinKey, strconv.Itoa(int(mid)))
	s.dao.Expire(ctx, joinKey, cacheTTL)

	year, month, day := time.Now().Date()
	ts := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location()).Add(7 * 86400 * time.Second).Unix()

	if err := s.task.Do(metadata.WithContext(ctx), func(ctx context.Context) {
		s.callGiftAddFreeGift(ctx, mid, 1, 6 /*亿圆*/, ts)
	}); err == fanout.ErrFull {
		log.Info("add_free_gift_task_is_full mid= %d ", mid)
		s.callGiftAddFreeGift(ctx, mid, 1, 6 /*亿圆*/, ts)
	}

	if int64(countIncr) >= num {
		jsonMap := map[string]interface{}{
			"cmd": "SPECIAL_GIFT",
			"data": map[string]interface{}{
				"39": map[string]interface{}{
					"id":     stormID,
					"action": "end",
				},
			},
		}
		jsonByte, _ := json.Marshal(jsonMap)
		err := s.callBroadCastRoom(ctx, string(jsonByte), roomid)
		if err != nil {
			log.Error("send_end_broadcase_err:%s", err.Error())
		}
		s.delCache(metadata.WithContext(ctx), stormID, roomid)
	}

	log.Info("#storm_join_success#%v#%v#%v", stormID, req.GetMid(), countIncr)
	resp = &v1.JoinStormResp{Join: &v1.JoinData{
		GiftId:        39,
		Title:         "节奏风暴",
		Content:       fmt.Sprintf("<p>你是前 %d 位跟风大师<br />恭喜你获得一个亿圆(7天有效期)</p>", countIncr),
		MobileContent: fmt.Sprintf("你是前 %d 位跟风大师", countIncr),
		GiftImg:       "http://static.hdslb.com/live-static/live-room/images/gift-section/gift-6.png?2017011901",
		GiftNum:       1,
		GiftName:      "亿圆",
	}}
	return resp, nil
}

// Check 检查用户能否加入
func (s *StormService) Check(ctx context.Context, req *v1.CheckStormReq) (*v1.CheckStormResp, error) {
	m := s.getByRoomID(ctx, req.GetRoomid(), req.GetUid())
	if m == nil {
		return &v1.CheckStormResp{
			Code: int32(StormExpireErr),
			Msg:  errCodeMsg[StormExpireErr],
		}, nil
	}
	return mapTOCheckStormResp(m), nil
}

func mapTOCheckStormResp(m map[string]string) *v1.CheckStormResp {
	ret := &v1.CheckData{}
	for k, v := range m {
		switch k {
		case "id":
			ret.Id, _ = strconv.ParseInt(v, 10, 64)
		case "roomid":
			ret.Roomid, _ = strconv.ParseInt(v, 10, 64)
		case "num":
			ret.Num, _ = strconv.ParseInt(v, 10, 64)
		case "send_num":
			ret.SendNum = v
		case "time":
			ret.Time, _ = strconv.ParseInt(v, 10, 64)
			ret.Time = ret.Time - time.Now().Unix()
		case "content":
			ret.Content = v
		case "hasJoin":
			r, _ := strconv.ParseInt(v, 10, 32)
			ret.HasJoin = int32(r)
		case "storm_gift":
			ret.StormGif = v
		}

	}

	return &v1.CheckStormResp{Check: ret}

}

// preCheck 返回的error！=nil 表示 check 失败
func (s *StormService) preCheck(ctx context.Context, roomid, uid, beatid, ruid, skipExternalCheck int64, userShield bool) (content string, ec errCode) {
	m := s.getByRoomID(ctx, roomid, uid)
	if m != nil {
		ec = JustSupportOneStorm
		return
	}
	content, ec = s.getBeat(ctx, beatid, uid, ruid, skipExternalCheck, userShield)
	return

}

func (s *StormService) getBeat(ctx context.Context, beatid, uid, ruid, skipExternalCheck int64, userShield bool) (string, errCode) {

	if v, isexit := publicBeats[beatid]; isexit {
		content := v
		if userShield {
			if ec := s.checkShield(ctx, uid, content); ec != 0 {
				return "", ec
			}
		}
		return content, 0

	}
	if skipExternalCheck != 0 {
		//判断账户是不是vip，不是vip就return 无法使用自定义
		resp, err := s.VipClient.Info(ctx, &xuserclient.UidReq{Uid: uid})
		if err != nil {
			return "", InnerErr
		}
		if resp.GetInfo().GetSvip() <= 0 {
			return "", UserNotVip
		}
	}
	beat, err := s.dao.FindBeatByBeatIDAndUID(beatid, uid)
	if err != nil {
		log.Error("find_beat_by_beatid = %d and uid %d  err :%s ", beatid, uid, err.Error())
		return "", BeatNotExsit
	}
	if beat.Content == "" {
		return "", StormContentEmpty
	}
	if beat.Status == -1 {
		return "", BeatFailureAudit
	}
	if beat.Status == 1 {
		return "", BeatAuditing
	}

	if skipExternalCheck != 0 {
		resp, err := s.FilterClient.Filter(ctx, &filter.FilterReq{
			Area:    "live_danmu",
			Message: beat.Content,
		})
		if err != nil {
			return "", InnerErr
		}
		if resp.GetLevel() > 15 { // level 值大于15 就要被过滤
			return "", BeatNotAudited
		}
		log.Info("filter_pass ")

	}
	if userShield {
		if ec := s.checkShield(ctx, uid, beat.Content); ec != 0 {
			return "", ec
		}
	}
	return beat.Content, 0
}

func (s *StormService) checkShield(ctx context.Context, uid int64, content string) errCode {
	resp, err := s.BanClient.V1Shield.IsShieldContent(ctx, &ban.ShieldIsShieldContentReq{
		Uid:     uid,
		Content: content,
	})
	if err != nil {
		log.Error("call_shield_uid = %d content= %s  fail err: %s", uid, content, err.Error())
		return InnerErr
	}
	if resp.GetData().GetIsShieldContent() {
		log.Info("call_shield_IsShieldContent is true ")
		return BeatBanned
	}
	log.Info("call_shield_success ")
	return 0
}

// getByRoomID 获取正在进行的节奏风暴信息和当前用户的参加信息
func (s *StormService) getByRoomID(ctx context.Context, roomID int64, uid int64) map[string]string {
	ret, err := s.dao.HGetAll(ctx, fmt.Sprintf(roomStormInfo, roomID))
	if err != nil {
		if err != dao.ErrEmptyMap {
			log.Info("redis_hgetall_roomid = %d uid = %d  err : %s", roomID, uid, err.Error())
		}
		return nil
	}
	var id, t int64
	if id, err = strconv.ParseInt(ret["id"], 10, 64); err != nil {
		log.Info("ParseInt ret[id] err : %s", err.Error())
		return nil
	}
	if t, err = strconv.ParseInt(ret["time"], 10, 64); err != nil {
		log.Info("ParseInt ret[id] err : %s", err.Error())
		return nil
	}
	// 过期了
	if t < time.Now().Unix() {
		log.Info("the_roomid %d beat_storm_is_expired ", roomID)
		roomid, _ := strconv.ParseInt(ret["roomid"], 10, 64)
		s.delCache(ctx, id, roomid)
		return nil
	}
	ret["hasJoin"] = "0"
	joinKey := fmt.Sprintf(joinFlag, id)
	if joined, _ := s.dao.SIsMember(ctx, joinKey, strconv.Itoa(int(uid))); joined {
		ret["hasJoin"] = "1"
	}

	ret["storm_gift"] = stormGif

	return ret
}

var localCache map[int64]map[string]string
var localCacheLock sync.RWMutex

func init() {
	localCache = make(map[int64]map[string]string)
}

func expireLocalCache() {
	go func() {
		ticker := time.NewTicker(time.Second * 50)
		for range ticker.C {
			localCacheLock.Lock()
			for k, v := range localCache {
				if t, _ := strconv.ParseInt(v["time"], 10, 64); t < time.Now().Unix() {
					delete(localCache, k)
				}
			}
			localCacheLock.Unlock()
		}
	}()

}
func (s *StormService) getStormByID(ctx context.Context, id int64) map[string]string {
	localCacheLock.RLock()
	m, ok := localCache[id]
	localCacheLock.RUnlock()
	if ok {
		log.Info("load_data_from_localcache_id = %d", id)
		return m
	}

	m, _ = s.dao.HGetAll(ctx, fmt.Sprintf(stormInfo, id))
	if m != nil {
		localCacheLock.Lock()
		localCache[id] = m
		localCacheLock.Unlock()
		log.Info("load_data_from_redis_id = %d", id)
		return m
	}
	return nil
}
func (s *StormService) setCache(ctx context.Context, id, roomID, num int64, content string) {
	sr := map[string]interface{}{
		"id":       id,
		"roomid":   roomID,
		"num":      num * awardNum, //奖品数量
		"send_num": num,            //节奏风暴数量
		"time":     time.Now().Unix() + ttl,
		"content":  content,
	}

	roomKey := fmt.Sprintf(roomStormInfo, roomID)
	s.dao.HMSet(ctx, roomKey, sr)
	s.dao.Expire(ctx, roomKey, ttl)

	stormKey := fmt.Sprintf(stormInfo, id)
	s.dao.HMSet(ctx, stormKey, sr)
	s.dao.Expire(ctx, stormKey, ttl)

	awardKey := fmt.Sprintf(stormAwardCount, id)
	s.dao.SetEx(ctx, awardKey, 0, ttl)

}

func (s *StormService) delCache(ctx context.Context, _id, roomID int64) {

	roomKey := fmt.Sprintf(roomStormInfo, roomID)
	s.dao.Del(ctx, roomKey)

	stormKey := fmt.Sprintf(stormInfo, _id)
	s.dao.Del(ctx, stormKey)

	awardKey := fmt.Sprintf(stormAwardCount, _id)
	s.dao.Expire(ctx, awardKey, ttl)

	joinFlagKey := fmt.Sprintf(joinFlag, _id)
	s.dao.Del(ctx, joinFlagKey)
}

func (s *StormService) isVerity(ctx context.Context, uid int64) bool {
	key := fmt.Sprintf("user:identification:%d", uid)
	ret, _ := s.dao.Get(ctx, key)
	if ret != "" {
		return true
	}
	pr, err := s.AccountClient.Profile3(ctx, &accountApi.MidReq{Mid: uid})
	if err != nil {
		return false
	}

	if pr.GetProfile().GetIdentification() == 1 {
		s.dao.SetEx(ctx, key, 1, 3600)
		return true
	}
	return false
}

// 发送room弹幕
func (s *StormService) callBroadCastRoom(ctx context.Context, jsonStr string, roomID int64) error {
	_, err := s.DanmakuClient.RoomMessage(ctx, &danmuku.RoomMessageRequest{
		RoomId:  int32(roomID),
		Ensure:  0,
		Message: jsonStr,
	})
	if err != nil {
		log.Error("call_DanmakuClient_roomid = %d  err :%s", roomID, err.Error())
		return err
	}
	log.Info("call_DanmakuClient_success_roomid=%d message=%s ", roomID, jsonStr)
	return nil

}

func (s *StormService) callGiftAddFreeGift(ctx context.Context, uid, num, giftid, expireat int64) error {
	_, err := s.GiftClient.V1Gift.AddFreeGift(ctx, &gift.GiftAddFreeGiftReq{
		Uid:      uid,
		Num:      num,
		Giftid:   giftid,
		Expireat: expireat,
	})
	if err != nil {
		log.Error("call_gift_liverpc_AddFreeGift_err: %s uid=%d num=%d giftid=%d expireAt=%d", err.Error(), uid, num, giftid, expireat)
		return err
	}
	log.Info("call_gift_liverpc_AddFreeGift_success uid=%d num=%d giftid=%d expireAt=%d", uid, num, giftid, expireat)
	return nil

}

// 验证码验证
func (s *StormService) callCaptcha(ctx context.Context, token, phrase string) (bool, error) {
	resp, err := s.CaptchaClient.V0Captcha.Check(ctx, &captcha.CaptchaCheckReq{
		Token:  token,
		Phrase: phrase,
	})
	if err != nil {
		log.Error("call_captcha_err:%s  token=%s ,phrase=%s", err.Error(), token, phrase)
		return false, err
	}
	if resp.Code != 0 {
		log.Info("captcha_fail_token=%s ,phrase=%s", token, phrase)
		return false, errors.New(resp.Msg)
	}
	log.Info("captcha_success_token=%s ,phrase=%s", token, phrase)
	return true, nil

}

// 用户抽奖成功了发弹幕
func (s *StormService) sendDamu(ctx context.Context, roomid, uid int64, content, platform string) error {
	if _, err := s.DMClient.SendMsg(ctx, &dm.SendMsgReq{
		Uid:      uid,
		Roomid:   roomid,
		Msg:      content,
		Msgtype:  1,
		Platform: platform,
		Mode:     1,
		Fontsize: 25,
		Ip:       metadata.String(ctx, metadata.RemoteIP),
	}); err != nil {
		log.Error("send_damu_err: %s", err.Error())
		return err
	}
	log.Info("#send_damu_success#content= %s roomid=%d uid=%d", content, roomid, uid)
	return nil
}

// 检查是否可以发送弹幕
func (s *StormService) canSendDamu(ctx context.Context, id, uid int64) bool {
	key := fmt.Sprintf(stormDanmuLimit, id, uid)
	_, err := s.dao.SetWithNxEx(ctx, key, "1", stormDanmuInterval)
	return err == nil

}
