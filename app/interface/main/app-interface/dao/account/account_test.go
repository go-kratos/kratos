package account

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"go-common/app/interface/main/app-interface/conf"

	. "github.com/smartystreets/goconvey/convey"
)

//go test -conf="../../app-interface-example.toml"  -v -test.run TestCard
func TestCard(t *testing.T) {
	Convey("TestCard", t, func() {
		err := conf.Init()
		if err != nil {
			return
		}
		dao := New(conf.Conf)
		card, err := dao.Profile3(context.TODO(), 28009145)
		if err != nil {
			t.Errorf("dao.Profile3 error(%v)", err)
			return
		}
		result, err := json.Marshal(card)
		if err != nil {
			t.Errorf("json.Marshal error(%v)", err)
			return
		}
		fmt.Printf("test card (%v) \n", string(result))
	})
}

//go test -conf="../../app-interface-example.toml"  -v -test.run TestCardByName
func TestCardByName(t *testing.T) {
	err := conf.Init()
	if err != nil {
		return
	}
	dao := New(conf.Conf)
	card, err := dao.ProfileByName3(context.TODO(), "冠冠爱看书")
	if err != nil {
		t.Errorf("dao.ProfileByName3 error(%v)", err)
		return
	}
	result, err := json.Marshal(card)
	if err != nil {
		t.Errorf("json.Marshal error(%v)", err)
		return
	}
	fmt.Printf("test card (%v) \n", string(result))
}
