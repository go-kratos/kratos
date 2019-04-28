package dao

import (
	"context"
	"math/rand"
	"testing"

	"go-common/app/admin/main/usersuit/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_AddPendantGroup(t *testing.T) {
	Convey("return someting", t, func() {
		pg := &model.PendantGroup{
			Name: "dasdasd",
			Rank: 22,
		}
		gid, err := d.AddPendantGroup(context.Background(), pg)
		So(err, ShouldBeNil)
		So(gid, ShouldNotBeNil)
	})
}

func Test_TxAddPendantGroupRef(t *testing.T) {
	Convey("return someting", t, func() {
		pr := &model.PendantGroupRef{
			GID: 11,
			PID: int64(rand.Int31()),
		}
		tx, err := d.BeginTran(context.Background())
		So(err, ShouldBeNil)
		effect, err := d.TxAddPendantGroupRef(tx, pr)
		So(err, ShouldBeNil)
		defer func() {
			if err != nil || effect == 0 {
				tx.Rollback()
				return
			}
			tx.Commit()
		}()
		So(err, ShouldBeNil)
	})
}

func Test_TxAddPendantInfo(t *testing.T) {
	Convey("return someting", t, func() {
		pi := &model.PendantInfo{
			Name:       "dasdasdsads",
			Image:      "dasdds",
			ImageModel: "xxsss",
			Rank:       11,
		}
		tx, err := d.BeginTran(context.Background())
		if err != nil {
			So(err, ShouldBeNil)
		}
		id, err := d.TxAddPendantInfo(tx, pi)
		if err != nil {
			So(err, ShouldBeNil)
		}
		defer func() {
			if err != nil || id == 0 {
				tx.Rollback()
				return
			}
			tx.Commit()
		}()
		So(err, ShouldBeNil)
	})
}

func Test_TxAddPendantPrices(t *testing.T) {
	Convey("return someting", t, func() {
		pp := &model.PendantPrice{
			PID:   22,
			TP:    1,
			Price: 22,
		}
		tx, err := d.BeginTran(context.Background())
		if err != nil {
			So(err, ShouldBeNil)
		}
		effect, err := d.TxAddPendantPrices(tx, pp)
		if err != nil {
			So(err, ShouldBeNil)
		}
		defer func() {
			if err != nil || effect == 0 {
				tx.Rollback()
				return
			}
			tx.Commit()
		}()
		So(err, ShouldBeNil)
	})
}

func Test_AddPendantPKG(t *testing.T) {
	Convey("return someting", t, func() {
		pkg := &model.PendantPKG{
			UID:     int64(rand.Int31()),
			PID:     11,
			Expires: 12312323,
		}
		_, err := d.AddPendantPKG(context.Background(), pkg)
		So(err, ShouldBeNil)
	})
}

func Test_TxAddPendantPKGs(t *testing.T) {
	uid := int64(rand.Int31())
	Convey("return someting", t, func() {
		var pkgs []*model.PendantPKG
		pkgs = append(pkgs, &model.PendantPKG{UID: uid, PID: 11, Expires: 12312323}, &model.PendantPKG{UID: uid, PID: 22, Expires: 12312323})
		tx, err := d.BeginTran(context.Background())
		if err != nil {
			So(err, ShouldBeNil)
		}
		effect, err := d.TxAddPendantPKGs(tx, pkgs)
		if err != nil {
			So(err, ShouldBeNil)
		}
		defer func() {
			if err != nil || effect == 0 {
				tx.Rollback()
				return
			}
			tx.Commit()
		}()
		So(err, ShouldBeNil)
	})
}

func Test_AddPendantEquip(t *testing.T) {
	Convey("return someting", t, func() {
		pkg := &model.PendantPKG{
			UID:     22,
			PID:     11,
			Expires: 12312323,
		}
		_, err := d.AddPendantEquip(context.Background(), pkg)
		So(err, ShouldBeNil)
	})
}

func Test_AddPendantOperLog(t *testing.T) {
	Convey("return someting", t, func() {
		_, err := d.AddPendantOperLog(context.Background(), 1, []int64{1}, 1, "sdsadasd")
		So(err, ShouldBeNil)
	})
}

func Test_TxUpPendantGroupRef(t *testing.T) {
	Convey("return someting", t, func() {
		tx, err := d.BeginTran(context.Background())
		if err != nil {
			So(err, ShouldBeNil)
		}
		effect, err := d.TxUpPendantGroupRef(tx, 22, 11)
		if err != nil {
			So(err, ShouldBeNil)
		}
		defer func() {
			if err != nil || effect == 0 {
				tx.Rollback()
				return
			}
			tx.Commit()
		}()
		So(err, ShouldBeNil)
	})
}

func Test_TxUpPendantPKGs(t *testing.T) {
	Convey("return someting", t, func() {
		var pkgs []*model.PendantPKG
		pkgs = append(pkgs, &model.PendantPKG{UID: 22, PID: 11, Expires: 12312323}, &model.PendantPKG{UID: 11, PID: 22, Expires: 12312323})
		tx, err := d.BeginTran(context.Background())
		if err != nil {
			So(err, ShouldBeNil)
		}
		effect, err := d.TxUpPendantPKGs(tx, pkgs)
		if err != nil {
			So(err, ShouldBeNil)
		}
		defer func() {
			if err != nil || effect == 0 {
				tx.Rollback()
				return
			}
			tx.Commit()
		}()
		So(err, ShouldBeNil)
	})
}

