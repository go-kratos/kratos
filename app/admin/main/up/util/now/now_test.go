package now

import (
	"testing"
	"time"
)

var (
	format          = "2006-01-02 15:04:05.999999999"
	locationCaracas *time.Location
	locationBerlin  *time.Location
	timeCaracas     time.Time
)

func init() {
	var err error
	if locationCaracas, err = time.LoadLocation("America/Caracas"); err != nil {
		panic(err)
	}

	if locationBerlin, err = time.LoadLocation("Europe/Berlin"); err != nil {
		panic(err)
	}

	timeCaracas = time.Date(2016, 1, 1, 12, 10, 0, 0, locationCaracas)
}

func assertT(t *testing.T) func(time.Time, string, string) {
	return func(actual time.Time, expected string, msg string) {
		actualStr := actual.Format(format)
		if actualStr != expected {
			t.Errorf("Failed %s: actual: %v, expected: %v", msg, actualStr, expected)
		}
	}
}

func TestBeginningOf(t *testing.T) {
	assert := assertT(t)

	n := time.Date(2013, 11, 18, 17, 51, 49, 123456789, time.UTC)

	assert(New(n).BeginningOfMinute(), "2013-11-18 17:51:00", "BeginningOfMinute")

	WeekStartDay = time.Monday
	assert(New(n).BeginningOfWeek(), "2013-11-18 00:00:00", "BeginningOfWeek, FirstDayMonday")

	WeekStartDay = time.Tuesday
	assert(New(n).BeginningOfWeek(), "2013-11-12 00:00:00", "BeginningOfWeek, FirstDayTuesday")

	WeekStartDay = time.Wednesday
	assert(New(n).BeginningOfWeek(), "2013-11-13 00:00:00", "BeginningOfWeek, FirstDayWednesday")

	WeekStartDay = time.Thursday
	assert(New(n).BeginningOfWeek(), "2013-11-14 00:00:00", "BeginningOfWeek, FirstDayThursday")

	WeekStartDay = time.Friday
	assert(New(n).BeginningOfWeek(), "2013-11-15 00:00:00", "BeginningOfWeek, FirstDayFriday")

	WeekStartDay = time.Saturday
	assert(New(n).BeginningOfWeek(), "2013-11-16 00:00:00", "BeginningOfWeek, FirstDaySaturday")

	WeekStartDay = time.Sunday
	assert(New(n).BeginningOfWeek(), "2013-11-17 00:00:00", "BeginningOfWeek, FirstDaySunday")

	assert(New(n).BeginningOfHour(), "2013-11-18 17:00:00", "BeginningOfHour")

	// Truncate with hour bug
	assert(New(timeCaracas).BeginningOfHour(), "2016-01-01 12:00:00", "BeginningOfHour Caracas")

	assert(New(n).BeginningOfDay(), "2013-11-18 00:00:00", "BeginningOfDay")

	location, err := time.LoadLocation("Japan")
	if err != nil {
		t.Fatalf("Error loading location: %v", err)
	}
	beginningOfDay := time.Date(2015, 05, 01, 0, 0, 0, 0, location)
	assert(New(beginningOfDay).BeginningOfDay(), "2015-05-01 00:00:00", "BeginningOfDay")

	// DST
	dstBeginningOfDay := time.Date(2017, 10, 29, 10, 0, 0, 0, locationBerlin)
	assert(New(dstBeginningOfDay).BeginningOfDay(), "2017-10-29 00:00:00", "BeginningOfDay DST")

	assert(New(n).BeginningOfWeek(), "2013-11-17 00:00:00", "BeginningOfWeek")

	dstBegginingOfWeek := time.Date(2017, 10, 30, 12, 0, 0, 0, locationBerlin)
	assert(New(dstBegginingOfWeek).BeginningOfWeek(), "2017-10-29 00:00:00", "BeginningOfWeek")

	dstBegginingOfWeek = time.Date(2017, 10, 29, 12, 0, 0, 0, locationBerlin)
	assert(New(dstBegginingOfWeek).BeginningOfWeek(), "2017-10-29 00:00:00", "BeginningOfWeek")

	WeekStartDay = time.Monday
	assert(New(n).BeginningOfWeek(), "2013-11-18 00:00:00", "BeginningOfWeek, FirstDayMonday")
	dstBegginingOfWeek = time.Date(2017, 10, 24, 12, 0, 0, 0, locationBerlin)
	assert(New(dstBegginingOfWeek).BeginningOfWeek(), "2017-10-23 00:00:00", "BeginningOfWeek, FirstDayMonday")

	dstBegginingOfWeek = time.Date(2017, 10, 29, 12, 0, 0, 0, locationBerlin)
	assert(New(dstBegginingOfWeek).BeginningOfWeek(), "2017-10-23 00:00:00", "BeginningOfWeek, FirstDayMonday")

	WeekStartDay = time.Sunday

	assert(New(n).BeginningOfMonth(), "2013-11-01 00:00:00", "BeginningOfMonth")

	// DST
	dstBeginningOfMonth := time.Date(2017, 10, 31, 0, 0, 0, 0, locationBerlin)
	assert(New(dstBeginningOfMonth).BeginningOfMonth(), "2017-10-01 00:00:00", "BeginningOfMonth DST")

	assert(New(n).BeginningOfQuarter(), "2013-10-01 00:00:00", "BeginningOfQuarter")

	// DST
	assert(New(dstBeginningOfMonth).BeginningOfQuarter(), "2017-10-01 00:00:00", "BeginningOfQuarter DST")
	dstBeginningOfQuarter := time.Date(2017, 11, 24, 0, 0, 0, 0, locationBerlin)
	assert(New(dstBeginningOfQuarter).BeginningOfQuarter(), "2017-10-01 00:00:00", "BeginningOfQuarter DST")

	assert(New(n.AddDate(0, -1, 0)).BeginningOfQuarter(), "2013-10-01 00:00:00", "BeginningOfQuarter")

	assert(New(n.AddDate(0, 1, 0)).BeginningOfQuarter(), "2013-10-01 00:00:00", "BeginningOfQuarter")

	// DST
	assert(New(dstBeginningOfQuarter).BeginningOfYear(), "2017-01-01 00:00:00", "BeginningOfYear DST")

	assert(New(timeCaracas).BeginningOfYear(), "2016-01-01 00:00:00", "BeginningOfYear Caracas")
}

