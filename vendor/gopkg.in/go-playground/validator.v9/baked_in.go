package validator

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"
)

// Func accepts a FieldLevel interface for all validation needs. The return
// value should be true when validation succeeds.
type Func func(fl FieldLevel) bool

// FuncCtx accepts a context.Context and FieldLevel interface for all
// validation needs. The return value should be true when validation succeeds.
type FuncCtx func(ctx context.Context, fl FieldLevel) bool

// wrapFunc wraps noramal Func makes it compatible with FuncCtx
func wrapFunc(fn Func) FuncCtx {
	if fn == nil {
		return nil // be sure not to wrap a bad function.
	}
	return func(ctx context.Context, fl FieldLevel) bool {
		return fn(fl)
	}
}

var (
	restrictedTags = map[string]struct{}{
		diveTag:           {},
		keysTag:           {},
		endKeysTag:        {},
		structOnlyTag:     {},
		omitempty:         {},
		skipValidationTag: {},
		utf8HexComma:      {},
		utf8Pipe:          {},
		noStructLevelTag:  {},
		requiredTag:       {},
		isdefault:         {},
	}

	// BakedInAliasValidators is a default mapping of a single validation tag that
	// defines a common or complex set of validation(s) to simplify
	// adding validation to structs.
	bakedInAliases = map[string]string{
		"iscolor": "hexcolor|rgb|rgba|hsl|hsla",
	}

	// BakedInValidators is the default map of ValidationFunc
	// you can add, remove or even replace items to suite your needs,
	// or even disregard and use your own map if so desired.
	bakedInValidators = map[string]Func{
		"required":        hasValue,
		"isdefault":       isDefault,
		"len":             hasLengthOf,
		"min":             hasMinOf,
		"max":             hasMaxOf,
		"eq":              isEq,
		"ne":              isNe,
		"lt":              isLt,
		"lte":             isLte,
		"gt":              isGt,
		"gte":             isGte,
		"eqfield":         isEqField,
		"eqcsfield":       isEqCrossStructField,
		"necsfield":       isNeCrossStructField,
		"gtcsfield":       isGtCrossStructField,
		"gtecsfield":      isGteCrossStructField,
		"ltcsfield":       isLtCrossStructField,
		"ltecsfield":      isLteCrossStructField,
		"nefield":         isNeField,
		"gtefield":        isGteField,
		"gtfield":         isGtField,
		"ltefield":        isLteField,
		"ltfield":         isLtField,
		"alpha":           isAlpha,
		"alphanum":        isAlphanum,
		"alphaunicode":    isAlphaUnicode,
		"alphanumunicode": isAlphanumUnicode,
		"numeric":         isNumeric,
		"number":          isNumber,
		"hexadecimal":     isHexadecimal,
		"hexcolor":        isHEXColor,
		"rgb":             isRGB,
		"rgba":            isRGBA,
		"hsl":             isHSL,
		"hsla":            isHSLA,
		"email":           isEmail,
		"url":             isURL,
		"uri":             isURI,
		"base64":          isBase64,
		"contains":        contains,
		"containsany":     containsAny,
		"containsrune":    containsRune,
		"excludes":        excludes,
		"excludesall":     excludesAll,
		"excludesrune":    excludesRune,
		"isbn":            isISBN,
		"isbn10":          isISBN10,
		"isbn13":          isISBN13,
		"uuid":            isUUID,
		"uuid3":           isUUID3,
		"uuid4":           isUUID4,
		"uuid5":           isUUID5,
		"ascii":           isASCII,
		"printascii":      isPrintableASCII,
		"multibyte":       hasMultiByteCharacter,
		"datauri":         isDataURI,
		"latitude":        isLatitude,
		"longitude":       isLongitude,
		"ssn":             isSSN,
		"ipv4":            isIPv4,
		"ipv6":            isIPv6,
		"ip":              isIP,
		"cidrv4":          isCIDRv4,
		"cidrv6":          isCIDRv6,
		"cidr":            isCIDR,
		"tcp4_addr":       isTCP4AddrResolvable,
		"tcp6_addr":       isTCP6AddrResolvable,
		"tcp_addr":        isTCPAddrResolvable,
		"udp4_addr":       isUDP4AddrResolvable,
		"udp6_addr":       isUDP6AddrResolvable,
		"udp_addr":        isUDPAddrResolvable,
		"ip4_addr":        isIP4AddrResolvable,
		"ip6_addr":        isIP6AddrResolvable,
		"ip_addr":         isIPAddrResolvable,
		"unix_addr":       isUnixAddrResolvable,
		"mac":             isMAC,
		"hostname":        isHostname,
		"fqdn":            isFQDN,
		"unique":          isUnique,
	}
)

