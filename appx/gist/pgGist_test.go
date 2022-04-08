package gist

import (
	"fmt"
	"github.com/fushuilu/golibrary/appx/db"
	"sort"
	"testing"

	"github.com/fushuilu/golibrary/appx/datax"
	"github.com/stretchr/testify/assert"
	"xorm.io/builder"
)
import _ "github.com/lib/pq"

var (
	eg = db.CreateEngineGroup(db.Config{
		Name: "postgres",
		Cons: []string{
			"postgres://shutao:123456@localhost/pg_demo?sslmode=disable",
		},
		MaxIdleConn: 0,
		MaxOpenConn: 10,
	})
)

type MyPoints struct {
	Id       int64   `json:"id"`
	Lng      float64 `json:"lng" xorm:"notnull default(0)"`               // 经度 longitude
	Lat      float64 `json:"lat" xorm:"notnull default(0)"`               // 纬度 latitude
	Title    string  `json:"title" xorm:"title notnull default('')"`      // 名称
	Distance int     `json:"distance" xorm:"distance notnull default(0)"` // 计算出来的距离(冗余的字段)
}

func (MyPoints) TableName() string {
	return "my_point"
}

func TestPgXorm(t *testing.T) {
	if err := eg.Ping(); err != nil {
		fmt.Println("test pgGist failed, couldn't link the db")
		return
	}
	err := eg.Sync2(MyPoints{})
	//err := Sync(eg, MyPoints{}, "202112")
	assert.Nil(t, err)

	mp := MyPoints{}
	isEmpty, err := eg.IsTableEmpty(mp)
	assert.Nil(t, err)
	if isEmpty {
		// https://api.map.baidu.com/lbsapi/getpoint/index.html
		points := []MyPoints{
			{Lng: 116.580763, Lat: 23.667547, Title: "古巷镇政府"},
			{Lng: 116.580578, Lat: 23.668101, Title: "古巷镇派出所"},
			{Lng: 116.577351, Lat: 23.670325, Title: "古巷中心卫生院"},
			{Lng: 116.580864, Lat: 23.652601, Title: "枫三小学"},
			{Lng: 116.581111, Lat: 23.661911, Title: "枫一村委会"},
			{Lng: 116.556069, Lat: 23.697566, Title: "横溪村"},
			{Lng: 116.576259, Lat: 23.683425, Title: "古五小学"},
			{Lng: 116.581089, Lat: 23.657517, Title: "枫洋综合市场"},
			{Lng: 116.578383, Lat: 23.667917, Title: "古二村"},
			{Lng: 116.565094, Lat: 23.692632, Title: "港华燃气"},
			{Lng: 116.580986, Lat: 23.668029, Title: "中国移动(古巷服务厅)"},
			{Lng: 116.580066, Lat: 23.665645, Title: "中国工商银行(古巷支行)"},
		}
		_, err := eg.Insert(&points)
		assert.Nil(t, err)
	}
	err = CreateGistExtensions(eg)
	assert.Nil(t, err)

	err = CreateGistIndex(eg, mp.TableName())
	assert.Nil(t, err)
	// 重复
	err = CreateGistIndex(eg, mp.TableName())
	assert.Nil(t, err)

	center := datax.Point{Lng: 116.580578, Lat: 23.668101} // 古巷镇派出所

	pg := NewPgGist(eg)
	rows, err := pg.Nearest(&mp, center, 1000, builder.NewCond(), db.DefaultPagination)
	assert.Nil(t, err)

	ids := pg.GetGistRecordsIds(rows)

	var points []MyPoints
	err = eg.In("id", ids).Find(&points)
	assert.Nil(t, err)

	for i := range points {
		points[i].Distance = int(pg.GetGistRecordDistance(points[i].Id, rows))
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].Distance < points[j].Distance
	})

	for i := range points {
		fmt.Printf("points:%+v\n", points[i])
	}
}

func TestNearest(t *testing.T) {

	gist := NewPgGist(eg)

	center := datax.Point{Lng: 116.580578, Lat: 23.668101} // 古巷镇派出所
	var points []MyPoints
	err := gist.NearestWith(center, 1000, builder.NewCond(), db.DefaultPagination, &points, "title", "lat", "lng")
	assert.Nil(t, err)
	fmt.Printf("points:%+v\n", points)
}
