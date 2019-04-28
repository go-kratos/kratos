package service

import (
	"context"
	"fmt"
	"testing"
)

func TestService_Limits(t *testing.T) {
	if res, err := svr.Limits(context.TODO(), "msm-service", ""); err != nil {
		t.Logf("svr.Limits() error(%v)", err)
		t.FailNow()
	} else {
		fmt.Println(res)
	}
}
