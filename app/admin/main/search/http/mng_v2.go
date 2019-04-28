package http

import (
	"go-common/app/admin/main/search/model"
	bm "go-common/library/net/http/blademaster"
)

func businessAllV2(c *bm.Context) {
	c.JSON(svr.BusinessAllV2(c))
}

func businessInfoV2(c *bm.Context) {
	p := new(struct {
		Name string `form:"name" validate:"required"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(svr.BusinessInfoV2(c, p.Name))
}

func businessAdd(c *bm.Context) {
	p := new(struct {
		Pid         int64  `form:"pid" validate:"required,min=1"`
		Name        string `form:"name" validate:"required"`
		Description string `form:"description" validate:"required"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(svr.BusinessAdd(c, p.Pid, p.Name, p.Description))
}

func businessUpdate(c *bm.Context) {
	p := new(struct {
		Name  string `form:"name" validate:"required"`
		Field string `form:"field" validate:"required"`
		Value string `form:"value"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(svr.BusinessUpdate(c, p.Name, p.Field, p.Value))
}

func assetDBTables(c *bm.Context) {
	c.JSON(svr.AssetDBTables(c))
}

func assetDBConnect(c *bm.Context) {
	p := new(struct {
		Host     string `form:"host" validate:"required"`
		Port     string `form:"port" validate:"required"`
		User     string `form:"user" validate:"required"`
		Password string `form:"password" validate:"required"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(svr.AssetDBConnect(c, p.Host, p.Port, p.User, p.Password))
}

func assetDBAdd(c *bm.Context) {
	p := new(struct {
		Name        string `form:"name" validate:"required"`
		Description string `form:"description"`
		Host        string `form:"host" validate:"required"`
		Port        string `form:"port" validate:"required"`
		User        string `form:"user" validate:"required"`
		Password    string `form:"password" validate:"required"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(svr.AssetDBAdd(c, p.Name, p.Description, p.Host, p.Port, p.User, p.Password))
}

func assetTableAdd(c *bm.Context) {
	p := new(struct {
		DB          string `form:"db" validate:"required"`
		Regex       string `form:"regex" validate:"required"`
		Fields      string `form:"fields" validate:"required"`
		Description string `form:"description"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(svr.AssetTableAdd(c, p.DB, p.Regex, p.Fields, p.Description))
}

func updateAssetTable(c *bm.Context) {
	p := new(struct {
		Name   string `form:"name" validate:"required"`
		Fields string `form:"fields" validate:"required"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(svr.UpdateAssetTable(c, p.Name, p.Fields))
}

func assetInfoV2(c *bm.Context) {
	p := new(struct {
		Name string `form:"name" validate:"required"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(svr.AssetInfoV2(c, p.Name))
}

func assetShowTables(c *bm.Context) {
	p := new(struct {
		DB string `form:"db" validate:"required"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(svr.AssetShowTables(c, p.DB))
}

func assetTableFields(c *bm.Context) {
	p := new(struct {
		DB    string `form:"db" validate:"required"`
		Regex string `form:"regex" validate:"required"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	fs, count, err := svr.AssetTableFields(c, p.DB, p.Regex)
	data := &struct {
		Fields []*model.TableField `json:"fields"`
		Count  int                 `json:"count"`
	}{
		Fields: fs,
		Count:  count,
	}
	c.JSON(data, err)
}

func clusterOwners(c *bm.Context) {
	c.JSON(svr.ClusterOwners(), nil)
}
