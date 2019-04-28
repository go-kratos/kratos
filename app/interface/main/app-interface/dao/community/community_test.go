package community

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"go-common/app/interface/main/app-interface/conf"

	. "github.com/smartystreets/goconvey/convey"
)

// go test -conf="../../app-interface-example.toml"  -v -test.run TestCommunity
func TestCommunity(t *testing.T) {
	Convey("TestCommuity", t, func() {

	})
	err := conf.Init()
	if err != nil {
		return
	}
	dao := New(conf.Conf)
	community, _, err := dao.Community(context.TODO(), 28009145, "2e3950631afd879592de5e2ee34c7293", "android", 1, 20)
	if err != nil {
		t.Errorf("dao.Community error(%v)", err)
		return
	}
	result, err := json.Marshal(community)
	if err != nil {
		t.Errorf("json.Marshal error(%v)", err)
		return
	}
	fmt.Printf("test community (%v) \n", string(result))
}
