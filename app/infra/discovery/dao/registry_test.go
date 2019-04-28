package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"go-common/app/infra/discovery/model"
	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

var reg = &model.ArgRegister{Appid: "main.arch.test", Hostname: "reg", RPC: "127.0.0.1:8080", Region: "shsb", Zone: "sh0001", Env: "pre", Status: 1}
var regH1 = &model.ArgRegister{Appid: "main.arch.test", Hostname: "regH1", RPC: "127.0.0.1:8080", Region: "shsb", Zone: "sh0001", Env: "pre", Status: 1}

var reg2 = &model.ArgRegister{Appid: "main.arch.test2", Hostname: "reg2", RPC: "127.0.0.1:8080", Region: "shsb", Zone: "sh0001", Env: "pre", Status: 1}

var arg = &model.ArgRenew{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: "main.arch.test", Hostname: "reg"}
var cancel = &model.ArgCancel{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: "main.arch.test", Hostname: "reg"}
var cancel2 = &model.ArgCancel{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: "main.arch.test", Hostname: "regH1"}

func TestReigster(t *testing.T) {
	i := model.NewInstance(reg)
	register(t, i)
}

func TestDiscovery(t *testing.T) {
	i1 := model.NewInstance(reg)
	i2 := model.NewInstance(regH1)
	fmt.Println(_evictThreshold)
	r := register(t, i1, i2)
	Convey("test discovery", t, func() {
		pollArg := &model.ArgPolls{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: []string{"main.arch.test"}, Hostname: "test"}
		fetchArg := &model.ArgFetch{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: "main.arch.test", Status: 3}
		info, err := r.Fetch(fetchArg.Zone, fetchArg.Env, fetchArg.Appid, 0, fetchArg.Status)
		So(err, ShouldBeNil)
		So(len(info.Instances), ShouldEqual, 2)
		ch, _, err := r.Polls(pollArg)
		So(err, ShouldBeNil)
		apps := <-ch
		So(len(apps["main.arch.test"].Instances), ShouldEqual, 2)
		pollArg.LatestTimestamp[0] = apps["main.arch.test"].LatestTimestamp
		fmt.Println(apps["main.arch.test"])
		r.Cancel(cancel)
		ch, _, err = r.Polls(pollArg)
		So(err, ShouldBeNil)
		apps = <-ch
		So(len(apps["main.arch.test"].Instances), ShouldEqual, 1)
		pollArg.LatestTimestamp[0] = apps["main.arch.test"].LatestTimestamp
		r.Cancel(cancel2)
	})
}

func TestRenew(t *testing.T) {
	src := model.NewInstance(reg)
	r := register(t, src)
	Convey("test renew", t, func() {
		i, ok := r.Renew(arg)
		So(ok, ShouldBeTrue)
		So(i, ShouldResemble, src)
	})
}

func BenchmarkRenew(b *testing.B) {
	var (
		i  *model.Instance
		ok bool
	)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r, src := benchRegister(b)
			if i, ok = r.Renew(arg); !ok {
				b.Errorf("Renew(%v)", src.Appid)
			}
			benchCompareInstance(b, src, i)
		}
	})
}

func TestCancel(t *testing.T) {
	src := model.NewInstance(reg)
	r := register(t, src)
	Convey("test cancel", t, func() {
		i, ok := r.Cancel(cancel)
		So(ok, ShouldBeTrue)
		So(i, ShouldResemble, src)
		fetchArg := &model.ArgFetch{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: "main.arch.test", Status: 3}
		_, err := r.Fetch(fetchArg.Zone, fetchArg.Env, fetchArg.Appid, 0, fetchArg.Status)
		So(err, ShouldResemble, ecode.NothingFound)
	})
}

func BenchmarkCancel(b *testing.B) {
	var (
		i   *model.Instance
		ok  bool
		err error
	)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r, src := benchRegister(b)
			if i, ok = r.Cancel(cancel); !ok {
				b.Errorf("Cancel(%v) error", src.Appid)
			}
			benchCompareInstance(b, src, i)
			fetchArg := &model.ArgFetch{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: "main.arch.test", Status: 3}
			if _, err = r.Fetch(fetchArg.Zone, fetchArg.Env, fetchArg.Appid, 0, fetchArg.Status); err != ecode.NothingFound {
				b.Errorf("Fetch(%v) error(%v)", src.Appid, err)
			}
		}
	})
}

