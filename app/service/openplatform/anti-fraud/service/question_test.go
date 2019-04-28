package service

import (
	"context"
	"testing"

	"go-common/app/service/openplatform/anti-fraud/model"
	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	_qid      = 1527477525681
	_bid      = 1527233672941
	_page     = 0
	_pageSize = 20
)

func TestGetQusInfo(t *testing.T) {
	Convey("TestGetQusInfo", t, func() {
		res, err := svr.GetQusInfo(context.TODO(), 1527479929734)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestGetAnswerList(t *testing.T) {
	Convey("GetAnswerList", t, func() {
		res, err := svr.GetAnswerList(context.TODO(), _qid)
		data := make([]*model.Answer, 0)
		So(err, ShouldBeNil)
		So(res, ShouldResemble, data)
	})
}

func TestAddQus(t *testing.T) {

	anlist := []model.Answer{
		0: {QsID: 0, AnswerContent: "answer_content1", IsCorrect: 1, AnswerID: 0},
		1: {QsID: 0, AnswerContent: "answer_content2", IsCorrect: 1, AnswerID: 1},
		2: {QsID: 0, AnswerContent: "answer_content3", IsCorrect: 1, AnswerID: 0},
	}
	addmodel := &model.AddQus{
		BId:     _bid,
		Type:    1,
		Name:    "ceshi",
		Dif:     2,
		Answers: anlist,
		QsID:    0,
		AnType:  1,
	}

	Convey("AddQus", t, func() {
		res, err := svr.AddQus(context.TODO(), addmodel, anlist)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestDelQus(t *testing.T) {
	Convey("DelQus", t, func() {
		res, err := svr.DelQus(context.TODO(), _qid)
		So(err, ShouldBeIn, ecode.QusIDInvalid, ecode.BankUsing, ecode.ParamInvalid)
		So(res, ShouldBeIn, false, true)

	})
}

func TestGetQuslist(t *testing.T) {
	Convey("GetQuslist", t, func() {
		res, err := svr.GetQuslist(context.TODO(), _page, _pageSize, _bid)
		data := make([]*model.Answer, 0)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeIn, data)
	})
}

func TestGetQusTotal(t *testing.T) {
	Convey("GetQusTotal", t, func() {
		res, err := svr.GetQusTotal(context.TODO(), _bid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestUpdateQus(t *testing.T) {

	anlist := []model.Answer{
		{QsID: 0, AnswerContent: "answer_content1", IsCorrect: 1, AnswerID: 1},
		{QsID: 0, AnswerContent: "answer_content2", IsCorrect: 2, AnswerID: 2},
		{QsID: 0, AnswerContent: "answer_content3", IsCorrect: 3, AnswerID: 3},
	}
	list := &model.ArgUpdateQus{
		ArgBaseQus: model.ArgBaseQus{Type: 1, AnType: 1, Name: "wlt", BId: _bid, Dif: 1, Answer: []model.Answer{}}, QsID: _bid,
	}

	Convey("UpdateQus", t, func() {
		res, err := svr.UpdateQus(context.TODO(), list, anlist)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestCheckAnswer(t *testing.T) {
	var (
		qtype int8
		qid   int64
	)
	qid = 1527480273423
	answser := []model.Answer{
		{QsID: qid, AnswerContent: "xxx", IsCorrect: 0, AnswerID: 1527480273454},
		{QsID: qid, AnswerContent: "xxx", IsCorrect: 0, AnswerID: 1212},
	}
	Convey("CheckAnswer", t, func() {
		res, err := svr.CheckAnswer(context.TODO(), qid, qtype, answser)
		So(err, ShouldBeIn, ecode.QusIDInvalid, ecode.ParamInvalid, ecode.AnswerError, nil)
		So(res, ShouldNotBeNil)
	})
}

func TestRandQuestion(t *testing.T) {
	in := &model.ArgGetQuestion{
		UID:            "111111",
		TargetItem:     "1111",
		TargetItemType: 1,
		Source:         1,
		Platform:       1,
		ComponentID:    123,
	}

	var qusIds []int64
	qusIds = append(qusIds, _qid)
	Convey("CheckAnswer", t, func() {
		res, err := svr.randQuestion(context.TODO(), qusIds, in)
		So(err, ShouldBeIn, ecode.QusIDInvalid, ecode.AnswerError, ecode.BindBankNotFound, nil)
		So(res, ShouldBeNil)
	})
}

func TestGetQuestion(t *testing.T) {
	args := &model.ArgGetQuestion{
		UID:            "1",
		TargetItemType: 1,
		TargetItem:     "122",
		Source:         1,
		Platform:       1,
		ComponentID:    9999,
	}

	Convey("GetQuestion", t, func() {
		_, err := svr.GetQuestion(context.TODO(), args)
		So(err, ShouldBeIn, ecode.QusIDInvalid, ecode.AnswerError, ecode.BindBankNotFound, ecode.SameCompentErr, nil)
	})
}
