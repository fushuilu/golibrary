package datax

import (
	"fmt"
	"testing"
)

func TestStructToMap(t *testing.T) {
	data := struct {
		Name string `json:"name"`
		Age  int    `json:"age,omitempty"`
	}{
		Name: "aa", Age: 0,
	}

	rst := StructToMap(data).Encode()
	fmt.Printf("rst:%+v", rst)
}
