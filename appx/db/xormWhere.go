package db

import (
	"fmt"
	"github.com/fushuilu/golibrary"
	"strings"
	"time"

	"xorm.io/builder"
	"xorm.io/xorm"
)

type WhereElement struct {
	Name  string      // 字段名称
	Cmp   string      // 比较符号 特殊的有（IN）
	Value interface{} // 值
}

type Where struct {
	session *xorm.Session
	stack   []WhereElement // 方便复制条件，比如在联表的 count 跟 find 中，第二次 builder.Cond 可能就需要前缀
}

/// 使用注意
/// 哪些情况下需要调用 UnDeleted ?
/// 使用了 xorm:deleted，并且使用 .Table(tableName) 之类的查询，即非数据表，间接的结构体查询
func NewWhere(eg *xorm.EngineGroup) *Where {
	return NewWhereWithSession(eg.Where("1=1"))
}

func NewWhereWithSession(se *xorm.Session) *Where {
	return &Where{session: se}
}

func (w *Where) In(column string, args interface{}) *Where {
	if args != nil {
		//w.session.In(column, args)
		w.stack = append(w.stack, WhereElement{
			Name:  column,
			Cmp:   "IN",
			Value: args,
		})
	}
	return w
}

func (w *Where) Int64(name string, v int64) *Where {
	if v > 0 {
		w.stack = append(w.stack, WhereElement{
			Name:  name,
			Cmp:   "=?",
			Value: v,
		})
	}
	return w
}

func (w *Where) Int(name string, v int) *Where {
	if v > 0 {
		w.stack = append(w.stack, WhereElement{
			Name:  name,
			Cmp:   "=?",
			Value: v,
		})
	}
	return w
}

func (w *Where) IntGTE(name string, v int) *Where {
	w.stack = append(w.stack, WhereElement{
		Name:  name,
		Cmp:   ">=?",
		Value: v,
	})
	return w
}
func (w *Where) IntLTE(name string, v int) *Where {
	w.stack = append(w.stack, WhereElement{
		Name:  name,
		Cmp:   "<=?",
		Value: v,
	})
	return w
}

func (w *Where) Int64GTE(name string, v int64) *Where {
	w.stack = append(w.stack, WhereElement{
		Name:  name,
		Cmp:   ">=?",
		Value: v,
	})
	return w
}
func (w *Where) Int64LTE(name string, v int64) *Where {
	w.stack = append(w.stack, WhereElement{
		Name:  name,
		Cmp:   "<=?",
		Value: v,
	})
	return w
}

// And name=?(v)
func (w *Where) And(name string, v interface{}) *Where {
	w.stack = append(w.stack, WhereElement{
		Name:  name,
		Cmp:   "=?",
		Value: v,
	})
	return w
}

// 替换值
func (w *Where) Replace(name string, v interface{}, cmp ...string) *Where {
	for i := range w.stack {
		if name == w.stack[i].Name {
			w.stack[i].Value = v
			if len(cmp) > 0 {
				w.stack[i].Cmp = cmp[0]
			}
		}
	}
	return w
}

func (w *Where) True(name string, condition bool) *Where {
	if condition {
		w.stack = append(w.stack, WhereElement{
			Name:  name,
			Cmp:   "=?",
			Value: true,
		})
	}
	return w
}

func (w *Where) False(name string, condition bool) *Where {
	if !condition {
		w.stack = append(w.stack, WhereElement{
			Name:  name,
			Cmp:   "=?",
			Value: false,
		})
	}
	return w
}

func (w *Where) MockBoolean(name string, condition string) *Where {
	switch strings.ToLower(condition) {
	case "true", "yes", "t":
		w.stack = append(w.stack, WhereElement{
			Name:  name,
			Cmp:   "=?",
			Value: true,
		})
	case "false", "no", "f":
		w.stack = append(w.stack, WhereElement{
			Name:  name,
			Cmp:   "=?",
			Value: false,
		})
	}
	return w
}

// name >= t
func (w *Where) TimeGTE(name string, t time.Time) *Where {
	if !t.IsZero() {
		w.stack = append(w.stack, WhereElement{
			Name:  name,
			Cmp:   ">=?",
			Value: t,
		})
	}
	return w
}

// name <= t
func (w *Where) TimeLTE(name string, t time.Time) *Where {
	if !t.IsZero() {
		w.stack = append(w.stack, WhereElement{
			Name:  name,
			Cmp:   "<=?",
			Value: t,
		})
	}
	return w
}

// name > t
func (w *Where) TimeGT(name string, t time.Time) *Where {
	if !t.IsZero() {
		w.stack = append(w.stack, WhereElement{
			Name:  name,
			Cmp:   ">?",
			Value: t,
		})
	}
	return w
}

// name < t
func (w *Where) TimeLT(name string, t time.Time) *Where {
	if !t.IsZero() {
		w.stack = append(w.stack, WhereElement{
			Name:  name,
			Cmp:   "<?",
			Value: t,
		})
	}
	return w
}

