package router

import (
	v1 "ginProject/api/v1"
	"ginProject/middleware"
	"github.com/gin-gonic/gin"
)

func InitSysOperationRecordRouter(Router *gin.RouterGroup) {
	OperationRecordRouter := Router.Group("sysOperationRecord")
	{
		OperationRecordRouter.POST("createSysOperationRecord", v1.CreateSysOperationRecord)                                               // 新建SysOperationRecord
		OperationRecordRouter.DELETE("deleteSysOperationRecord", v1.DeleteSysOperationRecord)                                             // 删除SysOperationRecord
		OperationRecordRouter.Use(middleware.OperationRecord()).DELETE("deleteBatchSysOperationRecord", v1.DeleteBatchSysOperationRecord) // 批量删除SysOperationRecord
		OperationRecordRouter.GET("findSysOperationRecord", v1.FindSysOperationRecord)                                                    // 根据ID获取SysOperationRecord
		OperationRecordRouter.POST("getSysOperationRecordList", v1.GetSysOperationRecordList)                                             // 获取SysOperationRecord列表
		OperationRecordRouter.POST("excelOut", v1.ExcelOutSysOperationRecord)                                                             // 获取SysOperationRecord列表
	}
}
