package service

import (
	"testing"
	"time"

	"github.com/Shopify/sarama"
	. "github.com/smartystreets/goconvey/convey"
)

// TestStartConsume .
func TestStartConsume(t *testing.T) {
	Convey("start consume", t, func() {
		err := svr.StartConsume()
		So(err, ShouldNotBeNil)
	})
}

func TestStartHandle(t *testing.T) {
	go svr.handleMsg()
}

// TestHandle .
func TestHandle(t *testing.T) {
	Convey("handle msg", t, func() {
		var l = `a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z|1|2|3|4|5|6|7|8`
		msg := &sarama.ConsumerMessage{
			Value: []byte(l),
		}
		svr.consumer.messages <- msg
		time.Sleep(time.Second)
		So(len(svr.consumer.messages), ShouldEqual, 0)
	})
}

// TestClose .
func TestClose(t *testing.T) {
	svr.Close()
}
