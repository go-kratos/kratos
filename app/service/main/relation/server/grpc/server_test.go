package grpc

import (
	"context"
	"flag"
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"go-common/app/service/main/relation/model"
	"os"
	"testing"
	"time"

	"go-common/app/service/main/relation/api"
	pb "go-common/app/service/main/relation/api"
	"go-common/app/service/main/relation/conf"
	"go-common/app/service/main/relation/service"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"
)

var (
	cli api.RelationClient
	svr *service.Service
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.relation-service")
		flag.Set("conf_token", "8hm3I5rWzuhChxrBI6VTqmCs7TpJwFhO")
		flag.Set("tree_id", "2139")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	svr = service.New(conf.Conf)
	cfg := &warden.ClientConfig{
		Dial:    xtime.Duration(time.Second * 3),
		Timeout: xtime.Duration(time.Second * 3),
	}
	var err error
	cli, err = api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	m.Run()
	os.Exit(0)
}

func TestServerRelation(t *testing.T) {
	var (
		c = context.Background()
	)

	convey.Convey("Relation", t, func(cv convey.C) {
		rr := pb.RelationReq{Mid: 1, Fid: 2, RealIp: "127.0.0.1"}

		fr, err := cli.Relation(c, &rr)
		cv.Convey("Then err should be nil.reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(fr, convey.ShouldNotBeNil)
		})
		fmt.Println(fr)

		f, err2 := svr.Relation(c, 1, 2)
		cv.Convey("Then err should be nil. reslut should not be nil.", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(f, convey.ShouldNotBeNil)
		})
		fmt.Println(f)

		cv.Convey("the grpc result should be equal to service result", func(cv convey.C) {
			cv.So(f.Mid, convey.ShouldEqual, fr.Mid)
			cv.So(f.Attribute, convey.ShouldEqual, fr.Attribute)
			cv.So(f.CTime, convey.ShouldEqual, fr.CTime)
		})
	})
}

func TestServerRelations(t *testing.T) {
	var (
		c = context.Background()
	)

	convey.Convey("Relations", t, func(cv convey.C) {
		rr := pb.RelationsReq{Mid: 1, Fid: []int64{2, 3}, RealIp: "127.0.0.1"}

		fr, err := cli.Relations(c, &rr)
		cv.Convey("Then err should be nil.reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(fr, convey.ShouldNotBeNil)
		})
		fmt.Println(fr)

		f, err2 := svr.Relations(c, 1, []int64{2, 3})
		cv.Convey("Then err should be nil. reslut should not be nil.", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(f, convey.ShouldNotBeNil)
		})
		fmt.Println(f)

		cv.Convey("the grpc result should be equal to service result", func(cv convey.C) {
			f1 := fr.FollowingMap[2]
			cv.So(f[2].Attribute, convey.ShouldEqual, f1.Attribute)
			cv.So(f[2].CTime, convey.ShouldEqual, f1.CTime)
		})
	})
}

func TestServerStat(t *testing.T) {
	var (
		c = context.Background()
	)

	convey.Convey("Stat", t, func(cv convey.C) {
		st, err := cli.Stat(c, &pb.MidReq{Mid: 2})
		cv.Convey("Then err should be nil.reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(st, convey.ShouldNotBeNil)
		})
		fmt.Println(st)

		s, err2 := svr.Stat(c, 2)
		cv.Convey("Then err should be nil. reslut should not be nil.", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(s, convey.ShouldNotBeNil)
		})
		fmt.Println(s)

		cv.Convey("the grpc result should be equal to service result", func(cv convey.C) {
			cv.So(st.CTime, convey.ShouldEqual, s.CTime)
			cv.So(st.Mid, convey.ShouldEqual, s.Mid)
			cv.So(st.Black, convey.ShouldEqual, s.Black)
			cv.So(st.Follower, convey.ShouldEqual, s.Follower)
			cv.So(st.Following, convey.ShouldEqual, s.Following)
			cv.So(st.MTime, convey.ShouldEqual, s.MTime)
			cv.So(st.Whisper, convey.ShouldEqual, s.Whisper)
		})
	})
}

