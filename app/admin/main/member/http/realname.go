package http

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"go-common/app/admin/main/member/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

func realnameList(ctx *bm.Context) {
	var (
		arg   = &model.ArgRealnameList{}
		list  []*model.RespRealnameApply
		total int
		err   error
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if arg.PN <= 0 {
		arg.PN = 1
	}
	if arg.PS <= 0 || arg.PS > 100 {
		arg.PS = 10
	}
	// 没有指定实名认证申请状态State，默认查询所有
	if arg.State == "" {
		arg.State = model.RealnameApplyStateAll
	}
	// 如果没有使用 mid || card 则取默认7天前作为 TSFrom
	if arg.MID == 0 && arg.Card == "" {
		if arg.TSFrom <= 0 {
			arg.TSFrom = time.Now().Add(-time.Hour * 24 * 7).Unix()
		}
	}
	// 如果使用 mid 或 card 进行查询，则查询所有时间段对应的数据
	if arg.MID != 0 || arg.Card != "" {
		arg.TSFrom = 0
		arg.TSTo = 0
	}

	if list, total, err = svc.RealnameList(ctx, arg); err != nil {
		ctx.JSON(nil, err)
		return
	}
	res := map[string]interface{}{
		"list": list,
		"page": map[string]int{
			"num":   arg.PN,
			"size":  arg.PS,
			"total": total,
		},
	}
	ctx.JSON(res, nil)
}

func realnamePendingList(ctx *bm.Context) {
	var (
		arg   = &model.ArgRealnamePendingList{}
		list  []*model.RespRealnameApply
		total int
		err   error
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if arg.PN <= 0 {
		arg.PN = 1
	}
	if arg.PS <= 0 || arg.PS > 100 {
		arg.PS = 10
	}
	if list, total, err = svc.RealnamePendingList(ctx, arg); err != nil {
		ctx.JSON(nil, err)
		return
	}
	res := map[string]interface{}{
		"list": list,
		"page": map[string]int{
			"num":   arg.PN,
			"size":  arg.PS,
			"total": total,
		},
	}
	ctx.JSON(res, nil)
}

func realnameAuditApply(ctx *bm.Context) {
	var (
		arg       = &model.ArgRealnameAuditApply{}
		adminName string
		adminID   int64
		data      interface{}
		ok        bool
		err       error
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	// 驳回必须有reason
	if arg.Action == model.RealnameActionReject && arg.Reason == "" {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if data, ok = ctx.Get("username"); !ok {
		ctx.JSON(nil, ecode.AccessDenied)
		return
	}
	if adminName, ok = data.(string); !ok {
		ctx.JSON(nil, ecode.ServerErr)
		return
	}
	if data, ok = ctx.Get("uid"); !ok {
		ctx.JSON(nil, ecode.AccessDenied)
		return
	}
	if adminID, ok = data.(int64); !ok {
		ctx.JSON(nil, ecode.ServerErr)
		return
	}
	ctx.JSON(nil, svc.RealnameAuditApply(ctx, arg, adminName, adminID))
}

func realnameReasonList(ctx *bm.Context) {
	var (
		arg   = &model.ArgRealnameReasonList{}
		list  []string
		total int
		err   error
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if arg.PN <= 0 {
		arg.PN = 1
	}
	if arg.PS <= 0 || arg.PS > 100 {
		arg.PS = 20
	}
	if list, total, err = svc.RealnameReasonList(ctx, arg); err != nil {
		ctx.JSON(nil, err)
		return
	}
	res := map[string]interface{}{
		"list": list,
		"page": map[string]int{
			"num":   arg.PN,
			"size":  arg.PS,
			"total": total,
		},
	}
	ctx.JSON(res, nil)
}

func realnameSetReason(ctx *bm.Context) {
	var (
		arg = &model.ArgRealnameSetReason{}
	)
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(nil, svc.RealnameSetReason(ctx, arg))
}

func realnameImage(ctx *bm.Context) {
	var (
		arg = &model.ArgRealnameImage{}
		img []byte
		err error
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if img, err = svc.FetchRealnameImage(ctx, arg.Token); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.Bytes(http.StatusOK, http.DetectContentType(img), img)
}

func realnameImagePreview(ctx *bm.Context) {
	var (
		arg = &model.ArgRealnameImagePreview{}
		img []byte
		err error
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if arg.BorderSize <= 0 {
		arg.BorderSize = 1000
	}
	if img, err = svc.RealnameImagePreview(ctx, arg.Token, arg.BorderSize); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.Bytes(http.StatusOK, http.DetectContentType(img), img)
}

func realnameSearchCard(ctx *bm.Context) {
	var (
		arg = &model.ArgRealnameSearchCard{}
		err error
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if len(arg.Cards) > 50 {
		err = ecode.RequestErr
		return
	}
	ctx.JSON(svc.RealnameSearchCard(ctx, arg.Cards, arg.CardType, arg.Country))
}

func realnameUnbind(ctx *bm.Context) {
	var (
		arg       = &model.ArgMid{}
		adminName string
		adminID   int64
		data      interface{}
		ok        bool
		err       error
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if data, ok = ctx.Get("username"); !ok {
		ctx.JSON(nil, ecode.AccessDenied)
		return
	}
	if adminName, ok = data.(string); !ok {
		ctx.JSON(nil, ecode.ServerErr)
		return
	}
	if data, ok = ctx.Get("uid"); !ok {
		ctx.JSON(nil, ecode.AccessDenied)
		return
	}
	if adminID, ok = data.(int64); !ok {
		ctx.JSON(nil, ecode.ServerErr)
		return
	}
	ctx.JSON(nil, svc.RealnameUnbind(ctx, arg.Mid, adminName, adminID))
}

func realnameExport(ctx *bm.Context) {
	arg := &model.ArgMids{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	if len(arg.Mid) > 1000 {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}

	realnameExports, err := svc.RealnameExcel(ctx, arg.Mid)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	data := make([][]string, 0, len(realnameExports))
	data = append(data, []string{"用户uid", "用户名", "昵称", "姓名", "手机号", "证件类型", "证件号"})
	for _, user := range realnameExports {
		fields := []string{
			strconv.FormatInt(user.Mid, 10),
			user.UserID,
			user.Uname,
			user.Realname,
			user.Tel,
			model.CardTypeString(user.CardType),
			avoidescape(user.CardNum),
		}
		data = append(data, fields)
	}
	buf := bytes.NewBuffer(nil)
	w := csv.NewWriter(buf)
	for _, record := range data {
		if err := w.Write(record); err != nil {
			ctx.JSON(nil, err)
			return
		}
	}
	w.Flush()
	res := buf.Bytes()
	ctx.Writer.Header().Set("Content-Type", "application/csv")
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", "实名认证导出"))
	ctx.Writer.Write([]byte("\xEF\xBB\xBF")) // 写 UTF-8 的 BOM 头，别删，删了客服就找上门了
	ctx.Writer.Write(res)
}

func avoidescape(in string) string {
	return fmt.Sprintf("%s,", in)
}

func realnameSubmit(ctx *bm.Context) {
	arg := &model.ArgRealnameSubmit{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(nil, svc.RealnameSubmit(ctx, arg))
}

func realnameFileUpload(c *bm.Context) {
	arg := &model.ArgMid{}
	if err := c.Bind(arg); err != nil {
		return
	}
	mid := arg.Mid
	defer c.Request.Form.Del("img") // 防止日志不出现
	c.Request.ParseMultipartForm(32 << 20)
	imgBytes, err := func() ([]byte, error) {
		img := c.Request.FormValue("img")
		if img != "" {
			log.Info("Succeeded to parse img file from form value: mid: %d, length: %d", mid, len(img))
			return []byte(img), nil
		}
		log.Warn("Failed to parse img file from form value, fallback to form file: mid: %d", mid)
		f, _, err := c.Request.FormFile("img")
		if err != nil {
			return nil, errors.Wrapf(err, "parse img form file: mid: %d", mid)
		}
		defer f.Close()
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, errors.Wrapf(err, "read img form file: mid: %d", mid)
		}
		if len(data) <= 0 {
			return nil, errors.Wrapf(err, "form file data: mid: %d, length: %d", mid, len(data))
		}
		log.Info("Succeeded to parse file from form file: mid: %d, length: %d", mid, len(data))
		return data, nil
	}()
	if err != nil {
		log.Error("Failed to parse realname upload file: mid: %d: %+v", mid, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var resData struct {
		SRC string `json:"token"`
	}
	if resData.SRC, err = svc.RealnameFileUpload(c, arg.Mid, imgBytes); err != nil {
		log.Error("%+v", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(resData, nil)
}
