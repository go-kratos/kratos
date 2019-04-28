package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/main/push/dao"
	"go-common/app/service/main/push/dao/apns2"
	"go-common/app/service/main/push/dao/fcm"
	"go-common/app/service/main/push/dao/huawei"
	"go-common/app/service/main/push/dao/jpush"
	"go-common/app/service/main/push/dao/oppo"
	"go-common/app/service/main/push/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"
)

const (
	_pushLimitAndroid = 1000
)

func (s *Service) pushAPNSproc() {
	for {
		v, ok := <-s.apnsCh
		if !ok {
			log.Info("apnsCh has been closed.")
			return
		}
		s.pushIOS(v.Info, v.Item)
		s.setChCounter(v.Info.TaskID, -1)
		time.Sleep(time.Millisecond)
	}
}

func (s *Service) pushMiproc() {
	for {
		v, ok := <-s.miCh
		if !ok {
			log.Info("miCh has been closed.")
			return
		}
		for scheme, items := range dispatchByScheme(v.Info.LinkType, v.Info.LinkValue, v.Items) {
			s.pushMi(v.Info, scheme, items)
		}
		s.setChCounter(v.Info.TaskID, -1)
		time.Sleep(time.Millisecond)
	}
}

func (s *Service) pushHuaweiproc() {
	for {
		v, ok := <-s.huaweiCh
		if !ok {
			log.Info("huaweiCh has been closed.")
			return
		}
		for scheme, items := range dispatchByScheme(v.Info.LinkType, v.Info.LinkValue, v.Items) {
			s.pushHuawei(v.Info, scheme, items)
		}
		s.setChCounter(v.Info.TaskID, -1)
		time.Sleep(time.Millisecond)
	}
}

func (s *Service) pushOppoproc() {
	for {
		v, ok := <-s.oppoCh
		if !ok {
			log.Info("oppoCh has been closed.")
			return
		}
		s.pushOppoOne(v.Info, v.Item)
		s.setChCounter(v.Info.TaskID, -1)
		time.Sleep(time.Millisecond)
	}
}

func (s *Service) pushJpushproc() {
	for {
		v, ok := <-s.jpushCh
		if !ok {
			log.Info("jpushCh has been closed.")
			return
		}
		for scheme, items := range dispatchByScheme(v.Info.LinkType, v.Info.LinkValue, v.Items) {
			s.pushJpush(v.Info, scheme, items)
		}
		s.setChCounter(v.Info.TaskID, -1)
		time.Sleep(time.Millisecond)
	}
}

func (s *Service) pushFCMproc() {
	for {
		v, ok := <-s.fcmCh
		if !ok {
			log.Info("fcmCh has been closed")
			return
		}
		for scheme, items := range dispatchByScheme(v.Info.LinkType, v.Info.LinkValue, v.Items) {
			s.pushFCM(v.Info, scheme, items)
		}
		s.setChCounter(v.Info.TaskID, -1)
		time.Sleep(time.Millisecond)
	}
}

func dispatchByScheme(linkType int8, linkValue string, items []*model.PushItem) (res map[string][]*model.PushItem) {
	var scheme string
	res = make(map[string][]*model.PushItem)
	for _, item := range items {
		scheme = model.Scheme(linkType, linkValue, model.PlatformAndroid, item.Build)
		res[scheme] = append(res[scheme], item)
	}
	return
}

func (s *Service) pushInfo(task *model.Task) *model.PushInfo {
	info := &model.PushInfo{
		Job:         task.Job,
		APPID:       task.APPID,
		TaskID:      task.ID,
		Title:       task.Title,
		Summary:     task.Summary,
		LinkType:    task.LinkType,
		LinkValue:   task.LinkValue,
		PushTime:    task.PushTime,
		ExpireTime:  task.ExpireTime,
		Sound:       task.Sound,
		Vibration:   task.Vibration,
		PassThrough: s.c.Push.PassThrough,
		// PassThrough: task.PassThrough,
		ImageURL: task.ImageURL,
	}
	if info.Title == "" {
		info.Title = model.DefaultMessageTitle
	}
	return info
}

