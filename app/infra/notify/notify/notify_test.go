package notify

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"testing"

	"go-common/app/infra/notify/conf"
	"go-common/app/infra/notify/dao"
	"go-common/app/infra/notify/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *dao.Dao
)

func TestMain(m *testing.M) {
	var err error
	flag.Set("conf", "../cmd/notify-test.toml")
	if err = conf.Init(); err != nil {
		log.Println(err)
		return
	}
	d = dao.New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func TestNotify(t *testing.T) {
	var (
		nt  *Sub
		err error
		pub *Pub
	)
	Convey("test notify", t, func() {
		pub, err = NewPub(&model.Pub{
			Cluster: "test",
			Group:   "pub",
			Topic:   "test1",
		}, conf.Conf)
		So(err, ShouldBeNil)
		err = pub.Send([]byte("test"), []byte(`{"test":"123"}`))
		So(err, ShouldBeNil)
		nt, err = NewSub(&model.Watcher{
			Cluster:  "test",
			Group:    "test",
			Topic:    "test1",
			Callback: string([]byte(`{"http://127.0.0.1:18888/push1": 1}`)),
		}, d, conf.Conf)
		So(err, ShouldBeNil)
		So(nt, ShouldNotBeNil)
		err = nt.dial()
		So(err, ShouldBeNil)
		//	go nt.serve()
		//fmt.Println(nt.consumer.))
	})
	nt, _ = NewSub(&model.Watcher{
		Cluster: "test",
		Group:   "test",
		Topic:   "test1",
	}, d, conf.Conf)
	b, err := json.Marshal(map[string]string{
		"http://127.0.0.1:18888/push1": "1",
		"http://127.0.0.1:18888/push2": "2",
	})
	nt.w.Callback = string(b)
	mockhttp()
	Convey("test push", t, func() {
		nt.push([]byte("push1"))
	})
}

func mockhttp() {
	http.HandleFunc("/push1", testpush1)
	http.HandleFunc("/push2", testpush2)
	go http.ListenAndServe(":18888", nil)
}

func testpush1(resp http.ResponseWriter, req *http.Request) {
}

func testpush2(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(500)
}

func TestSub(t *testing.T) {
	nt, _ := NewSub(&model.Watcher{
		Cluster: "test",
		Group:   "test",
		Topic:   "test1",
		Filters: []*model.Filter{
			{
				Field:     "table",
				Condition: model.ConditionEq,
				Value:     "yes",
			},
		},
	}, d, conf.Conf)

	Convey("test filter eq ", t, func() {
		fy := []byte(`{"table":"yes"}`)
		So(nt.filter(fy), ShouldBeFalse)
		fy = []byte(`{"table":"no"}`)
		So(nt.filter(fy), ShouldBeTrue)
	})
	Convey("test filter prefix ", t, func() {
		nt.w.Filters = []*model.Filter{
			{
				Field:     "table",
				Condition: model.ConditionPre,
				Value:     "yes",
			},
		}
		nt.parseFilter()
		fy := []byte(`{"table":"yes1"}`)
		So(nt.filter(fy), ShouldBeFalse)
		fy = []byte(`{"table":"no"}`)
		So(nt.filter(fy), ShouldBeTrue)
	})
}

func TestParseNotifyURL(t *testing.T) {
	Convey("test url parse result", t, func() {
		u := "liverpc://live.bannedservice?version=0&cmd=Message.synUser"
		notifyURL, err := parseNotifyURL(u)
		So(err, ShouldBeNil)
		So(notifyURL.RawURL, ShouldEqual, u)
		So(notifyURL.Schema, ShouldEqual, "liverpc")
		So(notifyURL.Host, ShouldEqual, "live.bannedservice")
		So(notifyURL.Query.Get("version"), ShouldEqual, "0")
		So(notifyURL.Query.Get("cmd"), ShouldEqual, "Message.synUser")
	})
}
