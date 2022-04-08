package dbx

import (
	"fmt"
	"github.com/fushuilu/golibrary/appx/datax"
	"github.com/fushuilu/golibrary/appx/db"
	"github.com/fushuilu/golibrary/appx/errorx"
	"github.com/fushuilu/golibrary/lerror"

	"xorm.io/builder"
	"xorm.io/xorm"
)

type ModelAction struct {
	ma db.ModelAction
}

func NewModelAction(ma db.ModelAction) ModelAction {
	return ModelAction{ma: ma}
}

func (a *ModelAction) DBMa() *db.ModelAction {
	return &a.ma
}

func (a *ModelAction) NewWhere() *db.Where  {
	return a.ma.NewWhere()
}

func (a *ModelAction) Eg() *xorm.EngineGroup {
	return a.ma.Eg()
}

// 修改系统状态，注意模型必须要有 state_index 和 state_msg 字段
func (a *ModelAction) ChangeStateWith(cond builder.Cond, data datax.ChangeState) error {
	if data.Id < 1 {
		return lerror.New(fmt.Sprintf("修改 %s 记录系统状态时 id 为空", a.ma.Title()))
	}
	if _, err := a.ma.Eg().Table(a.ma.TableName()).Where(cond).Cols("state_index", "state_msg").
		Update(&struct {
			StateIndex int    `json:"state_index" xorm:"state_index notnull default(0)"`
			StateMsg   string `json:"state_msg" xorm:"state_msg notnull default('')"`
		}{StateIndex: data.StateIndex, StateMsg: data.StateMsg}); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("更新%s记录系统状态时错误", a.ma.Title()))
	} else {
		return nil
	}
}

// 修改系统状态，注意模型必须要有 state_index 和 state_msg 字段
func (a *ModelAction) ChangeState(data datax.ChangeState) error {
	return a.ChangeStateWith(a.ma.NewWhere().Int64("id", data.Id).Finish(), data)
}

func (a *ModelAction) ChangeStatusWith(cond builder.Cond, data datax.ChangeStatus) error {
	if data.Id < 1 {
		return errorx.New(fmt.Sprintf("修改 %s 记录状态时 id 为空", a.ma.Title()))
	}

	if _, err := a.ma.Eg().Table(a.ma.TableName()).Where(cond).Cols("status_index").
		Update(&struct {
			StatusIndex int `json:"status_index" xorm:"status_index notnull default(0)"`
		}{StatusIndex: data.StatusIndex}); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("更新%s记录状态时错误", a.ma.Title()))
	} else {
		return nil
	}
}

func (a *ModelAction) ChangeStatus(data datax.ChangeStatus) error {
	return a.ChangeStatusWith(a.ma.NewWhere().
		Int64("id", data.Id).
		Int64("user_id", data.UserId).Finish(), data)
}

//
//  InsertOne
//  @Description: 添加一条记录
//  @param bean 模型
//  @param callback 回调函数 id, colName
//  @return error
//
func (a *ModelAction) InsertOne(bean db.Model, callback func() (int64, string)) error {
	return db.TransactionWrapper(a.ma.Eg(), func(se *xorm.Session) error {
		if _, err := se.Table(a.ma.Title()).InsertOne(bean); err != nil {
			return lerror.Wrap(err, fmt.Sprintf("添加 %s 记录时错误", a.ma.Title()))
		} else if !bean.IsGet() {
			return errorx.New(fmt.Sprintf("添加 %s 时错误", a.ma.Title()))
		}
		id, name := callback()
		if id < 1 {
			return errorx.New("必须返回记录主键 ID")
		} else if name == "" {
			return errorx.New("必须指定更新的字符串 id")
		}
		num, err := se.Table(a.ma.TableName()).ID(id).Cols(name).Update(bean)
		return a.ma.UpdateRst(num, err)

	})
}
