package router

import (
	v1 "ginProject/api/v1"
	"ginProject/middleware"
	"github.com/gin-gonic/gin"
)

func InitSysUserRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	UserRouter := Router.Group("sysUser")
	{
		UserRouter.POST("getSysUserList", v1.GetSysUserList)
		UserRouter.Use(middleware.OperationRecord()).PUT("updateSysUserInfo", v1.UpdateSysUserInfo)
		UserRouter.Use(middleware.OperationRecord()).POST("addSysUserInfo", v1.AddSysUserInfo)
		UserRouter.Use(middleware.OperationRecord()).DELETE("deleteBatchSysUser", v1.DeleteBatchSysUser)
		UserRouter.POST("excelOut", v1.ExcelOutSysUser)
		UserRouter.POST("excelInPreview", v1.ExcelInPreviewSysUser)
		UserRouter.Use(middleware.OperationRecord()).POST("excelIn", v1.ExcelInSysUser)
		UserRouter.Use(middleware.OperationRecord()).PUT("resetPassword", v1.ResetSysUserPassword)
		UserRouter.Use(middleware.OperationRecord()).PUT("updatePasswordByToken", v1.UpdatePasswordByToken)
		UserRouter.Use(middleware.OperationRecord()).PUT("updateBasicInfoByToken", v1.UpdateBasicInfoByToken)
	}
	return UserRouter
}
