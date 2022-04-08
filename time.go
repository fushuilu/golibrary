package golibrary

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	HalfHourSeconds = 60 * 30
	MinuteSeconds   = 60
	HourSeconds     = 60 * 60
	DaySeconds      = HourSeconds * 24
)

var zero, _ = time.ParseInLocation("15:04", "00:00", time.Local)

func Time2Seconds(dt time.Time) int {
	return int(dt.Sub(zero).Seconds())
}

func Seconds2Time(seconds int) time.Time {
	return zero.Add(time.Second * time.Duration(seconds))
}

var CnWeekDays = map[time.Weekday]string{
	time.Sunday:    "周日",
	time.Monday:    "周一",
	time.Tuesday:   "周二",
	time.Wednesday: "周三",
	time.Thursday:  "周四",
	time.Friday:    "周五",
	time.Saturday:  "周六",
}

var numsReg = regexp.MustCompile("^\\d+$")
var hmReg = regexp.MustCompile("^\\d\\d:\\d\\d$")

func FormatTimeHM(text string) (hh, mm int, err error) {
	if hmReg.MatchString(text) {
		if t, err := time.ParseInLocation("15:04", text, time.Local); err != nil {
			return 0, 0, err
		} else {
			return t.Hour(), t.Minute(), nil
		}
	}
	return 0, 0, errors.New("不是有效的 hh:mm 格式")
}

func SecondsFormatTimeHM(seconds int) string {
	if seconds > 0 {
		hour := seconds / 3600
		min := seconds % 3600 / 60
		return fmt.Sprintf("%02d:%02d", hour, min)
	}
	return ""
}
func DateTimeParse(text string) (time.Time, error) {
	if text == "" {
		return time.Time{}, nil
	}
	switch len(text) {
	case 5: // hh:mm
		return time.ParseInLocation("15:04", text, time.Local) // 注意可能无法存入数据库中
	case 6:
		return time.ParseInLocation("200601", text, time.Local)
	case 7:
		return time.ParseInLocation("2006-01", text, time.Local)
	case 8:
		if strings.Index(text, ":") > 0 { // hh:mm:ss
			return time.ParseInLocation("15:04:05", text, time.Local)
		}
		return time.ParseInLocation("20060102", text, time.Local)
	case 10:
		// date1[4:5] == "-" ?
		if numsReg.Match([]byte(text)) { // 时间戳
			ts, _ := strconv.ParseInt(text, 10, 0)
			return time.Unix(ts, 0), nil
		}
		return time.ParseInLocation("2006-01-02", text, time.Local)
	case 13: // js 时间戳 1592496000000
		if numsReg.Match([]byte(text)) {
			ts, _ := strconv.ParseInt(text[:10], 10, 0)
			ns, _ := strconv.ParseInt(text[11:], 10, 0)
			return time.Unix(ts, ns), nil
		}
	case 14:
		return time.ParseInLocation("20060102150405", text, time.Local)
	case 16:
		return time.ParseInLocation("2006-01-02 15:04", text, time.Local)
	case 19:
		return time.ParseInLocation("2006-01-02 15:04:05", text, time.Local)
	case 20: // 0001-01-01T00:00:00Z
		if text == "0001-01-01T00:00:00Z" {
			return time.Time{}, nil
		}
	case 24: // 2020-03-30T02:07:21.664Z
		return time.ParseInLocation("2006-01-02T15:04:05.000Z", text, time.Local)
	case 25: // 2020-06-19T00:00:00+08:00
		return time.ParseInLocation("2006-01-02T15:04:05-07:00", text, time.Local)
	case 29:
		// Thu, 02 Apr 2020 21:30:00 GMT
		// js toUTCString
		if strings.HasSuffix(text, "GMT") {
			return time.ParseInLocation("Mon, 02 Jan 2006 15:04:05 GMT", text, time.Local)
		} else if strings.HasSuffix(text, "MST") {
			return time.ParseInLocation("Mon, 02 Jan 2006 15:04:05 MST", text, time.Local)
		}
	case 33:
		return time.ParseInLocation("Mon Jan 02 2006 15:04:05 GMT-0700", text, time.Local)
	}
	return time.Time{}, errors.New("不支持的时间格式")
}

