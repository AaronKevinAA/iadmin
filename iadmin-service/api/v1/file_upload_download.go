package v1

import (
	"ginProject/global"
	"ginProject/model"
	"ginProject/model/response"
	"ginProject/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Tags File
// @Summary 上传文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "上传文件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"上传成功"}"
// @Router /api/file/upload [post]
func UploadFile(c *gin.Context) {
	var file model.File
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GvaLog.Error("接收文件失败!", zap.Any("err", err))
		response.FailWithMessage("接收文件失败！", c)
		return
	}
	if header != nil {
		file, err = service.UploadFile(header)
		if err != nil {
			global.GvaLog.Error("上传文件失败!", zap.Any("err", err))
			response.FailWithMessage("上传文件失败！", c)
			return
		}
	}
	response.OkWithDetailed(gin.H{"file": file}, "上传成功", c)

}
