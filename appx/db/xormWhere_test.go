package db

import (
	"fmt"
	"testing"
)

func Test_join(t *testing.T) {
	keys, vs := join("1,2,3", []string{"a", "b", "c"}, "%s IN (?)")
	fmt.Println("keys:", keys)
	fmt.Println("vs", vs)
}
