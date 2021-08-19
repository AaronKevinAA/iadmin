package router

import (
	v1 "ginProject/api/v1"
	"ginProject/middleware"
	"github.com/gin-gonic/gin"
)

func InitSysApiRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	ApiRouter := Router.Group("sysApi")
	{
		ApiRouter.POST("getSysApiList", v1.GetSysApiList)
		ApiRouter.GET("getSysApiTree", v1.GetSysApiTree)
		ApiRouter.Use(middleware.OperationRecord()).PUT("updateSysApiInfo", v1.UpdateSysApiInfo)
		ApiRouter.Use(middleware.OperationRecord()).POST("addSysApiInfo", v1.AddSysApiInfo)
		ApiRouter.Use(middleware.OperationRecord()).DELETE("deleteBatchSysApi", v1.DeleteBatchSysApi)
		ApiRouter.POST("excelOut", v1.ExcelOutSysApi)
		ApiRouter.POST("excelInPreview", v1.ExcelInPreviewSysApi)
		ApiRouter.Use(middleware.OperationRecord()).POST("excelIn", v1.ExcelInSysApi)

	}
	return ApiRouter
}
