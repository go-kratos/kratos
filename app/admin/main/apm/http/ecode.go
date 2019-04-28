package http

import (
	"encoding/json"
	"go-common/app/admin/main/apm/model/ecode"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func ecodeList(c *bm.Context) {
	v := new(struct {
		Code string `form:"code"`
		Pn   int    `form:"pn" default:"1"`
		Ps   int    `form:"ps" default:"20"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	var (
		cs    []*codes.Codes
		count int
		lk    = "%" + v.Code + "%"
	)
	if v.Code != "" {
		err = apmSvc.DB.Where("CAST(code AS CHAR) LIKE ? OR message LIKE ?", lk, lk).Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&cs).Error
	} else {
		err = apmSvc.DB.Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&cs).Error
	}
	if err != nil {
		log.Error("apmSvc.EcodeList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if v.Code != "" {
		err = apmSvc.DB.Where("CAST(code AS CHAR) LIKE ? OR message LIKE ?", lk, lk).Model(&codes.Codes{}).Count(&count).Error
	} else {
		err = apmSvc.DB.Model(&codes.Codes{}).Count(&count).Error
	}
	if err != nil {
		log.Error("apmSvc.EcodeList count error(%v)", err)
		c.JSON(nil, err)
		return
	}
	data := &Paper{
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: cs,
		Total: count,
	}
	c.JSON(data, nil)
}

func ecodeAdd(c *bm.Context) {
	v := new(struct {
		Code        int32  `form:"code" validate:"required"`
		Message     string `form:"message" validate:"required"`
		HantMessage string `form:"hant_message"`
		Level       int8   `form:"level" default:"2"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	cnt := 0
	if err = apmSvc.DB.Model(&codes.Codes{}).Where("code=?", v.Code).Count(&cnt).Error; err != nil {
		log.Error("apmSvc.EcodeAdd count error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	username := name(c)
	code := &codes.Codes{
		Code:        v.Code,
		Message:     v.Message,
		HantMessage: v.HantMessage,
		Level:       v.Level,
		Operator:    username,
	}
	if err = apmSvc.DB.Create(&code).Error; err != nil {
		log.Error("apmSvc.EcodeAdd create error(%v)", err)
		c.JSON(nil, err)
		return
	}
	apmSvc.SendLog(*c, username, 0, 1, code.ID, "apmSvc.EcodeAdd", code)
	c.JSON(nil, err)
}

func ecodeEdit(c *bm.Context) {
	v := new(struct {
		ID          int64  `form:"id" validate:"required"`
		Code        int32  `form:"code" validate:"required"`
		Message     string `form:"message" validate:"required"`
		HantMessage string `form:"hant_message"`
		Level       int8   `form:"level" default:"2"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	code := &codes.Codes{}
	if err = apmSvc.DB.Where("id = ?", v.ID).Find(code).Error; err != nil {
		log.Error("apmSvc.EcodeEdit find(%d) error(%v)", v.ID, err)
		c.JSON(nil, err)
		return
	}
	username := name(c)
	ups := map[string]interface{}{
		"code":         v.Code,
		"message":      v.Message,
		"operator":     username,
		"hant_message": v.HantMessage,
		"level":        v.Level,
	}
	if err = apmSvc.DB.Model(&codes.Codes{}).Where("id = ?", v.ID).Updates(ups).Error; err != nil {
		log.Error("apmSvc.EcodeEdit updates error(%v)", err)
		c.JSON(nil, err)
		return
	}
	sqlLog := &map[string]interface{}{
		"SQLType": "update",
		"Where":   "id = ?",
		"Value1":  v.ID,
		"Update":  ups,
		"Old":     code,
	}
	apmSvc.SendLog(*c, username, 0, 2, v.ID, "apmSvc.EcodeEdit", sqlLog)
	c.JSON(nil, err)
}

func ecodeDelete(c *bm.Context) {
	v := new(struct {
		ID int64 `form:"id" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	code := &codes.Codes{}
	if err = apmSvc.DB.Where("id = ?", v.ID).Find(code).Error; err != nil {
		log.Error("apmSvc.EcodeDelete find(%d) error(%v)", v.ID, err)
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.DB.Delete(code).Error; err != nil {
		log.Error("apmSvc.EcodeDelete delete(%d) error(%v)", v.ID, err)
		c.JSON(nil, err)
		return
	}
	username := name(c)
	apmSvc.SendLog(*c, username, 0, 3, code.ID, "apmSvc.EcodeDelete", code)
	c.JSON(nil, err)
}

func getEcodes(c *bm.Context) {
	v := new(struct {
		Interval1 int `form:"interval1" default:"0"`
		Interval2 int `form:"interval2" default:"0"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	codes := []*codes.Codes{}
	if err = apmSvc.DB.Where("code >= ? and code <= ?", v.Interval1, v.Interval2).Find(&codes).Error; err != nil {
		log.Error("getEcodes error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(codes, nil)
}

func syncEcodes(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(struct {
		Data string `form:"data"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	username := name(c)
	var Codes []*codes.Codes
	err = json.Unmarshal([]byte(v.Data), &Codes)
	if err != nil {
		res["message"] = "解析失败"
		c.JSONMap(res, err)
		return
	}
	var in []int32
	for _, val := range Codes {
		in = append(in, val.Code)
	}
	oldCodes := []*codes.Codes{}
	if err = apmSvc.DB.Where("code in (?)", in).Find(&oldCodes).Error; err != nil {
		log.Error("getEcodes error(%v)", err)
		c.JSON(nil, err)
		return
	}
	tmp := make(map[int32]struct{})
	for _, val := range oldCodes {
		tmp[val.Code] = struct{}{}
	}
	var code *codes.Codes
	tx := apmSvc.DB.Begin()
	for _, val := range Codes {
		code = &codes.Codes{}
		code.Code = val.Code
		code.Message = val.Message
		code.HantMessage = val.HantMessage
		code.Level = val.Level
		code.Operator = username
		if _, ok := tmp[val.Code]; ok {
			//update
			if err = tx.Model(codes.Codes{}).Where("code = ?", val.Code).Updates(code).Error; err != nil {
				tx.Rollback()
				c.JSON(nil, err)
				return
			}
		} else {
			//create
			if err = tx.Create(code).Error; err != nil {
				tx.Rollback()
				c.JSON(nil, err)
				return
			}
		}
	}
	tx.Commit()
	c.JSON(nil, err)
}

func getProdEcodes(c *bm.Context) {
	v := new(struct {
		Interval1 string `form:"interval1" default:"0"`
		Interval2 string `form:"interval2" default:"0"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	Data, err := apmSvc.GetCodes(c, v.Interval1, v.Interval2)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(Data, nil)
}

func codeLangsList(c *bm.Context) {
	v := new(struct {
		Code string `form:"code"`
		Pn   int    `form:"pn" default:"1"`
		Ps   int    `form:"ps" default:"20"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	var (
		cs    []*codes.Codes
		count int
		lk    = "%" + v.Code + "%"
		in    []int64
		ts    []*codes.CodeMsg
		res   []*codes.NewCodes
	)

	if v.Code != "" {
		err = apmSvc.DB.Where("CAST(code AS CHAR) LIKE ? OR message LIKE ?", lk, lk).Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&cs).Error
	} else {
		err = apmSvc.DB.Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&cs).Error
	}
	if err != nil {
		log.Error("apmSvc.EcodeList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	for _, val := range cs {
		in = append(in, val.ID)
	}
	err = apmSvc.DB.Where("code_id in(?)", in).Find(&ts).Error
	if err != nil {
		log.Error("apmSvc.EcodeList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	for _, val := range cs {
		list := make([]*codes.CodeMsg, 0)
		for _, vv := range ts {
			if vv.CodeID == val.ID {
				list = append(list, vv)
			}
		}
		res = append(res, &codes.NewCodes{Codes: val, List: list})
	}

	if v.Code != "" {
		err = apmSvc.DB.Where("CAST(code AS CHAR) LIKE ? OR message LIKE ?", lk, lk).Model(&codes.Codes{}).Count(&count).Error
	} else {
		err = apmSvc.DB.Model(&codes.Codes{}).Count(&count).Error
	}
	if err != nil {
		log.Error("apmSvc.EcodeList count error(%v)", err)
		c.JSON(nil, err)
		return
	}
	data := &Paper{
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: res,
		Total: count,
	}
	c.JSON(data, nil)
}

func ecodeLangsAdd(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(struct {
		Data string `form:"data" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	var Codes []*codes.CodeMsg
	err = json.Unmarshal([]byte(v.Data), &Codes)
	if err != nil {
		res["message"] = "解析失败"
		c.JSONMap(res, err)
		return
	}
	username := name(c)
	for _, val := range Codes {
		if !(val.CodeID > 0 && len(val.Locale) > 0 && len(val.Msg) > 0) {
			err = ecode.RequestErr
			res["message"] = "数据有误"
			c.JSONMap(res, err)
			return
		}
		code := &codes.CodeMsg{
			CodeID:   val.CodeID,
			Locale:   val.Locale,
			Msg:      val.Msg,
			Operator: username,
		}
		if err = apmSvc.DB.Create(&code).Error; err != nil {
			log.Error("apmSvc.ecodeLangsAdd create error(%v)", err)
			c.JSON(nil, err)
			return
		}
		apmSvc.SendLog(*c, username, 0, 1, code.ID, "apmSvc.ecodeLangsAdd", code)
	}
	c.JSON(nil, err)
}

func ecodeLangsEdit(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(struct {
		Data string `form:"data" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	var Codes []*codes.CodeMsg
	err = json.Unmarshal([]byte(v.Data), &Codes)
	if err != nil {
		res["message"] = "解析失败"
		c.JSONMap(res, err)
		return
	}
	username := name(c)
	for _, val := range Codes {
		code := &codes.CodeMsg{}
		if err = apmSvc.DB.Where("id = ?", val.ID).First(code).Error; err != nil {
			log.Error("apmSvc.ecodeLangsEdit find(%d) error(%v)", val.ID, err)
			c.JSON(nil, err)
			return
		}
		ups := map[string]interface{}{
			"operator": username,
			"msg":      val.Msg,
			"locale":   val.Locale,
		}
		if err = apmSvc.DB.Model(&codes.CodeMsg{}).Where("id = ?", val.ID).Updates(ups).Error; err != nil {
			log.Error("apmSvc.ecodeLangsEdit updates error(%v)", err)
			c.JSON(nil, err)
			return
		}
		sqlLog := &map[string]interface{}{
			"SQLType": "update",
			"Where":   "id = ?",
			"Value1":  val.ID,
			"Update":  ups,
			"Old":     code,
		}
		apmSvc.SendLog(*c, username, 0, 2, val.ID, "apmSvc.ecodeLangsEdit", sqlLog)
	}
	c.JSON(nil, err)
}

func ecodeLangsDelete(c *bm.Context) {
	v := new(struct {
		CID int64 `form:"cid" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	code := &codes.CodeMsg{}
	if err = apmSvc.DB.Where("id = ?", v.CID).First(code).Error; err != nil {
		log.Error("apmSvc.ecodeLangsDelete find(%d) error(%v)", v.CID, err)
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.DB.Delete(code).Error; err != nil {
		log.Error("apmSvc.ecodeLangsDelete delete(%d) error(%v)", v.CID, err)
		c.JSON(nil, err)
		return
	}
	username := name(c)
	apmSvc.SendLog(*c, username, 0, 3, code.ID, "apmSvc.ecodeLangsDelete", code)
	c.JSON(nil, err)
}

func ecodeLangsSave(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(struct {
		Data string `form:"data" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	var Codes []*codes.CodeMsg
	err = json.Unmarshal([]byte(v.Data), &Codes)
	if err != nil {
		res["message"] = "解析失败"
		c.JSONMap(res, err)
		return
	}
	username := name(c)
	for _, val := range Codes {
		code := &codes.CodeMsg{}
		if val.ID == 0 {
			if !(val.CodeID > 0 && len(val.Locale) > 0 && len(val.Msg) > 0) {
				err = ecode.RequestErr
				res["message"] = "数据有误"
				c.JSONMap(res, err)
				return
			}
			code = &codes.CodeMsg{
				CodeID:   val.CodeID,
				Locale:   val.Locale,
				Msg:      val.Msg,
				Operator: username,
			}
			if err = apmSvc.DB.Create(&code).Error; err != nil {
				log.Error("apmSvc.ecodeLangsSvae create error(%v)", err)
				c.JSON(nil, err)
				return
			}
			apmSvc.SendLog(*c, username, 0, 1, code.ID, "apmSvc.ecodeLangsAdd", code)
		} else {
			if err = apmSvc.DB.Where("id = ?", val.ID).First(code).Error; err != nil {
				log.Error("apmSvc.ecodeLangsSvae find(%d) error(%v)", val.ID, err)
				c.JSON(nil, err)
				return
			}
			ups := map[string]interface{}{
				"operator": username,
				"msg":      val.Msg,
				"locale":   val.Locale,
			}
			if err = apmSvc.DB.Model(&codes.CodeMsg{}).Where("id = ?", val.ID).Updates(ups).Error; err != nil {
				log.Error("apmSvc.ecodeLangsSave updates error(%v)", err)
				c.JSON(nil, err)
				return
			}
			sqlLog := &map[string]interface{}{
				"SQLType": "update",
				"Where":   "id = ?",
				"Value1":  val.ID,
				"Update":  ups,
				"Old":     code,
			}
			apmSvc.SendLog(*c, username, 0, 2, val.ID, "apmSvc.ecodeLangsSave", sqlLog)
		}
	}
	c.JSON(nil, err)
}
