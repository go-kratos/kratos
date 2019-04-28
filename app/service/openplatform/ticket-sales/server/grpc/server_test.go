package grpc

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	sales "go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-sales/conf"
	"go-common/app/service/openplatform/ticket-sales/service"

	"go-common/app/service/openplatform/ticket-sales/api/grpc/type"
	"go-common/library/conf/paladin"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *service.Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/ticket-sales.toml")
	flag.Set("conf", dir)
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	if err := paladin.Watch("ticket-sales.toml", conf.Conf); err != nil {
		panic(err)
	}
	s = service.New(conf.Conf)
}

//Test_Info
func TestInfo(t *testing.T) {
	Convey("get data", t, func() {

		data := &sales.UpBuyerRequest{}

		data.OrderID = 10000012002571
		data.Buyers = &_type.OrderBuyer{
			ID:         1,
			Name:       "wlt",
			Tel:        "1388888888",
			PersonalID: "360822199207227275",
		}
		res, err := s.UpdateBuyer(context.TODO(), data)
		fmt.Println(res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res, ShouldNotBeEmpty)
	})

}