// isUnique is the validation function for validating if each array|slice element is unique
func isUnique(fl FieldLevel) bool {

	field := fl.Field()
	v := reflect.ValueOf(struct{}{})

	switch field.Kind() {
	case reflect.Slice, reflect.Array:
		m := reflect.MakeMap(reflect.MapOf(fl.Field().Type().Elem(), v.Type()))

		for i := 0; i < field.Len(); i++ {
			m.SetMapIndex(field.Index(i), v)
		}
		return field.Len() == m.Len()
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}
}

// IsMAC is the validation function for validating if the field's value is a valid MAC address.
func isMAC(fl FieldLevel) bool {

	_, err := net.ParseMAC(fl.Field().String())

	return err == nil
}

// IsCIDRv4 is the validation function for validating if the field's value is a valid v4 CIDR address.
func isCIDRv4(fl FieldLevel) bool {

	ip, _, err := net.ParseCIDR(fl.Field().String())

	return err == nil && ip.To4() != nil
}

// IsCIDRv6 is the validation function for validating if the field's value is a valid v6 CIDR address.
func isCIDRv6(fl FieldLevel) bool {

	ip, _, err := net.ParseCIDR(fl.Field().String())

	return err == nil && ip.To4() == nil
}

// IsCIDR is the validation function for validating if the field's value is a valid v4 or v6 CIDR address.
func isCIDR(fl FieldLevel) bool {

	_, _, err := net.ParseCIDR(fl.Field().String())

	return err == nil
}

// IsIPv4 is the validation function for validating if a value is a valid v4 IP address.
func isIPv4(fl FieldLevel) bool {

	ip := net.ParseIP(fl.Field().String())

	return ip != nil && ip.To4() != nil
}

// IsIPv6 is the validation function for validating if the field's value is a valid v6 IP address.
func isIPv6(fl FieldLevel) bool {

	ip := net.ParseIP(fl.Field().String())

	return ip != nil && ip.To4() == nil
}

// IsIP is the validation function for validating if the field's value is a valid v4 or v6 IP address.
func isIP(fl FieldLevel) bool {

	ip := net.ParseIP(fl.Field().String())

	return ip != nil
}

// IsSSN is the validation function for validating if the field's value is a valid SSN.
func isSSN(fl FieldLevel) bool {

	field := fl.Field()

	if field.Len() != 11 {
		return false
	}

	return sSNRegex.MatchString(field.String())
}

// IsLongitude is the validation function for validating if the field's value is a valid longitude coordinate.
func isLongitude(fl FieldLevel) bool {
	return longitudeRegex.MatchString(fl.Field().String())
}

// IsLatitude is the validation function for validating if the field's value is a valid latitude coordinate.
func isLatitude(fl FieldLevel) bool {
	return latitudeRegex.MatchString(fl.Field().String())
}

// IsDataURI is the validation function for validating if the field's value is a valid data URI.
func isDataURI(fl FieldLevel) bool {

	uri := strings.SplitN(fl.Field().String(), ",", 2)

	if len(uri) != 2 {
		return false
	}

	if !dataURIRegex.MatchString(uri[0]) {
		return false
	}

	return base64Regex.MatchString(uri[1])
}

// HasMultiByteCharacter is the validation function for validating if the field's value has a multi byte character.
func hasMultiByteCharacter(fl FieldLevel) bool {

	field := fl.Field()

	if field.Len() == 0 {
		return true
	}

	return multibyteRegex.MatchString(field.String())
}

// IsPrintableASCII is the validation function for validating if the field's value is a valid printable ASCII character.
func isPrintableASCII(fl FieldLevel) bool {
	return printableASCIIRegex.MatchString(fl.Field().String())
}

// IsASCII is the validation function for validating if the field's value is a valid ASCII character.
func isASCII(fl FieldLevel) bool {
	return aSCIIRegex.MatchString(fl.Field().String())
}

