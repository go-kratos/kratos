package http

import (
	"strconv"
	"strings"

	"go-common/app/service/main/workflow/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

// AddChallenge add challenge
func addChallenge(c *bm.Context) {
	ap := new(model.ChallengeParam)
	if err := c.BindWith(ap, binding.FormPost); err != nil {
		return
	}
	if ap.AttachmentsStr != "" {
		ap.Attachments = strings.Split(ap.AttachmentsStr, ",")
	}
	if wkfSvc.TagMap(ap.Business, ap.Tid) == nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, ctrl := range wkfSvc.TagMap(ap.Business, ap.Tid).Controls {
		if ctrlValue := c.Request.PostForm.Get(ctrl.Name); ctrlValue != "" {
			ap.MetaData += ctrl.Name + ": " + ctrlValue + "\n"
		} else if ctrl.Required {
			log.Error("http addChallenge() control parms error ctrl.Name(%s) is required!  ap(%+v)", ctrl.Name, ap)
			c.JSON(nil, ecode.RequestErr)
			return
		} else {
			log.Info("http addChallenge() control parms missing ctrl.Name(%s) but not required ap(%+v)", ctrl.Name, ap)
			continue
		}
	}
	if !ap.CheckAdd() {
		log.Error("s.AddChallenge() params(%+v) error", ap)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	challengeNo, err := wkfSvc.AddChallenge(c, ap)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]int64{
		"challengeNo": challengeNo,
	}
	c.JSON(data, nil)
}

// ListChallenge get challenge list
func listChallenge(c *bm.Context) {
	ap := new(model.ChallengeParam)
	if err := c.Bind(ap); err != nil {
		return
	}
	if !ap.CheckList() {
		log.Error("s.Challenges() params(%+v) error", ap)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(wkfSvc.Challenges(c, ap))
}

// ReplyAddChallenge add reply to challenge
func replyAddChallenge(c *bm.Context) {
	rp := new(struct {
		Cid         int32  `form:"cid" validate:"required"`
		Event       int8   `form:"event" validate:"required"`
		Content     string `form:"content" validate:"required"`
		Attachments string `form:"attachments"`
	})
	if err := c.BindWith(rp, binding.FormPost); err != nil {
		return
	}
	_, err := wkfSvc.AddEvent(c, rp.Cid, rp.Content, rp.Attachments, rp.Event)
	c.JSON(nil, err)
}

// ChallengeInfo get challenge info
func challengeInfo(c *bm.Context) {
	ap := new(model.ChallengeParam)
	if err := c.Bind(ap); err != nil {
		return
	}
	if !ap.CheckInfo() {
		log.Error("s.Challenge() params(%+v) error", ap)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(wkfSvc.Challenge(c, ap))
}

// upChallengeState update challenge business state
func upChallengeState(c *bm.Context) {
	var role int8
	ap := new(struct {
		ID            int32 `form:"id" validate:"required"`
		Mid           int64 `form:"mid" validate:"required"`
		Business      int8  `form:"business" validate:"required"`
		BusinessState int8  `form:"business_state"`
	})
	roleStr := c.Request.PostForm.Get("role")
	if roleStr == "" {
		role = model.CustomerServiceRole
	} else {
		result, err := strconv.ParseUint(roleStr, 10, 8)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			c.Abort()
			return
		}
		role = int8(result)
	}
	if err := c.BindWith(ap, binding.FormPost); err != nil {
		return
	}
	c.JSON(nil, wkfSvc.UpChallengeState(c, ap.ID, ap.Mid, ap.Business, role, ap.BusinessState))
}

// CloseChallenge make challenge business state closed
func closeChallenge(c *bm.Context) {
	ap := new(struct {
		Cid           int32  `form:"cid" validate:"required"`
		Business      int8   `form:"business" validate:"required"`
		Role          int8   `form:"role" validate:"required"`
		BusinessState int8   `form:"business_state"`
		Note          string `form:"note" validate:"required"`
	})
	if err := c.Bind(ap); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, wkfSvc.CloseChallenge(c, ap.Cid, ap.Business, ap.Role, ap.BusinessState, ap.Note))
}

// untreatedChallenge get untreated challenge
func untreatedChallenge(c *bm.Context) {
	ap := new(struct {
		Oid  int64 `form:"oid" validate:"required"`
		Role int8  `form:"role" validate:"required"`
	})
	if err := c.Bind(ap); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(wkfSvc.UntreatedChallenge(c, ap.Oid, ap.Role))
}

// addChallenge3 add challange v3
func addChallenge3(c *bm.Context) {
	cp3 := &model.ChallengeParam3{}
	if err := c.Bind(cp3); err != nil {
		return
	}
	challengeNo, err := wkfSvc.AddChallenge3(c, cp3)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]int64{
		"challengeNo": challengeNo,
	}
	c.JSON(data, nil)
}

// listChallenge3 .
func listChallenge3(c *bm.Context) {
	cp3 := &model.ChallengeParam3{}
	if err := c.Bind(cp3); err != nil {
		return
	}
	c.JSON(wkfSvc.Challenges3(c, cp3))
}

// groupState3 .
func groupState3(c *bm.Context) {
	cp3 := &model.ChallengeParam3{}
	if err := c.Bind(cp3); err != nil {
		return
	}
	state, err := wkfSvc.GroupState3(c, cp3)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"state": state,
	}, err)
}
