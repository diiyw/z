package stdlib

import (
	"time"

	"github.com/diiyw/z"
)

var timesModule = map[string]z.Object{
	"format_ansic":        &z.String{Value: time.ANSIC},
	"format_unix_date":    &z.String{Value: time.UnixDate},
	"format_ruby_date":    &z.String{Value: time.RubyDate},
	"format_rfc822":       &z.String{Value: time.RFC822},
	"format_rfc822z":      &z.String{Value: time.RFC822Z},
	"format_rfc850":       &z.String{Value: time.RFC850},
	"format_rfc1123":      &z.String{Value: time.RFC1123},
	"format_rfc1123z":     &z.String{Value: time.RFC1123Z},
	"format_rfc3339":      &z.String{Value: time.RFC3339},
	"format_rfc3339_nano": &z.String{Value: time.RFC3339Nano},
	"format_kitchen":      &z.String{Value: time.Kitchen},
	"format_stamp":        &z.String{Value: time.Stamp},
	"format_stamp_milli":  &z.String{Value: time.StampMilli},
	"format_stamp_micro":  &z.String{Value: time.StampMicro},
	"format_stamp_nano":   &z.String{Value: time.StampNano},
	"nanosecond":          &z.Int{Value: int64(time.Nanosecond)},
	"microsecond":         &z.Int{Value: int64(time.Microsecond)},
	"millisecond":         &z.Int{Value: int64(time.Millisecond)},
	"second":              &z.Int{Value: int64(time.Second)},
	"minute":              &z.Int{Value: int64(time.Minute)},
	"hour":                &z.Int{Value: int64(time.Hour)},
	"january":             &z.Int{Value: int64(time.January)},
	"february":            &z.Int{Value: int64(time.February)},
	"march":               &z.Int{Value: int64(time.March)},
	"april":               &z.Int{Value: int64(time.April)},
	"may":                 &z.Int{Value: int64(time.May)},
	"june":                &z.Int{Value: int64(time.June)},
	"july":                &z.Int{Value: int64(time.July)},
	"august":              &z.Int{Value: int64(time.August)},
	"september":           &z.Int{Value: int64(time.September)},
	"october":             &z.Int{Value: int64(time.October)},
	"november":            &z.Int{Value: int64(time.November)},
	"december":            &z.Int{Value: int64(time.December)},
	"sleep": &z.UserFunction{
		Name:  "sleep",
		Value: timesSleep,
	}, // sleep(int)
	"parse_duration": &z.UserFunction{
		Name:  "parse_duration",
		Value: timesParseDuration,
	}, // parse_duration(str) => int
	"since": &z.UserFunction{
		Name:  "since",
		Value: timesSince,
	}, // since(time) => int
	"until": &z.UserFunction{
		Name:  "until",
		Value: timesUntil,
	}, // until(time) => int
	"duration_hours": &z.UserFunction{
		Name:  "duration_hours",
		Value: timesDurationHours,
	}, // duration_hours(int) => float
	"duration_minutes": &z.UserFunction{
		Name:  "duration_minutes",
		Value: timesDurationMinutes,
	}, // duration_minutes(int) => float
	"duration_nanoseconds": &z.UserFunction{
		Name:  "duration_nanoseconds",
		Value: timesDurationNanoseconds,
	}, // duration_nanoseconds(int) => int
	"duration_seconds": &z.UserFunction{
		Name:  "duration_seconds",
		Value: timesDurationSeconds,
	}, // duration_seconds(int) => float
	"duration_string": &z.UserFunction{
		Name:  "duration_string",
		Value: timesDurationString,
	}, // duration_string(int) => string
	"month_string": &z.UserFunction{
		Name:  "month_string",
		Value: timesMonthString,
	}, // month_string(int) => string
	"date": &z.UserFunction{
		Name:  "date",
		Value: timesDate,
	}, // date(year, month, day, hour, min, sec, nsec) => time
	"now": &z.UserFunction{
		Name:  "now",
		Value: timesNow,
	}, // now() => time
	"parse": &z.UserFunction{
		Name:  "parse",
		Value: timesParse,
	}, // parse(format, str) => time
	"unix": &z.UserFunction{
		Name:  "unix",
		Value: timesUnix,
	}, // unix(sec, nsec) => time
	"add": &z.UserFunction{
		Name:  "add",
		Value: timesAdd,
	}, // add(time, int) => time
	"add_date": &z.UserFunction{
		Name:  "add_date",
		Value: timesAddDate,
	}, // add_date(time, years, months, days) => time
	"sub": &z.UserFunction{
		Name:  "sub",
		Value: timesSub,
	}, // sub(t time, u time) => int
	"after": &z.UserFunction{
		Name:  "after",
		Value: timesAfter,
	}, // after(t time, u time) => bool
	"before": &z.UserFunction{
		Name:  "before",
		Value: timesBefore,
	}, // before(t time, u time) => bool
	"time_year": &z.UserFunction{
		Name:  "time_year",
		Value: timesTimeYear,
	}, // time_year(time) => int
	"time_month": &z.UserFunction{
		Name:  "time_month",
		Value: timesTimeMonth,
	}, // time_month(time) => int
	"time_day": &z.UserFunction{
		Name:  "time_day",
		Value: timesTimeDay,
	}, // time_day(time) => int
	"time_weekday": &z.UserFunction{
		Name:  "time_weekday",
		Value: timesTimeWeekday,
	}, // time_weekday(time) => int
	"time_hour": &z.UserFunction{
		Name:  "time_hour",
		Value: timesTimeHour,
	}, // time_hour(time) => int
	"time_minute": &z.UserFunction{
		Name:  "time_minute",
		Value: timesTimeMinute,
	}, // time_minute(time) => int
	"time_second": &z.UserFunction{
		Name:  "time_second",
		Value: timesTimeSecond,
	}, // time_second(time) => int
	"time_nanosecond": &z.UserFunction{
		Name:  "time_nanosecond",
		Value: timesTimeNanosecond,
	}, // time_nanosecond(time) => int
	"time_unix": &z.UserFunction{
		Name:  "time_unix",
		Value: timesTimeUnix,
	}, // time_unix(time) => int
	"time_unix_nano": &z.UserFunction{
		Name:  "time_unix_nano",
		Value: timesTimeUnixNano,
	}, // time_unix_nano(time) => int
	"time_format": &z.UserFunction{
		Name:  "time_format",
		Value: timesTimeFormat,
	}, // time_format(time, format) => string
	"time_location": &z.UserFunction{
		Name:  "time_location",
		Value: timesTimeLocation,
	}, // time_location(time) => string
	"time_string": &z.UserFunction{
		Name:  "time_string",
		Value: timesTimeString,
	}, // time_string(time) => string
	"is_zero": &z.UserFunction{
		Name:  "is_zero",
		Value: timesIsZero,
	}, // is_zero(time) => bool
	"to_local": &z.UserFunction{
		Name:  "to_local",
		Value: timesToLocal,
	}, // to_local(time) => time
	"to_utc": &z.UserFunction{
		Name:  "to_utc",
		Value: timesToUTC,
	}, // to_utc(time) => time
	"in_location": &z.UserFunction{
		Name:  "in_location",
		Value: timesInLocation,
	}, // in_location(time, location) => time
}

