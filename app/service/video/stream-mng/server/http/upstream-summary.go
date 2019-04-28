package http

import (
	"fmt"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"strconv"
	"time"
)

// getSummaryUpStreamRtmp 查询统计信息
func getSummaryUpStreamRtmp(c *bm.Context) {
	st, ed, err := analysisTimeParams(c)
	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	res, err := srv.GetSummaryUpStreamRtmp(c, st, ed)
	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.JSONMap(map[string]interface{}{"data": res}, nil)
}

// getSummaryUpStreamISP 获取运营商信息统计
func getSummaryUpStreamISP(c *bm.Context) {
	st, ed, err := analysisTimeParams(c)
	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	res, err := srv.GetSummaryUpStreamISP(c, st, ed)
	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.JSONMap(map[string]interface{}{"data": res}, nil)
}

// getSummaryUpStreamISP 获取运营商信息统计
func getSummaryUpStreamCountry(c *bm.Context) {
	st, ed, err := analysisTimeParams(c)
	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	res, err := srv.GetSummaryUpStreamCountry(c, st, ed)
	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.JSONMap(map[string]interface{}{"data": res}, nil)
}

// getSummaryUpStreamPlatform 获取Platform信息统计
func getSummaryUpStreamPlatform(c *bm.Context) {
	st, ed, err := analysisTimeParams(c)
	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	res, err := srv.GetSummaryUpStreamPlatform(c, st, ed)
	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.JSONMap(map[string]interface{}{"data": res}, nil)
}

// getSummaryUpStreamCity 获取City信息统计
func getSummaryUpStreamCity(c *bm.Context) {
	st, ed, err := analysisTimeParams(c)
	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	res, err := srv.GetSummaryUpStreamCity(c, st, ed)
	if err != nil {
		c.Set("output_data", err.Error())
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}

	c.JSONMap(map[string]interface{}{"data": res}, nil)
}

// analysisTimeParams 分析传入参数
func analysisTimeParams(c *bm.Context) (int64, int64, error) {
	params := c.Request.URL.Query()
	start := params.Get("start")
	end := params.Get("end")

	var st int64
	var ed int64
	var err error

	if start == "" {
		t := time.Now()
		year, month, day := t.Date()
		st = time.Date(year, month, day, 0, 0, 0, 0, t.Location()).Unix()
	} else {
		st, err = strconv.ParseInt(start, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("start is not right")
		}
	}

	if end == "" {
		ed = time.Now().Unix()
	} else {
		ed, err = strconv.ParseInt(end, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("end is not right")
		}
	}

	return st, ed, nil
}
