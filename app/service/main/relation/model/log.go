package model

// Reverse is
func (rl *RelationLog) Reverse() *RelationLog {
	content := make(map[string]string, len(rl.Content))
	for k, v := range rl.Content {
		content[k] = v
	}
	reversed := &RelationLog{
		// reverse
		Mid:    rl.Fid,
		Fid:    rl.Mid,
		Ts:     rl.Ts,
		Source: rl.Source,
		Ip:     rl.Ip,
		Buvid:  rl.Buvid,
		// reverse
		FromAttr:    rl.FromRevAttr,
		ToAttr:      rl.ToRevAttr,
		FromRevAttr: rl.FromAttr,
		ToRevAttr:   rl.ToAttr,

		Content: content,
	}
	return reversed
}