// IsUUID5 is the validation function for validating if the field's value is a valid v5 UUID.
func isUUID5(fl FieldLevel) bool {
	return uUID5Regex.MatchString(fl.Field().String())
}

// IsUUID4 is the validation function for validating if the field's value is a valid v4 UUID.
func isUUID4(fl FieldLevel) bool {
	return uUID4Regex.MatchString(fl.Field().String())
}

// IsUUID3 is the validation function for validating if the field's value is a valid v3 UUID.
func isUUID3(fl FieldLevel) bool {
	return uUID3Regex.MatchString(fl.Field().String())
}

// IsUUID is the validation function for validating if the field's value is a valid UUID of any version.
func isUUID(fl FieldLevel) bool {
	return uUIDRegex.MatchString(fl.Field().String())
}

// IsISBN is the validation function for validating if the field's value is a valid v10 or v13 ISBN.
func isISBN(fl FieldLevel) bool {
	return isISBN10(fl) || isISBN13(fl)
}

// IsISBN13 is the validation function for validating if the field's value is a valid v13 ISBN.
func isISBN13(fl FieldLevel) bool {

	s := strings.Replace(strings.Replace(fl.Field().String(), "-", "", 4), " ", "", 4)

	if !iSBN13Regex.MatchString(s) {
		return false
	}

	var checksum int32
	var i int32

	factor := []int32{1, 3}

	for i = 0; i < 12; i++ {
		checksum += factor[i%2] * int32(s[i]-'0')
	}

	return (int32(s[12]-'0'))-((10-(checksum%10))%10) == 0
}

// IsISBN10 is the validation function for validating if the field's value is a valid v10 ISBN.
func isISBN10(fl FieldLevel) bool {

	s := strings.Replace(strings.Replace(fl.Field().String(), "-", "", 3), " ", "", 3)

	if !iSBN10Regex.MatchString(s) {
		return false
	}

	var checksum int32
	var i int32

	for i = 0; i < 9; i++ {
		checksum += (i + 1) * int32(s[i]-'0')
	}

	if s[9] == 'X' {
		checksum += 10 * 10
	} else {
		checksum += 10 * int32(s[9]-'0')
	}

	return checksum%11 == 0
}

// ExcludesRune is the validation function for validating that the field's value does not contain the rune specified within the param.
func excludesRune(fl FieldLevel) bool {
	return !containsRune(fl)
}

// ExcludesAll is the validation function for validating that the field's value does not contain any of the characters specified within the param.
func excludesAll(fl FieldLevel) bool {
	return !containsAny(fl)
}

// Excludes is the validation function for validating that the field's value does not contain the text specified within the param.
func excludes(fl FieldLevel) bool {
	return !contains(fl)
}

// ContainsRune is the validation function for validating that the field's value contains the rune specified within the param.
func containsRune(fl FieldLevel) bool {

	r, _ := utf8.DecodeRuneInString(fl.Param())

	return strings.ContainsRune(fl.Field().String(), r)
}

// ContainsAny is the validation function for validating that the field's value contains any of the characters specified within the param.
func containsAny(fl FieldLevel) bool {
	return strings.ContainsAny(fl.Field().String(), fl.Param())
}

// Contains is the validation function for validating that the field's value contains the text specified within the param.
func contains(fl FieldLevel) bool {
	return strings.Contains(fl.Field().String(), fl.Param())
}

// IsNeField is the validation function for validating if the current field's value is not equal to the field specified by the param's value.
func isNeField(fl FieldLevel) bool {

	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()

	if !ok || currentKind != kind {
		return true
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() != currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() != currentField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() != currentField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) != int64(currentField.Len())

	case reflect.Struct:

		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() {
			return true
		}

		if fieldType == timeType {

			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)

			return !fieldTime.Equal(t)
		}

	}

	// default reflect.String:
	return field.String() != currentField.String()
}

// IsNe is the validation function for validating that the field's value does not equal the provided param value.
func isNe(fl FieldLevel) bool {
	return !isEq(fl)
}

// IsLteCrossStructField is the validation function for validating if the current field's value is less than or equal to the field, within a separate struct, specified by the param's value.
func isLteCrossStructField(fl FieldLevel) bool {

	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() <= topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() <= topField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() <= topField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) <= int64(topField.Len())

	case reflect.Struct:

		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() {
			return false
		}

		if fieldType == timeType {

			fieldTime := field.Interface().(time.Time)
			topTime := topField.Interface().(time.Time)

			return fieldTime.Before(topTime) || fieldTime.Equal(topTime)
		}
	}

	// default reflect.String:
	return field.String() <= topField.String()
}

