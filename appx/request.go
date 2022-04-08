package appx

import "errors"

type ParamInvalid interface {
	Invalid(created bool) (err error)
}

// 获取参数使用
func InvalidId(isCreated bool, id int64) error {
	if isCreated {
		if id > 0 {
			return errors.New("添加记录时 ID 必须为空")
		}
	} else {
		if id < 1 {
			return errors.New("更新记录时 ID 不能为空")
		}
	}
	return nil
}