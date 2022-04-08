package db

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)
import _ "github.com/lib/pq"

var (
	eg = CreateEngineGroup(Config{
		Name: "postgres",
		Cons: []string{
			"postgres://shutao:123456@localhost/pg_demo?sslmode=disable",
		},
		MaxIdleConn: 0,
		MaxOpenConn: 10,
	})
)

type PgDemo struct {
	Id   int64  `json:"id"`
	Name string `json:"name" xorm:"index name notnull default('')"`
}

func (PgDemo) TableName() string {
	return "pg_demo"
}

func TestPGDB(t *testing.T) {
	if err := eg.Ping(); err != nil {
		fmt.Println("ping test database failed, skip the test.")
		return
	}
	err := eg.Sync2(PgDemo{})
	assert.Nil(t, err)

	info, err := eg.TableInfo(&PgDemo{})
	assert.Nil(t, err)

	fmt.Println("info", info)

	tableName := PgDemo{}.TableName()

	//eg.ShowSQL(true)
	hasCol, err := HasCol(eg, tableName, "name")
	assert.Nil(t, err)
	assert.True(t, hasCol)

	indexName := "IDX_pg_demo_name"
	hasIndex, err := HasIndex(eg, tableName, indexName)
	assert.Nil(t, err)
	assert.True(t, hasIndex)

	indexName = "IDX_pg_demo_none"
	hasIndex, err = HasIndex(eg, tableName, indexName)
	assert.Nil(t, err)
	assert.False(t, hasIndex)
}