// IsLtCrossStructField is the validation function for validating if the current field's value is less than the field, within a separate struct, specified by the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func isLtCrossStructField(fl FieldLevel) bool {

	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() < topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() < topField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() < topField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) < int64(topField.Len())

	case reflect.Struct:

		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() {
			return false
		}

		if fieldType == timeType {

			fieldTime := field.Interface().(time.Time)
			topTime := topField.Interface().(time.Time)

			return fieldTime.Before(topTime)
		}
	}

	// default reflect.String:
	return field.String() < topField.String()
}

// IsGteCrossStructField is the validation function for validating if the current field's value is greater than or equal to the field, within a separate struct, specified by the param's value.
func isGteCrossStructField(fl FieldLevel) bool {

	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() >= topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() >= topField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() >= topField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) >= int64(topField.Len())

	case reflect.Struct:

		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() {
			return false
		}

		if fieldType == timeType {

			fieldTime := field.Interface().(time.Time)
			topTime := topField.Interface().(time.Time)

			return fieldTime.After(topTime) || fieldTime.Equal(topTime)
		}
	}

	// default reflect.String:
	return field.String() >= topField.String()
}

// IsGtCrossStructField is the validation function for validating if the current field's value is greater than the field, within a separate struct, specified by the param's value.
func isGtCrossStructField(fl FieldLevel) bool {

	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() > topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() > topField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() > topField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) > int64(topField.Len())

	case reflect.Struct:

		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() {
			return false
		}

		if fieldType == timeType {

			fieldTime := field.Interface().(time.Time)
			topTime := topField.Interface().(time.Time)

			return fieldTime.After(topTime)
		}
	}

	// default reflect.String:
	return field.String() > topField.String()
}

// IsNeCrossStructField is the validation function for validating that the current field's value is not equal to the field, within a separate struct, specified by the param's value.
func isNeCrossStructField(fl FieldLevel) bool {

	field := fl.Field()
	kind := field.Kind()

	topField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return true
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return topField.Int() != field.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return topField.Uint() != field.Uint()

	case reflect.Float32, reflect.Float64:
		return topField.Float() != field.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(topField.Len()) != int64(field.Len())

	case reflect.Struct:

		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() {
			return true
		}

		if fieldType == timeType {

			t := field.Interface().(time.Time)
			fieldTime := topField.Interface().(time.Time)

			return !fieldTime.Equal(t)
		}
	}

	// default reflect.String:
	return topField.String() != field.String()
}

// IsEqCrossStructField is the validation function for validating that the current field's value is equal to the field, within a separate struct, specified by the param's value.
func isEqCrossStructField(fl FieldLevel) bool {

	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return topField.Int() == field.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return topField.Uint() == field.Uint()

	case reflect.Float32, reflect.Float64:
		return topField.Float() == field.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(topField.Len()) == int64(field.Len())

	case reflect.Struct:

		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() {
			return false
		}

		if fieldType == timeType {

			t := field.Interface().(time.Time)
			fieldTime := topField.Interface().(time.Time)

			return fieldTime.Equal(t)
		}
	}

	// default reflect.String:
	return topField.String() == field.String()
}

// IsEqField is the validation function for validating if the current field's value is equal to the field specified by the param's value.
func isEqField(fl FieldLevel) bool {

	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() == currentField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() == currentField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) == int64(currentField.Len())

	case reflect.Struct:

		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() {
			return false
		}

		if fieldType == timeType {

			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)

			return fieldTime.Equal(t)
		}

	}

	// default reflect.String:
	return field.String() == currentField.String()
}

