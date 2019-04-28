package pay

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/faq"
	"go-common/app/interface/main/creative/dao/pay"
	"go-common/app/interface/main/creative/dao/up"
	faqMdl "go-common/app/interface/main/creative/model/faq"
	"go-common/app/interface/main/creative/service"
	"go-common/library/log"
	"go-common/library/sync/pipeline/fanout"
	"time"
)

//Service struct
type Service struct {
	c               *conf.Config
	faq             *faq.Dao
	cache           *fanout.Fanout
	pay             *pay.Dao
	up              *up.Dao
	exemptUgcPayUps map[int64]int64
}

//New get service
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:     c,
		faq:   faq.New(c),
		pay:   pay.New(c),
		up:    rpcdaos.Up,
		cache: fanout.New("service_ugcpay", fanout.Worker(1), fanout.Buffer(10240)),
	}
	s.loadExemptUgcPayUps()
	go s.loadproc()
	return s
}

// loadExemptUgcPayUps
func (s *Service) loadExemptUgcPayUps() {
	ups, err := s.up.UpSpecial(context.TODO(), 17)
	if err != nil {
		return
	}
	s.exemptUgcPayUps = ups
}

// loadproc
func (s *Service) loadproc() {
	for {
		time.Sleep(5 * time.Minute)
		s.loadExemptUgcPayUps()
	}
}

// Pre fn
func (s *Service) Pre(c context.Context, mid int64) (PrePay map[string]interface{}, err error) {
	PrePay = map[string]interface{}{
		"protocol_id":        s.c.UgcPay.ProtocolID,
		"protocol_need_read": false,
		"white":              false,
	}
	log.Warn("s.exemptUgcPayUps (%+v)", s.exemptUgcPayUps)
	if _, ok := s.exemptUgcPayUps[mid]; ok {
		PrePay["protocol_need_read"] = true
		PrePay["white"] = true
		if accept, _ := s.pay.UserAcceptProtocol(c, s.c.UgcPay.ProtocolID, mid); accept {
			PrePay["protocol_need_read"] = false
		}
	}
	return
}

// Protocol fn
func (s *Service) Protocol(c context.Context, protocolID string) (pd *faqMdl.Detail, err error) {
	resD := make([]*faqMdl.Detail, 0)
	var (
		total               int
		protocolLimitLength = int(1000)
	)
	if resD, total, err = s.faq.DetailCache(c, faqMdl.FaqUgcProtocolQuesTypeID, 1, 1, protocolLimitLength); err != nil || len(resD) == 0 {
		if resD, total, err = s.faq.Detail(context.Background(), faqMdl.FaqUgcProtocolQuesTypeID, 1, 1, protocolLimitLength); err != nil {
			log.Error("s.faq.Detail(%s,%d,%d,%d) error(%v)", faqMdl.FaqUgcProtocolQuesTypeID, 1, 1, protocolLimitLength, err)
			return
		}
	}
	if len(resD) > 0 {
		s.cache.Do(c, func(c context.Context) {
			s.faq.SetDetailCache(context.Background(), faqMdl.FaqUgcProtocolQuesTypeID, 1, 1, protocolLimitLength, total, resD)
		})
		for _, v := range resD {
			if v.AnswerID == protocolID {
				pd = v
				return
			}
		}
	}
	return
}

// White fn
func (s *Service) White(c context.Context, mid int64) (white bool, err error) {
	if _, ok := s.exemptUgcPayUps[mid]; ok {
		white = true
	}
	return
}
