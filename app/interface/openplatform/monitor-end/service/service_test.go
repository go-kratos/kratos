package service

import (
	"context"

	"go-common/app/interface/openplatform/monitor-end/conf"
	"go-common/app/interface/openplatform/monitor-end/model/kafka"
	"go-common/app/interface/openplatform/monitor-end/model/monitor"
	"go-common/library/cache/redis"
	"go-common/library/container/pool"
	"go-common/library/database/sql"
	"go-common/library/log/infoc"
	"go-common/library/time"
)

var (
	ctx = context.Background()
	svr *Service
)

type TestData map[string]string

type TestCase struct {
	tag      string
	testData string
	expected int
}

/*[mysql]
	dsn = "root:123456@tcp(172.16.33.203:3306)/public_monitor?timeout=500s&readTimeout=500s&writeTimeout=500s&parseTime=true&loc=Local&charset=utf8,utf8mb4"
    active = 5
    idle = 2
    idleTimeout ="4h"
    queryTimeout = "1000s"
    execTimeout = "200s"#
    tranTimeout = "2000s"

[redis]
    name = "article"
    proto = "tcp"
    addr = "172.16.33.203:6379"
    idle = 10
    active = 10
    dialTimeout = "1s"
    readTimeout = "1s"
    writeTimeout = "1s"
    idleTimeout = "10s"*/
func init() {
	c := conf.Conf
	if c.Monitor == nil {
		c = &conf.Config{
			Monitor: &monitor.MonitorConfig{Proto: "tcp", Addr: "127.0.0.1:9988"},
			Kafka: &kafka.Config{
				Addr:  []string{"1.1.1.1"},
				Topic: "test_topic",
			},
			NeedConsume: false,
			Redis: &redis.Config{
				Name:  "article",
				Proto: "tcp",
				Addr:  "172.16.33.203:6379",
				Config: &pool.Config{
					Idle:   2,
					Active: 5,
				},
				DialTimeout:  time.Duration(int64(1000000000)),
				ReadTimeout:  time.Duration(int64(1000000000)),
				WriteTimeout: time.Duration(int64(1000000000)),
			},
			MySQL: &sql.Config{
				DSN:          "root:123456@tcp(172.16.33.203:3306)/public_monitor?timeout=500s&readTimeout=500s&writeTimeout=500s&parseTime=true&loc=Local&charset=utf8,utf8mb4",
				QueryTimeout: time.Duration(int64(10000000000)),
				ExecTimeout:  time.Duration(int64(10000000000)),
				TranTimeout:  time.Duration(int64(20000000000)),
			},
			Prom:         &conf.Prom{Limit: 520},
			CollectInfoc: &infoc.Config{},
		}
	}
	svr = New(c)

	if err := svr.dao.Ping(context.Background()); err != nil {
		panic(err)
	}
}
