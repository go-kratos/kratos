package ai

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	ai  *AI
	ctx = context.TODO()
)

func TestMain(m *testing.M) {
	var err error
	ai = New()
	var loader = func(context.Context) (m map[int64]int64, err error) {
		m = map[int64]int64{
			46333: 1,
			35858: 1,
		}
		return
	}
	if err = ai.LoadWhite(ctx, loader); err != nil {
		panic(err)
	}

	m.Run()
}

func TestAI(t *testing.T) {
	Convey("ai", t, func() {
		i, err := ai.White(46333)
		So(i, ShouldBeGreaterThan, 0)
		So(err, ShouldBeNil)
	})
}
