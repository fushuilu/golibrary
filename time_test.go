package golibrary

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeParse(t *testing.T) {
	data := []struct {
		text string
		ts   int
		ok   bool
	}{
		{text: "01:00", ts: 60 * 60, ok: true},
	}

	for _, d := range data {
		dt := DateTimeTryParse(d.text)
		assert.Equal(t, d.ts, Time2Seconds(dt)) // 将 01:00 转为 3600

		// 将 3600 转为 01:00
		st := Seconds2Time(d.ts)
		assert.Equal(t, d.text, HumanTimeHM(st))
	}
}
func TestDateTimeParse(t *testing.T) {
	data := []struct {
		text   string
		err    bool
		year   int
		month  int
		day    int
		hour   int
		minute int
		second int
	}{
		{
			text: "20131010", err: false,
			year: 2013, month: 10, day: 10,
			hour: 0, minute: 0, second: 0,
		}, {
			text: "2013-10-10", err: false,
			year: 2013, month: 10, day: 10,
			hour: 0, minute: 0, second: 0,
		}, {
			text: "20131010090909", err: false,
			year: 2013, month: 10, day: 10,
			hour: 9, minute: 9, second: 9,
		}, {
			text: "2013-10-10 09:09:09", err: false,
			year: 2013, month: 10, day: 10,
			hour: 9, minute: 9, second: 9,
		}, {
			text: "2013", err: true,
		},
	}
	for i := range data {
		parse, err := DateTimeParse(data[i].text)
		if data[i].err {
			assert.Error(t, err)
		} else {
			expect := time.Date(data[i].year, time.Month(data[i].month), data[i].day,
				data[i].hour, data[i].minute, data[i].second,
				0, time.Local)
			assert.Equal(t, expect, parse)
		}
	}

	date1 := "2020-02-04"
	fmt.Printf("use slice:%s~~~%+v\n", date1[4:5], date1[4:5] == "-")
	parse, err := DateTimeParse(date1)
	assert.Nil(t, err)
	assert.True(t, 1580745600 == parse.Unix())

	date2 := "2020-03-30T02:07:21.664Z"
	parse, err = DateTimeParse(date2)
	assert.Nil(t, err)
	fmt.Println("z2:", parse.Year(), parse.Month(), parse.Day(), parse.Hour(), parse.Minute(), parse.Second())

	date3 := "2020-06-19T00:00:00+08:00"
	parse, err = DateTimeParse(date3)
	assert.Nil(t, err)
	fmt.Println("z3:", parse.Year(), parse.Month(), parse.Day(), parse.Hour(), parse.Minute(), parse.Second())

	date4 := "1592496000000"
	parse, err = DateTimeParse(date4)
	assert.Nil(t, err)
	fmt.Println("z4:", parse.Year(), parse.Month(), parse.Day(), parse.Hour(), parse.Minute(), parse.Second())

}

func TestHumanTimeTime(t *testing.T) {
	date := time.Time{}
	r := HumanTimeYMDHMS(date)
	assert.Equal(t, r, "")

	text := "2013-10-10 09:09:09"
	date, err := DateTimeParse(text)
	assert.Nil(t, err)
	actual := HumanTimeYMDHMS(date)
	assert.Equal(t, text, actual)
}

func TestMonthDateRange(t *testing.T) {
	date5 := "2020-07"
	parse, err := DateTimeParse(date5)
	assert.Nil(t, err)
	fmt.Println("z5:", parse.Year(), parse.Month(), parse.Day())
	begin, end := Int8RangeOfMonth(parse)
	assert.True(t, begin == 20200701)
	assert.True(t, end == 20200731)
}

func TestGMT(t *testing.T) {
	jsDate := "Mon Mar 30 2020 08:03:05 GMT+0800"

	jsTime, err := time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700", jsDate)
	assert.Nil(t, err)
	fmt.Printf("time:%+v\n", jsTime)
	fmt.Println(jsTime)
	assert.False(t, jsTime.IsZero())
	assert.Equal(t, jsTime.Year(), 2020)
}

type Stu struct {
	Name  string
	Birth time.Time
}

