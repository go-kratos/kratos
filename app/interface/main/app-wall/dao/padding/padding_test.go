package padding

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	p *pkcs5
)

func init() {
	p = &pkcs5{}
	time.Sleep(time.Second)
}

func TestPadding(t *testing.T) {
	Convey("Padding", t, func() {
		res := p.Padding([]byte{}, 1)
		So(res, ShouldNotBeEmpty)
	})
}
