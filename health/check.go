package health

import (
	"context"
	"errors"
)

var ErrServiceNotFind = errors.New("service not find")

type Status string

const (
	Up   Status = "UP"
	Down Status = "DOWN"
)

type Checker interface {
	Check(ctx context.Context) error
}

type Result struct {
	Status  Status            `json:"status"`
	Details map[string]Detail `json:"details"`
}

type Detail struct {
	Status Status `json:"status"`
	Error  string `json:"error"`
}
