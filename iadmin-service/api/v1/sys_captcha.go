package v1

import (
	"ginProject/global"
	"ginProject/model/response"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"log"
)
var store = base64Captcha.DefaultMemStore

// @Tags Base
// @Summary 生成验证码
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"验证码获取成功"}"
// @Router /api/base/captcha [get]
func Captcha(c *gin.Context)  {
	driver := base64Captcha.NewDriverDigit(global.GvaConfig.Captcha.ImgHeight,global.GvaConfig.Captcha.ImgWidth,global.GvaConfig.Captcha.KeyLong,0.7,80)
	cp := base64Captcha.NewCaptcha(driver,store)
	if id, b64s, err := cp.Generate(); err != nil{
		log.Println("验证码获取失败！")
		response.FailWithMessage("验证码获取失败！",c)
	}else{
		response.OkWithDetailed(response.SysCaptchaResponse{
			CaptchaId: id,
			PicPath: b64s,
		},"验证码获取成功",c)
	}
}