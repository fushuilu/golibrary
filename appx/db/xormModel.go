package db

import (
	"errors"
	"fmt"
	"github.com/fushuilu/golibrary"
	"github.com/fushuilu/golibrary/appx/datax"
	"github.com/fushuilu/golibrary/appx/errorx"
	"github.com/fushuilu/golibrary/lerror"

	"xorm.io/builder"
	"xorm.io/xorm"
)

type Model interface {
	TableName() string
	RecordFormat()
	Invalid(isCreated bool) error
	IsGet() bool
}

type ModelAction struct {
	eg    *xorm.EngineGroup
	bean  Model  // 模型
	title string // 模型名称
}

// 注意: bean 必须是一个 point，例如 &model.Page{}
func NewModelAction(eg *xorm.EngineGroup, title string, bean Model) ModelAction {
	if !golibrary.IsPointer(bean) {
		panic(fmt.Sprintf("NewModelAction[%s]'s param[bean] is not a pointer", title))
	}
	return ModelAction{eg: eg, title: title,
		bean: bean}
}

func (a *ModelAction) TableName() string {
	return a.bean.TableName()
}

func (a *ModelAction) NewWhere() *Where {
	return NewWhere(a.eg)
}

func (a *ModelAction) Eg() *xorm.EngineGroup {
	return a.eg
}

func (a *ModelAction) SetTitle(title string) {
	a.title = title
}

func (a *ModelAction) Title() string {
	return a.title
}

func (a *ModelAction) Bean() Model {
	return a.bean
}

func (a *ModelAction) MustId(id int64) error {
	if id < 1 {
		return errors.New(fmt.Sprintf("必须指定 %s 记录的 ID", a.title))
	}
	return nil
}

func (a *ModelAction) MustExist(cond builder.Cond) error {
	if exist, err := a.eg.Where(cond).NoAutoCondition().NoAutoTime().Exist(a.bean); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("检查%s记录时错误", a.title))
	} else if exist {
		//Debug("mustExist")
		//fmt.Printf("%+v\n", a.bean)
		return nil
	}
	return errorx.New(fmt.Sprintf("待检查%s不存在", a.title))
}

func (a *ModelAction) MustNotExist(cond builder.Cond) error {
	if exist, err := a.eg.Where(cond).NoAutoCondition().NoAutoTime().Exist(a.bean); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("检查%s记录时错误", a.title))
	} else if exist {
		return errorx.New(fmt.Sprintf("重复的%s记录", a.title))
	}
	return nil
}

func (a *ModelAction) Exist(cond builder.Cond) (bool, error) {
	if exist, err := a.eg.Where(cond).NoAutoCondition().NoAutoTime().Exist(a.bean); err != nil {
		return false, lerror.Wrap(err, fmt.Sprintf("检查%s记录时错误", a.title))
	} else {
		return exist, nil
	}
}

func (a *ModelAction) GetOf(cond builder.Cond, bean interface{}, with func(se *xorm.Session)) error {
	where := a.eg.Table(a.bean.TableName()).Where(cond)
	with(where)
	if _, err := where.Get(bean); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("获取%s记录时错误", a.title))
	}
	return nil
}

func (a *ModelAction) Get(cond builder.Cond, bean Model) error {
	if _, err := a.eg.Table(a.bean.TableName()).Where(cond).Get(bean); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("获取%s记录时错误", a.title))
	}
	bean.RecordFormat()
	return nil
}

func (a *ModelAction) MustGet(cond builder.Cond, bean Model) error {
	if err := a.Get(cond, bean); err != nil {
		return err
	}
	if !bean.IsGet() {
		return errorx.New(fmt.Sprintf("没有找到符合要求的%s记录", a.title))
	}
	return nil
}
func (a *ModelAction) MustGetWith(cond builder.Cond, bean Model, with func(se *xorm.Session)) error {
	if err := a.GetWith(cond, bean, with); err != nil {
		return err
	}
	if !bean.IsGet() {
		return errorx.New(fmt.Sprintf("没有找到符合要求的%s记录", a.title))
	}
	return nil
}

func (a *ModelAction) MustGetById(id int64, bean Model) error {
	if id < 1 {
		return errors.New(fmt.Sprintf("查询 %s 记录时 id 为空", a.title))
	}
	if err := a.Get(a.NewWhere().Int64("id", id).Finish(), bean); err != nil {
		return err
	}
	if !bean.IsGet() {
		return errorx.New(fmt.Sprintf("没有找到符合要求的%s记录", a.title))
	}
	return nil
}
func (a *ModelAction) MustGetByIdWith(id int64, bean Model, with func(se *xorm.Session)) error {
	if id < 1 {
		return errors.New(fmt.Sprintf("查询 %s 记录时 id 为空", a.title))
	}
	if err := a.GetWith(a.NewWhere().Int64("id", id).Finish(), bean, with); err != nil {
		return err
	}
	if !bean.IsGet() {
		return errorx.New(fmt.Sprintf("没有找到符合要求的%s记录", a.title))
	}
	return nil
}

