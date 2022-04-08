package qq

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fushuilu/golibrary/appx/db"

	"github.com/fushuilu/golibrary/appx/datax"
)

// https://lbs.qq.com/service/placeCloud/placeCloudGuide/cloudDataManage
const (
	pathPlaceCloudTableList = "/place_cloud/table/list" // get

	pathPlaceCloudDataCreate = "/place_cloud/data/create" // post
	pathPlaceCloudDataUpdate = "/place_cloud/data/update" // post
	pathPlaceCloudDataDelete = "/place_cloud/data/delete" // post
	pathPlaceCloudDataList   = "/place_cloud/data/list"   // 数据列表 get
)

type ResponsePlaceCloudTableInfo struct {
	CreateTime int    `json:"create_time"`  // 表创建时间
	UpdateTime int    `json:"update_time"`  // //表修改时间
	DataCount  int    `json:"data_count"`   // 表中数据量
	TableId    string `json:"table_id"`     // 表的 ID
	TableName  string `json:"table_name"`   // 表名
	UniqueUdId int    `json:"unique_ud_id"` //启用ud_id唯一性约束
	X          map[string]struct {
		CreateTime int    `json:"create_time"` // 字段创建时间
		UpdateTime int    `json:"update_time"` // 字段更新时间
		Comment    string `json:"comment"`     // 字段注释
		Default    int    `json:"default"`     // 默认值
		IsIndex    int    `json:"is_index"`    // 是否为云搜索筛选排序字段
		Type       string `json:"type"`        // 字段类型
	} `json:"x"` // 自定义字段

}

func (c *Lbs) PlaceCloudTableList() (tables []ResponsePlaceCloudTableInfo, err error) {
	var resp struct {
		CommonResponse
		Result struct {
			Tables []ResponsePlaceCloudTableInfo `json:"tables"` // 查询到的有权访问的表及其结果
		} `json:"result"`
	}
	if err = c.http.Get(domain+pathPlaceCloudTableList, map[string]string{
		"key": c.cf.Key,
	}, &resp); err != nil {
		return
	}
	if resp.IsOk() {
		return resp.Result.Tables, nil
	}
	err = resp.Error()
	return
}

type RequestPlaceCloudData struct {
	Id       int64                  // 记录 ID
	Title    string                 // 标题
	Location datax.Point            // 位置(优先)
	Address  string                 // 地址信息
	Tel      string                 // 联系电话
	X        map[string]interface{} // 自定义字段,其中的 kind 可索引
}

func (pd *RequestPlaceCloudData) Invalid() error {
	if pd.Id < 1 {
		return errors.New("必须提供主键 ID")
	}
	if pd.Title == "" {
		return errors.New("必须提供标题")
	}
	return nil
}

func (pd *RequestPlaceCloudData) toDataRow() dataRow {
	return dataRow{
		UdId:     fmt.Sprintf("%d", pd.Id),
		Title:    pd.Title,
		Location: pd.Location,
		Tel:      pd.Tel,
		Address:  pd.Address,
		X:        pd.X,
	}
}

type dataRow struct {
	UdId     string      `json:"ud_id"`
	Title    string      `json:"title"`
	Location datax.Point `json:"location"`
	Tel      string      `json:"tel"`
	Address  string      `json:"address,omitempty"`
	Polygon  string      `json:"polygon,omitempty"`
	X        interface{} `json:"x,omitempty"`
}

func (c *Lbs) getTableId(tableId ...string) (string, error) {
	table := c.cf.TableId
	if len(tableId) > 0 {
		table = tableId[0]
	}
	if table == "" {
		return "", errors.New("必须指定 table id")
	}
	return table, nil
}

