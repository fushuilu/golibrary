package golibrary

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"strconv"
	"strings"
)

/*
copy from tencent gin
经常用于处理微信支付
*/

type Params map[string]string

// map本来已经是引用类型了，所以不需要 *Params
func (p Params) SetString(k, s string) Params {
	p[k] = s
	return p
}

func (p Params) GetString(k string) string {
	s, _ := p[k]
	return s
}

func (p Params) SetInt64(k string, i int64) Params {
	p[k] = strconv.FormatInt(i, 10)
	return p
}

func (p Params) GetInt64(k string) int64 {
	i, _ := strconv.ParseInt(p.GetString(k), 10, 64)
	return i
}

// 判断key是否存在
func (p Params) ContainsKey(key string) bool {
	_, ok := p[key]
	return ok
}

func XmlCDATAToMap(xmlData []byte) Params {

	params := make(Params)
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))

	var (
		key   string
		value string
	)

	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement: // 开始标签
			key = token.Name.Local
		case xml.CharData: // 标签内容
			content := string([]byte(token))
			value = content
		}
		if key != "xml" {
			if value != "\n" && strings.TrimRight(value, " ") != "\n" {
				params.SetString(key, value)
			}
		}
	}

	return params
}

func XmlCDATAToStruct(xmlData []byte, data interface{}) error {
	params := XmlCDATAToMap(xmlData)
	if marshal, err := json.Marshal(params); err != nil {
		return err
	} else {
		return json.Unmarshal(marshal, data)
	}
}

type InnerXML struct {
	Value string `xml:",innerxml"`
}

func getInnerXML(v interface{}, indent string) (string, error) {
	b, err := xml.MarshalIndent(v, "", indent)
	if err != nil {
		return "", err
	}
	var in InnerXML
	err = xml.Unmarshal(b, &in)
	if strings.HasSuffix(in.Value, "\n") {
		return in.Value[:len(in.Value)-1], err
	}
	return in.Value, err
}

// 继承
type CDATA string

func (cdata CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	type innerData struct {
		Value string `xml:",cdata"`
	}
	err = e.EncodeElement(innerData{Value: string(cdata)}, start)
	return
}

func XmlHtmlSpecialChars(data interface{}) (string, error) {
	base, err := getInnerXML(data, "  ")
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	en := xml.NewEncoder(buf)
	en.Indent("", "  ")
	start := xml.StartElement{Name: xml.Name{Local: "xml"}}
	if err = en.EncodeElement(base, start); err != nil {
		return "", err
	}
	if err = en.Flush(); err != nil {
		return "", err
	}

	reqBody := append([]byte(xml.Header), buf.Bytes()...)
	return string(reqBody), nil
}
