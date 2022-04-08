package dbx

import (
	"github.com/fushuilu/golibrary"
	"reflect"
)

// 注意使用在模型上时，不要使用 *
type ModelCache interface {
	GetId() int64
	GetTitle() string
}

//
//  GetIdsFromOther
//  @Description: 从其它结构体中组成 id 集合
//  @param rows 普通切片数据
//  @param get 返回 id
//  @return []int64
//
func GetIdsFromOther(rows interface{}, get func(row interface{}) int64) []int64 {
	ids := make([]int64, 0)
	if golibrary.IsSlice(rows) {
		s := reflect.ValueOf(rows)
		for i := 0; i < s.Len(); i++ {
			ids = append(ids, get(s.Index(i).Interface()))
		}
	}
	return ids
}

//
//  GetIdsFromModelCache
//  @Description: 从 ModelCache 中获取 id
//  @param rows 实现 ModelCache 的切片
//  @return []int64
//
func GetIdsFromModelCache(rows interface{}) []int64 {
	ids := make([]int64, 0)
	if golibrary.IsSlice(rows) {
		s := reflect.ValueOf(rows)
		for i := 0; i < s.Len(); i++ {
			if e, ok := s.Index(i).Interface().(ModelCache); ok {
				if !golibrary.IsInInt64(ids, e.GetId()) {
					ids = append(ids, e.GetId())
				}
			}
		}
	}
	return ids
}

type MapInt64String map[int64]string

/*
使用示例
func (c *DbAction) GetArcTypeTitles(ids []int64) (dd gfcommon.MapInt64String, err error) {
	var rows []model.ArcType
	if err = c.ma.ListWith(c.ma.NewWhere().In("id", ids).Finish(), &rows, func(se *xorm.Session) {
		se.Cols("id", "title")
	}); err != nil {
		return
	}
	var ok bool
	if dd, ok = gfcommon.NewMapInt64String(rows); !ok {
		return
	}
	return
}
*/

func NewMapInt64String(rows interface{}) (dd MapInt64String, ok bool) {
	if golibrary.IsSlice(rows) {
		s := reflect.ValueOf(rows)
		dd = make(map[int64]string, s.Len())
		for i := 0; i < s.Len(); i++ {
			if e, ok := s.Index(i).Interface().(ModelCache); ok {
				dd.Set(e.GetId(), e.GetTitle())
			}
		}
		return dd, true
	}
	return dd, false
}

func (c MapInt64String) Get(id int64, def ...string) string {
	if e, ok := c[id]; ok {
		return e
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (c *MapInt64String) Set(id int64, data string) {
	(*c)[id] = data
}

type MapInt64Interface map[int64]interface{}

func (c MapInt64Interface) Get(id int64) (ok bool, data interface{}) {
	if e, ok := c[id]; ok {
		return true, e
	}
	return false, nil
}

func (c *MapInt64Interface) Set(id int64, data interface{}) {
	(*c)[id] = data
}
