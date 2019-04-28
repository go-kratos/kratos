package realname

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"time"

	"go-common/app/interface/main/account/conf"
	"go-common/app/interface/main/account/model"
	memmodel "go-common/app/service/main/member/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

var (
	_idCard18Regexp = regexp.MustCompile(`^\d{17}[\d|x]$`)
	_idCard15Regexp = regexp.MustCompile(`^\d{15}$`)
)

// AlipayApply 提交芝麻认证申请
func (s *Service) AlipayApply(c context.Context, mid int64, param *model.ParamRealnameAlipayApply) (res *model.RealnameAlipayApply, err error) {
	if !s.alipayAntispamCheck(c, mid) {
		err = ecode.RealnameAlipayAntispam
		return
	}
	if !s.checkID(param.CardNum) {
		err = ecode.RealnameCardNumErr
		return
	}
	var bizno string
	if bizno, err = s.alipayInit(c, param.Realname, param.CardNum); err != nil {
		return
	}
	res = &model.RealnameAlipayApply{}
	if res.URL, err = s.alipayCertifyURL(c, bizno); err != nil {
		return
	}
	var (
		arg = &memmodel.ArgRealnameAlipayApply{
			MID:         mid,
			CaptureCode: param.Capture,
			Realname:    param.Realname,
			CardCode:    param.CardNum,
			IMGToken:    param.ImgToken,
			Bizno:       bizno,
		}
	)
	if err = s.memRPC.RealnameAlipayApply(c, arg); err != nil {
		res = nil
		return
	}
	s.addmiss(func() {
		if missErr := s.alipayAntispamIncrease(context.Background(), mid); missErr != nil {
			log.Error("%+v", err)
		}
	})
	return
}

func (s *Service) checkID(card string) bool {
	if !_idCard15Regexp.MatchString(card) && !_idCard18Regexp.MatchString(card) {
		return false
	}
	return true
}

func (s *Service) alipayInit(c context.Context, realname, cardNum string) (bizno string, err error) {
	var (
		biz struct {
			TransID   string `json:"transaction_id"`
			ProdCode  string `json:"product_code"`
			BizCode   string `json:"biz_code"`
			IdenParam struct {
				IdentityType string `json:"identity_type"`
				CertType     string `json:"cert_type"`
				CertName     string `json:"cert_name"`
				CertNo       string `json:"cert_no"`
			} `json:"identity_param"`
		}
		param url.Values
	)
	biz.TransID = s.alipayTransactionID() // 商户请求的唯一标志，32位长度的字母数字下划线组合。该标识作为对账的关键信息，商户要保证其唯一性.
	biz.ProdCode = "w1010100000000002978" // 芝麻认证产品码
	biz.BizCode = "FACE"                  // 认证场景码，支持的场景码有： FACE：多因子活体人脸认证， SMART_FACE：多因子快捷活体人脸认证， FACE_SDK：SDK活体人脸认证 签约的协议决定了可以使用的场景
	biz.IdenParam.IdentityType = "CERT_INFO"
	biz.IdenParam.CertType = "IDENTITY_CARD"
	biz.IdenParam.CertName = realname
	biz.IdenParam.CertNo = cardNum
	if param, err = s.alipayParam("zhima.customer.certification.initialize", biz, ""); err != nil {
		return
	}
	if bizno, err = s.realnameDao.AlipayInit(c, param); err != nil {
		log.Error("%+v", err)
		err = ecode.RealnameAlipayErr
		return
	}
	if bizno == "" {
		err = ecode.RealnameAlipayErr
		return
	}
	return
}

func (s *Service) alipayCertifyURL(c context.Context, bizno string) (u string, err error) {
	var (
		param url.Values
		biz   struct {
			Bizno string `json:"biz_no"`
		}
	)
	biz.Bizno = bizno
	if param, err = s.alipayParam("zhima.customer.certification.certify", biz, "bilibili://auth.zhima"); err != nil {
		return
	}
	u = conf.Conf.Realname.Alipay.Gateway + "?" + param.Encode()
	return
}

// AlipayConfirm 查询芝麻认证状态
func (s *Service) AlipayConfirm(c context.Context, mid int64) (res *model.RealnameAlipayConfirm, err error) {
	var (
		pass   bool
		reason string
		rpcarg = &memmodel.ArgMemberMid{
			Mid: mid,
		}
		rpcres *memmodel.RealnameAlipayInfo
	)
	if rpcres, err = s.memRPC.RealnameAlipayBizno(c, rpcarg); err != nil {
		return
	}
	if pass, reason, err = s.alipayQuery(c, rpcres.Bizno); err != nil {
		log.Error("%+v", err)
		err = ecode.RealnameAlipayErr
		return
	}
	res = &model.RealnameAlipayConfirm{
		Reason: reason,
	}
	if pass {
		res.Passed = model.RealnameTrue
	} else {
		res.Passed = model.RealnameFalse
	}
	// rpc call
	var (
		rpcConfirmArg = &memmodel.ArgRealnameAlipayConfirm{
			MID:    mid,
			Pass:   pass,
			Reason: reason,
		}
	)
	if err = s.memRPC.RealnameAlipayConfirm(c, rpcConfirmArg); err != nil {
		return
	}
	s.addmiss(func() {
		if missErr := s.setAlipayAntispamPassFlag(context.Background(), mid, false); missErr != nil {
			log.Error("%+v", err)
		}
	})
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
	if pass, reason, err = s.realnameDao.AlipayQuery(c, param); err != nil {
		log.Error("%+v", err)
		err = ecode.RealnameAlipayErr
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
	p.Set("app_id", conf.Conf.Realname.Alipay.AppID)
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

func (s *Service) alipayTransactionID() string {
	return fmt.Sprintf("BILI%d%d", time.Now().UnixNano(), rand.Intn(100000))
}
