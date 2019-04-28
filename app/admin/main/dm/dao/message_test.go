package dao

import (
	"testing"

	"go-common/app/admin/main/dm/model"

	. "github.com/smartystreets/goconvey/convey"
)

//func TestDao_SendMsgToReporter1(t *testing.T) {
//	var (
//		d      = New(conf.Conf)
//		c      = context.TODO()
//		err    error
//		rptMsg = &model.ReportMsg{
//			Aid:         1,
//			Title:       "test title1",
//			Did:         1,
//			Msg:         "dm msg1",
//			Uids:        "150781",
//			State:       model.StatFirstDelete,
//			RptReason:   1,
//			BlockReason: 1,
//			Block:       -1,
//		}
//	)
//	err = d.SendMsgToReporter(c, rptMsg)
//	if err != nil {
//		t.Errorf("dao.SendMsgToReporter(rptMsg:%v) err(%v)", rptMsg, err)
//		t.Fail()
//	}
//}
//
//func TestDao_SendMsgToReporter2(t *testing.T) {
//	var (
//		d      = New(conf.Conf)
//		c      = context.TODO()
//		err    error
//		rptMsg = &model.ReportMsg{
//			Aid:         1,
//			Title:       "test title2",
//			Did:         1,
//			Msg:         "dm msg2",
//			Uids:        "150781",
//			State:       model.StatSecondDelete,
//			RptReason:   3,
//			BlockReason: 5,
//			Block:       -1,
//		}
//	)
//	err = d.SendMsgToReporter(c, rptMsg)
//	if err != nil {
//		t.Errorf("dao.SendMsgToReporter(rptMsg:%v) err(%v)", rptMsg, err)
//		t.Fail()
//	}
//}
//
//func TestDao_SendMsgToReporter3(t *testing.T) {
//	var (
//		d      = New(conf.Conf)
//		c      = context.TODO()
//		err    error
//		rptMsg = &model.ReportMsg{
//			Aid:         1,
//			Title:       "test title3",
//			Did:         1,
//			Msg:         "dm msg3",
//			Uids:        "150781",
//			State:       model.StatFirstIgnore,
//			RptReason:   1,
//			BlockReason: 1,
//			Block:       -1,
//		}
//	)
//	err = d.SendMsgToReporter(c, rptMsg)
//	if err != nil {
//		t.Errorf("dao.SendMsgToReporter(rptMsg:%v) err(%v)", rptMsg, err)
//		t.Fail()
//	}
//}
//
//func TestDao_SendMsgToReporter4(t *testing.T) {
//	var (
//		d      = New(conf.Conf)
//		c      = context.TODO()
//		err    error
//		rptMsg = &model.ReportMsg{
//			Aid:         1,
//			Title:       "test title4",
//			Did:         1,
//			Msg:         "dm msg4",
//			Uids:        "150781",
//			State:       model.StatSecondIgnore,
//			RptReason:   1,
//			BlockReason: 1,
//			Block:       -1,
//		}
//	)
//	err = d.SendMsgToReporter(c, rptMsg)
//	if err != nil {
//		t.Errorf("dao.SendMsgToReporter(rptMsg:%v) err(%v)", rptMsg, err)
//		t.Fail()
//	}
//}
//
//func TestDao_SendMsgToReporter5(t *testing.T) {
//	var (
//		d      = New(conf.Conf)
//		c      = context.TODO()
//		err    error
//		rptMsg = &model.ReportMsg{
//			Aid:         1,
//			Title:       "test title5",
//			Did:         1,
//			Msg:         "dm msg5",
//			Uids:        "150781",
//			State:       model.StatSecondAutoDelete,
//			RptReason:   1,
//			BlockReason: 5,
//			Block:       -1,
//		}
//	)
//	err = d.SendMsgToReporter(c, rptMsg)
//	if err != nil {
//		t.Errorf("dao.SendMsgToReporter(rptMsg:%v) err(%v)", rptMsg, err)
//		t.Fail()
//	}
//}
//
//func TestDao_SendMsgToPoster1(t *testing.T) {
//	var (
//		d      = New(conf.Conf)
//		c      = context.TODO()
//		err    error
//		rptMsg = &model.ReportMsg{
//			Aid:         1,
//			Title:       "test title1",
//			Did:         1,
//			Msg:         "dm msg1",
//			Uids:        "150781",
//			State:       model.StatSecondAutoDelete,
//			RptReason:   1,
//			BlockReason: 5,
//			Block:       10,
//		}
//	)
//	err = d.SendMsgToPoster(c, rptMsg)
//	if err != nil {
//		t.Errorf("dao.SendMsgToPoster(rptMsg:%v) err(%v)", rptMsg, err)
//		t.Fail()
//	}
//}
func TestCreatePosterContent(t *testing.T) {
	msg := &model.ReportMsg{
		Aid:         1,
		Uids:        "1,2",
		Did:         11,
		Title:       "test title",
		Msg:         "test dm content",
		State:       1,
		RptReason:   16,
		Block:       0,
		BlockReason: 0,
	}
	Convey("test poster content", t, func() {
		_, err := testDao.createPosterContent(msg)
		So(err, ShouldBeNil)
	})
}
