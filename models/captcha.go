package models

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
)

// 创建store
// var store = base64Captcha.DefaultMemStore
// 配置RedisStore  RedisStore实现base64Captcha.Store接口
var store = base64Captcha.DefaultMemStore

// MackCaptcha
// @Description 生成验证码
// @Author xYuan 2024-04-15 12:31:34
// @Return string
// @Return string
// @Return string
// @Return error
func MackCaptcha(height int, width int, length int) (string, string, string, error) {
	var driver base64Captcha.Driver
	driverString := &base64Captcha.DriverString{
		Height:          height,
		Width:           width,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          length,
		Source:          "1234567890qwertyuiopasdfghjklzxcvbnm",
		BgColor:         &color.RGBA{R: 3, G: 102, B: 214, A: 125},
		Fonts:           []string{"wqy-microhei.ttc"},
	}
	driver = driverString.ConvertFonts()
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, answer, err := c.Generate()
	return id, b64s, answer, err
}

// VerifyCaptcha
// @Description 校验验证码
// @Author xYuan 2024-04-15 12:31:21
// @Param id
// @Param VerifyValue
// @Return bool
func VerifyCaptcha(id string, VerifyValue string) bool {
	if store.Verify(id, VerifyValue, true) {
		return true
	} else {
		return false
	}
}
