package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_View3(t *testing.T) {
	Convey("View3", t, func() {
		d.View3(context.TODO(), 1)
	})
}

func Test_ViewCache(t *testing.T) {
	Convey("ViewCache", t, func() {
		d.ViewCache(context.TODO(), 1)
	})
}

func Test_Description(t *testing.T) {
	Convey("Description", t, func() {
		d.Description(context.TODO(), 2)
	})
}
