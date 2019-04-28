package main

import (
	"context"
	"flag"

	"go-common/app/service/ep/footman/conf"
	"go-common/app/service/ep/footman/service"
	"go-common/library/cache/memcache"
	"go-common/library/container/pool"
	"go-common/library/database/orm"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
	"go-common/library/net/netutil/breaker"
	"go-common/library/time"
)

func main() {
	var (
		versionPath string
		cookiePath  string
		tokenPath   string
		action      string
	)

	flag.StringVar(&versionPath, "v", "", "版本批次文件路径")
	flag.StringVar(&cookiePath, "c", "", "cookie文件路径")
	flag.StringVar(&tokenPath, "t", "", "token文件路径")
	flag.StringVar(&action, "a", "", "操作类型")
	flag.Parse()

	c := &conf.Config{
		HTTPClient: &xhttp.ClientConfig{
			App: &xhttp.App{
				Key:    "c05dd4e1638a8af0",
				Secret: "7daa7f8c06cd33c5c3067063c746fdcb",
			},
			Dial:      time.Duration(20000000000),
			Timeout:   time.Duration(100000000000),
			KeepAlive: time.Duration(600000000000),
			Breaker: &breaker.Config{
				Window:  time.Duration(100000000000),
				Sleep:   time.Duration(20000000000),
				Bucket:  10,
				Ratio:   0.5,
				Request: 100,
			},
		},
		Bugly: &conf.BuglyConf{
			Host:    "https://bugly.qq.com",
			Cookie:  cookiePath,
			Token:   tokenPath,
			Version: versionPath,
		},
		ORM: &orm.Config{
			DSN:         "root:123456@tcp(172.18.33.130:3306)/footman?timeout=200ms&readTimeout=2000ms&writeTimeout=2000ms&parseTime=true&loc=Local&charset=utf8,utf8mb4",
			Active:      5,
			Idle:        5,
			IdleTimeout: time.Duration(20000000000),
		},
		Mail: &conf.Mail{
			Host:        "smtp.exmail.qq.com",
			Port:        465,
			Username:    "merlin@bilibili.com",
			Password:    "",
			NoticeOwner: []string{"fengyifeng@bilibili.com"},
		},
		Memcache: &conf.Memcache{
			Expire: time.Duration(10000000),
			Config: &memcache.Config{
				Name:         "merlin",
				Proto:        "tcp",
				Addr:         "172.22.33.137:11216",
				DialTimeout:  time.Duration(1000),
				ReadTimeout:  time.Duration(1000),
				WriteTimeout: time.Duration(1000),
				Config: &pool.Config{
					Active:      10,
					IdleTimeout: time.Duration(1000),
				},
			},
		},
		Bugly2Tapd: &conf.Bugly2Tapd{
			ProjectIds: []string{"900028525"},
		},
	}
	s := service.New(c)
	log.Info("v1.0.40")

	switch action {
	case "insertTapd":
		s.BuglyInsertTapd(context.Background())
	default:
		s.GetSaveIssuesWithMultiVersion(context.Background())
		s.UpdateBuglyStatusInTapd(context.Background())
		s.UpdateBugInTapd(context.Background())
	}
	defer s.Close()

}
