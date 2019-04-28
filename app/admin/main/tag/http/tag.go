package http

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	regHTML      = regexp.MustCompile(`(?i)\<.*?(script|href|img|src)+.*?\>`)
	regSymbol    = regexp.MustCompile(`^[g\pP|g\pS]+$`)
	regZeroWidth = regexp.MustCompile(`[\x{200b}]+`)
)

func tagList(c *bm.Context) {
	var (
		err   error
		tags  *model.MngSearchTagList
		esTag = new(model.ESTag)
	)
	if err = c.Bind(esTag); err != nil {
		return
	}
	if tags, err = svc.TagList(c, esTag); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	data["tag"] = tags.Result
	data["page"] = map[string]interface{}{
		"page":     esTag.Pn,
		"pagesize": esTag.Ps,
		"total":    tags.Page.Total,
	}
	c.JSON(data, nil)
}

func tagEdit(c *bm.Context) {
	var (
		err   error
		param = new(model.ParamTagEdit)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.TName, err = checkName(param.TName, model.TNameMaxLen); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, svc.TagEdit(c, param.Tid, param.TP, param.TName, param.Content))
}

func tagInfo(c *bm.Context) {
	var (
		err   error
		info  *model.TagInfo
		param = new(struct {
			Tid   int64  `form:"tid"`
			TName string `form:"tname"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.Tid <= 0 && param.TName == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if param.Tid > 0 {
		if info, err = svc.TagInfo(c, param.Tid); err != nil {
			c.JSON(nil, err)
			return
		}
	} else {
		if param.TName, err = checkName(param.TName, model.TNameMaxLen); err != nil {
			c.JSON(nil, err)
			return
		}
		if info, err = svc.TagInfoByName(c, param.TName); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	c.JSON(info, nil)
}

func tagState(c *bm.Context) {
	var (
		err error
		v   = new(model.ParamTagState)
	)
	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, svc.TagState(c, v.Tid, v.State))
}

func tagVerify(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			Tid int64 `form:"tid" validate:"required,gt=0"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, svc.TagVerify(c, param.Tid))
}

func tagCheck(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			Tid   int64  `form:"tid"`
			TName string `form:"tname"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.Tid <= 0 && param.TName == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if param.TName != "" {
		if param.TName, err = checkName(param.TName, model.TNameMaxLen); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	c.JSON(svc.TagCheck(c, param.Tid, param.TName))
}

// checkName check if name illegal
func checkName(name string, nameLen int) (dst string, err error) {
	if !utf8.ValidString(name) {
		log.Error("utf8.ValidString(%s)", name)
		err = ecode.RequestErr
		return
	}
	index := regHTML.FindAllString(name, -1)
	if len(index) > 0 {
		err = ecode.RequestErr
		return
	}
	dst = regZeroWidth.ReplaceAllString(name, "")
	dst = strings.TrimSpace(dst)
	dst = replace(dst)
	if dst == "" || len([]rune(dst)) > nameLen {
		log.Error("name == nil or length > max length: %s", name)
		err = ecode.RequestErr
		return
	}
	if regSymbol.MatchString(dst) {
		log.Error("cant not contain continuous symbol(%s)", name)
		err = ecode.RequestErr
	}
	return
}

func replace(name string) string {
	var rb []byte
	sb := []byte(name)
	for _, b := range sb {
		if b < 0x20 || b == 0x7f {
			continue
		}
		rb = append(rb, b)
	}
	return string(rb)
}
