package amap

import (
	"errors"
	"github.com/fushuilu/golibrary"
)

type Lbs struct {
	cf   Conf
	http golibrary.ClientHttp
}

type Conf struct {
	Key   string
	Debug bool
}

func NewAMapLbs(cf Conf, http golibrary.ClientHttp) Lbs {
	return Lbs{
		cf:   cf,
		http: http,
	}
}

type CommonResponse struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}

func (r CommonResponse) IsOk() bool {
	return r.Status == "1"
}

func (r CommonResponse) Error() error {
	return errors.New(r.Info)
}
