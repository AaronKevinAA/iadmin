package v1

import (
	"ginProject/global"
	"ginProject/model"
	"ginProject/model/response"
	"ginProject/service"
	"ginProject/utils"
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
			response.FailWithMessage(err.Error(), c)
			return
		}
	}
	response.OkWithDetailed(gin.H{"file": file}, "上传成功", c)
}

// @Tags Base
// @Summary 下载批量导入的模板
// @accept application/json
// @Produce application/json
// @Param databaseName query string true "数据库表名称"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"下载成功"}"
// @Router /api/base/downloadExcelInTemplate [GET]
func DownloadExcelInTemplate(c *gin.Context) {
	databaseName := c.Query("databaseName")
	filePath, err := utils.GenerateExcelInTemplate(databaseName)
	if err != nil {
		global.GvaLog.Error("下载失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	c.Writer.Header().Add("Content-Disposition", "excelInTemplateTemp.xlsx")
	c.File(filePath)
	// 删除Excel文件
	err = utils.DeleteLocalFile(filePath)
	if err != nil {
		global.GvaLog.Error("删除文件失败!", zap.Any("err", err))
		return
	}
}
