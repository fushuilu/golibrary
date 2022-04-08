package datax

import (
	"fmt"
	"testing"
)

type direction = int

const (
	_    direction = iota
	east           = 1 << 0
	south           = 1 << 1
	west            = 1 << 2
	north           = 1 << 3
)

func TestEnum(t *testing.T) {
	fmt.Println("east:", east, ";south:", south, ";west:", west, ";north:", north)
	d := east + west

	res1 := d & east  // 1
	res2 := d & north // 0
	fmt.Println(res1, res2)
}
