package feedback

import (
	"context"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/feedback"
	fdMdl "go-common/app/interface/main/creative/model/feedback"
	"go-common/app/interface/main/creative/service"
)

//Service struct.
type Service struct {
	c        *conf.Config
	feedback *feedback.Dao
}

//New get service.
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:        c,
		feedback: feedback.New(c),
	}
	return s
}

// Feedbacks get user feedback list.
func (s *Service) Feedbacks(c context.Context, mid, ps, pn, tagID int64, state, start, end, platform, ip string) (data []*fdMdl.Feedback, count int64, err error) {
	data, count, err = s.feedback.Feedbacks(c, mid, ps, pn, tagID, state, start, end, platform, ip)
	return
}

// AddFeedback add feedback
func (s *Service) AddFeedback(c context.Context, mid, tagID, sessionID int64, qq, content, aid, imgURL, browser, platform, ip string) (err error) {
	err = s.feedback.AddFeedback(c, mid, tagID, sessionID, qq, content, aid, browser, imgURL, platform, ip)
	return
}

// Tags get tags.
func (s *Service) Tags(c context.Context, mid int64, ip string) (tags []*fdMdl.Tag, err error) {
	var tls *fdMdl.TagList
	if tls, err = s.feedback.Tags(c, mid, ip); err != nil {
		return
	}
	if tls == nil || len(tls.Platforms) == 0 {
		return
	}
	for _, v := range tls.Platforms {
		if v.EN == "ugc" {
			tags = v.Tags
		}
	}
	return
}

// NewTags get new tags.
func (s *Service) NewTags(c context.Context, mid int64, ip string) (tags *fdMdl.TagList, err error) {
	tags, err = s.feedback.Tags(c, mid, ip)
	return
}

// Detail get feedback detail
func (s *Service) Detail(c context.Context, mid, sessionID int64, ip string) (data []*fdMdl.Reply, err error) {

	data, err = s.feedback.Detail(c, mid, sessionID, ip)
	return
}

// CloseSession get feedback detail
func (s *Service) CloseSession(c context.Context, sessionID int64, ip string) (err error) {
	err = s.feedback.CloseSession(c, sessionID, ip)
	return
}