func TestServerStats(t *testing.T) {
	var (
		c = context.Background()
	)

	convey.Convey("Stats", t, func(cv convey.C) {
		st, err := cli.Stats(c, &pb.MidsReq{Mids: []int64{2, 3}, RealIp: "127.0.0.1"})
		cv.Convey("Then err should be nil.reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(st, convey.ShouldNotBeNil)
		})
		fmt.Println(st)

		s, err2 := svr.Stats(c, []int64{2, 3})
		cv.Convey("Then err should be nil. reslut should not be nil.", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(s, convey.ShouldNotBeNil)
		})
		fmt.Println(s)

		cv.Convey("the grpc result should be equal to service result", func(cv convey.C) {
			st2 := st.StatReplyMap[2]
			s2 := s[2]

			cv.So(s2.CTime, convey.ShouldEqual, st2.CTime)
			cv.So(s2.Mid, convey.ShouldEqual, st2.Mid)
			cv.So(s2.Black, convey.ShouldEqual, st2.Black)
			cv.So(s2.Follower, convey.ShouldEqual, st2.Follower)
			cv.So(s2.Following, convey.ShouldEqual, st2.Following)
			cv.So(s2.MTime, convey.ShouldEqual, st2.MTime)
			cv.So(s2.Whisper, convey.ShouldEqual, st2.Whisper)
		})
	})
}

func TestServerAttentions(t *testing.T) {
	var (
		c = context.Background()
	)

	convey.Convey("Attentions", t, func(cv convey.C) {
		at, err := cli.Attentions(c, &pb.MidReq{Mid: 2})
		cv.Convey("Then err should be nil.reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(at, convey.ShouldNotBeNil)
		})
		fmt.Println(at)

		a, err2 := svr.Attentions(c, 2)
		cv.Convey("Then err should be nil. reslut should not be nil.", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(a, convey.ShouldNotBeNil)
		})
		fmt.Println(a)

		cv.Convey("the grpc result should be equal to service result", func(cv convey.C) {
			f1 := at.FollowingList[0]
			f2 := a[0]
			cv.So(f1.Mid, convey.ShouldEqual, f2.Mid)
			cv.So(f1.Attribute, convey.ShouldEqual, f2.Attribute)
			cv.So(f1.CTime, convey.ShouldEqual, f2.CTime)
		})
	})
}

func TestServerFollowings(t *testing.T) {
	var (
		c = context.Background()
	)

	convey.Convey("Followings", t, func(cv convey.C) {
		at, err := cli.Followings(c, &pb.MidReq{Mid: 2})
		cv.Convey("Then err should be nil.reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(at, convey.ShouldNotBeNil)
		})
		fmt.Println(at)

		a, err2 := svr.Followings(c, 2)
		cv.Convey("Then err should be nil. reslut should not be nil.", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(a, convey.ShouldNotBeNil)
		})
		fmt.Println(a)

		cv.Convey("the grpc result should be equal to service result", func(cv convey.C) {
			f1 := at.FollowingList[0]
			f2 := a[0]
			cv.So(f1.Mid, convey.ShouldEqual, f2.Mid)
			cv.So(f1.Attribute, convey.ShouldEqual, f2.Attribute)
			cv.So(f1.CTime, convey.ShouldEqual, f2.CTime)
		})
	})
}

func TestServerWhispers(t *testing.T) {
	var (
		c = context.Background()
	)

	convey.Convey("Whispers", t, func(cv convey.C) {
		at, err := cli.Whispers(c, &pb.MidReq{Mid: 2231365})
		cv.Convey("Then err should be nil.reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(at, convey.ShouldNotBeNil)
		})
		fmt.Println(at)

		a, err2 := svr.Whispers(c, 2231365)
		cv.Convey("Then err should be nil. reslut should not be nil.", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(a, convey.ShouldNotBeNil)
		})
		fmt.Println(a)

		cv.Convey("the grpc result should be equal to service result", func(cv convey.C) {
			f1 := at.FollowingList[0]
			f2 := a[0]
			cv.So(f1.Mid, convey.ShouldEqual, f2.Mid)
			cv.So(f1.Attribute, convey.ShouldEqual, f2.Attribute)
			cv.So(f1.CTime, convey.ShouldEqual, f2.CTime)
		})
	})
}

func TestServerFollowers(t *testing.T) {
	var (
		c = context.Background()
	)

	convey.Convey("Followers", t, func(cv convey.C) {
		at, err := cli.Followers(c, &pb.MidReq{Mid: 2})
		cv.Convey("Then err should be nil.reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(at, convey.ShouldNotBeNil)
		})
		fmt.Println(at)

		a, err2 := svr.Followers(c, 2)
		cv.Convey("Then err should be nil. reslut should not be nil.", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(a, convey.ShouldNotBeNil)
		})
		fmt.Println(a)

		cv.Convey("the grpc result should be equal to service result", func(cv convey.C) {
			f1 := at.FollowingList[0]
			f2 := a[0]
			cv.So(f1.Mid, convey.ShouldEqual, f2.Mid)
			cv.So(f1.Attribute, convey.ShouldEqual, f2.Attribute)
			cv.So(f1.CTime, convey.ShouldEqual, f2.CTime)
		})
	})
}