func (a *ModelAction) GetWith(cond builder.Cond, bean Model, with func(se *xorm.Session)) error {
	where := a.eg.Table(a.bean.TableName()).Where(cond)
	with(where)
	if _, err := where.Get(bean); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("获取%s记录时错误", a.title))
	}
	bean.RecordFormat()
	return nil
}

func (a *ModelAction) GetById(id int64, bean Model) error {
	if id < 1 {
		return errors.New(fmt.Sprintf("查询 %s 记录时 id 为空", a.title))
	}
	if _, err := a.eg.Table(a.bean.TableName()).ID(id).Get(bean); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("查询%s记录时错误", a.title))
	}
	bean.RecordFormat()
	return nil
}
func (a *ModelAction) GetByIdWith(id int64, bean Model, with func(se *xorm.Session)) error {
	if id < 1 {
		return errors.New(fmt.Sprintf("查询 %s 记录时 id 为空", a.title))
	}
	where := a.eg.Table(a.bean.TableName()).ID(id)
	with(where)
	if _, err := where.Get(bean); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("查询%s记录时错误", a.title))
	}
	bean.RecordFormat()
	return nil
}

func (a *ModelAction) Count(cond builder.Cond) (int64, error) {
	if count, err := a.eg.Where(cond).NoAutoCondition().NoAutoTime().Count(a.bean); err != nil {
		return 0, lerror.Wrap(err, fmt.Sprintf("统计%s记录时错误", a.title))
	} else {
		return count, nil
	}
}

func (a *ModelAction) List(cond builder.Cond, beans interface{}) error {
	if err := a.eg.Table(a.bean.TableName()).Where(cond).After(func(i interface{}) {
		if e, ok := i.(Model); ok {
			e.RecordFormat()
		}
	}).Find(beans); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("查询%s列表记录时错误", a.title))
	}
	return nil
}

func (a *ModelAction) ListWith(cond builder.Cond, beans interface{}, with func(se *xorm.Session)) error {
	where := a.eg.Table(a.bean.TableName()).Where(cond)
	if with != nil {
		with(where)
	}
	if err := where.After(func(i interface{}) {
		if e, ok := i.(Model); ok {
			e.RecordFormat()
		}
	}).Find(beans); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("查询%s列表记录时错误", a.title))
	}
	return nil
}

func (a *ModelAction) Pagination(cond builder.Cond, pag Pagination, beans interface{}) error {
	if err := a.eg.Table(a.bean.TableName()).Where(cond).Limit(pag.Limit(), pag.Offset()).After(func(i interface{}) {
		if e, ok := i.(Model); ok {
			e.RecordFormat()
		}
	}).Find(beans); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("查询%s记录时错误", a.title))
	}
	return nil
}

func (a *ModelAction) PaginationWith(cond builder.Cond, pag Pagination, beans interface{}, with func(se *xorm.Session)) error {
	where := a.eg.Table(a.bean.TableName()).Where(cond).Limit(pag.Limit(), pag.Offset())
	if with != nil {
		with(where)
	}
	if err := where.After(func(i interface{}) {
		if e, ok := i.(Model); ok {
			e.RecordFormat()
		}
	}).Find(beans); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("查询%s列表记录时错误", a.title))
	}
	return nil
}

func (a *ModelAction) ListResult(cond builder.Cond, pag Pagination, beans interface{}) (*datax.ListResult, error) {
	if count, err := a.Count(cond); err != nil {
		return nil, lerror.Wrap(err, fmt.Sprintf("统计 %s 记录时错误", a.title))
	} else if count < 1 {
		return &datax.ListResultEmpty, nil
	} else {
		if err := a.Pagination(cond, pag, beans); err != nil {
			return nil, err
		}
		return &datax.ListResult{Total: count, Rows: beans}, nil
	}
}
func (a *ModelAction) ListResultWith(cond builder.Cond, pag Pagination, beans interface{}, with func(se *xorm.Session)) (*datax.ListResult, error) {
	if count, err := a.Count(cond); err != nil {
		return nil, lerror.Wrap(err, fmt.Sprintf("统计 %s 记录时错误", a.title))
	} else if count < 1 {
		return &datax.ListResultEmpty, nil
	} else {
		if err = a.PaginationWith(cond, pag, beans, with); err != nil {
			return nil, lerror.Wrap(err, "查询记录错误")
		}
		return &datax.ListResult{Total: count, Rows: beans}, nil
	}
}

