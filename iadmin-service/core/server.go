package core

import (
	"fmt"
	"ginProject/global"
	"ginProject/initialize"
	"log"
	"time"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	if global.GvaConfig.System.UseMultipoint {
		// 初始化redis服务
		initialize.Redis()
	}
	Router := initialize.Router()
	//Router.Static("/form-generator", "./resource/page")
	Router.Static("/upload", "./uploads/file")
	address := fmt.Sprintf(":%d", global.GvaConfig.System.Addr)
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.GvaLog.Info(fmt.Sprintf("server run success on http://127.0.0.1%s", address))
	global.GvaLog.Info(fmt.Sprintf("swagger page on http://127.0.0.1%s/swagger/index.html", address))
	log.Printf(s.ListenAndServe().Error())
}
