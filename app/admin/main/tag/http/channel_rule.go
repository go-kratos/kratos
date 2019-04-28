package http

import (
	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func channelRuleCheck(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			InRule    []string `form:"in_name,split"`
			InRuleID  []int64  `form:"in_ids,split"`
			NotInRule string   `form:"notin_name"`
			NotInID   int64    `form:"notin_id"`
		})
		inRules   []string
		inRuleIDs []int64
		rule      *model.ChannelRule
	)
	if err = c.Bind(param); err != nil {
		return
	}
	lenthNames := len(param.InRule)
	lenthIDs := len(param.InRuleID)
	if lenthNames == 0 && lenthIDs == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if lenthNames > 0 && lenthIDs > 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if lenthNames > 0 {
		for _, name := range param.InRule {
			if name, err = checkName(name, model.TNameMaxLen); err != nil {
				c.JSON(nil, err)
				return
			}
			inRules = append(inRules, name)
		}
		if param.NotInRule != "" {
			if lenthNames+1 > model.ChannelRuleMaxLen {
				c.JSON(nil, ecode.RequestErr)
				return
			}
			if param.NotInRule, err = checkName(param.NotInRule, model.TNameMaxLen); err != nil {
				c.JSON(nil, err)
				return
			}
		}
		rule, err = svc.ChannelRuleCheckNames(c, inRules, param.NotInRule)
	} else {
		for _, id := range param.InRuleID {
			if id <= 0 {
				c.JSON(nil, ecode.RequestErr)
				return
			}
			inRuleIDs = append(inRuleIDs, id)
		}
		if param.NotInID > 0 {
			if lenthIDs+1 > model.ChannelRuleMaxLen {
				c.JSON(nil, ecode.RequestErr)
				return
			}
		}
		rule, err = svc.ChannelRuleCheckIDs(c, inRuleIDs, param.NotInID)
	}
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(rule, nil)
}
