package service

import (
	"context"
	"github.com/bouk/monkey"
	"reflect"
	"testing"
	"time"

	"go-common/app/service/main/usersuit/dao/pendant"
	"go-common/app/service/main/usersuit/model"
	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {
	now := time.Now()
	t.Logf("%d, %d", now.Unix(), now.UnixNano()/1e6)
	var x = []string{"A", "B", "C"}
	for index, value := range x {
		t.Logf("index is (%+v) ,value is (%+v)", index, value)
	}
}

// TestService_PendantAll test get all pendant infomartion
func TestService_PendantAll(t *testing.T) {
	var (
		res = make([]*model.Pendant, 0)
		err error
	)
	Convey("need return something", t, func() {
		res, err = s.PendantAll(context.Background())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestService_PendantInfoByID(t *testing.T) {
	var (
		pid = int64(1)
		res *model.Pendant
		err error
	)
	Convey("need return something", t, func() {
		res, err = s.PendantInfo(context.Background(), pid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

// TestService_PendantPoint test PendantPoint
func TestService_PendantPoint(t *testing.T) {
	var (
		mid int64 = 20606508
		err error
		res = make([]*model.Pendant, 0)
	)
	Convey("need return something", t, func() {
		res, err = s.PendantPoint(context.Background(), mid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

// TestService_OrderHistory test order info
func TestService_OrderHistory(t *testing.T) {
	var (
		arg = &model.ArgOrderHistory{
			OrderID:   "201711141510654974657997844",
			PayID:     "201711141822545083995814",
			Mid:       650454,
			Pid:       17,
			Status:    2,
			PayType:   1,
			StartTime: 1500311684,
			EndTime:   1520311684,
			Page:      1,
		}
		res   = make([]*model.PendantOrderInfo, 0)
		count = make(map[string]int64)
		err   error
	)
	Convey("need return something", t, func() {
		res, count, err = s.OrderHistory(context.Background(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(count, ShouldNotBeNil)
	})
}

// TestService_PackageByMid test package
func TestService_PackageByMid(t *testing.T) {
	var (
		mid int64 = 88889021
		res       = make([]*model.PendantPackage, 0)
		err error
	)
	Convey("need return something", t, func() {
		res, err = s.PackageInfo(context.Background(), mid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestService_Equipment(t *testing.T) {
	var (
		mid        int64 = 111001965
		vipPendant       = &model.Pendant{
			ID:     102,
			Name:   "文豪",
			Status: 1, // 挂件状态:0.下线;1.在线
			Gid:    1, // 挂件所在组
			Rank:   1,
		}
		pendantEquip = &model.PendantEquip{
			Mid:     111001965,
			Pid:     102,
			Expires: 1594076284,
			Pendant: vipPendant,
		}
	)
	Convey("need return something", t, func() {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "EquipCache", func(_ *pendant.Dao, _ context.Context, _ int64) (*model.PendantEquip, error) {
			return pendantEquip, nil
		})
		defer guard.Unpatch()
		res, err := s.Equipment(context.Background(), mid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestService_GroupInfo(t *testing.T) {
	var (
		res []*model.PendantGroupInfo
		err error
	)
	Convey("need return something .", t, func() {
		res, err = s.GroupInfo(context.Background())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

// OrderPendant(c context.Context, mid, pid, expires, tp int64, ip string)

func TestService_OrderPendant(t *testing.T) {
	var (
		mid int64 = 88889021
		pid int64 = 10
		res *model.PayInfo
		err error
	)
	Convey("need return something", t, func() {
		res, err = s.OrderPendant(context.Background(), mid, pid, time.Now().Unix(), 1)
		So(err, ShouldNotBeNil)
		So(res, ShouldBeNil)
	})

}

func TestService_TakeOffPendant(t *testing.T) {
	var (
		mid int64 = 111001965
		err error
	)
	Convey("need return something", t, func() {
		monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "UpEquipMID", func(_ *pendant.Dao, _ context.Context, _ int64) (int64, error) {
			return 0, nil
		})
		monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "DelEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64) error {
			return nil
		})
		defer monkey.UnpatchAll()
		err = s.TakeOffPendant(context.Background(), mid)
		So(err, ShouldBeNil)
	})
}

func TestService_PendantInfo(t *testing.T) {
	var (
		pid int64 = 171
	)
	Convey("when everything goes well", t, func() {
		res, err := s.PendantInfo(context.Background(), pid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})

	Convey("when occur an err", t, func() {
		res, err := s.PendantInfo(context.Background(), 102)
		So(err, ShouldNotBeNil)
		So(res, ShouldBeNil)
	})
}

func TestService_PackageByID(t *testing.T) {
	var (
		mid        int64 = 111001965
		pid        int64 = 102
		pendantPkg       = &model.PendantPackage{
			ID:      16281,
			Mid:     111001965,
			Pid:     98,
			Expires: 1594076284,
			Type:    0,
			Status:  1,
			IsVIP:   0,
			Pendant: &model.Pendant{
				ID:         1,
				Name:       "文豪",
				Image:      "/bfs/face/2e699edf6b51f61fc501247d4e826c97eebfb4fe.png",
				ImageModel: "/bfs/face/95aa23fa00e619330a10e55cff328a7cace08e32.png",
				Status:     1,  // 挂件状态:0.下线;1.在线
				Gid:        30, // 挂件所在组
				Rank:       1,
			},
		}
	)
	Convey("when everything goes well", t, func() {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "PackageByID", func(_ *pendant.Dao, _ context.Context, _ int64, _ int64) (*model.PendantPackage, error) {
			return pendantPkg, nil
		})
		defer guard.Unpatch()
		res, err := s.PackageByID(context.Background(), mid, pid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})

	Convey("when occur an err", t, func() {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "PackageByID", func(_ *pendant.Dao, _ context.Context, _ int64, _ int64) (*model.PendantPackage, error) {
			return pendantPkg, nil
		})
		defer guard.Unpatch()
		res, err := s.PackageByID(context.Background(), mid, pid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

//func TestService_PendantInfoAA(t *testing.T) {
//	var (
//		pid int64 = 102
//	)
//	Convey("when everything goes well", t, func() {
//		guard := monkey.PatchInstanceMethod(reflect.TypeOf(s), "PendantInfo", func(_ *Service, _ context.Context, _ int64) (*model.Pendant, error) {
//			return nil, nil
//		})
//		defer guard.Unpatch()
//		res, err := s.PendantInfoAA(context.Background(), pid)
//		println(err == nil)
//		println(res == nil)
//	})
//}

func TestService_WearPendant(t *testing.T) {
	mid, vipPendant, vipPkg, pkgPendant, notVipPkg, vipInfo := preEquipPendant()
	Convey("test WearPendant", t, func() {
		Convey("when everything goes well", func() {
			// 穿戴挂件与来源一致 ：vip挂件（挂件在背包里面）
			Convey("1) test vip pendent", func() {
				// 1）查询挂件信息    //pendants,err := s.PendantInfo(context.Background(),pid)
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "PendantInfo", func(_ *Service, _ context.Context, _ int64) (*model.Pendant, error) {
					return vipPendant, nil
				})
				// 2）根据挂件ID查询挂件所在背包的信息
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "PackageByID", func(_ *pendant.Dao, _ context.Context, _ int64, _ int64) (*model.PendantPackage, error) {
					return vipPkg, nil
				})
				// 3) 穿戴
				// vip挂件需要查询vip信息
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "VipInfo", func(_ *pendant.Dao, _ context.Context, _ int64, _ string) (*model.VipInfo, error) {
					return vipInfo, nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquip", func(_ *pendant.Dao, _ context.Context, _ *model.PendantEquip) (int64, error) {
					return 1, nil
				})
				// 4）穿戴成功 删除缓存
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "DelEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64) error {
					return nil
				})
				// 5）重新添加缓存
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64, _ *model.PendantEquip) error {
					return nil
				})
				defer monkey.UnpatchAll()
				err := s.WearPendant(context.Background(), mid, 206, model.EquipFromVIP)
				Convey("pendant(vip) source(EquipFromVIP) ,then err should be nil", func() {
					So(err, ShouldBeNil)
				})
			})
			// 穿戴挂件与来源一致 ：vip挂件（挂件不在背包里面）
			Convey("2) test vip pendant and this pendant not exist in package", func() {
				// 1）查询挂件信息    //pendants,err := s.PendantInfo(context.Background(),pid)
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "PendantInfo", func(_ *Service, _ context.Context, _ int64) (*model.Pendant, error) {
					return vipPendant, nil
				})
				// 2）根据挂件ID查询挂件所在背包的信息
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "PackageByID", func(_ *pendant.Dao, _ context.Context, _ int64, _ int64) (*model.PendantPackage, error) {
					return nil, ecode.PendantPackageNotFound
				})
				// 3) 穿戴
				// vip挂件需要查询vip信息
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "VipInfo", func(_ *pendant.Dao, _ context.Context, _ int64, _ string) (*model.VipInfo, error) {
					return vipInfo, nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquip", func(_ *pendant.Dao, _ context.Context, _ *model.PendantEquip) (int64, error) {
					return 1, nil
				})
				// 4）穿戴成功 删除缓存
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "DelEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64) error {
					return nil
				})
				// 5）重新添加缓存
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64, _ *model.PendantEquip) error {
					return nil
				})
				defer monkey.UnpatchAll()
				err := s.WearPendant(context.Background(), mid, 0, model.EquipFromVIP)
				Convey("pendant(vip) source(EquipFromVIP), but allow it not exist in pkg, then err should be nil", func() {
					So(err, ShouldBeNil)
				})
			})

			// 穿戴挂件与来源一致 ：背包挂件(挂件必须在背包中)
			Convey("3) test pkg pendent", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "PendantInfo", func(_ *Service, _ context.Context, _ int64) (*model.Pendant, error) {
					return pkgPendant, nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "PackageByID", func(_ *pendant.Dao, _ context.Context, _ int64, _ int64) (*model.PendantPackage, error) {
					return notVipPkg, nil
				})

				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquip", func(_ *pendant.Dao, _ context.Context, _ *model.PendantEquip) (int64, error) {
					return 1, nil
				})

				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "DelEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64) error {
					return nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64, _ *model.PendantEquip) error {
					return nil
				})
				defer monkey.UnpatchAll()
				err := s.WearPendant(context.Background(), mid, 0, model.EquipFromPackage)
				Convey("pendant(pkg) source(EquipFromPackage) ,then err should be nil", func() {
					So(err, ShouldBeNil)
				})
			})

			// 来源未知 ：是背包挂件(挂件必须在背包中)
			Convey("4) test pkg pendent", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "PendantInfo", func(_ *Service, _ context.Context, _ int64) (*model.Pendant, error) {
					return pkgPendant, nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "PackageByID", func(_ *pendant.Dao, _ context.Context, _ int64, _ int64) (*model.PendantPackage, error) {
					return notVipPkg, nil
				})

				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquip", func(_ *pendant.Dao, _ context.Context, _ *model.PendantEquip) (int64, error) {
					return 1, nil
				})

				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "DelEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64) error {
					return nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64, _ *model.PendantEquip) error {
					return nil
				})
				defer monkey.UnpatchAll()
				err := s.WearPendant(context.Background(), mid, 0, model.UnknownEquipSource)
				Convey("pendant(pkg) source(UnknownEquipSource)->source(EquipFromPackage) ,then err should be nil", func() {
					So(err, ShouldBeNil)
				})
			})
			// 来源未知 ：是vip挂件(挂件在背包里面)
			Convey("5) test vip pendent", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "PendantInfo", func(_ *Service, _ context.Context, _ int64) (*model.Pendant, error) {
					return vipPendant, nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "PackageByID", func(_ *pendant.Dao, _ context.Context, _ int64, _ int64) (*model.PendantPackage, error) {
					return vipPkg, nil
				})
				// vip挂件需要查询vip信息
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "VipInfo", func(_ *pendant.Dao, _ context.Context, _ int64, _ string) (*model.VipInfo, error) {
					return vipInfo, nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquip", func(_ *pendant.Dao, _ context.Context, _ *model.PendantEquip) (int64, error) {
					return 1, nil
				})

				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "DelEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64) error {
					return nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64, _ *model.PendantEquip) error {
					return nil
				})
				defer monkey.UnpatchAll()
				err := s.WearPendant(context.Background(), mid, 0, model.UnknownEquipSource)
				Convey("pendant(vip) source(UnknownEquipSource)->source(EquipFromVIP) ,then err should be nil", func() {
					So(err, ShouldBeNil)
				})
			})
			// 来源未知 ：是vip挂件(挂件不在背包里面)
			Convey("6) test vip pendent", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "PendantInfo", func(_ *Service, _ context.Context, _ int64) (*model.Pendant, error) {
					return vipPendant, nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "PackageByID", func(_ *pendant.Dao, _ context.Context, _ int64, _ int64) (*model.PendantPackage, error) {
					return vipPkg, ecode.PendantPackageNotFound
				})
				// vip挂件需要查询vip信息
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "VipInfo", func(_ *pendant.Dao, _ context.Context, _ int64, _ string) (*model.VipInfo, error) {
					return vipInfo, nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquip", func(_ *pendant.Dao, _ context.Context, _ *model.PendantEquip) (int64, error) {
					return 1, nil
				})

				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "DelEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64) error {
					return nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64, _ *model.PendantEquip) error {
					return nil
				})
				defer monkey.UnpatchAll()
				err := s.WearPendant(context.Background(), mid, 0, model.UnknownEquipSource)
				Convey("pendant(vip) source(UnknownEquipSource)->source(EquipFromVIP), vip pendant allow not exist in package, then err should be nil", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("when occur an err", func() {
			// 穿戴一个背包挂件（非vip挂件），挂件来源是vip，报错
			Convey("1) test pkg pendent and source(EquipFromVIP)", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "PendantInfo", func(_ *Service, _ context.Context, _ int64) (*model.Pendant, error) {
					return pkgPendant, nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "PackageByID", func(_ *pendant.Dao, _ context.Context, _ int64, _ int64) (*model.PendantPackage, error) {
					return notVipPkg, nil
				})

				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquip", func(_ *pendant.Dao, _ context.Context, _ *model.PendantEquip) (int64, error) {
					return 1, nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "DelEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64) error {
					return nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64, _ *model.PendantEquip) error {
					return nil
				})
				defer monkey.UnpatchAll()
				err := s.WearPendant(context.Background(), mid, 0, model.EquipFromVIP)
				Convey("pendant(pkg) source(EquipFromPackage), then err should not be nil", func() {
					So(err, ShouldNotBeNil)
				})
			})

			// 穿戴一个背包挂件（非vip挂件），挂件来源是背包，但是挂件不在背包里，报错
			Convey("2) test pkg pendent and source(pkg), but not exist in package", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "PendantInfo", func(_ *Service, _ context.Context, _ int64) (*model.Pendant, error) {
					return pkgPendant, nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "PackageByID", func(_ *pendant.Dao, _ context.Context, _ int64, _ int64) (*model.PendantPackage, error) {
					return notVipPkg, ecode.PendantPackageNotFound
				})

				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquip", func(_ *pendant.Dao, _ context.Context, _ *model.PendantEquip) (int64, error) {
					return 1, nil
				})

				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "DelEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64) error {
					return nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64, _ *model.PendantEquip) error {
					return nil
				})
				defer monkey.UnpatchAll()
				err := s.WearPendant(context.Background(), mid, 0, model.EquipFromPackage)
				Convey("pendant(pkg) source(EquipFromPackage), but not exist in package, then err should not be nil", func() {
					So(err, ShouldNotBeNil)
				})
			})

			// 穿戴一个vip挂件，挂件来源是背包，但是挂件不在背包里，报错
			Convey("3) test vip pendent and source(pkg), but not exist in package", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "PendantInfo", func(_ *Service, _ context.Context, _ int64) (*model.Pendant, error) {
					return vipPendant, nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "PackageByID", func(_ *pendant.Dao, _ context.Context, _ int64, _ int64) (*model.PendantPackage, error) {
					return notVipPkg, ecode.PendantPackageNotFound
				})

				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquip", func(_ *pendant.Dao, _ context.Context, _ *model.PendantEquip) (int64, error) {
					return 1, nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "DelEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64) error {
					return nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64, _ *model.PendantEquip) error {
					return nil
				})
				defer monkey.UnpatchAll()
				err := s.WearPendant(context.Background(), mid, 0, model.EquipFromPackage)
				Convey("pendant(vip) source(EquipFromPackage), but not exist in package, then err should not be nil", func() {
					So(err, ShouldNotBeNil)
				})
			})

			// 来源未知 ：是背包挂件(挂件不在背包里面)，报错
			Convey("4) test pkg pendent and source(pkg), but not exist in package", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "PendantInfo", func(_ *Service, _ context.Context, _ int64) (*model.Pendant, error) {
					return pkgPendant, nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "PackageByID", func(_ *pendant.Dao, _ context.Context, _ int64, _ int64) (*model.PendantPackage, error) {
					return notVipPkg, ecode.PendantPackageNotFound
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquip", func(_ *pendant.Dao, _ context.Context, _ *model.PendantEquip) (int64, error) {
					return 1, nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "DelEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64) error {
					return nil
				})

				monkey.PatchInstanceMethod(reflect.TypeOf(s.pendantDao), "AddEquipCache", func(_ *pendant.Dao, _ context.Context, _ int64, _ *model.PendantEquip) error {
					return nil
				})
				defer monkey.UnpatchAll()
				err := s.WearPendant(context.Background(), mid, 0, model.UnknownEquipSource)
				Convey("pendant(pkg) source(UnknownEquipSource)->source(EquipFromPackage), but not exist in package, then err should not be nil", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}

func preEquipPendant() (mid int64, vipPendant *model.Pendant, vipPkg *model.PendantPackage, pkgPendant *model.Pendant, notVipPkg *model.PendantPackage, vipInfo *model.VipInfo) {
	mid = 111001965
	vipPendant = &model.Pendant{
		ID:     102,
		Name:   "文豪",
		Status: 1,  // 挂件状态:0.下线;1.在线
		Gid:    31, // 挂件所在组
		Rank:   1,
	}

	pkgPendant = &model.Pendant{
		ID:     98,
		Name:   "bilibili春",
		Status: 1,  // 挂件状态:0.下线;1.在线
		Gid:    30, // 挂件所在组
		Rank:   1,
	}

	vipPkg = &model.PendantPackage{
		ID:      16281,
		Mid:     111001965,
		Pid:     98,
		Expires: 1594076284,
		Type:    0,
		Status:  1,
		IsVIP:   0,
		Pendant: pkgPendant,
	}
	notVipPkg = &model.PendantPackage{
		ID:      16282,
		Mid:     111001965,
		Pid:     98,
		Expires: 1594076284,
		Type:    0,
		Status:  1,
		IsVIP:   0,
		Pendant: vipPendant,
	}
	vipInfo = &model.VipInfo{
		Mid:        111001965,
		VipType:    2,
		VipStatus:  1,
		VipDueDate: 1594076284,
	}
	return mid, vipPendant, vipPkg, pkgPendant, notVipPkg, vipInfo
}

//func TestService_EquipPendant(t *testing.T) {
//	var (
//		mid int64 = 111001965
//		pid int64 = 10
//		err error
//	)
//	Convey("need return something", func() {
//		err = s.EquipPendant(context.Background(), mid, pid, 1, 1)
//		So(err, ShouldNotBeNil)
//	})
//}

func TestService_EquipPendant(t *testing.T) {
	var (
		mid       int64 = 111001965
		vipPid    int64 = 102
		notVipPid int64 = 98
		wear      int8  = 2 //1：卸载 2：装备
		takeOff   int8  = 1 //1：卸载 2：装备
	)
	Convey("test EquipPendant", t, func() {
		Convey("when everything goes well", func() {
			Convey("1) takeOff pendant(vip) source(EquipFromVIP)", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "TakeOffPendant", func(_ *Service, _ context.Context, _ int64) error {
					return nil
				})
				defer monkey.UnpatchAll()
				err := s.EquipPendant(context.Background(), mid, vipPid, takeOff, model.EquipFromVIP)
				Convey("takeOff pendant(vip) source(EquipFromVIP),then err should be nil", func() {
					So(err, ShouldBeNil)
				})
			})
			Convey("2) takeOff pendant(pkg) source(EquipFromPackage)", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "TakeOffPendant", func(_ *Service, _ context.Context, _ int64) error {
					return nil
				})
				defer monkey.UnpatchAll()
				err := s.EquipPendant(context.Background(), mid, notVipPid, takeOff, model.EquipFromPackage)
				Convey("takeOff pendant(vip) source(EquipFromVIP),then err should be nil", func() {
					So(err, ShouldBeNil)
				})
			})

			// 穿戴一个vip挂件，挂件来源是vip，正确
			Convey("3) wear pendant(vip) source(EquipFromVIP)", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "WearPendant", func(_ *Service, _ context.Context, _ int64, _ int64, _ int64) error {
					return nil
				})
				defer monkey.UnpatchAll()
				err := s.EquipPendant(context.Background(), mid, vipPid, wear, model.EquipFromVIP)
				Convey("pendant(vip) source(EquipFromVIP), then err should be nil", func() {
					So(err, ShouldBeNil)
				})
			})

			// 穿戴一个vip挂件，挂件来源是背包（来源是背包的，挂件必须存在背包中），正确
			Convey("4) wear pendant(vip) source(EquipFromPackage)", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "WearPendant", func(_ *Service, _ context.Context, _ int64, _ int64, _ int64) error {
					return nil
				})
				defer monkey.UnpatchAll()
				err := s.EquipPendant(context.Background(), mid, vipPid, wear, model.EquipFromPackage)
				Convey("pendant(vip) source(EquipFromPackage), then err should be nil", func() {
					So(err, ShouldBeNil)
				})
			})

			// 穿戴一个背包挂件（非vip挂件），挂件来源是背包，正确
			Convey("5) wear pendant(pkg) source(EquipFromPackage)", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "WearPendant", func(_ *Service, _ context.Context, _ int64, _ int64, _ int64) error {
					return nil
				})
				defer monkey.UnpatchAll()
				err := s.EquipPendant(context.Background(), mid, notVipPid, wear, model.EquipFromPackage)
				Convey("pendant(pkg) source(EquipFromPackage), then err should be nil", func() {
					So(err, ShouldBeNil)
				})
			})

			// 穿戴一个vip挂件，挂件来源是 未知，正确
			Convey("6) wear pendant(vip) source(UnknownEquipSource)", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "WearPendant", func(_ *Service, _ context.Context, _ int64, _ int64, _ int64) error {
					return nil
				})
				defer monkey.UnpatchAll()
				err := s.EquipPendant(context.Background(), mid, vipPid, wear, model.UnknownEquipSource)
				Convey("pendant(vip) source(UnknownEquipSource)->source(EquipFromVIP), then err should be nil", func() {
					So(err, ShouldBeNil)
				})
			})

			// 穿戴一个背包挂件（存在背包里面），挂件来源是 未知，正确
			Convey("7) wear pendant(pkg) source(UnknownEquipSource)", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "WearPendant", func(_ *Service, _ context.Context, _ int64, _ int64, _ int64) error {
					return nil
				})
				defer monkey.UnpatchAll()
				err := s.EquipPendant(context.Background(), mid, notVipPid, wear, model.UnknownEquipSource)
				Convey("pendant(pkg) source(UnknownEquipSource)->source(EquipFromPackage), then err should be nil", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("wear pendant when occur an err", func() {
			// 穿戴一个背包挂件（非vip挂件），挂件来源是vip，报错
			Convey("1) wear pendant(pkg) source(EquipFromVIP)", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "WearPendant", func(_ *Service, _ context.Context, _ int64, _ int64, _ int64) error {
					return ecode.PendantPackageNotFound
				})
				defer monkey.UnpatchAll()
				err := s.EquipPendant(context.Background(), mid, notVipPid, wear, model.EquipFromVIP)
				Convey("pendant(pkg) source(EquipFromVIP),then err should not be nil", func() {
					So(err, ShouldNotBeNil)
				})
			})
			// 穿戴一个背包挂件（非vip挂件），挂件来源是背包，但是挂件不在背包里，报错
			Convey("2) wear pendant(pkg) source(EquipFromPackage)", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "WearPendant", func(_ *Service, _ context.Context, _ int64, _ int64, _ int64) error {
					return ecode.PendantPackageNotFound
				})
				defer monkey.UnpatchAll()
				err := s.EquipPendant(context.Background(), mid, notVipPid, wear, model.EquipFromPackage)
				Convey("pendant(pkg) source(EquipFromPackage),but it not exist in package, then err should not be nil", func() {
					So(err, ShouldNotBeNil)
				})
			})

			// 穿戴一个vip挂件，挂件来源是背包，但是挂件不在背包里，报错
			Convey("3) wear pendant(vip) source(EquipFromPackage)", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "WearPendant", func(_ *Service, _ context.Context, _ int64, _ int64, _ int64) error {
					return ecode.PendantPackageNotFound
				})
				defer monkey.UnpatchAll()
				err := s.EquipPendant(context.Background(), mid, vipPid, wear, model.EquipFromPackage)
				Convey("pendant(vip) source(EquipFromPackage), then err should not be nil", func() {
					So(err, ShouldNotBeNil)
				})
			})
			// 来源未知 ：是背包挂件(挂件不在背包里面)，报错
			Convey("4) wear pendant(pkg) source(UnknownEquipSource)", func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(s), "WearPendant", func(_ *Service, _ context.Context, _ int64, _ int64, _ int64) error {
					return ecode.PendantPackageNotFound
				})
				defer monkey.UnpatchAll()
				err := s.EquipPendant(context.Background(), mid, notVipPid, wear, model.UnknownEquipSource)
				Convey("pendant(pkg) source(UnknownEquipSource)->source(EquipFromPackage), but it not exist in package, then err should not be nil", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
	})

}
