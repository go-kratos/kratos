package sobot

// Check the TicketParam
func (tp *TicketParam) Check() bool {
	if tp.TicketTitle == "" || tp.TicketID == 0 || tp.TicketContent == "" || tp.CustomerEmail == "" {
		return false
	}
	return true
}

// CheckModify the TicketParam modiy
func (tp *TicketParam) CheckModify() bool {
	if tp.TicketID == 0 || tp.CustomerEmail == "" {
		return false
	}
	return true
}

// Check the ReplyParam
func (rp *ReplyParam) Check() bool {
	if rp.CustomerEmail == "" || rp.ReplyContent == "" || rp.TicketID == 0 {
		return false
	}
	return true
}
