package cmn

import (
	"errors"
	"github.com/dchest/captcha"
	"io"
	"time"
)

// https://github.com/dchest/captcha
type Captcha struct {
}

func NewCaptcha() Captcha {
	captcha.SetCustomStore(captcha.NewMemoryStore(100, 2*time.Minute))
	return Captcha{}
}
func (c *Captcha) CaptchaId(id string) string {
	if id == "" {
		return captcha.New()
	} else {
		captcha.Reload(id)
		return id
	}
}

// 生成一个图片验证码
/*
func (c *Captcha) Image(r *ghttp.Request) interface{} {
	width := r.GetQueryInt("width", 240)
	height := r.GetQueryInt("height", 80)
	id := r.GetString("id", "")
	return captcha.WriteImage(r.Response.Writer, id, width, height)
}
*/

func (c *Captcha) WriteImage(w io.Writer, id string, width, height int) error {
	return captcha.WriteImage(w, id, width, height)
}

func (c *Captcha) Verify(id, code string) error {
	if id == "" {
		return errors.New("captcha id is empty")
	}
	if code == "" {
		return errors.New("captcha code is empty")
	}
	if captcha.VerifyString(id, code) {
		return nil
	} else {
		return errors.New("验证码错误")
	}
}
