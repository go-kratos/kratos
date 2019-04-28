package model

import (
	"errors"
)

const (
	bit1 = int8(1)
	bit2 = int8(1) << 1

	// PubOnly pub only
	PubOnly = bit2 | int8(0)
	// SubOnly sub only
	SubOnly = int8(0) | bit1
	// PubSub pub and sub
	PubSub = bit2 | bit1
)

var (
	errGroup  = errors.New("error group")
	errTopic  = errors.New("error topic")
	errKey    = errors.New("error key")
	errSecret = errors.New("error secret")
)

// Auth databus auth info accordance with table:bilibili_databus_v2.auth
type Auth struct {
	Group     string
	Topic     string
	Operation int8
	Key       string
	Secret    string
	Batch     int64
	Cluster   string
}

// CanPub judge producer auth
func (a *Auth) CanPub() bool {
	return a.Operation&bit2 == bit2
}

// CanSub judge consumer auth
func (a *Auth) CanSub() bool {
	return a.Operation&bit1 == bit1
}

// Auth judge auth
func (a *Auth) Auth(group, topic, key, secret string) error {
	if a.Group != group {
		return errGroup
	}
	if a.Topic != topic {
		return errTopic
	}
	if a.Key != key {
		return errKey
	}
	if a.Secret != secret {
		return errSecret
	}
	return nil
}
