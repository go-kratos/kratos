package dao

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"go-common/app/service/main/dapper/conf"
	"go-common/app/service/main/dapper/model"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var cfg *conf.Config
var flagMap = map[string]string{
	"app_id":       "main.common-arch.dapper-service",
	"conf_token":   "528dd7e00bb411e894c14a552f48fef8",
	"tree_id":      "5172",
	"conf_version": "server-1",
	"deploy_env":   "uat",
	"conf_host":    "config.bilibili.co",
	"conf_path":    os.TempDir(),
	"region":       "sh",
	"zone":         "sh001",
}

func TestMain(m *testing.M) {
	for key, val := range flagMap {
		flag.Set(key, val)
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Printf("init config from remote error: %s", err)
	}
	cfg = conf.Conf
	if cfg.InfluxDB != nil {
		cfg.InfluxDB.Database = "dapper_ut"
	}
	if cfg.HBase != nil {
		cfg.HBase.Namespace = "dapperut"
	}
	if hbaseAddrs := os.Getenv("TEST_HBASE_ADDRS"); hbaseAddrs != "" {
		cfg.HBase = &conf.HBaseConfig{Addrs: hbaseAddrs, Namespace: "dapperut"}
		if influxdbAddr := os.Getenv("TEST_INFLUXDB_ADDR"); influxdbAddr != "" {
			cfg.InfluxDB = &conf.InfluxDBConfig{Addr: influxdbAddr, Database: "dapper_ut"}
		}
	}
	os.Exit(m.Run())
}

func TestDao(t *testing.T) {
	if cfg == nil {
		t.Skipf("no config provide skipped")
	}
	daoImpl, err := New(cfg)
	if err != nil {
		t.Fatalf("new dao error: %s", err)
	}
	ctx := context.Background()
	Convey("test fetch serviceName and operationName", t, func() {
		serviceNames, err := daoImpl.FetchServiceName(ctx)
		So(err, ShouldBeNil)
		So(serviceNames, ShouldNotBeEmpty)
		for _, serviceName := range serviceNames {
			operationNames, err := daoImpl.FetchOperationName(ctx, serviceName)
			So(err, ShouldBeNil)
			t.Logf("%s operationNames: %v", serviceName, operationNames)
		}
	})
	Convey("test write rawtrace", t, func() {
		if err := daoImpl.WriteRawTrace(
			context.Background(),
			strconv.FormatUint(rand.Uint64(), 16),
			map[string][]byte{strconv.FormatUint(rand.Uint64(), 16): []byte("hello world")},
		); err != nil {
			t.Error(err)
		}
	})
	Convey("test batchwrite span point", t, func() {
		points := []*model.SpanPoint{
			&model.SpanPoint{
				ServiceName:   "service_a",
				OperationName: "opt1",
				PeerService:   "peer_service_a",
				SpanKind:      "client",
				Timestamp:     time.Now().Unix() - rand.Int63n(3600),
				MaxDuration: model.SamplePoint{
					SpanID:  rand.Uint64(),
					TraceID: rand.Uint64(),
					Value:   rand.Int63n(1024),
				},
				MinDuration: model.SamplePoint{
					SpanID:  rand.Uint64(),
					TraceID: rand.Uint64(),
					Value:   rand.Int63n(1024),
				},
				AvgDuration: model.SamplePoint{
					SpanID:  rand.Uint64(),
					TraceID: rand.Uint64(),
					Value:   rand.Int63n(1024),
				},
				Errors: []model.SamplePoint{
					model.SamplePoint{
						SpanID:  rand.Uint64(),
						TraceID: rand.Uint64(),
						Value:   1,
					},
					model.SamplePoint{
						SpanID:  rand.Uint64(),
						TraceID: rand.Uint64(),
						Value:   1,
					},
				},
			},
			&model.SpanPoint{
				ServiceName:   "service_b",
				OperationName: "opt2",
				PeerService:   "peer_service_b",
				SpanKind:      "server",
				Timestamp:     time.Now().Unix() - rand.Int63n(3600),
			},
			&model.SpanPoint{
				ServiceName:   "service_c",
				OperationName: "opt3",
				PeerService:   "peer_service_c",
				SpanKind:      "client",
				Timestamp:     time.Now().Unix() - rand.Int63n(3600),
			},
		}
		err := daoImpl.BatchWriteSpanPoint(context.Background(), points)
		if err != nil {
			t.Error(err)
		}
	})
}
