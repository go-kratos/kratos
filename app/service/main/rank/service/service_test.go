package service

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"go-common/app/service/main/rank/conf"
	"go-common/app/service/main/rank/model"
	"go-common/library/log"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/test.toml")
	flag.Set("conf", dir)
	err := conf.Init()
	if err != nil {
		fmt.Printf("conf.Init() error(%v)", err)
	}
	s = New(conf.Conf)
}

func Test_dump(t *testing.T) {
	Convey("dump", t, func() {
		s.rmap = make(map[int][]*model.Field)
		for i := 1; i < 30; i++ {
			f := new(model.Field)
			f.Oid = int64(i)
			f.Pid = int16(i)
			f.Click = i * i
			f.Pubtime = xtime.Time(time.Now().Unix())
			s.setField(int64(i), f)
		}
		err := s.dump()
		t.Logf("err:%+v", err)
		So(err, ShouldBeNil)
	})
}

func Test_Marshal(t *testing.T) {
	Convey("Dump", t, func() {
		slic := make([]*model.Field, 10)
		for i := 1; i < 30; i++ {
			f := new(model.Field)
			f.Oid = int64(i)
			f.Pid = int16(i)
			f.Click = i * i
			f.Pubtime = xtime.Time(time.Now().Unix())
			slic = append(slic, f)
		}
		fields := new(model.Fields)
		fields.Fields = slic
		_, err := fields.Marshal()
		if err != nil {
			log.Error("fs.Marshal() error(%v)", err)
		}
		t.Logf("err:%+v", err)
		So(err, ShouldBeNil)
	})
}

func Test_field(t *testing.T) {
	Convey("field", t, func() {
		var oid int64 = 11
		f := &model.Field{
			Oid:     123,
			Pid:     22,
			Click:   33,
			Pubtime: 1551231231,
		}
		for i := 0; i < 30; i++ {
			f.Oid = int64(i)
			s.setField(int64(i), f)
		}
		ff := s.field(11)
		t.Logf("field:%+v", ff)
		So(s.rmap[s.bucket(oid)][s.mod(oid)], ShouldNotBeNil)
	})
}

func Test_setField(t *testing.T) {
	Convey("setField", t, func() {
		var oid int64 = 11
		f := &model.Field{
			Oid:     123,
			Pid:     22,
			Click:   33,
			Pubtime: 1551231231,
		}
		for i := 0; i < 30; i++ {
			f.Oid = int64(i)
			s.setField(int64(i), f)
		}
		t.Logf("field:%+v", s.rmap[s.bucket(oid)][s.mod(oid)])
		So(s.rmap[s.bucket(oid)][s.mod(oid)], ShouldNotBeNil)
	})
}

func Test_timestamp(t *testing.T) {
	Convey("timestamp", t, func() {
		now := time.Now().Unix()
		if err := ioutil.WriteFile(s.c.Rank.FilePath+"timestamp.txt", []byte(fmt.Sprintf("%d", now)), 0644); err != nil {
			log.Error("ioutil.WriteFile(%d) error(%v)", now, err)
		}
		fi, err := os.Open(s.c.Rank.FilePath + "timestamp.txt")
		if err != nil {
			log.Error(" os.Open(%s) error(%v)", s.c.Rank.FilePath+"timestamp.txt", err)
		}
		defer fi.Close()
		data, err := ioutil.ReadAll(fi)
		if err != nil {
			log.Error("ioutil.ReadAll() error(%v)", err)
		}
		begin, _ := strconv.ParseInt(string(data[:]), 10, 64)
		fmt.Println(string(data[:]), begin)
		So(begin, ShouldNotEqual, 0)
	})
}
