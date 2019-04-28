// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package shell implements calculate signature for shell(pay system) query parameters.
//
// As a simple example:
//
// 	type Options struct {
// 		Query   string `json:"q"`
// 		ShowAll bool   `json:"all"`
// 		Page    int    `json:"page,omitempty"`
// 	}
// use json tags
//
package shell

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})

//var encoderType = reflect.TypeOf(new(Encoder)).Elem()

// Encoder is an interface implemented by any type that wishes to encode
// itself into URL values in a non-standard way.
type Encoder interface {
	EncodeValues(key string, v *url.Values) error
}

//Sign sign for shell request
func Sign(v interface{}, token string) (sign string, err error) {
	var param string
	param, err = Encode(v)
	if err != nil {
		return
	}
	var final = param + "&token=" + token
	var h = md5.New()
	io.WriteString(h, final)
	var result = h.Sum(nil)
	sign = fmt.Sprintf("%x", result)
	return
}

// Encode encode with dictinary ascending order
func Encode(v interface{}) (res string, err error) {
	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return "", nil
		}
		val = val.Elem()
	}

	if v == nil {
		return "", nil
	}

	if val.Kind() != reflect.Struct {
		return "", fmt.Errorf("query: Values() expects struct input. Got %v", val.Kind())
	}

	res, err = encodeValue(val, tagOptions{}, 0)
	return
}

// encodeValue populates the values parameter from the struct fields in val.
// Embedded structs are followed recursively (using the rules defined in the
// Values function documentation) breadth-first.
func encodeValue(val reflect.Value, opts tagOptions, deep int) (res string, err error) {
	typ := val.Type()
	switch typ.Kind() {
	case reflect.Struct:
		return encodeStruct(val, deep+1)
	case reflect.Slice, reflect.Array:
		return encodeSlice(val, opts, deep+1)
	default:
		res = valueString(val, opts)
	}
	return
}

func encodeSlice(val reflect.Value, opts tagOptions, deep int) (res string, err error) {
	var del byte = ','
	if opts.Contains("comma") {
		del = ','
	} else if opts.Contains("space") {
		del = ' '
	} else if opts.Contains("semicolon") {
		del = ';'
	}

	s := new(bytes.Buffer)
	first := true
	for i := 0; i < val.Len(); i++ {
		if first {
			first = false
		} else {
			s.WriteByte(del)
		}
		var r string
		r, err = encodeValue(val.Index(i), opts, deep)
		if err != nil {
			return
		}
		s.WriteString(r)
	}
	return "[" + s.String() + "]", nil

}
func encodeStruct(val reflect.Value, deep int) (res string, err error) {
	var values = make(url.Values)
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if sf.PkgPath != "" && !sf.Anonymous { // unexported
			continue
		}

		sv := val.Field(i)
		tag := sf.Tag.Get("json")
		if tag == "-" {
			continue
		}
		name, opts := parseTag(tag)

		if opts.Contains("omitempty") && isEmptyValue(sv) {
			continue
		}

		for sv.Kind() == reflect.Ptr {
			if sv.IsNil() {
				break
			}
			sv = sv.Elem()
		}

		var r string
		r, err = encodeValue(sv, opts, deep)
		if err != nil {
			return
		}
		values.Add(name, r)
	}

	if deep > 1 {
		res = "{" + encode(values) + "}"
	} else {
		res = encode(values)
	}
	return res, nil
}

// valueString returns the string representation of a value.
func valueString(v reflect.Value, opts tagOptions) string {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.Bool && opts.Contains("int") {
		if v.Bool() {
			return "1"
		}
		return "0"
	}

	if v.Type() == timeType {
		t := v.Interface().(time.Time)
		if opts.Contains("unix") {
			return strconv.FormatInt(t.Unix(), 10)
		}
		return t.Format(time.RFC3339)
	}

	return fmt.Sprint(v.Interface())
}

// tagOptions is the string following a comma in a struct field's "url" tag, or
// the empty string. It does not include the leading comma.
type tagOptions []string

// parseTag splits a struct field's url tag into its name and comma-separated
// options.
func parseTag(tag string) (string, tagOptions) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}

// Contains checks whether the tagOptions contains the specified option.
func (o tagOptions) Contains(option string) bool {
	for _, s := range o {
		if s == option {
			return true
		}
	}
	return false
}

func encode(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := escape(k) + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(escape(v))
		}
	}
	return buf.String()
}

func escape(s string) string {
	return s
}

// isEmptyValue checks if a value should be considered empty for the purposes
// of omitting fields with the "omitempty" option.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	if v.Type() == timeType {
		return v.Interface().(time.Time).IsZero()
	}

	return false
}