func TestEndOf(t *testing.T) {
	assert := assertT(t)

	n := time.Date(2013, 11, 18, 17, 51, 49, 123456789, time.UTC)

	assert(New(n).EndOfMinute(), "2013-11-18 17:51:59.999999999", "EndOfMinute")

	assert(New(n).EndOfHour(), "2013-11-18 17:59:59.999999999", "EndOfHour")

	assert(New(timeCaracas).EndOfHour(), "2016-01-01 12:59:59.999999999", "EndOfHour Caracas")

	assert(New(n).EndOfDay(), "2013-11-18 23:59:59.999999999", "EndOfDay")

	dstEndOfDay := time.Date(2017, 10, 29, 1, 0, 0, 0, locationBerlin)
	assert(New(dstEndOfDay).EndOfDay(), "2017-10-29 23:59:59.999999999", "EndOfDay DST")

	WeekStartDay = time.Tuesday
	assert(New(n).EndOfWeek(), "2013-11-18 23:59:59.999999999", "EndOfWeek, FirstDayTuesday")

	WeekStartDay = time.Wednesday
	assert(New(n).EndOfWeek(), "2013-11-19 23:59:59.999999999", "EndOfWeek, FirstDayWednesday")

	WeekStartDay = time.Thursday
	assert(New(n).EndOfWeek(), "2013-11-20 23:59:59.999999999", "EndOfWeek, FirstDayThursday")

	WeekStartDay = time.Friday
	assert(New(n).EndOfWeek(), "2013-11-21 23:59:59.999999999", "EndOfWeek, FirstDayFriday")

	WeekStartDay = time.Saturday
	assert(New(n).EndOfWeek(), "2013-11-22 23:59:59.999999999", "EndOfWeek, FirstDaySaturday")

	WeekStartDay = time.Sunday
	assert(New(n).EndOfWeek(), "2013-11-23 23:59:59.999999999", "EndOfWeek, FirstDaySunday")

	WeekStartDay = time.Monday
	assert(New(n).EndOfWeek(), "2013-11-24 23:59:59.999999999", "EndOfWeek, FirstDayMonday")

	dstEndOfWeek := time.Date(2017, 10, 24, 12, 0, 0, 0, locationBerlin)
	assert(New(dstEndOfWeek).EndOfWeek(), "2017-10-29 23:59:59.999999999", "EndOfWeek, FirstDayMonday")

	dstEndOfWeek = time.Date(2017, 10, 29, 12, 0, 0, 0, locationBerlin)
	assert(New(dstEndOfWeek).EndOfWeek(), "2017-10-29 23:59:59.999999999", "EndOfWeek, FirstDayMonday")

	WeekStartDay = time.Sunday
	assert(New(n).EndOfWeek(), "2013-11-23 23:59:59.999999999", "EndOfWeek")

	dstEndOfWeek = time.Date(2017, 10, 29, 0, 0, 0, 0, locationBerlin)
	assert(New(dstEndOfWeek).EndOfWeek(), "2017-11-04 23:59:59.999999999", "EndOfWeek")

	dstEndOfWeek = time.Date(2017, 10, 29, 12, 0, 0, 0, locationBerlin)
	assert(New(dstEndOfWeek).EndOfWeek(), "2017-11-04 23:59:59.999999999", "EndOfWeek")

	assert(New(n).EndOfMonth(), "2013-11-30 23:59:59.999999999", "EndOfMonth")

	assert(New(n).EndOfQuarter(), "2013-12-31 23:59:59.999999999", "EndOfQuarter")

	assert(New(n.AddDate(0, -1, 0)).EndOfQuarter(), "2013-12-31 23:59:59.999999999", "EndOfQuarter")

	assert(New(n.AddDate(0, 1, 0)).EndOfQuarter(), "2013-12-31 23:59:59.999999999", "EndOfQuarter")

	assert(New(n).EndOfYear(), "2013-12-31 23:59:59.999999999", "EndOfYear")

	n1 := time.Date(2013, 02, 18, 17, 51, 49, 123456789, time.UTC)
	assert(New(n1).EndOfMonth(), "2013-02-28 23:59:59.999999999", "EndOfMonth for 2013/02")

	n2 := time.Date(1900, 02, 18, 17, 51, 49, 123456789, time.UTC)
	assert(New(n2).EndOfMonth(), "1900-02-28 23:59:59.999999999", "EndOfMonth")
}

