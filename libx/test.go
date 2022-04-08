package cmn

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fushuilu/golibrary"
	"github.com/fushuilu/golibrary/libx/cachex"

	"xorm.io/xorm"
)

func NewTestSQLite(db ...interface{}) (eg *xorm.EngineGroup) {
	var err error
	eg, err = xorm.NewEngineGroup("sqlite3", []string{"./test.db"},
		xorm.RandomPolicy())
	if err != nil {
		panic(err)
	}
	if len(db) > 0 {
		if err = eg.Sync2(db...); err != nil {
			panic(err)
		}
	}
	return eg
}

func NewTestXRedis() *cachex.XRedis {
	return cachex.NewXRedis(&cachex.RedisOpts{Host: "127.0.0.1:6379", Database: 10})
}

func NewTestPanicTime(err error, created time.Time, message string) {
	if err != nil || created.IsZero() {
		fmt.Println("|<=======", message)
		golibrary.PanicIfError(err)

		if created.IsZero() {
			panic("could not get created time")
		}
	}
}

func NewTestPanicId(err error, id int64, message string) {
	if err != nil || id < 1 {
		fmt.Println("|<=======", message)
		golibrary.PanicIfError(err)

		if id < 1 {
			panic("could not get id")
		}
	}
}

func ReUnmarshal(from interface{}, to interface{}) error {
	d, err := json.Marshal(from)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(d, to); err != nil {
		return err
	}
	return nil
}
