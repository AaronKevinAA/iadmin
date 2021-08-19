package router

import (
	v1 "ginProject/api/v1"
	"ginProject/middleware"
	"github.com/gin-gonic/gin"
)

func InitSysMenuRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	MenuRouter := Router.Group("sysMenu")
	{
		MenuRouter.POST("getSysRouteList", v1.GetSysRouteList)
		MenuRouter.GET("getSysMenuByToken", v1.GetSysMenuByToken)
		MenuRouter.Use(middleware.OperationRecord()).PUT("updateSysMenuInfo", v1.UpdateSysMenuInfo)
		MenuRouter.Use(middleware.OperationRecord()).POST("addSysMenuInfo", v1.AddSysMenuInfo)
		MenuRouter.Use(middleware.OperationRecord()).DELETE("deleteBatchSysMenu", v1.DeleteBatchSysMenu)
		MenuRouter.POST("excelOut", v1.ExcelOutSysMenu)
	}
	return MenuRouter
}
