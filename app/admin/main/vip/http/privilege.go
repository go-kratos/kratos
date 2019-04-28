package http

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

const (
	_maxnamelen    = 30
	_maxtitlelen   = 60
	_maxexplainlen = 1200
)

// regexp utf8 char 0x0e0d~0e4A
var (
	_emptyUnicodeReg = []*regexp.Regexp{
		regexp.MustCompile(`[\x{202e}]+`),  // right-to-left override
		regexp.MustCompile(`[\x{200b}]+`),  // zeroWithChar
		regexp.MustCompile(`[\x{1f6ab}]+`), // no_entry_sign
	}
	// trim
	returnReg  = regexp.MustCompile(`[\n]{3,}`)
	returnReg2 = regexp.MustCompile(`(\r\n){3,}`)
	spaceReg   = regexp.MustCompile(`[　]{5,}`) // Chinese quanjiao space character
)

func privileges(c *bm.Context) {
	var err error
	arg := new(struct {
		Langtype int8 `form:"lang_type"`
	})
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.Privileges(c, arg.Langtype))
}

func updatePrivilegeState(c *bm.Context) {
	var err error
	arg := new(model.ArgStatePrivilege)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, vipSvc.UpdatePrivilegeState(c, &model.Privilege{
		ID:    arg.ID,
		State: arg.Status,
	}))
}

func deletePrivilege(c *bm.Context) {
	var err error
	arg := new(model.ArgPivilegeID)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, vipSvc.DeletePrivilege(c, arg.ID))
}

func updateOrder(c *bm.Context) {
	var err error
	arg := new(model.ArgOrder)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, vipSvc.UpdateOrder(c, arg))
}

func addPrivilege(c *bm.Context) {
	var err error
	arg := new(model.ArgAddPrivilege)
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	arg.Operator = username.(string)
	if err = c.BindWith(arg, binding.FormMultipart); err != nil {
		return
	}
	if len(arg.Name) > _maxnamelen {
		c.JSON(nil, ecode.VipPrivilegeNameTooLongErr)
		return
	}
	if len(arg.Title) > _maxtitlelen {
		c.JSON(nil, ecode.VipPrivilegeTitleTooLongErr)
		return
	}
	if len(arg.Explain) > _maxexplainlen {
		c.JSON(nil, ecode.VipPrivilegeExplainTooLongErr)
		return
	}
	img := new(model.ArgImage)
	if img.IconBody, img.IconFileType, err = file(c, "icon"); err != nil {
		c.JSON(nil, err)
		return
	}
	if img.IconFileType == "" {
		c.JSON(nil, ecode.VipFileImgEmptyErr)
		return
	}
	if img.IconGrayBody, img.IconGrayFileType, err = file(c, "gray_icon"); err != nil {
		c.JSON(nil, err)
		return
	}
	if img.IconGrayFileType == "" {
		c.JSON(nil, ecode.VipFileImgEmptyErr)
		return
	}
	if img.WebImageBody, img.WebImageFileType, err = file(c, "web_image"); err != nil {
		c.JSON(nil, err)
		return
	}
	if img.AppImageBody, img.AppImageFileType, err = file(c, "app_image"); err != nil {
		c.JSON(nil, err)
		return
	}
	arg.Explain = filterContent(arg.Explain)
	c.JSON(nil, vipSvc.AddPrivilege(c, arg, img))
}

func updatePrivilege(c *bm.Context) {
	var (
		err error
	)
	arg := new(model.ArgUpdatePrivilege)
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	arg.Operator = username.(string)
	if err = c.BindWith(arg, binding.FormMultipart); err != nil {
		return
	}
	if len(arg.Name) > _maxnamelen {
		c.JSON(nil, ecode.VipPrivilegeNameTooLongErr)
		return
	}
	if len(arg.Title) > _maxtitlelen {
		c.JSON(nil, ecode.VipPrivilegeTitleTooLongErr)
		return
	}
	if len(arg.Explain) > _maxexplainlen {
		c.JSON(nil, ecode.VipPrivilegeExplainTooLongErr)
		return
	}
	img := new(model.ArgImage)
	if img.IconBody, img.IconFileType, err = file(c, "icon"); err != nil {
		c.JSON(nil, err)
		return
	}
	if img.IconGrayBody, img.IconGrayFileType, err = file(c, "gray_icon"); err != nil {
		c.JSON(nil, err)
		return
	}
	if img.WebImageBody, img.WebImageFileType, err = file(c, "web_image"); err != nil {
		c.JSON(nil, err)
		return
	}
	if img.AppImageBody, img.AppImageFileType, err = file(c, "app_image"); err != nil {
		c.JSON(nil, err)
		return
	}
	arg.Explain = filterContent(arg.Explain)
	c.JSON(nil, vipSvc.UpdatePrivilege(c, arg, img))
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

func filterContent(str string) string {
	tmp := str
	// check params
	tmp = strings.TrimSpace(tmp)
	tmp = spaceReg.ReplaceAllString(tmp, "　　　")
	tmp = returnReg.ReplaceAllString(tmp, "\n\n\n")
	tmp = returnReg2.ReplaceAllString(tmp, "\n\n\n")
	// checkout empty
	for _, reg := range _emptyUnicodeReg {
		tmp = reg.ReplaceAllString(tmp, "")
	}
	return tmp
}