func timesSleep(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	i1, ok := z.ToInt64(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	time.Sleep(time.Duration(i1))
	ret = z.UndefinedValue

	return
}

func timesParseDuration(args ...z.Object) (
	ret z.Object,
	err error,
) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	s1, ok := z.ToString(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	dur, err := time.ParseDuration(s1)
	if err != nil {
		ret = wrapError(err)
		return
	}

	ret = &z.Int{Value: int64(dur)}

	return
}

func timesSince(args ...z.Object) (
	ret z.Object,
	err error,
) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Int{Value: int64(time.Since(t1))}

	return
}

func timesUntil(args ...z.Object) (
	ret z.Object,
	err error,
) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Int{Value: int64(time.Until(t1))}

	return
}

func timesDurationHours(args ...z.Object) (
	ret z.Object,
	err error,
) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	i1, ok := z.ToInt64(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Float{Value: time.Duration(i1).Hours()}

	return
}

func timesDurationMinutes(args ...z.Object) (
	ret z.Object,
	err error,
) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	i1, ok := z.ToInt64(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Float{Value: time.Duration(i1).Minutes()}

	return
}

func timesDurationNanoseconds(args ...z.Object) (
	ret z.Object,
	err error,
) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	i1, ok := z.ToInt64(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Int{Value: time.Duration(i1).Nanoseconds()}

	return
}

func timesDurationSeconds(args ...z.Object) (
	ret z.Object,
	err error,
) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	i1, ok := z.ToInt64(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Float{Value: time.Duration(i1).Seconds()}

	return
}

func timesDurationString(args ...z.Object) (
	ret z.Object,
	err error,
) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	i1, ok := z.ToInt64(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.String{Value: time.Duration(i1).String()}

	return
}

func timesMonthString(args ...z.Object) (
	ret z.Object,
	err error,
) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	i1, ok := z.ToInt64(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.String{Value: time.Month(i1).String()}

	return
}

func timesDate(args ...z.Object) (
	ret z.Object,
	err error,
) {
	if len(args) < 7 || len(args) > 8 {
		err = z.ErrWrongNumArguments
		return
	}

	i1, ok := z.ToInt(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}
	i2, ok := z.ToInt(args[1])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}
	i3, ok := z.ToInt(args[2])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "third",
			Expected: "int(compatible)",
			Found:    args[2].TypeName(),
		}
		return
	}
	i4, ok := z.ToInt(args[3])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "fourth",
			Expected: "int(compatible)",
			Found:    args[3].TypeName(),
		}
		return
	}
	i5, ok := z.ToInt(args[4])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "fifth",
			Expected: "int(compatible)",
			Found:    args[4].TypeName(),
		}
		return
	}
	i6, ok := z.ToInt(args[5])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "sixth",
			Expected: "int(compatible)",
			Found:    args[5].TypeName(),
		}
		return
	}
	i7, ok := z.ToInt(args[6])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "seventh",
			Expected: "int(compatible)",
			Found:    args[6].TypeName(),
		}
		return
	}

	var loc *time.Location
	if len(args) == 8 {
		i8, ok := z.ToString(args[7])
		if !ok {
			err = z.ErrInvalidArgumentType{
				Name:     "eighth",
				Expected: "string(compatible)",
				Found:    args[7].TypeName(),
			}
			return
		}
		loc, err = time.LoadLocation(i8)
		if err != nil {
			ret = wrapError(err)
			return
		}
	} else {
		loc = time.Now().Location()
	}

	ret = &z.Time{
		Value: time.Date(i1,
			time.Month(i2), i3, i4, i5, i6, i7, loc),
	}

	return
}

