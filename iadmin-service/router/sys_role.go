package router

import (
	v1 "ginProject/api/v1"
	"ginProject/middleware"
	"github.com/gin-gonic/gin"
)

func InitSysRoleRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	RoleRouter := Router.Group("sysRole")
	{
		RoleRouter.GET("getSysRoleList", v1.GetSysRoleList)
		RoleRouter.Use(middleware.OperationRecord()).PUT("updateSysRoleInfo", v1.UpdateSysRoleInfo)
		RoleRouter.Use(middleware.OperationRecord()).POST("addSysRoleInfo", v1.AddSysRoleInfo)
		RoleRouter.Use(middleware.OperationRecord()).DELETE("deleteBatchSysRole", v1.DeleteBatchSysRole)
		RoleRouter.Use(middleware.OperationRecord()).PUT("updateSysRoleMenuConfig", v1.UpdateSysRoleMenuConfig)
		RoleRouter.Use(middleware.OperationRecord()).PUT("setRoleDefaultRouter", v1.SetRoleDefaultRouter)
		RoleRouter.POST("excelOut", v1.ExcelOutSysRole)
	}
	return RoleRouter
}
