package router

import (
	v1 "ginProject/api/v1"
	"ginProject/middleware"
	"github.com/gin-gonic/gin"
)

func InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", v1.Captcha)
		BaseRouter.Use(middleware.OperationRecord()).POST("login", v1.Login)
		BaseRouter.Use(middleware.OperationRecord()).POST("register", v1.Register)
		BaseRouter.GET("downloadExcelInTemplate", v1.DownloadExcelInTemplate)
		BaseRouter.GET("generateRSAKey", v1.GenerateRSAKey)
	}
	return BaseRouter
}
