package dao

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/main/dapper-query/conf"
	"golang.org/x/sys/unix"
)

var cfg *conf.Config
var flagMap = map[string]string{
	"app_id":       "main.common-arch.dapper-query",
	"conf_appid":   "main.common-arch.dapper-query",
	"conf_token":   "ed3241c850735df94d24d7b49f69ddd7",
	"tree_id":      "60617",
	"conf_version": "docker-1",
	"deploy_env":   "uat",
	"conf_env":     "uat",
	"conf_host":    "config.bilibili.co",
	"conf_path":    os.TempDir(),
	"region":       "sh",
	"zone":         "sh001",
}

// only for ut runner
func hackHosts() {
	hostsPath := "/etc/hosts"
	if unix.Access(hostsPath, unix.W_OK) != nil {
		return
	}
	fp, err := os.OpenFile(hostsPath, os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("open hosts file error: %s", err)
	}
	defer fp.Close()
	fmt.Fprintf(fp, "\n")
	fmt.Fprintln(fp, "172.22.33.146   nvm-test-dapper-influxdb-01")
}

func TestMain(m *testing.M) {
	hackHosts()
	for key, val := range flagMap {
		flag.Set(key, val)
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Printf("init config from remote error: %s", err)
	}
	if hbaseAddrs := os.Getenv("TEST_HBASE_ADDRS"); hbaseAddrs != "" {
		cfg = new(conf.Config)
		cfg.HBase = &conf.HBaseConfig{Addrs: hbaseAddrs, Namespace: "ugc"}
		if influxdbAddr := os.Getenv("TEST_INFLUXDB_ADDR"); influxdbAddr != "" {
			cfg.InfluxDB = &conf.InfluxDBConfig{Addr: influxdbAddr, Database: "dapper_uat"}
		}
	}
	if cfg == nil {
		cfg = conf.Conf
		if cfg.InfluxDB != nil {
			cfg.InfluxDB.Database = "dapper_uat"
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
	serviceName := "main.community.tag"
	operationName := "/x/internal/tag/archive/tags"
	Convey("query serviceNames", t, func() {
		serviceNames, err := daoImpl.ServiceNames(ctx)
		So(err, ShouldBeNil)
		So(serviceNames, ShouldNotBeEmpty)
		t.Logf("serviceNames: %v", serviceNames)
		Convey("query operationNames", func() {
			// FIXME: make mock data frist
			operationNames, err := daoImpl.OperationNames(ctx, serviceName)
			So(err, ShouldBeNil)
			So(operationNames, ShouldNotBeEmpty)
			t.Logf("operationNames for %s :%v", serviceName, operationNames)
		})
	})
	Convey("test QuerySpanListTime Asc", t, func() {
		// FIXME: make mock data frist
		spanListRefs, err := daoImpl.QuerySpanList(ctx, serviceName, operationName, &Selector{
			Start:  time.Now().Unix() - 3600,
			End:    time.Now().Unix(),
			Limit:  10,
			Offset: 10,
		}, TimeAsc)
		So(err, ShouldBeNil)
		So(spanListRefs, ShouldNotBeEmpty)
		t.Logf("spanListRefs: %v", spanListRefs)
	})
	Convey("test QuerySpanListTime Desc", t, func() {
		// FIXME: make mock data frist
		spanListRefs, err := daoImpl.QuerySpanList(ctx, serviceName, operationName, &Selector{
			Start:  time.Now().Unix() - 3600*12,
			End:    time.Now().Unix(),
			Limit:  10,
			Offset: 10,
		}, TimeDesc)
		So(err, ShouldBeNil)
		So(spanListRefs, ShouldNotBeEmpty)
		t.Logf("spanListRefs: %v", spanListRefs)
		Convey("test get trace", func() {
			spanListRef := spanListRefs[0]
			spans, err := daoImpl.Trace(ctx, spanListRef.TraceID)
			So(err, ShouldBeNil)
			So(spans, ShouldNotBeEmpty)
			t.Logf("spans %v", spans)
		})
	})
	Convey("test QuerySpanListDuration Desc", t, func() {
		// FIXME: make mock data frist
		spanListRefs, err := daoImpl.QuerySpanList(ctx, serviceName, operationName, &Selector{
			Start:  time.Now().Unix() - 3600*12,
			End:    time.Now().Unix(),
			Limit:  10,
			Offset: 10,
		}, DurationDesc)
		So(err, ShouldBeNil)
		So(spanListRefs, ShouldNotBeEmpty)
		t.Logf("spanListRefs: %v", spanListRefs)
		Convey("test get trace", func() {
			spanListRef := spanListRefs[len(spanListRefs)-1]
			spans, err := daoImpl.Trace(ctx, spanListRef.TraceID)
			So(err, ShouldBeNil)
			So(spans, ShouldNotBeEmpty)
			t.Logf("spans %v", spans)
		})
	})
	Convey("test QuerySpanListDuration Asc", t, func() {
		// FIXME: make mock data frist
		spanListRefs, err := daoImpl.QuerySpanList(ctx, serviceName, operationName, &Selector{
			Start:  time.Now().Unix() - 3600*12,
			End:    time.Now().Unix(),
			Limit:  10,
			Offset: 10,
		}, DurationAsc)
		So(err, ShouldBeNil)
		So(spanListRefs, ShouldNotBeEmpty)
		t.Logf("spanListRefs: %v", spanListRefs)
		Convey("test get trace", func() {
			spanListRef := spanListRefs[len(spanListRefs)-1]
			spans, err := daoImpl.Trace(ctx, spanListRef.TraceID)
			So(err, ShouldBeNil)
			So(spans, ShouldNotBeEmpty)
			t.Logf("spans %v", spans)
		})
	})
	Convey("test MeanOperationNameField", t, func() {
		start := time.Now().Unix() - 3600
		end := time.Now().Unix()
		values, err := daoImpl.MeanOperationNameField(ctx, map[string]string{"service_name": serviceName}, "max_duration", start, end, []string{"operation_name"})
		if err != nil {
			t.Error(err)
		}
		So(values, ShouldNotBeEmpty)
	})
	Convey("test SpanSeriesMean", t, func() {
		start := time.Now().Unix() - 3600
		end := time.Now().Unix()
		series, err := daoImpl.SpanSeriesMean(ctx, serviceName, operationName, []string{"max_duration", "min_duration"}, start, end, 30)
		if err != nil {
			t.Error(err)
		}
		So(series.Timestamps, ShouldNotBeEmpty)
		So(series.Items, ShouldNotBeEmpty)
		// FIXME
		//for _, item := range series.Items {
		//	So(len(series.Timestamps), ShouldEqual, len(item.Rows))
		//}
		t.Logf("%#v\n", series)
	})
	Convey("test SpanSeriesCount", t, func() {
		start := time.Now().Unix() - 3600
		end := time.Now().Unix()
		series, err := daoImpl.SpanSeriesCount(ctx, serviceName, operationName, []string{"max_duration", "min_duration"}, start, end, 30)
		if err != nil {
			t.Error(err)
		}
		So(series.Timestamps, ShouldNotBeEmpty)
		So(series.Items, ShouldNotBeEmpty)
		// FIXME
		//for _, item := range series.Items {
		//	So(len(series.Timestamps), ShouldEqual, len(item.Rows))
		//}
		t.Logf("%#v\n", series)
	})
	Convey("test PeerService", t, func() {
		serviceName := "main.bangumi.season-service"
		peerServices, err := daoImpl.PeerService(ctx, serviceName)
		if err != nil {
			t.Error(err)
		}
		So(peerServices, ShouldNotBeEmpty)
	})
	Convey("test trace", t, func() {
		traceID, _ := strconv.ParseUint("100056da7886666c", 16, 64)
		spans, err := daoImpl.Trace(ctx, traceID)
		if err != nil {
			t.Error(err)
		}
		So(spans, ShouldNotBeEmpty)
	})
	Convey("test ping close", t, func() {
		So(daoImpl.Ping(ctx), ShouldBeNil)
		So(daoImpl.Close(ctx), ShouldBeNil)
	})
}
