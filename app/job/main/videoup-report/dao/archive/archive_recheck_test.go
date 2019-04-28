package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/videoup-report/model/archive"
)

func TestDao_TxAddRecheckAID(t *testing.T) {
	Convey("TxAddRecheckAID", t, func() {
		c := context.TODO()
		a, _ := d.ArchiveByAid(c, 1)

		tx, _ := d.BeginTran(c)
		id, err := d.TxAddRecheckAID(tx, archive.TypeChannelRecheck, a.ID)
		tx.Commit()
		So(err, ShouldBeNil)
		Println(id)
	})
}
