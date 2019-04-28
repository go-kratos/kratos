package http

import (
	"go-common/app/service/openplatform/anti-fraud/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

// qusBankInfo 题库信息
func qusBankInfo(c *bm.Context) {
	params := new(model.ArgGetQusBank)
	if err := c.Bind(params); err != nil {
		return
	}
	c.JSON(svc.GetQusBankInfo(c, params.QsBId))
}

// qusBankList 题库列表
func qusBankList(c *bm.Context) {
	params := new(model.ArgBankList)
	if err := c.Bind(params); err != nil {
		return
	}

	result := model.RespList{}
	list, err := svc.GetQusBanklist(c, params.PageNo, params.PageSize, params.Name)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	total, err := svc.GetQusBankTotal(c, params.Name)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	result.Items = list
	result.PageSize = params.PageSize
	result.PageNo = params.PageNo
	result.Total = total
	c.JSON(result, nil)
}

// qusBankAdd 添加题库
func qusBankAdd(c *bm.Context) {
	params := new(model.ArgAddQusBank)
	if err := c.BindWith(params, binding.JSON); err != nil {
		return
	}
	c.JSON(svc.AddQusBank(c, params))
}

// qusBankEdit 更新题库
func qusBankEdit(c *bm.Context) {
	params := new(model.ArgUpdateQusBank)
	if err := c.BindWith(params, binding.JSON); err != nil {
		return
	}
	c.JSON(svc.UpdataQusBank(c, params.QsBId, params.QBName, params.MaxRetryTime, params.CdTime))
}

// qusBankDel 删除题库
func qusBankDel(c *bm.Context) {
	params := new(model.ArgBaseBank)
	if err := c.BindWith(params, binding.JSON); err != nil {
		return
	}
	c.JSON(svc.DelQusBank(c, params.QsBId, model.DeledStatus))
}

// qusInfo 题目详情
func qusInfo(c *bm.Context) {
	params := new(model.ArgGetQus)
	if err := c.Bind(params); err != nil {
		return
	}
	data := &model.QuestionAll{}

	answers, err := svc.GetAnswerList(c, params.QsID)
	if err != nil {
		c.JSON(nil, ecode.QusIDInvalid)
		return
	}

	info, err := svc.GetQusInfo(c, params.QsID)
	if err != nil || info == nil {
		c.JSON(nil, ecode.QusIDInvalid)
		return
	}
	data.Question = *info
	data.AnswersList = answers

	c.JSON(data, err)
}

// qusList 题目列表
func qusList(c *bm.Context) {
	params := new(model.ArgQusList)
	if err := c.Bind(params); err != nil {
		return
	}
	result := model.RespList{}
	list, err := svc.GetQuslist(c, params.PageNo, params.PageSize, params.QsBId)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	total, err := svc.GetQusTotal(c, params.QsBId)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	result.Items = list
	result.PageSize = params.PageSize
	result.PageNo = params.PageNo
	result.Total = total
	c.JSON(result, nil)
}

// qusAdd 题目添加
func qusAdd(c *bm.Context) {
	params := new(model.ArgAddQus)
	if err := c.BindWith(params, binding.JSON); err != nil {
		return
	}
	anlist := params.Answer
	msg, err := svc.CheckAnswer(c, 0, params.Type, anlist)
	if err != nil {
		c.JSON(msg, err)
		return
	}
	addQus := &model.AddQus{
		BId:     params.BId,
		Type:    params.Type,
		Name:    params.Name,
		Dif:     params.Dif,
		AnType:  params.AnType,
		Answers: anlist,
		QsID:    0,
	}
	c.JSON(svc.AddQus(c, addQus, anlist))
}

// qusUpdate 题目更新
func qusUpdate(c *bm.Context) {
	params := new(model.ArgUpdateQus)
	if err := c.BindWith(params, binding.JSON); err != nil {
		return
	}
	answers := params.Answer
	msg, err := svc.CheckAnswer(c, params.QsID, params.Type, answers)
	if err != nil {
		c.JSON(msg, err)
		return
	}

	c.JSON(svc.UpdateQus(c, params, answers))
}

// questionBankBind 项目关联题库
func questionBankBind(c *bm.Context) {
	params := &model.ArgQuestionBankBinds{}
	if err := c.BindWith(params, binding.JSON); err != nil {
		return
	}

	err := svc.QuestionBankBind(c, params)
	if err != nil {
		log.Error("questionBankBind(%v) error(%v)", params, err)
	}

	c.JSON(nil, err)
}

// questionBankUnbind 解绑
func questionBankUnbind(c *bm.Context) {
	params := &model.ArgQuestionBankUnbind{}
	if err := c.BindWith(params, binding.JSON); err != nil {
		return
	}

	c.JSON(nil, svc.QuestionBankUnbind(c, params))
}

// qusDel 删除题目
func qusDel(c *bm.Context) {
	params := new(model.ArgGetQus)
	if err := c.BindWith(params, binding.JSON); err != nil {
		return
	}
	c.JSON(svc.DelQus(c, params.QsID))
}

// getBankBind 绑定关系
func getBankBind(c *bm.Context) {
	params := &model.ArgGetBankBind{}
	if err := c.BindWith(params, binding.JSON); err != nil {
		return
	}

	c.JSON(svc.GetQuestionBankBind(c, params))
}

// getBindItems 查询绑定到题库的
func getBindItems(c *bm.Context) {
	params := &model.ArgGetBindItems{}
	if err := c.Bind(params); err != nil {
		return
	}

	c.JSON(svc.GetBindItem(c, params))
}

// getQuestion 随机获取一个问题
func getQuestion(c *bm.Context) {
	params := &model.ArgGetQuestion{}
	if err := c.Bind(params); err != nil {
		return
	}

	c.JSON(svc.GetQuestion(c, params))
}

// answerQuestion 答题
func answerQuestion(c *bm.Context) {
	// 判断是否可答 返回冷却时间
	params := &model.ArgCheckAnswer{}
	if err := c.BindWith(params, binding.JSON); err != nil {
		return
	}

	c.JSON(svc.UserAnswer(c, params))

}

// qusBankCheck 题库检查
func qusBankCheck(c *bm.Context) {
	// 判断是否可答 返回冷却时间
	params := &model.ArgCheckQus{}
	if err := c.BindWith(params, binding.JSON); err != nil {
		return
	}
	c.JSON(svc.QusBankCheck(c, params))
}
