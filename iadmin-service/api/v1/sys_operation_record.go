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

// @Tags SysOperationRecord
// @Summary 创建SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysOperationRecord true "创建SysOperationRecord"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysOperationRecord/createSysOperationRecord [post]
func CreateSysOperationRecord(c *gin.Context) {
	var sysOperationRecord model.SysOperationRecord
	_ = c.ShouldBindJSON(&sysOperationRecord)
	if err := service.CreateSysOperationRecord(sysOperationRecord); err != nil {
		global.GvaLog.Error("创建失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// @Tags SysOperationRecord
// @Summary 删除SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysOperationRecord true "SysOperationRecord模型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysOperationRecord/deleteSysOperationRecord [delete]
func DeleteSysOperationRecord(c *gin.Context) {
	var sysOperationRecord model.SysOperationRecord
	_ = c.ShouldBindJSON(&sysOperationRecord)
	if err := service.DeleteSysOperationRecord(sysOperationRecord); err != nil {
		global.GvaLog.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags SysOperationRecord
// @Summary 批量删除SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除SysOperationRecord"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /sysOperationRecord/deleteBatchSysOperationRecord [DELETE]
func DeleteBatchSysOperationRecord(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := service.DeleteBatchSysOperationRecord(IDS); err != nil {
		global.GvaLog.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// @Tags SysOperationRecord
// @Summary 用id查询SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysOperationRecord true "Id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sysOperationRecord/findSysOperationRecord [get]
func FindSysOperationRecord(c *gin.Context) {
	var sysOperationRecord model.SysOperationRecord
	_ = c.ShouldBindQuery(&sysOperationRecord)
	if err := utils.Verify(sysOperationRecord, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, resysOperationRecord := service.GetSysOperationRecord(sysOperationRecord.ID); err != nil {
		global.GvaLog.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(gin.H{"resysOperationRecord": resysOperationRecord}, "查询成功", c)
	}
}

// @Tags SysOperationRecord
// @Summary 分页获取SysOperationRecord列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysOperationRecordSearch true "页码, 每页大小, 搜索条件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysOperationRecord/getSysOperationRecordList [POST]
func GetSysOperationRecordList(c *gin.Context) {
	var pageInfo request.SysOperationRecordSearch
	_ = c.ShouldBindJSON(&pageInfo)
	if err, list, total := service.GetSysOperationRecordInfoList(pageInfo); err != nil {
		global.GvaLog.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Current:  pageInfo.Current,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// @Tags SysOperationRecord
// @Summary 批量导出SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysOperationRecordExcelOut true "批量导出请求模型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量导出成功"}"
// @Router /api/sysOperationRecord/excelOut [POST]
func ExcelOutSysOperationRecord(c *gin.Context) {
	var E request.SysOperationRecordExcelOut
	_ = c.ShouldBindJSON(&E)
	err, excelFilePath := service.ExcelOutSysOperationRecord(E)
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
	c.Writer.Header().Add("Content-Disposition", "sysOperationRecordTable.xlsx")
	c.File(excelFilePath)
	// 删除Excel文件
	err = utils.DeleteLocalFile(excelFilePath)
	if err != nil {
		global.GvaLog.Error("删除文件失败!", zap.Any("err", err))
		return
	}
	return
}