// Pushs push some mids.
func (s *Service) Pushs(c context.Context, task *model.Task, mids []int64) (err error) {
	if len(mids) == 0 {
		log.Warn("s.Pushs(%d) no mids", task.ID)
		dao.PromInfo("push:没有mid")
		return
	}
	s.pushByPart(c, task, mids)
	return
}

// SinglePush push one mid.
func (s *Service) SinglePush(c context.Context, token string, task *model.Task, mid int64) (err error) {
	if err = s.checkBusiness(task.BusinessID, token); err != nil {
		return
	}
	var rs map[int][]*model.PushItem
	if rs, err = s.tokensByMid(task, mid); err != nil {
		return
	}
	if len(rs) == 0 {
		log.Warn("no tokens. mid(%d) task(%+v)", mid, task)
		return
	}
	info := s.pushInfo(task)
	for p, v := range rs {
		s.pushByPlatform(info, p, v)
	}
	for {
		if s.chCounterVal(task.ID) == 0 {
			log.Info("Pushs done. task(%+v) result(%+v)", task, s.fetchProgress(task.ID))
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// TestToken for test via push token.
func (s *Service) TestToken(c context.Context, info *model.PushInfo, token string) (err error) {
	r, err := s.dao.Report(context.Background(), token)
	if err != nil {
		return
	}
	if r == nil {
		log.Warn("test token(%s) not exist in db", token)
		return
	}
	var (
		res  *model.HTTPResponse
		item = &model.PushItem{
			Platform:  r.PlatformID,
			Mid:       r.Mid,
			Token:     token,
			Sound:     info.Sound,
			Vibration: info.Vibration,
			Build:     r.Build,
		}
		androidScheme = model.Scheme(info.LinkType, info.LinkValue, model.PlatformAndroid, item.Build)
	)
	switch item.Platform {
	case model.PlatformIPhone, model.PlatformIPad:
		if item.Platform == model.PlatformIPhone {
			res, err = s.dao.PushIPhone(c, info, item)
		} else {
			res, err = s.dao.PushIPad(c, info, item)
		}
		if err == nil && res.Code != apns2.StatusCodeSuccess {
			err = errors.New(res.Msg)
		}
	case model.PlatformXiaomi:
		res, err = s.dao.PushMi(c, info, androidScheme, token)
		if err == nil && res.Code != model.HTTPCodeOk {
			err = errors.New(res.Msg)
		}
	case model.PlatformHuawei:
		var resp *huawei.Response
		resp, err = s.dao.PushHuawei(c, info, androidScheme, []string{token})
		if err == nil && resp.Code != huawei.ResponseCodeSuccess {
			err = errors.New(resp.Msg)
		}
		log.Info("huawei s.TestToken(%+v,%s) result(%+v)", info, token, resp)
	case model.PlatformOppo:
		err = s.pushOppoOne(info, item)
		log.Info("oppo pushOne s.TestToken(%+v,%s)", info, token)
	case model.PlatformJpush:
		_, err = s.dao.PushJpush(c, info, androidScheme, []string{token})
	case model.PlatformFCM:
		_, err = s.dao.PushFCM(c, info, androidScheme, []string{token})
		log.Info("fcm s.dao.PushFCM(%+v,%s)", info, token)
	default:
		err = errors.New("平台类型错误")
	}
	if err != nil {
		log.Error("s.TestToken(%+v,%s) error(%v)", info, token, err)
		return
	}
	log.Info("s.TestToken(%+v,%s) result(%+v)", info, token, res)
	return
}

func (s *Service) pushByPart(c context.Context, task *model.Task, mids []int64) (err error) {
	var (
		counter int
		group   = errgroup.Group{}
		info    = s.pushInfo(task)
	)
	for {
		l := len(mids)
		if l == 0 {
			break
		}
		n := s.c.Push.PushPartSize
		if l < n {
			n = l
		}
		part := mids[:n]
		mids = mids[n:]
		group.Go(func() error {
			var rs map[int][]*model.PushItem
			var missed []int64
			if rs, missed, err = s.tokensByMids(c, task, part); err != nil {
				return nil
			}
			log.Info("missed mid task(%s) %s", task.ID, xstr.JoinInts(missed))
			for p, v := range rs {
				s.pushByPlatform(info, p, v)
			}
			_ = missed
			return nil
		})
		time.Sleep(time.Duration(s.c.Push.PushPartInterval))
		counter++
		if counter > s.c.Push.PushPartChanSize {
			group.Wait()
			counter = 0
		}
	}
	if counter > 0 {
		group.Wait()
	}
	for {
		if s.chCounterVal(task.ID) == 0 {
			log.Info("Pushs done. task(%+v) result(%+v)", task, s.fetchProgress(task.ID))
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (s *Service) pushByPlatform(info *model.PushInfo, platform int, rs []*model.PushItem) {
	switch platform {
	case model.PlatformIPhone, model.PlatformIPad:
		for _, v := range rs {
			s.apnsCh <- &model.PushChanItem{Info: info, Item: v}
			dao.PromChanLen("apns_chan_len", int64(len(s.apnsCh)))
			s.setChCounter(info.TaskID, 1)
		}
		return
	case model.PlatformOppo:
		for _, v := range rs {
			s.oppoCh <- &model.PushChanItem{Info: info, Item: v}
			dao.PromChanLen("oppo_chan_len", int64(len(s.oppoCh)))
			s.setChCounter(info.TaskID, 1)
		}
		return
	}
	n := _pushLimitAndroid
	if platform == model.PlatformHuawei {
		n = s.c.Android.PushHuaweiPart
	}
	for len(rs) > 0 {
		if n > len(rs) {
			n = len(rs)
		}
		part := rs[:n]
		rs = rs[n:]
		switch platform {
		case model.PlatformHuawei:
			s.huaweiCh <- &model.PushChanItems{Info: info, Items: part}
			dao.PromChanLen("huawei_chan_len", int64(len(s.huaweiCh)))
		case model.PlatformJpush:
			s.jpushCh <- &model.PushChanItems{Info: info, Items: part}
			dao.PromChanLen("jpush_chan_len", int64(len(s.jpushCh)))
		case model.PlatformFCM:
			s.fcmCh <- &model.PushChanItems{Info: info, Items: part}
			dao.PromChanLen("fcm_chan_len", int64(len(s.fcmCh)))
		default:
			s.miCh <- &model.PushChanItems{Info: info, Items: part}
			dao.PromChanLen("mi_chan_len", int64(len(s.miCh)))
		}
		s.setChCounter(info.TaskID, 1)
	}
}

func (s *Service) pushMi(info *model.PushInfo, scheme string, items []*model.PushItem) (err error) {
	all := len(items)
	// ts := make([]string, all)
	tokenBuf := bytes.Buffer{}
	for _, item := range items {
		tokenBuf.WriteString(item.Token)
		tokenBuf.WriteString(",")
		// ts = append(ts, item.Token)
	}
	if tokenBuf.Len() == 0 {
		return nil
	}
	tokenBuf.Truncate(tokenBuf.Len() - 1)
	tokens := tokenBuf.String()
	ctx := context.Background()
	var res *model.HTTPResponse
	for c := 0; c <= s.c.Push.RetryTimes; c++ {
		res, err = s.dao.PushMi(ctx, info, scheme, tokens)
		if err == nil {
			break
		}
		dao.PromInfo("retry push mi")
		s.setProgress(info.TaskID, _pgRetryTimes, 1)
	}
	if err != nil || res.Code != model.HTTPCodeOk {
		s.setProgress(info.TaskID, _pgTokenFailed, int64(all))
		return
	}
	s.setProgress(info.TaskID, _pgTokenSuccess, int64(all))
	s.cache.Save(func() { logPushed(info.TaskID, model.PlatformXiaomi, items) })
	var success int
	msg := strings.Split(res.Msg, " ")
	if len(msg) == 6 {
		success, _ = strconv.Atoi(msg[4])
	} else {
		log.Warn("push mi result msg(%s)", res.Msg)
	}
	s.setProgress(info.TaskID, _pgTokenValid, int64(success))
	return
}

func (s *Service) pushIOS(info *model.PushInfo, item *model.PushItem) (err error) {
	var (
		res *model.HTTPResponse
		ctx = context.Background()
	)
	for c := 0; c <= s.c.Push.RetryTimes; c++ {
		if item.Platform == model.PlatformIPhone {
			res, err = s.dao.PushIPhone(ctx, info, item)
		} else {
			res, err = s.dao.PushIPad(ctx, info, item)
		}
		if err == nil {
			break
		}
		dao.PromInfo("retry push iOS")
		s.setProgress(info.TaskID, _pgRetryTimes, 1)
	}
	if err != nil || res == nil {
		s.setProgress(info.TaskID, _pgTokenFailed, 1)
		return
	}
	s.setProgress(info.TaskID, _pgTokenSuccess, 1)
	if res.Code == apns2.StatusCodeSuccess {
		s.setProgress(info.TaskID, _pgTokenValid, 1)
	} else if res.Code == apns2.StatusCodeNoActive || res.Code == apns2.StatusCodeNotForTopic {
		log.Warn("invalid token. mid(%d) token(%s) response(%+v)", item.Mid, item.Token, res)
		s.reportCache.Save(func() { s.DelReport(context.TODO(), info.APPID, item.Mid, item.Token) })
	} else {
		log.Error("apns response(%+v)", res)
	}
	return
}

func (s *Service) pushHuawei(info *model.PushInfo, scheme string, items []*model.PushItem) (err error) {
	var (
		tokens []string
		res    *huawei.Response
		all    = int64(len(items))
	)
	for _, item := range items {
		tokens = append(tokens, item.Token)
	}
	for c := 0; c <= s.c.Push.RetryTimes; c++ {
		if res, err = s.dao.PushHuawei(context.TODO(), info, scheme, tokens); err == nil {
			break
		}
		dao.PromInfo("retry push huawei")
		s.setProgress(info.TaskID, _pgRetryTimes, 1)
		if err == huawei.ErrLimit {
			time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond)
		}
	}
	if err != nil {
		dao.PromError("push:推送华为")
		s.setProgress(info.TaskID, _pgTokenFailed, all)
		if err == huawei.ErrLimit {
			for _, t := range tokens {
				log.Error("push huawei task(%s) token(%s) error(%v)", info.TaskID, t, err)
			}
		}
		return
	}
	s.setProgress(info.TaskID, _pgTokenSuccess, all)
	s.cache.Save(func() { logPushed(info.TaskID, model.PlatformHuawei, items) })
	switch res.Code {
	case huawei.ResponseCodeSuccess:
		s.setProgress(info.TaskID, _pgTokenValid, all)
		log.Info("push huawei success task(%s) success(%d)", info.TaskID, all)
	case huawei.ResponseCodeSomeTokenInvalid:
		itr := &huawei.InvalidTokenResponse{}
		if err = json.Unmarshal([]byte(res.Msg), itr); err != nil {
			log.Error("json.Unmarshal() error(%v)", err)
			return
		}
		s.setProgress(info.TaskID, _pgTokenValid, int64(itr.Success))
		log.Warn("push huawei success task(%s) failed(%d) illegal(%v)", info.TaskID, itr.Failure, itr.IllegalTokens)
		if len(itr.IllegalTokens) == 0 {
			return
		}
		s.reportCache.Save(func() {
			m := make(map[string]int64, all)
			for _, i := range items {
				m[i.Token] = i.Mid
			}
			for _, t := range itr.IllegalTokens {
				s.DelReport(context.TODO(), info.APPID, m[t], t)
			}
		})
	case huawei.ResponseCodeAllTokenInvalid, huawei.ResponseCodeAllTokenInvalidNew:
		s.cache.Save(func() {
			for _, i := range items {
				s.DelReport(context.TODO(), info.APPID, i.Mid, i.Token)
			}
		})
		log.Error("push huawei failed task(%s) failed(%d) illegal(%v)", info.TaskID, all, tokens)
	default:
		log.Error("huawei push response task(%s) error(%v)", info.TaskID, res)
	}
	return
}

func (s *Service) pushOppoOne(info *model.PushInfo, item *model.PushItem) (err error) {
	params, _ := json.Marshal(map[string]string{
		"task_id": info.TaskID,
		"scheme":  model.Scheme(info.LinkType, info.LinkValue, model.PlatformAndroid, item.Build),
	})
	m := &oppo.Message{
		Title:        info.Title,
		Content:      info.Summary,
		ActionType:   oppo.ActionTypeInner,
		ActionParams: string(params),
		OfflineTTL:   int(int64(info.ExpireTime) - time.Now().Unix()),
		CallbackURL:  oppo.CallbackURL(info.APPID, info.TaskID),
	}
	var res *oppo.Response
	for c := 0; c <= s.c.Push.RetryTimes; c++ {
		if res, err = s.dao.PushOppoOne(context.TODO(), info, m, item.Token); err == nil {
			break
		}
		dao.PromInfo("retry push oppo")
		s.setProgress(info.TaskID, _pgRetryTimes, 1)
	}
	if err != nil || res == nil {
		s.setProgress(info.TaskID, _pgTokenFailed, 1)
		dao.PromError("push:推送oppo")
		log.Error("oppo push response task(%s) error(%v)", info.TaskID, res)
		return
	}
	s.setProgress(info.TaskID, _pgTokenSuccess, 1)
	if res.Code == oppo.ResponseCodeInvalidToken || res.Code == oppo.ResponseCodeUnsubscribeToken || res.Code == oppo.ResponseCodeRepeatToken {
		if item.Mid > 0 {
			s.reportCache.Save(func() { s.DelReport(context.TODO(), info.APPID, item.Mid, item.Token) })
		}
		return
	}
	s.setProgress(info.TaskID, _pgTokenValid, 1)
	return
}

// func (s *Service) pushOppo(info *model.PushInfo, items []*model.PushItem) (err error) {
// 	params, _ := json.Marshal(map[string]string{
// 		"task_id": info.TaskID,
// 		"scheme":  model.Scheme(info.LinkType, info.LinkValue, model.PlatformAndroid),
// 	})
// 	m := &oppo.Message{
// 		Title:        info.Title,
// 		Content:      info.Summary,
// 		ActionType:   oppo.ActionTypeInner,
// 		ActionParams: string(params),
// 		OfflineTTL:   int(int64(info.ExpireTime) - time.Now().Unix()),
// 		CallbackURL:  oppo.CallbackURL(info.APPID, info.TaskID),
// 	}
// 	res, err := s.dao.OppoMessage(context.TODO(), info, m)
// 	if err != nil || res == nil || res.Data.MsgID == "" {
// 		return
// 	}
// 	var (
// 		ts  []string
// 		all = int64(len(items))
// 		tm  = make(map[string]int64, all)
// 	)
// 	for _, i := range items {
// 		ts = append(ts, i.Token)
// 		tm[i.Token] = i.Mid
// 	}
// 	for c := 0; c <= s.c.Push.RetryTimes; c++ {
// 		if res, err = s.dao.PushOppo(context.TODO(), info, res.Data.MsgID, ts); err == nil {
// 			break
// 		}
// 		dao.PromInfo("retry push oppo")
// 		s.setProgress(info.TaskID, _pgRetryTimes, 1, model.PlatformOppo)
// 	}
// 	if err != nil || res == nil || res.Code != oppo.ResponseCodeSuccess {
// 		s.setProgress(info.TaskID, _pgTokenFailed, all, model.PlatformOppo)
// 		dao.PromError("push:推送oppo")
// 		log.Error("oppo push response task(%s) error(%v)", info.TaskID, res)
// 		if res.Code == oppo.ResponseCodeInvalidToken {
// 			s.delOppoReports(tm, info.APPID, ts)
// 		}
// 		return
// 	}
// 	s.setProgress(info.TaskID, _pgTokenSuccess, all, model.PlatformOppo)
// 	s.cache.Save(func() { logPushed(model.PlatformOppo, items) })
// 	var (
// 		invalid     = len(res.Data.TokenInvalid)
// 		unsubscribe = len(res.Data.TokenUnsubscribe)
// 		repeat      = len(res.Data.TokenRepeat)
// 		valid       = int(all) - invalid - unsubscribe - repeat
// 	)
// 	if invalid+unsubscribe+repeat > 0 {
// 		if invalid > 0 {
// 			s.delOppoReports(tm, info.APPID, res.Data.TokenInvalid)
// 		}
// 		if unsubscribe > 0 {
// 			s.delOppoReports(tm, info.APPID, res.Data.TokenUnsubscribe)
// 		}
// 		if repeat > 0 {
// 			s.delOppoReports(tm, info.APPID, res.Data.TokenRepeat)
// 		}
// 	}
// 	s.setProgress(info.TaskID, _pgTokenValid, int64(valid), model.PlatformOppo)
// 	s.cache.Save(func() { s.dao.AddTokensCache(context.TODO(), info.TaskID, ts) })
// 	log.Info("push oppo success task(%s) success(%d)", info.TaskID, all)
// 	return
// }

// func (s *Service) delOppoReports(m map[string]int64, appid int64, tokens []string) {
// 	if len(tokens) == 0 {
// 		return
// 	}
// 	s.reportCache.Save(func() {
// 		for _, t := range tokens {
// 			s.DelReport(context.TODO(), appid, m[t], t)
// 		}
// 	})
// }

func (s *Service) pushJpush(info *model.PushInfo, scheme string, items []*model.PushItem) (err error) {
	var (
		tokens []string
		all    = int64(len(items))
		valid  = all
		resp   *jpush.PushResponse
	)
	for _, item := range items {
		tokens = append(tokens, item.Token)
	}
	for c := 0; c <= s.c.Push.RetryTimes; c++ {
		if resp, err = s.dao.PushJpush(context.TODO(), info, scheme, tokens); err == nil && !resp.Retry {
			break
		}
		dao.PromInfo("retry push jpush")
		s.setProgress(info.TaskID, _pgRetryTimes, 1)
	}
	if err != nil {
		dao.PromError("push:推送极光")
		s.setProgress(info.TaskID, _pgTokenFailed, all)
		return
	}
	if resp.Error.Code != 0 {
		s.setProgress(info.TaskID, _pgTokenFailed, all)
		log.Error("jpush task(%s) tokens(%d) invalid code(%+v)", info.TaskID, len(tokens), resp)
		return
	}
	valid -= int64(len(resp.IllegalTokens))
	s.setProgress(info.TaskID, _pgTokenSuccess, all)
	s.cache.Save(func() { logPushed(info.TaskID, model.PlatformJpush, items) })
	s.setProgress(info.TaskID, _pgTokenValid, valid)
	log.Info("push jpush success task(%s) success(%d)", info.TaskID, all)
	return
}

func (s *Service) pushFCM(info *model.PushInfo, scheme string, items []*model.PushItem) (err error) {
	var (
		tokens []string
		all    = int64(len(items))
		valid  = all
		resp   *fcm.Response
	)
	for _, item := range items {
		tokens = append(tokens, item.Token)
	}
	for c := 0; c <= s.c.Push.RetryTimes; c++ {
		if resp, err = s.dao.PushFCM(context.Background(), info, scheme, tokens); err == nil {
			break
		}
		// 报错的情况下，如果 RetryAfter 有值才需要重试
		if resp != nil && resp.RetryAfter == "" {
			break
		}
		if d, e := resp.GetRetryAfterTime(); e == nil {
			time.Sleep(d)
		}
		dao.PromInfo("retry push fcm")
		s.setProgress(info.TaskID, _pgRetryTimes, 1)
	}
	if err != nil {
		dao.PromError("push:push fcm")
		s.setProgress(info.TaskID, _pgTokenFailed, all)
		return
	}
	s.setProgress(info.TaskID, _pgTokenSuccess, all)
	s.cache.Save(func() { logPushed(info.TaskID, model.PlatformFCM, items) })
	valid -= int64(resp.Fail)
	s.setProgress(info.TaskID, _pgTokenValid, valid)
	log.Info("push fcm success task(%s) success(%d)", info.TaskID, all)
	return
}

func (s *Service) tokensByMid(task *model.Task, mid int64) (res map[int][]*model.PushItem, err error) {
	var rs []*model.Report
	if rs, err = s.dao.ReportsCacheByMid(context.TODO(), mid); err != nil {
		if rs, err = s.dao.ReportsByMid(context.TODO(), mid); err != nil {
			return
		}
	}
	if len(rs) == 0 {
		return
	}
	p := map[int64][]*model.Report{mid: rs}
	res = s.distribute(task, p)
	return
}

func (s *Service) tokensByMids(c context.Context, task *model.Task, mids []int64) (res map[int][]*model.PushItem, missed []int64, err error) {
	rs, missed, err := s.dao.ReportsCacheByMids(c, mids)
	if err != nil {
		dao.PromInfo("report:查缓存失败回源")
		if rs, err = s.dao.ReportsByMids(c, mids); err != nil {
			return
		}
	}
	res = s.distribute(task, rs)
	return
}

func (s *Service) distribute(task *model.Task, rs map[int64][]*model.Report) (res map[int][]*model.PushItem) {
	var (
		validMid     int64
		validMidPlat = make(map[int]map[int64]bool)
		buildCount   = len(task.Build)
		// platformCount = len(task.Platform)
		brands = make(map[string]int64)
	)
	res = make(map[int][]*model.PushItem)
	for mid, rr := range rs {
		var valid bool
		for _, r := range rr {
			if r.APPID != task.APPID {
				// log.Info("task(%s) token(%s) mid(%d) app not match", task.ID, r.DeviceToken, mid)
				continue
			}
			if r.NotifySwitch == model.SwitchOff {
				// log.Info("task(%s) token(%s) mid(%d) switchoff", task.ID, r.DeviceToken, mid)
				continue
			}
			realTime := model.RealTime(r.TimeZone)
			if realTime.Unix() > int64(task.ExpireTime) {
				// log.Info("task(%s) token(%s) mid(%d) expire_time(%s) expired", task.ID, r.DeviceToken, mid, task.ExpireTime)
				continue
			}
			// if platformCount > 0 && !validatePlatform(r.PlatformID, task.Platform) {
			// 	// log.Info("task(%s) token(%s) mid(%d) platform forbid", task.ID, r.DeviceToken, mid)
			// 	continue
			// }
			if buildCount > 0 && !validateBuild(r.PlatformID, r.Build, task.Build) {
				// log.Info("task(%s) token(%s) mid(%d) build forbid", task.ID, r.DeviceToken, mid)
				continue
			}
			p := &model.PushItem{Platform: r.PlatformID, Token: r.DeviceToken, Mid: mid, Sound: task.Sound, Vibration: task.Vibration}
			res[r.PlatformID] = append(res[r.PlatformID], p)
			valid = true
			brands[r.DeviceBrand]++
			if _, ok := validMidPlat[r.PlatformID]; !ok {
				validMidPlat[r.PlatformID] = make(map[int64]bool)
			}
			validMidPlat[r.PlatformID][r.Mid] = true
		}
		if valid {
			validMid++
		}
	}
	s.setProgress(task.ID, _pgMidValid, validMid)
	for br, v := range brands {
		s.setBrandProgress(task.ID, model.DeviceBrand(br), v)
	}
	return
}

// func validatePlatform(platform int, set []int) bool {
// 	for _, v := range set {
// 		if v == platform {
// 			return true
// 		}
// 	}
// 	if platform == model.PlatformIPhone || platform == model.PlatformIPad {
// 		return false
// 	}
// 	if platform == model.PlatformAndroid {
// 		return true
// 	}
// 	return false
// }

func logPushed(task string, platform int, items []*model.PushItem) {
	for _, v := range items {
		log.Info("push done, task(%s) platform(%d) mid(%d) token(%s)", task, platform, v.Mid, v.Token)
	}
}

func validateBuild(platform, build int, builds map[int]*model.Build) bool {
	if len(builds) == 0 {
		return true
	}
	if builds[platform] == nil {
		return true
	}
	c := builds[platform].Condition
	b := builds[platform].Build
	switch c {
	case "gt":
		return build > b
	case "gte":
		return build >= b
	case "lt":
		return build < b
	case "lte":
		return build <= b
	case "eq":
		return build == b
	case "ne":
		return build != b
	}
	return false
}

func (s *Service) pushTokens(task *model.Task) (err error) {
	bs, err := ioutil.ReadFile(task.MidFile)
	if err != nil {
		log.Error("ioutil.ReadFile(%s) error(v)", task.MidFile, err)
		return
	}
	var (
		counter int
		group   = errgroup.Group{}
		info    = s.pushInfo(task)
		tokens  = strings.Split(string(bs), "\n")
	)
	for {
		l := len(tokens)
		if l == 0 {
			break
		}
		n := s.c.Push.PushPartSize
		if l < n {
			n = l
		}
		part := tokens[:n]
		tokens = tokens[n:]
		group.Go(func() error {
			res, missed, e := s.dao.TokensCache(context.Background(), part)
			if e != nil {
				log.Error("s.dao.TokensCache error(%v)", err)
			}
			if len(missed) > 0 {
				log.Info("task(%s) tokens cache missed(%d)", task.ID, len(missed))
			}
			var rs []*model.PushItem
			brands := make(map[int]int64)
			for _, t := range part {
				build := model.UnknownBuild
				if v, ok := res[t]; ok {
					build = v.Build
					brands[model.DeviceBrand(v.DeviceBrand)]++
				}
				rs = append(rs, &model.PushItem{Platform: task.PlatformID, Token: t, Sound: task.Sound, Vibration: task.Vibration, Build: build})
			}
			log.Info("push tokens task(%+v) tokens(%d)", task, len(rs))
			s.setProgress(task.ID, _pgTokenTotal, int64(len(rs)))
			for k, v := range brands {
				s.setBrandProgress(task.ID, k, v)
			}
			s.pushByPlatform(info, task.PlatformID, rs)
			return nil
		})
		time.Sleep(time.Duration(s.c.Push.PushPartInterval))
		counter++
		if counter > s.c.Push.PushPartChanSize {
			group.Wait()
			counter = 0
		}
	}
	if counter > 0 {
		group.Wait()
	}
	for {
		if s.chCounterVal(task.ID) == 0 {
			log.Info("Pushs done. task(%+v) result(%+v)", task, s.fetchProgress(task.ID))
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
}
