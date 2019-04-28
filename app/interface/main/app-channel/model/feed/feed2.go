package feed

import (
	cardm "go-common/app/interface/main/app-card/model/card"
)

type Show2 struct {
	Topic cardm.Handler   `json:"topic,omitempty"`
	Feed  []cardm.Handler `json:"feed"`
}

type Tab struct {
	Items []cardm.Handler `json:"items"`
}
