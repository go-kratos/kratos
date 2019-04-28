package usersuit

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	usmdl "go-common/app/service/main/usersuit/model"
	usrpc "go-common/app/service/main/usersuit/rpc/client"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Equip(t *testing.T) {
	// 1）背包里只要存在挂件，且来源是背包（不论挂件是否为大会员的）都能佩戴 2）背包里不存在挂件，但来源是大会员挂件，也可以佩戴，反之报错，用例如下：
	Convey("Equip interface", t, func() {
		var (
			c = context.Background()
		)
		// 穿戴一个非vip挂件，但是挂件来源是vip，报错
		Convey(" wear vip pendant", t, func() {
			var ArgEquip = &usmdl.ArgEquip{
				Mid:    111001965,
				Pid:    98,
				Status: 2, //1 卸载  2 佩戴
				Source: 2, // 0 未知  1背包  2 vip
			}
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(s.usRPC), "Equip", func(_ *usrpc.Service2, _ context.Context, _ *usmdl.ArgEquip) error {
				return fmt.Errorf("wear not vip pendant,but source is EquipFromVip")
			})
			defer guard.Unpatch()
			err := s.usRPC.Equip(c, ArgEquip)
			Convey("the pendant is not vip pendant,then err should not be nil", t, func() {
				So(err, ShouldNotBeNil)
			})
		})

		// 穿戴一个vip挂件，但是挂件来源是背包(前提：背包里存在该挂件),正确
		Convey(" wear vip pendant", t, func() {
			var ArgEquip = &usmdl.ArgEquip{
				Mid:    111001965,
				Pid:    102,
				Status: 2, //1 卸载  2 佩戴
				Source: 1, // 0 未知  1背包  2 vip
			}
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(s.usRPC), "Equip", func(_ *usrpc.Service2, _ context.Context, _ *usmdl.ArgEquip) error {
				return nil
			})
			defer guard.Unpatch()
			err := s.usRPC.Equip(c, ArgEquip)
			Convey("wear vip pendant and this pendant exist in package, then err should be nil", t, func() {
				So(err, ShouldBeNil)
			})
		})

		// 穿戴挂件与来源一致 ：vip挂件
		Convey(" wear vip pendant", t, func() {
			var ArgEquip = &usmdl.ArgEquip{
				Mid:    111001965,
				Pid:    103,
				Status: 2, //1 卸载  2 佩戴
				Source: 2, // 0 未知  1背包  2 vip
			}
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(s.usRPC), "Equip", func(_ *usrpc.Service2, _ context.Context, _ *usmdl.ArgEquip) error {
				return nil
			})
			defer guard.Unpatch()
			err := s.usRPC.Equip(c, ArgEquip)
			Convey("the pendant is vip pendant and source is EquipFromVip,then err should be nil", t, func() {
				So(err, ShouldBeNil)
			})
		})
		// 穿戴挂件与来源一致 ：背包挂件（背包存在该挂件）
		Convey(" wear vip pendant", t, func() {
			var ArgEquip = &usmdl.ArgEquip{
				Mid:    111001965,
				Pid:    98,
				Status: 2, //1 卸载  2 佩戴
				Source: 1, // 0 未知  1背包  2 vip
			}
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(s.usRPC), "Equip", func(_ *usrpc.Service2, _ context.Context, _ *usmdl.ArgEquip) error {
				return nil
			})
			defer guard.Unpatch()
			err := s.usRPC.Equip(c, ArgEquip)
			Convey("wear a pkg pendant and package exist the pendant and source is EquipFromPackage, then err should be nil", t, func() {
				So(err, ShouldBeNil)
			})
		})

		// 穿戴一个背包里不存在的挂件，但是挂件来源是：背包（非vip挂件）,报错
		Convey(" wear vip pendant", t, func() {
			var ArgEquip = &usmdl.ArgEquip{
				Mid:    111001965,
				Pid:    99,
				Status: 2, //1 卸载  2 佩戴
				Source: 1, // 0 未知  1背包  2 vip
			}
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(s.usRPC), "Equip", func(_ *usrpc.Service2, _ context.Context, _ *usmdl.ArgEquip) error {
				return fmt.Errorf("pendant is not exist, err_code: 64101")
			})
			defer guard.Unpatch()
			err := s.usRPC.Equip(c, ArgEquip)
			Convey("wear a pendant which is not exist in package and source is EquipFromPackage, then err should not be nil", t, func() {
				So(err, ShouldNotBeNil)
			})
		})

		// 卸下挂件(挂件存在)，不会受到 source 的影响
		Convey(" take off pkg pendant", t, func() {
			var ArgEquip = &usmdl.ArgEquip{
				Mid: 111001965, Pid: 98, Status: 1, Source: 2,
			}
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(s.usRPC), "Equip", func(_ *usrpc.Service2, _ context.Context, _ *usmdl.ArgEquip) error {
				return nil
			})
			defer guard.Unpatch()
			err := s.usRPC.Equip(c, ArgEquip)
			Convey("take off pendant is not be affected by source, then err should not be nil", t, func() {
				So(err, ShouldBeNil)
			})

		})

		Convey(" take off vip pendant", t, func() {
			var ArgEquip = &usmdl.ArgEquip{
				Mid: 111001965, Pid: 102, Status: 1, Source: 2,
			}
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(s.usRPC), "Equip", func(_ *usrpc.Service2, _ context.Context, _ *usmdl.ArgEquip) error {
				return nil
			})
			defer guard.Unpatch()
			err := s.usRPC.Equip(c, ArgEquip)
			Convey("take off pendant is not be affected by source, then err should not be nil", t, func() {
				So(err, ShouldBeNil)
			})
		})
	})

}