func TestFetchAll(t *testing.T) {
	i := model.NewInstance(reg)
	r := register(t, i)
	Convey("test fetch all", t, func() {
		am := r.FetchAll()
		So(len(am), ShouldResemble, 1)
	})
}

func BenchmarkFetchAll(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r, _ := benchRegister(b)
			if am := r.FetchAll(); len(am) != 1 {
				b.Errorf("FetchAll() error")
			}
		}
	})
}

func TestFetch(t *testing.T) {
	i := model.NewInstance(reg)
	r := register(t, i)
	Convey("test fetch", t, func() {
		fetchArg2 := &model.ArgFetch{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: "main.arch.test", Status: 1}
		c, err := r.Fetch(fetchArg2.Zone, fetchArg2.Env, fetchArg2.Appid, 0, fetchArg2.Status)
		So(err, ShouldBeNil)
		So(len(c.Instances), ShouldResemble, 1)
	})
}

func BenchmarkFetch(b *testing.B) {
	var (
		err error
		c   *model.InstanceInfo
	)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r, _ := benchRegister(b)
			fetchArg := &model.ArgFetch{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: "main.arch.test", Status: 1}
			if c, err = r.Fetch(fetchArg.Zone, fetchArg.Env, fetchArg.Appid, 0, fetchArg.Status); err != nil {
				b.Errorf("Fetch(%v) error(%v)", arg.Appid, err)
			}
			fetchArg2 := &model.ArgFetch{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: "main.arch.test", Status: 2}
			if c, err = r.Fetch(fetchArg2.Zone, fetchArg2.Env, fetchArg2.Appid, 0, fetchArg2.Status); err != nil {
				b.Errorf("Fetch(%v) error(%v)", arg.Appid, err)
			}
			_ = c
		}
	})
}

func TestPoll(t *testing.T) {
	i := model.NewInstance(reg)
	r := register(t, i)
	Convey("test poll", t, func() {
		pollArg := &model.ArgPolls{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: []string{"main.arch.test"}, Hostname: "csq"}
		ch, _, err := r.Polls(pollArg)
		So(err, ShouldBeNil)
		c := <-ch
		So(len(c[pollArg.Appid[0]].Instances), ShouldEqual, 1)
	})
}

func TestPolls(t *testing.T) {
	i1 := model.NewInstance(reg)
	i2 := model.NewInstance(reg2)
	r := register(t, i1, i2)
	Convey("test polls", t, func() {
		pollArg := &model.ArgPolls{Region: "shsb", Zone: "sh0001", Env: "pre", LatestTimestamp: []int64{0, 0}, Appid: []string{"main.arch.test", "main.arch.test2"}, Hostname: "csq"}
		ch, new, err := r.Polls(pollArg)
		So(err, ShouldBeNil)
		So(new, ShouldBeTrue)
		c := <-ch
		So(len(c), ShouldResemble, 2)
	})
}

func TestPollsParallel(t *testing.T) {

	i1 := model.NewInstance(reg)
	i2 := model.NewInstance(reg2)
	r := register(t, i1, i2)

	Convey("test polls parallel", t, func(c C) {
		var (
			wg       sync.WaitGroup
			ch1, ch2 chan map[string]*model.InstanceInfo
			new      bool
			err      error
		)
		pollArg := &model.ArgPolls{Region: "shsb", Zone: "sh0001", Env: "pre", LatestTimestamp: []int64{time.Now().UnixNano(), time.Now().UnixNano()}, Appid: []string{"main.arch.test", "main.arch.test2"}, Hostname: "csq"}
		ch1, new, err = r.Polls(pollArg)
		c.So(err, ShouldEqual, ecode.NotModified)
		c.So(new, ShouldBeFalse)
		c.So(ch1, ShouldNotBeNil)
		ch2, new, err = r.Polls(pollArg)
		c.So(err, ShouldEqual, ecode.NotModified)
		c.So(new, ShouldBeFalse)
		c.So(ch2, ShouldNotBeNil)
		// wait group
		wg.Add(2)
		go func() {
			res := <-ch1
			c.So(len(res), ShouldResemble, 1)
			ress, _ := json.Marshal(res)
			fmt.Println("chenggongle 1!!!", string(ress))
			wg.Done()
		}()
		go func() {
			res := <-ch2
			c.So(len(res), ShouldResemble, 1)
			ress, _ := json.Marshal(res)
			fmt.Println("chenggongle 2!!!", string(ress))
			wg.Done()
		}()
		// re register when 1s later, make sure latest_timestamp changed
		time.Sleep(time.Second)
		h1 := model.NewInstance(regH1)
		r.Register(h1, 0)
		// wait
		wg.Wait()
	})
}

