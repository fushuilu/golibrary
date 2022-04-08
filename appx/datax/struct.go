package datax

import (
	"time"
)

type KV struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
	Value string `json:"value"`
}

// 用于分页输出
type ListResult struct {
	Total int64       `json:"total"`
	Tag   int64       `json:"tag"`
	Rows  interface{} `json:"rows"`
}

// 别名
type RowResult ListResult
type SearchResult ListResult

var EmptyArray = make([]int64, 0)
var EmptyTime = time.Time{}
var EmptyMap = map[string]interface{}{}
var ListResultEmpty = ListResult{Total: 0, Tag: 0, Rows: EmptyArray}

type CodeResponse struct {
	Code    int         `json:"code"`    // 响应码，0 或 200 表示成功
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 响应数据
}

func (rd CodeResponse) IsOK() bool {
	return rd.Code == 0 || rd.Code == 200
}

type Location307Response struct {
	Location                 string `json:"location"`
	AccessControlAllowOrigin string `json:"access_control_allow_origin"`
}

// url 响应
type OauthResp struct {
	Url string `json:"url"`
}

// 通常在控制器中使用， 用于代替 if err = xxxx; err != nil { return err }
// ResponseIfError(r, err)
type DoResponse func(err error)

type ChangeState struct {
	Id       int64  `json:"id" form:"id"`               // 主键 id
	Sid      string `json:"sid" form:"sid"`             // 字符串 id
	State    string `json:"state" form:"state"`         // 审核状态
	StateMsg string `json:"state_msg" form:"state_msg"` // 审核信息

	StateIndex int `json:"-" form:"-"`
}

type ChangeStatus struct {
	Id     int64  `json:"id" form:"id"`
	Sid    string `json:"sid" form:"sid"`
	Status string `json:"status" form:"status"`
	UserId int64  `json:"-" form:"-"`

	StatusIndex int `json:"-" form:"-"`
}
