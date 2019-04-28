## Now

Now is a time toolkit for golang

[![wercker status](https://app.wercker.com/status/a350da4eae6cb28a35687ba41afb565a/s/master "wercker status")](https://app.wercker.com/project/byKey/a350da4eae6cb28a35687ba41afb565a)

## Install

```
go get -u github.com/jinzhu/now
```

## Usage

Calculating time based on current time

```go
import "github.com/jinzhu/now"

time.Now() // 2013-11-18 17:51:49.123456789 Mon

now.BeginningOfMinute()        // 2013-11-18 17:51:00 Mon
now.BeginningOfHour()          // 2013-11-18 17:00:00 Mon
now.BeginningOfDay()           // 2013-11-18 00:00:00 Mon
now.BeginningOfWeek()          // 2013-11-17 00:00:00 Sun
now.BeginningOfMonth()         // 2013-11-01 00:00:00 Fri
now.BeginningOfQuarter()       // 2013-10-01 00:00:00 Tue
now.BeginningOfYear()          // 2013-01-01 00:00:00 Tue

now.WeekStartDay = time.Monday // Set Monday as first day, default is Sunday
now.BeginningOfWeek()          // 2013-11-18 00:00:00 Mon

now.EndOfMinute()              // 2013-11-18 17:51:59.999999999 Mon
now.EndOfHour()                // 2013-11-18 17:59:59.999999999 Mon
now.EndOfDay()                 // 2013-11-18 23:59:59.999999999 Mon
now.EndOfWeek()                // 2013-11-23 23:59:59.999999999 Sat
now.EndOfMonth()               // 2013-11-30 23:59:59.999999999 Sat
now.EndOfQuarter()             // 2013-12-31 23:59:59.999999999 Tue
now.EndOfYear()                // 2013-12-31 23:59:59.999999999 Tue

now.WeekStartDay = time.Monday // Set Monday as first day, default is Sunday
now.EndOfWeek()                // 2013-11-24 23:59:59.999999999 Sun
```

Calculating time based on another time

```go
t := time.Date(2013, 02, 18, 17, 51, 49, 123456789, time.Now().Location())
now.New(t).EndOfMonth()   // 2013-02-28 23:59:59.999999999 Thu
```

### Monday/Sunday

Don't be bothered with the `WeekStartDay` setting, you can use `Monday`, `Sunday`

```go
now.Monday()              // 2013-11-18 00:00:00 Mon
now.Sunday()              // 2013-11-24 00:00:00 Sun (Next Sunday)
now.EndOfSunday()         // 2013-11-24 23:59:59.999999999 Sun (End of next Sunday)

t := time.Date(2013, 11, 24, 17, 51, 49, 123456789, time.Now().Location()) // 2013-11-24 17:51:49.123456789 Sun
now.New(t).Monday()       // 2013-11-18 00:00:00 Sun (Last Monday if today is Sunday)
now.New(t).Sunday()       // 2013-11-24 00:00:00 Sun (Beginning Of Today if today is Sunday)
now.New(t).EndOfSunday()  // 2013-11-24 23:59:59.999999999 Sun (End of Today if today is Sunday)
```

### Parse String to Time

```go
time.Now() // 2013-11-18 17:51:49.123456789 Mon

// Parse(string) (time.Time, error)
t, err := now.Parse("2017")                // 2017-01-01 00:00:00, nil
t, err := now.Parse("2017-10")             // 2017-10-01 00:00:00, nil
t, err := now.Parse("2017-10-13")          // 2017-10-13 00:00:00, nil
t, err := now.Parse("1999-12-12 12")       // 1999-12-12 12:00:00, nil
t, err := now.Parse("1999-12-12 12:20")    // 1999-12-12 12:20:00, nil
t, err := now.Parse("1999-12-12 12:20:21") // 1999-12-12 12:20:00, nil
t, err := now.Parse("10-13")               // 2013-10-13 00:00:00, nil
t, err := now.Parse("12:20")               // 2013-11-18 12:20:00, nil
t, err := now.Parse("12:20:13")            // 2013-11-18 12:20:13, nil
t, err := now.Parse("14")                  // 2013-11-18 14:00:00, nil
t, err := now.Parse("99:99")               // 2013-11-18 12:20:00, Can't parse string as time: 99:99

// MustParse must parse string to time or it will panic
now.MustParse("2013-01-13")             // 2013-01-13 00:00:00
now.MustParse("02-17")                  // 2013-02-17 00:00:00
now.MustParse("2-17")                   // 2013-02-17 00:00:00
now.MustParse("8")                      // 2013-11-18 08:00:00
now.MustParse("2002-10-12 22:14")       // 2002-10-12 22:14:00
now.MustParse("99:99")                  // panic: Can't parse string as time: 99:99
```

Extend `now` to support more formats is quite easy, just update `now.TimeFormats` with other time layouts, e.g:

```go
now.TimeFormats = append(now.TimeFormats, "02 Jan 2006 15:04")
```

Please send me pull requests if you want a format to be supported officially

## Contributing

You can help to make the project better, check out [http://gorm.io/contribute.html](http://gorm.io/contribute.html) for things you can do.

# Author

**jinzhu**

* <http://github.com/jinzhu>
* <wosmvp@gmail.com>
* <http://twitter.com/zhangjinzhu>

## License

Released under the [MIT License](http://www.opensource.org/licenses/MIT).
