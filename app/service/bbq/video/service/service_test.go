package service

import (
	"context"
	"flag"
	"go-common/app/service/bbq/video/api/grpc/v1"
	"go-common/app/service/bbq/video/conf"
	"os"
	"testing"
)

var (
	s *Service
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "")
		flag.Set("conf_token", "")
		flag.Set("tree_id", "")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/")
		flag.Set("conf_name", "test.toml")
		flag.Set("deploy.env", "uat")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	conf.Conf.Log.Stdout = false
	s = New(conf.Conf)
	os.Exit(m.Run())
}

func BenchmarkService_CreateID(b *testing.B) {
	req := &v1.CreateIDRequest{Mid: 3}
	for i := 0; i < b.N; i++ {
		s.CreateID(context.Background(), req)
	}
}
