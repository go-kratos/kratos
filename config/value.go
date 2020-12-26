package config

import (
	"errors"
	"time"
)

// ErrNotFound is value not found.
var ErrNotFound = errors.New("error key not found")

// Value is config value interface.
type Value interface {
	Bool() (bool, error)
	Int() (int, error)
	Int32() (int32, error)
	Int64() (int64, error)
	Float32() (float32, error)
	Float64() (float64, error)
	Duration() (time.Duration, error)
	String() (string, error)
	Scan(interface{}) error
}

type errValue struct {
	err error
}

func (v errValue) Bool() (bool, error)              { return false, v.err }
func (v errValue) Int() (int, error)                { return 0, v.err }
func (v errValue) Int32() (int32, error)            { return 0, v.err }
func (v errValue) Int64() (int64, error)            { return 0, v.err }
func (v errValue) Float32() (float32, error)        { return 0, v.err }
func (v errValue) Float64() (float64, error)        { return 0, v.err }
func (v errValue) Duration() (time.Duration, error) { return 0, v.err }
func (v errValue) String() (string, error)          { return "", v.err }
func (v errValue) Scan(interface{}) error           { return v.err }
