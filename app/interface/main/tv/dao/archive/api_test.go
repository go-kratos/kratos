package archive

import (
	"context"
	"fmt"
	"testing"

	"go-common/library/database/sql"

	. "github.com/smartystreets/goconvey/convey"
)

func getPassAid(db *sql.DB) (aid int64, err error) {
	if err = db.QueryRow(context.Background(), "select aid from ugc_archive where deleted = 0 and valid = 1 and result = 1 limit 1").Scan(&aid); err != nil {
		fmt.Println(err)
	}
	return
}

func TestDao_RelateAids(t *testing.T) {
	Convey("TestDao_RelateAids", t, WithDao(func(d *Dao) {
		httpMock("GET", d.relateURL).Reply(200).JSON(`{"data":[{"key": "1234", "value": "23999,515522,55,10099960,10099978,10099959,10099959,10099959,10099959,10099959,10099959,10099959,10099959,10099959,10099959,10099959,10099959,10099959,10099959,10099959,10099959,10099959"}]}`)
		aids, err := d.RelateAids(context.Background(), 333444, "")
		So(err, ShouldBeNil)
		So(len(aids), ShouldBeGreaterThan, 0)
		fmt.Println(aids)
	}))
}
