package global

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type GvaModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt int64          `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt int64          `json:"updated_at" gorm:"autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Timestamp2DateTime 写在utils会有导包循环的问题，先写在这里解决问题
// Timestamp2DateTime 时间戳转日期格式
func Timestamp2DateTime(timestamp int64) (datetime string) {
	if timestamp == 0 {
		return ""
	}
	// 时区
	Loc, _ := time.LoadLocation("Asia/Shanghai")
	// 13位时间戳需要/1000
	timeObj := time.Unix(int64(timestamp/1000), 0).In(Loc)
	// 2006 表示年份，好像是 Go 开始设计的年份，而后面 1 2 3 4 5 6 7，通过简单的顺序就可以记忆
	// 无语
	datetime = timeObj.Format("2006-01-02 15:04:05")
	return datetime
}

// Uint2String uint 转 string
func Uint2String(number uint) string {
	return strconv.Itoa(int(number))
}

// Int2String int 转 string
func Int2String(number int) string {
	return strconv.Itoa(number)
}

// Bool2String bool 转 string
func Bool2String(flag *bool) string {
	if *flag {
		return "是"
	} else {
		return "否"
	}
}

// Int642string int64 转 string
func Int642string(number int64) string {
	return strconv.FormatInt(number, 10)
}

// String2uint string 转 uint
func String2uint(str string) uint {
	intNum, _ := strconv.Atoi(str)
	return uint(intNum)
}