func (a *ModelAction) Session() *xorm.Session {
	return a.eg.Table(a.bean.TableName())
}

func (a *ModelAction) DeleteWith(cond builder.Cond) error {
	i, err := a.eg.Where(cond).NoAutoCondition().NoAutoTime().Delete(a.bean) // deleted_at 会被赋值
	return a.DeleteRst(i, err)
}

func (a *ModelAction) DeleteOne(cond builder.Cond) error {
	i, err := a.eg.Where(cond).NoAutoCondition().NoAutoTime().Limit(0, 1).Delete(a.bean) // deleted_at 会被赋值
	return a.DeleteRst(i, err)
}

func (a *ModelAction) DeleteById(id int64) error {
	if id < 1 {
		return errorx.New(fmt.Sprintf("移除 %s 记录时 ID 参数错误", a.title))
	}
	return a.DeleteWith(NewWhere(a.eg).Int64("id", id).Finish())
}

func (a *ModelAction) InsertOne(bean interface{}) error {
	_, err := a.eg.InsertOne(bean)
	if err != nil {
		return lerror.Wrap(err, fmt.Sprintf("添加 %s 记录错误", a.title))
	}
	if e, ok := bean.(Model); ok {
		if !e.IsGet() {
			return errorx.New(fmt.Sprintf("添加%s记录失败", a.title))
		}
	}
	return nil
}

func (a *ModelAction) InsertRst(id int64, err error) error {
	if err != nil {
		return lerror.Wrap(err, fmt.Sprintf("添加%s记录时错误", a.title))
	}
	if id < 1 {
		return errorx.New(fmt.Sprintf("添加%s记录失败", a.title))
	}
	return nil
}

func (a *ModelAction) UpdateRst(num int64, err error) error {
	if err != nil {
		return lerror.Wrap(err, fmt.Sprintf("更新%s记录时错误", a.title))
	}
	if num < 1 {
		return errorx.New(fmt.Sprintf("没有任何%s记录被更新", a.title))
	}
	return nil
}

func (a *ModelAction) DeleteRst(num int64, err error) error {
	if err != nil {
		return lerror.Wrap(err, fmt.Sprintf("移除%s记录时错误", a.title))
	}
	if num < 1 {
		return errorx.New(fmt.Sprintf("没有任何%s记录被移除", a.title))
	}
	return nil
}

func (a *ModelAction) GetRst(id int64, err error) error {
	return a.CheckGetRst(id > 0, err)
}

func (a *ModelAction) CheckGetRst(isGet bool, err error) error {
	if err != nil {
		return lerror.Wrap(err, fmt.Sprintf("查询%s记录时错误", a.title))
	}
	if !isGet {
		return errorx.New(fmt.Sprintf("没有找到符合要求的%s", a.title))
	}
	return nil
}

func (a *ModelAction) Incr(id int64, colName string) error {
	if id < 1 {
		return errorx.New(fmt.Sprintf("递增 %s 记录列 %s 时 id 为空", a.title, colName))
	}
	if _, err := a.eg.ID(id).NoAutoCondition().NoAutoTime().
		Cols(colName).Incr(colName).Update(a.bean); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("递增 %s 记录列 %s 时错误", a.title, colName))
	}
	return nil
}

func (a *ModelAction) Decr(id int64, colName string) error {
	if id < 1 {
		return errorx.New(fmt.Sprintf("递减 %s 记录列 %s 时 id 为空", a.title, colName))
	}
	if _, err := a.eg.ID(id).NoAutoCondition().NoAutoTime().
		Cols(colName).Decr(colName).Update(a.bean); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("递减 %s 记录列 %s 时错误", a.title, colName))
	}
	return nil
}
func (a *ModelAction) DecrZero(id int64, colName string) error {
	if id < 1 {
		return errorx.New(fmt.Sprintf("递减 %s 记录列 %s 时 id 为空", a.Title(), colName))
	}
	if _, err := a.Eg().ID(id).And(colName + ">0").NoAutoCondition().NoAutoTime().
		Cols(colName).Decr(colName).Update(a.bean); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("递减 %s 记录列 %s 时错误", a.Title(), colName))
	}
	return nil
}

// 记录 id 集合
func (a *ModelAction) Ids(cond builder.Cond) ([]int64, error) {
	var rows []struct {
		Id int64 `json:"id"`
	}
	if err := a.eg.Table(a.bean.TableName()).Where(cond).Cols("id").Find(&rows); err != nil {
		return datax.EmptyArray, lerror.Wrap(err, fmt.Sprintf("查询 %s id 集合时错误", a.title))
	}
	ids := make([]int64, len(rows))
	for i := range rows {
		ids[i] = rows[i].Id
	}
	return ids, nil
}
