package golibrary

import (
	"fmt"
	"strconv"
)

func FormatFloat(data float64, maxDecimalLen int) float64 {
	if data == 0 {
		return 0
	}
	tmp := strconv.FormatFloat(data, 'g', maxDecimalLen, 64)
	f, _ := strconv.ParseFloat(tmp, 64)
	return f
}

//
//  FormatPrice
//  @Description: 金额格式化，将 100 分转为 1.00 元， 101 分转为 1.01 元
//
func FormatPrice(price int) string {
	if price < 0 {
		return "NAN"
	} else if price == 0 {
		return "0"
	}
	return fmt.Sprintf("%.2f", FormatFloat(float64(price)/100, 3))
}

//
//  FormatMoney
//  @Description: 金额格式化，将 100 分转为 1.00 元， 101 分转为 1.01 元
//
func FormatMoney(price int) string {
	return FormatPrice(price)
}