func DateTimeTryParse(text string) time.Time {
	parse, _ := DateTimeParse(text)
	return parse
}

func DateTimeIfNotEmpty(text string, set func(t time.Time)) {
	if text != "" {
		if parse, _ := DateTimeParse(text); !parse.IsZero() {
			set(parse)
		}
	}
}

type Format int

const (
	FormatGmt           Format = iota
	FormatDateTime             // yyyy-MM-dd hh:mm:ss
	FormatDate                 // yyyy-MM-dd
	FormatTime                 // hh:mm:ss
	FormatDateTimeLong         // yyyy-MM-dd hh:mm:ss 或者 yyyy-MM-dd
	FormatDateTimeShort        // yyyy-MM-dd hh:mm 或者 yyyy-MM-dd
	FormatDateTimeYMDHM        // yyyy-MM-dd hh:mm
	FormatDateTimeYMD          // yyyyMMdd
	FormatDateTimeYM           // yyyyMM
	FormatDateTimeHMS          // hh:mm:ss
	FormatDateTimeHM           // hh:mm
)

func DateTimeFormat(t time.Time, f Format) string {
	if t.IsZero() {
		return ""
	}
	switch f {
	case FormatGmt:
		return t.Format("Mon Jan 02 2006 15:04:05 GMT-0700")
	case FormatDateTime:
		return t.Format("2006-01-02 15:04:05")
	case FormatDateTimeYMDHM:
		return t.Format("2006-01-02 15:04")
	case FormatDate:
		return t.Format("2006-01-02")
	case FormatTime:
		return t.Format("15:04:05")
	case FormatDateTimeLong:
		text := t.Format("2006-01-02 15:04:05")
		return strings.Replace(text, " 00:00:00", "", 1)
	case FormatDateTimeShort:
		text := t.Format("2006-01-02 15:04")
		return strings.Replace(text, " 00:00", "", 1)
	case FormatDateTimeYMD:
		return t.Format("20060102")
	case FormatDateTimeYM:
		return t.Format("200601")
	case FormatDateTimeHMS:
		return t.Format("15:04:05")
	case FormatDateTimeHM:
		return t.Format("15:04")
	}
	return ""
}

// JS 能够直接处理 0001-01-01T00:00:00Z 格式的时间字符串
func HumanTimeGMT(t time.Time) string {
	return DateTimeFormat(t, FormatGmt)
}

// 将 time.Time 转为 YYYY-MM-DD hh:mm:ss
func HumanTimeYMDHMS(t time.Time) string {
	return DateTimeFormat(t, FormatDateTime)
}

// 将 time.Time 转为 YYYY-MM-DD hh:mm
func HumanTimeYMDHM(t time.Time) string {
	return DateTimeFormat(t, FormatDateTimeYMDHM)
}

// 将 time.Time 转为 YYYY-MM-DD；如果需要 YYYYMMDD 格式，请使用 DateTimeFormat(t, FormatDateTimeYMD)
func HumanTimeYMD(t time.Time) string {
	return DateTimeFormat(t, FormatDate)
}

func HumanTimeHMS(t time.Time) string {
	return DateTimeFormat(t, FormatDateTimeHMS)
}

func HumanTimeHM(t time.Time) string {
	return DateTimeFormat(t, FormatDateTimeHM)
}

// 转为 YYYY-MM-DD 或者 YYYY-MM-DD hh:mm:ss
func HumanDateTimeLong(t time.Time) string {
	return DateTimeFormat(t, FormatDateTimeLong)
}

// 转为 YYYY-MM-DD 或者 YYYY-MM-DD hh:mm
func HumanDateTimeShort(t time.Time) string {
	return DateTimeFormat(t, FormatDateTimeShort)
}

