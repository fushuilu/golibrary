package db

import (
	"time"

	"github.com/fushuilu/golibrary"
	"github.com/fushuilu/golibrary/lerror"
	"xorm.io/xorm"
)

type Migration struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"-" xorm:"created"`                           // CreateAt 创建时间
	Version   string    `json:"version" xorm:"unique notnull default('')"`  // 日期
	Upgrade   string    `json:"upgrade" xorm:"upgrade notnull default('')"` // 更新的内容
}

func (Migration) TableName() string {
	return "internal_migration"
}

type MigrationAction struct {
	eg *xorm.EngineGroup
}

type MigrationHandler func(se *xorm.Session) error

var todoSyncMigration = true

func NewMigration(eg *xorm.EngineGroup) MigrationAction {
	if todoSyncMigration {
		todoSyncMigration = false
		err := eg.Sync2(Migration{})
		golibrary.PanicIfError(err)
	}
	return MigrationAction{eg: eg}
}

func (m *MigrationAction) Eg() *xorm.EngineGroup {
	return m.eg
}

func (m *MigrationAction) Version(version string, handler MigrationHandler, upgrade ...string) error {
	model := Migration{Version: version}
	if exist, err := m.eg.Where("version=?", version).Exist(&model); err != nil {
		return err
	} else if exist {
		return nil
	}

	se := m.eg.NewSession()
	defer se.Close()
	if err := se.Begin(); err != nil {
		return lerror.Wrap(err, "事务开启失败")
	}

	if err := handler(se); err != nil {
		_ = se.Rollback()
		return lerror.Wrap(err, "处理事务失败")
	}
	// 添加进版本号
	if upgrade != nil {
		model.Upgrade = upgrade[0]
	}
	if _, err := se.Insert(&model); err != nil {
		_ = se.Rollback()
		return lerror.Wrap(err, "添加版本号失败")
	}

	if err := se.Commit(); err != nil {
		return lerror.Wrap(err, "提交事务失败")
	}
	return nil
}
