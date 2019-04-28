package http

import (
	"strconv"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func subtitleLanAdd(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			Code     uint8  `form:"code" validate:"required"`
			Lan      string `form:"lan" validate:"required"`
			DocEn    string `form:"doc_en" validate:"required"`
			DocZh    string `form:"doc_zh" validate:"required"`
			IsDelete bool   `form:"is_delete"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, dmSvc.SubtitleLanOp(c, v.Code, v.Lan, v.DocZh, v.DocEn, v.IsDelete))
}

func subtitleFilter(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			Words string `form:"words" validate:"required"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(dmSvc.SubtitleFilter(c, v.Words))
}

func draftSave(c *bm.Context) {
	var (
		v = new(struct {
			Oid              int64  `form:"oid" validate:"required"`
			Type             int32  `form:"type" validate:"required"`
			Aid              int64  `form:"aid" validate:"required" `
			Lan              string `form:"lan" validate:"required"`
			Submit           bool   `form:"submit"`
			Sign             bool   `form:"sign"`
			OriginSubtitleID int64  `form:"origin_subtitle_id"`
			Data             string `form:"data" validate:"required"`
		})
	)
	mid, _ := c.Get("mid")
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(dmSvc.SaveSubtitleDraft(c, v.Aid, v.Oid, v.Type, mid.(int64), v.Lan, v.Submit, v.Sign, v.OriginSubtitleID, []byte(v.Data)))
}

func assitAudit(c *bm.Context) {
	var (
		v = new(struct {
			Oid           int64  `form:"oid" validate:"required"`
			SubtitleID    int64  `form:"subtitle_id" validate:"required"`
			Pass          bool   `form:"pass"`
			RejectComment string `form:"reject_comment"`
		})
	)
	mid, _ := c.Get("mid")
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, dmSvc.AuditSubtitle(c, v.Oid, v.SubtitleID, mid.(int64), v.Pass, v.RejectComment))

}

func subtitleDel(c *bm.Context) {
	var (
		v = new(struct {
			Oid        int64 `form:"oid" validate:"required"`
			SubtitleID int64 `form:"subtitle_id" validate:"required"`
		})
	)
	mid, _ := c.Get("mid")
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, dmSvc.DelSubtitle(c, v.Oid, v.SubtitleID, mid.(int64)))
}

// 锁定发布字幕
func subtitleLock(c *bm.Context) {
	var (
		v = new(struct {
			Oid        int64 `form:"oid" validate:"required"`
			Type       int32 `form:"type" validate:"required"`
			SubtitleID int64 `form:"subtitle_id" validate:"required"`
			Lock       bool  `form:"lock"`
		})
	)
	mid, _ := c.Get("mid")
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, dmSvc.SubtitleLock(c, v.Oid, v.Type, mid.(int64), v.SubtitleID, v.Lock))
}

// 署名字幕
func subtitleSign(c *bm.Context) {
	var (
		v = new(struct {
			Oid        int64 `form:"oid" validate:"required"`
			Type       int32 `form:"type" validate:"required"`
			SubtitleID int64 `form:"subtitle_id" validate:"required"`
			Sign       bool  `form:"sign"`
		})
	)
	mid, _ := c.Get("mid")
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, dmSvc.SubtitleSign(c, v.Oid, v.Type, mid.(int64), v.SubtitleID, v.Sign))
}

func subtitleArchiveName(c *bm.Context) {
	var (
		v = new(struct {
			Aid int64 `form:"aid" validate:"required"`
		})
	)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(dmSvc.ArchiveName(c, v.Aid))
}

// 查看字幕内容
func subtitleShow(c *bm.Context) {
	var (
		v = new(struct {
			Oid        int64 `form:"oid" validate:"required"`
			SubtitleID int64 `form:"subtitle_id" validate:"required"`
		})
	)
	mid, _ := c.Get("mid")
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(dmSvc.SubtitleShow(c, v.Oid, v.SubtitleID, mid.(int64)))
}

func subtitleLans(c *bm.Context) {
	var (
		v = new(struct {
			Oid  int64 `form:"oid" validate:"required"`
			Type int32 `form:"type" validate:"required"`
		})
	)
	mid, _ := c.Get("mid")
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(dmSvc.SubtitleLans(c, v.Oid, v.Type, mid.(int64)))
}

func searchAssist(c *bm.Context) {
	var (
		v = new(struct {
			Oid    int64 `form:"oid"`
			Type   int32 `form:"type"`
			Aid    int64 `form:"aid"`
			Status int32 `form:"status"`
			Page   int32 `form:"page"`
			Size   int32 `form:"size" validate:"max=100"`
		})
	)
	mid, _ := c.Get("mid")
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(dmSvc.SearchAssist(c, v.Aid, v.Oid, v.Type, mid.(int64), v.Status, v.Page, v.Size))
}

func authorList(c *bm.Context) {
	var (
		v = new(struct {
			Status int32 `form:"status"`
			Page   int32 `form:"page"`
			Size   int32 `form:"size" validate:"max=100"`
		})
	)
	mid, _ := c.Get("mid")
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(dmSvc.SearchAuthor(c, mid.(int64), v.Status, v.Page, v.Size))
}

func subtitlePermission(c *bm.Context) {
	var (
		v = new(struct {
			Aid  int64 `form:"aid" validate:"required"`
			Oid  int64 `form:"oid" validate:"required"`
			Type int32 `form:"type" validate:"required"`
		})
	)
	mid, _ := c.Get("mid")
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, dmSvc.SubtitlePermission(c, v.Aid, v.Oid, v.Type, mid.(int64)))
}

// waveForm .
func waveForm(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			Aid  int64 `form:"aid" validate:"required"`
			Oid  int64 `form:"oid" validate:"required"`
			Type int32 `form:"type" validate:"required"`
		})
	)
	mid, _ := c.Get("mid")
	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(dmSvc.WaveForm(c, v.Aid, v.Oid, v.Type, mid.(int64)))
}

// waveFormCallBack .
func waveFormCallBack(c *bm.Context) {
	var (
		oid int64
		err error
		v   = new(struct {
			OK   int32  `json:"ok"`
			Info string `json:"info"`
		})
	)
	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}
	if oid, err = strconv.ParseInt(c.Request.URL.Query().Get("oid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, dmSvc.WaveFormCallBack(c, oid, 1, v.OK, v.Info))
}

func subtitleReportAdd(c *bm.Context) {
	var (
		err error
		v   = new(model.SubtitleReportAddParam)
	)
	mid, _ := c.Get("mid")
	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, dmSvc.SubtitleReportAdd(c, mid.(int64), v))
}

func subtitleReportTag(c *bm.Context) {
	c.JSON(dmSvc.SubtitleReportList(c))
}
