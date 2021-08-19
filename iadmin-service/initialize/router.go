package initialize

import (
	"ginProject/global"
	"ginProject/middleware"
	"ginProject/router"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// 初始化总路由

func Router() *gin.Engine {
	var Router = gin.Default()
	//Router.StaticFS(global.GVA_CONFIG.Local.Path, http.Dir(global.GVA_CONFIG.Local.Path)) // 为用户头像和文件提供静态地址

	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.GvaLog.Info("swagger注册成功,swagger页面为：/swagger/index.html")
	// 配置跨域请求
	Router.Use(middleware.Cors())
	// 方便统一添加路由组前缀 多服务器上线使用
	PublicGroup := Router.Group("/api")
	{
		router.InitBaseRouter(PublicGroup)

	}
	PrivateGroup := Router.Group("/api")
	PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		router.InitSysUserRouter(PrivateGroup)
		router.InitSysMenuRouter(PrivateGroup)
		router.InitSysRoleRouter(PrivateGroup)
		router.InitSysApiRouter(PrivateGroup)
		router.InitSysOperationRecordRouter(PrivateGroup) // 操作记录
		router.InitCasbinRouter(PrivateGroup)             // 权限相关路由
		router.InitFileRouter(PrivateGroup)

	}
	global.GvaLog.Info("router 注册成功")
	return Router
}
