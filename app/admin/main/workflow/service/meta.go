package service

import "strconv"

// BusinessRoleState business->role->state
var (
	_state = map[int8]map[int8]map[int8]string{
		1:  {0: {0: "投诉刚创建", 1: "被接受", 2: "被驳回"}},
		2:  {0: {0: "未处理", 1: "有效", 2: "无效"}, 1: {1: "未处理", 2: "已回复已读", 3: "管理员关闭", 4: "用户已解决", 5: "过期自动关闭", 6: "已回复未读", 7: "用户关闭", 8: "通过关闭"}},
		3:  {0: {0: "投诉刚创建", 1: "被接受", 2: "被驳回"}},
		4:  {0: {0: "投诉刚创建", 1: "被接受", 2: "被驳回"}},
		5:  {0: {0: "未处理", 1: "有效", 2: "无效"}, 1: {1: "未处理", 2: "已回复已读", 3: "管理员关闭", 4: "用户已解决", 5: "过期自动关闭", 6: "已回复未读", 7: "用户关闭", 8: "通过关闭"}},
		6:  {0: {0: "未处理", 1: "通过", 2: "驳回", 3: "失效"}},
		8:  {0: {0: "未处理", 1: "有效", 2: "无效"}},
		9:  {0: {0: "未处理", 1: "有效", 2: "无效"}},
		13: {0: {0: "未处理", 1: "有效", 2: "无效", 9: "已删除", 10: "移交众裁"}},
		14: {0: {0: "未处理", 1: "有效", 2: "无效"}},
	}

	// BusinessIDName bid->name
	_business = map[int8]string{
		1:  "稿件投诉",
		2:  "稿件申诉",
		3:  "短点评投诉",
		4:  "长点评投诉",
		5:  "小黑屋申诉",
		6:  "稿件审核",
		7:  "任务质检",
		8:  "音频Tag",
		9:  "频道举报",
		13: "评论举报",
		14: "字幕举报",
	}

	_role = map[int8]string{
		0: "处理",
		1: "反馈",
	}

	_flow = map[int8]string{
		0: "审核流",
		1: "投诉流",
		2: "申诉流",
	}
)

// StateDescr state description
func (s *Service) StateDescr(business, role, state int8) string {
	//TODO: describe all states
	stateStr, ok := _state[business][role][state]
	if !ok {
		return strconv.Itoa(int(state))
	}
	return stateStr
}

// StateDescV3 .
func (s *Service) StateDescV3(business, fid, state int8) string {
	return ""
}

// BusinessDesc business description
func (s *Service) BusinessDesc(business int8) string {
	busStr, ok := _business[business]
	if !ok {
		return strconv.Itoa(int(business))
	}
	return busStr
}

// RoleDesc role description
func (s *Service) RoleDesc(role int8) string {
	roleStr, ok := _role[role]
	if !ok {
		return strconv.Itoa(int(role))
	}
	return roleStr
}

// FlowDesc .
func (s *Service) FlowDesc(fid int8) string {
	flowStr, ok := _flow[fid]
	if !ok {
		return strconv.Itoa(int(fid))
	}
	return flowStr
}

// RidDesc .
func (s *Service) RidDesc(bid, rid int8) string {
	roleStr, ok := s.roleCache[bid][rid]
	if !ok {
		return strconv.Itoa(int(rid))
	}
	return roleStr
}