func TestJS(t *testing.T) {
	/*
		   在 ng8 下测试了 js

		   const d = {
		     id: 5, date: new Date(),
		   };
		   this.http.http.post('http://127.0.0.1:8011/date', d).subscribe(res => {
		     console.log('date:', res);
		   });

		   在控制台下可以看到 Request Payload 提交的数据为
		   	{id: 5, date: "2020-03-30T02:07:21.664Z"}

		   type Param struct {
		   	Id   int64     `json:"id"`
		   	Date time.Time `json:"date"`
		   }
		   打印 go 下接收到的数据为
		   	{"id":5,"date":"2020-03-30T02:07:21.664Z"}

		   结论：
		   	对于前端提交的时间，后端使用如果使用 time.Time 的话，可能会限制了前端提交的时间格式;
			另外，提交空字符串，会导致解析失败
		parsing time """" as ""2006-01-02T15: 04: 05Z07: 00"": cannot parse """ as "2006"{"id":15,"date":"0001-01-01T00: 00: 00Z"}
		1. 提交参数时，使用字符串 —— 格式多样化
		2. 响应时，直接 time.Time
	*/
	d := Stu{Name: "aaa", Birth: time.Now()}
	bytes, err := json.Marshal(d)
	assert.Nil(t, err)
	fmt.Println("json.Marshal:", string(bytes))

	d2 := Stu{}
	err = json.Unmarshal(bytes, &d2) // 2020-03-30 09:57:39.241223 +0800 CST
	assert.Nil(t, err)
	fmt.Printf("d2:%+v\n", d2)
}

func TestWeek(t *testing.T) {
	tt := time.Unix(1586320504, 0)
	text := tt.Format("2006-01-02 (Mon)")
	fmt.Println(text, tt.Weekday())
}

func TestTimeToInt(t *testing.T) {
	d := 20200712
	tt := Int8ToTime(d)
	assert.True(t, tt.Year() == 2020)
	assert.True(t, tt.Month() == 7)
	assert.True(t, tt.Day() == 12)
	rst := Int8FromTime(tt)
	assert.True(t, rst == d)

	// 过去的时间
	text := "202005"
	timeParse, err := DateTimeParse(text)
	assert.Nil(t, err)

	date := Int8FromTime(timeParse)
	assert.Equal(t, 20200501, date)

	month := IntYm6FromTime(timeParse)
	assert.Equal(t, 202005, month)
}

func TestDays(t *testing.T) {
	t1 := time.Date(2020, 5, 1, 0, 0, 0, 0, time.Local) // 2020-05-01 00:00:00
	t2 := time.Date(2020, 5, 0, 0, 0, 0, 0, time.Local) // 2020-04-30 00:00:00
	fmt.Println("t1:", HumanTimeYMDHMS(t1), ";t2:", HumanTimeYMDHMS(t2))
	//assert.Equal(t, t1, t2)

	for _, v := range []struct {
		month int
		days  int
	}{
		{month: 1, days: 31},
		{month: 2, days: 29},
		{month: 3, days: 31},
		{month: 4, days: 30},

		{month: 5, days: 31},
		{month: 6, days: 30},
		{month: 7, days: 31},
		{month: 8, days: 31},

		{month: 9, days: 30},
		{month: 10, days: 31},
		{month: 11, days: 30},
		{month: 12, days: 31},
	} {
		d1 := MonthDays(time.Date(2020, time.Month(v.month), 1, 0, 0, 0, 0, time.Local))
		assert.Equal(t, v.days, d1)
	}

	t3 := time.Date(2020, 5, 3, 0, 0, 0, 0, time.Local)
	assert.Equal(t, 1, RealDays(t3, Date))
	assert.Equal(t, 31, RealDays(t3, Month))
	assert.Equal(t, 91, RealDays(t3, Quarter))
	assert.Equal(t, (31+29+31)+(30+31+30)+(31+31+30)+(31+30+31), RealDays(t3, Year))
}

func TestFrom(t *testing.T) {
	date := 201701
	dt := IntYm6ToTime(date)
	fmt.Printf("dt:%+v\n", dt)
	rst := Int8FromTime(dt)
	assert.Equal(t, 20170101, rst)
}

func TestPassDays(t *testing.T) {

	for _, v := range []struct {
		month int
		days  int
	}{
		{month: 2, days: 46},
		{month: 10, days: 15},
	} {
		d1 := PassDays(time.Date(2020, time.Month(v.month), 15, 0, 0, 0, 0, time.Local),
			Quarter)
		assert.Equal(t, v.days, d1)
	}

	date := time.Date(2020, 5, 2, 0, 0, 1, 0, time.Local)
	items := []struct {
		Type string
		Rst  int
	}{
		{Type: Year, Rst: 31 + 29 + 31 + 30 + 2},
		{Type: Month, Rst: 2},
		{Type: Quarter, Rst: 30 + 2},
		{Type: Date, Rst: 1},
	}

	for _, v := range items {
		rst := PassDays(date, v.Type)
		assert.Equal(t, v.Rst, rst, fmt.Sprintf("%s wrong", v.Type))
	}

}

