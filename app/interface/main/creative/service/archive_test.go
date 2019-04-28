package service

import (
	"testing"

	"context"
	"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_BGM(t *testing.T) {
	Convey("bgm", t, WithService(func(p *Public) {
		Convey("bgm", func() {
			res, err := p.BgmBindList(context.TODO(), 10110983, 10135130, 3, false)
			if err == nil {
				spew.Dump(res)
			}
			spew.Dump("bgm")
		})
	}))
}

func Test_Type(t *testing.T) {
	Convey("Type", t, WithService(func(p *Public) {
		Convey("loadTypes", func() {
			p.loadTypes()
			spew.Dump(p.TopTypesCache)
		})
		// Convey("loadDescFormat", func() {
		// 	p.loadDescFormat()
		// 	spew.Dump(p.DescFmtsCache)
		// })
	}))
}

func Test_Staff(t *testing.T) {
	Convey("Staff", t, WithService(func(p *Public) {
		Convey("Staff", func() {
			res, err := p.StaffList(context.TODO(), 10110983, false)
			if err == nil {
				spew.Dump(res)
			}
			spew.Dump("Staff")
		})
	}))
}
