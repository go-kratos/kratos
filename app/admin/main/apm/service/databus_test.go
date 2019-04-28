package service

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"go-common/app/admin/main/apm/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	svr *Service
)

func TestMain(m *testing.M) {
	var (
		err error
	)
	dir, _ := filepath.Abs("../cmd/apm-admin-test.toml")
	if err = flag.Set("conf", dir); err != nil {
		panic(err)
	}
	if err = conf.Init(); err != nil {
		panic(err)
	}
	svr = New(conf.Conf)
	os.Exit(m.Run())
}
func TestFake(t *testing.T) {
	Convey("fake", t, func() {
		t.Log("fake test")
	})
}

func TestService_NewClient(t *testing.T) {
	Convey("should new client all", t, func() {
		c, err := NewClient(conf.Conf.Kafka["test_kafka_9092-266"].Brokers, "Archive-T", "Archive-Live-S")
		t.Log(err, c)
		So(err, ShouldBeNil)
	})
}

func TestService_OffsetNew(t *testing.T) {
	Convey("should offset new all", t, func() {
		c, err := NewClient(conf.Conf.Kafka["test_kafka_9092-266"].Brokers, "Archive-T", "Archive-Live-S")
		So(err, ShouldBeNil)
		info, err := c.OffsetNew()
		t.Log(err, info)
		So(err, ShouldBeNil)
	})
}

func TestService_OffsetOld(t *testing.T) {
	Convey("should offset old all", t, func() {
		c, err := NewClient(conf.Conf.Kafka["test_kafka_9092-266"].Brokers, "Archive-T", "Archive-Live-S")
		So(err, ShouldBeNil)
		info, err := c.OffsetOld()
		t.Log(err, info)
		So(err, ShouldBeNil)
	})
}

func TestService_SeekBegin(t *testing.T) {
	Convey("should seek begin all", t, func() {
		c, err := NewClient(conf.Conf.Kafka["test_kafka_9092-266"].Brokers, "Archive-T", "Archive-Live-S")
		c.SeekBegin()
		t.Log(err)
		So(err, ShouldBeNil)
	})
}

func TestService_SeekEnd(t *testing.T) {
	Convey("should seek end all", t, func() {
		c, err := NewClient(conf.Conf.Kafka["test_kafka_9092-266"].Brokers, "Archive-T", "Archive-Live-S")
		So(err, ShouldBeNil)
		err = c.SeekEnd()
		t.Log(err)
		So(err, ShouldBeNil)
	})
}

func TestCreateTopic(t *testing.T) {
	Convey("test create topic", t, func() {
		err := CreateTopic([]string{"172.18.33.51:9098", "172.18.33.52:9098", "172.18.33.50:9098"}, "testcreate11", 1, 1)
		So(err, ShouldBeNil)
	})
}

func TestService_OffsetMarked(t *testing.T) {
	Convey("should offset marked all", t, func() {
		c, err := NewClient(conf.Conf.Kafka["test_kafka_9092-266"].Brokers, "Archive-T", "Archive-Live-S")
		So(err, ShouldBeNil)
		_, err = c.OffsetMarked()
		t.Log(err)
		So(err, ShouldBeNil)
	})
}

func TestService_MsgFetch(t *testing.T) {
	Convey("should msg fetch", t, func() {
		res, err := FetchMessage(context.Background(), "test_kafka_9092-266", "Archive-T", "ArchiveAPM-MainCommonArch-S", "", 0, 0, 10)
		So(err, ShouldBeNil)
		for _, r := range res {
			t.Logf("fetch key:%s value:%s partition:%d offset:%d timestamp:%d", r.Key, r.Value, r.Partition, r.Offset, r.Timestamp)
		}
	})
}
