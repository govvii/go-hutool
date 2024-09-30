package datetime

import (
	"errors"
	"fmt"
	"math"
	"time"
)

// 常用时间格式常量
const (
	ISO8601          = "2006-01-02T15:04:05Z07:00"
	ISO8601NoTZ      = "2006-01-02T15:04:05"
	RFC3339          = time.RFC3339
	RFC3339Nano      = time.RFC3339Nano
	RFC822           = time.RFC822
	RFC822Z          = time.RFC822Z
	RFC850           = time.RFC850
	RFC1123          = time.RFC1123
	RFC1123Z         = time.RFC1123Z
	ANSIC            = time.ANSIC
	UnixDate         = time.UnixDate
	RubyDate         = time.RubyDate
	DateOnly         = "2006-01-02"
	TimeOnly         = "15:04:05"
	DateTimeMinute   = "2006-01-02 15:04"
	DateTimeSecond   = "2006-01-02 15:04:05"
	DateTimeMillis   = "2006-01-02 15:04:05.000"
	DateTimeMicros   = "2006-01-02 15:04:05.000000"
	DateTimeNanos    = "2006-01-02 15:04:05.000000000"
	YearMonth        = "2006-01"
	MonthDay         = "01-02"
	YearOnly         = "2006"
	MonthOnly        = "01"
	DayOnly          = "02"
	HourMinute       = "15:04"
	HourMinuteSecond = "15:04:05"
)

// DateTime 提供了一系列日期时间相关的工具方法
type DateTime struct {
	location *time.Location
}

// New 创建一个新的 DateTime 实例
func New(location *time.Location) *DateTime {
	if location == nil {
		location = time.Local
	}
	return &DateTime{location: location}
}

// Now 返回当前时间
func (dtu *DateTime) Now() time.Time {
	return time.Now().In(dtu.location)
}

// Parse 解析字符串为时间
func (dtu *DateTime) Parse(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, dtu.location)
}

// Format 格式化时间为字符串
func (dtu *DateTime) Format(t time.Time, layout string) string {
	return t.In(dtu.location).Format(layout)
}

// AddDuration 增加或减少一段时间
func (dtu *DateTime) AddDuration(t time.Time, d time.Duration) time.Time {
	return t.Add(d)
}

// DiffDays 计算两个日期之间的天数差
func (dtu *DateTime) DiffDays(t1, t2 time.Time) int {
	duration := t2.Sub(t1)
	return int(math.Abs(duration.Hours() / 24))
}

// IsLeapYear 判断是否为闰年
func (dtu *DateTime) IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// GetWeekday 获取指定日期是星期几
func (dtu *DateTime) GetWeekday(t time.Time) time.Weekday {
	return t.Weekday()
}

// GetWeekOfYear 获取指定日期是一年中的第几周
func (dtu *DateTime) GetWeekOfYear(t time.Time) int {
	_, week := t.ISOWeek()
	return week
}

// GetDaysInMonth 获取指定年月的天数
func (dtu *DateTime) GetDaysInMonth(year, month int) int {
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, dtu.location).Day()
}

// StartOfDay 获取指定日期的开始时间（0点0分0秒）
func (dtu *DateTime) StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, dtu.location)
}

// EndOfDay 获取指定日期的结束时间（23点59分59秒）
func (dtu *DateTime) EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, dtu.location)
}

// StartOfWeek 获取指定日期所在周的开始时间（周一0点）
func (dtu *DateTime) StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return dtu.StartOfDay(t.AddDate(0, 0, -weekday+1))
}

// EndOfWeek 获取指定日期所在周的结束时间（周日23:59:59）
func (dtu *DateTime) EndOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return dtu.EndOfDay(t.AddDate(0, 0, 7-weekday))
}

// StartOfMonth 获取指定日期所在月的开始时间
func (dtu *DateTime) StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, dtu.location)
}

// EndOfMonth 获取指定日期所在月的结束时间
func (dtu *DateTime) EndOfMonth(t time.Time) time.Time {
	return dtu.StartOfMonth(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// StartOfYear 获取指定日期所在年的开始时间
func (dtu *DateTime) StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, dtu.location)
}

// EndOfYear 获取指定日期所在年的结束时间
func (dtu *DateTime) EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 12, 31, 23, 59, 59, 999999999, dtu.location)
}

// IsBetween 判断指定日期是否在两个日期之间
func (dtu *DateTime) IsBetween(t, start, end time.Time) bool {
	return (t.After(start) || t.Equal(start)) && (t.Before(end) || t.Equal(end))
}

// AddWorkdays 增加指定的工作日天数（不包括周末）
func (dtu *DateTime) AddWorkdays(t time.Time, days int) time.Time {
	for days > 0 {
		t = t.AddDate(0, 0, 1)
		if t.Weekday() != time.Saturday && t.Weekday() != time.Sunday {
			days--
		}
	}
	return t
}

// ToUTC 将时间转换为UTC时间
func (dtu *DateTime) ToUTC(t time.Time) time.Time {
	return t.UTC()
}

// ToTimeZone 将时间转换为指定时区的时间
func (dtu *DateTime) ToTimeZone(t time.Time, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}

// Age 计算年龄
func (dtu *DateTime) Age(birthDate time.Time) int {
	now := dtu.Now()
	years := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		years--
	}
	return years
}

// NextOccurrence 获取下一个指定星期几的日期
func (dtu *DateTime) NextOccurrence(t time.Time, weekday time.Weekday) time.Time {
	daysUntil := int(weekday - t.Weekday())
	if daysUntil <= 0 {
		daysUntil += 7
	}
	return t.AddDate(0, 0, daysUntil)
}

// DurationString 将持续时间格式化为人类可读的字符串
func (dtu *DateTime) DurationString(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	parts := []string{}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d days", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d hours", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d minutes", minutes))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%d seconds", seconds))
	}

	return dtu.joinStrings(parts, ", ")
}

// ParseDuration 解析人类可读的持续时间字符串
func (dtu *DateTime) ParseDuration(s string) (time.Duration, error) {
	d, err := time.ParseDuration(s)
	if err == nil {
		return d, nil
	}

	// 自定义解析逻辑，支持天、周、月、年
	var total time.Duration
	var value int
	var unit string

	_, err = fmt.Sscanf(s, "%d%s", &value, &unit)
	if err != nil {
		return 0, errors.New("invalid duration format")
	}

	switch unit {
	case "d", "day", "days":
		total = time.Duration(value) * 24 * time.Hour
	case "w", "week", "weeks":
		total = time.Duration(value) * 7 * 24 * time.Hour
	case "m", "month", "months":
		total = time.Duration(value) * 30 * 24 * time.Hour // 近似值
	case "y", "year", "years":
		total = time.Duration(value) * 365 * 24 * time.Hour // 近似值
	default:
		return 0, errors.New("unknown time unit")
	}

	return total, nil
}

// joinStrings 连接字符串切片
func (dtu *DateTime) joinStrings(strings []string, separator string) string {
	if len(strings) == 0 {
		return ""
	}
	if len(strings) == 1 {
		return strings[0]
	}
	return strings[0] + separator + dtu.joinStrings(strings[1:], separator)
}
