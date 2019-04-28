package archive

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDao_GetUserCard(t *testing.T) {
	Convey("GetUserCard", t, WithDao(func(d *Dao) {
		httpMock("GET", d.userCardURL).Reply(200).JSON(`{"code":0,"data":{"has":100}}`)
		card, err := d.GetUserCard(context.Background(), 27515615)
		t.Logf("GetUserCard error(%v)\r\n card(%+v)", err, card)
		So(err, ShouldBeNil)
	}))
}
