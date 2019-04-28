package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go-common/app/job/main/passport-game-data/model"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	_timeFormat = "2006-01-02 15:04:05"
)

var (
	_loc = time.Now().Location()
)

func TestDao_AddAsoAccountsCloud(t *testing.T) {
	once.Do(startDao)

	Convey("batch add aso account to cloud", t, func() {
		Convey("single", func() {
			as := make([]*model.AsoAccount, 0)

			a := &model.AsoAccount{
				Mid:            12047569,
				UserID:         "bili_1710676855",
				Uname:          "Bili_12047569",
				Pwd:            "3686c9d96ae6896fe117319ba6c07087",
				Salt:           "pdMXF856",
				Email:          "62fe0d616162f56ecab3e12a2de83ea6",
				Tel:            "bdb27b0300e3984e48e7aea5c672a243",
				CountryID:      1,
				MobileVerified: 1,
				Isleak:         0,
			}

			as = append(as, a)
			err := d.AddAsoAccountsCloud(context.TODO(), as)

			So(err, ShouldBeNil)
		})

		Convey("multiple", func() {
			as := make([]*model.AsoAccount, 0)

			a := &model.AsoAccount{
				Mid:            12047569,
				UserID:         "bili_1710676855",
				Uname:          "Bili_12047569",
				Pwd:            "3686c9d96ae6896fe117319ba6c07087",
				Salt:           "pdMXF856",
				Email:          "62fe0d616162f56ecab3e12a2de83ea6",
				Tel:            "bdb27b0300e3984e48e7aea5c672a243",
				CountryID:      1,
				MobileVerified: 1,
				Isleak:         0,
			}

			as = append(as, a)
			as = append(as, a)
			err := d.AddAsoAccountsCloud(context.TODO(), as)

			So(err, ShouldBeNil)
		})
	})
}

func TestDao_AsoAccountRangeCloud(t *testing.T) {
	once.Do(startDao)

	Convey("get a aso account from cloud range start time and end time", t, func() {
		Convey("when start time after end time", func() {
			st, err := time.ParseInLocation(_timeFormat, "2018-01-22 12:50:41", _loc)
			ed, err := time.ParseInLocation(_timeFormat, "2018-01-22 12:49:41", _loc)
			So(err, ShouldBeNil)

			res, err := d.AsoAccountRangeCloud(context.TODO(), st, ed)
			So(err, ShouldBeNil)
			So(len(res), ShouldEqual, 0)
		})

		Convey("when res is not empty", func() {
			st, err := time.ParseInLocation(_timeFormat, "2018-01-22 12:50:41", _loc)
			ed, err := time.ParseInLocation(_timeFormat, "2018-01-22 12:51:41", _loc)
			So(err, ShouldBeNil)

			res, err := d.AsoAccountRangeCloud(context.TODO(), st, ed)
			So(err, ShouldBeNil)
			So(len(res), ShouldBeGreaterThan, 0)

			mid := int64(88888970)
			ok := false
			var target *model.AsoAccount
			for _, a := range res {
				if a.Mid == mid {
					target = a
					ok = true
					break
				}
			}
			So(ok, ShouldBeTrue)
			So(target.Email, ShouldNotBeNil)
			So(target.Tel, ShouldBeEmpty)

			str, _ := json.Marshal(target)
			t.Logf("res: %s", str)
		})
	})
}

func TestDao_AsoAccountRangeLocal(t *testing.T) {
	once.Do(startDao)

	Convey("get a aso account from local range start time and end time", t, func() {
		Convey("when start time is after end time", func() {
			st, err := time.ParseInLocation(_timeFormat, "2018-01-22 12:50:41", _loc)
			ed, err := time.ParseInLocation(_timeFormat, "2018-01-22 12:49:41", _loc)
			So(err, ShouldBeNil)

			res, err := d.AsoAccountRangeLocal(context.TODO(), st, ed)
			So(err, ShouldBeNil)
			So(len(res), ShouldEqual, 0)
		})

		Convey("when res is not empty", func() {
			st, err := time.ParseInLocation(_timeFormat, "2018-01-22 12:50:41", _loc)
			ed, err := time.ParseInLocation(_timeFormat, "2018-01-22 12:51:41", _loc)
			So(err, ShouldBeNil)

			res, err := d.AsoAccountRangeLocal(context.TODO(), st, ed)
			So(err, ShouldBeNil)
			So(len(res), ShouldBeGreaterThan, 0)

			mid := int64(88888970)
			ok := false
			var target *model.OriginAsoAccount
			for _, a := range res {
				if a.Mid == mid {
					target = a
					ok = true
					break
				}
			}
			So(ok, ShouldBeTrue)
			So(target.Email, ShouldNotBeNil)
			So(target.Tel, ShouldNotBeEmpty)

			str, _ := json.Marshal(target)
			t.Logf("res: %s", str)
		})
	})
}