func Test_UpPendantGroup(t *testing.T) {
	Convey("return someting", t, func() {
		pg := &model.PendantGroup{
			Name:   "weqweqw",
			Rank:   2,
			Status: 1,
			ID:     22,
		}
		_, err := d.UpPendantGroup(context.Background(), pg)
		So(err, ShouldBeNil)
	})
}

func Test_UpPendantGroupStatus(t *testing.T) {
	Convey("return someting", t, func() {
		_, err := d.UpPendantGroupStatus(context.Background(), 22, 1)
		So(err, ShouldBeNil)
	})
}

func Test_TxUpPendantInfo(t *testing.T) {
	Convey("return someting", t, func() {
		pi := &model.PendantInfo{
			Name:       "dasdasdsads",
			Image:      "dasdds",
			ImageModel: "xxsss",
			Rank:       11,
			ID:         22,
		}
		tx, err := d.BeginTran(context.Background())
		if err != nil {
			So(err, ShouldBeNil)
		}
		effect, err := d.TxUpPendantInfo(tx, pi)
		if err != nil {
			So(err, ShouldBeNil)
		}
		defer func() {
			if err != nil || effect == 0 {
				tx.Rollback()
				return
			}
			tx.Commit()
		}()
		So(err, ShouldBeNil)
	})
}

func Test_UpPendantInfoStatus(t *testing.T) {
	Convey("return someting", t, func() {
		_, err := d.UpPendantInfoStatus(context.Background(), 22, 1)
		So(err, ShouldBeNil)
	})
}

func Test_PendantInfoAll(t *testing.T) {
	Convey("return someting", t, func() {
		_, _, err := d.PendantInfoAll(context.Background(), 1, 2)
		So(err, ShouldBeNil)
	})
}

func Test_PendantGroupIDs(t *testing.T) {
	Convey("return someting", t, func() {
		_, err := d.PendantGroupIDs(context.Background(), []int64{11, 22})
		So(err, ShouldBeNil)
	})
}

func Test_PendantGroupID(t *testing.T) {
	Convey("return someting", t, func() {
		_, err := d.PendantGroupID(context.Background(), 12)
		So(err, ShouldBeNil)
	})
}

func Test_PendantInfoIDs(t *testing.T) {
	Convey("return someting", t, func() {
		_, _, err := d.PendantInfoIDs(context.Background(), []int64{11, 22})
		So(err, ShouldBeNil)
	})
}

func Test_PendantPriceIDs(t *testing.T) {
	Convey("return someting", t, func() {
		ppm, err := d.PendantPriceIDs(context.Background(), []int64{11, 22})
		So(err, ShouldBeNil)
		So(ppm, ShouldNotBeNil)
	})
}

func Test_PendantGroupRefRanks(t *testing.T) {
	Convey("return someting", t, func() {
		_, err := d.PendantGroupRefRanks(context.Background(), 1, 2)
		So(err, ShouldBeNil)
	})
}

func Test_PendantGroupPIDs(t *testing.T) {
	Convey("return someting", t, func() {
		_, err := d.PendantGroupPIDs(context.Background(), 11, 1, 2)
		So(err, ShouldBeNil)
	})
}

func Test_PendantInfoID(t *testing.T) {
	Convey("return someting", t, func() {
		pi, err := d.PendantInfoID(context.Background(), 11)
		So(err, ShouldBeNil)
		So(pi, ShouldNotBeNil)
	})
}

func Test_PendantInfoAllOnSale(t *testing.T) {
	Convey("return someting", t, func() {
		_, err := d.PendantInfoAllNoPage(context.Background())
		So(err, ShouldBeNil)
	})
}

func Test_CountOrderHistory(t *testing.T) {
	Convey("return someting", t, func() {
		arg := &model.ArgPendantOrder{}
		_, err := d.CountOrderHistory(context.Background(), arg)
		So(err, ShouldBeNil)
	})
}

func Test_OrderHistorys(t *testing.T) {
	Convey("return someting", t, func() {
		arg := &model.ArgPendantOrder{}
		_, _, err := d.OrderHistorys(context.Background(), arg)
		So(err, ShouldBeNil)
	})
}

func Test_PendantPKGs(t *testing.T) {
	Convey("return someting", t, func() {
		_, err := d.PendantPKGs(context.Background(), 112)
		So(err, ShouldBeNil)
	})
}

func Test_PendantPKGUIDs(t *testing.T) {
	Convey("return someting", t, func() {
		_, err := d.PendantPKGUIDs(context.Background(), []int64{11, 22}, 112)
		So(err, ShouldBeNil)
	})
}

func Test_PendantPKG(t *testing.T) {
	Convey("return someting", t, func() {
		_, err := d.PendantPKG(context.Background(), 11, 112)
		So(err, ShouldBeNil)
	})
}

func Test_PendantEquipUID(t *testing.T) {
	Convey("return someting", t, func() {
		_, err := d.PendantEquipUID(context.Background(), 11)
		So(err, ShouldBeNil)
	})
}

func Test_PendantOperLog(t *testing.T) {
	Convey("return someting", t, func() {
		_, _, err := d.PendantOperLog(context.Background(), 1, 2)
		So(err, ShouldBeNil)
	})
}

func Test_PendantOperationLogTotal(t *testing.T) {
	Convey("return someting", t, func() {
		total, err := d.PendantOperationLogTotal(context.Background())
		So(err, ShouldBeNil)
		So(total, ShouldNotBeNil)
	})
}
