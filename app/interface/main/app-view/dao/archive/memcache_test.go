package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_AddRelatesCache(t *testing.T) {
	Convey("AddRelatesCache", t, func() {
		d.AddRelatesCache(1, nil)
	})
}

func Test_RelatesCache(t *testing.T) {
	Convey("RelatesCache", t, func() {
		d.RelatesCache(context.TODO(), 1)
	})
}

func Test_ViewContributeCache(t *testing.T) {
	Convey("ViewContributeCache", t, func() {
		d.ViewContributeCache(context.TODO(), 2)
	})
}