func TestDao_AsoAccountsCloud(t *testing.T) {
	once.Do(startDao)

	Convey("get aso accounts from cloud", t, func() {
		Convey("when res is empty", func() {
			res, err := d.AsoAccountsCloud(context.TODO(), []int64{10000000000})
			So(err, ShouldBeNil)
			So(len(res), ShouldEqual, 0)
		})

		Convey("when res has single item", func() {
			mids := []int64{88888970}
			res, err := d.AsoAccountsCloud(context.TODO(), mids)
			So(err, ShouldBeNil)
			So(len(res), ShouldEqual, 1)

			m := make(map[int64]*model.AsoAccount)
			for _, a := range res {
				m[a.Mid] = a
			}

			for _, mid := range mids {
				a, ok := m[mid]
				So(ok, ShouldBeTrue)
				So(a.Email, ShouldBeEmpty)
				So(a.Tel, ShouldNotBeEmpty)

				str, _ := json.Marshal(a)
				t.Logf("a: %s", str)
			}
		})

		Convey("when res has multiple items", func() {
			mids := []int64{88888970, 110000784}
			res, err := d.AsoAccountsCloud(context.TODO(), mids)
			So(err, ShouldBeNil)
			So(len(res), ShouldEqual, 2)

			m := make(map[int64]*model.AsoAccount)
			for _, a := range res {
				m[a.Mid] = a
			}

			for _, mid := range mids {
				a, ok := m[mid]
				So(ok, ShouldBeTrue)
				So(a.Email, ShouldBeEmpty)
				So(a.Tel, ShouldNotBeEmpty)

				str, _ := json.Marshal(a)
				t.Logf("a: %s", str)
			}
		})
	})
}

func TestDao_AsoAccountsLocal(t *testing.T) {
	once.Do(startDao)

	Convey("get aso accounts from local", t, func() {
		Convey("when res is empty", func() {
			res, err := d.AsoAccountsCloud(context.TODO(), []int64{10000000000})
			So(err, ShouldBeNil)
			So(len(res), ShouldEqual, 0)
		})

		Convey("when res has single item", func() {
			mids := []int64{88888970}
			res, err := d.AsoAccountsCloud(context.TODO(), mids)
			So(err, ShouldBeNil)
			So(len(res), ShouldEqual, 1)

			m := make(map[int64]*model.AsoAccount)
			for _, a := range res {
				m[a.Mid] = a
			}

			for _, mid := range mids {
				a, ok := m[mid]
				So(ok, ShouldBeTrue)
				So(a.Email, ShouldNotBeNil)
				So(a.Tel, ShouldBeEmpty)

				str, _ := json.Marshal(a)
				t.Logf("a: %s", str)
			}
		})

		Convey("when res has multiple items", func() {
			mids := []int64{88888970, 110000784}
			res, err := d.AsoAccountsCloud(context.TODO(), mids)
			So(err, ShouldBeNil)
			So(len(res), ShouldEqual, 2)

			m := make(map[int64]*model.AsoAccount)
			for _, a := range res {
				m[a.Mid] = a
			}

			for _, mid := range mids {
				a, ok := m[mid]
				So(ok, ShouldBeTrue)
				So(a.Email, ShouldNotBeNil)
				So(a.Tel, ShouldBeEmpty)

				str, _ := json.Marshal(a)
				t.Logf("a: %s", str)
			}
		})
	})
}

func TestDao_UpdateAsoAccountCloud(t *testing.T) {
	once.Do(startDao)
	Convey("update aso account", t, func() {
		Convey("when mtime matches", func() {
			mid := int64(12047569)

			as, err := d.AsoAccountsCloud(context.TODO(), []int64{mid})
			So(err, ShouldBeNil)
			So(len(as), ShouldEqual, 1)

			account := as[0]
			account.MobileVerified = 1 - account.MobileVerified

			affected, err := d.UpdateAsoAccountCloud(context.TODO(), account, account.Mtime)
			So(err, ShouldBeNil)
			So(affected, ShouldEqual, 1)
		})

		Convey("when mtime not matches", func() {
			mid := int64(12047569)

			as, err := d.AsoAccountsCloud(context.TODO(), []int64{mid})
			So(err, ShouldBeNil)
			So(len(as), ShouldEqual, 1)

			account := as[0]
			account.MobileVerified = 1 - account.MobileVerified

			affected, err := d.UpdateAsoAccountCloud(context.TODO(), account, time.Now())
			So(err, ShouldBeNil)
			So(affected, ShouldEqual, 0)
		})
	})
}

