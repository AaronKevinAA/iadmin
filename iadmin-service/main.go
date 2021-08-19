package main

import (
	"ginProject/core"
	_ "ginProject/docs"
	"ginProject/global"
	"ginProject/initialize"
	_ "github.com/go-sql-driver/mysql"
)

// @title Swagger Example API
// @version 0.0.1
// @description This is a sample Server pets
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-token
// @BasePath
func main() {
	// 读取配置文件
	global.GvaVp = core.Viper()
	// 初始化zap日志库
	global.GvaLog = core.Zap()
	// 初始化gorm
	global.GvaDb = initialize.GormMysql()

	core.RunWindowsServer()
}
