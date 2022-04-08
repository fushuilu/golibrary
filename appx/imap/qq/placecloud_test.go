package qq

import (
	"fmt"
	"github.com/fushuilu/golibrary/appx/datax"
	"github.com/fushuilu/golibrary/appx/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	key     = ""
	tableId = "0or5oQ1m6LRZuGGja1"
	id      = 1000000000
)

type mockHttp struct {
}

func (m *mockHttp) Post(url string, mapData interface{}, resp interface{}) error {
	panic("implement me")
}

func (m *mockHttp) PostByte(url string, data []byte, resp interface{}) (err error) {
	panic("implement me")
}

func (m *mockHttp) Get(url string, mapData interface{}, resp interface{}) error {
	panic("implement me")
}

func TestQQ_PlaceCloudTableList(t *testing.T) {
	if key == "" {
		return
	}

	lbs := NewQQLbs(Conf{
		Key: key,
	}, &mockHttp{})

	tables, err := lbs.PlaceCloudTableList()
	assert.Nil(t, err)

	fmt.Printf("tables:%+v", tables)
}

func TestQQ_PlaceCloudDataCreate(t *testing.T) {
	if key == "" {
		return
	}

	lbs := NewQQLbs(Conf{
		Key:     key,
		Debug:   true,
		TableId: tableId,
	}, &mockHttp{})

	err := lbs.PlaceCloudDataCreate(RequestPlaceCloudData{
		Id:    id,
		Title: "横溪工业区",
		Location: datax.Point{
			Lat: 23.685475,
			Lng: 116.556922,
		},
		Address: "潮州市古巷镇横溪工业区",
		Tel:     "13420958290",
		X: map[string]interface{}{
			"thumb":   "http://localhost/logo.jpg",
			"summary": "测试简介",
			"remark":  "测试记录备注",
			"kind":    100, // 测试数据
		},
	})
	assert.Nil(t, err)
}

func TestQQ_PlaceCloudDataUpdate(t *testing.T) {
	if key == "" {
		return
	}

	lbs := NewQQLbs(Conf{
		Key:   key,
		Debug: true,
	}, &mockHttp{})

	err := lbs.PlaceCloudDataUpdate(RequestPlaceCloudData{
		Id:    id,
		Title: "横溪新工业区",
		Location: datax.Point{
			Lat: 23.685475,
			Lng: 116.556922,
		},
		Address: "潮州市古巷镇横溪工业区",
		Tel:     "13420959290",
		X: map[string]interface{}{
			"thumb":   "http://localhost/logo2.jpg",
			"summary": "测试简介2",
			"remark":  "测试记录备注2",
		},
	}, tableId)
	assert.Nil(t, err)
}

func TestQQ_PlaceCloudDataDelete(t *testing.T) {
	if key == "" {
		return
	}
	lbs := NewQQLbs(Conf{
		Key: key, Debug: true,
	}, &mockHttp{})

	err := lbs.PlaceCloudDataDelete(1, tableId)
	assert.Nil(t, err)
}
func TestLbs_PlaceCloudDataList(t *testing.T) {
	if key == "" {
		return
	}
	lbs := NewQQLbs(Conf{
		Key:   key,
		Debug: true,
	}, &mockHttp{})
	_, err := lbs.PlaceCloudDataList(db.DefaultPagination, tableId)
	assert.Nil(t, err)
}

func TestQQ_PlaceCloudSearch(t *testing.T) {
	if key == "" {
		return
	}

	lbs := NewQQLbs(Conf{Key: key, Debug: true}, &mockHttp{})
	rst, err := lbs.PlaceCloudSearch(RequestPlaceSearch{
		//Region: "广东省,潮州市,潮安区",
		Center: datax.Point{
			Lat: 23.667331,
			Lng: 116.571802,
		},
		Keyword:         "",
		Pagination:      db.DefaultPagination,
		Kind:            100,
		OrderByDistance: true,
	}, tableId)
	assert.Nil(t, err)
	fmt.Printf("rst:%+v", rst.Count)
}
