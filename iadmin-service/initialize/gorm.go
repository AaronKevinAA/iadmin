package initialize

import (
	"ginProject/global"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
)

func GormMysql() *gorm.DB {
	dsn := global.GvaConfig.Mysql.Dsn()
	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	}); err != nil {
		global.GvaLog.Error("mysql 连接异常, err:",zap.Any("err",err))
		os.Exit(0)
		return nil
	} else {
		global.GvaLog.Info("mysql 连接成功",zap.String("dsn",dsn))
		return db
	}
}