// 按小到大排列
func DateTimeOrderAsc(t1 string, t2 string) (t11 time.Time, t22 time.Time, err error) {
	if t11, err = DateTimeParse(t1); err != nil {
		return
	}
	if t22, err = DateTimeParse(t2); err != nil {
		return
	}
	if t11.IsZero() || t22.IsZero() {
		return
	}
	if t11.After(t22) {
		return t22, t11, nil
	}
	return
}

func Int8FromTime(t time.Time) int {
	return t.Year()*10000 + int(t.Month())*100 + t.Day()
}
func Int8ToTime(d int) time.Time {
	return time.Date(d/10000, time.Month(d%10000/100), d%100, 0, 0, 0, 0, time.Local)
}

func IntYm6FromTime(t time.Time) int {
	return t.Year()*100 + int(t.Month())
}

func IntYm6ToTime(m int) time.Time {
	return time.Date(m/100, time.Month(m%100), 1, 0, 0, 0, 0, time.Local)
}

// 将 yyyymm 转为 yyyy-mm
func IntYmd6Format(d int) string {
	if d < 1 {
		return ""
	}
	return strings.Join([]string{
		strconv.Itoa(d / 100),
		PadStartZero(d%100, 2),
	}, "-")
}

// 季度 quarter 转换 yyyy01, yyyy02, yyyy03, yyyy04
func IntQuarter6FromTime(t time.Time) int {
	return t.Year()*100 + (int(t.Month())+2)/3
}
func IntQuarter6FromInt6(int6 int) int {
	return int6/100*100 + (int6%100+2)/3
}

func IntQuarter6FromInt8(date int) int {
	q := (date/100%100 + 2) / 3
	return date/10000*100 + q
}

// 指定时间所在月份的第1天和最后1天
func Int8RangeOfMonth(t time.Time) (mBegin int, mEnd int) {
	first := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	last := first.AddDate(0, 1, -1)
	return Int8FromTime(first), Int8FromTime(last)
}

func Int8RangeOfQuarter(t time.Time) (dBegin int, dEnd int) {
	year := t.Year() * 10000
	switch t.Month() {
	case time.January, time.February, time.March:
		dBegin = 101
		dEnd = 331
	case time.April, time.May, time.June:
		dBegin = 401
		dEnd = 630
	case time.July, time.August, time.September:
		dBegin = 701
		dEnd = 930
	case time.October, time.November, time.December:
		dBegin = 1001
		dEnd = 1231
	}
	return year + dBegin, year + dEnd
}

func Int8DateSub(date int, sub int, kind string) (int, error) {
	if date == 0 || sub == 0 {
		return date, nil
	}
	switch kind {
	case Year, "年":
		return date + sub*10000, nil
	case Month, "月":
		oldMonth := date / 100 % 100
		oldDate := date % 100
		//fmt.Printf("date:%d, y:%d, m:%d, d:%d\n", date, oldYear, oldMonth, oldDate)

		newYear := date/10000 + sub/12
		newMonth := sub%12 + oldMonth
		if newMonth > 12 {
			newMonth -= 12
			newYear += 1
		} else if newMonth < 0 {
			newMonth += 12
			newYear -= 1
		}
		//fmt.Printf("after: y:%d, m:%d\n", newYear, newMonth)
		return oldDate + newMonth*100 + newYear*10000, nil
	case Date:
		t := Int8ToTime(date)
		t = t.AddDate(0, 0, sub)
		return Int8FromTime(t), nil
	default:
		return 0, errors.New("暂不支持的 Int8Date 求和")
	}
}

func Int6RangeOfQuarter(t time.Time) (mBegin int, mEnd int) {
	year := t.Year() * 100
	switch t.Month() {
	case time.January, time.February, time.March:
		mBegin = 1
		mEnd = 3
	case time.April, time.May, time.June:
		mBegin = 4
		mEnd = 6
	case time.July, time.August, time.September:
		mBegin = 7
		mEnd = 9
	case time.October, time.November, time.December:
		mBegin = 10
		mEnd = 12
	}
	return year + mBegin, year + mEnd
}