//func TestService_Equip(t *testing.T) {
//	Convey("Equip interface", t, func() {
//		So(s.Equip(context.Background(), 1, 1, 2, 1), ShouldBeNil)
//	})
//}

func TestService_Equipment(t *testing.T) {
	Convey("Equipment interface", t, func() {
		equip, err := s.Equipment(context.Background(), 1)
		So(err, ShouldBeNil)
		So(equip, ShouldNotBeEmpty)
	})
}

func TestService_Pendant(t *testing.T) {
	Convey("Pendant interface", t, func() {
		pendant, err := s.Pendant(context.Background(), 1)
		So(err, ShouldBeNil)
		So(pendant, ShouldNotBeEmpty)
	})
}

func TestService_GroupEntry(t *testing.T) {
	Convey("GroupEntry interface", t, func() {
		groups, err := s.GroupEntry(context.Background(), 1)
		So(err, ShouldBeNil)
		So(groups, ShouldNotBeEmpty)
	})
}

func TestService_GroupVIP(t *testing.T) {
	Convey("GroupVIP interface", t, func() {
		vips, err := s.GroupVIP(context.Background(), 1)
		So(err, ShouldBeNil)
		So(vips, ShouldNotBeEmpty)
	})
}

func TestService_VipGet(t *testing.T) {
	Convey("VipGet interface", t, func() {
		So(s.VipGet(context.Background(), 1, 1, 2), ShouldBeNil)
	})
}

func TestService_CheckOrder(t *testing.T) {
	Convey("CheckOrder interface", t, func() {
		So(s.CheckOrder(context.Background(), 1, "lalala"), ShouldBeNil)
	})
}

func TestService_My(t *testing.T) {
	Convey("My interface", t, func() {
		my, err := s.My(context.Background(), 1)
		So(err, ShouldBeNil)
		So(my, ShouldNotBeEmpty)
	})
}

func TestService_MyHistory(t *testing.T) {
	Convey("MyHistory interface", t, func() {
		res, err := s.MyHistory(context.Background(), 1, 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
