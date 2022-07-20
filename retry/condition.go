package retry

type Condition interface {
	Judge(r Resp) bool
}

type Resp struct {
	MD   map[string][]string
	Code int
}

type ByCode struct {
	Codes []int
}

func (c *ByCode) Judge(r Resp) bool {
	if len(c.Codes) == 0 {
		return false
	}
	if len(c.Codes) == 1 {
		return r.Code == c.Codes[0]
	}

	return c.Codes[0] <= r.Code && r.Code <= c.Codes[1]
}

func NewByCode(code ...int) *ByCode {
	return &ByCode{Codes: code}
}

type ByMetedata struct {
	Key  string
	Vals map[string]struct{}
}

func (c *ByMetedata) Judge(r Resp) bool {
	value, ok := r.MD[c.Key]
	if !ok {
		return false
	}
	for _, v := range value {
		_, ok := c.Vals[v]
		if ok {
			return true
		}
	}
	return false
}

func NewByMetedata(key string, vals ...string) *ByMetedata {
	parsedVals := map[string]struct{}{}
	for _, v := range vals {
		parsedVals[v] = struct{}{}
	}
	return &ByMetedata{Key: key, Vals: parsedVals}
}