func TestDao_AddIgnoreAsoAccount(t *testing.T) {
	once.Do(startDao)
	Convey("add ignore a aso account when not exist", t, func() {
		//Convey("when not exists", func() {
		//	account := &model.AsoAccount{
		//		Mid:            12047569,
		//		UserID:         "bili_1710676855",
		//		Uname:          "Bili_12047569",
		//		Pwd:            "3686c9d96ae6896fe117319ba6c07087",
		//		Salt:           "pdMXF856",
		//		Email:          "62fe0d616162f56ecab3e12a2de83ea6",
		//		Tel:            "bdb27b0300e3984e48e7aea5c672a243",
		//		CountryID:      1,
		//		MobileVerified: 1,
		//		Isleak:         0,
		//	}
		//	affected, err := d.AddIgnoreAsoAccount(context.TODO(), account)
		//	So(err, ShouldBeNil)
		//	So(affected, ShouldEqual, 1)
		//})
		Convey("when not exists", func() {
			account := &model.AsoAccount{
				Mid:            12047569,
				UserID:         "bili_1710676855",
				Uname:          "Bili_12047569",
				Pwd:            "3686c9d96ae6896fe117319ba6c07087",
				Salt:           "pdMXF856",
				Email:          "62fe0d616162f56ecab3e12a2de83ea6",
				Tel:            "bdb27b0300e3984e48e7aea5c672a243",
				CountryID:      1,
				MobileVerified: 1,
				Isleak:         0,
			}
			affected, err := d.AddIgnoreAsoAccount(context.TODO(), account)
			So(err, ShouldBeNil)
			So(affected, ShouldEqual, 0)
		})
	})
}

const (
	_pattern = "INSERT INTO aso_account (mid,userid,uname,pwd,salt,email,tel,country_id,mobile_verified,isleak) VALUES(%d,'%s','%s','%s','%s',%s,%s,%d,%d,%d) ON DUPLICATE KEY UPDATE userid='%s',uname='%s',pwd='%s',salt='%s',email=%s,tel=%s,country_id=%d,mobile_verified=%d,isleak=%d;"
)

func TestDao_AddAsoAccount(t *testing.T) {
	oldStr := `{
        "mid": 255554277,
        "userid": "bili_93079136999",
        "uname": "白又寻",
        "pwd": "8489c2cbddb7ee1438698a4f21ee1d78",
        "salt": "r0MHcs5M",
        "email": "",
        "tel": "ca6d0469ca340f67f4635425dcd11581",
        "country_id": 1,
        "mobile_verified": 2,
        "isleak": 0,
        "ctime": "2017-11-25T12:03:38+08:00",
        "mtime": "2017-12-04T14:19:59+08:00"
      }`

	old := new(model.AsoAccount)
	err := json.Unmarshal([]byte(oldStr), &old)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	sql := getSQL(old)
	t.Logf("sql: %s", sql)

	afterStr := `{
        "mid": 255554277,
        "userid": "bili_93079136999",
        "uname": "白又寻",
        "pwd": "5f064b2ddb4d8cd5f9e01507ab1d34c6",
        "salt": "ggr58PEs",
        "email": "",
        "tel": "ca6d0469ca340f67f4635425dcd11581",
        "country_id": 1,
        "mobile_verified": 2,
        "isleak": 0,
        "ctime": "0001-01-01T00:00:00Z",
        "mtime": "2017-11-25T19:14:20+08:00"
      }`

	after := new(model.AsoAccount)

	err = json.Unmarshal([]byte(afterStr), &after)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	afterSQL := getSQL(after)
	t.Logf("after sql: %s", afterSQL)

}

func getSQL(a *model.AsoAccount) string {
	email := "NULL"
	tel := "NULL"
	if len(a.Email) > 0 {
		email = "'" + a.Email + "'"
	}

	if len(a.Tel) > 0 {
		tel = "'" + a.Tel + "'"
	}

	return fmt.Sprintf(_pattern, a.Mid, a.UserID, a.Uname, a.Pwd, a.Salt, email, tel, a.CountryID, a.MobileVerified, a.Isleak, a.UserID, a.Uname, a.Pwd, a.Salt, email, tel, a.CountryID, a.MobileVerified, a.Isleak)

}
