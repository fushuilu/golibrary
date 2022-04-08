package datax

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fushuilu/golibrary"
)

type Address struct {
	Province     string  `json:"province" xorm:"province notnull default('')"`           // 省份
	ProvinceCode string  `json:"province_code" xorm:"province_code notnull default('')"` // 省份区码
	City         string  `json:"city" xorm:"city notnull default('')"`                   // 城市
	CityCode     string  `json:"city_code" xorm:"city_code notnull default('')"`         // 城市区码
	County       string  `json:"county" xorm:"county notnull default('')"`               // 区县
	CountyCode   string  `json:"county_code" xorm:"county_code notnull default('')"`     // 区县区码 adcode
	Detail       string  `json:"detail" xorm:"detail notnull default('')"`               // 详细地址
	Lat          float64 `json:"lat" xorm:"notnull default(0)"`                          // 纬度 latitude 0~90
	Lng          float64 `json:"lng" xorm:"notnull default(0)"`                          // 经度 longitude 0~180
	Distance     int     `json:"distance" xorm:"distance notnull default(0)"`            // 只读：计算出来的距离(冗余的字段)
}

func AddressCols() []string {
	return []string{"province", "province_code", "city", "city_code", "county", "county_code", "detail", "lat", "lng"}
}

func (pd *Address) Invalid() error {
	if pd.Province == "" || pd.City == "" || pd.County == "" {
		return errors.New("请填写地址")
	} else if pd.Detail == "" {
		return errors.New("请填写详细地址")
	}
	return nil
}

func (pd *Address) JoinDesc() string {
	return strings.Join([]string{
		pd.Province, pd.City, pd.County, pd.Detail,
	}, "/")
}

func (pd *Address) Desc() string {
	return strings.Join([]string{
		pd.Province, pd.City, pd.County, pd.Detail,
	}, "")
}

func (pd *Address) HasPoint() bool {
	return !(pd.Lat == 0 && pd.Lng == 0)
}

type Point struct {
	Lat float64 `json:"lat" xorm:"notnull default(0)"` // 纬度 latitude -90~90
	Lng float64 `json:"lng" xorm:"notnull default(0)"` // 经度 longitude -180~180
}

func (p *Point) Invalid() error {
	if p.Lat < -90 || p.Lat > 90 {
		return errors.New("纬度 lat 范围 -90 ~ 90")
	} else if p.Lat < -180 || p.Lat > 180 {
		return errors.New("经度 lng 范围 -180 ~ 180")
	}
	return nil
}

// 纬经度，适用于腾讯地址
func (p *Point) ToQQMapLatLng() string {
	if p.HasData() {
		return fmt.Sprintf("%f,%f", p.Lat, p.Lng)
	}
	return ""
}

// 经纬度，适用于高德地址
func (p *Point) ToAMapLngLat() string {
	if p.HasData() {
		return fmt.Sprintf("%f,%f", p.Lng, p.Lat)
	}
	return ""
}

func (p *Point) HasData() bool {
	return !(p.Lat == 0 && p.Lng == 0)
}

func (p *Point) FromLatLng(points string) {
	if points != "" {
		if ds := strings.Split(points, ","); len(ds) == 2 {
			p.Lat = golibrary.AnyToFloat64(ds[0])
			p.Lng = golibrary.AnyToFloat64(ds[1])
		}
	}
}

//// 将坐标转为字符串标记
//func GeoHash(lat, lng float64) string {
//	return geohash.EncodeWithPrecision(lat, lng, geoChars2000Meters)
//}
//
//// 计算 hash 坐标周围的8个区域，形成8个区域，在 db 中使用此 9 个区域查询
//func GeoHashes(pHash string) []string {
//	cs := geohash.Neighbors(pHash)
//	cs = append(cs, pHash)
//	return cs
//}
//
//const geoBits1000Meters = 30 // 30 位是 1000 米; 32 位是 500米
//const geoChars2000Meters = 5 // 5 位是 2000 米; 6 位是 600 米
//
//// 将坐标转为整型标记
//func GeoHashInt64(lat, lng float64) int64 {
//	return int64(geohash.EncodeIntWithPrecision(lat, lng, geoBits1000Meters))
//}
//
//// 计算 hash 坐标周围的 8 个区，共9个区域，在 db 中使用此 9 个区域查询
//func GeoHashesInt64(pHash int64) []int64 {
//	cs := geohash.NeighborsIntWithPrecision(uint64(pHash), geoBits1000Meters)
//	rst := make([]int64, len(cs)+1)
//	for i := range cs {
//		rst[i] = int64(cs[i])
//	}
//	rst = append(rst, pHash)
//	return rst
//}

/*
// 实现driver.Valuer接口
func (p *Point) Value() (driver.Value, error) {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "(%f,%f)", p.Lat, p.Lng)
	return buf.Bytes(), nil
}

func (p *Point) String() string {
	return fmt.Sprintf("(%v,%v)", p.Lat, p.Lng)
}

// 实现sql.Scanner接口
func (p *Point) Scan(val interface{}) (err error) {
	if ll, ok := val.([]uint8); ok {
		tmp := ll[1 : len(ll)-1]
		coors := strings.Split(string(tmp[:]), ",")
		if p.Lat, err = strconv.ParseFloat(coors[0], 64); err != nil {
			return err
		}
		if p.Lng, err = strconv.ParseFloat(coors[1], 64); err != nil {
			return err
		}
	}
	return nil
}

// 实现 xorm 接口
func (p *Point) FromDB(bytes []byte) error {
	return p.Scan(bytes)
}

func (p *Point) ToDB(bytes []byte, err error) {
	bytes = []byte(p.String())
}
*/