func timesNow(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 0 {
		err = z.ErrWrongNumArguments
		return
	}

	ret = &z.Time{Value: time.Now()}

	return
}

func timesParse(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 2 {
		err = z.ErrWrongNumArguments
		return
	}

	s1, ok := z.ToString(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	s2, ok := z.ToString(args[1])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}

	parsed, err := time.Parse(s1, s2)
	if err != nil {
		ret = wrapError(err)
		return
	}

	ret = &z.Time{Value: parsed}

	return
}

func timesUnix(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 2 {
		err = z.ErrWrongNumArguments
		return
	}

	i1, ok := z.ToInt64(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	i2, ok := z.ToInt64(args[1])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}

	ret = &z.Time{Value: time.Unix(i1, i2)}

	return
}

func timesAdd(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 2 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	i2, ok := z.ToInt64(args[1])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}

	ret = &z.Time{Value: t1.Add(time.Duration(i2))}

	return
}

func timesSub(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 2 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	t2, ok := z.ToTime(args[1])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "time(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}

	ret = &z.Int{Value: int64(t1.Sub(t2))}

	return
}

func timesAddDate(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 4 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	i2, ok := z.ToInt(args[1])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}

	i3, ok := z.ToInt(args[2])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "third",
			Expected: "int(compatible)",
			Found:    args[2].TypeName(),
		}
		return
	}

	i4, ok := z.ToInt(args[3])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "fourth",
			Expected: "int(compatible)",
			Found:    args[3].TypeName(),
		}
		return
	}

	ret = &z.Time{Value: t1.AddDate(i2, i3, i4)}

	return
}

func timesAfter(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 2 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	t2, ok := z.ToTime(args[1])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "time(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}

	if t1.After(t2) {
		ret = z.TrueValue
	} else {
		ret = z.FalseValue
	}

	return
}

func timesBefore(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 2 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	t2, ok := z.ToTime(args[1])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	if t1.Before(t2) {
		ret = z.TrueValue
	} else {
		ret = z.FalseValue
	}

	return
}

func timesTimeYear(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Int{Value: int64(t1.Year())}

	return
}

func timesTimeMonth(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Int{Value: int64(t1.Month())}

	return
}

func timesTimeDay(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Int{Value: int64(t1.Day())}

	return
}

func timesTimeWeekday(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Int{Value: int64(t1.Weekday())}

	return
}

func timesTimeHour(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Int{Value: int64(t1.Hour())}

	return
}

func timesTimeMinute(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Int{Value: int64(t1.Minute())}

	return
}

func timesTimeSecond(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Int{Value: int64(t1.Second())}

	return
}

func timesTimeNanosecond(args ...z.Object) (
	ret z.Object,
	err error,
) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Int{Value: int64(t1.Nanosecond())}

	return
}

func timesTimeUnix(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Int{Value: t1.Unix()}

	return
}

func timesTimeUnixNano(args ...z.Object) (
	ret z.Object,
	err error,
) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Int{Value: t1.UnixNano()}

	return
}

func timesTimeFormat(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 2 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	s2, ok := z.ToString(args[1])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}

	s := t1.Format(s2)
	if len(s) > z.MaxStringLen {

		return nil, z.ErrStringLimit
	}

	ret = &z.String{Value: s}

	return
}

func timesIsZero(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	if t1.IsZero() {
		ret = z.TrueValue
	} else {
		ret = z.FalseValue
	}

	return
}

func timesToLocal(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Time{Value: t1.Local()}

	return
}

func timesToUTC(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.Time{Value: t1.UTC()}

	return
}

func timesTimeLocation(args ...z.Object) (
	ret z.Object,
	err error,
) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.String{Value: t1.Location().String()}

	return
}

func timesInLocation(args ...z.Object) (
	ret z.Object,
	err error,
) {
	if len(args) != 2 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	s2, ok := z.ToString(args[1])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}

	location, err := time.LoadLocation(s2)
	if err != nil {
		ret = wrapError(err)
		return
	}

	ret = &z.Time{Value: t1.In(location)}

	return
}

func timesTimeString(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		err = z.ErrWrongNumArguments
		return
	}

	t1, ok := z.ToTime(args[0])
	if !ok {
		err = z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "time(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	ret = &z.String{Value: t1.String()}

	return
}
