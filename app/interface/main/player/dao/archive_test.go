package dao

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_PvData(t *testing.T) {
	convey.Convey("test pvdata", t, func(ctx convey.C) {
		url := "http://i3.hdslb.com/bfs/videoshot/10135146.bin?vsign=5d5c80b9c583b1ce49f0fb8ab8dd178568efa66c&ver=108653418"
		res, err := d.PvData(context.Background(), url)
		ctx.So(err, convey.ShouldBeNil)
		var (
			v   uint16
			pvs []uint16
			buf = bytes.NewReader(res)
		)
		for {
			if err = binary.Read(buf, binary.BigEndian, &v); err != nil {
				if err != io.EOF {
					ctx.Printf("binary.Read err(%v)", err)
				}
				err = nil
				break
			}
			pvs = append(pvs, v)
		}
		ctx.So(err, convey.ShouldBeNil)
		ctx.Printf("%+v", pvs)
	})
}
