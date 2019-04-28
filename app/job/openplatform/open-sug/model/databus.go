package model

import (
	"encoding/json"
	"strings"
)

// Message ...
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// Season ...
type Season struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	SeasonTitle string `json:"season_title"`
	Mode        int    `json:"mode"`
	Type        int    `json:"type"`
	Alias       string `json:"alias"`
	AliasSearch string `json:"alias_search"`
	Brief       string `json:"brief"`
	Evaluate    string `json:"evaluate"`
	Actors      string `json:"actors"`
	Staff       string `json:"staff"`
	SquareCover string `json:"square_cover"`
	Cover       string `json:"cover"`
	EpCover     string `json:"epcover"`
	Area        int    `json:"area"`
	Ctime       string `json:"ctime"`
	Mtime       string `json:"mtime"`
}

// EsSeason ...
type EsSeason struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Alias       string   `json:"alias"`
	AliasSearch []string `json:"alias_search"`
	Actors      []string `json:"actors"`
}

//EsFormat ...
func (s *Season) EsFormat() (es *EsSeason) {
	es = &EsSeason{
		ID:          s.ID,
		Title:       s.Title,
		Alias:       s.Alias,
		AliasSearch: make([]string, 0),
		Actors:      make([]string, 0),
	}

	if s.AliasSearch != "" {
		es.AliasSearch = strings.Split(s.AliasSearch, ",")
	}

	if s.Actors != "" {
		for _, act := range strings.Split(s.Actors, "\n") {
			if act == "" {
				continue
			}

			act := strings.Split(act, "：")
			if len(act) > 0 {
				es.Actors = append(es.Actors, act[0])
			}
		}
	}
	return
}

// FieldDiff 检查指定字段是有变化
func (s *Season) FieldDiff(season *Season) bool {
	return s.Title != season.Title ||
		s.Alias != season.Alias ||
		s.AliasSearch != season.AliasSearch ||
		s.Actors != season.Actors ||
		s.Mtime != season.Mtime
}