// IsEq is the validation function for validating if the current field's value is equal to the param's value.
func isEq(fl FieldLevel) bool {

	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {

	case reflect.String:
		return field.String() == param

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) == p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return field.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() == p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return field.Float() == p
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// IsBase64 is the validation function for validating if the current field's value is a valid base 64.
func isBase64(fl FieldLevel) bool {
	return base64Regex.MatchString(fl.Field().String())
}

// IsURI is the validation function for validating if the current field's value is a valid URI.
func isURI(fl FieldLevel) bool {

	field := fl.Field()

	switch field.Kind() {

	case reflect.String:

		s := field.String()

		// checks needed as of Go 1.6 because of change https://github.com/golang/go/commit/617c93ce740c3c3cc28cdd1a0d712be183d0b328#diff-6c2d018290e298803c0c9419d8739885L195
		// emulate browser and strip the '#' suffix prior to validation. see issue-#237
		if i := strings.Index(s, "#"); i > -1 {
			s = s[:i]
		}

		if len(s) == 0 {
			return false
		}

		_, err := url.ParseRequestURI(s)

		return err == nil
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// IsURL is the validation function for validating if the current field's value is a valid URL.
func isURL(fl FieldLevel) bool {

	field := fl.Field()

	switch field.Kind() {

	case reflect.String:

		var i int
		s := field.String()

		// checks needed as of Go 1.6 because of change https://github.com/golang/go/commit/617c93ce740c3c3cc28cdd1a0d712be183d0b328#diff-6c2d018290e298803c0c9419d8739885L195
		// emulate browser and strip the '#' suffix prior to validation. see issue-#237
		if i = strings.Index(s, "#"); i > -1 {
			s = s[:i]
		}

		if len(s) == 0 {
			return false
		}

		url, err := url.ParseRequestURI(s)

		if err != nil || url.Scheme == "" {
			return false
		}

		return err == nil
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// IsEmail is the validation function for validating if the current field's value is a valid email address.
func isEmail(fl FieldLevel) bool {
	return emailRegex.MatchString(fl.Field().String())
}

// IsHSLA is the validation function for validating if the current field's value is a valid HSLA color.
func isHSLA(fl FieldLevel) bool {
	return hslaRegex.MatchString(fl.Field().String())
}

// IsHSL is the validation function for validating if the current field's value is a valid HSL color.
func isHSL(fl FieldLevel) bool {
	return hslRegex.MatchString(fl.Field().String())
}

// IsRGBA is the validation function for validating if the current field's value is a valid RGBA color.
func isRGBA(fl FieldLevel) bool {
	return rgbaRegex.MatchString(fl.Field().String())
}

// IsRGB is the validation function for validating if the current field's value is a valid RGB color.
func isRGB(fl FieldLevel) bool {
	return rgbRegex.MatchString(fl.Field().String())
}

// IsHEXColor is the validation function for validating if the current field's value is a valid HEX color.
func isHEXColor(fl FieldLevel) bool {
	return hexcolorRegex.MatchString(fl.Field().String())
}

// IsHexadecimal is the validation function for validating if the current field's value is a valid hexadecimal.
func isHexadecimal(fl FieldLevel) bool {
	return hexadecimalRegex.MatchString(fl.Field().String())
}

// IsNumber is the validation function for validating if the current field's value is a valid number.
func isNumber(fl FieldLevel) bool {
	return numberRegex.MatchString(fl.Field().String())
}

// IsNumeric is the validation function for validating if the current field's value is a valid numeric value.
func isNumeric(fl FieldLevel) bool {
	return numericRegex.MatchString(fl.Field().String())
}

// IsAlphanum is the validation function for validating if the current field's value is a valid alphanumeric value.
func isAlphanum(fl FieldLevel) bool {
	return alphaNumericRegex.MatchString(fl.Field().String())
}

// IsAlpha is the validation function for validating if the current field's value is a valid alpha value.
func isAlpha(fl FieldLevel) bool {
	return alphaRegex.MatchString(fl.Field().String())
}

// IsAlphanumUnicode is the validation function for validating if the current field's value is a valid alphanumeric unicode value.
func isAlphanumUnicode(fl FieldLevel) bool {
	return alphaUnicodeNumericRegex.MatchString(fl.Field().String())
}

// IsAlphaUnicode is the validation function for validating if the current field's value is a valid alpha unicode value.
func isAlphaUnicode(fl FieldLevel) bool {
	return alphaUnicodeRegex.MatchString(fl.Field().String())
}

// isDefault is the opposite of required aka hasValue
func isDefault(fl FieldLevel) bool {
	return !hasValue(fl)
}

// HasValue is the validation function for validating if the current field's value is not the default static value.
func hasValue(fl FieldLevel) bool {

	field := fl.Field()

	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:

		if fl.(*validate).fldIsPointer && field.Interface() != nil {
			return true
		}

		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

// IsGteField is the validation function for validating if the current field's value is greater than or equal to the field specified by the param's value.
func isGteField(fl FieldLevel) bool {

	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return field.Int() >= currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return field.Uint() >= currentField.Uint()

	case reflect.Float32, reflect.Float64:

		return field.Float() >= currentField.Float()

	case reflect.Struct:

		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() {
			return false
		}

		if fieldType == timeType {

			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)

			return fieldTime.After(t) || fieldTime.Equal(t)
		}
	}

	// default reflect.String
	return len(field.String()) >= len(currentField.String())
}

// IsGtField is the validation function for validating if the current field's value is greater than the field specified by the param's value.
func isGtField(fl FieldLevel) bool {

	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return field.Int() > currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return field.Uint() > currentField.Uint()

	case reflect.Float32, reflect.Float64:

		return field.Float() > currentField.Float()

	case reflect.Struct:

		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() {
			return false
		}

		if fieldType == timeType {

			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)

			return fieldTime.After(t)
		}
	}

	// default reflect.String
	return len(field.String()) > len(currentField.String())
}

// IsGte is the validation function for validating if the current field's value is greater than or equal to the param's value.
func isGte(fl FieldLevel) bool {

	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {

	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) >= p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) >= p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return field.Int() >= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() >= p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return field.Float() >= p

	case reflect.Struct:

		if field.Type() == timeType {

			now := time.Now().UTC()
			t := field.Interface().(time.Time)

			return t.After(now) || t.Equal(now)
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// IsGt is the validation function for validating if the current field's value is greater than the param's value.
func isGt(fl FieldLevel) bool {

	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {

	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) > p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) > p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return field.Int() > p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() > p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return field.Float() > p
	case reflect.Struct:

		if field.Type() == timeType {

			return field.Interface().(time.Time).After(time.Now().UTC())
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// HasLengthOf is the validation function for validating if the current field's value is equal to the param's value.
func hasLengthOf(fl FieldLevel) bool {

	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {

	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) == p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) == p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return field.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() == p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return field.Float() == p
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// HasMinOf is the validation function for validating if the current field's value is greater than or equal to the param's value.
func hasMinOf(fl FieldLevel) bool {
	return isGte(fl)
}

// IsLteField is the validation function for validating if the current field's value is less than or equal to the field specified by the param's value.
func isLteField(fl FieldLevel) bool {

	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return field.Int() <= currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return field.Uint() <= currentField.Uint()

	case reflect.Float32, reflect.Float64:

		return field.Float() <= currentField.Float()

	case reflect.Struct:

		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() {
			return false
		}

		if fieldType == timeType {

			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)

			return fieldTime.Before(t) || fieldTime.Equal(t)
		}
	}

	// default reflect.String
	return len(field.String()) <= len(currentField.String())
}

// IsLtField is the validation function for validating if the current field's value is less than the field specified by the param's value.
func isLtField(fl FieldLevel) bool {

	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return field.Int() < currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return field.Uint() < currentField.Uint()

	case reflect.Float32, reflect.Float64:

		return field.Float() < currentField.Float()

	case reflect.Struct:

		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() {
			return false
		}

		if fieldType == timeType {

			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)

			return fieldTime.Before(t)
		}
	}

	// default reflect.String
	return len(field.String()) < len(currentField.String())
}

