package service

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"go-common/app/job/main/answer/conf"
	"go-common/app/job/main/answer/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "")
		flag.Set("conf_appid", "")
		flag.Set("conf_token", "")
		flag.Set("tree_id", "")
		flag.Set("conf_version", "server-1")
		flag.Set("deploy_env", "dev")
		flag.Set("conf_env", "10")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/answer-job-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
	os.Exit(m.Run())
}

func TestUpdateBfs(t *testing.T) {
	Convey("create img", t, func(s *Service) {
		err := s.createBFSImg(context.TODO(), &model.LabourQs{
			Question: "测试附加题3",
			ID:       3,
		})
		fmt.Println("err:", err)
		So(err, ShouldBeNil)
	})
}

func TestPing(t *testing.T) {
	Convey("Test_Ping", t, func(s *Service) {
		err := s.Ping(context.TODO())
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestAddQue
func TestCreateBFSImg(t *testing.T) {
	Convey("TestAddQue add question data", t, func(s *Service) {
		for i := 90; i < 100; i++ {
			que := &model.LabourQs{
				Question: fmt.Sprintf("测试d=====(￣▽￣*)b厉害 %d", i),
				Ans:      int8(i%2 + 1),
				AvID:     int64(i),
				Status:   int8(i%2 + 1),
				Source:   int8(1),
				State:    model.HadCreateImg,
				ID:       int64(i),
			}
			err := s.createBFSImg(context.TODO(), que)
			So(err, ShouldBeNil)
			err = s.dao.AddQs(context.TODO(), que)
			So(err, ShouldBeNil)
		}
	})
}
