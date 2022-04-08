# 数据库备份

从 `https://github.com/keighl/barkup` 复制过来的，因为部分依赖无法下载

使用示例


```toml
[pgsql]
    Host="127.0.0.1"
    Port=5432
    DB="db name"
    Username="user name"
    Password="password"
    Path="/backup/" # 移动到此目录下，绝对路径
    Options = ["--no-owner"]
    Exclude = ["demos"] # 排除的表
    Format = "sql"
```

```
func main() {

	var pg dbBackup.Postgres
	if err := g.Config().GetStruct("pgsql", &pg); err != nil {
		panic(err)
	}
	if err := pg.Invalid(); err != nil {
		panic(err)
	}

	rst := pg.Export()
	if rst.Error != nil {
		panic(rst.Error)
	}

	// 移动到指定目录
	path := g.Config().GetString("pgsql.Path", "")
	if path != "" {
		moveRst := rst.To(path, nil)
		if moveRst != nil {
			panic(moveRst)
		}
	}

	glog.Debug("rst:[MIME:", rst.MIME, "][Path:", rst.Path, "][Filename:", rst.Filename(), "]")
}
```

