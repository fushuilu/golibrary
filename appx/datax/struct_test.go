package datax

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type SubjectType int

const (
	_ SubjectType = iota

	SubjectTypeCLPersonal = 210
	SubjectTypeCLCompany  = 220
)

var MapSubjectType = Dict{
	Name: "科目类型",
	Data: map[string]interface{}{
		"210": SubjectTypeCLPersonal,
		"220": SubjectTypeCLCompany,
	},
}

func TestDict_ToText(t *testing.T) {
	text := MapSubjectType.ToText(SubjectTypeCLPersonal)
	assert.Equal(t, "210", text)
}