// IsLte is the validation function for validating if the current field's value is less than or equal to the param's value.
func isLte(fl FieldLevel) bool {

	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {

	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) <= p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) <= p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return field.Int() <= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() <= p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return field.Float() <= p

	case reflect.Struct:

		if field.Type() == timeType {

			now := time.Now().UTC()
			t := field.Interface().(time.Time)

			return t.Before(now) || t.Equal(now)
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// IsLt is the validation function for validating if the current field's value is less than the param's value.
func isLt(fl FieldLevel) bool {

	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {

	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) < p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) < p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return field.Int() < p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() < p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return field.Float() < p

	case reflect.Struct:

		if field.Type() == timeType {

			return field.Interface().(time.Time).Before(time.Now().UTC())
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// HasMaxOf is the validation function for validating if the current field's value is less than or equal to the param's value.
func hasMaxOf(fl FieldLevel) bool {
	return isLte(fl)
}

// IsTCP4AddrResolvable is the validation function for validating if the field's value is a resolvable tcp4 address.
func isTCP4AddrResolvable(fl FieldLevel) bool {

	if !isIP4Addr(fl) {
		return false
	}

	_, err := net.ResolveTCPAddr("tcp4", fl.Field().String())
	return err == nil
}

// IsTCP6AddrResolvable is the validation function for validating if the field's value is a resolvable tcp6 address.
func isTCP6AddrResolvable(fl FieldLevel) bool {

	if !isIP6Addr(fl) {
		return false
	}

	_, err := net.ResolveTCPAddr("tcp6", fl.Field().String())

	return err == nil
}

// IsTCPAddrResolvable is the validation function for validating if the field's value is a resolvable tcp address.
func isTCPAddrResolvable(fl FieldLevel) bool {

	if !isIP4Addr(fl) && !isIP6Addr(fl) {
		return false
	}

	_, err := net.ResolveTCPAddr("tcp", fl.Field().String())

	return err == nil
}

// IsUDP4AddrResolvable is the validation function for validating if the field's value is a resolvable udp4 address.
func isUDP4AddrResolvable(fl FieldLevel) bool {

	if !isIP4Addr(fl) {
		return false
	}

	_, err := net.ResolveUDPAddr("udp4", fl.Field().String())

	return err == nil
}

// IsUDP6AddrResolvable is the validation function for validating if the field's value is a resolvable udp6 address.
func isUDP6AddrResolvable(fl FieldLevel) bool {

	if !isIP6Addr(fl) {
		return false
	}

	_, err := net.ResolveUDPAddr("udp6", fl.Field().String())

	return err == nil
}

// IsUDPAddrResolvable is the validation function for validating if the field's value is a resolvable udp address.
func isUDPAddrResolvable(fl FieldLevel) bool {

	if !isIP4Addr(fl) && !isIP6Addr(fl) {
		return false
	}

	_, err := net.ResolveUDPAddr("udp", fl.Field().String())

	return err == nil
}

// IsIP4AddrResolvable is the validation function for validating if the field's value is a resolvable ip4 address.
func isIP4AddrResolvable(fl FieldLevel) bool {

	if !isIPv4(fl) {
		return false
	}

	_, err := net.ResolveIPAddr("ip4", fl.Field().String())

	return err == nil
}

// IsIP6AddrResolvable is the validation function for validating if the field's value is a resolvable ip6 address.
func isIP6AddrResolvable(fl FieldLevel) bool {

	if !isIPv6(fl) {
		return false
	}

	_, err := net.ResolveIPAddr("ip6", fl.Field().String())

	return err == nil
}

// IsIPAddrResolvable is the validation function for validating if the field's value is a resolvable ip address.
func isIPAddrResolvable(fl FieldLevel) bool {

	if !isIP(fl) {
		return false
	}

	_, err := net.ResolveIPAddr("ip", fl.Field().String())

	return err == nil
}

// IsUnixAddrResolvable is the validation function for validating if the field's value is a resolvable unix address.
func isUnixAddrResolvable(fl FieldLevel) bool {

	_, err := net.ResolveUnixAddr("unix", fl.Field().String())

	return err == nil
}

func isIP4Addr(fl FieldLevel) bool {

	val := fl.Field().String()

	if idx := strings.LastIndex(val, ":"); idx != -1 {
		val = val[0:idx]
	}

	ip := net.ParseIP(val)

	return ip != nil && ip.To4() != nil
}

func isIP6Addr(fl FieldLevel) bool {

	val := fl.Field().String()

	if idx := strings.LastIndex(val, ":"); idx != -1 {
		if idx != 0 && val[idx-1:idx] == "]" {
			val = val[1 : idx-1]
		}
	}

	ip := net.ParseIP(val)

	return ip != nil && ip.To4() == nil
}

func isHostname(fl FieldLevel) bool {
	return hostnameRegex.MatchString(fl.Field().String())
}

func isFQDN(fl FieldLevel) bool {
	val := fl.Field().String()

	if val == "" {
		return false
	}

	if val[len(val)-1] == '.' {
		val = val[0 : len(val)-1]
	}

	return (strings.IndexAny(val, ".") > -1) &&
		hostnameRegex.MatchString(val)
}
