package http

import (
	"go-common/app/admin/main/app/model/language"
	bm "go-common/library/net/http/blademaster"
)

// languages select language all
func languages(c *bm.Context) {
	c.JSON(langSvc.Languages(c))
}

// langByID select language by id
func langByID(c *bm.Context) {
	v := &language.Param{}
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(langSvc.LangByID(c, v.ID))
}

// addOrup insert or update language
func addOrup(c *bm.Context) {
	var (
		err error
		v   = &language.Param{}
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID > 0 {
		err = langSvc.Update(c, v)
	} else {
		err = langSvc.Insert(c, v)
	}
	c.JSON(nil, err)
}
