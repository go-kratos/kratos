package http

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"go-common/app/admin/main/card/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func groups(c *bm.Context) {
	var err error
	arg := new(model.ArgQueryGroup)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.GroupList(c, arg))
}

func cards(c *bm.Context) {
	var err error
	arg := new(model.ArgQueryCards)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.CardsByGid(c, arg.GroupID))
}

func addGroup(c *bm.Context) {
	var err error
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	arg := new(model.AddGroup)
	arg.Operator = username.(string)
	if err = c.Bind(arg); err != nil {
		return
	}
	if len(arg.Name) > _maxnamelen {
		c.JSON(nil, ecode.CardNameTooLongErr)
		return
	}
	c.JSON(nil, srv.AddGroup(c, arg))
}

func updateGroup(c *bm.Context) {
	var err error
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	arg := new(model.UpdateGroup)
	arg.Operator = username.(string)
	if err = c.Bind(arg); err != nil {
		return
	}
	if len(arg.Name) > _maxnamelen {
		c.JSON(nil, ecode.CardNameTooLongErr)
		return
	}
	c.JSON(nil, srv.UpdateGroup(c, arg))
}

func addCard(c *bm.Context) {
	var err error
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	arg := new(model.AddCard)
	arg.Operator = username.(string)
	if err = c.BindWith(arg, binding.FormMultipart); err != nil {
		return
	}
	if len(arg.Name) > _maxnamelen {
		c.JSON(nil, ecode.CardNameTooLongErr)
		return
	}
	if arg.CardBody, arg.CardFileType, err = file(c, "card_url"); err != nil {
		c.JSON(nil, err)
		return
	}
	if arg.CardFileType == "" {
		c.JSON(nil, ecode.CardImageEmptyErr)
		return
	}
	if arg.BigCardBody, arg.BigCardFileType, err = file(c, "big_crad_url"); err != nil {
		c.JSON(nil, err)
		return
	}
	if arg.BigCardFileType == "" {
		c.JSON(nil, ecode.CardImageEmptyErr)
		return
	}
	c.JSON(nil, srv.AddCard(c, arg))
}

func updateCard(c *bm.Context) {
	var err error
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	arg := new(model.UpdateCard)
	arg.Operator = username.(string)
	if err = c.BindWith(arg, binding.FormMultipart); err != nil {
		return
	}
	if len(arg.Name) > _maxnamelen {
		c.JSON(nil, ecode.CardNameTooLongErr)
		return
	}
	if arg.CardBody, arg.CardFileType, err = file(c, "card_url"); err != nil {
		c.JSON(nil, err)
		return
	}
	if arg.BigCardBody, arg.BigCardFileType, err = file(c, "big_crad_url"); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, srv.UpdateCard(c, arg))
}

func cardStateChange(c *bm.Context) {
	var err error
	arg := new(model.ArgState)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, srv.UpdateCardState(c, arg))
}

func deleteCard(c *bm.Context) {
	var err error
	arg := new(model.ArgID)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, srv.DeleteCard(c, arg.ID))
}

func groupStateChange(c *bm.Context) {
	var err error
	arg := new(model.ArgState)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, srv.UpdateGroupState(c, arg))
}

func deleteGroup(c *bm.Context) {
	var err error
	arg := new(model.ArgID)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, srv.DeleteGroup(c, arg.ID))
}

func cardOrderChange(c *bm.Context) {
	var err error
	arg := new(model.ArgIds)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, srv.CardOrderChange(c, arg))
}

func groupOrderChange(c *bm.Context) {
	var err error
	arg := new(model.ArgIds)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, srv.GroupOrderChange(c, arg))
}

func file(c *bm.Context, name string) (body []byte, filetype string, err error) {
	var file multipart.File
	if file, _, err = c.Request.FormFile(name); err != nil {
		if err == http.ErrMissingFile {
			err = nil
			return
		}
		err = ecode.RequestErr
		return
	}
	if file == nil {
		return
	}
	defer file.Close()
	if body, err = ioutil.ReadAll(file); err != nil {
		err = ecode.RequestErr
		return
	}
	filetype = http.DetectContentType(body)
	if err = checkImgFileType(filetype); err != nil {
		return
	}
	err = checkFileBody(body)
	return
}

func checkImgFileType(filetype string) error {
	switch filetype {
	case "image/jpeg", "image/jpg":
	case "image/png":
	default:
		return ecode.VipFileTypeErr
	}
	return nil
}

func checkFileBody(body []byte) error {
	if len(body) == 0 {
		return ecode.FileNotExists
	}
	if len(body) > cf.Bfs.MaxFileSize {
		return ecode.FileTooLarge
	}
	return nil
}
