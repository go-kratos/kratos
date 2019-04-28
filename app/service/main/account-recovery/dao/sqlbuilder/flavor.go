// Copyright 2018 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sqlbuilder

import "fmt"

// Supported flavors.
const (
	invalidFlavor Flavor = iota

	MySQL
	PostgreSQL
)

var (
	// DefaultFlavor is the default flavor for all builders.
	DefaultFlavor = MySQL
)

// Flavor is the flag to control the format of compiled sql.
type Flavor int

// String returns the name of f.
func (f Flavor) String() string {
	switch f {
	case MySQL:
		return "MySQL"
	case PostgreSQL:
		return "PostgreSQL"
	}

	return "<invalid>"
}

// NewSelectBuilder creates a new SELECT builder with flavor.
func (f Flavor) NewSelectBuilder() *SelectBuilder {
	b := newSelectBuilder()
	b.SetFlavor(f)
	return b
}

// Quote adds quote for name to make sure the name can be used safely
// as table name or field name.
//
// * For MySQL, use back quote (`) to quote name;
// * For PostgreSQL, use double quote (") to quote name.
func (f Flavor) Quote(name string) string {
	switch f {
	case MySQL:
		return fmt.Sprintf("`%v`", name)
	case PostgreSQL:
		return fmt.Sprintf(`"%v"`, name)
	}

	return name
}
