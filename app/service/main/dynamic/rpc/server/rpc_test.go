package server

import (
	"net/rpc"
	"testing"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/dynamic/model"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	_addr           = "127.0.0.1:6235"
	_regionArcs3    = "RPC.RegionArcs3"
	_regionTagArcs3 = "RPC.RegionTagArcs3"
	_regionsArcs3   = "RPC.RegionsArcs3"
	_regionTotal    = "RPC.RegionTotal"
)

func Test_RPC(t *testing.T) {
	Convey("RPC Region Archives", t, func() {
		client, err := rpc.Dial("tcp", _addr)
		So(err, ShouldBeNil)
		defer client.Close()

		testRPCRegionArcs(t, client)
		testRPCRegionTagArcs(t, client)
		testRPCRegionsArcs(t, client)
		testRPCRegionTotal(t, client)
	})
}

func testRPCRegionArcs(t *testing.T, client *rpc.Client) {
	var (
		rid    int32 = 168
		pn, ps int   = 1, 10
		ip           = "127.0.0.1"
		res    *model.DynamicArcs3
	)
	Convey("RPC Region Archives", func() {
		arg := &model.ArgRegion3{RegionID: rid, Pn: pn, Ps: ps, RealIP: ip}
		res = new(model.DynamicArcs3)
		err := client.Call(_regionArcs3, arg, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func testRPCRegionTagArcs(t *testing.T, client *rpc.Client) {
	var (
		rid    int32 = 168
		tid    int64 = 123456
		pn, ps int   = 1, 10
		ip           = "127.0.0.1"
		res    *model.DynamicArcs3
	)
	Convey("RPC Region Tag Archives", func() {
		arg := &model.ArgRegionTag3{TagID: tid, RegionID: rid, Pn: pn, Ps: ps, RealIP: ip}
		res = new(model.DynamicArcs3)
		err := client.Call(_regionTagArcs3, arg, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func testRPCRegionsArcs(t *testing.T, client *rpc.Client) {
	var (
		rids  = []int32{1, 3, 4, 5, 13, 36, 129, 119, 23, 11, 155, 160, 165, 168}
		count = 10
		ip    = "127.0.0.1"
		res   *map[int32][]*api.Arc
	)
	Convey("RPC Regions Archives", func() {
		arg := &model.ArgRegions3{RegionIDs: rids, Count: count, RealIP: ip}
		err := client.Call(_regionsArcs3, arg, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})

}

func testRPCRegionTotal(t *testing.T, client *rpc.Client) {
	var (
		ip  = "127.0.0.1"
		res map[string]int
	)
	Convey("RPC Region Total", func() {
		arg := &model.ArgRegionTotal{RealIP: ip}
		err := client.Call(_regionTotal, arg, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
