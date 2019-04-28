package rpc

import (
	"net/rpc"
	"testing"

	coidel "go-common/app/service/main/coin/model"
)

const (
	addr  = "172.16.12.122:6159"
	mid   = 23675773
	aid   = 1
	added = 1
	ip    = "172.16.12.122"

	coinInfo = "RPC.ArchiveUserCoins"
	addCoin  = "RPC.AddCoins"
)

func TestAddCoinsRpc(t *testing.T) {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		t.Errorf("rpc.Dial(tcp, (%s)) error(%v)", addr, err)
		t.FailNow()
	}
	x := coidel.ArgAddCoin{Aid: aid, Mid: mid, Multiply: added, RealIP: ip}
	cf := &coidel.ArchiveUserCoins{}
	if err = client.Call(addCoin, x, cf); err != nil {
		t.Logf("call.addMoral error(%v)", err)
	}

	t.Logf("res: %v", cf.Multiply)
}

func TestArchiveUserCoinsRpc(t *testing.T) {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		t.Errorf("rpc.Dial(tcp, (%s)) error(%v)", addr, err)
		t.FailNow()
	}
	x := coidel.ArgCoinInfo{Aid: aid, Mid: mid}
	cf := &coidel.ArchiveUserCoins{}
	if err = client.Call(coinInfo, x, cf); err != nil {
		t.Logf("call.addMoral error(%v)", err)
	}

	t.Logf("res: %v", cf.Multiply)
}
