package v1

import (
	"context"
	"fmt"
	"testing"

	"go-common/library/net/rpc/warden"
)

var (
	err    error
	rpccli PushClient
)

func init() {
	rpccli, err = NewClient(&warden.ClientConfig{})
	if err != nil {
		fmt.Printf("new push grpc client error(%v)", err)
	}
}

func Test_AddReport(t *testing.T) {
	_, err := rpccli.AddReport(context.Background(), &AddReportRequest{Report: &ModelReport{
		APPID:        1,
		Mid:          91221505,
		DeviceToken:  "tototototototo",
		NotifySwitch: 1,
	}})
	if err != nil {
		t.Errorf("AddReport error(%v)", err)
	}
}

func Test_AddTokenCache(t *testing.T) {
	_, err := rpccli.AddTokenCache(context.Background(), &AddTokenCacheRequest{Report: &ModelReport{
		APPID:        1,
		Mid:          91221505,
		DeviceToken:  "tototototototo",
		NotifySwitch: 1,
	}})
	if err != nil {
		t.Errorf("AddTokenCache error(%v)", err)
	}
}

func Test_AddTokensCache(t *testing.T) {
	_, err := rpccli.AddTokensCache(context.Background(), &AddTokensCacheRequest{Reports: []*ModelReport{{
		APPID:        1,
		Mid:          91221505,
		DeviceToken:  "tototototototo",
		NotifySwitch: 1,
	}}})
	if err != nil {
		t.Errorf("AddTokensCache error(%v)", err)
	}
}
