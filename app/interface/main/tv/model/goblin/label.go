package goblin

import (
	"strings"

	"go-common/app/interface/main/tv/conf"
)

// Label def.
type Label struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Param     string `json:"param"`
	ParamName string `json:"param_name"`
	Value     string `json:"value"`
}

// TypeLabels def.
type TypeLabels struct {
	ParamName string   `json:"param_name"`
	Param     string   `json:"param"`
	Labels    []*Label `json:"labels"`
}

// FromLabels def.
func (v *TypeLabels) FromLabels(labels []*Label) {
	if len(labels) == 0 {
		return
	}
	v.Param = labels[0].Param
	v.ParamName = labels[0].ParamName
	v.Labels = labels
}

// IndexLabels is used to combine the data in memory
type IndexLabels struct {
	PGC map[int][]*TypeLabels // key is category, value is all the param and their labels
	UGC map[int][]*TypeLabels
}

// YearVDur def.
type YearVDur struct {
	Dur string `json:"dur"`
}

// TransYear transforms the value of year type labels
func (v *Label) TransYear(cfg *conf.IndexLabel) {
	if !cfg.IsYear(v.Param) {
		return
	}
	if len(cfg.YearV) == 0 {
		return
	}
	if newV, ok := cfg.YearV[v.Value]; ok { // replace the value
		v.Value = newV.Dur
	}
	if !strings.Contains(v.Value, "-") {
		v.Value = v.Value + "-" + v.Value
	} else { // transform 2004-2000 to 2000-2004
		years := strings.Split(v.Value, "-")
		if len(years) == 2 && years[0] != "" && years[1] != "" {
			v.Value = years[1] + "-" + years[0]
		}
	}
}
