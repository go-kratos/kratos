package http

import (
	"fmt"
	"strings"
	"time"

	"go-common/app/interface/bbq/common/model"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/trace"

	jsoniter "github.com/json-iterator/go"
)

// UILog .
type UILog struct {
	infoc *infoc.Infoc
}

// New .
func New(c *infoc.Config) *UILog {
	return &UILog{
		infoc: infoc.New(c),
	}
}

// Infoc .
func (u *UILog) Infoc(ctx *bm.Context, action int, ext interface{}) {
	arg := new(model.Base)
	ctx.Bind(arg)
	log.V(5).Info("uilog [%+v]", arg)

	if action == model.ActionPlay {
		// data type 1 心跳上报 2 暂停 3 home
		switch arg.DataType {
		case 2:
			action = model.ActionPlayPause
		case 3:
			action = model.ActionPlayOut
		}
	}

	// interface base field
	app := arg.App
	client := arg.Client
	version := arg.Version
	channel := arg.Channel
	loc := arg.Location
	cip := ctx.Request.RemoteAddr
	if ips := ctx.Request.Header.Get("X-Forwarded-For"); ips != "" {
		ipArr := strings.Split(ips, ",")
		if len(ipArr) > 0 {
			cip = ipArr[0]
		}
	}

	// session id
	sid, _ := ctx.Get("SessionID")

	// trace id
	tracer, _ := trace.FromContext(ctx.Context)
	tid := tracer

	// recsys query id
	qid := arg.QueryID
	if action == model.ActionRecommend {
		qid = fmt.Sprintf("%s", tracer)
	}

	// video
	svid := arg.SVID

	// user
	var mid int64
	if tmp, ok := ctx.Get("mid"); ok {
		mid = tmp.(int64)
	}
	buvid := ctx.Request.Header.Get("Buvid")
	totalDuration := arg.TotalDuration
	playDuration := arg.PlayDuration
	ctime := time.Now().Unix()

	b, _ := jsoniter.Marshal(ext)
	u.infoc.Info(client, app, version, channel, loc, cip, qid, sid, tid, svid, mid, buvid, action, totalDuration, playDuration, ctime, string(b), arg.From, arg.FromID, arg.PFrom, arg.PFromID)
}
