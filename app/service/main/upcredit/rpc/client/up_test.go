package client

import (
	"testing"
	"time"
)

func TestRpcClient(t *testing.T) {

	var _ = New(nil)

	time.Sleep(1 * time.Second)

}