func TestInt8DateSub(t *testing.T) {
	datas := []struct {
		date int
		sub  int
		kind string
		rst  int
		err  bool
	}{
		{date: 20190508, sub: 60, kind: "月", rst: 20240508, err: false},
		{date: 20190508, sub: 3, kind: "月", rst: 20190808, err: false},
		{date: 20190508, sub: 8, kind: "月", rst: 20200108, err: false},
		{date: 20190508, sub: 1, kind: "年", rst: 20200508, err: false},
		{date: 20200401, sub: -1, kind: Date, rst: 20200331, err: false},
		{date: 20200701, sub: -1, kind: Date, rst: 20200630, err: false},
		{date: 20201230, sub: 48, kind: "月", rst: 20241230, err: false},
		{date: 20201230, sub: 49, kind: "月", rst: 20250130, err: false},
		{date: 20201230, sub: 47, kind: "月", rst: 20241130, err: false},
	}

	for _, v := range datas {
		sub, err := Int8DateSub(v.date, v.sub, v.kind)
		if v.err {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
		}
		assert.Equal(t, v.rst, sub)
	}
}

func TestIntQuarter6FromInt8(t *testing.T) {
	datas := []struct {
		date int
		q    int
	}{
		{date: 20200101, q: 202001},
		{date: 20200331, q: 202001},
		{date: 20200401, q: 202002},
		{date: 20200630, q: 202002},
		{date: 20200701, q: 202003},
		{date: 20200930, q: 202003},
		{date: 20201001, q: 202004},
		{date: 20201231, q: 202004},
	}
	for _, v := range datas {
		rst := IntQuarter6FromInt8(v.date)
		assert.Equal(t, v.q, rst, fmt.Sprintf("error date:%d, expect:%d", v.date, v.q))
	}
}

func TestIntQuarter6FromInt6(t *testing.T) {

	datas := []struct {
		Month int
		Rst   int
	}{
		{Month: 202001, Rst: 202001},
		{Month: 202002, Rst: 202001},
		{Month: 202003, Rst: 202001},
		{Month: 202004, Rst: 202002},
		{Month: 202005, Rst: 202002},
		{Month: 202006, Rst: 202002},
		{Month: 202007, Rst: 202003},
		{Month: 202008, Rst: 202003},
		{Month: 202009, Rst: 202003},
		{Month: 202010, Rst: 202004},
		{Month: 202011, Rst: 202004},
		{Month: 202012, Rst: 202004},
	}

	for _, v := range datas {
		rst := IntQuarter6FromInt6(v.Month)
		assert.Equal(t, v.Rst, rst, fmt.Sprintf("%d expect %d but get %d", v.Month, v.Rst, rst))
	}
}

func TestInt8DatesBetween(t *testing.T) {
	dates := DatesBetween(20200115, 20200302)
	assert.Equal(t, 17+29+1, len(dates))
	//fmt.Println(dates)
}

func TestMonthsBetween(t *testing.T) {
	months := MonthsBetween(20200201, 20200301)
	fmt.Println(months)
}

func TestTimestamp(t *testing.T) {
	now := time.Now()
	fmt.Println("timeStamp:", now.Unix())           // 1601109774
	fmt.Println(now.Format("20060102150405"))       // 20200926164254
	fmt.Println("timeStamp:", now.UTC().Unix())     // 1601109774
	fmt.Println(now.UTC().Format("20060102150405")) // 20200926084254
	fmt.Println("时区时间戳:", TimeStamp())
}

func TestFormatTimeHM(t *testing.T) {
	args := []struct {
		text string
		hh   int
		mm   int
		err  bool
	}{
		{text: "01:50", hh: 1, mm: 50, err: false},
		{text: "12:30", hh: 12, mm: 30, err: false},
		{text: "24:01", hh: 0, mm: 0, err: true},
	}

	for _, v := range args {
		hh, mm, err := FormatTimeHM(v.text)
		assert.Equal(t, v.hh, hh)
		assert.Equal(t, v.mm, mm)
		if v.err {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestSecondsFormatTimeHM(t *testing.T) {
	type args struct {
		seconds int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "s70", args: args{seconds: 70}, want: "00:01"},
		{name: "s3670", args: args{seconds: 3670}, want: "01:01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SecondsFormatTimeHM(tt.args.seconds); got != tt.want {
				t.Errorf("SecondsFormatTimeHM() = %v, want %v", got, tt.want)
			}
		})
	}
}