func TestServerTag(t *testing.T) {
	var (
		c = context.Background()
	)

	convey.Convey("Tag", t, func(cv convey.C) {
		t, err := cli.Tag(c, &pb.TagIdReq{Mid: 1, TagId: -10, RealIp: "127.0.0.1"})
		cv.Convey("Then err should be nil.reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(t, convey.ShouldNotBeNil)
		})
		fmt.Println(t)

		tg, err2 := svr.Tag(c, 1, -10, "127.0.0.1")
		cv.Convey("Then err should be nil. reslut should not be nil.", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(tg, convey.ShouldNotBeNil)
		})
		fmt.Println(tg)

		cv.Convey("the grpc result should be equal to service result", func(cv convey.C) {
			cv.So(len(t.Mids), convey.ShouldEqual, len(tg))
			cv.So(t.Mids[0], convey.ShouldEqual, tg[0])
		})
	})
}

func TestServerTags(t *testing.T) {
	var (
		c = context.Background()
	)

	convey.Convey("Tags", t, func(cv convey.C) {
		t, err := cli.Tags(c, &pb.MidReq{Mid: 1})
		cv.Convey("Then err should be nil.reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(t, convey.ShouldNotBeNil)
		})
		fmt.Println(t)

		tg, err2 := svr.Tags(c, 1, "127.0.0.1")
		cv.Convey("Then err should be nil. reslut should not be nil.", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(tg, convey.ShouldNotBeNil)
		})
		fmt.Println(tg)

		cv.Convey("the grpc result should be equal to service result", func(cv convey.C) {
			cv.So(t.TagCountList[0].Tagid, convey.ShouldEqual, tg[0].Tagid)
			cv.So(t.TagCountList[0].Count, convey.ShouldEqual, tg[0].Count)
			cv.So(t.TagCountList[0].Name, convey.ShouldEqual, tg[0].Name)
		})
	})
}

func TestUserTag(t *testing.T) {
	var (
		c = context.Background()
	)

	convey.Convey("UserTag", t, func(cv convey.C) {
		tt, err := cli.UserTag(c, &pb.RelationReq{Mid: 1, Fid: 2, RealIp: "127.0.0.1"})
		cv.Convey("Then err should be nil.reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(tt, convey.ShouldNotBeNil)
		})
		fmt.Println(tt)

		tg, err2 := svr.UserTag(c, 1, 2, "127.0.0.1")
		cv.Convey("Then err should be nil. reslut should not be nil.", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(tg, convey.ShouldNotBeNil)
		})
		fmt.Println(tg)

		cv.Convey("the grpc result should be equal to service result", func(cv convey.C) {
			cv.So(len(tt.Tags), convey.ShouldEqual, len(tg))
		})
	})
}

func TestSpecial(t *testing.T) {
	var (
		c = context.Background()
	)

	convey.Convey("Special", t, func(cv convey.C) {
		s, err := cli.Special(c, &pb.MidReq{Mid: 1})
		cv.Convey("Then err should be nil.reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(s, convey.ShouldNotBeNil)
		})
		fmt.Println(s)

		sg, err2 := svr.Special(c, 1)
		cv.Convey("Then err should be nil. reslut should not be nil.", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(sg, convey.ShouldNotBeNil)
		})
		fmt.Println(sg)

		cv.Convey("the grpc result should be equal to service result", func(cv convey.C) {
			cv.So(len(sg), convey.ShouldEqual, len(s.Mids))
		})
	})
}

func TestFollowersUnread(t *testing.T) {
	var (
		c = context.Background()
	)

	convey.Convey("FollowersUnread", t, func(cv convey.C) {
		s, err := cli.FollowersUnread(c, &pb.MidReq{Mid: 1})
		cv.Convey("Then err should be nil.reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(s, convey.ShouldNotBeNil)
		})
		fmt.Println(s)

		h, err2 := svr.Unread(c, 1)
		cv.Convey("Then err should be nil. reslut should not be nil.", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(h, convey.ShouldNotBeNil)
		})
		fmt.Println(h)

		cv.Convey("the grpc result should be equal to service result", func(cv convey.C) {
			cv.So(s.HasUnread, convey.ShouldEqual, h)
		})
	})
}

