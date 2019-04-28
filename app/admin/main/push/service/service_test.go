package service

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-common/app/admin/main/push/conf"
	"go-common/app/admin/main/push/model"
	pushmdl "go-common/app/service/main/push/model"

	. "github.com/smartystreets/goconvey/convey"
)

var svr *Service

func init() {
	dir, _ := filepath.Abs("../cmd/push-admin-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	svr = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() {})
		f(svr)
	}
}

func Test_Service(t *testing.T) {
	Convey("service test", t, WithService(func(s *Service) {
		s.Wait()
		s.Close()
	}))
}

func Test_AddTaskAll(t *testing.T) {
	Convey("get last report id", t, WithService(func(s *Service) {
		var maxID int64
		err := s.dao.DB.Raw("SELECT MAX(id) FROM push_reports").Row().Scan(&maxID)
		So(err, ShouldBeNil)
		t.Logf("maxid(%d)", maxID)
	}))

	Convey("get tokens", t, WithService(func(s *Service) {
		sqlStr := fmt.Sprintf("SELECT platform_id,device_token FROM push_reports WHERE id>%d and id<=%d and app_id=%d and dtime=0 and notify_switch=1", 1, 100, 1)
		rows, err := s.dao.DB.Raw(sqlStr).Rows()
		So(err, ShouldBeNil)
		var tokens []string
		for rows.Next() {
			var (
				platformID int
				token      string
			)
			err = rows.Scan(&platformID, &token)
			So(err, ShouldBeNil)
			tokens = append(tokens, fmt.Sprintf("%d\t%s", platformID, token))
		}
		t.Logf("tokens(%d)", len(tokens))
		if len(tokens) > 0 {
			t.Logf("token one (%s)", tokens[0])
		}
	}))
}

func Test_CheckUploadMid(t *testing.T) {
	Convey("CheckUploadMid", t, WithService(func(s *Service) {
		data := []byte("1\n2\n3")
		err := s.CheckUploadMid(context.TODO(), data)
		So(err, ShouldBeNil)
		data = []byte("1\nabc\n3")
		err = s.CheckUploadMid(context.TODO(), data)
		So(err, ShouldNotBeNil)
		t.Logf("check mid error(%v)", err)
	}))
}

func Test_CheckUploadToken(t *testing.T) {
	Convey("CheckUploadToken", t, WithService(func(s *Service) {
		data := []byte("2	fdsahjfkdshaj\n3	hjkhjhjkhj")
		err := s.CheckUploadToken(context.TODO(), data)
		So(err, ShouldBeNil)

		data = []byte("2	fdsahjfkdshaj\n3	hjkhjhjkhj\n4\n")
		err = s.CheckUploadToken(context.TODO(), data)
		So(err, ShouldNotBeNil)
		t.Logf("check token error(%v)", err)
	}))
}

func Test_parseQuery(t *testing.T) {
	Convey("parse query", t, WithService(func(s *Service) {
		str := `{"age":1,"sex":1,"is_up":0,"is_formal_member":1,"user_active_day":0,"user_new_day":0,"user_silentDay":0,"area":"2","level":"2,3,4","platforms":"1,2,3","like":"1,2,3","channel":"huawei,xiaomi","vip_expire":[{"begin":"2018-07-01 00:00:00","end":"2018-07-27 00:00:00"}],"self_attention":null,"self_attention_type":0,"active":null,"ActivePeriod":0}`
		p := new(model.DPParams)
		err := json.Unmarshal([]byte(str), p)
		So(err, ShouldBeNil)

		p.Area = pushmdl.SplitInts(p.AreaStr)
		p.Level = pushmdl.SplitInts(p.LevelStr)
		p.Platforms = pushmdl.SplitInts(p.PlatformStr)
		p.Like = pushmdl.SplitInts(p.LikeStr)
		p.Channel = strings.Split(p.ChannelStr, ",")
		if p.VipExpireStr != "" {
			err = json.Unmarshal([]byte(p.VipExpireStr), &p.VipExpires)
			So(err, ShouldBeNil)
		}
		if p.AttentionStr != "" {
			err = json.Unmarshal([]byte(p.AttentionStr), &p.Attentions)
			So(err, ShouldBeNil)
		}
		if p.ActivePeriodStr != "" {
			err = json.Unmarshal([]byte(p.ActivePeriodStr), &p.ActivePeriods)
			So(err, ShouldBeNil)
		}
		sql := s.parseQuery(pushmdl.TaskTypeDataPlatformMid, p)
		t.Logf("sql(%s)", sql)
	}))
}

func TestCheckDpData(t *testing.T) {
	Convey("test check data platform data", t, WithService(func(s *Service) {
		err := s.CheckDpData(context.Background())
		So(err, ShouldBeNil)
	}))
}