func (c *Lbs) PlaceCloudDataCreate(pd RequestPlaceCloudData, tableId ...string) (err error) {
	if err = pd.Invalid(); err != nil {
		return err
	}

	var id string
	if id, err = c.getTableId(tableId...); err != nil {
		return err
	}

	data, _ := json.Marshal(struct {
		Key     string    `json:"key"`
		TableId string    `json:"table_id"`
		Data    []dataRow `json:"data"`
	}{
		Key:     c.cf.Key,
		TableId: id,
		Data: []dataRow{
			pd.toDataRow(),
		},
	})

	var resp struct {
		CommonResponse
		Result struct {
			Count   int `json:"count"` // 成功创建数据条数
			Failure []struct {
				Message string `json:"message"` // 创建该条数据时产生的错误信息
				RowIdx  int    `json:"row_idx"` // 该条数据在data数组中的下标位置（从0开始）
				Status  int    `json:"status"`  // 错误码
				UdId    string `json:"ud_id"`   // 自定义id
			} `json:"failure"`
			Success []struct {
				Id     string `json:"id"`      // 数据创建成功，返回系统生成的唯一标识（id）
				RowIdx int    `json:"row_idx"` // 该条数据在data数组中的下标位置（从0开始
				UdId   string `json:"ud_id"`   // 自定义id
			}
		} `json:"result"`
	}
	err = c.http.PostByte(domain+pathPlaceCloudDataCreate, data, &resp)
	if err != nil {
		return err
	}

	if resp.IsOk() {
		if resp.Result.Count < 1 && len(resp.Result.Failure) > 0 {
			return errors.New(resp.Result.Failure[0].Message)
		}
		return nil
	}
	return resp.Error()
}

func (c *Lbs) PlaceCloudDataUpdate(pd RequestPlaceCloudData, tableId ...string) (err error) {
	if err = pd.Invalid(); err != nil {
		return err
	}

	var id string
	if id, err = c.getTableId(tableId...); err != nil {
		return err
	}

	data, _ := json.Marshal(struct {
		Key     string  `json:"key"`
		TableId string  `json:"table_id"`
		Data    dataRow `json:"data"`
		Filter  string  `json:"filter"`
	}{
		Key:     c.cf.Key,
		TableId: id,
		Data:    pd.toDataRow(),
		Filter:  fmt.Sprintf(`ud_id="%d"`, pd.Id),
	})

	var resp struct {
		CommonResponse
		Result struct {
			Count int `json:"count"`
		} `json:"result"`
	}

	if err = c.http.PostByte(domain+pathPlaceCloudDataUpdate, data, &resp); err != nil {
		return
	}
	if resp.IsOk() {
		return nil
	}
	return resp.Error()
}

func (c *Lbs) PlaceCloudDataDelete(id int64, tableId ...string) (err error) {
	if id < 1 {
		return errors.New("必须指定删除 id")
	}
	var table string
	if table, err = c.getTableId(tableId...); err != nil {
		return err
	}
	q, _ := json.Marshal(struct {
		Key     string `json:"key"`
		TableId string `json:"table_id"`
		Filter  string `json:"filter"`
	}{
		Key:     c.cf.Key,
		TableId: table,
		Filter:  fmt.Sprintf(`ud_id="%d"`, id),
	})

	var resp struct {
		CommonResponse
		Result struct {
			Count int `json:"count"`
		} `json:"result"`
	}

	if err := c.http.PostByte(domain+pathPlaceCloudDataDelete, q, &resp); err != nil {
		return err
	}

	if resp.IsOk() {
		return nil
	}
	return resp.Error()
}

type ResponsePlaceCloudList struct {
	Count    int        `json:"count"`
	Data     []PlaceRow `json:"data"`
	PageNext string     `json:"page_next"` // 下一页，不传 page_index 时才有
}

func (c *Lbs) PlaceCloudDataList(pag db.Pagination, tableId ...string) (rst ResponsePlaceCloudList, err error) {
	var id string
	if id, err = c.getTableId(tableId...); err != nil {
		return
	}
	params := map[string]interface{}{
		"table_id":  id,
		"orderby":   "id desc",
		"page_size": pag.Limit(),
		"page_next": pag.PageNext,
		"key":       c.cf.Key,
	}
	if pag.PageIndex > 0 {
		params["page_index"] = pag.PageIndex + 1
	}

	var resp struct {
		CommonResponse
		ResponsePlaceCloudList
	}

	if err = c.http.Get(domain+pathPlaceCloudDataList, params, &resp); err != nil {
		return
	}
	if resp.IsOk() {
		return resp.ResponsePlaceCloudList, nil
	}
	err = resp.Error()
	return
}

