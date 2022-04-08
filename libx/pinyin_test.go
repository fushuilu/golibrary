package cmn

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFirstPinYin(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "多字", args: args{name: "中国人"}, want: "z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FirstPinYin(tt.args.name); got != tt.want {
				t.Errorf("FirstPinYin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPinYin(t *testing.T) {
	word := "中国人"
	assert.Equal(t, "zgr", PinYin(word))
}
