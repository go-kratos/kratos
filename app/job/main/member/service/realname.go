package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"go-common/app/job/main/member/conf"
	"go-common/app/job/main/member/model"
	memmodel "go-common/app/service/main/member/model"
	"go-common/library/log"
	"go-common/library/net/ip"

	"github.com/pkg/errors"
)

// constrs for gender
const (
	_genderMale   = "male"
	_genderFemale = "female"
)

// realname alipay polling

func (s *Service) realnamealipaycheckproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("%+v", errors.WithStack(fmt.Errorf("service.realnamealipaycheckproc panic(%v)", x)))
			go s.realnamealipaycheckproc()
			log.Info("service.realnamealipaycheckproc recover")
		}
	}()
	for {
		var (
			to          = time.Now()
			from        = to.Add(-2 * time.Duration(conf.Conf.Biz.RealnameAlipayCheckTick))
			expiredTime = from
			startTime   = expiredTime.AddDate(0, -1, 0)
		)
		log.Info("realname alipay check start from : %s , to : %s", from, to)
		s.realnameAlipayCheckHandler(context.Background(), from, to)
		// to = from
		// from = to.Add(-2 * time.Duration(conf.Conf.Biz.RealnameAlipayCheckTick))
		log.Info("realname alipay handle expired end startTime : %s , expiredTime : %s", startTime, expiredTime)
		s.realnameAlipayExpiredHandler(context.Background(), startTime, expiredTime)
		time.Sleep(time.Duration(conf.Conf.Biz.RealnameAlipayCheckTick))
	}
}

// realnameAlipayCheckHandler 轮询时间段 [from,to] 中，未完成阿里实名的实名申请
func (s *Service) realnameAlipayCheckHandler(c context.Context, from, to time.Time) {
	if conf.Conf.Biz.RealnameAlipayCheckLimit <= 0 {
		log.Error("conf.Conf.Property.realnameAlipayCheckHandler [%d] <= 0", conf.Conf.Biz.RealnameAlipayCheckLimit)
		return
	}
	var (
		applys  = make([]*model.RealnameAlipayApply, conf.Conf.Biz.RealnameAlipayCheckLimit)
		startID int64
		err     error
	)
	for len(applys) >= conf.Conf.Biz.RealnameAlipayCheckLimit {
		if startID, applys, err = s.dao.RealnameAlipayApplyList(c, startID, model.RealnameApplyStatusPending, from, to, conf.Conf.Biz.RealnameAlipayCheckLimit); err != nil {
			log.Error("%+v", err)
			return
		}
		for _, apply := range applys {
			log.Info("Start check realname alipay apply mid (%d) bizno (%s)", apply.MID, apply.Bizno)
			if err = s.realnameAlipayConfirm(c, apply); err != nil {
				log.Error("%+v", err)
				continue
			}
		}
	}

	for len(applys) >= conf.Conf.Biz.RealnameAlipayCheckLimit {
		if startID, applys, err = s.dao.RealnameAlipayApplyList(c, startID, model.RealnameApplyStatusBack, from, to, conf.Conf.Biz.RealnameAlipayCheckLimit); err != nil {
			log.Error("%+v", err)
			return
		}
		for _, apply := range applys {
			log.Info("Start check realname alipay apply mid (%d) bizno (%s)", apply.MID, apply.Bizno)
			if err = s.realnameAlipayConfirm(c, apply); err != nil {
				log.Error("%+v", err)
				continue
			}
		}
	}
}

func (s *Service) realnameAlipayConfirm(c context.Context, apply *model.RealnameAlipayApply) (err error) {
	if apply.Bizno == "" {
		return
	}
	var (
		pass   bool
		reason string
	)
	if pass, reason, err = s.alipayQuery(c, apply.Bizno); err != nil {
		return
	}
	// rpc call
	var (
		rpcConfirmArg = &memmodel.ArgRealnameAlipayConfirm{
			MID:    apply.MID,
			Pass:   pass,
			Reason: reason,
		}
	)
	if err = s.memrpc.RealnameAlipayConfirm(c, rpcConfirmArg); err != nil {
		return
	}
	log.Info("Succeed to confirm realname alipay with arg: %+v", rpcConfirmArg)
	if pass {
		expArg := &model.AddExp{
			Mid:   apply.MID,
			IP:    ip.InternalIP(),
			Ts:    time.Now().Unix(),
			Event: "identify",
		}
		if expErr := s.addExp(context.TODO(), expArg); expErr != nil {
			log.Error("realname exp error(%+v) ", expErr)
			return
		}
		log.Info("realname exp success(%+v)", expArg)
	}

	return
}

