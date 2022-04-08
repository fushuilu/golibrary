package db

import (
	"fmt"
	"strings"

	"github.com/fushuilu/golibrary/lerror"
	"xorm.io/xorm"
)

// 只支持 PgSQL，不支持 MySQL

// 检查是否有列
func HasCol(eg *xorm.EngineGroup, table, field string, schema ...string) (bool, error) {
	s := "public"
	if schema != nil {
		s = schema[0]
	}
	sql := fmt.Sprintf(`SELECT column_name
FROM information_schema.columns
WHERE table_schema = '%s'
AND table_name = '%s';`, s, table)

	if rows, err := eg.QueryString(sql); err != nil {
		return false, err
	} else {
		for _, v := range rows {
			if v["column_name"] == field {
				return true, nil
			}
		}
	}
	return false, nil
}

func HasIndex(eg *xorm.EngineGroup, tableName string, idxName string) (bool, error) {
	session := eg.NewSession()
	return hasIndex2(session, tableName, idxName)
}

// 暂时不支持命名空间
func hasIndex2(session *xorm.Session, tableName, idxName string) (bool, error) {
	rows, err := session.QueryString(`SELECT indexname FROM pg_indexes WHERE tablename = ? AND indexname = ?`, tableName, idxName)
	//rows, err := session.QueryString(`SELECT indexname FROM pg_indexes WHERE tablename = ?`, tableName)
	if err != nil {
		return false, err
	}
	//fmt.Println("rows:", rows)
	return len(rows) == 1, nil
}

func RefreshPGDBIdSeq(session *xorm.Session, tableName string) error {
	sql := fmt.Sprintf("select setval('%s_id_seq', (select max(id) from %s))", tableName, tableName)
	_, err := session.Exec(sql)
	return err
}

/// 清空 SQLite 表
func Truncate(eg *xorm.EngineGroup, tables ...string) error {
	for _, table := range tables {
		if _, err := eg.Exec("DELETE FROM " + table); err != nil {
			return err
		}
	}
	return nil
}

// pgsql 经常会出现 id 错误，烦
func PgTablesIdSeq(eg *xorm.EngineGroup) (string, error) {
	sql := `select relname as tabname,cast(obj_description(relfilenode,'pg_class') as varchar) as comment from pg_class c 
where  relkind = 'r' and relname not like 'pg_%' and relname not like 'sql_%' order by relname`

	if rows, err := eg.QueryString(sql); err != nil {
		return "", lerror.Wrap(err, "query pgsql tables failed")
	} else {
		vals := make([]string, len(rows))
		for i := range rows {
			table := rows[i]["tabname"]
			vals[i] = fmt.Sprintf("select setval('%s_id_seq', (select max(id) from %s));", table, table)
		}
		return strings.Join(vals, "\n"), nil
	}

}
