package realname

import (
	"context"
	"strings"

	"go-common/app/interface/main/account/conf"
	"go-common/app/interface/main/account/dao/realname"
	"go-common/app/interface/main/account/model"
	"go-common/app/interface/main/account/service/realname/crypto"
	memrpc "go-common/app/service/main/member/api/gorpc"
	memmodel "go-common/app/service/main/member/model"
	"go-common/library/log"
)

// Service is
type Service struct {
	c             *conf.Config
	memRPC        *memrpc.Service
	realnameDao   *realname.Dao
	mainCryptor   *crypto.Main
	alipayCryptor *crypto.Alipay
	missch        chan func()
}

// New create service instance and return.
func New(c *conf.Config, rsapub, rsapriv, alipub, alibilipriv string) (s *Service) {
	s = &Service{
		c:             c,
		memRPC:        memrpc.New(c.RPCClient2.Member),
		realnameDao:   realname.New(c),
		mainCryptor:   crypto.NewMain(rsapub, rsapriv),
		alipayCryptor: crypto.NewAlipay(alipub, alibilipriv),
		missch:        make(chan func(), 1024),
	}
	go s.missproc()
	return
}

// Status get status of realname
func (s *Service) Status(c context.Context, mid int64) (status int8, err error) {
	var (
		arg = &memmodel.ArgMemberMid{
			Mid: mid,
		}
		res *memmodel.RealnameStatus
	)
	if res, err = s.memRPC.RealnameStatus(c, arg); err != nil {
		return
	}
	status = int8(*res)
	return
}

// ApplyStatus return realname apply status
func (s *Service) ApplyStatus(c context.Context, mid int64) (status int8, remark string, realname string, card string, err error) {
	var (
		arg = &memmodel.ArgMemberMid{
			Mid: mid,
		}
		res *memmodel.RealnameApplyStatusInfo
	)
	if res, err = s.memRPC.RealnameApplyStatus(c, arg); err != nil {
		return
	}
	status = int8(res.Status)
	remark = res.Remark
	realname, card = maskRealnameInfo(res.Realname, res.Card)
	return
}

func maskRealnameInfo(realname, card string) (r, c string) {
	var (
		rStrs = strings.Split(realname, "")
		cStrs = strings.Split(card, "")
	)
	if len(rStrs) > 0 {
		r = "*" + strings.Join(rStrs[1:], "")
	}
	if len(cStrs) > 0 {
		if len(cStrs) == 1 {
			c = "*"
		} else if len(cStrs) > 5 {
			c = cStrs[0] + strings.Repeat("*", len(cStrs)-3) + strings.Join(cStrs[len(cStrs)-2:], "")
		} else {
			c = cStrs[0] + strings.Repeat("*", len(cStrs)-1)
		}
	}
	return
}

// CountryList .
func (s *Service) CountryList(c context.Context) (list []*model.RealnameCountry, err error) {
	list = countryList
	return
}

// CardTypes .
func (s *Service) CardTypes(c context.Context, platform string, mobiapp string, device string, build int) (list []*model.RealnameCardType, err error) {
	if (platform == "android" && build < 512000) || (platform == "ios" && build <= 5990) {
		list = cardTypeOldList
		return
	}
	// IOS粉暂返回 5.32 版本数据，待 5.36 IOS 重新接入后，根据 build 号，返回 cardTypeList
	if platform == "ios" && mobiapp == "iphone" && device != "pad" {
		list = cardTypeOldIOSList
		return
	}
	list = cardTypeList
	return
}

// CardTypesV2 .
func (s *Service) CardTypesV2(c context.Context) (list []*model.RealnameCardType, err error) {
	list = cardTypeList
	return
}

// TelCapture .
func (s *Service) TelCapture(c context.Context, mid int64) (err error) {
	var (
		arg = &memmodel.ArgMemberMid{
			Mid: mid,
		}
	)
	if err = s.memRPC.RealnameTelCapture(c, arg); err != nil {
		return
	}
	return
}

// TelInfo .
func (s *Service) TelInfo(c context.Context, mid int64) (tel string, err error) {
	if tel, err = s.realnameDao.TelInfo(c, mid); err != nil {
		return
	}
	if len(tel) == 0 {
		return
	}
	if len(tel) < 4 {
		tel = tel[:1] + "****"
		return
	}
	tel = tel[:3] + "****" + tel[len(tel)-4:]
	return
}

// Apply .
func (s *Service) Apply(c context.Context, mid int64, realname string, cardType int, cardNum string, countryID int, captureCode int, handIMGToken, frontIMGToken, backIMGToken string) (err error) {
	var (
		arg = &memmodel.ArgRealnameApply{
			MID:           mid,
			CaptureCode:   captureCode,
			Realname:      realname,
			CardType:      int8(cardType),
			CardCode:      cardNum,
			Country:       int16(countryID),
			HandIMGToken:  handIMGToken,
			FrontIMGToken: frontIMGToken,
			BackIMGToken:  backIMGToken,
		}
	)
	if err = s.memRPC.RealnameApply(c, arg); err != nil {
		return
	}
	return
}

// Channel .
func (s *Service) Channel(c context.Context) (channels []*model.RealnameChannel, err error) {
	for _, c := range conf.Conf.Realname.Channel {
		var (
			channel = &model.RealnameChannel{
				Name: c.Name,
				Flag: model.RealnameFalse,
			}
		)
		if c.Flag {
			channel.Flag = model.RealnameTrue
		}
		channels = append(channels, channel)
	}
	return
}

func (s *Service) addmiss(f func()) {
	select {
	case s.missch <- f:
	default:
		log.Error("eventproc chan full")
	}
}

// missproc is a routine for executing closure.
func (s *Service) missproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("missproc panic %+v", x)
		}
		go s.missproc()
	}()
	for {
		f := <-s.missch
		f()
	}
}