func TestMondayAndSunday(t *testing.T) {
	assert := assertT(t)

	n := time.Date(2013, 11, 19, 17, 51, 49, 123456789, time.UTC)
	n2 := time.Date(2013, 11, 24, 17, 51, 49, 123456789, time.UTC)
	nDst := time.Date(2017, 10, 29, 10, 0, 0, 0, locationBerlin)

	assert(New(n).Monday(), "2013-11-18 00:00:00", "Monday")

	assert(New(n2).Monday(), "2013-11-18 00:00:00", "Monday")

	assert(New(timeCaracas).Monday(), "2015-12-28 00:00:00", "Monday Caracas")

	assert(New(nDst).Monday(), "2017-10-23 00:00:00", "Monday DST")

	assert(New(n).Sunday(), "2013-11-24 00:00:00", "Sunday")

	assert(New(n2).Sunday(), "2013-11-24 00:00:00", "Sunday")

	assert(New(timeCaracas).Sunday(), "2016-01-03 00:00:00", "Sunday Caracas")

	assert(New(nDst).Sunday(), "2017-10-29 00:00:00", "Sunday DST")

	assert(New(n).EndOfSunday(), "2013-11-24 23:59:59.999999999", "EndOfSunday")

	assert(New(timeCaracas).EndOfSunday(), "2016-01-03 23:59:59.999999999", "EndOfSunday Caracas")

	assert(New(nDst).EndOfSunday(), "2017-10-29 23:59:59.999999999", "EndOfSunday DST")

	assert(New(n).BeginningOfWeek(), "2013-11-17 00:00:00", "BeginningOfWeek, FirstDayMonday")

	WeekStartDay = time.Monday
	assert(New(n).BeginningOfWeek(), "2013-11-18 00:00:00", "BeginningOfWeek, FirstDayMonday")
}

