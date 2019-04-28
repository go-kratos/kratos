package v1

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go-common/app/service/live/rtc/common"
	v1pb "go-common/app/service/live/xcaptcha/api/grpc/v1"
	"go-common/app/service/live/xcaptcha/conf"
	"go-common/app/service/live/xcaptcha/dao"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/sync/pipeline/fanout"
	"math/rand"
	"strconv"
	"time"
)

// XCaptchaService struct
type XCaptchaService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao           *dao.Dao
	task          *fanout.Fanout
	captchaLancer *common.LancerLogStream
}

var (
	_geeCaptchaType = int64(1)
	_liveImageType  = int64(0)

	_actionCreate   = 1 // 创建行为
	_actionValidate = 2 // 校验行为

	_responseSuc    = int64(0) // 返回成功
	_responseFailed = int64(1) // 返回失败

	_validateSuc       = int64(0) // 校验成功
	_validateFalid     = int64(1) // 校验失败
	_validateParamsErr = int64(3) // 参数非法

	_xCaptchaCallGeeNum = "xcaptcha_call_gee_time:%d_num:%d" // qps计数，分片减轻热key
	_xCaptchaToken      = "xcaptcha_type:%d_token:%s"        // 二次验证的token缓存
)

// NewXCaptchaService init
func NewXCaptchaService(c *conf.Config) (s *XCaptchaService) {
	s = &XCaptchaService{
		conf:          c,
		dao:           dao.New(c),
		task:          fanout.New("task", fanout.Worker(2), fanout.Buffer(10240)),
		captchaLancer: common.NewLancerLogStream(c.LogStream.Address, c.LogStream.Capacity, time.Duration(c.LogStream.Timeout)),
	}
	return s
}

// xCaptcha 相关服务

// Create implementation
// 创建验证码
func (s *XCaptchaService) Create(ctx context.Context, req *v1pb.XCreateCaptchaReq) (resp *v1pb.XCreateCaptchaResp, err error) {
	// lancer data
	beginTime := time.Now().UnixNano()

	resp = &v1pb.XCreateCaptchaResp{
		Type: req.GetType(),
	}
	resp.Geetest = &v1pb.GeeTest{}
	resp.Image = &v1pb.Image{}

	dataBusM := &dao.PubMessage{
		Uid:     req.GetUid(),
		Ip:      req.GetClientIp(),
		Action:  _actionCreate,
		ReqType: req.GetType(),
		ResType: req.GetType(),
		ResCode: _responseSuc,
	}
	// 校验是否call geeTest
	callGee := false
	var (
		callGeeBegin int64
		callGeeEnd   int64
	)
	degrade := 0
	if req.GetType() == _geeCaptchaType {
		callGee = s.numCheck(ctx, req.GetUid())
		if !callGee { // 触发极验qps限制
			degrade = 1
		}
	}

	// call liveCaptcha
	if req.GetType() == _liveImageType || !callGee {
		resp.Type = _liveImageType
		dataBusM.ResType = _liveImageType
		begin := time.Now().UnixNano()
		resp.Image, err = s.liveCreate(ctx, req.GetWidth(), req.GetHeight())
		end := time.Now().UnixNano()
		if err != nil {
			dataBusM.ResCode = _responseFailed
			log.Error("[xCaptcha][Create][create liveCaptcha error] err:%v", err)
		}
		// 投递消息
		s.syncPush(ctx, dataBusM, beginTime, int(end-begin), degrade)
		return
	}

	// call geeTest
	if req.GetType() == _geeCaptchaType && callGee {
		var (
			geeErr  error
			liveErr error
		)
		groupTimeout := time.Duration(500 * time.Millisecond)
		wg := &errgroup.Group{}
		wg.Go(func() error {
			// call geeTest
			callGeeBegin = time.Now().UnixNano()
			cxtHttp, cancel := context.WithTimeout(ctx, groupTimeout)
			resp.Geetest, geeErr = s.geeCreate(cxtHttp, req.GetClientType(), req.GetClientIp())
			cancel()
			callGeeEnd = time.Now().UnixNano()
			return geeErr
		})
		wg.Go(func() error {
			cxtLive, cancel := context.WithTimeout(ctx, groupTimeout)
			// 创建live captcha
			resp.Image, liveErr = s.liveCreate(cxtLive, req.GetWidth(), req.GetHeight())
			cancel()
			return liveErr
		})
		wg.Wait()

		if geeErr != nil || resp.Geetest.Challenge == "" {
			resp.Geetest = &v1pb.GeeTest{}
			if liveErr != nil || resp.Image.Token == "" {
				dataBusM.ResCode = _responseFailed
				log.Error("[xCaptcha][Create][after geeTest error create liveCaptcha error] err:%v", liveErr)
				err = ecode.Error(-400, "创建验证码失败,请重试~")
			} else {
				degrade = 2 // 极验失败时的降级
				resp.Type = _liveImageType
				dataBusM.ResType = _liveImageType
			}
		} else {
			resp.Image = &v1pb.Image{}
			dataBusM.ResType = _geeCaptchaType
		}
		s.syncPush(ctx, dataBusM, beginTime, int(callGeeEnd-callGeeBegin), degrade)
		return
	}

	return
}