func DateTimeRangeOfQuarter(t time.Time) (begin time.Time, end time.Time) {
	switch t.Month() {
	case time.January, time.February, time.March:
		begin = time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.Local)
		end = time.Date(t.Year(), 3, 31, 0, 0, 0, 0, time.Local)
	case time.April, time.May, time.June:
		begin = time.Date(t.Year(), 4, 1, 0, 0, 0, 0, time.Local)
		end = time.Date(t.Year(), 6, 31, 0, 0, 0, 0, time.Local)
	case time.July, time.August, time.September:
		begin = time.Date(t.Year(), 7, 1, 0, 0, 0, 0, time.Local)
		end = time.Date(t.Year(), 9, 30, 0, 0, 0, 0, time.Local)
	case time.October, time.November, time.December:
		begin = time.Date(t.Year(), 10, 1, 0, 0, 0, 0, time.Local)
		end = time.Date(t.Year(), 12, 31, 0, 0, 0, 0, time.Local)
	}
	return
}

func DateTimeRangeOfDate(date time.Time) (begin time.Time, end time.Time) {
	if date.IsZero() {
		return
	}
	begin = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	end = begin.AddDate(0, 0, 1)
	return
}

// 返回从 [start, end) 之间的日期，注意 end 不包含
func DatesBetween(start, end int) (dates []int) {
	if start < end {
		startTime := Int8ToTime(start)
		endTime := Int8ToTime(end)
		for startTime.Before(endTime) {
			dates = append(dates, Int8FromTime(startTime))
			startTime = startTime.AddDate(0, 0, 1)
		}
	}
	return
}

func MonthsBetween(start, end int) (months []int) {
	if start < end {
		months = append(months, start/100)
		months = append(months, end/100)
		// 计算中间的
		startTime := Int8ToTime(start)
		endTime := Int8ToTime(end)
		firstStartTime := FirstDate(startTime, Month)

		if firstStartTime.Before(endTime) {
			months = append(months, IntYm6FromTime(startTime))
			firstStartTime = firstStartTime.AddDate(0, 1, 0)
		}

		return IntUnique(months)
	}
	return
}

// 时间天数
func Days(t1, t2 time.Time) int {
	if t1.Before(t2) {
		return int(t2.Sub(t1).Hours() / 24)
	} else {
		return int(t1.Sub(t2).Hours() / 24)
	}
}

// 指定月份的天数
func MonthDays(t1 time.Time) int {
	t11 := time.Date(t1.Year(), t1.Month(), 1, 0, 0, 0, 0, time.Local)
	t22 := t11.AddDate(0, 1, 0)
	return Days(t11, t22)
}

const (
	Year    = "year"
	Month   = "month"
	Quarter = "quarter"
	Date    = "date"
)

// 当前年份天数
func YearRealDays(t1 time.Time) int {
	return RealDays(t1, Year)
}

// 当前季度天数
func QuarterRealDays(t1 time.Time) int {
	return RealDays(t1, Quarter)
}

// 当前月份天数
func MonthRealDays(t1 time.Time) int {
	return RealDays(t1, Month)
}

// 至少要返回 1 天
func RealDays(t time.Time, kind string) int {
	switch kind {
	case Year, "y":
		return Days(
			time.Date(t.Year(), 1, 0, 0, 0, 0, 0, time.Local),
			time.Date(t.Year()+1, 1, 0, 0, 0, 0, 0, time.Local),
		)
	case Quarter, "q":
		var begin time.Time
		var end time.Time
		switch t.Month() {
		case time.January, time.February, time.March:
			begin = time.Date(t.Year(), 1, 0, 0, 0, 0, 0, time.Local)
			end = time.Date(t.Year(), 3, 31, 0, 0, 0, 0, time.Local)
		case time.April, time.May, time.June:
			begin = time.Date(t.Year(), 4, 0, 0, 0, 0, 0, time.Local)
			end = time.Date(t.Year(), 6, 30, 0, 0, 0, 0, time.Local)

		case time.July, time.August, time.September:
			begin = time.Date(t.Year(), 7, 0, 0, 0, 0, 0, time.Local)
			end = time.Date(t.Year(), 9, 30, 0, 0, 0, 0, time.Local)
		case time.October, time.November, time.December:
			begin = time.Date(t.Year(), 10, 0, 0, 0, 0, 0, time.Local)
			end = time.Date(t.Year(), 12, 31, 0, 0, 0, 0, time.Local)
		}
		return Days(begin, end)
	case Month, "m":
		t1 := time.Date(t.Year(), t.Month(), 0, 0, 0, 0, 0, time.Local)
		t2 := time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.Local)
		return Days(t1, t2)
	default:
		return 1
	}
}