func TestParse(t *testing.T) {
	assert := assertT(t)

	n := time.Date(2013, 11, 18, 17, 51, 49, 123456789, time.UTC)

	assert(New(n).MustParse("2002"), "2002-01-01 00:00:00", "Parse 2002")

	assert(New(n).MustParse("2002-10"), "2002-10-01 00:00:00", "Parse 2002-10")

	assert(New(n).MustParse("2002-10-12"), "2002-10-12 00:00:00", "Parse 2002-10-12")

	assert(New(n).MustParse("2002-10-12 22"), "2002-10-12 22:00:00", "Parse 2002-10-12 22")

	assert(New(n).MustParse("2002-10-12 22:14"), "2002-10-12 22:14:00", "Parse 2002-10-12 22:14")

	assert(New(n).MustParse("2002-10-12 2:4"), "2002-10-12 02:04:00", "Parse 2002-10-12 2:4")

	assert(New(n).MustParse("2002-10-12 02:04"), "2002-10-12 02:04:00", "Parse 2002-10-12 02:04")

	assert(New(n).MustParse("2002-10-12 22:14:56"), "2002-10-12 22:14:56", "Parse 2002-10-12 22:14:56")

	assert(New(n).MustParse("2002-10-12 00:14:56"), "2002-10-12 00:14:56", "Parse 2002-10-12 00:14:56")

	assert(New(n).MustParse("2013-12-19 23:28:09.999999999 +0800 CST"), "2013-12-19 23:28:09.999999999", "Parse two strings 2013-12-19 23:28:09.999999999 +0800 CST")

	assert(New(n).MustParse("10-12"), "2013-10-12 00:00:00", "Parse 10-12")

	assert(New(n).MustParse("18"), "2013-11-18 18:00:00", "Parse 18 as hour")

	assert(New(n).MustParse("18:20"), "2013-11-18 18:20:00", "Parse 18:20")

	assert(New(n).MustParse("00:01"), "2013-11-18 00:01:00", "Parse 00:01")

	assert(New(n).MustParse("00:00:00"), "2013-11-18 00:00:00", "Parse 00:00:00")

	assert(New(n).MustParse("18:20:39"), "2013-11-18 18:20:39", "Parse 18:20:39")

	assert(New(n).MustParse("18:20:39", "2011-01-01"), "2011-01-01 18:20:39", "Parse two strings 18:20:39, 2011-01-01")

	assert(New(n).MustParse("2011-1-1", "18:20:39"), "2011-01-01 18:20:39", "Parse two strings 2011-01-01, 18:20:39")

	assert(New(n).MustParse("2011-01-01", "18"), "2011-01-01 18:00:00", "Parse two strings 2011-01-01, 18")

	TimeFormats = append(TimeFormats, "02 Jan 15:04")
	assert(New(n).MustParse("04 Feb 12:09"), "2013-02-04 12:09:00", "Parse 04 Feb 12:09 with specified format")

	assert(New(n).MustParse("23:28:9 Dec 19, 2013 PST"), "2013-12-19 23:28:09", "Parse 23:28:9 Dec 19, 2013 PST")

	if New(n).MustParse("23:28:9 Dec 19, 2013 PST").Location().String() != "PST" {
		t.Errorf("Parse 23:28:9 Dec 19, 2013 PST shouldn't lose time zone")
	}

	n2 := New(n).MustParse("23:28:9 Dec 19, 2013 PST")
	if New(n2).MustParse("10:20").Location().String() != "PST" {
		t.Errorf("Parse 10:20 shouldn't change time zone")
	}

	TimeFormats = append(TimeFormats, "2006-01-02T15:04:05.0")
	if MustParseInLocation(time.UTC, "2018-02-13T15:17:06.0").String() != "2018-02-13 15:17:06 +0000 UTC" {
		t.Errorf("ParseInLocation 2018-02-13T15:17:06.0")
	}

	TimeFormats = append(TimeFormats, "2006-01-02 15:04:05.000")
	assert(New(n).MustParse("2018-04-20 21:22:23.473"), "2018-04-20 21:22:23.473", "Parse 2018/04/20 21:22:23.473")

	TimeFormats = append(TimeFormats, "15:04:05.000")
	assert(New(n).MustParse("13:00:01.365"), "2013-11-18 13:00:01.365", "Parse 13:00:01.365")

	TimeFormats = append(TimeFormats, "2006-01-02 15:04:05.000000")
	assert(New(n).MustParse("2010-01-01 07:24:23.131384"), "2010-01-01 07:24:23.131384", "Parse 2010-01-01 07:24:23.131384")
	assert(New(n).MustParse("00:00:00.182736"), "2013-11-18 00:00:00.182736", "Parse 00:00:00.182736")
}

