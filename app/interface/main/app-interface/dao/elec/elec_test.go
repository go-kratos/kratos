package elec

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"go-common/app/interface/main/app-interface/conf"

	. "github.com/smartystreets/goconvey/convey"
)

// go test -conf="../../app-interface-example.toml"  -v -test.run TestElec
func TestElec(t *testing.T) {
	Convey("TestElec", t, func() {

	})
	err := conf.Init()
	if err != nil {
		return
	}
	dao := New(conf.Conf)
	elec, err := dao.Info(context.TODO(), 5461533, 15555180)
	if err != nil {
		t.Errorf("dao.Elec error(%v)", err)
		return
	}
	result, err := json.Marshal(elec)
	if err != nil {
		t.Errorf("json.Marshal error(%v)", elec)
		return
	}
	fmt.Printf("test elec (%v) \n", string(result))
}
