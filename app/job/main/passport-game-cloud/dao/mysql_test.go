package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/passport-game-cloud/model"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_PingMySQL(t *testing.T) {
	once.Do(startDao)
	if err := d.cloudDB.Ping(context.TODO()); err != nil {
		t.Errorf("dao.cloudDB.Ping() error(%v)", err)
		t.FailNow()
	}
}

func TestDao_AddMemberInfo(t *testing.T) {
	once.Do(startDao)
	info := &model.Info{
		Mid:  110000130,
		Face: "/bfs/face/bbc031c4b7bdabb6635a246ce7386ccb587c5214811111.jpg",
	}
	if a, err := d.AddMemberInfo(context.TODO(), info); err != nil {
		t.FailNow()
	} else {
		t.Logf("a: %d", a)
	}
}

func TestDao_UpdateAsoAccount(t *testing.T) {
	once.Do(startDao)
	Convey("update a aso account", t, func() {
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
		affected, err := d.UpdateAsoAccount(context.TODO(), account)
		So(err, ShouldBeNil)

		So(affected, ShouldEqual, 1)
	})
}

func TestDao_AddAsoAccount(t *testing.T) {
	once.Do(startDao)
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
	if _, err := d.AddAsoAccount(context.TODO(), account); err != nil {
		switch nErr := errors.Cause(err).(type) {
		case *mysql.MySQLError:
			if nErr.Number != 1062 {
				t.Errorf("expected MySQL error 1062 but got error(%v)", err)
				t.FailNow()
			}
		}
	} else {
		t.FailNow()
	}
}
