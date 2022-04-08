少量数据可以使用 `pg_gist` 代替


## [腾讯地图](https://lbs.qq.com/)

* [控制台](https://lbs.qq.com/dev/console/user/info)
* [地点云开发指南](https://lbs.qq.com/service/placeCloud/placeCloudGuide/cloudOverview)
* [地点云数据管理](https://lbs.qq.com/dev/console/dataManage/myData)
* [坐标拾取](https://lbs.qq.com/getPoint/)
* [服务端签名](https://lbs.qq.com/FAQ/server_faq.html#4)

使用步骤

1. 创建数据表 `没有 api, 需要登录地点云数据管理`
2. 存入数据
3. 使用云搜索与地图生成

当前数据表自定义字段为

```
x.remark 备注 string
x.summary 简介 string
x.thumb 缩略图 string
x.kind 类型 number 可搜索
```

## 高德地图(https://lbs.amap.com/)


