package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/passport-encrypt/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_PingMySQL(t *testing.T) {
	once.Do(startDao)
	if err := d.encryptDB.Ping(context.TODO()); err != nil {
		t.Errorf("dao.cloudDB.Ping() error(%v)", err)
		t.FailNow()
	}
}

func TestDao_UpdateAsoAccount(t *testing.T) {
	once.Do(startDao)
	Convey("update a aso account", t, func() {
		account := &model.EncryptAccount{
			Mid:            12047569,
			UserID:         "bili_1710676855",
			Uname:          "Bili_12047569",
			Pwd:            "3686c9d96ae6896fe117319ba6c07087",
			Salt:           "pdMXF856",
			Email:          "62fe0d616162f56ecab3e12a2de83ea6",
			Tel:            []byte("bdb27b0300e3984e48e7aea5c672a243"),
			CountryID:      1,
			MobileVerified: 1,
			Isleak:         0,
		}
		affected, err := d.UpdateAsoAccount(context.TODO(), account)
		So(err, ShouldBeNil)

		So(affected, ShouldEqual, 1)
	})
}
