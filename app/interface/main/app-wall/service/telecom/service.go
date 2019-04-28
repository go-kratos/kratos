package telecom

import (
	"fmt"

	"go-common/app/interface/main/app-wall/conf"
	seqDao "go-common/app/interface/main/app-wall/dao/seq"
	telecomDao "go-common/app/interface/main/app-wall/dao/telecom"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

const (
	_initIPlimitKey = "iplimit_%v_%v"
	_telecomKey     = "telecom"
)

type Service struct {
	c                  *conf.Config
	dao                *telecomDao.Dao
	seqdao             *seqDao.Dao
	flowPercentage     int
	smsTemplate        string
	smsMsgTemplate     string
	smsFlowTemplate    string
	smsOrderTemplateOK string
	operationIPlimit   map[string]struct{}
	telecomArea        map[string]struct{}
	// prom
	pHit  *prom.Prom
	pMiss *prom.Prom
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                  c,
		dao:                telecomDao.New(c),
		seqdao:             seqDao.New(c),
		flowPercentage:     c.Telecom.FlowPercentage,
		smsTemplate:        c.Telecom.SMSTemplate,
		smsMsgTemplate:     c.Telecom.SMSMsgTemplate,
		smsFlowTemplate:    c.Telecom.SMSFlowTemplate,
		smsOrderTemplateOK: c.Telecom.SMSOrderTemplateOK,
		operationIPlimit:   map[string]struct{}{},
		telecomArea:        map[string]struct{}{},
		// prom
		pHit:  prom.CacheHit,
		pMiss: prom.CacheMiss,
	}
	go s.loadIPlimit(c)
	go s.loadTelecomArea(c)
	return
}

func (s *Service) loadIPlimit(c *conf.Config) {
	hosts := make(map[string]struct{}, len(c.IPLimit.Addrs))
	for k, v := range c.IPLimit.Addrs {
		for _, ipStr := range v {
			key := fmt.Sprintf(_initIPlimitKey, k, ipStr)
			if _, ok := hosts[key]; !ok {
				hosts[key] = struct{}{}
			}
		}
	}
	s.operationIPlimit = hosts
	log.Info("loadTelecomIPCache success")
}

func (s *Service) loadTelecomArea(c *conf.Config) {
	areas := make(map[string]struct{}, len(c.Telecom.Area))
	for _, v := range c.Telecom.Area {
		for _, area := range v {
			if _, ok := areas[area]; !ok {
				areas[area] = struct{}{}
			}
		}
	}
	s.telecomArea = areas
	log.Info("loadTelecomArea success")
}