// 日期记录统计
func (w *Where) TimeDate(name string, t time.Time, days int) *Where {
	begin := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	return w.TimeGTE(name, begin).TimeLT(name, begin.AddDate(0, 0, days))
}

func (w *Where) Like(name, keyword string) *Where {
	if keyword != "" {
		w.stack = append(w.stack, WhereElement{
			Name:  name,
			Cmp:   " LIKE ?",
			Value: "%" + keyword + "%",
		})
	}
	return w
}

func (w *Where) LeftLike(name, keyword string) *Where {
	if keyword != "" {
		w.stack = append(w.stack, WhereElement{
			Name:  name,
			Cmp:   " LIKE ?",
			Value: keyword + "%",
		})
	}
	return w
}

func (w *Where) RightLike(name, keyword string) *Where {
	if keyword != "" {
		w.stack = append(w.stack, WhereElement{
			Name:  name,
			Cmp:   " LIKE ?",
			Value: "%" + keyword,
		})
	}
	return w
}

func (w *Where) String(name, v string) *Where {
	if v != "" {
		w.stack = append(w.stack, WhereElement{
			Name:  name,
			Cmp:   "=?",
			Value: v,
		})
	}
	return w
}

// Element 示例 ("id<?", 5) => id<5
func (w *Where) Element(query string, value interface{}) *Where {
	w.stack = append(w.stack, WhereElement{
		Name:  query,
		Cmp:   "",
		Value: value,
	})
	return w
}

/// not equal
func (w *Where) NotEqual(name string, value interface{}) *Where {
	w.stack = append(w.stack, WhereElement{
		Name:  name,
		Cmp:   "!=?",
		Value: value,
	})
	return w
}
func (w *Where) Today(colName string) *Where {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	w.GetSession().And(fmt.Sprintf("%s >= ? AND %s < ?", colName, colName),
		start, start.AddDate(0, 0, 1),
	)
	return w
}

// 注意，这个无法添加到前缀中，可以使用 Element 代替
func (w *Where) Session(run func(se *xorm.Session)) *Where {
	run(w.session)
	return w
}

func (w *Where) OrIn(v []int64, col string, cols ...string) {
	if len(cols) == 0 {
		w.In(col, v)
	} else {
		cols = append(cols, col)
		values := golibrary.Int64SToString(v)
		w.session.And(joinOrSQL(values, cols))
	}
}

func joinOrSQL(value string, cols []string) string {
	formatSQL := make([]string, len(cols))
	for i := range cols {
		formatSQL[i] = fmt.Sprintf("%s IN (%s)", cols[i], value)
	}
	return strings.Join(formatSQL, " OR ")
}

// 不保存到 stack 中
func (w *Where) Or(v interface{}, col string, cols ...string) *Where {

	switch v.(type) {
	case string:
		if v.(string) == "" || strings.TrimSpace(v.(string)) == "" {
			return w
		}
	case int:
		if v.(int) == 0 {
			return w
		}
	case int64:
		if v.(int64) == 0 {
			return w
		}
	default:
		fmt.Println("warning: nonsupport type in where.Or")
		return w
	}
	totalLen := 1 + len(cols)
	keys := make([]string, totalLen)
	keys = append(keys, fmt.Sprintf("%s=?", col))

	vs := make([]interface{}, totalLen)
	vs = append(vs, v)
	for _, v1 := range cols {
		keys = append(keys, fmt.Sprintf("%s=?", v1)) // 列
		vs = append(vs, v)                           // 值
	}

	w.session.And(strings.Join(keys, " OR "), vs...)

	return w
}

func (w *Where) GetSession() *xorm.Session {
	return w.session
}

func (w *Where) UnDeleted(cols ...string) *Where {
	var name = "deleted"
	if len(cols) > 0 {
		name = cols[0]
	}
	w.stack = append(w.stack, WhereElement{
		Name:  fmt.Sprintf("%s=? OR %s IS NULL", name, name), // 这里使用大写
		Cmp:   "",
		Value: "0001-01-01 00:00:00",
	})
	return w
}

func (w *Where) IncludeDeleted() *Where {
	w.session.Unscoped()
	return w
}

// 保存并清空全部条件
func (w *Where) Save() *Where {
	for _, v := range w.stack {
		if v.Cmp == "IN" {
			w.session.In(v.Name, v.Value)
			continue
		}
		w.session.And(v.Name+v.Cmp, v.Value)
	}
	w.stack = []WhereElement{}
	return w
}

// 使用 builder.ToSQL 可以打印 builder.Cond
func (w *Where) Finish() builder.Cond {
	return w.Save().session.Conds()
}

// 复制条件
func (w *Where) CloneStack(from *Where) *Where {
	w.stack = from.stack
	return w
}

func (w *Where) Done() *xorm.Session {
	return w.session.Where(w.Finish())
}

func NewBuildCond() builder.Cond {
	return builder.NewCond()
}