func (s *Service) alipayQuery(c context.Context, bizno string) (pass bool, reason string, err error) {
	var (
		param url.Values
		biz   struct {
			Bizno string `json:"biz_no"`
		}
	)
	biz.Bizno = bizno
	if param, err = s.alipayParam("zhima.customer.certification.query", biz, ""); err != nil {
		return
	}
	if pass, reason, err = s.dao.AlipayQuery(c, param); err != nil {
		return
	}
	return
}

// alipayParam 构造阿里请求param，biz为 biz_content struct
func (s *Service) alipayParam(method string, biz interface{}, returnURL string) (p url.Values, err error) {
	var (
		sign     string
		bizBytes []byte
	)
	if bizBytes, err = json.Marshal(biz); err != nil {
		err = errors.WithStack(err)
		return
	}
	p = url.Values{}
	p.Set("app_id", conf.Conf.Biz.RealnameAlipayAppID)
	p.Set("method", method)
	p.Set("charset", "utf-8")
	p.Set("sign_type", "RSA2")
	p.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	p.Set("version", "1.0")
	p.Set("biz_content", string(bizBytes))
	if returnURL != "" {
		p.Set("return_url", returnURL)
	}
	if sign, err = s.alipayCryptor.SignParam(p); err != nil {
		return
	}
	p.Set("sign", sign)
	return
}

// rejectExpiredRealnameAlipay 自动驳回超过两天还没有通过芝麻认证的实名认证
func (s *Service) realnameAlipayExpiredHandler(c context.Context, startTime, expiredTime time.Time) {
	if conf.Conf.Biz.RealnameAlipayCheckLimit <= 0 {
		log.Error("conf.Conf.Property.realnameAlipayCheckHandler [%d] <= 0", conf.Conf.Biz.RealnameAlipayCheckLimit)
		return
	}
	var (
		applys  []*model.RealnameAlipayApply
		startID int64
		err     error
	)
	// 每次查询（一个月里）100条过期的未处理的位处理的芝麻认证数据，进行驳回
	for {
		log.Info("realname handle startID (%d)", startID)
		startID, applys, err = s.dao.RealnameAlipayApplyList(c, startID, model.RealnameApplyStatusPending, startTime, expiredTime, conf.Conf.Biz.RealnameAlipayCheckLimit)
		if err != nil {
			log.Error("realnameAlipayExpiredHandler search err(%+v)", err)
			return
		}
		// 没有查询到预期的过期数据，则停止循环，等待下一次检查
		if len(applys) == 0 {
			log.Error("realnameAlipayExpiredHandler search no row in result")
			return
		}
		// 循环驳回验证超时的芝麻认证
		for _, apply := range applys {
			log.Info("Start expire realname alipay apply mid (%d) bizno (%s)", apply.MID, apply.Bizno)
			var (
				rpcConfirmArg = &memmodel.ArgRealnameAlipayConfirm{
					MID:    apply.MID,
					Pass:   false,
					Reason: "超时自动驳回",
				}
			)
			if err = s.memrpc.RealnameAlipayConfirm(c, rpcConfirmArg); err != nil {
				log.Error("realnameAlipayExpiredHandler reject err(%+v)", err)
				continue
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// ParseIdentity to birthday and gender
func ParseIdentity(id string) (birthday time.Time, gender string, err error) {
	var (
		ystr, mstr, dstr, gstr string
		y, m, d, g             int
	)
	switch len(id) {
	case 15:
		ystr, mstr, dstr = "19"+id[6:8], id[8:10], id[10:12]
		gstr = id[14:15]
	case 18:
		ystr, mstr, dstr = id[6:10], id[10:12], id[12:14]
		gstr = id[16:17]
	default:
		err = errors.Errorf("identity id invalid : %s", id)
		return
	}
	if y, err = strconv.Atoi(ystr); err != nil {
		err = errors.WithStack(err)
		return
	}
	if m, err = strconv.Atoi(mstr); err != nil {
		err = errors.WithStack(err)
		return
	}
	if d, err = strconv.Atoi(dstr); err != nil {
		err = errors.WithStack(err)
		return
	}
	if g, err = strconv.Atoi(gstr); err != nil {
		err = errors.WithStack(err)
		return
	}
	birthday = time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
	if g%2 == 1 {
		gender = _genderMale
	} else {
		gender = _genderFemale
	}
	return
}
