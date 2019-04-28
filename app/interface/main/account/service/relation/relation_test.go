package relation

import (
	"context"
	"sync"
	"testing"

	"go-common/app/interface/main/account/conf"
	mrl "go-common/app/service/main/relation/model"
	"go-common/library/log"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	once sync.Once
	//ip   = "127.0.0.1"
	s *Service
)

func startService() {
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	s = New(conf.Conf)
}

func TestRelation(t *testing.T) {
	once.Do(startService)
	Convey("relation", t, func() {
		testBlacks(t)
		testFollowers(t)
		testFollowings(t)
		testRelation(t)
		testRelations(t)
		testStat(t)
		testWhispers(t)
	})
}

func testBlacks(t *testing.T) {
	res, _, total, err := s.Blacks(context.TODO(), 18552813, 0, 1, 100)
	if err != nil {
		t.Errorf("s.Black err(%v)", err)
	}
	t.Logf("black %v, total:%v", res, total)
}

func testRelation(t *testing.T) {
	res, err := s.Relation(context.TODO(), 500, 100)
	if err != nil {
		t.Errorf("s.Relation err(%v)", err)
		return
	}
	t.Logf("Relation res(%v)", res)
}
func testRelations(t *testing.T) {
	res, err := s.Relations(context.TODO(), 500, []int64{100, 200})
	if err != nil {
		t.Errorf("s.Relations err(%v)", err)
		return
	}
	t.Logf("Relations res(%v)", res)
}

func testWhispers(t *testing.T) {
	res, _, err := s.Whispers(context.TODO(), 500, 1, 100, 0)
	if err != nil {
		t.Errorf("s.Whispers err(%v)", err)
		return
	}
	t.Logf("Whispers res(%v)", res)
}

func testFollowers(t *testing.T) {
	res, _, _, err := s.Followers(context.TODO(), 500, 1, 10, 1, 0)
	if err != nil {
		t.Errorf("s.Followers err(%v)", err)
		return
	}
	t.Logf("Followers res(%v)", res)
}

func testFollowings(t *testing.T) {
	res, _, _, err := s.Followings(context.TODO(), 500, 1, 10, 1, 0, "asc")
	if err != nil {
		t.Errorf("s.Followings err(%v)", err)
		return
	}
	t.Logf("Followings res(%v)", res)
}

func testStat(t *testing.T) {
	res, err := s.Stat(context.TODO(), 500, true)
	if err != nil {
		t.Errorf("s.Stat err(%v)", err)
		return
	}
	t.Logf("Stat  self res(%+v)", res)
	res, err = s.Stat(context.TODO(), 500, false)
	if err != nil {
		t.Errorf("s.Stat err(%v)", err)
		return
	}
	t.Logf("Stat res(%+v)", res)
}

func TestUnread(t *testing.T) {
	Convey("Unread", t, func() {
		_, err := s.Unread(context.TODO(), 1, false)
		So(err, ShouldBeNil)
	})
}

func TestUnreadCount(t *testing.T) {
	Convey("UnreadCount", t, func() {
		_, err := s.UnreadCount(context.TODO(), 1, false)
		So(err, ShouldBeNil)
	})
}

func TestSpecial(t *testing.T) {
	Convey("Special", t, func() {
		_, err := s.Special(context.TODO(), 1)
		So(err, ShouldBeNil)
	})
}

func TestDelSpecial(t *testing.T) {
	Convey("DelSpecial", t, func() {
		err := s.DelSpecial(context.TODO(), &mrl.ArgFollowing{Mid: 1})
		So(err, ShouldBeNil)
	})
}

func TestAddSpecial(t *testing.T) {
	Convey("AddSpecial", t, func() {
		err := s.AddSpecial(context.TODO(), &mrl.ArgFollowing{Mid: 1})
		So(err, ShouldBeNil)
	})
}

func TestClosePrompt(t *testing.T) {
	Convey("ClosePrompt", t, func() {
		err := s.ClosePrompt(context.TODO(), &mrl.ArgPrompt{Mid: 1})
		So(err, ShouldBeNil)
	})
}

func TestPrompt(t *testing.T) {
	Convey("Prompt", t, func() {
		_, err := s.Prompt(context.TODO(), &mrl.ArgPrompt{Mid: 1})
		So(err, ShouldBeNil)
	})
}

func TestTagsMoveUsers(t *testing.T) {
	Convey("TagsMoveUsers", t, func() {
		err := s.TagsMoveUsers(context.TODO(), 1, 1, "", "")
		So(err, ShouldBeNil)
	})
}

func TestTagsCopyUsers(t *testing.T) {
	Convey("TagsCopyUsers", t, func() {
		err := s.TagsCopyUsers(context.TODO(), 1, "", "")
		So(err, ShouldBeNil)
	})
}
