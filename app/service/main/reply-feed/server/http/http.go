package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-common/app/service/main/reply-feed/conf"
	"go-common/app/service/main/reply-feed/model"
	"go-common/app/service/main/reply-feed/service"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	svc *service.Service
	vfy *verify.Verify
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	svc = s
	vfy = verify.New(c.Verify)
	engine := bm.DefaultServer(c.BM)
	router(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func router(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/x/reply-feed")
	{
		g.POST("/strategy/new", newEGroup)
		g.POST("/strategy/edit", editEgroup)
		g.POST("/strategy/state", modifyState)
		g.POST("/strategy/resize", resizeSlots)
		g.POST("/strategy/reset", resetEgroup)
		g.GET("/strategy/list", listEGroup)

		g.GET("/statistics/list", statistics)
	}
}

func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

func validate(algorithm string, weight string) (err error) {
	switch algorithm {
	case model.WilsonLHRRAlgorithm:
		w := &model.WilsonLHRRWeight{}
		if err = json.Unmarshal([]byte(weight), &w); err != nil {
			return
		}
		return w.Validate()
	case model.WilsonLHRRFluidAlgorithm:
		w := &model.WilsonLHRRFluidWeight{}
		if err = json.Unmarshal([]byte(weight), &w); err != nil {
			return
		}
		return w.Validate()
	case model.OriginAlgorithm, model.LikeDescAlgorithm:
		return
	default:
		log.Error("unknown algorithm accepted (%s)", algorithm)
		err = ecode.RequestErr
	}
	if err != nil {
		return
	}
	return
}

func listEGroup(c *bm.Context) {
	stats, err := svc.SlotStatsManager(c)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(stats, nil)
}

func newEGroup(c *bm.Context) {
	v := new(struct {
		Name      string `form:"name" validate:"required"`
		Percent   int    `form:"percent" validate:"required"`
		Algorithm string `form:"algorithm" validate:"required"`
		Weight    string `form:"weight" validate:"required"`
	})
	var (
		err error
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if err = validate(v.Algorithm, v.Weight); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.NewEGroup(c, v.Name, v.Algorithm, v.Weight, v.Percent))
}

func editEgroup(c *bm.Context) {
	v := new(struct {
		Name      string  `form:"name" validate:"required"`
		Algorithm string  `form:"algorithm" validate:"required"`
		Weight    string  `form:"weight" validate:"required"`
		Slots     []int64 `form:"slots,split" validate:"required"`
	})
	var (
		err error
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if err = validate(v.Algorithm, v.Weight); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.EditSlotsStat(c, v.Name, v.Algorithm, v.Weight, v.Slots))
}

func modifyState(c *bm.Context) {
	v := new(struct {
		Name  string `form:"name" validate:"required"`
		State int    `form:"state"`
	})
	var (
		err error
	)
	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, svc.ModifyState(c, v.Name, v.State))
}

func resizeSlots(c *bm.Context) {
	v := new(struct {
		Name    string `form:"name" validate:"required"`
		Percent int    `form:"percent" validate:"required"`
	})
	var (
		err error
	)
	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, svc.ResizeSlots(c, v.Name, v.Percent))
}

func resetEgroup(c *bm.Context) {
	v := new(struct {
		Name string `form:"name" validate:"required"`
	})
	var (
		err error
	)
	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, svc.ResetEGroup(c, v.Name))
}

func statistics(c *bm.Context) {
	v := new(model.SSReq)
	if err := c.Bind(v); err != nil {
		return
	}
	if v.Hour {
		res := new(model.SSHourRes)
		data, err := svc.StatisticsByHour(c, v)
		if err != nil {
			c.JSON(nil, err)
			return
		}
		xAxisMap := make(map[string]struct{})
		res.Series = make(map[string][]*model.StatisticsStat)
		for legend, m := range data {
			res.Legend = append(res.Legend, legend)
			for xAxis := range m {
				xAxisMap[xAxis] = struct{}{}
			}
		}
		for k := range xAxisMap {
			res.XAxis = append(res.XAxis, k)
		}
		for _, legend := range res.Legend {
			if hourStatistics, ok := data[legend]; ok {
				for _, hour := range res.XAxis {
					if stat, exists := hourStatistics[hour]; exists {
						res.Series[legend] = append(res.Series[legend], stat)
					} else {
						t := strings.Split(hour, "-")
						if len(t) < 2 {
							c.Abort()
							return
						}
						date, _ := strconv.Atoi(t[0])
						hour, _ := strconv.Atoi(t[1])
						res.Series[legend] = append(res.Series[legend], &model.StatisticsStat{Date: date, Hour: hour})
					}
				}
			}
		}
		res.Sort()
		c.JSON(res, nil)
		return
	}
	res := new(model.SSDateRes)
	data, err := svc.StatisticsByDate(c, v)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	xAxisMap := make(map[int]struct{})
	res.Series = make(map[string][]*model.StatisticsStat)
	for legend, m := range data {
		res.Legend = append(res.Legend, legend)
		for xAxis := range m {
			xAxisMap[xAxis] = struct{}{}
		}
	}
	for k := range xAxisMap {
		res.XAxis = append(res.XAxis, k)
	}
	for _, legend := range res.Legend {
		if dateStatistics, ok := data[legend]; ok {
			for _, date := range res.XAxis {
				if stat, exists := dateStatistics[date]; exists {
					res.Series[legend] = append(res.Series[legend], stat)
				} else {
					res.Series[legend] = append(res.Series[legend], &model.StatisticsStat{Date: date})
				}
			}
		}
	}
	res.Sort()
	c.JSON(res, nil)
}
