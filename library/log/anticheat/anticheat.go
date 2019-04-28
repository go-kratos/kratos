package anticheat

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

// AntiCheat send anti-cheating info to berserker.
type AntiCheat struct {
	infoc *infoc.Infoc
}

// New new AntiCheat logger.
func New(c *infoc.Config) (a *AntiCheat) {
	return &AntiCheat{infoc: infoc.New(c)}
}

// antiCheat 尽可能多的提供信息.
type antiCheat struct {
	Buvid  string
	Build  string
	Client string // for example ClientWeb
	IP     string
	UID    string
	Aid    string
	Mid    string

	Sid      string
	Refer    string
	URL      string
	From     string
	ItemID   string
	ItemType string // for example ItemTypeAv

	Action   string // for example ActionClick
	ActionID string
	UA       string
	TS       string
	Extra    string
}

// anti-cheat const.
const (
	ClientWeb     = "web"
	ClientIphone  = "iphone"
	ClientIpad    = "ipad"
	ClientAndroid = "android"
	// AntiCheat ItemType
	ItemTypeAv       = "av"
	ItemTypeBangumi  = "bangumi"
	ItemTypeLive     = "live"
	ItemTypeTopic    = "topic"
	ItemTypeRank     = "rank"
	ItemTypeActivity = "activity"
	ItemTypeTag      = "tag"
	ItemTypeAD       = "ad"
	ItemTypeLV       = "lv"

	// AntiCheat Action
	ActionClick     = "click"
	ActionPlay      = "play"
	ActionFav       = "fav"
	ActionCoin      = "coin"
	ActionDM        = "dm"
	ActionToView    = "toview"
	ActionShare     = "share"
	ActionSpace     = "space"
	Actionfollow    = "follow"
	ActionHeartbeat = "heartbeat"
	ActionAnswer    = "answer"
)

func (a *antiCheat) toSlice() (as []interface{}) {
	as = make([]interface{}, 0, 18)
	as = append(as, a.Buvid, a.Build, a.Client, a.IP, a.UID, a.Aid, a.Mid)
	as = append(as, a.Sid, a.Refer, a.URL, a.From, a.ItemID, a.ItemType)
	as = append(as, a.Action, a.ActionID, a.UA, a.TS, a.Extra)
	return
}

// InfoAntiCheat2 for new http framework(bm).
func (a *AntiCheat) InfoAntiCheat2(ctx *bm.Context, uid, aid, mid, itemID, itemType, action, actionID string) error {
	return a.infoAntiCheat(ctx, ctx.Request, metadata.String(ctx, metadata.RemoteIP), uid, aid, mid, itemID, itemType, action, actionID)
}

// infoAntiCheat common logic.
func (a *AntiCheat) infoAntiCheat(ctx context.Context, req *http.Request, IP, uid, aid, mid, itemID, itemType, action, actionID string) error {
	params := req.Form
	ac := &antiCheat{
		UID:      uid,
		Aid:      aid,
		Mid:      mid,
		ItemID:   itemID,
		ItemType: itemType,
		Action:   action,
		ActionID: actionID,
		IP:       IP,
		URL:      req.URL.Path,
		Refer:    req.Header.Get("Referer"),
		UA:       req.Header.Get("User-Agent"),
		TS:       strconv.FormatInt(time.Now().Unix(), 10),
	}
	ac.From = params.Get("from")
	if csid, err := req.Cookie("sid"); err == nil {
		ac.Sid = csid.Value
	}
	var cli string
	switch {
	case len(params.Get("access_key")) == 0:
		cli = ClientWeb
		if ck, err := req.Cookie("buvid3"); err == nil {
			ac.Buvid = ck.Value
		}
	case params.Get("platform") == "ios":
		cli = ClientIphone
		if params.Get("device") == "pad" {
			cli = ClientIpad
		}
	case params.Get("platform") == "android":
		cli = ClientAndroid
	default:
		log.Warn("unkown plat(%s)", params.Get("platform"))
	}
	ac.Client = cli
	if cli != ClientWeb {
		ac.Buvid = req.Header.Get("buvid")
		ac.Build = params.Get("build")
	}
	return a.infoc.Infov(ctx, ac.toSlice()...)
}

// ServiceAntiCheat common anti-cheat.
func (a *AntiCheat) ServiceAntiCheat(p map[string]string) error {
	return a.infoc.Info(convertBase(p)...)
}

// ServiceAntiCheatBus for answer anti-cheat.
func (a *AntiCheat) ServiceAntiCheatBus(p map[string]string, bus []interface{}) error {
	ac := append(convertBase(p), bus...)
	return a.infoc.Info(ac...)
}

// ServiceAntiCheatv support mirror request
func (a *AntiCheat) ServiceAntiCheatv(ctx context.Context, p map[string]string) error {
	return a.infoc.Infov(ctx, convertBase(p)...)
}

// ServiceAntiCheatBusv support mirror request
func (a *AntiCheat) ServiceAntiCheatBusv(ctx context.Context, p map[string]string, bus []interface{}) error {
	ac := append(convertBase(p), bus...)
	return a.infoc.Infov(ctx, ac...)
}

func convertBase(p map[string]string) (res []interface{}) {
	ac := &antiCheat{
		ItemType: p["itemType"],
		Action:   p["action"],
		IP:       p["ip"],
		Mid:      p["mid"],
		UID:      p["fid"],
		Aid:      p["aid"],
		Sid:      p["sid"],
		UA:       p["ua"],
		Buvid:    p["buvid"],
		Refer:    p["refer"],
		URL:      p["url"],
		TS:       strconv.FormatInt(time.Now().Unix(), 10),
	}
	res = ac.toSlice()
	return
}
