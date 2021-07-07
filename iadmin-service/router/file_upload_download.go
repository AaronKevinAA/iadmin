package router

import (
	v1 "ginProject/api/v1"
	"ginProject/middleware"
	"github.com/gin-gonic/gin"
)

func InitFileRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	FileRouter := Router.Group("file")
	{
		FileRouter.Use(middleware.OperationRecord()).POST("upload", v1.UploadFile)
	}
	return FileRouter
}
