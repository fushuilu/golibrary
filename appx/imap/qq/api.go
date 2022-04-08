package qq

import (
	"github.com/fushuilu/golibrary/appx/datax"
)

// https://lbs.qq.com/service/webService/webServiceGuide/webServiceGeocoder

const (
	urlGeocoder = "/ws/geocoder/v1/"
)

type ResponseGeocode struct {
	Title             string      `json:"title"`    // 最终用于坐标解析的地址或地点名称
	Location          datax.Point `json:"location"` // 解析到的坐标（GCJ02坐标系）
	AddressComponents struct {
		Province string `json:"province"` // 省
		City     string `json:"city"`     // 市，如果当前城市为省直辖县级区划，city与district字段均会返回此城市
		// 注：省直辖县级区划adcode第3和第4位分别为9、0，如济源市adcode为419001
		District     string `json:"district"`      // 区，可能为空
		Street       string `json:"street"`        // 街道/道路，可能为空字串
		StreetNumber string `json:"street_number"` // 门牌，可能为空字串
	} `json:"address_components"` // 解析后的地址部件

	AdInfo struct {
		Adcode string `json:"adcode"`
	} `json:"ad_info"`                   // 行政区划信息
	Reliability int `json:"reliability"` // 可信范围 1低可信 ~ 7可信 ~ 10高可信
	Level       int `json:"level"`       // 解析精度级别，分为11个级别，一般>=9即可采用（定位到点，精度较高）
}

//
//  Geocoder
//  @Description: 地址解析（地址转坐标）
//  @receiver c
//
func (c *Lbs) Geocoder(address string) (info ResponseGeocode, err error) {
	var resp struct {
		CommonResponse
		Result ResponseGeocode `json:"result"`
	}
	if err = c.http.Get(domain+urlGeocoder, map[string]string{
		"key":     c.cf.Key,
		"address": address,
	}, &resp); err != nil {
		return
	}
	if resp.IsOk() {
		return resp.Result, nil
	}
	err = resp.Error()
	return
}
