package param

import "go-common/app/interface/main/app-resource/model"

// Param struct
type Param struct {
	Name      string `json:"-"`
	Value     string `json:"-"`
	Plat      int8   `json:"-"`
	Build     int    `json:"-"`
	Condition string `json:"-"`
}

// Change func
func (p *Param) Change() {
	switch p.Plat {
	case 10:
		p.Plat = model.PlatAndroidB
	case 11:
		p.Plat = model.PlatAndroidTVYST
	case 12:
		p.Plat = model.PlatIPhoneB
	}
}