func TestBetween(t *testing.T) {
	tm := time.Date(2015, 06, 30, 17, 51, 49, 123456789, time.Now().Location())
	if !New(tm).Between("23:28:9 Dec 19, 2013 PST", "23:28:9 Dec 19, 2015 PST") {
		t.Errorf("Between")
	}

	if !New(tm).Between("2015-05-12 12:20", "2015-06-30 17:51:50") {
		t.Errorf("Between")
	}
}

func Example() {
	time.Now() // 2013-11-18 17:51:49.123456789 Mon

	BeginningOfMinute() // 2013-11-18 17:51:00 Mon
	BeginningOfHour()   // 2013-11-18 17:00:00 Mon
	BeginningOfDay()    // 2013-11-18 00:00:00 Mon
	BeginningOfWeek()   // 2013-11-17 00:00:00 Sun

	WeekStartDay = time.Monday // Set Monday as first day
	BeginningOfWeek()          // 2013-11-18 00:00:00 Mon
	BeginningOfMonth()         // 2013-11-01 00:00:00 Fri
	BeginningOfQuarter()       // 2013-10-01 00:00:00 Tue
	BeginningOfYear()          // 2013-01-01 00:00:00 Tue

	EndOfMinute() // 2013-11-18 17:51:59.999999999 Mon
	EndOfHour()   // 2013-11-18 17:59:59.999999999 Mon
	EndOfDay()    // 2013-11-18 23:59:59.999999999 Mon
	EndOfWeek()   // 2013-11-23 23:59:59.999999999 Sat

	WeekStartDay = time.Monday // Set Monday as first day
	EndOfWeek()                // 2013-11-24 23:59:59.999999999 Sun
	EndOfMonth()               // 2013-11-30 23:59:59.999999999 Sat
	EndOfQuarter()             // 2013-12-31 23:59:59.999999999 Tue
	EndOfYear()                // 2013-12-31 23:59:59.999999999 Tue

	// Use another time
	t := time.Date(2013, 02, 18, 17, 51, 49, 123456789, time.UTC)
	New(t).EndOfMonth() // 2013-02-28 23:59:59.999999999 Thu

	Monday()      // 2013-11-18 00:00:00 Mon
	Sunday()      // 2013-11-24 00:00:00 Sun
	EndOfSunday() // 2013-11-24 23:59:59.999999999 Sun
}
