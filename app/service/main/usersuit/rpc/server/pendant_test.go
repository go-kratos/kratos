package server

import (
	"go-common/app/service/main/usersuit/model"
	"testing"
)

const (
	_pendantEquip  = "RPC.Equipment"
	_pendantEquips = "RPC.Equipments"
)

// TestRPC_Equipment test
func TestRPC_Equipment(t *testing.T) {
	var (
		res = new(model.PendantEquip)
		err error
	)
	once.Do(startServer)
	arg := &model.ArgEquipment{
		Mid: 27515240,
	}
	if err = client.Call(_pendantEquip, arg.Mid, &res); err != nil {
		t.Errorf("client.Call(%s) error(%v)", _pendantEquip, err)
		t.FailNow()
	}
	t.Logf("res (%v)", res)
}

// TestRPC_Equipment test
func TestRPC_Equipments(t *testing.T) {
	var (
		res = make(map[int64]*model.PendantEquip)
		err error
	)
	once.Do(startServer)
	arg := &model.ArgEquipments{
		Mids: []int64{27515240, 100},
	}
	if err = client.Call(_pendantEquips, arg.Mids, &res); err != nil {
		t.Errorf("client.Call(%s) error(%v)", _pendantEquips, err)
		t.FailNow()
	}
	t.Logf("res (%v)", res)
}
