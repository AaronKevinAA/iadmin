package v1

import (
	"ginProject/global"
	"ginProject/model"
	"ginProject/model/request"
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

// @Tags File
// @Summary 分页获取文件列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.Pagination true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取数据成功"}"
// @Router /api/file/getSysFileList [POST]
func GetSysFileList(c *gin.Context) {
	var U request.Pagination
	_ = c.ShouldBindJSON(&U)
	err, total, list := service.GetSysFileList(&U)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Current:  U.Current,
			PageSize: U.PageSize,
		}, "获取数据成功", c)
	}
}

// @Tags File
// @Summary 下载文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param id query uint true "文件数据库表id"
// @Success 200
// @Router /api/file/download [GET]
func DownloadFile(c *gin.Context) {
	fileID := c.Query("id")
	fileName, filePath, err := service.DownloadFile(global.String2uint(fileID))
	ok, err := utils.PathExists(filePath)
	if !ok || err != nil {
		global.GvaLog.Error("文件不存在!", zap.Any("err", err))
		response.FailWithMessage("文件不存在!", c)
		return
	}
	c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	c.Writer.Header().Add("Content-Disposition", fileName)
	c.File(filePath)
}

// @Tags File
// @Summary 批量删除文件
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "id列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /api/file/deleteBatchFile [DELETE]
func DeleteBatchFile(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := service.DeleteBatchFile(IDS); err != nil {
		global.GvaLog.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}
