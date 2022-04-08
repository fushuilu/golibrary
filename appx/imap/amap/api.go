package amap

import (
	"strings"

	"github.com/fushuilu/golibrary"
	"github.com/fushuilu/golibrary/appx/datax"
)

// https://lbs.amap.com/api/webservice/guide/api/georegeo

const (
	urlGeocoder   = "https://restapi.amap.com/v3/geocode/geo"
	urlReGeocoder = "https://restapi.amap.com/v3/geocode/regeo"
)

type ResponseGeocodeResult struct {
	Count    string            `json:"count"`
	Geocodes []ResponseGeocode `json:"geocodes"`
}

func (rg *ResponseGeocodeResult) HasData() bool {
	return golibrary.AnyToInt(rg.Count) > 0
}

func (rg *ResponseGeocodeResult) First() ResponseGeocode {
	return rg.Geocodes[0]
}

type ResponseGeocode struct {
	FormattedAddress string      `json:"formatted_address"` // 结构化地址信息  省份＋城市＋区县＋城镇＋乡村＋街道＋门牌号码
	Country          string      `json:"country"`           // 国家
	Province         string      `json:"province"`          // 省份
	City             string      `json:"city"`              // 北京市
	Citycode         string      `json:"citycode"`          // 城市编码 010
	District         string      `json:"district"`          // 地址所在区
	//Street           interface{} `json:"street"`            // 街道
	//Number           interface{} `json:"number"`            // 门牌
	Adcode           string      `json:"adcode"`            // 区域编码 110101
	Location         string      `json:"location"`          // 经度，纬度 "116.483038,39.990633"
	Level            string      `json:"level"`             // 匹配级别 "门牌号"
}

func (rg ResponseGeocode) ToPoint() datax.Point {
	if pp := strings.Split(rg.Location, ","); len(pp) == 2 {
		return datax.Point{
			Lat: golibrary.AnyToFloat64(pp[0]),
			Lng: golibrary.AnyToFloat64(pp[1]),
		}
	}
	return datax.Point{}
}

type RequestGeocode struct {
	Address string `json:"address"` // 结构化地址信息: 北京市朝阳区阜通东大街6号
	City    string `json:"city"`    // 北京、beijing, 010, 110000 不支持县级市
	Batch   bool   `json:"batch"`   // 多点查询
}

func (c *Lbs) Geocode(pd RequestGeocode) (info ResponseGeocodeResult, err error) {
	var resp struct {
		CommonResponse
		ResponseGeocodeResult
	}

	if err = c.http.Get(urlGeocoder, map[string]interface{}{
		"key":     c.cf.Key,
		"address": pd.Address,
		"city":    pd.City,
		"batch":   pd.Batch,
	}, &resp); err != nil {
		return
	}

	if resp.IsOk() {
		return resp.ResponseGeocodeResult, nil
	}
	err = resp.Error()
	return
}

// 逆地址编码
type RequestReGeocode struct {
	Location datax.Point
}
