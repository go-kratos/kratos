package experiment

type Experiment struct {
	ID           int64    `json:"id,omitempty"`
	Plat         int8     `json:"-"`
	Name         string   `json:"name,omitempty"`
	Strategy     string   `json:"strategy,omitempty"`
	Desc         string   `json:"desc,omitempty"`
	TrafficGroup string   `json:"traffic_group,omitempty"`
	Limit        []*Limit `json:"-"`
}

type Limit struct {
	ExperimentID int64  `json:"-"`
	Build        int    `json:"-"`
	Condition    string `json:"-"`
}

type ABTestV2 struct {
	AutoPlay int `json:"autoplay"`
}