// 过去多少天
func PassDays(t time.Time, kind string) int {
	switch kind {
	case Year, "y":
		return Days(time.Date(t.Year(), 1, 0, 0, 0, 0, 0, time.Local), t)
	case Quarter, "q":
		var begin time.Time
		switch t.Month() {
		case time.January, time.February, time.March:
			begin = time.Date(t.Year(), 1, 0, 0, 0, 0, 0, time.Local)
		case time.April, time.May, time.June:
			begin = time.Date(t.Year(), 4, 0, 0, 0, 0, 0, time.Local)
		case time.July, time.August, time.September:
			begin = time.Date(t.Year(), 7, 0, 0, 0, 0, 0, time.Local)
		case time.October, time.November, time.December:
			begin = time.Date(t.Year(), 10, 0, 0, 0, 0, 0, time.Local)
		}
		return Days(begin, t)
	case Month, "m":
		t11 := time.Date(t.Year(), t.Month(), 0, 0, 0, 0, 0, time.Local)
		return Days(t11, t)
	default:
		return 1
	}
}

func FirstDate(t time.Time, kind string) time.Time {
	switch kind {
	case Year, "y":
		return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	case Quarter, "q":
		switch t.Month() {
		case time.January, time.February, time.March:
			return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.Local)
		case time.April, time.May, time.June:
			return time.Date(t.Year(), 4, 1, 0, 0, 0, 0, time.Local)
		case time.July, time.August, time.September:
			return time.Date(t.Year(), 7, 1, 0, 0, 0, 0, time.Local)
		case time.October, time.November, time.December:
			return time.Date(t.Year(), 10, 1, 0, 0, 0, 0, time.Local)
		}
	case Month, "m":
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	default:
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	}
	return t
}

func QuarterLastDate(t time.Time) time.Time {
	switch t.Month() {
	case time.January, time.February, time.March:
		return time.Date(t.Year(), 3, 31, 0, 0, 0, 0, time.Local)
	case time.April, time.May, time.June:
		return time.Date(t.Year(), 6, 30, 0, 0, 0, 0, time.Local)
	case time.July, time.August, time.September:
		return time.Date(t.Year(), 9, 30, 0, 0, 0, 0, time.Local)
	case time.October, time.November, time.December:
		return time.Date(t.Year(), 12, 31, 0, 0, 0, 0, time.Local)
	}
	return t
}

// 时区时间戳
func TimeStamp() int64 {
	t := time.Now()
	_, offset := t.Zone()
	return t.Unix() + int64(offset)
}

// 当前 time.Time 当天起始时间
func BeginDayTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

func Today() time.Time {
	return BeginDayTime(time.Now())
}

func DatetimeInvalid(begin, end string, mustBegin, mustEnd bool) (beginAt, endAt time.Time, err error) {
	if beginAt = DateTimeTryParse(begin); mustBegin && beginAt.IsZero() {
		err = errors.New("必须指定开始时间")
		return
	}
	if endAt = DateTimeTryParse(end); mustEnd && endAt.IsZero() {
		err = errors.New("必须指定结束时间")
		return
	}
	if !beginAt.IsZero() && !endAt.IsZero() && beginAt.After(endAt) {
		err = errors.New("开始时间不能晚于结束时间")
	}
	return
}

type GmtTime time.Time

func (t GmtTime) MarshalJSON() ([]byte, error) {
	if time.Time(t).IsZero() {
		return []byte(`""`), nil
	}
	return []byte(DateTimeFormat(time.Time(t), FormatGmt)), nil
}
