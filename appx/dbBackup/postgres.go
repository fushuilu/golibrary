package dbBackup

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// https://gist.github.com/vielhuber/96eefdb3aff327bdf8230d753aaee1e1
var (
	// PGDumpCmd is the path to the `pg_dump` executable
	PGDumpCmd = "pg_dump"
)

/*
pg_dump常用参数
-h host，指定数据库主机名，或者IP
-p port，指定端口号
-U user，指定连接使用的用户名
-W，按提示输入密码
dbname，指定连接的数据库名称，实际上也是要备份的数据库名称。
-a，–data-only，只导出数据，不导出表结构
-c，–clean，是否生成清理该数据库对象的语句，比如drop table
-C，–create，是否输出一条创建数据库语句
-f file，–file=file，输出到指定文件中
-n schema，–schema=schema，只转存匹配schema的模式内容
-N schema，–exclude-schema=schema，不转存匹配schema的模式内容
-O，–no-owner，不设置导出对象的所有权
-s，–schema-only，只导致对象定义模式，不导出数据
-t table，–table=table，只转存匹配到的表，视图，序列，可以使用多个-t匹配多个表
-T table，–exclude-table=table，不转存匹配到的表。
–inserts，使用insert命令形式导出数据，这种方式比默认的copy方式慢很多，但是可用于将数据导入到非PostgreSQL数据库。
–column-inserts，导出的数据，有显式列名
*/
// Postgres is an `Exporter` interface that backs up a Postgres database via the `pg_dump` command
type Postgres struct {
	// DB Host (e.g. 127.0.0.1)
	Host string
	// DB Port (e.g. 5432)
	Port string
	// DB Name
	DB string
	// Connection Username
	Username string
	// 密码
	Password string
	// Extra pg_dump options
	// e.g []string{"--inserts"}
	Options []string
	Include []string // 只转存这此表 (测试时命令错误)
	Exclude []string // 不转存这些表 (测试时不起作用)
	Format  string   // 格式，如果是 sql ，则保存了 sql
}

func (x *Postgres) Invalid() error {
	if x.Host == "" {
		x.Host = "127.0.0.1"
	}
	if x.Port == "" {
		x.Port = "5432"
	}
	if x.DB == "" {
		return errors.New("postgres DB name is empty")
	}
	if x.Username == "" {
		return errors.New("postgres Username is empty")
	}

	return nil
}

// Export produces a `pg_dump` of the specified database, and creates a gzip compressed tarball archive.
func (x Postgres) Export() *ExportResult {
	result := &ExportResult{MIME: "application/x-tar"}

	options := x.Options

	if x.DB != "" {
		options = append(options, fmt.Sprintf(`-d%v`, x.DB))
	}

	if x.Host != "" {
		options = append(options, fmt.Sprintf(`-h%v`, x.Host))
	}

	if x.Port != "" {
		options = append(options, fmt.Sprintf(`-p%v`, x.Port))
	}

	if x.Username != "" {
		options = append(options, fmt.Sprintf(`-U%v`, x.Username))
	}

	if x.Password != "" {
		if err := os.Setenv("PGPASSWORD", x.Password); err != nil {
			result.Error = makeErr(err, "无法设置 PGPASSWORD 变量")
			result.Error.Cmd = `os.Setenv("PGPASSWORD",xxx)`
			return result
		}
	}

	// -Fc 格式化为定制的格式(c), 明文(p)
	ymd := time.Now().Format("20060102-1504")
	if x.Format == "sql" {
		result.Path = fmt.Sprintf(`pg_%v_%v.sql`, x.DB, ymd)
		options = append(options, "-Fp")
	} else {
		result.Path = fmt.Sprintf(`pg_%v_%v.sql.tar.gz`, x.DB, ymd)
		options = append(options, "-Fc")
	}

	for _, v := range x.Include {
		options = append(options, fmt.Sprintf(`-t %v`, v))
	}

	for _, v := range x.Exclude {
		options = append(options, fmt.Sprintf(`-T %v`, v))
	}

	options = append(options, fmt.Sprintf(`-f%v`, result.Path))

	cmd := exec.Command(PGDumpCmd, options...)
	out, err := cmd.Output()
	if err != nil {
		result.Error = makeErr(err, string(out))
		result.Error.Cmd = cmd.String()
	}
	fmt.Println("OK:", cmd.String())
	return result
}

func (x Postgres) dumpOptions() []string {
	options := x.Options

	if x.DB != "" {
		options = append(options, fmt.Sprintf(`-d%v`, x.DB))
	}

	if x.Host != "" {
		options = append(options, fmt.Sprintf(`-h%v`, x.Host))
	}

	if x.Port != "" {
		options = append(options, fmt.Sprintf(`-p%v`, x.Port))
	}

	if x.Username != "" {
		options = append(options, fmt.Sprintf(`-U%v`, x.Username))
	}

	return options
}