// reqAnti format verify requst
type reqAnti struct {
	Type       int64
	Challenge  string
	Validate   string
	SecCode    string
	ClientType string
	Token      string
	Phrase     string
}

// Check implementation
// 验证码校验 `internal:"true"`
func (s *XCaptchaService) Check(ctx context.Context, req *v1pb.CheckReq) (resp *v1pb.CheckResp, err error) {
	// lance
	beginTime := time.Now().Nanosecond()

	resp = &v1pb.CheckResp{}
	antiStruct := &reqAnti{}

	err = json.Unmarshal([]byte(req.GetAnti()), &antiStruct)
	if err != nil {
		log.Error("[XCaptcha][Verify][unmarshal request error] err: %v", err)
		err = ecode.Error(-400, "校验失败,请重试~")
		return
	}

	if (antiStruct.Type == _liveImageType && (antiStruct.Token == "" || antiStruct.Phrase == "")) ||
		(antiStruct.Type == _geeCaptchaType &&
			(antiStruct.Challenge == "" || antiStruct.SecCode == "" || antiStruct.Validate == "")) {
		err = ecode.Error(-400, "校验失败,请重试~")
		return
	}

	resp.Type = antiStruct.Type

	message := &captchaLancerData{
		uid:      req.GetUid(),
		clientIp: req.GetClientIp(),
		action:   _actionValidate,
		reqType:  antiStruct.Type,
		resType:  antiStruct.Type,
	}
	var (
		callBegin int
		callEnd   int
	)

	if antiStruct.Type == _liveImageType {
		callBegin = time.Now().Nanosecond()
		err = s.liveCheck(ctx, antiStruct.Token, antiStruct.Phrase)
		callEnd = time.Now().Nanosecond()
		if err != nil {
			message.resCode = 1
			message.reqCost = callEnd - beginTime
			message.relyCost = callEnd - callBegin
			s.lancerData(ctx, message)
			return
		}
	}

	if antiStruct.Type == _geeCaptchaType {
		slice := md5.Sum([]byte(s.conf.GeeTest.Key + "geetest" + antiStruct.Challenge))

		if len(antiStruct.Validate) != 32 || antiStruct.Validate != hex.EncodeToString(slice[:]) {
			err = ecode.Error(-400, "校验失败,请重试~")
			return
		}
		callBegin = time.Now().Nanosecond()
		err = s.geeVerify(ctx, antiStruct.Challenge, antiStruct.SecCode, antiStruct.ClientType, req.GetClientIp(), req.GetUid())
		callEnd = time.Now().Nanosecond()
		if err != nil {
			message.resCode = 1
			message.reqCost = callEnd - beginTime
			message.relyCost = callEnd - callBegin
			s.lancerData(ctx, message)
			err = ecode.Error(-400, "校验失败,请重试~")
			return
		}
	}

	resp.Token, err = s.token(ctx, antiStruct.Type, antiStruct.Token)
	if err != nil {
		message.resCode = 1
		err = ecode.Error(-400, "校验失败,请重试~")
	}

	message.reqCost = time.Now().Nanosecond() - beginTime
	message.relyCost = callEnd - callBegin
	s.lancerData(ctx, message)

	return
}

// geeCreate geeCreate
func (s *XCaptchaService) geeCreate(ctx context.Context, clientType string, ip string) (resp *v1pb.GeeTest, err error) {
	resp = &v1pb.GeeTest{}
	id := conf.Conf.GeeTest.Id
	privateKey := conf.Conf.GeeTest.Key
	challenge, err := s.dao.Register(ctx, ip, clientType, 1, id)

	if err != nil || challenge == "" {
		log.Error("[xCaptcha][Create][call geetest error] err:%v", err)
		err = ecode.Error(-400, "创建验证码失败,请重试~")
		return
	}
	slice := md5.Sum([]byte(challenge + privateKey))
	resp.Challenge = hex.EncodeToString(slice[:])
	resp.Gt = id
	return
}

