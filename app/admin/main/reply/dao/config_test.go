package dao

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go-common/app/admin/main/reply/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAddConfig(t *testing.T) {
	var (
		c      = context.Background()
		now    = time.Now()
		config = &model.Config{
			Oid:      1,
			Type:     1,
			Category: 1,
			AdminID:  1,
			Operator: "admin",
		}
	)
	Convey("add a config", t, WithDao(func(d *Dao) {
		configValue := map[string]int64{
			"showentry": 0,
			"showadmin": 1,
		}
		bs, err := json.Marshal(configValue)
		So(err, ShouldBeNil)
		config.Config = string(bs)
		_, err = d.AddConfig(c, config.Type, config.Category, config.Oid, config.AdminID, config.Operator, config.Config, now)
		So(err, ShouldBeNil)
	}))
}

func TestLoadConfig(t *testing.T) {
	var (
		c      = context.Background()
		config = &model.Config{
			Oid:      1,
			Type:     1,
			Category: 1,
			AdminID:  1,
			Operator: "admin",
		}
	)
	Convey("load a config", t, WithDao(func(d *Dao) {
		var err error
		config, err = d.LoadConfig(c, config.Type, config.Category, config.Oid)
		So(err, ShouldBeNil)
		So(config, ShouldNotBeNil)
	}))
}

func TestPaginateConfig(t *testing.T) {
	var (
		config = &model.Config{
			Oid:      1,
			Type:     1,
			Category: 1,
			AdminID:  1,
			Operator: "admin",
		}
		c = context.Background()
	)
	Convey("load a config", t, WithDao(func(d *Dao) {
		configs, err := d.PaginateConfig(c, config.Type, config.Category, config.Oid, config.Operator, 0, 20)
		So(err, ShouldBeNil)
		So(len(configs), ShouldNotEqual, 0)
	}))
}

func TestDeleteConfig(t *testing.T) {
	var (
		id = int64(1)
		c  = context.Background()
	)
	Convey("load a config", t, WithDao(func(d *Dao) {
		_, err := d.DeleteConfig(c, id)
		So(err, ShouldBeNil)
	}))
}
