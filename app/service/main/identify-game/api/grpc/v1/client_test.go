package v1

import (
	"context"
	"testing"

	"go-common/library/net/rpc/warden"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClients_Post(t *testing.T) {
	var (
		ctx     = context.Background()
		err     error
		res     *DelCacheReply
		cookies *CreateCookieReply
	)
	Convey("test delete cache", t, func() {
		client, _ := NewClient(&warden.ClientConfig{})
		res, err = client.DelCache(ctx, &DelCacheReq{Token: "24325defrf"})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
	Convey("test get cookie by token", t, func() {
		client, _ := NewClient(&warden.ClientConfig{})
		cookies, err = client.GetCookieByToken(ctx, &CreateCookieReq{Token: "24325defrf"})
		So(err, ShouldNotBeNil)
		So(cookies, ShouldBeNil)
	})
}
