package golibrary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXml(t *testing.T) {
	// 注意，换行时前面没有空格
	closeData := `<xml>
<return_code><![CDATA[SUCCESS]]></return_code>
<return_msg><![CDATA[OK]]></return_msg>
<appid><![CDATA[wx2421b1c4370ec43b]]></appid>
<mch_id><![CDATA[10000100]]></mch_id>
<nonce_str><![CDATA[BFK89FC6rxKCOjLX]]></nonce_str>
<sign><![CDATA[72B321D92A7BFA0B2509F3D13C7B1631]]></sign>
<result_code><![CDATA[SUCCESS]]></result_code>
<result_msg><![CDATA[OK]]></result_msg>
</xml>`
	params := XmlCDATAToMap([]byte(closeData))
	assert.Equal(t, "wx2421b1c4370ec43b", params.GetString("appid"))
	assert.Equal(t, "10000100", params.GetString("mch_id"))
	assert.Equal(t, "BFK89FC6rxKCOjLX", params.GetString("nonce_str"))
	assert.Equal(t, "72B321D92A7BFA0B2509F3D13C7B1631", params.GetString("sign"))


	type closeorder struct {
		ReturnCode string `json:"return_code"`
		ReturnMsg  string `json:"return_msg"`
		Appid      string `json:"appid"`
		MchId      string `json:"mch_id"`
		NonceStr   string `json:"nonce_str"`
		Sign       string `json:"sign"`
		ResultCode string `json:"result_code"`
		ResultMsg  string `json:"result_msg"`
	}

	coData := closeorder{}
	err := XmlCDATAToStruct([]byte(closeData), &coData)
	assert.Nil(t, err)
	assert.Equal(t, "wx2421b1c4370ec43b", coData.Appid)
	assert.Equal(t, "10000100", coData.MchId)
	assert.Equal(t, "BFK89FC6rxKCOjLX", coData.NonceStr)
	assert.Equal(t, "72B321D92A7BFA0B2509F3D13C7B1631", coData.Sign)

}

type payNotifyRspOut struct {
	ReturnCode string `xml:"return_code"` // 返回状态码，SUCCESS表示商户接收通知成功并校验成功
	ReturnMsg  string `xml:"return_msg"`  // 返回信息
}

func TestXmlHtmlSpecialChars(t *testing.T) {

	rspOut := payNotifyRspOut{
		ReturnCode: "SUCCESS",
		ReturnMsg:  "OK",
	}
	_, err := XmlHtmlSpecialChars(rspOut)
	assert.Nil(t, err)
}
