package db

import (
	"time"

	"github.com/fushuilu/golibrary"
	"xorm.io/xorm"
	"xorm.io/xorm/caches"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	LeftJoin = "LEFT OUTER" // 记录列没有值，使用 null 表示
	// 通常不会使用 CROSS
)

/// 事务包裹器
func TransactionWrapper(eg *xorm.EngineGroup, runFunc func(se *xorm.Session) error) (err error) {
	session := eg.NewSession()
	if err = session.Begin(); err != nil {
		return
	}
	defer session.Close()

	err = runFunc(session)
	if err != nil {
		_ = session.Rollback()
		return
	}
	err = session.Commit()
	return
}

func MapCache(eg *xorm.EngineGroup, cache *caches.LRUCacher, beans ...interface{}) error {
	for _, v := range beans {
		if err := eg.MapCacher(v, cache); err != nil {
			return err
		}
	}
	return nil
}

// 载入主从配置
// https://gobook.io/read/gitea.com/xorm/manual-zh-CN/chapter-01/1.engine.html
/*
_ "github.com/go-sql-driver/mysql"
engine, err = xorm.NewEngine("mysql", "root:123@/test?charset=utf8")

_ "github.com/mattn/go-sqlite3"
engine, err = xorm.NewEngine("sqlite3", "./test.db")

[db]
    name = "postgres"
    master = "postgres://auth:123456@localhost/fsl_sso3?sslmode=disable"
    slaver = ["postgres://auth:123456@localhost/fsl_sso3?sslmode=disable"] # 数组
    maxIdleConn = 0
    maxOpenConn = 10
*/

type Config struct {
	Name        string   `json:"name"`
	Cons        []string `json:"cons"` // 第一个默认为 master
	MaxIdleConn int      `json:"maxIdleConn"`
	MaxOpenConn int      `json:"maxOpenConn"`
}

func CreateEngineGroup(conf Config) *xorm.EngineGroup {
	eg, err := xorm.NewEngineGroup(conf.Name, conf.Cons)
	golibrary.PanicIfError(err)

	if conf.MaxIdleConn > 0 {
		eg.SetMaxIdleConns(conf.MaxIdleConn)
	}
	if conf.MaxOpenConn > 0 {
		eg.SetMaxOpenConns(conf.MaxOpenConn)
	}
	eg.DatabaseTZ = time.Local
	eg.TZLocation = time.Local

	err = eg.Ping()
	golibrary.PanicIfError(err)
	return eg
}

func CreateEngineGroupWith(name string, dataSource ...string) *xorm.EngineGroup {
	return CreateEngineGroup(Config{
		Name: name, Cons: dataSource,
	})
}
