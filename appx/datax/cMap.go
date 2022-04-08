package datax

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// https://github.com/dreamlu/gt/blob/2c0230b701fdbec64113111b4c651d0f6e41b040/tool/type/cmap/cmap.go#L107
// CMap = url.Values{} = map[string][]string
// struct to CMap, maybe use Encode

type CMap map[string][]string

func StructToMap(v interface{}) (values CMap) {
	values = CMap{}
	el := reflect.ValueOf(v)
	if el.Kind() == reflect.Ptr {
		el = el.Elem()
	}
	iVal := el
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		fi := typ.Field(i)
		name := fi.Tag.Get("json")

		v := fmt.Sprint(iVal.Field(i))
		if strings.Index(name, "omitempty") > -1 {
			name = strings.ReplaceAll(strings.ReplaceAll(name, "omitempty", ""), ",", "")
			name = strings.TrimSpace(name)

			val := strings.ToLower(strings.TrimSpace(v))
			//fmt.Printf("value:%s:%s|\n", name, val)
			switch val {
			case "0", "", "undefined", "null":
				continue
			}
		}
		if name == "" {
			name = fi.Name
		}
		values.Set(name, v)
	}
	return
}

// Set sets the key to value. It replaces any existing
// values.
func (v CMap) Set(key, value string) CMap {
	v[key] = []string{value}
	return v
}

// Encode encodes the values into ``URL encoded'' form
// ("bar=baz&foo=quux") sorted by key.
func (v CMap) Encode() string {
	return url.Values(v).Encode()
}

