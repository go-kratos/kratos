package notify

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/infra/notify/model"
	"go-common/library/net/rpc/liverpc"
	"testing"
)

func initNc(urls []string) *LiverpcClients {
	w := &model.Watcher{}
	w.Callbacks = make([]*model.Callback, 0, len(urls))
	for i, u := range urls {
		notifyURL, err := parseNotifyURL(u)
		So(err, ShouldBeNil)
		cb := &model.Callback{
			URL:      notifyURL,
			Priority: int8(i),
		}
		w.Callbacks = append(w.Callbacks, cb)
	}
	nc := newLiverpcClients(w)
	return nc
}

func TestNewLiverpcClients(t *testing.T) {
	Convey("test new liverpc clients", t, func() {
		urls := []string{
			"http://www.bilibili.com",
			"liverpc://live.bannedservice?version=0&cmd=Message.synUser&addr=172.18.33.82:20822",
			"liverpc://live.room?version=1&cmd=Consumer/receiveGift&addr=172.18.33.82:20200",
		}
		nc := initNc(urls)
		So(len(nc.clients), ShouldEqual, 2)
	})
}

func TestFormatLiveMsg(t *testing.T) {
	var (
		msg    []byte
		header *liverpc.Header
		body   map[string]string
		err    error
		urls   = []string{}
		nc     = initNc(urls)
	)

	Convey("test live post message format result", t, func() {
		msg = []byte(`{"topic":"BannedUserSyn-T","msg_id":"a485ad75609304b920782462ce1c7632","msg_content":"{\"uid\":1734992,\"status\":0,\"begin\":\"2018-09-10 17:51:45\",\"uname\":\"\\u83ca\\u82b1\\u75db\",\"face\":\"http:\\\/\\\/i2.hdslb.com\\\/bfs\\\/face\\\/9eab4e877c83dd77bd010994e35e0d113ec7bf9d.jpg\",\"rank\":\"10000\",\"identification\":0,\"mobile_verify\":1,\"silence\":0,\"official_verify\":{\"type\":-1,\"desc\":\"\",\"role\":0}}","msg_key":1734992,"timestamp":1536573105.4239,"failure_cnt":0,"caller_header":{"platform":"","src":"","version":"","buvid":"AUTO3715365731053968","trace_id":"6172727a27c922e0:61727287a8ae0264:61727285b093e28e:0","uid":0,"caller":"user.user\\common\\logic\\Databus_Service.call-40","user_ip":"172.18.29.22","source_group":"default","sessdata2":"access_key=&SESSDATA=","group":"default"},"__ts1":1536573105.3978,"__ts2":1536573105.4032}`)
		header, body, err = nc.formatLiveMsg(string(msg))
		So(err, ShouldBeNil)
		So(header.TraceId, ShouldEqual, "6172727a27c922e0:61727287a8ae0264:61727285b093e28e:0")
		So(header.Caller, ShouldEqual, _liverpcCaller)
		So(body["msg"], ShouldNotEqual, "")
		So(body["msg_content"], ShouldNotEqual, "")
	})

	Convey("test non live post message format result", t, func() {
		msg = []byte(`{"action":"insert","table":"user_relation_fid_439","new":{"attribute":2,"ctime":"2018-09-10 18:56:26","fid":18021939,"id":4090654,"mid":354291964,"mtime":"2018-09-10 18:56:26","source":0,"status":0}}`)
		header, body, err = nc.formatLiveMsg(string(msg))
		So(err, ShouldBeNil)
		So(header.Caller, ShouldEqual, _liverpcCaller)
		So(body["msg"], ShouldEqual, "")
		So(body["msg_content"], ShouldNotEqual, "")
	})
}

func TestLiverpcClients_Post(t *testing.T) {
	Convey("test liverpc client post", t, func() {
		Convey("test post with invalid params", func() {
			u := "liverpc://live.bannedservice"
			urls := []string{u}
			nc := initNc(urls)
			notifyURL, err := parseNotifyURL(u)
			So(err, ShouldBeNil)
			err = nc.Post(context.Background(), notifyURL, "test")
			So(err, ShouldResemble, errLiverpcInvalidParams)

			u = "liverpc://live.bannedservice?version=error&cmd=Message.synUser"
			nc = initNc([]string{u})
			notifyURL, err = parseNotifyURL(u)
			So(err, ShouldBeNil)
			err = nc.Post(context.Background(), notifyURL, "test")
			So(err, ShouldResemble, errStrconvVersion)
		})

		Convey("test post with invalid msg", func() {
			u := "liverpc://live.bannedservice?version=0&cmd=Message.synUser"
			nc := initNc([]string{u})
			notifyURL, err := parseNotifyURL(u)
			So(err, ShouldBeNil)
			err = nc.Post(context.Background(), notifyURL, "not a json format message")
			So(err, ShouldResemble, nil)
		})

		Convey("test post call client failed with invalid addr", func() {
			u := "liverpc://live.bannedservice?version=0&cmd=Message.synUser&addr=1.1.1.1:11111"
			nc := initNc([]string{u})
			notifyURL, err := parseNotifyURL(u)
			So(err, ShouldBeNil)
			m := []byte(`{"topic":"BannedUserSyn-T","msg_id":"123456","msg_content":"{\"test\":123456}","msg_key":1734992,"timestamp":1536573105.4239}`)
			err = nc.Post(context.Background(), notifyURL, string(m))
			So(err, ShouldResemble, errCallRaw)
		})

		Convey("test success", func() {
			u := "liverpc://live.bannedservice?version=0&cmd=Message.synUser"
			nc := initNc([]string{u})
			notifyURL, err := parseNotifyURL(u)
			So(err, ShouldBeNil)
			m := []byte(`{"topic":"BannedUserSyn-T","msg_id":"a485ad75609304b920782462ce1c7632","msg_content":"{\"uid\":1734992,\"status\":0,\"begin\":\"2018-09-10 17:51:45\",\"uname\":\"\\u83ca\\u82b1\\u75db\",\"face\":\"http:\\\/\\\/i2.hdslb.com\\\/bfs\\\/face\\\/9eab4e877c83dd77bd010994e35e0d113ec7bf9d.jpg\",\"rank\":\"10000\",\"identification\":0,\"mobile_verify\":1,\"silence\":0,\"official_verify\":{\"type\":-1,\"desc\":\"\",\"role\":0}}","msg_key":1734992,"timestamp":1536573105.4239}`)
			err = nc.Post(context.Background(), notifyURL, string(m))
			So(err, ShouldBeNil)
		})
	})
}
