package faq

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/faq"
	faqMdl "go-common/app/interface/main/creative/model/faq"
	"go-common/app/interface/main/creative/service"
	"go-common/library/log"
	"go-common/library/sync/pipeline/fanout"
)

//Service struct
type Service struct {
	c     *conf.Config
	faq   *faq.Dao
	cache *fanout.Fanout
}

//New get service
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:     c,
		faq:   faq.New(c),
		cache: fanout.New("faq-service"),
	}
	return s
}

// Pre fn
func (s *Service) Pre(c context.Context) (res map[string]*faqMdl.Faq) {
	res = make(map[string]*faqMdl.Faq)
	var (
		total int
		err   error
		resD  []*faqMdl.Detail
	)
	res["app_edit_faq"] = &faqMdl.Faq{
		State:            false,
		QuestionTypeID:   faqMdl.PhoneFaqQuesTypeID,
		QuestionTypeName: faqMdl.PhoneFaqName,
		URL:              s.c.H5Page.FAQVideoEditor,
	}
	if resD, total, err = s.faq.DetailCache(c, faqMdl.PhoneFaqQuesTypeID, 1, 1, 2); err != nil || len(resD) == 0 {
		if _, total, err = s.faq.Detail(context.Background(), faqMdl.PhoneFaqQuesTypeID, 1, 1, 2); err != nil {
			log.Error("s.faq.Detail(%s,%d,%d,%d) error(%v)", faqMdl.PhoneFaqQuesTypeID, 1, 1, 2, err)
		}
	}
	if total > 0 {
		res["app_edit_faq"].State = true
	}
	return
}

// Detail fn
func (s *Service) Detail(c context.Context, qTypeID string, keyFlag, pn, ps int) (resD []*faqMdl.Detail, total int, err error) {
	if resD, total, err = s.faq.DetailCache(c, qTypeID, keyFlag, pn, ps); err != nil || len(resD) == 0 {
		if resD, total, err = s.faq.Detail(context.Background(), qTypeID, keyFlag, pn, ps); err != nil {
			log.Error("s.faq.Detail(%s,%d,%d,%d) error(%v)", qTypeID, keyFlag, pn, ps, err)
		}
		if len(resD) > 0 {
			s.cache.Do(c, func(c context.Context) {
				s.faq.SetDetailCache(context.Background(), qTypeID, keyFlag, pn, ps, total, resD)
			})
		}
	}
	return
}
