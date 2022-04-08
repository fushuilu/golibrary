package gist

import (
	"errors"
	"fmt"
	"github.com/fushuilu/golibrary/appx/datax"
	"github.com/fushuilu/golibrary/appx/db"
	"github.com/fushuilu/golibrary/lerror"
	"strings"
	"xorm.io/builder"
	"xorm.io/xorm/names"
	"xorm.io/xorm/schemas"

	"xorm.io/xorm"
)

// 注意：要使用 ll_to_earth ，必须先在当前 DB 中执行以下命令，安装扩展
// https://johanndutoit.net/searching-in-a-radius-using-postgres/
//<database-name>=# CREATE EXTENSION cube;
//<database-name>=# CREATE EXTENSION earthdistance;
// 表名
// 经度、纬度列名
func CreateGistExtensions(eg *xorm.EngineGroup) error {
	for _, v := range []string{"cube", "earthdistance"} {
		if _, err := eg.Exec(fmt.Sprintf(`CREATE EXTENSION IF NOT EXISTS "%s"`, v)); err != nil {
			return lerror.Wrap(err, fmt.Sprintf("create gist extension:%s error", v))
		}
	}
	return nil
}

//
//  CreateGistIndex
//  @Description: 创建索引; 由于 xorm.Sync2 的特性，需要单独 Sync2 并忽略错误
//  @param tableName string 表名，注意表中必须存在 lat (float64), lng (float64), distance(int) 列
//  @return error
//
func CreateGistIndex(eg *xorm.EngineGroup, tableName string) error {
	// create table tbl_point(id serial8, poi point);
	// create index IF NOT EXISTS idx_tbl_point on tbl_point using gist(poi);
	// CREATE INDEX indexName on events USING gits(ll_to_earth(lat, lng));
	indexName := gistIndexName(tableName)
	sql := fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s ON %s using gist(ll_to_earth(%s, %s));`,
		indexName, tableName, "lat", "lng")
	if _, err := eg.Exec(sql); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("create gist index(%s) for table(%s) error", indexName, tableName), sql)
	}
	return nil
}

func RemoveGistIndex(eg *xorm.EngineGroup, tableName string) error {
	indexName := gistIndexName(tableName)
	sql := fmt.Sprintf("drop index %s", indexName)
	if _, err := eg.Exec(sql); err != nil {
		return lerror.Wrap(err, fmt.Sprintf("drop gist index(%s) for table(%s) error",
			indexName, tableName))
	}
	return nil
}

// xorm 会忽略 _pkey 后缀的索引
func gistIndexName(tableName string) string {
	return fmt.Sprintf("idx_%v_gist_pkey", tableName)
}

/*
// 需要在表中添加以下字段
Lat          float64 `json:"lat" xorm:"notnull default(0)"`                          // 纬度 latitude 0~90
Lng          float64 `json:"lng" xorm:"notnull default(0)"`                          // 经度 longitude 0~180
Distance     int     `json:"distance" xorm:"distance notnull default(0)"`            // 只读：计算出来的距离(冗余的字段)
*/
// Deprecated: 直接使用 eg.Sync2 即可
func Sync(eg *xorm.EngineGroup, bean interface{}, version string) (err error) {
	if eg.Dialect().URI().DBType != schemas.POSTGRES {
		return errors.New("only support postgres engine")
	}
	migration := db.NewMigration(eg)
	tableName := bean.(names.TableName).TableName()

	return migration.Version(version, func(se *xorm.Session) error {

		indexName := gistIndexName(tableName)

		if yes, err := db.HasIndex(eg, tableName, indexName); err != nil {
			return lerror.Wrap(err, "检查 gist 索引时错误")
		} else if yes {
			if err = RemoveGistIndex(eg, tableName); err != nil {
				return err
			}
		} else {
			fmt.Println("|<==== gist index not exists:", tableName, indexName)
		}

		if err = eg.Sync2(bean); err != nil {
			return lerror.Wrap(err, "sync gist table failed")
		}

		return CreateGistIndex(eg, tableName)
	}, fmt.Sprintf("sync gist table:%s", tableName))
}

//
//  Nearest
//  @Description: 查询周边范围内的记录
//  @param eg *xorm.EngineGroup 数据库连接句柄
//  @param rows 保存记录，如果 cols = nil, 则可以直接使用 []gist.Record
//  @param center 中心点
//  @param miles 查询距离，范围为米
//  @param cond 查询条件
//  @param pag 分页
//  @param cols 除了 id, distance 外的列
//  @return err
//
func (pg *PgGist) NearestWith(center datax.Point, miles int, cond builder.Cond, pag db.Pagination, rows interface{}, cols ...string) (err error) {
	cols = append(cols, fmt.Sprintf("id, earth_distance(ll_to_earth(lat,lng), ll_to_earth(%f,%f)) as distance", center.Lat, center.Lng))
	err = pg.eg.Where(cond).
		And(fmt.Sprintf("earth_distance(ll_to_earth(lat, lng), ll_to_earth(%f, %f)) < %d", center.Lat, center.Lng, miles)).
		Select(strings.Join(cols, ",")).
		OrderBy("distance ASC").
		Limit(pag.Limit(), pag.Offset()).
		Find(rows)
	return
}

// xorm 单独对 PgSql 的支持

type PgGist struct {
	eg *xorm.EngineGroup
}

func NewPgGist(eg *xorm.EngineGroup) (gist PgGist) {
	return PgGist{
		eg: eg,
	}
}

type Record struct {
	Id       int64   `json:"id"`
	Distance float64 `json:"distance"` // 计算出来的距离，不可以 xorm:"-" 否则取不到值
}

//
//  NearestGistRecords
//  @Description: 查询附近的记录
//  @param table 表名
//  @param point 查询的点
//  @param miles 距离米
//  @param cond 条件
//  @param pag 分页
//  @return rows
//  @return err
//
func (pg *PgGist) Nearest(tableNameOrBean interface{}, center datax.Point, miles int, cond builder.Cond, pag db.Pagination) (rows []Record, err error) {
	err = pg.eg.Table(tableNameOrBean).Alias("z").
		Where(cond).
		And(fmt.Sprintf("earth_distance(ll_to_earth(z.lat, z.lng), ll_to_earth(%f, %f)) < %d", center.Lat, center.Lng, miles)).
		//And(fmt.Sprintf("earth_box(ll_to_earth(z.lat, z.lng),%f) @> ll_to_earth(%f, %f)", miles, point.Lat, point.Lng)).
		Select(fmt.Sprintf("id, earth_distance(ll_to_earth(z.lat,z.lng), ll_to_earth(%f,%f)) as distance", center.Lat, center.Lng)).
		OrderBy("distance ASC").
		Limit(pag.Limit(), pag.Offset()).
		Find(&rows)
	return
}

func (pg *PgGist) NearestCount(tableNameOrBean interface{}, center datax.Point, miles int, cond builder.Cond) (count int64, err error) {
	if count, err = pg.eg.Table(tableNameOrBean).Alias("z").
		Where(cond).
		And(fmt.Sprintf("earth_distance(ll_to_earth(z.lat, z.lng), ll_to_earth(%f, %f)) < %d", center.Lat, center.Lng, miles)).
		Count(); err != nil {
		err = lerror.Wrap(err, "统计记录错误")
	}
	return
}

// 从 NearestGistRecords 的查询结果中提取出 id
func (pg *PgGist) GetGistRecordsIds(rows []Record) []int64 {
	ids := make([]int64, len(rows))
	for i := range rows {
		ids[i] = rows[i].Id
	}
	return ids
}

// 从 NearestGistRecords 的查询结果中，查询指定 id 记录的距离
func (pg *PgGist) GetGistRecordDistance(id int64, rows []Record) float64 {
	for i := range rows {
		if rows[i].Id == id {
			return rows[i].Distance
		}
	}
	return 0
}

/*
https://rextester.com/WQAY4056

// 查询离用户最近的记录
SELECT events.id, events.name, earth_distance(
	ll_to_earth(user.lat, user.lng),
	ll_to_earth(events.lat, events.lng)
)
as distanceFromCurrentLocation FROM events
ORDER BY distanceFromCurrentLocation ASC
LIMIT 10 OFFSET 0;

// 指定范围内的记录
SELECT events.id, events.name FROM events
WHERE earth_box(
	 ll_to_earth(user.lat, user.lng), radiusInMetres
) @> ll_to_earth(events.lat, events.lng)
*/
