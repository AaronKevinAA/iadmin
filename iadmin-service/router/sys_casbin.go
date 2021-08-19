package router

import (
	v1 "ginProject/api/v1"
	"ginProject/middleware"
	"github.com/gin-gonic/gin"
)

func InitCasbinRouter(Router *gin.RouterGroup) {
	CasbinRouter := Router.Group("casbin")
	{
		CasbinRouter.Use(middleware.OperationRecord()).POST("updateCasbin", v1.UpdateCasbin)
		CasbinRouter.POST("getPolicyPathByRoleId", v1.GetPolicyPathByRoleId)
	}
}
