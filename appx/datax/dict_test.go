package datax

import "testing"

type DemoKind int

const (
	_ DemoKind = iota

	Company
	Personal
)

var MapDemoKind = Dict{
	Name: "Demo",
	Data: map[string]interface{}{
		"company":  Company,
		"personal": Personal,
	},
}

func TestDict(t *testing.T) {

}
