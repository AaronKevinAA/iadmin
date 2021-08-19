package v1

import (
	"ginProject/global"
	"ginProject/model/response"
	"ginProject/utils"
	"github.com/gin-gonic/gin"
	"time"
)

// @Tags Base
// @Summary 生成RSA密钥
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"生成密钥成功"}"
// @Router /api/base/generateRSAKey [get]
func GenerateRSAKey(c *gin.Context) {
	// 生成写入redis的key字段
	UnixNano := time.Now().UnixNano()
	redisKey := "RSAKey" + global.Int642string(UnixNano)
	// 生成随机公钥和私钥
	err, X509PrivateKey, X509PublicKey := utils.GenerateRSAKey()
	if err != nil {
		response.FailWithMessage("生成密钥失败！", c)
		return
	}
	// 将私钥写入redis
	global.GvaRedis.Append(redisKey, string(X509PrivateKey))
	// 将redis的Key返回给前端
	response.OkWithDetailed(gin.H{"publicKey": X509PublicKey, "redisKey": redisKey}, "生成密钥成功", c)
	return
}