func BenchmarkPoll(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var (
				err error
				ch  chan map[string]*model.InstanceInfo
				c   map[string]*model.InstanceInfo
			)
			r, _ := benchRegister(b)
			pollArg := &model.ArgPolls{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: []string{"main.arch.test"}, Hostname: "csq"}
			if ch, _, err = r.Polls(pollArg); err != nil {
				b.Errorf("Poll(%v) error(%v)", arg.Appid, err)
			}
			if c = <-ch; len(c[pollArg.Appid[0]].Instances) != 1 {
				b.Errorf("Poll(%v) lenth error", arg.Appid)
			}
		}
	})
}

func TestBroadcast(t *testing.T) {
	i := model.NewInstance(reg)
	r := register(t, i)
	Convey("test poll push connection", t, func() {
		go func() {
			Convey("must poll ahead of time", t, func() {
				time.Sleep(time.Microsecond * 5)
				var arg2 = &model.ArgRegister{Appid: "main.arch.test", Hostname: "go", RPC: "127.0.0.1:8080", Region: "shsb", Zone: "sh0001", Env: "pre", Status: 1}
				m2 := model.NewInstance(arg2)
				err2 := r.Register(m2, 0)
				So(err2, ShouldBeNil)
			})
		}()
		pollArg := &model.ArgPolls{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: []string{"main.arch.test"}, LatestTimestamp: []int64{time.Now().UnixNano()}}
		ch, _, err := r.Polls(pollArg)
		So(err, ShouldResemble, ecode.NotModified)
		c := <-ch
		So(len(c[pollArg.Appid[0]].Instances), ShouldResemble, 2)
		So(c[pollArg.Appid[0]].ZoneInstances, ShouldNotBeNil)
		So(len(c[pollArg.Appid[0]].ZoneInstances["sh0001"]), ShouldResemble, 2)
	})
}

func BenchmarkBroadcast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var (
			err  error
			err2 error
			ch   chan map[string]*model.InstanceInfo
			c    map[string]*model.InstanceInfo
		)
		r, _ := benchRegister(b)
		go func() {
			time.Sleep(time.Millisecond * 1)
			var arg2 = &model.ArgRegister{Appid: "main.arch.test", Hostname: "go", RPC: "127.0.0.1:8080", Region: "shsb", Zone: "sh0001", Env: "pre", Status: 1}
			m2 := model.NewInstance(arg2)
			if err2 = r.Register(m2, 0); err2 != nil {
				b.Errorf("Reigster(%v) error(%v)", m2.Appid, err2)
			}
		}()
		pollArg := &model.ArgPolls{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: []string{"main.arch.test"}, LatestTimestamp: []int64{time.Now().UnixNano()}}
		if ch, _, err = r.Polls(pollArg); err != nil && err != ecode.NotModified {
			b.Errorf("Poll(%v) error(%v)", pollArg.Appid, err)
		}
		c = <-ch
		if len(c[pollArg.Appid[0]].Instances) != 2 {
			b.Errorf("Poll(%v) length error", pollArg.Appid)
		}
		if c[pollArg.Appid[0]].ZoneInstances == nil {
			b.Errorf("Poll(%v) zone instances nil error", pollArg.Appid)
		}
		if len(c[pollArg.Appid[0]].ZoneInstances["sh0001"]) != 2 {
			b.Errorf("Poll(%v) zone instances length error", pollArg.Appid)
		}
	}
}

func TestRegistrySet(t *testing.T) {
	i := model.NewInstance(reg)
	r := register(t, i)
	changes := make(map[string]string)
	changes["reg"] = "1"
	Convey("test set weight to 1", t, func() {
		set := &model.ArgSet{
			Region:   "shsb",
			Env:      "pre",
			Appid:    "main.arch.test",
			Hostname: []string{"reg"},
			Metadata: []string{`{"weight":"1"}`},
		}
		ok := r.Set(context.TODO(), set)
		So(ok, ShouldBeTrue)
		fetchArg := &model.ArgFetch{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: "main.arch.test", Status: 3}
		c, err := r.Fetch(fetchArg.Zone, fetchArg.Env, fetchArg.Appid, 0, fetchArg.Status)
		So(err, ShouldBeNil)
		So(c.Instances[0].Metadata["weight"], ShouldResemble, "1")
	})
}