// geeVerify geeTest verify
func (s *XCaptchaService) geeVerify(ctx context.Context, challenge string, seccode string, clientType string, ip string, mid int64) (err error) {
	res, err := s.dao.Validate(ctx, clientType, seccode, clientType, ip, mid, s.conf.GeeTest.Id)
	if err != nil {
		log.Error("[xCaptcha][Verify][call geetest error] err:%v", err)
		err = ecode.Error(1, "校验失败,请重试~")
		return
	}
	seccodeMd5 := md5.Sum([]byte(seccode))
	if res == nil || res.Seccode != hex.EncodeToString(seccodeMd5[:]) {
		log.Error("[xCaptcha][Verify][call geetest error] err:%v, res.secCode:%v, seccode:%v", err, res.Seccode, seccode)
		err = ecode.Error(-400, "校验失败,请重试~")
		return
	}
	return
}

// liveCreate 创建直播验证码
func (s *XCaptchaService) liveCreate(ctx context.Context, width int64, height int64) (resp *v1pb.Image, err error) {
	resp = &v1pb.Image{}
	liveResp, err := s.dao.LiveCreate(ctx, width, height)
	if err != nil || liveResp.Token == "" || liveResp.Image == "" {
		err = ecode.Error(-400, "创建验证码失败,请重试~")
		return
	}
	resp.Tips = "请输入验证码"
	resp.Token = liveResp.Token
	resp.Content = liveResp.Image
	return
}

// liveCheck 校验直播验证码
func (s *XCaptchaService) liveCheck(ctx context.Context, token string, pharse string) (err error) {
	liveResp, err := s.dao.LiveCheck(ctx, token, pharse)
	if err != nil || liveResp != 0 {
		err = ecode.Error(1, "校验失败,请重试~")
		return
	}
	if liveResp != 0 {
		err = ecode.Error(-400, "校验失败,请重试~")
		return
	}
	return
}

// numCheck qps check
func (s *XCaptchaService) numCheck(ctx context.Context, uid int64) (check bool) {
	check = false
	if s.conf.GeeTest.On == 0 {
		return
	}
	qps := s.conf.GeeTest.Qps
	slice := s.conf.GeeTest.Slice
	if qps <= 0 {
		qps = 100
	}
	if slice <= 0 || slice > 10 {
		slice = 5
	}
	each := int64(qps / slice)
	rand.Seed(time.Now().UnixNano())
	slot := rand.Intn(int(slice))
	timestamp := time.Now().Unix()
	key := fmt.Sprintf(_xCaptchaCallGeeNum, timestamp, slot)
	count, err := s.dao.RedisIncr(ctx, key)
	if err != nil {
		return
	}
	if count < each {
		check = true
	}
	return
}

// syncPush
func (s *XCaptchaService) syncPush(ctx context.Context, message *dao.PubMessage, reqBegin int64, relyCost int, degrade int) {
	object := dao.PubMessage{
		Uid:     message.Uid,
		Ip:      message.Ip,
		Action:  message.Action,
		ReqType: message.ReqType,
		ResType: message.ResType,
		ResCode: message.ResCode,
	}
	if reqBegin != 0 {
		// 埋点
		now := time.Now().UnixNano()
		lancerData := &captchaLancerData{
			uid:      message.Uid,
			clientIp: message.Ip,
			action:   message.Action,
			reqType:  message.ReqType,
			resType:  message.ResType,
			resCode:  message.ResCode,
			reqCost:  int(now - reqBegin),
			relyCost: relyCost,
			degrade:  degrade,
		}
		s.lancerData(ctx, lancerData)
	}
	// sync push dataBus
	f := func(ctx context.Context) {
		err := s.dao.Pub(ctx, object)
		if err == nil {
			log.Info("[xCaptcha][PubMessage][success] message :%+v", object)
		}
	}
	sync := s.task.Do(ctx, func(c context.Context) {
		f(c)
	})
	if sync != nil {
		log.Error("[xCaptcha][PubMessage][task Full] err:%v, message:%+v", sync, object)
		f(ctx)
	}
}

// XAnti token校验参数
type XAnti struct {
	Type  int64
	Token string
}

