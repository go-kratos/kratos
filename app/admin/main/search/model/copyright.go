package model

import (
	"strconv"
)

// CopyRight .
type CopyRight struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	OName       string `json:"oname"`
	AkaNames    string `json:"aka_names"`
	Level       string `json:"level"`
	AVoid       string `json:"avoid"`
	Plan        string `json:"plan"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// IndexName .
func (c *CopyRight) IndexName() string {
	return "copyright"
}

// IndexType .
func (c *CopyRight) IndexType() string {
	return "base"
}

// IndexID .
func (c *CopyRight) IndexID() string {
	return strconv.FormatInt(c.ID, 10)
}