func TestFollowersUnreadCount(t *testing.T) {
	var (
		c = context.Background()
	)

	convey.Convey("FollowersUnreadCount", t, func(cv convey.C) {
		s, err := cli.FollowersUnreadCount(c, &pb.MidReq{Mid: 1})
		cv.Convey("Then err should be nil.reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(s, convey.ShouldNotBeNil)
		})
		fmt.Println(s)

		h, err2 := svr.UnreadCount(c, 1)
		cv.Convey("Then err should be nil. reslut should not be nil.", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(h, convey.ShouldNotBeNil)
		})
		fmt.Println(h)

		cv.Convey("the grpc result should be equal to service result", func(cv convey.C) {
			cv.So(s.UnreadCount, convey.ShouldEqual, h)
		})
	})
}

func TestAchieveGet(t *testing.T) {
	var (
		c = context.Background()
	)

	convey.Convey("FollowersUnreadCount", t, func(cv convey.C) {
		s, err := cli.AchieveGet(c, &pb.AchieveGetReq{Award: "10k", Mid: 3})
		cv.Convey("Then err should be nil.reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(s, convey.ShouldNotBeNil)
		})
		fmt.Println(s)

		s2, err := cli.Achieve(c, &pb.AchieveReq{AwardToken: s.AwardToken})
		cv.Convey("Then err should be nil.  reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(s2, convey.ShouldNotBeNil)
		})
		fmt.Println(s2)

		h, err2 := svr.AchieveGet(c, &model.ArgAchieveGet{Award: "10k", Mid: 3})
		cv.Convey("Then err should be nil. reslut should not be nil.", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(h, convey.ShouldNotBeNil)
		})
		fmt.Println(h)

		h2, err2 := svr.Achieve(c, &model.ArgAchieve{AwardToken: h.AwardToken})
		cv.Convey("Then err should be nil  reslut should not be nil. ", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(h2, convey.ShouldNotBeNil)
		})
		fmt.Println(h2)

		cv.Convey("Then err should be nil reslut should not be nil.", func(cv convey.C) {
			cv.So(s2.Mid, convey.ShouldEqual, h2.Mid)
			cv.So(s2.Award, convey.ShouldEqual, h2.Award)
		})
	})
}

func TestFollowerNotifySetting(t *testing.T) {
	var (
		c = context.Background()
	)

	convey.Convey("FollowerNotifySetting", t, func(cv convey.C) {
		s, err := cli.FollowerNotifySetting(c, &pb.MidReq{Mid: 3})
		cv.Convey("Then err should be nil.reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(s, convey.ShouldNotBeNil)
		})
		fmt.Println(s)

		h, err2 := svr.FollowerNotifySetting(c, &model.ArgMid{Mid: 3})
		cv.Convey("Then err should be nil. reslut should not be nil.", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(h, convey.ShouldNotBeNil)
		})
		fmt.Println(h)

		cv.Convey("the grpc result should be equal to service result", func(cv convey.C) {
			cv.So(s.Mid, convey.ShouldEqual, h.Mid)
			cv.So(s.Enabled, convey.ShouldEqual, h.Enabled)
		})
	})
}

func TestSameFollowings(t *testing.T) {
	var (
		c = context.Background()
	)

	convey.Convey("SameFollowings", t, func(cv convey.C) {
		s, err := cli.SameFollowings(c, &pb.SameFollowingReq{Mid: 3, Mid2: 4})
		cv.Convey("Then err should be nil.reslut should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(s, convey.ShouldNotBeNil)
		})
		fmt.Println(s)

		h, err2 := svr.SameFollowings(c, &model.ArgSameFollowing{Mid1: 3, Mid2: 4})
		cv.Convey("Then err should be nil. reslut should not be nil.", func(cv convey.C) {
			cv.So(err2, convey.ShouldBeNil)
			cv.So(h, convey.ShouldNotBeNil)
		})
		fmt.Println(h)

		cv.Convey("the grpc result should be equal to service result", func(cv convey.C) {
			cv.So(len(s.FollowingList), convey.ShouldEqual, len(h))
			f1 := s.FollowingList[0]
			f2 := h[0]
			cv.So(f1.Mid, convey.ShouldEqual, f2.Mid)
			cv.So(f1.MTime, convey.ShouldEqual, f2.MTime)
			cv.So(f1.CTime, convey.ShouldEqual, f2.CTime)
			cv.So(f1.Attribute, convey.ShouldEqual, f2.Attribute)
			cv.So(f1.Source, convey.ShouldEqual, f2.Source)
			cv.So(f1.Special, convey.ShouldEqual, f2.Special)
		})
	})
}