func BenchmarkSet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var (
				c   *model.InstanceInfo
				err error
				ok  bool
			)
			r, _ := benchRegister(b)
			set := &model.ArgSet{
				Region:   "shsb",
				Env:      "pre",
				Appid:    "main.arch.account-service",
				Hostname: []string{"test1"},
				Status:   []int64{1},
				Metadata: []string{`{"weight":"1"}`},
			}
			if ok = r.Set(context.TODO(), set); !ok {
				b.Errorf("SetWeight(%v) error", arg.Appid)
			}
			fetchArg := &model.ArgFetch{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: "main.arch.test", Status: 3}
			if c, err = r.Fetch(fetchArg.Zone, fetchArg.Env, fetchArg.Appid, 0, fetchArg.Status); err != nil {
				b.Errorf("Fetch(%v) error(%v)", fetchArg.Appid, err)
			}
			if c.Instances[0].Metadata["weight"] != "1" {
				b.Errorf("SetWeight(%v) change error", fetchArg.Appid)
			}
		}
	})
}

func TestResetExp(t *testing.T) {
	i := model.NewInstance(reg)
	r := register(t, i)
	Convey("test ResetExp", t, func() {
		r.resetExp()
		So(r.gd.expPerMin, ShouldResemble, int64(2))
	})
}

func benchCompareInstance(b *testing.B, src *model.Instance, i *model.Instance) {
	if src.Appid != i.Appid || src.Env != i.Env || src.Hostname != i.Hostname ||
		src.Region != i.Region {
		b.Errorf("instance compare error")
	}
}

func register(t *testing.T, is ...*model.Instance) (r *Registry) {
	Convey("test register", t, func() {
		r = NewRegistry()
		var num int
		for _, i := range is {
			err := r.Register(i, 0)
			So(err, ShouldBeNil)
			if i.Appid == "main.arch.test" {
				num++
			}
		}
		fetchArg := &model.ArgFetch{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: "main.arch.test", Status: 3}
		instancesInfo, err := r.Fetch(fetchArg.Zone, fetchArg.Env, fetchArg.Appid, 0, fetchArg.Status)
		So(err, ShouldBeNil)
		So(len(instancesInfo.Instances), ShouldResemble, num)
	})
	return r
}

func benchRegister(b *testing.B) (r *Registry, i *model.Instance) {
	r = NewRegistry()
	i = model.NewInstance(reg)
	if err := r.Register(i, 0); err != nil {
		b.Errorf("Reigster(%v) error(%v)", i.Appid, err)
	}
	return r, i
}

func TestEvict(t *testing.T) {
	Convey("test evict for protect", t, func() {
		r := NewRegistry()
		m := model.NewInstance(reg)
		// promise the renewtime of instance is expire
		m.RenewTimestamp -= 100
		err := r.Register(m, 0)
		So(err, ShouldBeNil)
		// move up the statistics of heartbeat for evict
		r.gd.facLastMin = r.gd.facInMin
		r.evict()
		fetchArg := &model.ArgFetch{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: "main.arch.test", Status: 3}
		c, err := r.Fetch(fetchArg.Zone, fetchArg.Env, fetchArg.Appid, 0, fetchArg.Status)
		So(err, ShouldBeNil)
		// protect
		So(len(c.Instances), ShouldResemble, 1)
	})
}

func TestEvict2(t *testing.T) {
	Convey("test evict for cancel", t, func() {
		r := NewRegistry()
		m := model.NewInstance(reg)
		err := r.Register(m, 0)
		So(err, ShouldBeNil)
		_, ok := r.Renew(arg)
		So(ok, ShouldBeTrue)
		// promise the renewtime of instance is expire
		m.RenewTimestamp -= int64(time.Second * 100)
		r.Register(m, 0)
		// move up the statistics of heartbeat for evict
		r.gd.facLastMin = r.gd.facInMin
		r.evict()
		fetchArg := &model.ArgFetch{Region: "shsb", Zone: "sh0001", Env: "pre", Appid: "main.arch.test", Status: 1}
		_, err = r.Fetch(fetchArg.Zone, fetchArg.Env, fetchArg.Appid, 0, fetchArg.Status)
		So(err, ShouldResemble, ecode.NothingFound)
	})
}
