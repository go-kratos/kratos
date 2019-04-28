package http

import (
	"encoding/json"

	"go-common/app/service/main/workflow/model/account"
	"go-common/app/service/main/workflow/model/sobot"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func sobotFetchUser(c *bm.Context) {
	data := []byte(`
    {
      "mid": 1,
      "uname": "biliuser",
      "tel": "132****1234",
      "email": "biliuser@qq.com",
      "status": 0,
      "formal": 0,
      "moral": 70,
      "level": 3,
      "exp": "4000",
      "coin": 300.12,
      "bcoin": 10.12,
      "medal": "青铜殿堂",
      "up": {
        "relation": {
          "following": 1,
          "whisper": 1,
          "black": 0,
          "follower": 1
        },
        "archive": 5,
        "identify": 1,
        "office": "bilibili认证",
        "shell": 10.12,
        "bank_card": "6227123412341234123"
      },
      "extra": {
        "arc_pubed": 312,
        "arc_not_pubed": 34,
        "arc_is_pubing": 175
      }
    }
    `)
	user := &account.User{}
	user.Extra = make(map[string]interface{})
	if err := json.Unmarshal(data, user); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(user, nil)
}

func sobotInfoTicket(c *bm.Context) {
	tp := new(struct {
		TicketID int32 `form:"ticket_id" validate:"required"`
	})
	if err := c.Bind(tp); err != nil {
		return
	}
	c.JSON(wkfSvc.SobotTicketInfo(c, tp.TicketID))
}

func sobotAddTicket(c *bm.Context) {
	tp := new(sobot.TicketParam)
	if err := c.BindWith(tp, binding.FormPost); err != nil {
		return
	}
	c.JSON(nil, wkfSvc.SobotTicketAdd(c, tp))
}

func sobotModifyTicket(c *bm.Context) {
	tp := new(sobot.TicketParam)
	if err := c.BindWith(tp, binding.FormPost); err != nil {
		return
	}
	c.JSON(nil, wkfSvc.SobotTicketModify(c, tp))
}

func sobotAddReply(c *bm.Context) {
	rp := new(sobot.ReplyParam)
	if err := c.BindWith(rp, binding.FormPost); err != nil {
		return
	}
	c.JSON(nil, wkfSvc.SobotReplyAdd(c, rp))
}

// func sobotCallback(c *bm.Context) {
// 	req := c.Request
// 	bs, err := ioutil.ReadAll(req.Body)
// 	if err != nil {
// 		log.Error("ioutil.ReadAll() error(%v)", err)
// 		c.JSON(nil, ecode.RequestErr)
// 		return
// 	}
// 	req.Body.Close()
// 	var jsbody map[string]interface{}
// 	if err := json.Unmarshal(bs, &jsbody); err != nil {
// 		c.JSON(nil, ecode.RequestErr)
// 		return
// 	}
// 	log.Info("sobotCallback(%s)", string(bs))
// 	c.JSON(jsbody, nil)
// }