const (
	pathPlaceCloudSearchNearby = "/place_cloud/search/nearby" // 附近搜索
	pathPlaceCloudSearchRegion = "/place_cloud/search/region" // 区域搜索
)

// https://lbs.qq.com/service/placeCloud/placeCloudGuide/cloudSearch
type RequestPlaceSearch struct {
	Center          datax.Point // 附近搜索
	Region          string      // 行政区域 region=北京市,海淀区 ,region=130681
	Radius          int         // 默认 5000 米，最大 10000米
	Keyword         string      // 关键字
	Pagination db.Pagination
	Kind       int // 类型过滤
	OrderByDistance bool
}

type PlaceRow struct {
	UdId       string      `json:"ud_id"`                  // 即自定义ID（user defined id) ，若您已有库表数据ID或编号，可填入此字段，以便关联管理。
	Title      string      `json:"title"`                  // 地点名称
	Location   datax.Point `json:"location"`               // 坐标
	Polygon    string      `json:"polygon"`                // 多边形轮廓坐标串
	Address    string      `json:"address"`                // 地址
	Tel        string      `json:"tel"`                    // 联系电话
	Id         string      `json:"id"`                     // 主键/标识，数据创建时自动生成，全库不重复
	CreateTime int         `json:"create_time"`            // 创建时间
	UpdateTime int         `json:"update_time"`            // 更新时间
	Province   string      `json:"province"`               // 省
	City       string      `json:"city"`                   // 市
	District   string      `json:"district"`               // 区
	Adcode     int         `json:"adcode"`                 // 行政区划代码
	Distance   int                    `json:"_distance"`   // 距离米
	X          map[string]interface{} `json:"x,omitempty"` // 自定义字段
}
type ResponsePlaceSearch struct {
	Count int64      `json:"count"`
	Data  []PlaceRow `json:"data"`
}

// 周边搜索
func (c *Lbs) PlaceCloudSearch(pd RequestPlaceSearch, tableId ...string) (rst ResponsePlaceSearch, err error) {
	var id string
	if id, err = c.getTableId(tableId...); err != nil {
		return
	}

	if pd.Radius < 1 {
		pd.Radius = 5000
	}
	param := struct {
		Region    string `json:"region,omitempty"`
		Location  string `json:"location"`
		Radius    int    `json:"radius"`
		Keyword   string `json:"keyword,omitempty"`
		PageSize  int    `json:"page_size,omitempty"`
		PageIndex int    `json:"page_index,omitempty"`
		TableId   string `json:"table_id"`
		Key       string `json:"key"`
		Filter    string `json:"filter,omitempty"`
		OrderBy   string `json:"orderby,omitempty"` // 排序
	}{
		Region:    pd.Region,
		Location:  pd.Center.ToQQMapLatLng(),
		Radius:    pd.Radius,
		Keyword:   pd.Keyword,
		PageSize:  pd.Pagination.Limit(),
		PageIndex: pd.Pagination.PageIndex,
		TableId:   id,
		Key:       c.cf.Key,
	}
	if pd.Kind > 0 {
		param.Filter = fmt.Sprintf("x.kind=%d", pd.Kind)
	}
	if pd.OrderByDistance && pd.Center.HasData() {
		param.OrderBy = fmt.Sprintf("distance(%f,%f)", pd.Center.Lat, pd.Center.Lng)
	}

	resp := struct {
		CommonResponse
		Result ResponsePlaceSearch `json:"result"`
	}{}
	url := pathPlaceCloudSearchNearby
	if pd.Region != "" {
		url = pathPlaceCloudSearchRegion
	}

	if err = c.http.Get(domain+url, datax.StructToMap(param).Encode(), &resp); err != nil {
		return
	}
	if resp.IsOk() {
		return resp.Result, nil
	}
	return rst, resp.Error()
}
