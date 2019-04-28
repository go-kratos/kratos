package discovery

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/exp/feature"
	"go-common/library/naming"
	"go-common/library/net/http/blademaster"

	. "github.com/smartystreets/goconvey/convey"
)

var appdID = "main.arch.test66"
var appID2 = "main.arch.test22"

func TestMain(m *testing.M) {
	feature.DefaultGate.AddFlag(flag.CommandLine)
	flag.Set("feature-gates", fmt.Sprintf("%s=true", _selfDiscoveryFeatrue))

	os.Exit(m.Run())
}

var c = &Config{
	Nodes:  []string{"172.18.33.51:7171"},
	Zone:   "sh001",
	Env:    "pre",
	Key:    "0c4b8fe3ff35a4b6",
	Secret: "b370880d1aca7d3a289b9b9a7f4d6812",
	Host:   "host_1",
}
var indis = &naming.Instance{
	AppID: appdID,
	Zone:  env.Zone,
	Addrs: []string{
		"grpc://172.18.33.51:8080",
		"http://172.18.33.51:7171",
	},
	Version: "1",
	Metadata: map[string]string{
		"test":   "1",
		"weight": "12",
		"color":  "blue",
	},
}

var in = &naming.Instance{
	AppID: appdID,
	Addrs: []string{
		"grpc://127.0.0.1:8080",
	},
	Version: "1",
	Metadata: map[string]string{
		"test":   "1",
		"weight": "12",
		"color":  "blue",
	},
}

var intest2 = &naming.Instance{
	AppID: appID2,
	Addrs: []string{
		"grpc://127.0.0.1:8080",
	},
	Version: "1",
	Metadata: map[string]string{
		"test":   "1",
		"weight": "12",
		"color":  "blue",
	},
}

var in2 = &naming.Instance{
	AppID: appdID,
	Addrs: []string{
		"grpc://127.0.0.1:8081",
	},
	Version: "1",
	Metadata: map[string]string{
		"test":   "2",
		"weight": "6",
		"color":  "red",
	},
}

func TestDiscoverySelf(t *testing.T) {
	Convey("test TestDiscoverySelf ", t, func() {
		So(feature.DefaultGate.Enabled(_selfDiscoveryFeatrue), ShouldBeTrue)
		d := New(c)
		So(len(d.node.Load().([]string)), ShouldNotEqual, 0)
	})
}

func TestRegister(t *testing.T) {
	Convey("test register and cancel", t, func() {
		env.Hostname = "host_1"
		d := New(c)
		defer d.Close()
		ctx := context.TODO()
		cancel, err := d.Register(ctx, in)
		defer cancel()
		So(err, ShouldBeNil)
		rs := d.Build(appdID)
		ch := rs.Watch()
		So(ch, ShouldNotBeNil)
		<-ch
		ins, ok := rs.Fetch(ctx)
		So(ok, ShouldBeTrue)
		So(err, ShouldBeNil)
		var count int
		for _, data := range ins {
			count += len(data)
		}
		So(count, ShouldEqual, 1)
		c.Host = "host_2"
		env.Hostname = "host_2"

		d2 := New(c)
		defer d2.Close()
		cancel2, err := d2.Register(ctx, in2)
		So(err, ShouldBeNil)
		<-ch
		ins, _ = rs.Fetch(ctx)
		So(err, ShouldBeNil)
		count = 0
		for _, data := range ins {
			count += len(data)
		}
		So(count, ShouldEqual, 2)
		time.Sleep(time.Millisecond * 500)
		cancel2()
		<-ch
		ins, _ = rs.Fetch(ctx)
		So(err, ShouldBeNil)
		count = 0
		for _, data := range ins {
			count += len(data)
		}
		So(count, ShouldEqual, 1)
		Convey("test discovery set", func() {
			c.Host = "host_1"
			inSet := &naming.Instance{
				AppID: appdID,
				Addrs: []string{
					"grpc://127.0.0.1:8080",
				},
				Status: 1,
				Metadata: map[string]string{
					"test":   "1",
					"weight": "111",
					"color":  "blue",
				},
			}
			ins, _ = rs.Fetch(context.TODO())
			fmt.Println("ins", ins["sh001"][0])
			err = d2.Set(inSet)
			//	So(err, ShouldBeNil)
			<-ch
			ins, _ = rs.Fetch(context.TODO())
			fmt.Println("ins1", ins["sh001"][0])
			So(ins["sh001"][0].Metadata["weight"], ShouldResemble, "111")
		})
	})
}

func TestMultiZone(t *testing.T) {
	Convey("test multi zone", t, func() {
		env.Hostname = "host_1"
		d := New(c)
		defer d.Close()
		ctx := context.TODO()
		cancel, err := d.Register(ctx, in)
		So(err, ShouldBeNil)
		defer cancel()
		rs := d.Build(appdID)
		ch := rs.Watch()
		So(ch, ShouldNotBeNil)
		<-ch
		ins, ok := rs.Fetch(ctx)
		So(ok, ShouldBeTrue)
		count := 0
		for _, data := range ins {
			count += len(data)
		}
		So(count, ShouldEqual, 1)
		env.Hostname = "host_2"
		env.Zone = "other_zone"
		d2 := New(c)
		defer d2.Close()
		cancel2, err := d2.Register(ctx, in2)
		So(err, ShouldBeNil)
		defer func() {
			cancel2()
			time.Sleep(time.Millisecond * 300)
		}()
		<-ch
		ins, ok = rs.Fetch(ctx)
		So(ok, ShouldBeTrue)
		count = 0
		zoneCount := 0
		for _, data := range ins {
			zoneCount++
			count += len(data)
		}
		So(count, ShouldEqual, 2)
		So(zoneCount, ShouldEqual, 2)
	})
}