// Verify implementation
// checkToken `internal:"true"`
func (s *XCaptchaService) Verify(ctx context.Context, req *v1pb.XVerifyReq) (resp *v1pb.XVerifyResp, err error) {
	resp = &v1pb.XVerifyResp{}
	dataBusM := &dao.PubMessage{
		RoomId: req.GetRoomId(),
		Uid:    req.GetUid(),
		Ip:     req.GetClientIp(),
		Action: _actionValidate,
	}
	antiBytes, err := base64.URLEncoding.DecodeString(req.GetXAnti())
	if err != nil {
		log.Error("[XCaptcha][TokenCheck][base64decode request error] err: %v", err)
		err = ecode.Error(-400, "验证口令失败~")
		dataBusM.ResCode = _validateParamsErr
		s.syncPush(ctx, dataBusM, 0, 0, 0)
		return
	}
	xAnti := &XAnti{}
	err = json.Unmarshal(antiBytes, &xAnti)
	if err != nil {
		log.Error("[XCaptcha][TokenCheck][unmarshal request error] err: %v", err)
		err = ecode.Error(-400, "验证口令失败~")
		dataBusM.ResCode = _validateParamsErr
		s.syncPush(ctx, dataBusM, 0, 0, 0)
		return
	}
	dataBusM.ReqType = xAnti.Type
	dataBusM.ResType = xAnti.Type

	key := fmt.Sprintf(_xCaptchaToken, xAnti.Type, xAnti.Token)
	value, err := s.dao.RedisGet(ctx, key)

	if err != nil || value == 0 {
		log.Error("[XCaptcha][TokenCheck][Token Wrong] err: %v", err)
		err = ecode.Error(-400, "验证口令失败~")
		dataBusM.ResCode = _validateFalid
		s.syncPush(ctx, dataBusM, 0, 0, 0)
		return
	}
	// 校验成功, 删除token
	err = s.dao.RedisDel(ctx, key)
	if err != nil {
		log.Error("[XCaptcha][TokenCheck][Token Check Over del key] err: %v", err)
		err = ecode.Error(-400, "验证口令失败~")
		return
	}
	dataBusM.ResCode = _validateSuc
	s.syncPush(ctx, dataBusM, 0, 0, 0)
	return
}

func (s *XCaptchaService) token(ctx context.Context, cType int64, base string) (token string, err error) {
	slice := md5.Sum([]byte(strconv.Itoa(int(cType)) + base))
	token = hex.EncodeToString(slice[:])
	err = s.dao.RedisSet(ctx, fmt.Sprintf(_xCaptchaToken, cType, token), 1, 300)
	return
}

// captchaLancerData lancer 上报
type captchaLancerData struct {
	uid      int64
	clientIp string
	action   int
	reqType  int64
	resType  int64
	resCode  int64
	reqCost  int
	relyCost int
	degrade  int
}

func (s *XCaptchaService) lancerData(ctx context.Context, data *captchaLancerData) {
	data.reqCost = int(data.reqCost / 1e6)
	data.relyCost = int(data.relyCost / 1e6)
	ld := s.captchaLancer.NewLancerData(s.conf.LogStream.LogId, s.conf.LogStream.Token)
	ld.PutInt(data.uid)
	ld.PutString(data.clientIp)
	ld.PutInt(int64(data.action))
	ld.PutInt(data.reqType)
	ld.PutInt(data.resType)
	ld.PutInt(data.resCode)
	ld.PutInt(int64(data.reqCost))
	ld.PutInt(int64(data.relyCost))
	ld.PutInt(int64(data.degrade))
	log.Info("[XCaptcha][lancerData][logStream] data:%+v", data)
	if err := ld.Commit(); err != nil {
		log.Error("[XCaptcha][lancerData][logStream] err: %v", err)
	}
	/**
	f := func(ctx context.Context) {
		log.Info("[XCaptcha][lancerData][info] data:%+v", data)
		infocErr := s.captchaLancer.Info(data.uid, data.clientIp, data.action, data.reqType, data.resType, data.resCode, data.reqCost, data.relyCost, data.degrade)
		if infocErr != nil {
			log.Error("[XCaptcha][lancerData][infoc error] err: %v", infocErr)
		}
	}
	sync := s.task.Do(ctx, func(c context.Context) {
		f(c)
	})
	if sync != nil {
		log.Error("[XCaptcha][lancerData][sync error] err: %v", sync)
	}
	**/
}
