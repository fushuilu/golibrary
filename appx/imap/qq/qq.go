package qq

import (
	"errors"
	"github.com/fushuilu/golibrary"
)

const (
	domain = "https://apis.map.qq.com"
)

type Lbs struct {
	cf   Conf
	http golibrary.ClientHttp
}

type Conf struct {
	Key     string `json:"key"`
	TableId string `json:"table_id"`
	Debug   bool   `json:"debug"`
}

func NewQQLbs(cf Conf, http golibrary.ClientHttp) Lbs {
	return Lbs{
		cf:   cf,
		http: http,
	}
}

type CommonResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (r CommonResponse) IsOk() bool {
	return r.Status == 0
}
func (r CommonResponse) Error() error {
	return errors.New(r.Message)
}