func TestDiscoveryFailOver(t *testing.T) {
	Convey("test failover", t, func() {
		var conf = &Config{
			Nodes:  []string{"127.0.0.1:8080"},
			Zone:   "sh001",
			Env:    "pre",
			Key:    "0c4b8fe3ff35a4b6",
			Secret: "b370880d1aca7d3a289b9b9a7f4d6812",
			Host:   "host_1",
		}
		for _, a := range []string{":8080"} {
			go func(addr string) {
				e := blademaster.DefaultServer(nil)
				e.GET("/discovery/nodes", func(ctx *blademaster.Context) {
					type v struct {
						Addr string `json:"addr"`
					}
					ctx.JSON([]v{{Addr: "127.0.0.1:8080"}}, nil)
				})
				e.GET("/discovery/polls", func(ctx *blademaster.Context) {
					params := ctx.Request.Form
					ts := params.Get("latest_timestamp")
					if ts == "0" {
						ctx.JSON(map[string]appData{
							appdID:            {LastTs: time.Now().UnixNano(), ZoneInstances: map[string][]*naming.Instance{"zone": {in}}},
							"infra.discovery": {LastTs: time.Now().UnixNano(), ZoneInstances: map[string][]*naming.Instance{conf.Zone: {indis}}},
						}, nil)
					} else {
						ctx.JSON(nil, ecode.ServerErr)
					}
				})
				e.Run(addr)
			}(a)
		}
		time.Sleep(time.Millisecond * 30)
		d := New(conf)
		defer d.Close()
		rs := d.Build(appdID)
		ch := rs.Watch()
		<-ch
		ins, _ := rs.Fetch(context.TODO())
		count := 0
		zoneCount := 0
		for _, data := range ins {
			zoneCount++
			count += len(data)
		}
		So(count, ShouldEqual, 1)
		So(zoneCount, ShouldEqual, 1)
	})
}

func TestWatchContinuosly(t *testing.T) {
	Convey("test TestWatchContinuosly ", t, func() {
		env.Hostname = "host_1"
		d := New(c)
		defer d.Close()
		in1 := *in
		in1.AppID = "test.test"
		ctx := context.TODO()
		cancel, err := d.Register(ctx, &in1)
		So(err, ShouldBeNil)
		defer cancel()
		in2 := *in
		in2.AppID = "test.test2"
		cancel2, err := d.Register(ctx, &in2)
		So(err, ShouldBeNil)
		defer cancel2()
		in3 := *in
		in3.AppID = "test.test3"
		cancel3, err := d.Register(ctx, &in3)
		So(err, ShouldBeNil)
		defer cancel3()
		rs := d.Build("test.test")
		ch := rs.Watch()
		<-ch
		ins, ok := rs.Fetch(ctx)
		So(ok, ShouldBeTrue)
		count := 0
		for _, data := range ins {
			count += len(data)
		}
		So(count, ShouldBeGreaterThanOrEqualTo, 1)
		time.Sleep(time.Millisecond * 10)
		rs = d.Build("test.test2")
		ch2 := rs.Watch()
		<-ch2
		ins, ok = rs.Fetch(ctx)
		So(ok, ShouldBeTrue)
		count = 0
		for _, data := range ins {
			count += len(data)
		}
		So(count, ShouldBeGreaterThanOrEqualTo, 1)
		rs = d.Build("test.test3")
		ch3 := rs.Watch()
		<-ch3
		ins, ok = rs.Fetch(ctx)
		So(ok, ShouldBeTrue)
		count = 0
		for _, data := range ins {
			count += len(data)
		}
		So(count, ShouldBeGreaterThanOrEqualTo, 1)
	})
}

func TestSameBuilder(t *testing.T) {
	Convey("test multi watch", t, func() {
		env.Hostname = "host_1"
		d := New(c)
		defer d.Close()
		ctx := context.TODO()
		cancel, err := d.Register(ctx, in)
		d.Register(ctx, intest2)
		defer cancel()
		So(err, ShouldBeNil)
		// first builder
		rs := d.Build(appdID)
		ch := rs.Watch()
		So(ch, ShouldNotBeNil)
		<-ch
		_, ok := rs.Fetch(ctx)
		So(ok, ShouldBeTrue)
		So(err, ShouldBeNil)
		var count int
		// for _, data := range ins {
		// 	count += len(data)
		// }
		// So(count, ShouldEqual, 1)
		// same appd builder
		rs2 := d.Build(appdID)
		ch2 := rs2.Watch()
		<-ch2
		_, ok = rs2.Fetch(ctx)
		So(ok, ShouldBeTrue)
		So(err, ShouldBeNil)
		rs3 := d.Build(appID2)
		ch3 := rs3.Watch()
		<-ch3
		ins, _ := rs3.Fetch(ctx)
		So(ok, ShouldBeTrue)
		So(err, ShouldBeNil)
		count = 0
		for _, data := range ins {
			count += len(data)
		}
		So(count, ShouldEqual, 1)
	})
}
