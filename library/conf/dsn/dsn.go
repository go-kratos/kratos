// Package dsn provide parse dsn and bind to struct
// see http://git.bilibili.co/platform/go-common/issues/279
package dsn

import (
	"net/url"
	"reflect"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

var _validator *validator.Validate

func init() {
	_validator = validator.New()
}

// DSN a DSN represents a parsed DSN as same as url.URL.
type DSN struct {
	*url.URL
}

// Bind dsn to specify struct and validate use use go-playground/validator format
//
// The bind of each struct field can be customized by the format string
// stored under the 'dsn' key in the struct field's tag. The format string
// gives the name of the field, possibly followed by a comma-separated
// list of options.  The name may be empty in order to specify options
// without overriding the default field name.
//
// A two type data you can bind to struct
// built-in values, use below keys to bind built-in value
//	username
//	password
//	address
//	network
// the value in query string, use query.{name} to bind value in query string
//
// As a special case, if the field tag is "-", the field is always omitted.
// NOTE: that a field with name "-" can still be generated using the tag "-,".
//
// Examples of struct field tags and their meanings:
//	// Field bind username
//	Field string `dsn:"username"`
//	// Field is ignored by this package.
//	Field string `dsn:"-"`
//	// Field bind value from query
//	Field string `dsn:"query.name"`
//
func (d *DSN) Bind(v interface{}) (url.Values, error) {
	assignFuncs := make(map[string]assignFunc)
	if d.User != nil {
		username := d.User.Username()
		password, ok := d.User.Password()
		if ok {
			assignFuncs["password"] = stringsAssignFunc(password)
		}
		assignFuncs["username"] = stringsAssignFunc(username)
	}
	assignFuncs["address"] = addressesAssignFunc(d.Addresses())
	assignFuncs["network"] = stringsAssignFunc(d.Scheme)
	query, err := bindQuery(d.Query(), v, assignFuncs)
	if err != nil {
		return nil, err
	}
	return query, _validator.Struct(v)
}

func addressesAssignFunc(addresses []string) assignFunc {
	return func(v reflect.Value, to tagOpt) error {
		if v.Kind() == reflect.String {
			if addresses[0] == "" && to.Default != "" {
				v.SetString(to.Default)
			} else {
				v.SetString(addresses[0])
			}
			return nil
		}
		if !(v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.String) {
			return &BindTypeError{Value: strings.Join(addresses, ","), Type: v.Type()}
		}
		vals := reflect.MakeSlice(v.Type(), len(addresses), len(addresses))
		for i, address := range addresses {
			vals.Index(i).SetString(address)
		}
		if v.CanSet() {
			v.Set(vals)
		}
		return nil
	}
}

// Addresses parse host split by ','
// For Unix networks, return ['path']
func (d *DSN) Addresses() []string {
	switch d.Scheme {
	case "unix", "unixgram", "unixpacket":
		return []string{d.Path}
	}
	return strings.Split(d.Host, ",")
}

// Parse parses rawdsn into a URL structure.
func Parse(rawdsn string) (*DSN, error) {
	u, err := url.Parse(rawdsn)
	return &DSN{URL: u}, err
}
