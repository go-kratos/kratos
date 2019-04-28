package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_AuditTypesConf(t *testing.T) {
	Convey("AuditTypesConf", t, func() {
		configs, err := d.AuditTypesConf(context.TODO())
		So(err, ShouldBeNil)
		Println(configs)
	})
}
