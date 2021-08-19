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

// @Tags SysApi
// @Summary 分页获取接口列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysApiListSearch true "页码, 每页大小，搜索条件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取数据成功"}"
// @Router /api/sysApi/getSysApiList [POST]
func GetSysApiList(c *gin.Context) {
	var U request.SysApiListSearch
	_ = c.ShouldBindJSON(&U)
	err, total, list := service.GetSysApiList(&U)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Current:  U.Pagination.Current,
			PageSize: U.Pagination.PageSize,
		}, "获取数据成功", c)
	}
}

// @Tags SysApi
// @Summary 更新接口信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysApi true "接口模型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /api/sysApi/updateSysApiInfo [PUT]
func UpdateSysApiInfo(c *gin.Context) {
	var Api model.SysApi
	_ = c.ShouldBindJSON(&Api)
	if err := service.UpdateSysApiInfo(Api); err != nil {
		global.GvaLog.Error("更新失败！", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags SysApi
// @Summary 新增接口
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysApi true "接口模型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"新增成功"}"
// @Router /api/sysApi/addSysApiInfo [POST]
func AddSysApiInfo(c *gin.Context) {
	var Api model.SysApi
	_ = c.ShouldBindJSON(&Api)
	if err, ret := service.AddSysApiInfo(Api); err != nil {
		global.GvaLog.Error("新增失败！", zap.Any("err", err))
		response.FailWithDetailed(response.SysApiResponse{Api: ret}, err.Error(), c)
	} else {
		response.OkWithDetailed(gin.H{"ApiInfo": ret}, "新增成功", c)
	}
}

// @Tags SysApi
// @Summary 批量删除接口
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "id列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /api/sysApi/deleteBatchSysApi [DELETE]
func DeleteBatchSysApi(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := service.DeleteBatchSysApi(IDS); err != nil {
		global.GvaLog.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// @Tags SysApi
// @Summary 获得接口树
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取数据成功"}"
// @Router /api/sysApi/getSysApiTree [GET]
func GetSysApiTree(c *gin.Context) {
	if err, data := service.GetSysApiTree(); err != nil {
		global.GvaLog.Error("获取数据失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(data, "获取数据成功", c)
	}
}

// @Tags SysApi
// @Summary 批量导出接口
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysApiExcelOut true "批量导出请求模型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量导出成功"}"
// @Router /api/sysApi/excelOut [POST]
func ExcelOutSysApi(c *gin.Context) {
	var E request.SysApiExcelOut
	_ = c.ShouldBindJSON(&E)
	err, excelFilePath := service.ExcelOutSysApi(E)
	if err != nil {
		global.GvaLog.Error("批量导出失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		// 如果文件路径不为空，则页需要删除文件
		if excelFilePath != "" {
			err = utils.DeleteLocalFile(excelFilePath)
			if err != nil {
				global.GvaLog.Error("删除文件失败!", zap.Any("err", err))
				return
			}
		}
		return
	}
	// 要设置这一条，要不然前端获取不到Content-Disposition
	c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	c.Writer.Header().Add("Content-Disposition", "apiTable.xlsx")
	c.File(excelFilePath)
	// 删除Excel文件
	err = utils.DeleteLocalFile(excelFilePath)
	if err != nil {
		global.GvaLog.Error("删除文件失败!", zap.Any("err", err))
		return
	}
	return
}

// @Tags SysApi
// @Summary 批量导入接口预览
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param file formData file true "上传文件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量导入预览成功"}"
// @Router /api/sysApi/excelInPreview [POST]
func ExcelInPreviewSysApi(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GvaLog.Error("接收文件失败!", zap.Any("err", err))
		response.FailWithMessage("接收文件失败！", c)
		return
	}
	if header != nil {
		errPreview, newFileName, dataList, allDataCorrect := service.ExcelInPreviewSysApi(header)
		if errPreview != nil {
			global.GvaLog.Error("批量导入预览失败!", zap.Any("err", errPreview))
			response.FailWithMessage(errPreview.Error(), c)
			return
		}
		response.OkWithDetailed(gin.H{"saveFileName": newFileName, "dataList": dataList, "allDataCorrect": allDataCorrect}, "批量导入预览成功", c)
		return
	}
	response.FailWithMessage("批量导入预览失败！", c)
}

// @Tags SysApi
// @Summary 批量导入接口
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.ExcelInRequest true "保存的文件名称"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量导入成功"}"
// @Router /api/sysApi/excelIn [POST]
func ExcelInSysApi(c *gin.Context) {
	var E request.ExcelInRequest
	_ = c.ShouldBindJSON(&E)
	err := service.ExcelInSysApi(E)
	if err != nil {
		global.GvaLog.Error("批量导入失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("批量导入成功", c)
}
