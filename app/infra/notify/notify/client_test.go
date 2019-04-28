package notify

import (
	"context"
	"flag"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/infra/notify/conf"
	"go-common/app/infra/notify/model"
)

func initConf() *conf.Config {
	var err error
	flag.Set("conf", "../cmd/notify-test.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}
	return conf.Conf
}

func TestClients_Post(t *testing.T) {
	var (
		ctx       = context.Background()
		c         = initConf()
		err       error
		notifyURL *model.NotifyURL
	)
	Convey("test call post with different urls", t, func() {
		w := &model.Watcher{}

		Convey("error schema", func() {
			u := "https://live.bilibili.com"
			notifyURL, err = parseNotifyURL(u)
			So(err, ShouldBeNil)
			nc := NewClients(c, w)
			err = nc.Post(ctx, notifyURL, "test1")
			So(err, ShouldResemble, errUnknownSchema)
		})

		Convey("test http url post", func() {
			u := "http://127.0.0.1:19999/push3"
			notifyURL, err = parseNotifyURL(u)
			So(err, ShouldBeNil)
			nc := NewClients(c, w)
			mockhttp2()
			err = nc.Post(ctx, notifyURL, "test2")
			So(err, ShouldBeNil)
		})

		Convey("test liverpc url post", func() {
			u := "liverpc://live.bannedservice?version=0&cmd=Message.synUser"
			notifyURL, err = parseNotifyURL(u)
			cb := &model.Callback{
				URL:      notifyURL,
				Priority: 1,
			}
			w.Callbacks = []*model.Callback{cb}
			So(err, ShouldBeNil)
			nc := NewClients(c, w)
			So(len(nc.liverpcClients.clients), ShouldEqual, 1)
			msg := []byte(`{"topic":"BannedUserSyn-T","msg_id":"a485ad75609304b920782462ce1c7632","msg_content":"{\"uid\":1734992,\"status\":0,\"begin\":\"2018-09-10 17:51:45\",\"uname\":\"\\u83ca\\u82b1\\u75db\",\"face\":\"http:\\\/\\\/i2.hdslb.com\\\/bfs\\\/face\\\/9eab4e877c83dd77bd010994e35e0d113ec7bf9d.jpg\",\"rank\":\"10000\",\"identification\":0,\"mobile_verify\":1,\"silence\":0,\"official_verify\":{\"type\":-1,\"desc\":\"\",\"role\":0}}","msg_key":1734992,"timestamp":1536573105.4239,"failure_cnt":0,"caller_header":{"platform":"","src":"","version":"","buvid":"AUTO3715365731053968","trace_id":"6172727a27c922e0:61727287a8ae0264:61727285b093e28e:0","uid":0,"caller":"user.user\\common\\logic\\Databus_Service.call-40","user_ip":"172.18.29.22","source_group":"qa01","sessdata2":"access_key=&SESSDATA=","group":"default"},"__ts1":1536573105.3978,"__ts2":1536573105.4032}`)
			err = nc.Post(ctx, notifyURL, string(msg))
			So(err, ShouldBeNil)
		})

		Convey("test liverpc url post with different message format", func() {
			u := "liverpc://live.relation?version=1&cmd=Notify.on_relation_changed"
			notifyURL, err = parseNotifyURL(u)
			cb := &model.Callback{
				URL:      notifyURL,
				Priority: 1,
			}
			w.Callbacks = []*model.Callback{cb}
			So(err, ShouldBeNil)
			nc := NewClients(c, w)
			So(len(nc.liverpcClients.clients), ShouldEqual, 1)
			msg := []byte(`{"action":"insert","table":"user_relation_fid_439","new":{"attribute":2,"ctime":"2018-09-10 18:56:26","fid":18021939,"id":4090654,"mid":354291964,"mtime":"2018-09-10 18:56:26","source":0,"status":0}}`)
			err = nc.Post(ctx, notifyURL, string(msg))
			So(err, ShouldBeNil)
		})
	})
}

func mockhttp2() {
	http.HandleFunc("/push3", testpush1)
	http.HandleFunc("/push4", testpush2)
	go http.ListenAndServe(":19999", nil)
}
