package http

import (
	"fmt"
	"net/http"
	"time"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

func cheatUps(c *bm.Context) {
	v := new(struct {
		MID      int64  `form:"mid"`
		Nickname string `form:"nickname"`
		From     int    `form:"from" validate:"min=0" default:"0"`
		Limit    int    `form:"limit" validate:"min=1" default:"20"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	total, spies, err := svr.CheatUps(c, v.MID, v.Nickname, v.From, v.Limit)
	if err != nil {
		log.Error("growup svr.CheatUps error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    spies,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func cheatArchives(c *bm.Context) {
	v := new(struct {
		MID       int64  `form:"mid"`
		ArchiveID int64  `form:"archive_id"`
		Nickname  string `form:"nickname"`
		From      int    `form:"from" validate:"min=0" default:"0"`
		Limit     int    `form:"limit" validate:"min=1" default:"20"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	total, spies, err := svr.CheatArchives(c, v.MID, v.ArchiveID, v.Nickname, v.From, v.Limit)
	if err != nil {
		log.Error("growup svr.CheatUps error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    spies,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func exportCheatUps(c *bm.Context) {
	v := new(struct {
		MID      int64  `form:"mid"`
		Nickname string `form:"nickname"`
		From     int    `form:"from" validate:"min=0" default:"0"`
		Limit    int    `form:"limit" validate:"min=1" default:"20"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	content, err := svr.ExportCheatUps(c, v.MID, v.Nickname, v.From, v.Limit)
	if err != nil {
		log.Error("s.exportSpyUp ExportSpyUp error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, CSV{
		Content: content,
		Title:   fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02"), "cheat_ups"),
	})
}

func exportCheatAvs(c *bm.Context) {
	v := new(struct {
		MID       int64  `form:"mid"`
		Nickname  string `form:"nickname"`
		ArchiveID int64  `form:"archive_id"`
		From      int    `form:"from" validate:"min=0" default:"0"`
		Limit     int    `form:"limit" validate:"min=1" default:"20"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	content, err := svr.ExportCheatAvs(c, v.MID, v.ArchiveID, v.Nickname, v.From, v.Limit)
	if err != nil {
		log.Error("s.exportSpyAV ExportSpyAV error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, CSV{
		Content: content,
		Title:   fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02"), "cheat_avs"),
	})
}

func queryCheatFans(c *bm.Context) {
	v := new(struct {
		From  int64 `form:"from" default:"0" validate:"min=0"`
		Limit int64 `form:"limit" default:"20" validate:"min=1"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	total, fans, err := svr.QueryCheatFans(c, v.From, v.Limit)
	if err != nil {
		log.Error("growup svr.QueryCheatFans error(%v)", err)
		c.JSON(nil, err)
		return
	}

	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code": 0,
		"data": fans,
		"paging": map[string]int64{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func cheatFans(c *bm.Context) {
	v := new(struct {
		MID int64 `form:"mid"`
	})

	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.CheatFans(c, v.MID)
	if err != nil {
		log.Error("Exec cheatFans error!(%v)", err)
	}
	c.JSON(nil, err)
}
