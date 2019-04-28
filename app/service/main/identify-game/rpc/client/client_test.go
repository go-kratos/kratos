package client

import (
	"context"
	"testing"
	"time"

	"go-common/app/service/main/identify-game/model"
)

func TestNew(t *testing.T) {
	cli := New(nil)
	time.Sleep(2 * time.Second)
	if err := cli.DelCache(context.Background(), &model.CleanCacheArgs{Token: "1234567890123456789012"}); err != nil {
		t.FailNow()
	}

	testDelCache(t, cli)
}

func testDelCache(t *testing.T, cli *Client) {
	if err := cli.DelCache(context.Background(), &model.CleanCacheArgs{Token: "1234567890123456789012"}); err != nil {
		t.Errorf("%v\n", err)
	}
}
