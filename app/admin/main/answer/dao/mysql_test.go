package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/app/admin/main/answer/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDaoIDsByState(t *testing.T) {
	Convey("IDsByState", t, func() {
		_, err := d.IDsByState(context.TODO())
		So(err, ShouldBeNil)
	})
}

func TestDaoQuestionAll(t *testing.T) {
	var id int64 = 1
	q := &model.QuestionDB{ID: id, State: 1}
	Convey("QuestionAdd", t, func() {
		aff, err := d.QuestionAdd(context.TODO(), q)
		So(err, ShouldBeNil)
		So(aff, ShouldNotBeNil)
	})
	Convey("QueByID", t, func() {
		que, err := d.QueByID(context.TODO(), id)
		So(err, ShouldBeNil)
		So(que, ShouldNotBeNil)
	})
	Convey("QuestionEdit", t, func() {
		q.Question = "testone"
		aff, err := d.QuestionEdit(context.TODO(), q)
		So(err, ShouldBeNil)
		So(aff, ShouldNotBeNil)
	})
	Convey("ByIDs", t, func() {
		res, err := d.ByIDs(context.TODO(), []int64{id})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})

}

func TestDaoQuestionPage(t *testing.T) {
	arg := &model.ArgQue{State: 1}
	Convey("QuestionList", t, func() {
		res, err := d.QuestionList(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
	Convey("QuestionCount", t, func() {
		res, err := d.QuestionCount(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestDaoUpdateStatus(t *testing.T) {
	Convey("UpdateStatus", t, func() {
		aff, err := d.UpdateStatus(context.TODO(), 0, 0, "")
		So(err, ShouldBeNil)
		So(aff, ShouldNotBeNil)
	})
}

func TestDaoTypes(t *testing.T) {
	Convey("Types", t, func() {
		res, err := d.Types(context.TODO())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestDaoTypeAdd(t *testing.T) {
	Convey("Types", t, func() {
		var allType = []*model.TypeInfo{
			{ID: 1, Parentid: 0, Name: "游戏"},
			{ID: 2, Parentid: 0, Name: "影视"},
			{ID: 3, Parentid: 0, Name: "科技"},
			{ID: 4, Parentid: 0, Name: "动画"},
			{ID: 5, Parentid: 0, Name: "艺术"},
			{ID: 6, Parentid: 0, Name: "流行前线"},
			{ID: 7, Parentid: 0, Name: "鬼畜"},
			{ID: 8, Parentid: 1, Name: "动作射击", LabelName: "游戏"},
			{ID: 9, Parentid: 1, Name: "冒险格斗", LabelName: "游戏"},
			{ID: 12, Parentid: 1, Name: "策略模拟 ", LabelName: "游戏"},
			{ID: 13, Parentid: 1, Name: "角色扮演 ", LabelName: "游戏"},
			{ID: 14, Parentid: 1, Name: "音乐体育 ", LabelName: "游戏"},
			{ID: 15, Parentid: 2, Name: "纪录片 ", LabelName: "影视"},
			{ID: 16, Parentid: 2, Name: "电影 ", LabelName: "影视"},
			{ID: 17, Parentid: 2, Name: "电视剧 ", LabelName: "影视"},
			{ID: 18, Parentid: 3, Name: "军事 ", LabelName: "科技"},
			{ID: 19, Parentid: 3, Name: "地理 ", LabelName: "科技"},
			{ID: 20, Parentid: 3, Name: "历史 ", LabelName: "科技"},
			{ID: 21, Parentid: 3, Name: "文学 ", LabelName: "科技"},
			{ID: 22, Parentid: 3, Name: "数学 ", LabelName: "科技"},
			{ID: 23, Parentid: 3, Name: "物理 ", LabelName: "科技"},
			{ID: 24, Parentid: 3, Name: "化学 ", LabelName: "科技"},
			{ID: 25, Parentid: 3, Name: "生物 ", LabelName: "科技"},
			{ID: 26, Parentid: 3, Name: "数码科技 ", LabelName: "科技"},
			{ID: 27, Parentid: 4, Name: "动画声优 ", LabelName: "动画"},
			{ID: 28, Parentid: 4, Name: "动漫内容 ", LabelName: "动画"},
			{ID: 29, Parentid: 5, Name: "ACG音乐 ", LabelName: "艺术"},
			{ID: 30, Parentid: 5, Name: "三次元音乐 ", LabelName: "艺术"},
			{ID: 31, Parentid: 5, Name: "绘画 ", LabelName: "艺术"},
			{ID: 32, Parentid: 6, Name: "娱乐 ", LabelName: "流行前线"},
			{ID: 33, Parentid: 6, Name: "时尚 ", LabelName: "流行前线"},
			{ID: 34, Parentid: 6, Name: "运动 ", LabelName: "流行前线"},
			{ID: 35, Parentid: 7, Name: "鬼畜 ", LabelName: "鬼畜"},
			{ID: 36, Parentid: 0, Name: "基础题", LabelName: "基础题"},
		}
		for _, v := range allType {
			res, err := d.TypeSave(context.TODO(), v)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
		}
	})
}

func TestDaoBaseQS(t *testing.T) {
	Convey("TestBaseQS", t, func() {
		res, err := d.BaseQS(context.Background())
		for k, re := range res {
			fmt.Printf("k:%d re:%+v	\n", k, re)
		}
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestDaoInsBaseQs(t *testing.T) {
	var qs = &model.QuestionDB{
		Mid:       1,
		IP:        "127.0.0.1",
		Question:  "qs",
		Ans1:      "Ans1",
		Ans2:      "Ans2",
		Ans3:      "Ans3",
		Ans4:      "Ans4",
		Tips:      "tips",
		AvID:      1,
		MediaType: 1,
		Source:    1,
		Ctime:     time.Now(),
		Mtime:     time.Now(),
		Operator:  "operator",
	}
	Convey("TestInsBaseQs", t, func() {
		af, err := d.InsBaseQs(context.Background(), qs)
		fmt.Printf("%+v \n", af)
		So(err, ShouldBeNil)
	})
}

func TestAllQS(t *testing.T) {
	Convey("TestAllQS", t, func() {
		res, err := d.AllQS(context.Background())
		for k, re := range res {
			fmt.Printf("k:%d re:%+v	\n", k, re)
		}
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
