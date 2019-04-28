package tcp

import (
	"fmt"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"go-common/app/infra/databus/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	pubCfg = &conf.Kafka{
		Cluster: "test_topic",
		Brokers: []string{"172.22.33.174:9092", "172.22.33.183:9092", "172.22.33.185:9092"},
	}
)

func TestDatabus(t *testing.T) {
	Convey("Test publish:", t, func() {
		l, _ := net.Listen("tcp", ":8888")
		go func() {
			for {
				conn, err := l.Accept()
				if err != nil {
					continue
				}
				b, err := ioutil.ReadAll(conn)
				if err == nil {
					fmt.Printf("test conn: %s", b)
				}
				conn.Close()
			}
		}()
		conn, err := net.Dial("tcp", ":8888")
		So(err, ShouldBeNil)
		p, err := NewPub(newConn(conn, time.Second, time.Second), "pub", "", _testTopic, pubCfg)
		So(err, ShouldBeNil)
		key := []byte("key")
		header := []byte("header")
		msg := []byte("message")
		err = p.publish(key, header, msg)
		So(err, ShouldBeNil)
		time.Sleep(time.Second)
		Convey("test sub", func() {
			conn, _ := net.Dial("tcp", ":8888")
			s, err := NewSub(newConn(conn, time.Second, time.Second), "sub", "", _testTopic, pubCfg, 1)
			So(err, ShouldBeNil)
			t.Logf("subscriptions: %v", s.consumer.Subscriptions())
			for {
				select {
				case msg := <-s.consumer.Messages():
					s.consumer.CommitOffsets()
					t.Logf("sub message: %s timestamp: %d", msg.Value, msg.Timestamp.Unix())
					return
				case err := <-s.consumer.Errors():
					t.Errorf("error: %v", err)
					So(err, ShouldBeNil)
				case n := <-s.consumer.Notifications():
					t.Logf("notify: %v", n)
					err := p.publish(key, header, msg)
					So(err, ShouldBeNil)
				}
			}
		})
	})
}
