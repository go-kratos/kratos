package client

import (
	"context"
	"testing"
	"time"

	model "go-common/app/service/main/assist/model/assist"

	"github.com/davecgh/go-spew/spew"
)

const (
	mid       = 27515256
	assistMid = 27515255
	realIP    = "127.0.0.1"
	logID     = 692
	subjectID = 111
	objectID  = "222"
	detail    = "testing"
	pn        = 1
	ps        = 20
)

func TestAssistRpcService(t *testing.T) {
	s := New(nil)
	time.Sleep(1 * time.Second)

	// test assist relation
	delAssist(t, s)
	addAssist(t, s)
	assist(t, s)
	assists(t, s)
	assistIDs(t, s)
	assistUps(t, s)
	assistExit(t, s)
	// test assistlog
	assistLogAdd(t, s)
	assistLogInfo(t, s)
	assistLogCancel(t, s)
	assistLogs(t, s)
}

func delAssist(t *testing.T, s *Service) {
	arg := &model.ArgAssist{
		Mid:       mid,
		AssistMid: assistMid,
		RealIP:    realIP,
	}
	if err := s.DelAssist(context.TODO(), arg); err != nil {
		t.Logf("call error(%v)", err)
	}
}

func addAssist(t *testing.T, s *Service) {
	arg := &model.ArgAssist{
		Mid:       mid,
		AssistMid: assistMid,
		RealIP:    realIP,
	}
	if err := s.AddAssist(context.TODO(), arg); err != nil {
		t.Logf("call error(%v)", err)
	}
}

func assistLogAdd(t *testing.T, s *Service) {
	arg := &model.ArgAssistLogAdd{
		Mid:       mid,
		AssistMid: assistMid,
		Type:      model.TypeComment,
		Action:    model.ActDelete,
		SubjectID: subjectID,
		ObjectID:  objectID,
		Detail:    detail,
		RealIP:    realIP,
	}
	if err := s.AssistLogAdd(context.TODO(), arg); err != nil {
		t.Logf("call error(%v)", err)
	}
}

func assists(t *testing.T, s *Service) {
	arg := &model.ArgAssists{
		Mid:    mid,
		RealIP: realIP,
	}
	if res, err := s.Assists(context.TODO(), arg); err != nil && res != nil {
		t.Logf("call error(%v)", err)
	}
}

func assistIDs(t *testing.T, s *Service) {
	arg := &model.ArgAssists{
		Mid:    mid,
		RealIP: realIP,
	}
	if res, err := s.AssistIDs(context.TODO(), arg); err != nil && res != nil {
		t.Logf("call error(%v)", err)
	}
}

func assistUps(t *testing.T, s *Service) {
	arg := &model.ArgAssistUps{
		AssistMid: assistMid,
		Ps:        20,
		Pn:        1,
		RealIP:    realIP,
	}
	if res, err := s.AssistUps(context.TODO(), arg); err != nil && res != nil {
		spew.Dump(res)
		t.Logf("call error(%v)", err)
	}
}

func assistExit(t *testing.T, s *Service) {
	arg := &model.ArgAssist{
		AssistMid: assistMid,
		Mid:       mid,
		RealIP:    realIP,
	}
	if err := s.AssistExit(context.TODO(), arg); err != nil {
		t.Logf("call error(%v)", err)
	}
}

func assistLogInfo(t *testing.T, s *Service) {
	arg := &model.ArgAssistLog{
		Mid:       mid,
		AssistMid: assistMid,
		LogID:     logID,
		RealIP:    realIP,
	}
	if res, err := s.AssistLogInfo(context.TODO(), arg); err != nil && res != nil {
		t.Logf("call error(%v)", err)
	}
}

func assist(t *testing.T, s *Service) {
	arg := &model.ArgAssist{
		Mid:       mid,
		AssistMid: assistMid,
		RealIP:    realIP,
	}
	if res, err := s.Assist(context.TODO(), arg); err != nil && res != nil {
		spew.Dump(res)
		t.Logf("call error(%v)", err)
	}
}

func assistLogs(t *testing.T, s *Service) {
	arg := &model.ArgAssistLogs{
		Mid:       mid,
		AssistMid: assistMid,
		Stime:     time.Unix(time.Now().Unix(), 0),
		Etime:     time.Unix(time.Now().Add(48*time.Hour).Unix(), 0),
		Pn:        ps,
		Ps:        pn,
		RealIP:    realIP,
	}
	if res, err := s.AssistLogs(context.TODO(), arg); err != nil && res != nil {
		t.Logf("call error(%v)", err)
	}
}

func assistLogCancel(t *testing.T, s *Service) {
	arg := &model.ArgAssistLog{
		Mid:       mid,
		AssistMid: assistMid,
		LogID:     logID,
		RealIP:    realIP,
	}
	if err := s.AssistLogCancel(context.TODO(), arg); err != nil {
		t.Logf("call error(%v)", err)
	}
}
